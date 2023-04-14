package giveaway

import (
	"context"
	"net/http"

	"github.com/opensaucerers/giveawaybot/helper"
	"github.com/opensaucerers/giveawaybot/logic/v1/giveaway"
	"github.com/opensaucerers/giveawaybot/repository/v1/user"
	"github.com/opensaucerers/giveawaybot/typing"
)

// Start is a controller for starting a giveaway
func Start(w http.ResponseWriter, r *http.Request) {
	// get user from context
	id := r.Context().Value(typing.AuthCtxKey{}).(*user.User).Twitter.ID

	giveaway, err := giveaway.Start(id)
	if err != nil {
		helper.SendJSONResponse(w, false, http.StatusBadRequest, err.Error(), nil)
		return
	}

	// add user to context
	ctx := context.WithValue(r.Context(), typing.AuthCtxKey{}, &giveaway.Author)

	// load request with context
	r = r.WithContext(ctx)

	// send response
	helper.SendJSONResponse(w, true, http.StatusOK, "Giveaway started", typing.M{"giveaway": giveaway}, true, r)
}

// Disrupt is a controller for disrupting a giveaway
func Disrupt(w http.ResponseWriter, r *http.Request) {
	// get user from context
	id := r.Context().Value(typing.AuthCtxKey{}).(*user.User).Twitter.ID

	giveaway, err := giveaway.Disrupt(id)
	if err != nil {
		helper.SendJSONResponse(w, false, http.StatusBadRequest, err.Error(), nil)
		return
	}

	// add user to context
	ctx := context.WithValue(r.Context(), typing.AuthCtxKey{}, &giveaway.Author)

	// load request with context
	r = r.WithContext(ctx)

	// send response
	helper.SendJSONResponse(w, true, http.StatusOK, "Giveaway disrupted", typing.M{"giveaway": giveaway}, true, r)
}

// End is a controller for ending a giveaway
func End(w http.ResponseWriter, r *http.Request) {
	// get user from context
	id := r.Context().Value(typing.AuthCtxKey{}).(*user.User).Twitter.ID

	giveaway, err := giveaway.End(id)
	if err != nil {
		helper.SendJSONResponse(w, false, http.StatusBadRequest, err.Error(), nil)
		return
	}

	// add user to context
	ctx := context.WithValue(r.Context(), typing.AuthCtxKey{}, &giveaway.Author)

	// load request with context
	r = r.WithContext(ctx)

	// send response
	helper.SendJSONResponse(w, true, http.StatusOK, "Giveaway ended", typing.M{"giveaway": giveaway}, true, r)
}
