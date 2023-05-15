package v1

import (
	"github.com/gorilla/mux"
	"github.com/opensaucerers/giveawaybot/route/v1/giveaway"
)

func RegisterGiveawayRoutes(r *mux.Router) {

	router := r.PathPrefix("/giveaway").Subrouter()

	giveaway.RegisterSimulateRoutes(router)
	giveaway.RegisterTweetRoutes(router)
	giveaway.RegisterRewardRoutes(router)
}
