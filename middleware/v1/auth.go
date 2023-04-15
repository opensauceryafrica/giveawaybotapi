package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/opensaucerers/giveawaybot/helper"
	"github.com/opensaucerers/giveawaybot/repository/v1/user"
	"github.com/opensaucerers/giveawaybot/typing"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Auth is checking if the user is logged in
func Auth(w http.ResponseWriter, r *http.Request) context.Context {

	// get auth header
	authHeader := r.Header.Get("Authorization")
	useBearer := true
	token := ""

	// TODO: remove this but keep for now since we have to UI and we're
	// defualting all calls to GET method
	if authHeader == "" {
		token = r.URL.Query().Get("token")
		useBearer = false
	}

	if useBearer {
		// validate auth header
		if authHeader == "" {
			helper.SendJSONResponse(w, false, http.StatusUnauthorized, "User not logged in", nil)
			return nil
		}
		// split auth header
		authHeaderSplit := strings.Split(authHeader, " ")
		// validate auth header split
		if len(authHeaderSplit) != 2 {
			helper.SendJSONResponse(w, false, http.StatusUnauthorized, "User session not invalid", nil)
			return nil
		}
		// validate auth header split
		if authHeaderSplit[0] != "Bearer" {
			helper.SendJSONResponse(w, false, http.StatusUnauthorized, "User session not valid", nil)
			return nil
		}
		// get token from auth header
		token = authHeaderSplit[1]
	}

	// validate token
	if token == "" {
		helper.SendJSONResponse(w, false, http.StatusUnauthorized, "User not logged in", nil)
		return nil
	}
	// validate jwt token
	claim, valid := helper.VerifyJWT(token)
	if !valid {
		helper.SendJSONResponse(w, false, http.StatusUnauthorized, "User session not valid", nil)
		return nil
	}

	user := user.User{Twitter: typing.Twitter{ID: claim.ID}}
	// get user from token
	if err := user.FindSocial(typing.Social{Twitter: true}); err != nil {
		helper.SendJSONResponse(w, false, http.StatusUnauthorized, "Error getting user from token: "+err.Error(), nil)
		return nil
	}

	// validate user
	if user.ID.String() == primitive.NilObjectID.String() {
		helper.SendJSONResponse(w, false, http.StatusUnauthorized, "User session not valid", nil)
		return nil
	}

	// set user in context
	return context.WithValue(r.Context(), typing.AuthCtxKey{}, &user)
}

// AuthH is a middleware that checks if the user is logged in
// by validating the token provided in the Authorization header
func AuthH(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get auth context
		ctx := Auth(w, r)
		if ctx == nil {
			return
		}
		// call next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AuthF is a middleware that checks if the user is logged in
// by validating the token provided in the Authorization header
func AuthF(next func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get auth context
		ctx := Auth(w, r)
		if ctx == nil {
			return
		}
		// call next handler
		next(w, r.WithContext(ctx))
	})
}
