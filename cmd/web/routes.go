package main

import (
	"net/http"

	"github.com/piotrj-rtbh/bookings/internal/config"
	"github.com/piotrj-rtbh/bookings/internal/handlers"

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

	mux.Use(NoSurf) // this will prevent POST requests without proper CSRF token

	// use sessions using the middleware SessionLoad
	mux.Use(SessionLoad)

	// defining paths
	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)
	mux.Get("/generals-quarters", handlers.Repo.Generals)
	mux.Get("/majors-suite", handlers.Repo.Majors)
	mux.Get("/search-availability", handlers.Repo.Availability)
	mux.Post("/search-availability", handlers.Repo.PostAvailability)
	mux.Post("/search-availability-json", handlers.Repo.AvailabilityJSON)
	mux.Get("/choose-room/{id}", handlers.Repo.ChooseRoom)
	mux.Get("/book-room", handlers.Repo.BookRoom)

	mux.Get("/contact", handlers.Repo.Contact)

	mux.Get("/make-reservation", handlers.Repo.Reservation)
	mux.Post("/make-reservation", handlers.Repo.PostReservation)
	mux.Get("/reservation-summary", handlers.Repo.ReservationSummary)

	mux.Get("/user/login", handlers.Repo.ShowLogin)
	mux.Post("/user/login", handlers.Repo.PostShowLogin)
	mux.Get("/user/logout", handlers.Repo.Logout)

	// in order to enable images loading we have to run a file server
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	// anything starting with /admin will be handled by this func which requires authentication:
	mux.Route("/admin", func(mux chi.Router) {
		// mux.Use(Auth)
		// internal path, so it will be /admin/dashboard
		mux.Get("/dashboard", handlers.Repo.AdminDashboard)

		mux.Get("/reservations-new", handlers.Repo.AdminNewReservations)
		mux.Get("/reservations-all", handlers.Repo.AdminAllReservations)
		mux.Get("/reservations-calendar", handlers.Repo.AdminReservationsCalendar)
		mux.Get("/reservations/{src}/{id}", handlers.Repo.AdminShowReservation)
	})

	return mux
}
