package user

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/opensaucerers/giveawaybot/config"
	"github.com/opensaucerers/giveawaybot/database"
	"github.com/opensaucerers/giveawaybot/service"
	"github.com/opensaucerers/giveawaybot/typing"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

//check if a user has an Avatar
func (u *User) SetDefaultAvatar() error {
	if u.Avatar != "" {
		return nil
	}
	u.Avatar = config.DefaultAvatar
	return nil
}

// FindUser finds user by social id
func (u *User) FindSocial(q typing.Social) error {
	query := []bson.M{}
	if q.Twitter {
		query = append(query, bson.M{"twitter.id": u.Twitter.ID})
	}
	if err := database.MongoDB.Collection(config.UserCollection).FindOne(context.Background(), bson.M{"$or": query}).Decode(&u); err != nil {
		if err != mongo.ErrNoDocuments {
			return err
		}
	}
	return nil
}

// Find finds user by id
func (u *User) Find() error {
	if err := database.MongoDB.Collection(config.UserCollection).FindOne(context.Background(), bson.M{"_id": u.ID}).Decode(&u); err != nil {
		if err != mongo.ErrNoDocuments {
			return err
		}
	}
	return nil
}

// Create inserts a new user into the database
func (u *User) Create() error {
	// insert new user
	r, err := database.MongoDB.Collection(config.UserCollection).InsertOne(context.Background(), u)
	if err != nil {
		return err
	}
	u.ID = r.InsertedID.(primitive.ObjectID)
	return nil
}

// Safe removes sensitive information from user
func (u *User) Safe() {
	u.Twitter.AccessToken = ""
	u.Twitter.RefreshToken = ""
}

// Save updates user in database
func (u *User) Save() error {
	// update user
	_, err := database.MongoDB.Collection(config.UserCollection).UpdateOne(context.Background(), bson.M{"_id": u.ID}, bson.M{"$set": u})
	if err != nil {
		return err
	}
	return nil
}

// RefreshTwitterAccessToken refreshes a twitter access token
func (u *User) RefreshTwitterAccessToken() error {
	b, err := service.RefreshTwitterAccessToken(u.Twitter.RefreshToken)
	if err != nil {
		return err
	}

	// parse response
	var authResponse typing.TwitterAuthResponse

	if err := json.Unmarshal(b, &authResponse); err != nil {
		return err
	}

	// if no access token, return error
	if authResponse.AccessToken == "" || authResponse.RefreshToken == "" {
		var twitterError typing.TwitterAuthError

		if err := json.Unmarshal(b, &twitterError); err != nil {
			return err
		}

		return fmt.Errorf(config.ErrTwitterUnauthorized)
	}

	// update user with new access token
	u.Twitter.AccessToken = authResponse.AccessToken
	u.Twitter.RefreshToken = authResponse.RefreshToken
	u.Save()
	return nil
}

// FindByUsername finds user by username
func (u *User) FindByUsername() error {
	if err := database.MongoDB.Collection(config.UserCollection).FindOne(context.Background(), bson.M{"username": u.Username}).Decode(&u); err != nil {
		if err != mongo.ErrNoDocuments {
			return err
		}
	}
	return nil
}
