package api

import (
	"net/http"

	"appengine"
	"appengine/datastore"
	"appengine/user"
)

type holdingAncestor struct {
	Key string
}

type HandlerWithDataStore func(http.ResponseWriter, *http.Request, appengine.Context)

var Ancestor *datastore.Key = nil

func CheckLogin(w http.ResponseWriter, r *http.Request) (appengine.Context, *user.User, bool) {
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

// we need a ancestor for consistent queries therefore we must make sure that it is in the datastore
func maybeCheckAncestor(w http.ResponseWriter, c appengine.Context) {
	if Ancestor != nil {
		return
	}

	q := datastore.NewQuery("holdingAncestor").Limit(1).KeysOnly()
	count, err := q.Count(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if count == 0 {

		Ancestor, err = datastore.Put(c, datastore.NewIncompleteKey(c, "holdingAncestor", nil), &holdingAncestor{Key: "ancestor"})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		Ancestor, err = q.Run(c).Next(nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func WithDatastore(h HandlerWithDataStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, _, login := CheckLogin(w, r)
		if !login {
			http.Error(w, "Not logged in correctly", 401)
		}

		maybeCheckAncestor(w, c)

		h(w, r, c)
	}
}
