package routes

import (
	"encoding/json"
	"net/http"

	"github.com/matthiase/warden/models"
	"github.com/matthiase/warden/verification"
)

type RegistrationRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

type RegistrationResponse struct {
	User              *User  `json:"user"`
	VerificationToken string `json:"verification_token"`
}

func registrationHandler(w http.ResponseWriter, r *http.Request) {
	var data RegistrationRequest
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		panic(err)
	}

	// TODO: validate the email address and name

	user, err := app.UserStore.Create(data.FirstName, data.LastName, data.Email)
	if err != nil {
		if err.Error() == models.ErrUserDuplicateEmail {
			ConflictError("Email address is already registered").Render(w, r)
			return
		} else {
			panic(err)
		}
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

	json.NewEncoder(w).Encode(RegistrationResponse{
		User: &User{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
		},
		VerificationToken: verificationToken,
	})
}
