package service

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"

	"github.com/opensaucerers/giveawaybot/config"
	"github.com/opensaucerers/giveawaybot/typing"
	"github.com/samperfect/goaxios"
)

// GetAccessToken gets the access token from twitter
func GetTwitterAccessToken(code string) ([]byte, error) {
	// send http request to twitter to get access token
	r := goaxios.GoAxios{
		Url:     "https://api.twitter.com/2/oauth2/token",
		Method:  "POST",
		Headers: map[string]string{
			// needs to be empty to prevent goaxios from setting content-type to application/json
		},
		Query: map[string]interface{}{
			"code":          code,
			"grant_type":    "authorization_code",
			"client_id":     config.Env.TwitterClientID,
			"redirect_uri":  config.Env.TwitterRedirectURL,
			"code_verifier": "challenge",
		},
	}

	// send request
	_, b, _, err := r.RunRest()
	if err != nil {
		return nil, err
	}

	return b, nil
}

// GetAuthenticatedUser gets the authenticated user from twitter
func GetAuthenticatedTwitterUser(accessToken string) ([]byte, error) {
	// get the twitter user
	u := goaxios.GoAxios{
		Url:         "https://api.twitter.com/2/users/me",
		BearerToken: accessToken,
		Query: map[string]interface{}{
			"user.fields": "public_metrics",
		},
	}

	// send request
	_, b, _, err := u.RunRest()
	if err != nil {
		return nil, err
	}

	return b, nil
}

// RefreshTwitterAccessToken refreshes the access token from twitter
func RefreshTwitterAccessToken(refreshToken string) ([]byte, error) {
	// refresh the access token
	r := goaxios.GoAxios{
		Url:     "https://api.twitter.com/2/oauth2/token",
		Method:  "POST",
		Headers: map[string]string{
			// needs to be empty to prevent goaxios from setting content-type to application/json
		},
		Query: map[string]interface{}{
			"grant_type":    "refresh_token",
			"refresh_token": refreshToken,
			"client_id":     config.Env.TwitterClientID,
		},
	}

	// send request
	_, b, _, err := r.RunRest()
	if err != nil {
		return nil, err
	}

	return b, nil
}

// Tweet tweets a message to twitter
func Tweet(accessToken, message, replyTo string) ([]byte, error) {

	body := map[string]interface{}{"text": message}
	if replyTo != "" {
		body["reply"] = map[string]interface{}{
			"in_reply_to_tweet_id": replyTo,
		}
	}

	// tweet the message
	r := goaxios.GoAxios{
		Url:         "https://api.twitter.com/2/tweets",
		Method:      "POST",
		BearerToken: accessToken,
		Body:        body,
	}

	// send request
	_, b, _, err := r.RunRest()
	if err != nil {
		return nil, err
	}

	return b, nil
}

// ListTweets lists tweets from twitter
func ListTweets(accessToken, userID, limit, cursor string) ([]byte, error) {
	// list tweets
	query := map[string]interface{}{
		"tweet.fields": "created_at",
	}
	if limit != "" {
		query["max_results"] = limit
	}
	if cursor != "" {
		query["pagination_token"] = cursor
	}

	r := goaxios.GoAxios{
		Url:         "https://api.twitter.com/2/users/" + userID + "/tweets",
		BearerToken: accessToken,
		Query:       query,
	}

	// send request
	_, b, _, err := r.RunRest()
	if err != nil {
		return nil, err
	}

	return b, nil
}

// LikeTweet likes a tweet
func LikeTweet(accessToken, userID, tweetID string) ([]byte, error) {
	// like tweet
	r := goaxios.GoAxios{
		Url:         "https://api.twitter.com/2/users/" + userID + "/likes",
		Method:      "POST",
		BearerToken: accessToken,
		Body: map[string]interface{}{
			"tweet_id": tweetID,
		},
	}

	// send request
	_, b, _, err := r.RunRest()
	if err != nil {
		return nil, err
	}

	return b, nil
}

// Retweet retweets a tweet
func Retweet(accessToken, userID, tweetID string) ([]byte, error) {
	// retweet
	r := goaxios.GoAxios{
		Url:         "https://api.twitter.com/2/users/" + userID + "/retweets",
		Method:      "POST",
		BearerToken: accessToken,
		Body: map[string]interface{}{
			"tweet_id": tweetID,
		},
	}

	// send request
	_, b, _, err := r.RunRest()
	if err != nil {
		return nil, err
	}

	return b, nil
}

// ListLikes lists likes of a specific tweet
func ListLikes(accessToken, tweetID, limit, cursor string) ([]byte, error) {
	// list likes
	query := map[string]interface{}{}
	if limit != "" {
		query["max_results"] = limit
	}
	if cursor != "" {
		query["pagination_token"] = cursor
	}

	r := goaxios.GoAxios{
		Url:         "https://api.twitter.com/2/tweets/" + tweetID + "/liking_users",
		BearerToken: accessToken,
		Query:       query,
	}

	// send request
	_, b, _, err := r.RunRest()
	if err != nil {
		return nil, err
	}

	return b, nil
}

// IsLiking verifies if a user is liking a tweet
func IsLiking(accessToken, userID, tweetID, cursor string) (bool, error) {

	// time.Sleep(20 * time.Second)

	// list likes
	b, err := ListLikes(accessToken, tweetID, "100", cursor)
	if err != nil {
		return true, err
	}

	// parse response
	var response typing.TwitterListResponse

	if err := json.Unmarshal(b, &response); err != nil {
		return true, err
	}

	var null *int = nil
	var zero *int = new(int)
	if !reflect.DeepEqual(response.Meta.ResultCount, null) && reflect.DeepEqual(response.Meta.ResultCount, zero) {
		return false, nil
	}

	if len(response.Data) == 0 {
		var response typing.TwitterTweetError
		if err := json.Unmarshal(b, &response); err != nil {
			return true, err
		}

		if len(response.Errors) > 0 {

			if response.Errors[0].Title == config.ErrTwitterNotFound {
				return true, nil
			}

			if response.Errors[0].Title != "" {
				return true, errors.New(response.Errors[0].Title)
			}
		} else if response.Title != "" {
			return true, errors.New(response.Title)
		}

		return true, errors.New("unknown error")
	}

	for _, user := range response.Data {
		if user.ID == userID {
			return true, nil
		}
	}

	if response.Meta.NextToken != "" {
		return IsLiking(accessToken, userID, tweetID, response.Meta.NextToken)
	}

	return false, nil
}

// ListRetweets lists retweets of a specific tweet
func ListRetweets(accessToken, tweetID, limit, cursor string) ([]byte, error) {
	// list likes
	query := map[string]interface{}{}
	if limit != "" {
		query["max_results"] = limit
	}
	if cursor != "" {
		query["pagination_token"] = cursor
	}

	r := goaxios.GoAxios{
		Url:         "https://api.twitter.com/2/tweets/" + tweetID + "/retweeted_by",
		BearerToken: accessToken,
		Query:       query,
	}

	// send request
	_, b, _, err := r.RunRest()
	if err != nil {
		return nil, err
	}

	return b, nil
}

// IsRetweeting verifies if a user is retweeting a tweet
func IsRetweeting(accessToken, userID, tweetID, cursor string) (bool, error) {

	// time.Sleep(20 * time.Second)

	// list likes
	b, err := ListRetweets(accessToken, tweetID, "100", cursor)
	if err != nil {
		return true, err
	}

	// parse response
	var response typing.TwitterListResponse

	if err := json.Unmarshal(b, &response); err != nil {
		return true, err
	}

	var null *int = nil
	var zero *int = new(int)
	if !reflect.DeepEqual(response.Meta.ResultCount, null) && reflect.DeepEqual(response.Meta.ResultCount, zero) {
		return false, nil
	}

	if len(response.Data) == 0 {
		var response typing.TwitterTweetError
		if err := json.Unmarshal(b, &response); err != nil {
			return true, err
		}

		if len(response.Errors) > 0 {

			if response.Errors[0].Title == config.ErrTwitterNotFound {
				return true, nil
			}

			if response.Errors[0].Title != "" {
				return true, errors.New(response.Errors[0].Title)
			}
		} else if response.Title != "" {
			return true, errors.New(response.Title)
		}

		return true, errors.New("unknown error")
	}

	for _, user := range response.Data {
		if user.ID == userID {
			return true, nil
		}
	}

	if response.Meta.NextToken != "" {
		return IsRetweeting(accessToken, userID, tweetID, response.Meta.NextToken)
	}

	return false, nil
}

// GetTweet gets a tweet
func GetTweet(accessToken, tweetID string) ([]byte, error) {
	// get tweet
	r := goaxios.GoAxios{
		Url:         "https://api.twitter.com/2/tweets",
		BearerToken: accessToken,
		Query: map[string]interface{}{
			"ids":        tweetID,
			"expansions": "author_id",
		},
		Method: "GET",
	}

	// send request
	_, b, _, err := r.RunRest()
	if err != nil {
		return nil, err
	}

	return b, nil
}

// IsReplying verifies if a user is replying to a tweet
func IsReplying(accessToken, userID, replyID string) (bool, error) {

	if replyID == "" {
		return false, nil
	}

	// time.Sleep(20 * time.Second)

	// list likes
	b, err := GetTweet(accessToken, replyID)
	if err != nil {
		return true, err
	}

	// parse response
	var response typing.TwitterListResponse

	if err := json.Unmarshal(b, &response); err != nil {
		return true, err
	}

	if len(response.Data) == 0 {
		var response typing.TwitterTweetError
		if err := json.Unmarshal(b, &response); err != nil {
			return true, err
		}

		if len(response.Errors) > 0 {

			if response.Errors[0].Title == config.ErrTwitterNotFound {
				return false, nil
			}

			if response.Errors[0].Title != "" {
				return true, errors.New(response.Errors[0].Title)
			}
		} else if response.Title != "" {
			return true, errors.New(response.Title)
		}
	}

	for _, tweet := range response.Data {
		if tweet.ID == replyID && tweet.AuthorID == userID {
			return true, nil
		}
	}

	return false, nil
}

// IsAlive verifies if a tweet is alive
func IsAlive(accessToken, replyID string) (bool, error) {

	// time.Sleep(20 * time.Second)

	// list likes
	b, err := GetTweet(accessToken, replyID)
	if err != nil {
		return true, err
	}

	// parse response
	var response typing.TwitterListResponse

	if err := json.Unmarshal(b, &response); err != nil {
		return true, err
	}

	if len(response.Data) == 0 {
		var response typing.TwitterTweetError
		if err := json.Unmarshal(b, &response); err != nil {
			return true, err
		}

		if len(response.Errors) > 0 {

			if strings.EqualFold(response.Errors[0].Title, config.ErrTwitterNotFound) {
				return false, nil
			}

			if response.Errors[0].Title != "" {
				return true, errors.New(response.Errors[0].Title)
			}

			return true, errors.New("unknown error")
		} else if response.Title != "" {
			return true, errors.New(response.Title)
		}
	}

	return true, nil
}

// GetTweetEmbed gets a tweet embed
func GetTweetEmbed(username, tweetID string) ([]byte, error) {
	/* https://publish.twitter.com/oembed?url=https://twitter.com/UnderworldsNFT/status/1638179205081826306?s=20&partner=&hide_thread=false */

	// get embed
	r := goaxios.GoAxios{
		Url: "https://publish.twitter.com/oembed",
		Query: map[string]interface{}{
			"url":         "https://twitter.com/" + username + "/status/" + tweetID,
			"hide_thread": "false",
			"partner":     "",
			"s":           "20",
		},
		Method:  "GET",
		Headers: map[string]string{},
	}

	// send request
	_, b, _, err := r.RunRest()
	if err != nil {
		return nil, err
	}

	return b, nil
}
