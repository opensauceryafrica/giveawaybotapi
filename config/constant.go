package config

import (
	"os"

	"github.com/opensaucerers/giveawaybot/typing"
)

const (

	// EnvTagName is the tag name for environment variables struct
	envTagName = "env"

	// ShutdownTimeout is the time to wait for the server to shutdown gracefully
	ShutdownTimeout = 5 // seconds

	//maxconnections is the maximum number of connections in the pgx pool
	MaxConnections = 15
)

const (
	// UserCollection is the name of the user collection
	UserCollection = "users"

	// RaidCollection is the name of the raid collection
	GiveawayCollection = "giveaways"
)

const (

	// DefaultAvatar is the default avatar for users
	DefaultAvatar = "https://e7.pngegg.com/pngimages/84/165/png-clipart-united-states-avatar-organization-information-user-avatar-service-computer-wallpaper-thumbnail.png"

	// DefaultPageLimit is the default page limit
	DefaultPageLimit = 10

	//DefaultPageOffset is the default page to skip
	DefaultPageOffset = 0
)

const (
	// TwitterScope is the scope for twitter oauth
	TwitterScope = "like.read%20like.write%20tweet.read%20tweet.write%20users.read%20offline.access%20follows.write%20follows.read"

	// TwitterGiveawayTweet is the tweet to quote for giveaways
	TwitterGiveawayTweet = "1646776307919798272"

	// TwitterGiveawayTweet is the comment to add to the giveaway tweet
	TwitterGiveawayComment = "Weekly data giveaway 🥂\nPlease find the rules below.\n1. The giveaway is only open for 5 hours\n2. Comment your username under this tweet like this @opensaucerer\n3. 10 usernames will be selected randomly as winners\n4. Don't comment twice\nDo like and share to others. Thank you 🫴"
)

var (
	// Env is the global environment variable
	Env = new(typing.Env) // global environment variable

	// ShutdownChan is the channel to listen for shutdown signals
	ShutdownChan = make(chan os.Signal, 1)
)
