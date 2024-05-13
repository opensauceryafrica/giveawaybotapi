package service

import (
	"encoding/json"
	"errors"
	"log"
	"reflect"
	"strings"

	"github.com/opensaucerers/giveawaybot/config"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/opensaucerer/goaxios"
	"github.com/opensaucerers/giveawaybot/typing"
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
		Query: map[string]string{
			"code":          code,
			"grant_type":    "authorization_code",
			"client_id":     config.Env.TwitterClientID,
			"redirect_uri":  config.Env.TwitterRedirectURL,
			"code_verifier": "challenge",
		},
	}

	// send request
	resp := r.RunRest()
	if resp.Error != nil {
		return nil, resp.Error
	}

	return resp.Bytes, nil
}

// GetAuthenticatedUser gets the authenticated user from twitter
func GetAuthenticatedTwitter(accessToken string) ([]byte, error) {
	// get the twitter user
	u := goaxios.GoAxios{
		Url:         "https://api.twitter.com/2/users/me",
		BearerToken: accessToken,
		Query: map[string]string{
			"user.fields": "public_metrics",
		},
	}

	// send request
	resp := u.RunRest()
	if resp.Error != nil {
		return nil, resp.Error
	}

	return resp.Bytes, nil
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
		Query: map[string]string{
			"grant_type":    "refresh_token",
			"refresh_token": refreshToken,
			"client_id":     config.Env.TwitterClientID,
		},
	}

	// send request
	resp := r.RunRest()
	if resp.Error != nil {
		return nil, resp.Error
	}

	return resp.Bytes, nil
}

// Tweet tweets a message to twitter
func Tweet(accessToken, message, replyTo, quote string) ([]byte, error) {

	body := map[string]interface{}{"text": message}
	if replyTo != "" {
		body["reply"] = map[string]interface{}{
			"in_reply_to_tweet_id": replyTo,
		}
	}
	if quote != "" {
		body["quote_tweet_id"] = quote
	}

	// tweet the message
	r := goaxios.GoAxios{
		Url:         "https://api.twitter.com/2/tweets",
		Method:      "POST",
		BearerToken: accessToken,
		Body:        body,
	}

	// send request
	resp := r.RunRest()
	if resp.Error != nil {
		return nil, resp.Error
	}

	return resp.Bytes, nil
}

// ListTweets lists tweets from twitter
func ListTweets(accessToken, userID, limit, cursor string) ([]byte, error) {
	// list tweets
	query := map[string]string{
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
	resp := r.RunRest()
	if resp.Error != nil {
		return nil, resp.Error
	}

	return resp.Bytes, nil
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
	resp := r.RunRest()
	if resp.Error != nil {
		return nil, resp.Error
	}

	return resp.Bytes, nil
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
	resp := r.RunRest()
	if resp.Error != nil {
		return nil, resp.Error
	}

	return resp.Bytes, nil
}

// ListLikes lists likes of a specific tweet
func ListLikes(accessToken, tweetID, limit, cursor string) ([]byte, error) {
	// list likes
	query := map[string]string{}
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
	resp := r.RunRest()
	if resp.Error != nil {
		return nil, resp.Error
	}

	return resp.Bytes, nil
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
	query := map[string]string{}
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
	resp := r.RunRest()
	if resp.Error != nil {
		return nil, resp.Error
	}

	return resp.Bytes, nil
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
		Query: map[string]string{
			"ids":        tweetID,
			"expansions": "author_id",
		},
		Method: "GET",
	}

	// send request
	resp := r.RunRest()
	if resp.Error != nil {
		return nil, resp.Error
	}

	return resp.Bytes, nil
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
		Query: map[string]string{
			"url":         "https://twitter.com/" + username + "/status/" + tweetID,
			"hide_thread": "false",
			"partner":     "",
			"s":           "20",
		},
		Method:  "GET",
		Headers: map[string]string{},
	}

	// send request
	resp := r.RunRest()
	if resp.Error != nil {
		return nil, resp.Error
	}

	return resp.Bytes, nil
}

// DeleteTweet deletes a tweet
func DeleteTweet(accessToken, tweetID string) (bool, error) {
	// delete tweet
	r := goaxios.GoAxios{
		Url:         "https://api.twitter.com/2/tweets/" + tweetID,
		BearerToken: accessToken,
		Method:      "DELETE",
	}

	// send request
	resp := r.RunRest()
	if resp.Error != nil {
		return false, resp.Error
	}

	// parse response
	var response typing.TwitterDeleteTweetResponse

	if err := json.Unmarshal(resp.Bytes, &response); err != nil {
		return false, err
	}

	if !response.Data.Deteted {
		var response typing.TwitterTweetError
		if err := json.Unmarshal(resp.Bytes, &response); err != nil {
			return true, err
		}

		if len(response.Errors) > 0 {

			if strings.EqualFold(response.Errors[0].Title, config.ErrTwitterNotFound) {
				return true, nil
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

// Mentions gets all mentions to a user
func Mentions(accessToken, userID, limit, cursor string) ([]byte, error) {
	// list replies
	query := map[string]string{
		"tweet.fields": "conversation_id,referenced_tweets",
		"expansions":   "author_id",
		"user.fields":  "username",
	}
	if limit != "" {
		query["max_results"] = limit
	}
	if cursor != "" {
		query["pagination_token"] = cursor
	}
	r := goaxios.GoAxios{
		Url:         "https://api.twitter.com/2/users/" + userID + "/mentions",
		BearerToken: accessToken,
		Query:       query,
		Method:      "GET",
	}

	// send request
	// send request
	resp := r.RunRest()
	if resp.Error != nil {
		return nil, resp.Error
	}

	return resp.Bytes, nil
}

// RetriveReplies gets all replies to a tweet
func RetriveReplies(accessToken, userID, tweetID string, giveawayid primitive.ObjectID) ([]typing.Reply, error) {

	var (
		replies = make([]typing.Reply, 0)
		cursor  = ""
		limit   = "100"
		err     error
		b       []byte
		i       = 0
	)

	for {

		log.Printf("In 5 seconds getting page %d with cursor %s", i, cursor)

		// time.Sleep(5 * time.Second)

		// list likes
		b, err = Mentions(accessToken, userID, limit, cursor)
		if err != nil {
			return replies, err
		}

		i++

		// parse response
		var response typing.TwitterListResponse

		if err := json.Unmarshal(b, &response); err != nil {
			return replies, err
		}

		var null *int = nil
		var zero *int = new(int)
		if !reflect.DeepEqual(response.Meta.ResultCount, null) && reflect.DeepEqual(response.Meta.ResultCount, zero) {
			return replies, nil
		}

		if len(response.Data) == 0 {
			var response typing.TwitterTweetError
			if err := json.Unmarshal(b, &response); err != nil {
				break
			}

			if len(response.Errors) > 0 {

				if response.Errors[0].Title == config.ErrTwitterNotFound {
					break
				}

				if response.Errors[0].Title != "" {
					err = errors.New(response.Errors[0].Title)
					break
				}
			} else if response.Title != "" {
				err = errors.New(response.Title)
				break
			}

			break
		}

		log.Printf("Found %d replies on page %d", len(response.Data), i)

		for _, tweet := range response.Data {

			if tweet.ConversationID == tweetID {

				username := ""
			Username:
				for _, user := range response.Includes.Users {
					if user.ID == tweet.AuthorID {
						username = user.Username
						break Username
					}
				}

				// each tweet text is expected to contain @username. Ideally, each should contain just two @s, the first one being the username of the person being replied to and the second one being the username the person replying tagged. However, there are cases where the tweet text contains more than two @s. To efficiently handle this, we get the last @ in the tweet text and use that to get the username the person replying tagged.

				at := strings.LastIndex(tweet.Text, "@")
				if at == -1 {
					continue
				}

				// get the username
				ftext := tweet.Text[at:]

				replies = append(replies, typing.Reply{
					ID:       tweet.AuthorID,
					Text:     tweet.Text,
					Username: username,
					TweetID:  tweet.ID,
					FText:    ftext,
					Giveaway: giveawayid,
				})
			}
		}

		if response.Meta.NextToken != "" {
			cursor = response.Meta.NextToken
		} else {
			cursor = ""
		}

		if cursor == "" {
			log.Println("Stopping iteration...no more data")
			break
		}

		log.Printf("Next cursor: %s", cursor)

	}

	return replies, err
}

// Message sends a message to a user
func Message(accessToken, userID, text string) ([]byte, error) {

	// send message
	r := goaxios.GoAxios{
		Url:         "https://api.twitter.com/2/dm_conversations/with/" + userID + "/messages",
		BearerToken: accessToken,
		Body: map[string]interface{}{
			"text": text,
		},
		Method: "POST",
	}

	// send request
	resp := r.RunRest()
	if resp.Error != nil {
		return nil, resp.Error
	}

	return resp.Bytes, nil
}
