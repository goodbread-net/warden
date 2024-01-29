package routes

import (
	"log"
	"net/http"

	"github.com/birdbox/authnz"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
)

var (
	app *authnz.Application
)

func NewHandler(application *authnz.Application) http.Handler {
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

	allowedOrigins := []string{"http://localhost:3000"}
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

	router.Get("/healthcheck", Healthcheck)
	router.Get("/sessions/me", getCurrentUser)
	router.Route("/registration", func(r chi.Router) {
		r.Post("/start", createUser)
		r.Post("/confirm", confirmUser)
	})

	return router
}

func methodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	MethodNotAllowedError().Render(w, r)
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	NotFoundError().Render(w, r)
}
