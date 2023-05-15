package giveaway

import (
	"github.com/gorilla/mux"
	"github.com/opensaucerers/giveawaybot/controller/v1/giveaway"
	"github.com/opensaucerers/giveawaybot/middleware/v1"
	giveawayr "github.com/opensaucerers/giveawaybot/repository/v1/giveaway"
)

func RegisterRewardRoutes(r *mux.Router) {

	router := r.PathPrefix("/reward").Subrouter()

	// making them GET methods for now while there's no UI
	router.Handle("/reward", middleware.AuthF(giveaway.Reward)).Methods("POST")
	router.Handle("/winners", middleware.AuthH(middleware.BodyF(&giveawayr.Giveaway{})(giveaway.Winners))).Methods("POST")

}
