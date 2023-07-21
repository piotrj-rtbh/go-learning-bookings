package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/piotrj-rtbh/bookings/pkg/config"
	"github.com/piotrj-rtbh/bookings/pkg/handlers"
	"github.com/piotrj-rtbh/bookings/pkg/render"

	"github.com/alexedwards/scs/v2"
)

const portNumber = "localhost:8080"

var app config.AppConfig        // have to define here not in main() because middleware.go uses app.InProduction !
var session *scs.SessionManager // have to define global bc config.go will also use sessions!

// main is the main function
func main() {

	// change this to true when in production
	app.InProduction = false

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

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)
	render.NewTemplates(&app)

	fmt.Println(fmt.Sprintf("Staring application on port %s", portNumber))
	// _ = http.ListenAndServe(portNumber, nil)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}
