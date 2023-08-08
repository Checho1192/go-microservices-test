package main

import (
	"errors"
	"fmt"
	"net/http"
	"stats-ms/data"
	"strconv"
	"time"
)

type MeterResponse struct {
	MeterId            int       `json:"meter_id"`
	Address            string    `json:"address"`
	Active             []float64 `json:"active"`
	ReactiveInductive  []float64 `json:"reactive_inductive"`
	ReactiveCapacitive []float64 `json:"reactive_capacitive"`
	Exported           []float64 `json:"exported"`
}

type DataResponse struct {
	Period    []string        `json:"period"`
	DataGraph []MeterResponse `json:"data_graph"`
}

type DBResponse struct {
	MeterId            int `json:"meter_id"`
	MeasurementsReport []data.MeasurementReport
}

func (app *Config) Consumption(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Consumption called with Query: ", r.URL.Query())

	err := app.checkQueryParams(w, r, "meters_ids", "start_date", "end_date", "kind_period")

	if err != nil {
		return
	}
	if r.URL.Query().Get("kind_period") == "daily" {
		app.ConsumptionReport(w, r, "day")
		return
	} else if r.URL.Query().Get("kind_period") == "weekly" {
		app.ConsumptionReport(w, r, "week")
		return
	} else if r.URL.Query().Get("kind_period") == "monthly" {
		app.ConsumptionReport(w, r, "month")
		return
	}

	app.errorJSON(w, errors.New("invalid kind_period. Valid periods daily, weekly, monthly"), http.StatusBadRequest)
}

func (app *Config) ConsumptionReport(w http.ResponseWriter, r *http.Request, period string) {

	start_date, end_date, meter_ids, err := app.validationInput(w, r)
	if err != nil {
		return
	}

	DBResponses := make([]DBResponse, 0)

	for _, meter_id := range meter_ids {
		fmt.Println("Meter ID: ", meter_id)
		measurement_report, err := app.Models.MeasurementReport.GetConsumptionFromMeterIdForPeriod(meter_id, start_date, end_date, period)
		if err != nil {
			app.errorJSON(w, err, http.StatusBadRequest)
			return
		}
		DBResponses = append(DBResponses, DBResponse{MeterId: meter_id, MeasurementsReport: measurement_report})
	}

	app.writeJSON(w, http.StatusOK, app.createAnswerResponse(DBResponses))
}

func (app *Config) createAnswerResponse(DBResponses []DBResponse) (answer DataResponse) {
	periods := make([]string, 0)
	for _, measurement_report := range DBResponses[0].MeasurementsReport {
		periods = append(periods, measurement_report.Period.Format("2006-01-02"))
	}

	dataGraph := make([]MeterResponse, 0)
	for _, DBResponse := range DBResponses {
		active := make([]float64, 0)
		reactive_inductive := make([]float64, 0)
		reactive_capacitive := make([]float64, 0)
		exported := make([]float64, 0)
		for _, measurement_report := range DBResponse.MeasurementsReport {
			active = append(active, measurement_report.Active)
			reactive_inductive = append(reactive_inductive, measurement_report.ReactiveInductive)
			reactive_capacitive = append(reactive_capacitive, measurement_report.ReactiveCapacitive)
			exported = append(exported, measurement_report.Exported)
		}
		dataGraph = append(dataGraph, MeterResponse{
			MeterId:            DBResponse.MeterId,
			Address:            "Direccion Mock",
			Active:             active,
			ReactiveInductive:  reactive_inductive,
			ReactiveCapacitive: reactive_capacitive,
			Exported:           exported,
		})
	}

	answer = DataResponse{
		Period:    periods,
		DataGraph: dataGraph,
	}

	return answer
}

func (app *Config) validationInput(w http.ResponseWriter, r *http.Request) (start_date time.Time, end_date time.Time, meter_ids []int, err error) {
	start_date, err = time.Parse("2006-01-02", r.URL.Query().Get("start_date"))
	if err != nil {
		app.errorJSON(w, errors.New("incorrect date format. Date format must be YYYY-MM-DD"), http.StatusBadRequest)
		return start_date, end_date, meter_ids, err
	}

	end_date, err = time.Parse("2006-01-02", r.URL.Query().Get("end_date"))
	if err != nil {
		app.errorJSON(w, errors.New("incorrect date format. Date format must be YYYY-MM-DD"), http.StatusBadRequest)
		return start_date, end_date, meter_ids, err
	}

	if end_date.Before(start_date) {
		app.errorJSON(w, errors.New("end_date must be after start_date"), http.StatusBadRequest)
		return start_date, end_date, meter_ids, err
	}

	r.ParseForm()
	meter_ids_str := r.Form["meters_ids"]

	if len(meter_ids_str) == 0 {
		app.errorJSON(w, errors.New("missing query parameter: meters_ids"), http.StatusBadRequest)
		return start_date, end_date, meter_ids, err
	}

	//convert meters_ids to int slice
	meter_ids_int := make([]int, len(meter_ids_str))
	for i, v := range meter_ids_str {
		meter_ids_int[i], err = strconv.Atoi(v)
		if err != nil {
			app.errorJSON(w, errors.New("meters_ids must be integers"), http.StatusBadRequest)
			return start_date, end_date, meter_ids, err
		}
	}

	for _, meter_id := range meter_ids {
		fmt.Println("Meter ID: ", meter_id)

	}

	return start_date, end_date, meter_ids_int, err
}
