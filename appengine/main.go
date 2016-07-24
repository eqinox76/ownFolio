package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"appengine"
	"appengine/datastore"
	"appengine/user"

	"github.com/eqinox76/ownFolio/api"
	"github.com/eqinox76/ownFolio/data"
)

func init() {
	http.HandleFunc("/stocks", getStock)
	http.HandleFunc("/timeseries", getTimeSeries)
	http.HandleFunc("/holding", showHoldings)
	http.HandleFunc("/holding/add", addHolding)
	http.HandleFunc("/holding/get", getHolding)
}


var chartTempl = template.Must(template.ParseFiles("templates/base.html", "templates/chart.html"))

var holdingTempl = template.Must(template.ParseFiles("templates/base.html", "templates/holding.html"))


func getHolding(w http.ResponseWriter, r *http.Request) {
	c, _, login := checkLogin(w, r)
	if !login {
		http.Error(w, "Not logged in correctly", 401)
	}

	q := datastore.NewQuery("holding")
	var results []data.Holding
	_, err := q.GetAll(c, &results)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsData, err := json.Marshal(results)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	fmt.Fprintf(w, "%s", jsData)

}



func showHoldings(w http.ResponseWriter, r *http.Request) {
	_, _, login := checkLogin(w, r)
	if !login {
		http.Error(w, "Not logged in correctly", 401)
	}

	err := holdingTempl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// http://localhost:8080/holding/add?isin=%22huhuh%22&price=42.4&volume=51223&date=%222015-01-04%22
func addHolding(w http.ResponseWriter, r *http.Request) {
	c, _, login := checkLogin(w, r)
	if !login {
		http.Error(w, "Not logged in correctly", 401)
	}

	isin := strings.Trim(r.URL.Query().Get("isin"), "\"")
	priceStr := strings.Trim(r.URL.Query().Get("price"), "\"")
	dateStr := strings.Trim(r.URL.Query().Get("date"), "\"")
	volumeStr := strings.Trim(r.URL.Query().Get("volume"), "\"")

	if isin == "" || priceStr == "" || dateStr == "" || volumeStr == "" {
		fmt.Fprintf(w, "%v\n", r.URL.Query())
		http.Error(w, "Empty/missing parameters. Need isin, price, date, volume", http.StatusBadRequest)
		return
	}

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		http.Error(w, "Could not parse "+priceStr, http.StatusBadRequest)
		return
	}

	volume, err := strconv.ParseFloat(volumeStr, 64)
	if err != nil {
		http.Error(w, "Could not parse "+volumeStr, http.StatusBadRequest)
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		http.Error(w, "Could not parse '"+dateStr+"'\n"+err.Error(), http.StatusBadRequest)
		return
	}

	// we add a new symbol
	a := data.Holding{ISIN: isin, Price: price, Volume: volume, BuyDate: date}
	_, err = datastore.Put(c, datastore.NewIncompleteKey(c, "holding", nil), &a)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func getStock(w http.ResponseWriter, r *http.Request) {
	_, _, login := checkLogin(w, r)
	if !login {
		http.Error(w, "Not logged in correctly", 401)
	}

	err := chartTempl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getTimeSeries(w http.ResponseWriter, r *http.Request) {
	c, _, login := checkLogin(w, r)
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

func checkLogin(w http.ResponseWriter, r *http.Request) (appengine.Context, *user.User, bool) {
	c := appengine.NewContext(r)
	u := user.Current(c)

	// if we are not logged in lets try it
	if u == nil {
		url, err := user.LoginURL(c, r.URL.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return c, u, false
		}
		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusFound)
		return c, u, false
	}

	return c, u, true

}
