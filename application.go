package authnz

import (
	"github.com/birdbox/authnz/config"
	"github.com/birdbox/authnz/data"
)

type Application struct {
	UserStore     data.UserStore
	SessionStore  data.SessionStore
	PasscodeStore data.PasscodeStore
}

func NewApplication(cfg *config.Config) (*Application, error) {
	userStore, err := data.NewUserStore()
	if err != nil {
		return nil, err
	}

	sessionStore, err := data.NewSessionStore()
	if err != nil {
		return nil, err
	}

	passcodeStore, err := data.NewPasscodeStore()
	if err != nil {
		return nil, err
	}

	return &Application{
		UserStore:     userStore,
		SessionStore:  sessionStore,
		PasscodeStore: passcodeStore,
	}, nil
}
