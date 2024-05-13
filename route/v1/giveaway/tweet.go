package giveaway

import (
	"github.com/gorilla/mux"
	"github.com/opensaucerers/giveawaybot/controller/v1/giveaway"
	"github.com/opensaucerers/giveawaybot/middleware/v1"
)

func RegisterTweetRoutes(r *mux.Router) {

	router := r.PathPrefix("/tweet").Subrouter()

	// making them GET methods for now while there's no UI
	router.Handle("/replies", middleware.AuthF(giveaway.Replies)).Methods("GET")
	router.Handle("/report/{id}", middleware.ParamsF(giveaway.Report)).Methods("GET")
}
