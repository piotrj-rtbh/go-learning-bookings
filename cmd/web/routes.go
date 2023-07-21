package main

import (
	"net/http"

	"github.com/piotrj-rtbh/bookings/pkg/config"
	"github.com/piotrj-rtbh/bookings/pkg/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// we need to get the routes out from main and put in here
// we'll use an external package pat from https://github.com/bmizerany/pat

func routes(app *config.AppConfig) http.Handler {
	// multiplexer
	// go mod tidy => cleans up the go.mod file if we had previously pat installed
	mux := chi.NewRouter()

	// middleware'y pomagają podpiąć się pod żądanie i zrobić coś pomiędzy
	// installing middleware
	mux.Use(middleware.Recoverer) // prevents from panic errors! middleware comes from chi
	// use our own middleware as well!
	// this middleware writes to console at every page hit
	// mux.Use(WriteToConsole)

	mux.Use(NoSurf)

	// use sessions using the middleware SessionLoad
	mux.Use(SessionLoad)

	// defining paths
	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)
	return mux
}
