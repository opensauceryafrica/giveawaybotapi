package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/opensaucerers/giveawaybot/config"
	"github.com/opensaucerers/giveawaybot/cron"
	"github.com/opensaucerers/giveawaybot/database"
	"github.com/opensaucerers/giveawaybot/typing"

	"github.com/opensaucerers/giveawaybot/middleware/v1"
	"github.com/opensaucerers/giveawaybot/version"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	ts "github.com/n0madic/twitter-scraper"
)

func createServer() (s *http.Server) {

	// create router
	r := mux.NewRouter()

	// wrap router into custom recover middleware
	rwr := middleware.Recover(r)

	// we should do more cross origin stuff here
	rc := middleware.CORS(rwr)

	// inject combined logger (apache & nginx style)
	logger := handlers.CombinedLoggingHandler(os.Stdout, rc)

	// register routes with versioning
	version.Version1Routes(r.StrictSlash(true))

	//connect to monogoDB and select database
	if err := database.NewMongoDBClient(config.Env.MongoDBURI, config.Env.MongoDBName); err != nil {
		log.Fatalf("Error connecting to MongoDBURI: %s\n", err)
	}

	//connect to PostgreSQL database
	// INFO: using mongodb instead
	// if err := database.NewPostgreSQLConnection(config.Env.PostgreSQLURI, config.MaxConnections); err != nil {
	// 	log.Fatalf("Error connecting to PostgreSQLURI: %s\n", err)
	// }

	s = &http.Server{
		Addr:           fmt.Sprintf(":%s", config.Env.Port),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		Handler:        logger,
	}

	go func() {
		go func() {
			// start cron jobs
			cron.S.StartBlocking()
		}()

		log.Printf("Starting at http://127.0.0.1%s", fmt.Sprintf(":%s", config.Env.Port))
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			cron.S.StopBlockingChan()
			log.Fatalf("error listening on port: %s\n", err)
		}
	}()

	return s
}

func main() {

	// load .env file
	env := config.MustGet("ENV_PATH", ".env") // might be wondering, get ENV_PATH when env is not loaded yet? the answer is that you can set env variables out of the .env file (e.g. in the terminal using export ENV_PATH=...)
	log.Printf("Loading %s file\n", env)
	if err := godotenv.Load(env); err != nil {
		if err := godotenv.Load(); err != nil {
			log.Printf("Error loading %s file\n", env)
		}
	}

	// verify env variables
	if err := config.VerifyEnvironment(typing.Env{}); err != nil {
		log.Fatalf("Error verifying environment variables: %s\n", err)
	}

	// append env variables to config.Env
	config.AppendEnvironment(config.Env)

	scraper := ts.New()

	err := scraper.Login("", "", "")
	if err != nil {
		log.Fatal(err)
	}

	cookies := scraper.GetCookies()
	// serialize to JSON
	js, _ := json.Marshal(cookies)
	// save to file
	f, _ := os.Create("cookies.json")
	f.Write(js)

	// scraper.WithReplies(true)

	// var cookies []*http.Cookie

	// if err := json.Unmarshal([]byte(os.Getenv("TWITTER_COOKIES")), &cookies); err != nil {
	// 	log.Printf(`[main.main] [json.Unmarshal([]byte(primer.ENV.TwitterCookies), &cookies)] %s`, err.Error())
	// 	os.Exit(1)
	// }

	// load cookies
	scraper.SetCookies(cookies)
	// check login status
	scraper.IsLoggedIn()

	// get tweet with replies
	tw := ""
	counter := 0
	for tweet := range scraper.SearchTweets(context.Background(), "@opensaucery @opensauceryf", 1800) {

		if tweet.Error != nil {
			log.Print(tweet.Error)
			continue
		}
		if tweet.InReplyToStatusID == tw {
			counter++
			log.Printf("Tweet %d", counter)
			log.Println("nil:", tweet.ID, tweet.Text, tweet.InReplyToStatusID)
			continue
		}
		// if tweet.InReplyToStatus.ID == tw {
		// log.Println(tw, tweet.InReplyToStatus.ID, tweet.Text)
		// }
	}

	return

	s := createServer()
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need to add it
	signal.Notify(config.ShutdownChan, syscall.SIGINT, syscall.SIGTERM)
	<-config.ShutdownChan
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), config.ShutdownTimeout*time.Second)
	defer func() {
		cancel()
	}()
	if err := s.Shutdown(ctx); err != nil {
		cron.S.StopBlockingChan()
		log.Fatal("Server forced to shut down...")
	}
	cron.S.StopBlockingChan()
	log.Println("Server exited!")
}
