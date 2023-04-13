package v1

import (
	"github.com/gorilla/mux"
	"github.com/opensaucerers/giveawaybot/route/v1/auth"
)

func RegisterAuthRoutes(r *mux.Router) {

	router := r.PathPrefix("/auth").Subrouter()

	auth.RegisterTwitterRoutes(router)
}
