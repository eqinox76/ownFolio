package portfolio

import (
	"fmt"
	"net/http"
	"time"

	"appengine"
	"appengine/datastore"
	"appengine/user"
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

func getStock(w http.ResponseWriter, r *http.Request) {
	c, _, login := checkLogin(w, r)
	if !login {
		http.Error(w, "Not logged in correctly", 401)
	}

	stock := &OwnedStock{
		Name: "test",
	}

	fmt.Fprintf(w, "Stored and retrieved the Employee named %q", stock)

	key, err := datastore.Put(c, datastore.NewIncompleteKey(c, "employee", nil), stock)
	fmt.Fprintf(w, "%q\n", key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	stock2 := new(OwnedStock)
	err = datastore.Get(c, key, stock2)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Stored and retrieved the Employee named %q", stock2)
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
