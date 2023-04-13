package v1

import (
	"github.com/gorilla/mux"
	"github.com/opensaucerers/giveawaybot/route/v1/base"
)

func RegisterBaseRoutes(r *mux.Router) {

	router := r.PathPrefix("/").Subrouter()

	base.RegisterHealthRoutes(router)
}
