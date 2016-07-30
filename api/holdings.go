package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"appengine"
	"appengine/datastore"

	"github.com/eqinox76/ownFolio/data"
)

type holdingAncestor struct {
	Key string
}

// we need a ancestor for consistent queries
var ancestor *datastore.Key = nil

func maybeCheckAncestor(w http.ResponseWriter, c appengine.Context) {
	if ancestor != nil {
		return
	}

	q := datastore.NewQuery("holdingAncestor").Limit(1).KeysOnly()
	count, err := q.Count(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if count == 0 {

		ancestor, err = datastore.Put(c, datastore.NewIncompleteKey(c, "holdingAncestor", nil), &holdingAncestor{Key: "ancestor"})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		ancestor, err = q.Run(c).Next(nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func GetHolding(w http.ResponseWriter, r *http.Request) {
	c, _, login := CheckLogin(w, r)
	if !login {
		http.Error(w, "Not logged in correctly", 401)
	}

	maybeCheckAncestor(w, c)

	q := datastore.NewQuery("holding").Ancestor(ancestor)

	var results []data.Holding
	keys, err := q.GetAll(c, &results)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// add keys to the holding data
	for i, _ := range results {
		results[i].Key = keys[i].Encode()
	}

	jsData, err := json.Marshal(results)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	fmt.Fprintf(w, "%s", jsData)

}

// http://localhost:8080/holding/add?isin=%22huhuh%22&price=42.4&volume=51223&date=%222015-01-04%22
func AddHolding(w http.ResponseWriter, r *http.Request) {
	c, _, login := CheckLogin(w, r)
	if !login {
		http.Error(w, "Not logged in correctly", 401)
	}

	maybeCheckAncestor(w, c)

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
	_, err = datastore.Put(c, datastore.NewIncompleteKey(c, "holding", ancestor), &a)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func DelHolding(w http.ResponseWriter, r *http.Request) {
	c, _, login := CheckLogin(w, r)
	if !login {
		http.Error(w, "Not logged in correctly", 401)
	}

	maybeCheckAncestor(w, c)

	key := strings.Trim(r.URL.Query().Get("key"), "\"")

	datastoreKey, err := datastore.DecodeKey(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = datastore.Delete(c, datastoreKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
