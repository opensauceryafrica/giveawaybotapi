package version

import (
	mux "github.com/gorilla/mux"
	"github.com/opensaucerers/giveawaybot/route"
	v1 "github.com/opensaucerers/giveawaybot/version/v1"
)

// Version1Routes registers all routes for the v1 version
func Version1Routes(r *mux.Router) {

	// this doesn't need versioning, yet
	route.RegisterHomeRoutes(r)

	// V1 routes
	router := r.PathPrefix("/v1").Subrouter()
	v1.RegisterBaseRoutes(router)
	v1.RegisterAuthRoutes(router)

}
