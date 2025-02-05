package routes

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/matthiase/warden/identity"
)

type AuthenticationResponse struct {
	*User       `json:"user"`
	AccessToken string `json:"access_token"`
}

func authenticationHandler(w http.ResponseWriter, r *http.Request) {
	// Get the authorization token from the request header
	authorizationHeader := r.Header.Get("Authorization")
	if authorizationHeader == "" {
		UnauthorizedError("Missing authorization header").Render(w, r)
		return
	}

	// Parse the authorization token and ensure that it is valid
	identityToken := strings.Split(authorizationHeader, " ")[1]
	identityClaims, err := identity.Parse(identityToken, []byte(app.Config.Server.Secret))
	if err != nil {
		UnauthorizedError("Invalid identity token").Render(w, r)
		return
	}

	// Get and log the expiration time of the token
	log.Printf("Token expires at: %s", identityClaims.ExpiresAt)

	// Get the user from the database
	user, err := app.UserStore.Find(identityClaims.Subject)
	if err != nil {
		panic(err)
	}

	if user == nil {
		NotFoundError().Render(w, r)
		return
	}

	json.NewEncoder(w).Encode(AuthenticationResponse{
		User: &User{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
		},
		AccessToken: identityToken,
	})
}
