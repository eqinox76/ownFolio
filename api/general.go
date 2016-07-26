package api

import (
	"net/http"

	"appengine"
	"appengine/user"
)

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
