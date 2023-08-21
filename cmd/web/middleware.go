package main

import (
	"net/http"

	"github.com/justinas/nosurf"
	"github.com/piotrj-rtbh/bookings/internal/helpers"
)

// func WriteToConsole(next http.Handler) http.Handler {

// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		fmt.Println("Hit the page")
// 		next.ServeHTTP(w, r)
// 	})
// }

// CSRFToken generator as a middleware: https://github.com/justinas/nosurf
// NoSurf adds CSRF protection to all POST requests
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})

	return csrfHandler
}

// we have to remember the user in some way (remember state)
// SessionLoad loads and saves sessions on every request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next) // globally defined session in main.go
}

// Auth ensures the user is authenticated
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !helpers.IsAuthenticated(r) {
			session.Put(r.Context(), "error", "Log in first!")
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}
