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
var Cache cache.Cache

func IdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	id := vars["id"]

	res, _ := Cache.Get(id, Database.GetById)

	fmt.Fprintf(w, "Response: %s", res)
}
