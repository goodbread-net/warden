package routes

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/matthiase/warden/identity"
	"github.com/matthiase/warden/session"
)

type TokenRefreshResponse struct {
	AccessToken string `json:"access_token"`
}

func refreshHandler(w http.ResponseWriter, r *http.Request) {
	authorization := r.Header.Get("Authorization")
	if authorization == "" {
		UnauthorizedError("Missing authorization header").Render(w, r)
		return
	}

	refreshToken := strings.Split(authorization, " ")[1]
	refreshClaims, err := session.Parse(refreshToken, []byte(app.Config.Server.Secret))
	if err != nil {
		log.Printf("%s", err)
		UnauthorizedError("Invalid refresh token").Render(w, r)
		return
	}

	// Look up the user id associated with session id stored in the refresh token
	sessionID := refreshClaims.Subject
	userID, err := app.SessionStore.Find(sessionID)
	if err != nil {
		panic(err)
	}

	if userID == "" {
		NotFoundError().Render(w, r)
		return
	}

	log.Printf("Found user id: %s associated with session %s", userID, sessionID)

	// Look up the user in the datebase
	user, err := app.UserStore.Find(userID)
	if err != nil {
		panic(err)
	}

	// Create a new access token for the user
	identityClaims := identity.NewIdentityClaims(sessionID, user, app.Config)
	accessToken, err := identityClaims.Sign([]byte(app.Config.Server.Secret))
	if err != nil {
		panic(err)
	}

	// Send the access token in the response
	json.NewEncoder(w).Encode(TokenRefreshResponse{
		AccessToken: accessToken,
	})
}
