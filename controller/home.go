package controller

import (
	"net/http"

	"github.com/opensaucerers/giveawaybot/helper"
	"github.com/opensaucerers/giveawaybot/typing"
)

func Home(w http.ResponseWriter, r *http.Request) {

	data := typing.Home{
		Status:      true,
		Version:     "1.0.0",
		Name:        "Giveaway Bot",
		Description: "It's a great thing to win but it's an even greater thing to help others win just as you already have.",
		Twitter:     "https://twitter.com/opensaucerers",
	}

	helper.SendJSONResponse(w, true, http.StatusOK, "Home", typing.M{"home": data})
}
