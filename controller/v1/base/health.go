package base

import (
	"net/http"

	"github.com/opensaucerers/giveawaybot/helper"
	"github.com/opensaucerers/giveawaybot/typing"
)

//Health checks health information of the app in relation to a specific version
func Health(w http.ResponseWriter, r *http.Request) {

	data := typing.Health{
		Name:    "Golang Server Template",
		Status:  true,
		Version: "1.0.0",
	}

	helper.SendJSONResponse(w, true, http.StatusOK, "Health", typing.M{"health": data})
}
