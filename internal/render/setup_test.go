package render

import (
	"encoding/gob"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/piotrj-rtbh/bookings/internal/config"
	"github.com/piotrj-rtbh/bookings/internal/models"
)

var session *scs.SessionManager
var testApp config.AppConfig // create our own copy of the app

func TestMain(m *testing.M) {
	// copy of main.go run()
	// but have to change app to our testApp

	// what am I going to put in the session?
	gob.Register(models.Reservation{})

	// change this to true when in production
	testApp.InProduction = false

	// create a session using scs package
	session = scs.New()
	// we'd like the session to live for a defined range of time
	session.Lifetime = 24 * time.Hour
	// by default it stores sessions in cookies but that may be changed
	session.Cookie.Persist = true                  // false only for session of web browser being opened
	session.Cookie.SameSite = http.SameSiteLaxMode // how strict you want to be about what site this cookie applies to
	session.Cookie.Secure = false                  // we have to set it hard to FALSE for tests

	// store this session in globally accessible app (from config.go)
	testApp.Session = session

	app = &testApp // we assign the global app to be our testApp

	os.Exit(m.Run())
}

// let's create a test ResponseWriter which satisfies the original ResponseWriter interface (look at the description of Response Writer)
type myWriter struct{}

func (tw *myWriter) Header() http.Header {
	var h http.Header
	return h
}

func (tw *myWriter) WriteHeader(i int) {

}

func (tw *myWriter) Write(b []byte) (int, error) {
	length := len(b)
	return length, nil
}
