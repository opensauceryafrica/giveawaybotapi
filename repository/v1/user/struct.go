package user

import (
	"time"

	"github.com/opensaucerers/giveawaybot/typing"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID             primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username       string             `json:"username" bson:"username"`
	Avatar         string             `json:"avatar" bson:"avatar"`
	UseTwitterAuth bool               `json:"use_twitter_auth" bson:"use_twitter_auth"`
	Twitter        typing.Twitter     `json:"twitter" bson:"twitter"`
	CreatedAt      time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at" bson:"updated_at"`
}
