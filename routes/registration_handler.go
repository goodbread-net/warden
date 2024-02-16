package routes

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/matthiase/warden/verification"
)

type RegistrationRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

func registrationHandler(w http.ResponseWriter, r *http.Request) {
	var data RegistrationRequest
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		panic(err)
	}

	// TODO: validate the email address and name

	// TODO: ensure the email address is not already registered

	user, err := app.UserStore.Create(data.FirstName, data.LastName, data.Email)
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

	if app.Config.Session.Secure {
		http.SetCookie(w, &http.Cookie{
			Name:     app.Config.Session.Name + "_vt",
			Value:    verificationToken,
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteNoneMode,
			Expires:  time.Now().UTC().Add(300 * time.Second),
		})
	} else {
		http.SetCookie(w, &http.Cookie{
			Name:     app.Config.Session.Name + "_vt",
			Value:    verificationToken,
			Path:     "/",
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
			Expires:  time.Now().UTC().Add(300 * time.Second),
		})
	}
}
