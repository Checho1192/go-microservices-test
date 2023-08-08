package main

import (
	"errors"
	"fmt"
	"net/http"
)

func (app *Config) Consumption(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Consumption called")
	fmt.Println(r.URL.Query())
	fmt.Println(r.URL.Query().Get("id"))

	measurement, err := app.Models.Measurement.GetById(r.URL.Query().Get("id"))

	if err != nil {
		app.errorJSON(w, errors.New("measurement not found"), http.StatusNotFound)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "Measurement found",
		Data:    measurement,
	}

	app.writeJSON(w, http.StatusOK, payload)

}
