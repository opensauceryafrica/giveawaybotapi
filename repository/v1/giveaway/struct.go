package giveaway

import (
	"time"

	"github.com/opensaucerers/giveawaybot/repository/v1/user"
	"github.com/opensaucerers/giveawaybot/typing"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Giveaways []Giveaway

type Giveaway struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Author       user.User          `json:"author"`
	Active       bool               `json:"active" bson:"active"`
	Replies      []typing.Reply     `json:"replies" bson:"replies"`
	Winners      []string           `json:"winners" bson:"winners"`
	TotalReplies int                `json:"total_replies" bson:"total_replies"`
	Completed    bool               `json:"completed" bson:"completed"`
	Tweet        string             `json:"tweet" bson:"tweet"`
	TweetID      string             `json:"tweet_id" bson:"tweet_id"`
	TwitterURL   string             `json:"twitter_url" bson:"twitter_url"`
	TwitterHTML  string             `json:"twitter_html" bson:"twitter_html"`
	Rewarded     bool               `json:"rewarded" bson:"rewarded"`
	ReportTweet  string             `json:"report_tweet" bson:"report_tweet"`
	WinnersTweet []string           `json:"winners_tweet" bson:"winners_tweet"`
	AmountSpent  float64            `json:"amount_spent" bson:"amount_spent"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at"`
	CompletedAt  time.Time          `json:"completed_at" bson:"completed_at"`
}

type Report struct {
	TotalReplies      int               `json:"total_replies"`
	TotalDuplicates   int               `json:"total_duplicates"`
	DuplicateReplies  map[string]string `json:"duplicate_replies"`
	TotalValidReplies int               `json:"total_valid_replies"`
	ValidReplies      map[string]string `json:"valid_replies"`
	ValidRepliesList  []string          `json:"valid_replies_list"`
}

type Replies []typing.Reply
