package main

import (
	"doccer/api"
	"doccer/model"
	"doccer/usecases"
	"net/http"
	"time"
)

func main() {
	m := model.NewModelImpl(&usecases.SimpleStorage{}, []byte("abacaba"))
	service := api.NewApi(&m)

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