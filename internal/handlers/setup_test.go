package handlers

import (
	"encoding/gob"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/justinas/nosurf"
	"github.com/piotrj-rtbh/bookings/internal/config"
	"github.com/piotrj-rtbh/bookings/internal/models"
	"github.com/piotrj-rtbh/bookings/internal/render"
)

var app config.AppConfig
var session *scs.SessionManager
var pathToTemplates = "./../../templates"
var functions = template.FuncMap{}

func getRoutes() http.Handler {

	// copy of main.go
	gob.Register(models.Reservation{})

	// change this to true when in production
	app.InProduction = false

	// we have to copy error logging from main.go here
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	// create a session using scs package
	session = scs.New()
	// we'd like the session to live for a defined range of time
	session.Lifetime = 24 * time.Hour
	// by default it stores sessions in cookies but that may be changed
	session.Cookie.Persist = true                  // false only for session of web browser being opened
	session.Cookie.SameSite = http.SameSiteLaxMode // how strict you want to be about what site this cookie applies to
	session.Cookie.Secure = app.InProduction

	// store this session in globally accessible app (from config.go)
	app.Session = session

	tc, err := CreateTestTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")

	}

	app.TemplateCache = tc
	// IMPORTANT CHANGE: MUST CHANGE THIS TO true
	// https://www.udemy.com/course/building-modern-web-applications-with-go/learn/lecture/23069034 9:40
	app.UseCache = true

	repo := NewRepo(&app)
	NewHandlers(repo)
	render.NewTemplates(&app)

	// copy of func routes() from routes.go
	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer) // prevents from panic errors! middleware comes from chi
	// IMPORTANT: WE COMMENT IT RIGHT NOW FOR A MOMENT
	// JUST TO ALLOW POST-ing with tests
	// THIS WILL BE CHANGED IN THE NEXT LECTURES
	// mux.Use(NoSurf)               // this will prevent POST requests without proper CSRF token
	mux.Use(SessionLoad)

	// defining paths
	mux.Get("/", Repo.Home)
	mux.Get("/about", Repo.About)
	mux.Get("/generals-quarters", Repo.Generals)
	mux.Get("/majors-suite", Repo.Majors)
	mux.Get("/search-availability", Repo.Availability)
	mux.Post("/search-availability", Repo.PostAvailability)
	mux.Post("/search-availability-json", Repo.AvailabilityJSON)

	mux.Get("/contact", Repo.Contact)

	mux.Get("/make-reservation", Repo.Reservation)
	mux.Post("/make-reservation", Repo.PostReservation)
	mux.Get("/reservation-summary", Repo.ReservationSummary)

	// in order to enable images loading we have to run a file server
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}

// copy of middleware.go because we can't import them from there
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

func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next) // globally defined session in main.go
}

// CreateTestTemplateCache creates a template cache as a map
// it's a copy of CreateTemplateCache from render.go
func CreateTestTemplateCache() (map[string]*template.Template, error) {

	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}

	return myCache, nil
}
