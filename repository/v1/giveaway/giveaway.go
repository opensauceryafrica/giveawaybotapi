package giveaway

import (
	"context"
	"encoding/json"
	"time"

	"github.com/opensaucerers/giveawaybot/config"
	"github.com/opensaucerers/giveawaybot/database"
	"github.com/opensaucerers/giveawaybot/repository/v1/user"
	"github.com/opensaucerers/giveawaybot/service"
	"github.com/opensaucerers/giveawaybot/typing"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Create inserts a new giveaway into the database
func (r *Giveaway) Create() error {
	// insert new giveaway
	result, err := database.MongoDB.Collection(config.GiveawayCollection).InsertOne(context.Background(), r)
	if err != nil {
		return err
	}
	r.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// Save updates the giveaway in database
func (g *Giveaway) Save() error {
	// update giveaway
	g.UpdatedAt = time.Now().UTC()
	_, err := database.MongoDB.Collection(config.GiveawayCollection).UpdateOne(context.Background(), bson.M{"_id": g.ID}, bson.M{"$set": g})
	if err != nil {
		return err
	}
	return nil
}

// Find finds a giveaway by id
func (g *Giveaway) Find() error {
	if err := database.MongoDB.Collection(config.GiveawayCollection).FindOne(context.Background(), bson.M{"_id": g.ID}).Decode(&g); err != nil {
		if err != mongo.ErrNoDocuments {
			return err
		}
		g.ID = primitive.NilObjectID
	}
	return nil
}

// FindUserGiveaways finds all giveaways for a user
func (g *Giveaways) FindUserGiveaways(user user.User, limit, offset, direction int64) error {
	//direction takes either -1(desc) and 1(asc)
	findOptions := options.Find().SetLimit(limit).SetSkip(offset).SetSort(bson.M{"created_at": direction})
	if cursor, err := database.MongoDB.Collection(config.GiveawayCollection).Find(context.Background(), bson.M{"author._id": user.ID}, findOptions); err != nil {
		return err
	} else {
		if err := cursor.All(context.Background(), g); err != nil {
			return err
		}
	}
	return nil
}

// Complete marks a giveaway as completed
func (g *Giveaway) Complete() error {
	// update giveaway
	_, err := database.MongoDB.Collection(config.GiveawayCollection).UpdateOne(context.Background(), bson.M{"_id": g.ID}, bson.M{"$set": bson.M{"active": false, "completed": true, "completed_at": time.Now().UTC(), "updated_at": time.Now().UTC()}})
	if err != nil {
		return err
	}
	return nil
}

// IsRunning checks if there exists at least one active giveaway
// for a user
func IsRunning(owner user.User) (bool, error) {
	// find active giveaways
	if count, err := database.MongoDB.Collection(config.GiveawayCollection).CountDocuments(context.Background(), bson.M{"active": true, "completed": false, "author._id": owner.ID}); err != nil {
		if err != mongo.ErrNoDocuments {
			return false, nil // no active giveaways
		}
		return false, err
	} else {
		return count > 0, nil
	}
}

// Running finds the active giveaway for a user
func Running(owner user.User) (*Giveaway, error) {
	// find active giveaways
	giveaway := Giveaway{}
	if err := database.MongoDB.Collection(config.GiveawayCollection).FindOne(context.Background(), bson.M{"active": true, "completed": false, "author._id": owner.ID}).Decode(&giveaway); err != nil {
		if err != mongo.ErrNoDocuments {
			return nil, err
		}
	}
	return &giveaway, nil
}

// Delete deletes a giveaway
func (g *Giveaway) Delete() error {
	// delete giveaway
	_, err := database.MongoDB.Collection(config.GiveawayCollection).DeleteOne(context.Background(), bson.M{"_id": g.ID})
	if err != nil {
		return err
	}
	return nil
}

// EmbedTweet embeds a tweet in the giveaway
func (g *Giveaway) EmbedTweet() error {
	b, err := service.GetTweetEmbed(g.Author.Username, g.TweetID)
	if err != nil {
		return err
	}
	// unmarshal response
	var e typing.TwitterEmbedResponse
	if err := json.Unmarshal(b, &e); err != nil {
		return err
	}
	// update giveaway
	g.TwitterHTML = e.HTML
	g.TwitterURL = e.URL
	return nil
}
