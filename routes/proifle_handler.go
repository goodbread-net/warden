package routes

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/matthiase/warden/identity"
)

type ProfileResponse struct {
	User *User `json:"user"`
}

func profileHandler(w http.ResponseWriter, r *http.Request) {
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

	// Get the user from the database
	user, err := app.UserStore.Find(identityClaims.Subject)
	if err != nil {
		panic(err)
	}

	if user == nil {
		NotFoundError().Render(w, r)
		return
	}

	json.NewEncoder(w).Encode(ProfileResponse{
		User: &User{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
		},
	})
}
