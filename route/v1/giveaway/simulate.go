package giveaway

import (
	"github.com/gorilla/mux"
	"github.com/opensaucerers/giveawaybot/controller/v1/giveaway"
	"github.com/opensaucerers/giveawaybot/middleware/v1"
)

func RegisterSimulateRoutes(r *mux.Router) {

	router := r.PathPrefix("/simulate").Subrouter()

	// making them GET methods for now while there's no UI
	router.Handle("/start", middleware.AuthF(giveaway.Start)).Methods("GET")
	router.Handle("/disrupt", middleware.AuthF(giveaway.Disrupt)).Methods("GET")
	router.Handle("/end", middleware.AuthF(giveaway.End)).Methods("GET")

}
