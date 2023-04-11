package route

import (
	"github.com/opensaucerers/giveawaybot/controller"

	mux "github.com/gorilla/mux"
)

func RegisterHomeRoutes(r *mux.Router) {

	r.HandleFunc("/", controller.Home).Methods("GET")
}
