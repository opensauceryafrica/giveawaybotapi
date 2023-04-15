package giveaway

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/opensaucerers/giveawaybot/config"
	"github.com/opensaucerers/giveawaybot/repository/v1/giveaway"
	"github.com/opensaucerers/giveawaybot/repository/v1/user"
	"github.com/opensaucerers/giveawaybot/service"
	"github.com/opensaucerers/giveawaybot/typing"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Start tweets about the giveaway
func Start(id string) (*giveaway.Giveaway, error) {

	/// build user
	user := user.User{
		Twitter: typing.Twitter{
			ID: id,
		},
	}

	// check if user exists with twitter
	if err := user.FindSocial(typing.Social{Twitter: true}); err != nil {
		return nil, err
	}

	// if user doesn't exist, return error
	if user.ID.String() == primitive.NilObjectID.String() {
		return nil, fmt.Errorf("user not found")
	}

	// only one active giveaway at a time per user
	isRunning, err := giveaway.IsRunning(user)
	if err != nil {
		return nil, err
	}

	if isRunning {
		return nil, fmt.Errorf("you already have a giveaway running")
	}

	// refresh user access token
	if err := user.RefreshTwitterAccessToken(); err != nil {
		return nil, err
	}

	// Tweet the tweet
	b, err := service.Tweet(user.Twitter.AccessToken, config.TwitterGiveawayComment, "", config.TwitterGiveawayTweet)
	if err != nil {
		return nil, err
	}

	// parse response
	var tweetResponse typing.TwitterTweetResponse

	if err := json.Unmarshal(b, &tweetResponse); err != nil {
		return nil, err
	}

	if tweetResponse.Data.ID == "" {
		var twitterError typing.TwitterTweetError

		if err := json.Unmarshal(b, &twitterError); err != nil {
			return nil, err
		}

		return nil, fmt.Errorf(twitterError.Errors[0].Message)
	}

	// create giveaway
	giveaway := giveaway.Giveaway{
		Author:    user,
		Tweet:     tweetResponse.Data.Text,
		TweetID:   tweetResponse.Data.ID,
		Replies:   []typing.Reply{},
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Active:    true,
	}
	giveaway.EmbedTweet()

	// save giveaway
	if err := giveaway.Create(); err != nil {
		return nil, err
	}

	return &giveaway, nil
}

// Disrupt deletes the giveaway tweet and the giveaway
func Disrupt(id string) (*giveaway.Giveaway, error) {

	// build user
	user := user.User{
		Twitter: typing.Twitter{
			ID: id,
		},
	}

	// check if user exists with twitter
	if err := user.FindSocial(typing.Social{Twitter: true}); err != nil {
		return nil, err
	}

	// if user doesn't exist, return error
	if user.ID.String() == primitive.NilObjectID.String() {
		return nil, fmt.Errorf("user not found")
	}

	// refresh user access token
	if err := user.RefreshTwitterAccessToken(); err != nil {
		return nil, err
	}

	// get giveaway
	giveaway, err := giveaway.Running(user)
	if err != nil {
		return nil, err
	}

	if giveaway.ID.String() == primitive.NilObjectID.String() {
		return nil, fmt.Errorf("no active giveaway found")
	}

	// delete tweet
	deleted, err := service.DeleteTweet(user.Twitter.AccessToken, giveaway.TweetID)
	if err != nil {
		return nil, err
	}

	if !deleted {
		return nil, fmt.Errorf("could not delete tweet: %s", err.Error())
	}

	// delete giveaway
	if err := giveaway.Delete(); err != nil {
		return nil, err
	}

	return giveaway, nil
}

// End ends the giveaway
func End(id string) (*giveaway.Giveaway, error) {

	// build user
	user := user.User{
		Twitter: typing.Twitter{
			ID: id,
		},
	}

	// check if user exists with twitter
	if err := user.FindSocial(typing.Social{Twitter: true}); err != nil {
		return nil, err
	}

	// if user doesn't exist, return error
	if user.ID.String() == primitive.NilObjectID.String() {
		return nil, fmt.Errorf("user not found")
	}

	// refresh user access token
	if err := user.RefreshTwitterAccessToken(); err != nil {
		return nil, err
	}

	// get giveaway
	giveaway, err := giveaway.Running(user)
	if err != nil {
		return nil, err
	}

	if giveaway.ID.String() == primitive.NilObjectID.String() {
		return nil, fmt.Errorf("no active giveaway found")
	}

	// don't end if replies are empty
	if len(giveaway.Replies) == 0 {
		return nil, fmt.Errorf("you need to trigger the bot to retrieve all replies before ending the giveaway")
	}

	// end giveaway
	if err := giveaway.Complete(); err != nil {
		return nil, err
	}

	return giveaway, nil
}
