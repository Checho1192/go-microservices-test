package data

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

const dbTimeout = time.Second * 3

var db *sql.DB

// New is the function used to create an instance of the data package. It returns the type
// Model, which embeds all the types we want to be available to our application.
func New(dbPool *sql.DB) Models {
	db = dbPool

	return Models{
		Measurement: Measurement{},
	}
}

// Models is the type for this package. Note that any model that is included as a member
// in this type is available to us throughout the application, anywhere that the
// app variable is used, provided that the model is also added in the New function.
type Models struct {
	Measurement Measurement
}

// Measurement is the structure which holds one measurement from the database.
type Measurement struct {
	ID                 string    `json:"id"`
	MeterID            int       `json:"meter_id"`
	ActiveEnergy       float64   `json:"active_energy"`
	ReactiveEnergy     float64   `json:"reactive_energy"`
	CapacitiveReactive float64   `json:"capacitive_reactive"`
	Solar              float64   `json:"solar"`
	Date               time.Time `json:"date"`
}

// GetById returns one measurement from the database, based on the id provided.
func (u *Measurement) GetById(id string) (*Measurement, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `SELECT id, meter_id, active_energy, reactive_energy, capacitive_reactive, solar, "date"
	FROM public.measurements
	WHERE id = $1;
	`

	var measurement Measurement
	row := db.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&measurement.ID,
		&measurement.MeterID,
		&measurement.ActiveEnergy,
		&measurement.ReactiveEnergy,
		&measurement.CapacitiveReactive,
		&measurement.Solar,
		&measurement.Date,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &measurement, nil
}
