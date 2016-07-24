package main

import (
	"fmt"
	"net/http"
	"html/template"
	"time"
	"encoding/json"
	"strconv"

	"appengine"
	"appengine/datastore"
	"appengine/user"

	"github.com/eqinox76/ownFolio/api"
	"github.com/eqinox76/ownFolio/data"
)

func init() {
	http.HandleFunc("/stocks", getStock)
	http.HandleFunc("/timeseries", getTimeSeries)
	http.HandleFunc("/observe", showObserved)
}

type Observed struct {
	Name string
	Id   string
}

func showObserved(w http.ResponseWriter, r *http.Request) {
	c, _, login := checkLogin(w, r)
	if !login {
		http.Error(w, "Not logged in correctly", 401)
	}

	add := r.URL.Query().Get("add")
	if len(add) != 0 {
		// we add a new symbol
		a := Observed{Id: add}
		_, err := datastore.Put(c, datastore.NewIncompleteKey(c, "observed", nil), &a)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		q := datastore.NewQuery("observed")
		var results []Observed
		_, err := q.GetAll(c, &results)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		fmt.Fprintf(w, "%v elements:\n", len(results))
		fmt.Fprintf(w, "%v", results)
	}
}

type OwnedStock struct {
	Name    string
	BuyDate time.Time
}

var chartTempl = template.Must(template.ParseFiles("templates/base.html", "templates/chart.html"))

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

func getTimeSeries(w http.ResponseWriter, r *http.Request){
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
	if err != nil{
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
