package routes

import (
	"log"
	"net/http"
	"time"

	"github.com/matthiase/warden/session"
)

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	var sessionClaims *session.SessionClaims
	if sessionTokenCookie, err := r.Cookie(app.Config.Session.Name + "_st"); err == nil {
		sessionClaims, _ = session.Parse(sessionTokenCookie.Value, []byte(app.Config.Server.Secret))
	}

	if sessionClaims.Subject != "" {
		if err := app.SessionStore.Revoke(sessionClaims.Subject); err != nil {
			log.Printf("Failed to revoke session: %v", err)
		}
	}

	// Clear the identity token and session token cookies
	http.SetCookie(w, &http.Cookie{
		Name:     app.Config.Session.Name + "_it",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   app.Config.Session.Secure,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().UTC().Add(-1 * time.Second),
	})

	http.SetCookie(w, &http.Cookie{
		Name:     app.Config.Session.Name + "_st",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   app.Config.Session.Secure,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().UTC().Add(-1 * time.Second),
	})
}
