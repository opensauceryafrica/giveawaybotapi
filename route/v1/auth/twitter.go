package auth

import (
	"github.com/opensaucerers/giveawaybot/controller/v1/auth"

	mux "github.com/gorilla/mux"
)

func RegisterTwitterRoutes(r *mux.Router) {

	router := r.PathPrefix("/twitter").Subrouter()

	// making them GET methods for now while there's no UI
	router.HandleFunc("/begin", auth.TwitterAuth).Methods("GET")
	router.HandleFunc("/signon", auth.TwitterSignon).Methods("GET")

}
