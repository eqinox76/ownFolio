package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/eqinox76/ownFolio/api"
	"github.com/eqinox76/ownFolio/data"
)

func init() {
	http.HandleFunc("/stocks", getStock)
	http.HandleFunc("/timeseries", getTimeSeries)
	http.HandleFunc("/holding", showHoldings)
	http.HandleFunc("/holding/add", api.AddHolding)
	http.HandleFunc("/holding/get", api.GetHolding)
}

var chartTempl = template.Must(template.ParseFiles("templates/base.html", "templates/chart.html"))

var holdingTempl = template.Must(template.ParseFiles("templates/base.html", "templates/holding.html"))

func showHoldings(w http.ResponseWriter, r *http.Request) {
	_, _, login := api.CheckLogin(w, r)
	if !login {
		http.Error(w, "Not logged in correctly", 401)
	}

	err := holdingTempl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getStock(w http.ResponseWriter, r *http.Request) {
	_, _, login := api.CheckLogin(w, r)
	if !login {
		http.Error(w, "Not logged in correctly", 401)
	}

	err := chartTempl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getTimeSeries(w http.ResponseWriter, r *http.Request) {
	c, _, login := api.CheckLogin(w, r)
	if !login {
		http.Error(w, "Not logged in correctly", 401)
	}

	instrument := r.URL.Query().Get("instrument")

	//TODO find a better way to do this optional limit parameter
	var timeseries data.DataPoints

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		timeseries = api.GetInstrument(c, instrument)
	} else {
		timeseries = api.GetInstrumentLimited(c, instrument, limit)
	}

	jsData, err := json.Marshal(timeseries)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	fmt.Fprintf(w, "%s", jsData)
}
