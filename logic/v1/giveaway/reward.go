package giveaway

import (
	"fmt"

	"github.com/opensaucerers/giveawaybot/config"
	"github.com/opensaucerers/giveawaybot/repository/v1/giveaway"
	"github.com/opensaucerers/giveawaybot/repository/v1/user"
	"github.com/opensaucerers/giveawaybot/service"
	"github.com/opensaucerers/giveawaybot/typing"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Reward sends a direct message to the winners of the giveaway
func Reward(id string) (*giveaway.Giveaway, error) {

	// build user
	user := user.User{
		Twitter: typing.TwitterUser{
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

	// if replies are already fetched, return
	if len(giveaway.Replies) > 0 {
		return giveaway, nil
	}

	replies, err := service.RetriveReplies(config.Env.TwitterBearerToken, user.Twitter.ID, giveaway.TweetID)
	if err != nil {
		return nil, err
	}

	giveaway.Replies = replies
	giveaway.TotalReplies = len(replies)

	giveaway.Save()

	return giveaway, nil
}
