package typing

type Env struct {
	//Port for the server to listen on
	Port string `env:"PORT"`
	//MongoDBURI is the monogoDB connection string for the app
	MongoDBURI string `env:"MONGO_DB_URI"`
	//MongoDBName is the MongoDB database name for the app
	MongoDBName string `env:"MONGO_DB_NAME"`
	//PostgreSQLURI is the PostgreSQL database uri string for App
	// PostgreSQLURI string `env:"POSTGRE_SQL_URI"`

	// TwitterClientID is the twitter client id
	TwitterClientID string `env:"TWITTER_CLIENT_ID"`
	// TwitterClientSecret is the twitter client secret
	TwitterClientSecret string `env:"TWITTER_CLIENT_SECRET"`
	// TwitterAPIKey is the twitter api key
	TwitterAPIKey string `env:"TWITTER_API_KEY"`
	// TwitterAPISecretKey is the twitter api secret key
	TwitterAPISecretKey string `env:"TWITTER_API_SECRET_KEY"`
	// TwitterAccessToken is the twitter access token
	TwitterAccessToken string `env:"TWITTER_ACCESS_TOKEN"`
	// TwitterBearerToken is the twitter bearer token
	TwitterBearerToken string `env:"TWITTER_BEARER_TOKEN"`
	// TwitterAccessTokenSecret is the twitter access token secret
	TwitterAccessTokenSecret string `env:"TWITTER_ACCESS_TOKEN_SECRET"`
	// TwitterRedirectURL is the twitter redirect url
	TwitterRedirectURL string `env:"TWITTER_REDIRECT_URL"`

	// JWTSecret is the jwt secret
	JWTSecret string `env:"JWT_SECRET"`
	// AppName is the app name
	AppName string `env:"APP_NAME"`
}
