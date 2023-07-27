package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/piotrj-rtbh/bookings/internal/config"
	"github.com/piotrj-rtbh/bookings/internal/forms"
	"github.com/piotrj-rtbh/bookings/internal/models"
	"github.com/piotrj-rtbh/bookings/internal/render"
)

// Repo the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the handler for the home page
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	// we'd like to grap the user's IP address and store it in the home page (in the session)
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.RenderTemplate(w, r, "home.page.tmpl", &models.TemplateData{})
}

// About is the handler for the about page
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	// perform some logic
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello, again"

	// use session
	// m.App.Session ...
	// pull out the IP from the session stored when user visited home page
	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP

	// send data to the template
	render.RenderTemplate(w, r, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}

// Reservation renders the make a reservation page and displays form
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	var emptyReservation models.Reservation
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation

	render.RenderTemplate(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		// for GET we make the Form object be empty
		Form: forms.New(nil),
		Data: data,
	})
}

// PostReservation handles the posting of a reservation form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}

	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Phone:     r.Form.Get("phone"),
		Email:     r.Form.Get("email"),
	}

	// now we have to check our data
	form := forms.New(r.PostForm) // r.PostForm is accessible only after ParseForm has been called

	//form.Has("first_name", r) // adds errors if failed
	// instead we use our shiny new function form.Required
	form.Required("first_name", "last_name", "email")
	// and some more validators as well :)
	form.MinLength("first_name", 3, r)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{}) // create empty data for TemplateData.Data
		// store the reservation as it is so that we can repost those data back to the user
		// because the data is not valid and the user has to correct data
		data["reservation"] = reservation
		// now render the form
		render.RenderTemplate(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			// for GET we make the Form object be empty
			Form: form,
			Data: data,
		})
		// stop processing further
		return
	}
}

// Generals renders the room page
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "generals.page.tmpl", &models.TemplateData{})
}

// Majors renders the room page
func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "majors.page.tmpl", &models.TemplateData{})
}

// Availability renders the search availability page
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

// PostAvailability renders the search availability page
// handler for data being POST-ed, so we change the body of this method
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	// let's do something with the values POST-ed to
	start := r.Form.Get("start")
	end := r.Form.Get("end")
	w.Write([]byte(fmt.Sprintf("Start date is %s and end date is %s", start, end)))
}

/*
the response will b in a form of

	{
		"ok": true/false,
		"message": "some message"
	}

We need to have a struct that reflects this JSON
*/
type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

// AvailabilityJSON handles request for availability and send JSON response
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	resp := jsonResponse{
		OK:      false,
		Message: "Available!",
	}

	// create JSON string to be written out - marshalling a struct into a string containing JSON
	out, err := json.MarshalIndent(resp, "", "     ")
	if err != nil {
		log.Println(err)
	}

	// just to monitor what's sent
	log.Println(string(out))

	// write headers
	w.Header().Set("Content-Type", "application/json")
	// write body
	w.Write(out)
}

// Contact renders the contact page
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "contact.page.tmpl", &models.TemplateData{})
}
