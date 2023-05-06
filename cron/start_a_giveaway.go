package cron

import (
	giveawayl "github.com/opensaucerers/giveawaybot/logic/v1/giveaway"
	"github.com/opensaucerers/giveawaybot/repository/v1/giveaway"
	"github.com/opensaucerers/giveawaybot/repository/v1/user"
)

// StartAGiveaway closes any active giveaways and starts a new one
func StartAGiveaway() {

	u := user.User{
		Username: "opensaucerers",
	}

	if err := u.FindByUsername(); err != nil {
		return
	}

	if u.ID.IsZero() {
		return
	}

	// refresh twitter token
	if err := u.RefreshTwitterAccessToken(); err != nil {
		return
	}

	// get active giveaway
	g, err := giveaway.Running(u)
	if err != nil {
		return
	}

	// close active giveaway
	if g != nil && !g.ID.IsZero() {
		if err := g.Close(); err != nil {
			return
		}
	}

	// create new giveaway
	if _, err = giveawayl.Start(u.Twitter.ID); err != nil {
		return
	}

}
