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
	http.HandleFunc("/stocks", api.LoggedIn(getStock))
	http.HandleFunc("/holding", api.LoggedIn(showHoldings))
	http.HandleFunc("/isinresolver", api.LoggedIn(showIsinResolver))

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

func showHoldings(w http.ResponseWriter, r *http.Request) {
	var tmpl = template.Must(template.ParseFiles("templates/base.html", "templates/holding.html"))

	err := tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getStock(w http.ResponseWriter, r *http.Request) {
	var tmpl = template.Must(template.ParseFiles("templates/base.html", "templates/chart.html"))

	err := tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func showIsinResolver(w http.ResponseWriter, r *http.Request) {
	var tmpl = template.Must(template.ParseFiles("templates/base.html", "templates/isinresolver.html"))

	err := tmpl.Execute(w, nil)
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
