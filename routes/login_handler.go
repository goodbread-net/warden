package routes

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/matthiase/warden/verification"
)

type AuthenticationRequest struct {
	Email string `json:"email"`
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var data AuthenticationRequest
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		panic(err)
	}

	user, err := app.UserStore.FindByEmail(data.Email)
	if err != nil {
		panic(err)
	}

	passcode, err := app.PasscodeStore.Create(user.ID)
	if err != nil {
		panic(err)
	}

	app.Mailer.Send(user.Email, "login", map[string]interface{}{
		"Application":   app.Config.Application,
		"RecipientName": user.FirstName,
		"Passcode":      passcode,
	})

	// Create a verification token for the user. This token will be used in
	// conjunction with the passcode to confirm the user's identity.
	verificationClaims := verification.NewVerificationClaims(user.ID, app.Config)
	verificationToken, err := verificationClaims.Sign([]byte(app.Config.Server.Secret))
	if err != nil {
		panic(err)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     app.Config.Session.Name + "_vt",
		Value:    verificationToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   app.Config.Session.Secure,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().UTC().Add(300 * time.Second),
	})
}
