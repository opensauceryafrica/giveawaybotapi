package giveaway

import (
	"context"
	"net/http"

	"github.com/opensaucerers/giveawaybot/helper"
	"github.com/opensaucerers/giveawaybot/logic/v1/giveaway"
	giveawayr "github.com/opensaucerers/giveawaybot/repository/v1/giveaway"
	"github.com/opensaucerers/giveawaybot/repository/v1/user"
	"github.com/opensaucerers/giveawaybot/typing"
)

// Reward tweets the winners and sends a direct message to the winners of the giveaway
func Reward(w http.ResponseWriter, r *http.Request) {
	// get user from context
	id := r.Context().Value(typing.AuthCtxKey{}).(*user.User).Twitter.ID

	giveaway, err := giveaway.Reward(id)
	if err != nil {
		helper.SendJSONResponse(w, false, http.StatusBadRequest, err.Error(), nil)
		return
	}

	// add user to context
	ctx := context.WithValue(r.Context(), typing.AuthCtxKey{}, &giveaway.Author)

	// load request with context
	r = r.WithContext(ctx)

	// send response
	helper.SendJSONResponse(w, true, http.StatusOK, "Reward process initiated", typing.M{"giveaway": giveaway}, true, r)
}

// Winners registers the winners of the giveaway
func Winners(w http.ResponseWriter, r *http.Request) {
	// get user from context
	id := r.Context().Value(typing.AuthCtxKey{}).(*user.User).Twitter.ID

	g := r.Context().Value(typing.BodyCtxKey{}).(*giveawayr.Giveaway)

	if len(g.Winners) == 0 {
		helper.SendJSONResponse(w, false, http.StatusBadRequest, "No winners specified", nil)
		return
	}

	giveaway, err := giveaway.Winners(id, g.Winners)
	if err != nil {
		helper.SendJSONResponse(w, false, http.StatusBadRequest, err.Error(), nil)
		return
	}

	// add user to context
	ctx := context.WithValue(r.Context(), typing.AuthCtxKey{}, &giveaway.Author)

	// load request with context
	r = r.WithContext(ctx)

	// send response
	helper.SendJSONResponse(w, true, http.StatusOK, "Winners recorded", typing.M{"giveaway": giveaway}, true, r)
}
