package giveaway

import (
	"context"
	"net/http"

	"github.com/opensaucerers/giveawaybot/helper"
	"github.com/opensaucerers/giveawaybot/logic/v1/giveaway"
	"github.com/opensaucerers/giveawaybot/repository/v1/user"
	"github.com/opensaucerers/giveawaybot/typing"
)

// Replies retrieves all replies to a giveaway
func Replies(w http.ResponseWriter, r *http.Request) {
	// get user from context
	id := r.Context().Value(typing.AuthCtxKey{}).(*user.User).Twitter.ID

	giveaway, err := giveaway.Replies(id)
	if err != nil {
		helper.SendJSONResponse(w, false, http.StatusBadRequest, err.Error(), nil)
		return
	}

	// add user to context
	ctx := context.WithValue(r.Context(), typing.AuthCtxKey{}, &giveaway.Author)

	// load request with context
	r = r.WithContext(ctx)

	// send response
	helper.SendJSONResponse(w, true, http.StatusOK, "Replies retrieved", typing.M{"giveaway": giveaway}, true, r)
}

// Replies retrieves all replies to a giveaway
func Report(w http.ResponseWriter, r *http.Request) {
	// get user from context
	id := r.Context().Value(typing.ParamsCtxKey{}).(map[string]string)["id"]
	if id == "" {
		helper.SendJSONResponse(w, false, http.StatusBadRequest, "Invalid giveaway ID", nil)
		return
	}

	report, err := giveaway.Report(id)
	if err != nil {
		helper.SendJSONResponse(w, false, http.StatusBadRequest, err.Error(), nil)
		return
	}

	// send response
	helper.SendJSONResponse(w, true, http.StatusOK, "Report generated", typing.M{"report": report})
}
