package typing

import (
	"time"
)

type TwitterAuthResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
	RefreshToken string `json:"refresh_token"`
}

type Twitter struct {
	ID           string    `json:"id" bson:"id"`
	Username     string    `json:"username" bson:"username"`
	AccessToken  string    `json:"access_token " bson:"access_token"`
	RefreshToken string    `json:"refresh_token" bson:"refresh_token"`
	TokenType    string    `json:"token_type" bson:"token_type"`
	RetrievedAt  time.Time `json:"retrieved_at" bson:"retrieved_at"`
	UpdatedAt    time.Time `json:"updated_at" bson:"updated_at"`
}

type TwitterResponse struct {
	Data struct {
		ID            string `json:"id"`
		Name          string `json:"name"`
		Username      string `json:"username"`
		PublicMetrics struct {
			FollowersCount int64 `json:"followers_count"`
			FollowingCount int64 `json:"following_count"`
			TweetCount     int64 `json:"tweet_count"`
		} `json:"public_metrics"`
	} `json:"data"`
}

type TwitterAuthError struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

type TwitterTweetResponse struct {
	Data struct {
		ID       string `json:"id"`
		Text     string `json:"text"`
		AuthorID string `json:"author_id"`
	} `json:"data"`
}

type TwitterTweetError struct {
	Errors []struct {
		Value     string `json:"value"`
		Title     string `json:"title"`
		Detail    string `json:"detail"`
		Parameter string `json:"parameter"`
		Message   string `json:"message"`
		Type      string `json:"type"`
	} `json:"errors"`
	Title     string `json:"title"`
	Detail    string `json:"detail"`
	Message   string `json:"message"`
	Parameter string `json:"parameter"`
	Type      string `json:"type"`
}

type TwitterListTweetsResponse struct {
	Data []struct {
		ID        string `json:"id"`
		Text      string `json:"text"`
		CreatedAt string `json:"created_at"`
	} `json:"data"`
	Meta struct {
		NewestID    string `json:"newest_id"`
		OldestID    string `json:"oldest_id"`
		ResultCount int    `json:"result_count"`
		NextToken   string `json:"next_token"`
	} `json:"meta"`
}

type TwitterLikeTweetResponse struct {
	Data struct {
		Liked bool `json:"liked"`
	} `json:"data"`
}

type TwitterDeleteTweetResponse struct {
	Data struct {
		Deteted bool `json:"deleted"`
	} `json:"data"`
}

type TwitterRetweetResponse struct {
	Data struct {
		Retweeted bool `json:"retweeted"`
	} `json:"data"`
}

type TwitterListResponse struct {
	Data []struct {
		ID             string `json:"id"`
		Name           string `json:"name"`
		Username       string `json:"username"`
		AuthorID       string `json:"author_id"`
		Text           string `json:"text"`
		ConversationID string `json:"conversation_id"`
	}
	Meta struct {
		ResultCount *int   `json:"result_count"`
		NextToken   string `json:"next_token"`
	}
	Includes struct {
		Users []struct {
			ID       string `json:"id"`
			Name     string `json:"name"`
			Username string `json:"username"`
		} `json:"users"`
	} `json:"includes"`
}

type TwitterEmbedResponse struct {
	HTML string `json:"html"`
	URL  string `json:"url"`
}

type Reply struct {
	ID       string `json:"id" bson:"id"` // twitter id
	Text     string `json:"text" bson:"text"`
	Username string `json:"username" bson:"username"`
	TweetID  string `json:"tweet_id" bson:"tweet_id"`
	FText    string `json:"f_text" bson:"f_text"`
}
