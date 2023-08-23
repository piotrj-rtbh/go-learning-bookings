package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/piotrj-rtbh/bookings/internal/config"
	"github.com/piotrj-rtbh/bookings/internal/driver"
	"github.com/piotrj-rtbh/bookings/internal/handlers"
	"github.com/piotrj-rtbh/bookings/internal/helpers"
	"github.com/piotrj-rtbh/bookings/internal/models"
	"github.com/piotrj-rtbh/bookings/internal/render"

	"github.com/alexedwards/scs/v2"
)

const portNumber = "localhost:8080"

var app config.AppConfig        // have to define here not in main() because middleware.go uses app.InProduction !
var session *scs.SessionManager // have to define global bc config.go will also use sessions!
var infoLog *log.Logger
var errorLog *log.Logger

// main is the main function
func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	// we're closing connection only when main() stops running!
	defer db.SQL.Close()
	// close mail channel only after main() stops running
	defer close(app.MailChan)

	fmt.Println("Starting mail listener...")
	listenForMail()

	/*
		// Just for testing email sending works
		msg := models.MailData{
			To:      "john@do.ca",
			From:    "me@here.com",
			Subject: "Some subject",
			Content: "",
		}
		app.MailChan <- msg
	*/

	// sending an email example (does not work now)
	/* from := "me@here.com"
	auth := smtp.PlainAuth("", from, "", "localhost")
	err = smtp.SendMail("localhost:1025", auth, from, []string{"you@there.com"}, []byte("Hello, world"))
	if err != nil {
		log.Println(err)
	} */

	fmt.Println(fmt.Sprintf("Staring application on port %s", portNumber))
	// _ = http.ListenAndServe(portNumber, nil)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}

func run() (*driver.DB, error) {
	// what am I going to put in the session?
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})
	gob.Register(map[string]int{}) // we store the block maps in session as well

	// creating mail channel - ideal tool for sending emails from the app (it's async)
	mailChan := make(chan models.MailData)
	app.MailChan = mailChan

	// change this to true when in production
	app.InProduction = false

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
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

	// connect to database
	log.Println("Connecting to database...")
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=bookings user=postgres password=postgres")
	if err != nil {
		log.Fatal("Cannot connect to database! Dying...")
	}

	// defer db.SQL.Close() // we can't have closing here because run() is called from main() and once run() finishes the DB conn will be closed

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
		return nil, err
	}

	app.TemplateCache = tc
	app.UseCache = false

	// we also MUST allow the db connection to be accessible by handlers
	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)
	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	return db, nil
}
