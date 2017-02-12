package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/eqinox76/ownFolio/api"
	"github.com/eqinox76/ownFolio/api/holdings"
	"github.com/eqinox76/ownFolio/api/isinresolver"
)

func init() {
	/*
	* js start pages
	 */
	http.HandleFunc("/stocks", getStock)
	http.HandleFunc("/holding", showHoldings)

	/*
	* REST api
	 */

	// access market data
	http.HandleFunc("/timeseries", getTimeSeries)

	// manage what the logged in user owns
	http.HandleFunc("/holding/add", api.WithDatastore(holdings.Add))
	http.HandleFunc("/holding/get", api.WithDatastore(holdings.Get))
	http.HandleFunc("/holding/del", api.DelFromDataStore)

	// manage retrieval options
	http.HandleFunc("/isinresolver/add", api.WithDatastore(isinresolver.Add))
	http.HandleFunc("/isinresolver/get", api.WithDatastore(isinresolver.Get))
	http.HandleFunc("/isinresolver/del", api.DelFromDataStore)
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
	var timeseries []byte
	var generateError error

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		timeseries, generateError = api.GetInstrument(c, instrument)
	} else {
		timeseries, generateError = api.GetInstrumentLimited(c, instrument, limit)
	}

	if generateError != nil {
		http.Error(w, generateError.Error(), 500)
	}

	fmt.Fprintf(w, "%s", timeseries)
}
