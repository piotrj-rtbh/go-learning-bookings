package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/piotrj-rtbh/bookings/internal/config"
	"github.com/piotrj-rtbh/bookings/internal/driver"
	"github.com/piotrj-rtbh/bookings/internal/forms"
	"github.com/piotrj-rtbh/bookings/internal/helpers"
	"github.com/piotrj-rtbh/bookings/internal/models"
	"github.com/piotrj-rtbh/bookings/internal/render"
	"github.com/piotrj-rtbh/bookings/internal/repository"
	"github.com/piotrj-rtbh/bookings/internal/repository/dbrepo"
)

// Repo the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the handler for the home page
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	// we'd like to grap the user's IP address and store it in the home page (in the session)
	render.Template(w, r, "home.page.tmpl", &models.TemplateData{})
}

// About is the handler for the about page
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	// send data to the template
	render.Template(w, r, "about.page.tmpl", &models.TemplateData{})
}

// Reservation renders the make a reservation page and displays form
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	var emptyReservation models.Reservation
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation

	render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		// for GET we make the Form object be empty
		Form: forms.New(nil),
		Data: data,
	})
}

// PostReservation handles the posting of a reservation form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	// err = errors.New("this is an error message") // <-- uncomment it to see the error console logs from the server
	// when trying to /make-reservation
	if err != nil {
		// log.Println(err) // <-- replace with our error handling package now!
		helpers.ServerError(w, err)
		return
	}

	// Formatting and parsing dates in Go:
	// https://www.pauladamsmith.com/blog/2011/05/go_time.html
	sd := r.Form.Get("start_date")
	ed := r.Form.Get("end_date")
	// 2020-01-01 -- 01/02 03:04:05PM '06 -0700

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, sd)
	if err != nil {
		helpers.ServerError(w, err)
	}
	endDate, err := time.Parse(layout, ed)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	roomID, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Phone:     r.Form.Get("phone"),
		Email:     r.Form.Get("email"),
		StartDate: startDate,
		EndDate:   endDate,
		RoomID:    roomID,
	}

	// now we have to check our data
	form := forms.New(r.PostForm) // r.PostForm is accessible only after ParseForm has been called

	//form.Has("first_name", r) // adds errors if failed
	// instead we use our shiny new function form.Required
	form.Required("first_name", "last_name", "email")
	// and some more validators as well :)
	form.MinLength("first_name", 3)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{}) // create empty data for TemplateData.Data
		// store the reservation as it is so that we can repost those data back to the user
		// because the data is not valid and the user has to correct data
		data["reservation"] = reservation
		// now render the form
		render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			// for GET we make the Form object be empty
			Form: form,
			Data: data,
		})
		// stop processing further
		return
	}

	// write the information to the DB
	newReservationID, err := m.DB.InsertReservation(reservation)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// we have to put room restrictions as well!
	restriction := models.RoomRestriction{
		StartDate:     startDate,
		EndDate:       endDate,
		RoomID:        roomID,
		ReservationID: newReservationID,
		RestrictionID: 1,
	}

	err = m.DB.InsertRoomRestriction(restriction)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// store the reservation object in the session
	m.App.Session.Put(r.Context(), "reservation", reservation)

	// we don't want users to submit the form twice! So we redirect to summary
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

// Generals renders the room page
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "generals.page.tmpl", &models.TemplateData{})
}

// Majors renders the room page
func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "majors.page.tmpl", &models.TemplateData{})
}

// Availability renders the search availability page
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

// PostAvailability renders the search availability page
// handler for data being POST-ed, so we change the body of this method
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	// let's do something with the values POST-ed to
	start := r.Form.Get("start")
	end := r.Form.Get("end")

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, start)
	if err != nil {
		helpers.ServerError(w, err)
	}
	endDate, err := time.Parse(layout, end)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	rooms, err := m.DB.SearchAvailabilityForAllRooms(startDate, endDate)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	/* for _, i := range rooms {
		m.App.InfoLog.Println("ROOM:", i.ID, i.RoomName)
	} */

	if len(rooms) == 0 {
		// no availability
		// m.App.InfoLog.Println("NO Avail")
		m.App.Session.Put(r.Context(), "error", "No availability")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})
	data["rooms"] = rooms

	res := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}

	m.App.Session.Put(r.Context(), "reservation", res)

	render.Template(w, r, "choose-room.page.tmpl", &models.TemplateData{
		Data: data,
	})

	// w.Write([]byte(fmt.Sprintf("Start date is %s and end date is %s", start, end)))
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
		// log.Println(err)
		helpers.ServerError(w, err)
		return
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
	render.Template(w, r, "contact.page.tmpl", &models.TemplateData{})
}

func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		// log.Println("cannot get item from session")
		m.App.ErrorLog.Println("Can't get error from session")
		// pass a value into session:
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		// we also need the error to be displayed

		return
	}
	// now we can get rid of the data stored in the session!
	m.App.Session.Remove(r.Context(), "reservation")

	// display the summary page
	data := make(map[string]interface{})
	data["reservation"] = reservation

	render.Template(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data: data,
	})
}
