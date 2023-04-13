package auth

import (
	"context"
	"net/http"

	"github.com/opensaucerers/giveawaybot/config"
	"github.com/opensaucerers/giveawaybot/helper"
	"github.com/opensaucerers/giveawaybot/logic/v1/auth"
	"github.com/opensaucerers/giveawaybot/typing"
)

// TwitterAuth is a controller for twitter auth
func TwitterAuth(w http.ResponseWriter, r *http.Request) {
	url := auth.TwitterAuth()
	helper.SendJSONResponse(w, true, http.StatusOK, "Twitter auth url", typing.M{"url": url})
}

// TwitterSignon is a controller for intercepting twitter auth redirect and creating a new user
func TwitterSignon(w http.ResponseWriter, r *http.Request) {
	// TODO: add validation for query params
	// get code and address from query params
	code := r.URL.Query().Get("code")
	if code == "" {
		helper.SendJSONResponse(w, false, http.StatusBadRequest, "Twitter code is required", nil)
		return
	}

	user, err := auth.TwitterSignon(code)
	if err != nil {
		if err.Error() == config.ErrTwitterUnauthorized {
			helper.SendJSONResponse(w, false, http.StatusUnauthorized, err.Error(), nil)
			return
		}
		helper.SendJSONResponse(w, false, http.StatusBadRequest, err.Error(), nil)
		return
	}

	// add user to context
	ctx := context.WithValue(r.Context(), typing.AuthCtxKey{}, user)

	// load request with context
	r = r.WithContext(ctx)

	// send response
	helper.SendJSONResponse(w, true, http.StatusOK, "Signon successful", typing.M{"user": user}, true, r)
}
