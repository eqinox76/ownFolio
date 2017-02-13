package isinresolver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"appengine"
	"appengine/datastore"

	"github.com/eqinox76/ownFolio/api"
	"github.com/eqinox76/ownFolio/data"
)

func Get(w http.ResponseWriter, r *http.Request, c appengine.Context) {

	q := datastore.NewQuery("isintranslation").Ancestor(api.Ancestor)

	var results []data.IsinTranslation
	keys, err := q.GetAll(c, &results)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// add keys to the holding data to be able to delete them later
	for i, _ := range results {
		results[i].Key = keys[i].Encode()
	}

	jsData, err := json.Marshal(results)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	fmt.Fprintf(w, "%s", jsData)

}

func Add(w http.ResponseWriter, r *http.Request, c appengine.Context) {

	isin := strings.Trim(r.URL.Query().Get("isin"), "\"")
	source := strings.Trim(r.URL.Query().Get("source"), "\"")
	identifier := strings.Trim(r.URL.Query().Get("identifier"), "\"")
	database := strings.Trim(r.URL.Query().Get("database"), "\"")

	if isin == "" || source == "" || identifier == "" || database == "" {
		fmt.Fprintf(w, "%v\n", r.URL.Query())
		http.Error(w, "Empty/missing parameters. Need isin, source, identifier, database", http.StatusBadRequest)
		return
	}

	a := data.IsinTranslation{ISIN: isin, Source: source, Identifier: identifier, Database: database}
	_, err := datastore.Put(c, datastore.NewIncompleteKey(c, "isintranslation", api.Ancestor), &a)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
