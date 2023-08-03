package dbrepo

import (
	"context"
	"time"

	"github.com/piotrj-rtbh/bookings/internal/models"
)

// we'll create any function that will be available to the interface repository.DatabaseRepo

func (m *postgresDBRepo) AllUsers() bool {
	return true
}

// InsertReservation inserts a reservation into the database
func (m *postgresDBRepo) InsertReservation(res models.Reservation) error {
	// in case the user loses his connection while a transaction is being executed in the backend
	// then make this connection timeout after 3 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into reservations (first_name, last_name, email, phone, start_date, 
					 end_date, room_id, created_at, updated_at)
					 values ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	// there is m.DB.Exec(stmt, params...) but with context version we assure we timeout
	// the connection when something wrong happens
	_, err := m.DB.ExecContext(ctx, stmt,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		res.RoomID,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return err
	}
	return nil
}
