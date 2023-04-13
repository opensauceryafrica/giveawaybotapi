package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/opensaucerers/giveawaybot/config"
	"github.com/opensaucerers/giveawaybot/repository/v1/user"
	"github.com/opensaucerers/giveawaybot/service"
	"github.com/opensaucerers/giveawaybot/typing"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TwitterAuth() string {

	return fmt.Sprintf(
		"https://twitter.com/i/oauth2/authorize?response_type=code&client_id=%s&scope=%s&redirect_uri=%s&state=state&code_challenge=challenge&code_challenge_method=plain",
		config.Env.TwitterClientID,
		config.TwitterScope,
		config.Env.TwitterRedirectURL,
	)
}

// TwitterSignon registers a new user with twitter
func TwitterSignon(code string) (*user.User, error) {

	// send http request to twitter to get access token
	b, err := service.GetTwitterAccessToken(code)
	if err != nil {
		return nil, errors.New(config.ErrTwitterUnauthorized)
	}

	// parse response
	var authResponse typing.TwitterAuthResponse

	if err := json.Unmarshal(b, &authResponse); err != nil {
		return nil, errors.New(config.ErrTwitterUnauthorized)
	}

	// if no access token, return error
	if authResponse.AccessToken == "" {
		var twitterError typing.TwitterAuthError

		if err := json.Unmarshal(b, &twitterError); err != nil {
			return nil, errors.New(config.ErrTwitterUnauthorized)
		}

		return nil, errors.New(config.ErrTwitterUnauthorized)
	}

	// get the twitter user
	b, err = service.GetAuthenticatedTwitterUser(authResponse.AccessToken)
	if err != nil {
		return nil, errors.New(config.ErrTwitterUnauthorized)
	}

	// parse response
	var twitterUser typing.TwitterUserResponse

	if err := json.Unmarshal(b, &twitterUser); err != nil {
		return nil, errors.New(config.ErrTwitterUnauthorized)
	}

	// build user
	user := user.User{
		Username:       twitterUser.Data.Username,
		UseTwitterAuth: true,
		Twitter: typing.TwitterUser{
			AccessToken:  authResponse.AccessToken,
			RefreshToken: authResponse.RefreshToken,
			TokenType:    authResponse.TokenType,
			RetrievedAt:  time.Now().UTC(),
			UpdatedAt:    time.Now().UTC(),
			ID:           twitterUser.Data.ID,
			Username:     twitterUser.Data.Username,
		},
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	// check if twitter user exists
	if err := user.FindSocial(typing.Social{Twitter: true}); err != nil {
		return nil, err
	}

	// if twitter user exists, return error
	if user.ID.String() != primitive.NilObjectID.String() {
		// if user exists, prompt to login
		return TwitterLogin(user.Twitter.ID, authResponse)
	}

	// if user doesn't exist, create user
	user.Avatar = config.DefaultAvatar // set default avatar
	if err := user.Create(); err != nil {
		return nil, err
	}

	// make safe
	user.Safe()

	return &user, nil
}

// TwitterLogin logs in a user with twitter
func TwitterLogin(id string, authResponse typing.TwitterAuthResponse) (*user.User, error) {

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

	// return user if exists
	if user.ID.String() == primitive.NilObjectID.String() {
		return nil, fmt.Errorf("%s user found for Twitter account. Please register first", "No")
	}

	// twitter sign in
	if !user.UseTwitterAuth {
		return nil, fmt.Errorf("%s authentication disabled for this account", "Twitter")
	}

	// update user with new access token
	user.Twitter.AccessToken = authResponse.AccessToken
	user.Twitter.RefreshToken = authResponse.RefreshToken
	user.Save()

	// make safe
	user.Safe()

	return &user, nil
}
