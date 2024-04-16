package routes

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/matthiase/warden"
	"github.com/rs/cors"
)

var (
	app *warden.Application
)

func NewHandler(application *warden.Application) http.Handler {
	app = application
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.SetHeader("Content-Type", "application/json"))
	router.Use(middleware.AllowContentType("application/json"))

	router.MethodNotAllowed(methodNotAllowedHandler)
	router.NotFound(notFoundHandler)

	allowedOrigins := []string{"http://localhost:3000", "https://goodbread.vercel.app", "https://preview.goodbread.net"}
	if len(allowedOrigins) == 0 {
		log.Fatal("HTTP server unable to start - expected ALLOWED_ORIGINS")
	}

	cors := cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	router.Use(cors.Handler)

	router.Get("/healthcheck", healthcheckHandler)
	router.Route("/api/v1", func(r chi.Router) {
		r.Route("/users", func(r chi.Router) {
			r.Post("/register", registrationHandler)
			r.Post("/login", loginHandler)
		})
		r.Route("/sessions", func(r chi.Router) {
			r.Post("/confirm", confirmationHandler)
			r.Get("/authenticate", authenticationHandler)
			r.Post("/refresh", refreshHandler)
			r.Get("/profile", profileHandler)
		})
	})

	return router
}

func methodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	MethodNotAllowedError().Render(w, r)
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	NotFoundError().Render(w, r)
}
