package giveaway

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/opensaucerers/giveawaybot/config"
	"github.com/opensaucerers/giveawaybot/database"
	"github.com/opensaucerers/giveawaybot/helper"
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

// FindByTweet finds a giveaway by tweet id
func (g *Giveaway) FindByTweet() error {
	if err := database.MongoDB.Collection(config.GiveawayCollection).FindOne(context.Background(), bson.M{"tweet_id": g.TweetID}).Decode(&g); err != nil {
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

// Reward marks a giveaway as rewarded
func (g *Giveaway) Reward() error {
	// update giveaway
	_, err := database.MongoDB.Collection(config.GiveawayCollection).UpdateOne(context.Background(), bson.M{"_id": g.ID}, bson.M{"$set": bson.M{"rewarded": true, "updated_at": time.Now().UTC()}})
	if err != nil {
		return err
	}
	return nil
}

// Close closes a giveaway by marking as completed and rewarding
func (g *Giveaway) Close() error {
	// update giveaway
	_, err := database.MongoDB.Collection(config.GiveawayCollection).UpdateOne(context.Background(), bson.M{"_id": g.ID}, bson.M{"$set": bson.M{"active": false, "completed": true, "completed_at": time.Now().UTC(), "rewarded": true, "updated_at": time.Now().UTC()}})
	if err != nil {
		return err
	}
	return nil
}

// IsRunning checks if there exists at least one active giveaway
// for a user
func IsRunning(owner user.User) (bool, error) {
	// find active giveaways
	if count, err := database.MongoDB.Collection(config.GiveawayCollection).CountDocuments(context.Background(), bson.M{"$or": []bson.M{
		{
			"active":     true,
			"completed":  false,
			"author._id": owner.ID,
		},
		{
			"active":     false,
			"completed":  true,
			"rewarded":   false,
			"author._id": owner.ID,
		},
	}}); err != nil {
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
	if err := database.MongoDB.Collection(config.GiveawayCollection).FindOne(context.Background(), bson.M{"$or": []bson.M{
		{
			"active":     true,
			"completed":  false,
			"author._id": owner.ID,
		},
		{
			"active":     false,
			"completed":  true,
			"rewarded":   false,
			"author._id": owner.ID,
		},
	}}).Decode(&giveaway); err != nil {
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

// Replies returns all replies for a giveaway
func (g *Giveaway) Replyies() error {
	// get replies
	if cursor, err := database.MongoDB.Collection(config.RepliesCollection).Find(context.Background(), bson.M{"giveaway": g.ID}); err != nil {
		return err
	} else {
		replies := []typing.Reply{}
		if err := cursor.All(context.Background(), &replies); err != nil {
			return err
		}
		g.Replies = replies
		return nil
	}
}

// ReplyInBatch saves replies in batch
func (g *Giveaway) ReplyInBatch() error {
	// save replies
	var replies []interface{}
	for _, reply := range g.Replies {
		replies = append(replies, reply)
	}
	if _, err := database.MongoDB.Collection(config.RepliesCollection).InsertMany(context.Background(), replies); err != nil {
		return err
	}
	return nil
}

// FindGiveawaysByMatch finds giveaways by match
func (g *Giveaways) FindGiveawaysByMatch(match bson.M, limit, offset int64) error {
	if cursor, err := database.MongoDB.Collection(config.GiveawayCollection).Find(context.Background(), match, &options.FindOptions{
		Limit: &limit,
		Skip:  &offset,
		// Sort:  bson.M{"created_at": -1},
	}); err != nil {
		return err
	} else {
		if err := cursor.All(context.Background(), g); err != nil {
			return err
		}
	}
	return nil
}

// FindRepliesByMatch finds giveaways by match
func (r *Replies) FindRepliesByMatch(match bson.M, limit, offset int64) error {
	if cursor, err := database.MongoDB.Collection(config.RepliesCollection).Find(context.Background(), match, &options.FindOptions{
		Limit: &limit,
		Skip:  &offset,
		// Sort:  bson.M{"created_at": -1},
	}); err != nil {
		return err
	} else {
		if err := cursor.All(context.Background(), r); err != nil {
			return err
		}
	}
	return nil
}

func YieldGiveawayByMatch(match bson.M, limit int64) (giveaways chan Giveaways) {
	// create a channel to send orders
	giveaways = make(chan Giveaways)
	done := make(chan bool, 1)
	var skip int64

	go func(skip, limit *int64) {
		// loop until done
	Loop:
		for {
			select {
			case <-done:
				close(giveaways) // close the channel
				break Loop       // break the loop
			default:
				// get completed orders
				g := Giveaways{}
				rs := g.FindGiveawaysByMatch(match, *limit, *skip)
				if rs != nil {
					done <- true
					break
				}
				if len(g) == 0 {
					done <- true
					break
				}
				*skip += int64(len(g))
				giveaways <- g
			}
		}
	}(&skip, &limit)
	return
}

func YieldRepliesByMatch(match bson.M, limit int64) (replies chan Replies) {
	// create a channel to send orders
	replies = make(chan Replies)
	done := make(chan bool, 1)
	var skip int64

	go func(skip, limit *int64) {
		// loop until done
	Loop:
		for {
			select {
			case <-done:
				close(replies) // close the channel
				break Loop     // break the loop
			default:
				// get completed orders
				r := Replies{}
				rs := r.FindRepliesByMatch(match, *limit, *skip)
				if rs != nil {
					done <- true
					break
				}
				if len(r) == 0 {
					done <- true
					break
				}
				*skip += int64(len(r))
				replies <- r
			}
		}
	}(&skip, &limit)
	return
}

// ClearReplies clears replies for a giveaway
func (g *Giveaway) ClearReplies() error {
	// delete replies
	_, err := database.MongoDB.Collection(config.GiveawayCollection).UpdateOne(context.Background(), bson.M{"_id": g.ID}, bson.M{"$set": bson.M{"replies": nil}})
	if err != nil {
		return err
	}
	return nil
}

// SaveReply saves the given giveaway reply
func (g *Giveaway) SaveReply(reply *typing.Reply) error {
	// save reply
	_, err := database.MongoDB.Collection(config.RepliesCollection).UpdateOne(context.Background(), bson.M{"_id": reply.MID}, bson.M{"$set": reply})
	if err != nil {
		return err
	}
	return nil
}

// InboxForReward sends a message to the giveaway winners
func (g *Giveaway) InboxForReward(user user.User) {
	for _, winner := range g.Winners {

		winner := strings.Trim(winner, "@")

		// get winner
		for _, r := range g.Replies {

			if strings.EqualFold(r.Username, winner) {
				if r.ConversationID == "" {
					// sign jwt for claim
					jwt, err := helper.SignJWT(r.TweetID, false)
					if err != nil {
						return
					}
					// send direct message
					b, err := service.Message(user.Twitter.AccessToken, r.ID, fmt.Sprintf(config.TwitterGiveawayMessage, "https://opensaucerersgiveaway.onrender.com/claim?token="+jwt))
					if err != nil {
						return
					}

					// parse response
					var messageResponse typing.TwitterMessageResponse

					if err := json.Unmarshal(b, &messageResponse); err != nil {
						return
					}

					if messageResponse.Data.DMConversationID == "" {
						var twitterError typing.TwitterTweetError

						if err := json.Unmarshal(b, &twitterError); err != nil {
							return
						}

						return
					}

					// save message id
					r.ConversationID = messageResponse.Data.DMEventID

					// save reply
					g.SaveReply(&r)
				}

			}
		}
	}
}
