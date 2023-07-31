package render

import (
	"net/http"
	"testing"

	"github.com/piotrj-rtbh/bookings/internal/models"
)

func TestAddDefaultData(t *testing.T) {
	var td models.TemplateData

	// r, err := http.NewRequest("GET", "/some-url", nil)
	// if err != nil {
	// 	t.Error(err)
	// }
	// ^^ above code we replace with:
	r, err := getSession()
	if err != nil {
		t.Error(err)
	}
	// let's store something in the session
	session.Put(r.Context(), "flash", "123")

	result := AddDefaultData(&td, r)

	if result.Flash != "123" {
		t.Error("flash value of 123 not found in session")
	}
}

// there's no context for http.Request so we have a request that has session data
func getSession() (*http.Request, error) {
	// move the code from TestAddDefaultData
	r, err := http.NewRequest("GET", "/some-url", nil)
	if err != nil {
		return nil, err
	}

	// create a context
	ctx := r.Context()
	ctx, _ = session.Load(ctx, r.Header.Get("X-Session"))
	r = r.WithContext(ctx)

	return r, nil
}
