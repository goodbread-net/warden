package routes

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/matthiase/warden/identity"
	"github.com/matthiase/warden/session"
)

type ConfirmationRequest struct {
	Passcode string `json:"passcode"`
}

type ConfirmationResponse struct {
	IdentityToken string `json:"identity_token"`
	SessionToken  string `json:"session_token"`
}

func confirmationHandler(w http.ResponseWriter, r *http.Request) {
	// Get the user-provided passcode from the request body and use it to
	// to look up the associated user id.
	var data ConfirmationRequest
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		BadRequestError(err.Error()).Render(w, r)
		return
	}

	if data.Passcode == "" {
		BadRequestError("Missing passcode").Render(w, r)
		return
	}

	expectedUserID, err := app.PasscodeStore.Find(data.Passcode)
	if err != nil {
		UnauthorizedError("Invalid passcode").Render(w, r)
		return
	}

	// Get the verification token from the request header
	authorizationHeader := r.Header.Get("Authorization")
	if authorizationHeader == "" {
		UnauthorizedError("Missing authorization header").Render(w, r)
		return
	}

	// Parse the verification token and ensure that it matches the user id
	verificationToken := strings.Split(authorizationHeader, " ")[1]
	verificationClaims, err := identity.Parse(verificationToken, []byte(app.Config.Server.Secret))
	if err != nil {
		UnauthorizedError("Invalid verification token").Render(w, r)
		return
	}

	providedUserID := verificationClaims.Subject

	if expectedUserID != providedUserID {
		UnauthorizedError("Invalid authentication token").Render(w, r)
		return
	}

	// Fetch the user record from the database and create a new session.
	user, err := app.UserStore.Find(providedUserID)
	if err != nil {
		UnauthorizedError("Invalid authentication token").Render(w, r)
		return
	}

	// Print the user's name to the console
	log.Printf("User %s %s (%s) has confirmed their identity", user.FirstName, user.LastName, user.Email)

	sessionID, err := app.SessionStore.Create(user.ID)
	if err != nil {
		panic(err)
	}

	// Generate the session and identity tokens and set them as http-only
	// cookies in the response.
	secret := []byte(app.Config.Server.Secret)

	// The session cookie will be used to refresh expired identity tokens
	sessionClaims := session.NewSessionClaims(sessionID, app.Config)
	sessionToken, err := sessionClaims.Sign(secret)
	if err != nil {
		panic(err)
	}

	// The identity token will be used to authenticate requests
	identityClaims := identity.NewIdentityClaims(sessionID, user, app.Config)
	identityToken, err := identityClaims.Sign(secret)
	if err != nil {
		panic(err)
	}

	json.NewEncoder(w).Encode(ConfirmationResponse{
		IdentityToken: identityToken,
		SessionToken:  sessionToken,
	})
}
