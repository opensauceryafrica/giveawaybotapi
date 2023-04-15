package giveaway

import (
	"errors"
	"fmt"
	"strings"

	"github.com/opensaucerers/giveawaybot/config"
	"github.com/opensaucerers/giveawaybot/repository/v1/giveaway"
	"github.com/opensaucerers/giveawaybot/repository/v1/user"
	"github.com/opensaucerers/giveawaybot/service"
	"github.com/opensaucerers/giveawaybot/typing"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Replies(id string) (*giveaway.Giveaway, error) {

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

func Report(id string) (*giveaway.Report, error) {
	// build giveaway
	giveawayID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid giveaway id")
	}

	ga := giveaway.Giveaway{
		ID: giveawayID,
	}

	// check if giveaway exists
	if err := ga.Find(); err != nil {
		return nil, err
	}

	if ga.ID.String() == primitive.NilObjectID.String() {
		return nil, fmt.Errorf("giveaway not found")
	}

	// if giveaway is not finished, return error
	if !ga.Completed {
		return nil, fmt.Errorf("unable to generate report for an active giveaway")
	}

	// generate report
	report := giveaway.Report{
		TotalReplies:     ga.TotalReplies,
		DuplicateReplies: make(map[string]string),
		ValidReplies:     make(map[string]string),
		ValidRepliesList: make([]string, 0),
	}

	// find duplicate replies
	responses := make(map[string]string)
	for _, reply := range ga.Replies {
		// check if reply is a duplicate
		if val, ok := responses[reply.FText+"g1v3@w@60+"+reply.Username]; ok {

			// duplicate could be more than 2
			if _, ok := report.DuplicateReplies[reply.Username]; !ok {
				report.TotalDuplicates++
				report.DuplicateReplies[reply.Username] = fmt.Sprintf(`Tweeted the username '%s' more than once.  One tweet at https://twitter.com/%s/status/%s   Another tweet at https://twitter.com/%s/status/%s`, reply.Username, reply.Username, val, reply.Username, reply.TweetID)
			}

			// remove duplicate reply
			delete(responses, reply.FText+"g1v3@w@60+"+reply.Username)

		} else {
			// because we are deleting duplicate replies above, we have to check if the reply is a duplicate in case there are more than 2
			if _, ok := report.DuplicateReplies[reply.Username]; ok {
				report.DuplicateReplies[reply.Username] = fmt.Sprintf(`%s  Another tweet at https://twitter.com/%s/status/%s`, report.DuplicateReplies[reply.Username], reply.Username, reply.TweetID)
			} else {
				// add non duplicate reply
				responses[reply.FText+"g1v3@w@60+"+reply.Username] = reply.TweetID
			}
		}
	}

	// at the end of the loop above, we have a map of non duplicate replies
	// we can generate report for non duplicate replies
	for key, reply := range responses {
		report.TotalValidReplies++
		report.ValidReplies[strings.Split(key, "g1v3@w@60+")[1]] = fmt.Sprintf(`Tweeted the username '%s' once.  Tweet at https://twitter.com/%s/status/%s`, strings.Split(key, "g1v3@w@60+")[0], strings.Split(key, "g1v3@w@60+")[1], reply)
		report.ValidRepliesList = append(report.ValidRepliesList, strings.Split(key, "g1v3@w@60+")[0])
	}

	return &report, nil
}
