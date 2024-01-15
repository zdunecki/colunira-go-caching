package handlers

import (
	"fmt"
	"gosession/caching/cache"
	"gosession/caching/database"
	"net/http"

	"github.com/gorilla/mux"
)

// TODO: Better way to make this "global"?
var Database database.Database
var Cache cache.CacheInterface

func IdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	id := vars["id"]

	res, found := Cache.Get(id)
	if !found {
		res, _ = Database.GetById(id)
		Cache.Set(id, res)
	}
	
	fmt.Fprintf(w, "Response: %s", res)
}
