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
		MeasurementReport: MeasurementReport{},
	}
}

// Models is the type for this package. Note that any model that is included as a member
// in this type is available to us throughout the application, anywhere that the
// app variable is used, provided that the model is also added in the New function.
type Models struct {
	MeasurementReport MeasurementReport
}

// MeasurementReport is the structure which holds one measurement from the database.
type MeasurementReport struct {
	Period             time.Time `json:"period"`
	Active             float64   `json:"active"`
	ReactiveInductive  float64   `json:"reactive_inductive"`
	ReactiveCapacitive float64   `json:"reactive_capacitive"`
	Exported           float64   `json:"exported"`
}

// GetConsumptionFromMeterId in a given period of time
func (u *MeasurementReport) GetConsumptionFromMeterIdForPeriod(meterId int, startDate time.Time, endDate time.Time, period string) ([]MeasurementReport, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
	SELECT 
		DATE_TRUNC($1, date) AS period,
		SUM(active_energy) AS active,
		SUM(reactive_energy) AS reactive_inductive,
		SUM(capacitive_reactive) AS reactive_capacitive,
		SUM(solar) AS exported
	FROM 
		MEASUREMENTS
	WHERE 
		meter_id = $2 AND
		date BETWEEN $3 AND $4
	GROUP BY 
		DATE_TRUNC($1, date)
	ORDER BY 
		period ASC;
	`

	var measurements []MeasurementReport
	rows, err := db.QueryContext(ctx, query, period, meterId, startDate, endDate)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	for rows.Next() {
		var measurement MeasurementReport
		err := rows.Scan(
			&measurement.Period,
			&measurement.Active,
			&measurement.ReactiveInductive,
			&measurement.ReactiveCapacitive,
			&measurement.Exported,
		)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		measurements = append(measurements, measurement)
	}

	return measurements, nil
}
