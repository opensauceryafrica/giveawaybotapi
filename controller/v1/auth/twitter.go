package user

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

// TwitterRegister is a controller for intercepting twitter auth redirect and creating a new user
func TwitterRegister(w http.ResponseWriter, r *http.Request) {
	// TODO: add validation for query params
	// get code and address from query params
	code := r.URL.Query().Get("code")

	user, err := auth.TwitterRegister(code)
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
	helper.SendJSONResponse(w, true, http.StatusOK, "Registration successful", typing.M{"user": user}, true, r)
}

// TwitterLogin is a controller for intercepting twitter auth redirect and logging in an existing user
func TwitterLogin(w http.ResponseWriter, r *http.Request) {
	// TODO: add validation for query params
	// get code and address from query params
	code := r.URL.Query().Get("code")
	address := r.URL.Query().Get("address")
	if code == "" || address == "" {
		helper.SendJSONResponse(w, false, http.StatusBadRequest, "Twitter code and wallet address are required", nil)
		return
	}

	user, err := auth.TwitterLogin(code)
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
	helper.SendJSONResponse(w, true, http.StatusOK, "Login successful", typing.M{"user": user}, true, r)
}
