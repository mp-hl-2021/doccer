package main

import (
	"doccer/api"
	"doccer/usecases"
	"net/http"
	"time"
)

func main() {
	service := api.NewApi(&usecases.SimpleUserInterface{})

	server := http.Server{
		Addr:         "localhost:8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      service.Router(),
	}
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}