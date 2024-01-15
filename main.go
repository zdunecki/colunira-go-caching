package main

import (
	// "fmt"
	"gosession/caching/cache"
	"gosession/caching/database"
	"gosession/caching/handlers"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	handlers.Database = database.Database{
		Data: map[string]interface{}{
			"123": "abc",
		},
	}

	expiresAfter, _ := time.ParseDuration("5s")
	handlers.Cache = cache.New(expiresAfter)

	router.HandleFunc("/{id}", handlers.IdHandler).Methods("GET")

	http.ListenAndServe(":8080", router)
}
