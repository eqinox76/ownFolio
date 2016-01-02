package main

import (
	"fmt"
	"net/http"
	"html/template"
	"time"
	"encoding/json"

	"appengine"
	"appengine/datastore"
	"appengine/user"

	"github.com/eqinox76/ownFolio/api"
)

func init() {
	http.HandleFunc("/stocks", getStock)
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
	c, _, login := checkLogin(w, r)
	if !login {
		http.Error(w, "Not logged in correctly", 401)
	}

	instrument := r.URL.Query().Get("instrument")

	timeseries := api.GetInstrument(c, instrument)

	jsData, err := json.Marshal(timeseries)
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	data := make(map[string]interface{})
	data["timeseriesdata"] = jsData

	err = chartTempl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}


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
