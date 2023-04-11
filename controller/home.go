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
		Description: "Golang Server Template",
	}

	helper.SendJSONResponse(w, true, http.StatusOK, "Home", typing.M{"home": data})
}
