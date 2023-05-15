package giveaway

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/opensaucerers/giveawaybot/config"
	"github.com/opensaucerers/giveawaybot/repository/v1/giveaway"
	"github.com/opensaucerers/giveawaybot/repository/v1/user"
	"github.com/opensaucerers/giveawaybot/service"
	"github.com/opensaucerers/giveawaybot/typing"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Winners registers the winners of the giveaway
func Winners(id string, winners []string) (*giveaway.Giveaway, error) {

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

	// if winners are already fetched, return
	if len(giveaway.Winners) > 0 {
		return giveaway, nil
	}

	// remove " prefix and " suffix
	for i, winner := range winners {
		winners[i] = strings.Trim(winner, "\"")
		// prefix @ if not present
		if !strings.HasPrefix(winners[i], "@") {
			winners[i] = "@" + winners[i]
		}
	}

	// get winners
	giveaway.Winners = winners
	giveaway.Save()

	return giveaway, nil
}

// Reward tweets the winners and sends a direct message to the winners of the giveaway
func Reward(id string) (*giveaway.Giveaway, error) {

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

	// can't double spend
	if giveaway.Rewarded {
		// send direct message to winners
		go giveaway.InboxForReward(user)
		return nil, fmt.Errorf("giveaway already rewarded. data claim message resent to winners")
	}

	// if no winners are fetched, return
	if len(giveaway.Winners) == 0 {
		return nil, fmt.Errorf("no winners found")
	}

	if err := giveaway.Replyies(); err != nil {
		return nil, err
	}

	// tweet report
	rb, err := service.Tweet(user.Twitter.AccessToken, fmt.Sprintf(config.TwitterGiveawayReport, "https://abbrefy.xyz/giveawayreport"), "", giveaway.TweetID)
	if err != nil {
		return nil, err
	}

	// parse response
	var rtweetResponse typing.TwitterTweetResponse

	if err := json.Unmarshal(rb, &rtweetResponse); err != nil {
		return nil, err
	}

	if rtweetResponse.Data.ID == "" {
		var twitterError typing.TwitterTweetError

		if err := json.Unmarshal(rb, &twitterError); err != nil {
			return nil, err
		}

		return nil, fmt.Errorf(twitterError.Errors[0].Message)
	}

	text := ""
	for _, winner := range giveaway.Winners {
		text += winner + "\n"
	}
	// remove last \n
	text = text[:len(text)-1]

	tweets := []string{}
	stop := false
	replyTo := rtweetResponse.Data.ID
	for {

		replacement := ""

		// check twitter character limit and remove from the last @ if it exceeds
		if len(text) > config.TwitterCharacterLimit {
			cutout := text[:config.TwitterCharacterLimit]
			at := strings.LastIndex(cutout, "@")
			if at != -1 {
				replacement = cutout[:at]
				text = text[at:]
			} else {
				replacement = text[:config.TwitterCharacterLimit]
				stop = true
			}
		} else {
			replacement = text
			stop = true
		}

		// Tweet winners
		wb, err := service.Tweet(user.Twitter.AccessToken, fmt.Sprintf(config.TwitterGiveawayWinners, replacement), replyTo, "")
		if err != nil {
			return nil, err
		}

		// parse response
		var wtweetResponse typing.TwitterTweetResponse

		if err := json.Unmarshal(wb, &wtweetResponse); err != nil {
			return nil, err
		}

		if wtweetResponse.Data.ID == "" {
			var twitterError typing.TwitterTweetError

			if err := json.Unmarshal(wb, &twitterError); err != nil {
				return nil, err
			}

			return nil, fmt.Errorf(twitterError.Errors[0].Message)
		}

		tweets = append(tweets, wtweetResponse.Data.ID)
		replyTo = wtweetResponse.Data.ID

		if stop {
			break
		}

	}

	giveaway.ReportTweet = rtweetResponse.Data.ID
	giveaway.WinnersTweet = tweets
	giveaway.Rewarded = true

	giveaway.Save()

	// send direct message to winners
	go giveaway.InboxForReward(user)

	return giveaway, nil
}
