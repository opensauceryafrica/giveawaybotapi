package helper

import (
	"encoding/json"
	"net/http"

	"github.com/opensaucerers/giveawaybot/repository/v1/user"

	"github.com/opensaucerers/giveawaybot/typing"
)

// SendJSONResponse sends a JSON response
func SendJSONResponse(w http.ResponseWriter, status bool, statusCode int, message string, data map[string]interface{}, opts ...interface{}) {
	w.Header().Set("Content-Type", "application/json")
	// if withNewToken, generate a new token and add it to the response
	if len(opts) > 0 && opts[0].(bool) {
		// request must be passed in as an option if withNewToken is true
		if len(opts) > 1 && opts[1].(*http.Request) != nil {
			// get user address from context
			id := opts[1].(*http.Request).Context().Value(typing.AuthCtxKey{}).(*user.User).Twitter.ID
			// generate new token
			token, err := SignJWT(id)
			if err != nil {
				// internal server error
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(typing.Response{
					Status:  false,
					Message: "Error generating new token: " + err.Error(),
					Data:    nil,
				})
			}
			// add token to response
			data["token"] = token
		}
	}
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(typing.Response{
		Status:  status,
		Message: message,
		Data:    data,
	})
}
