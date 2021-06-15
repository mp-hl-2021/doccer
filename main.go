package main

import (
	"database/sql"
	"doccer/api"
	"doccer/model"
	storage2 "doccer/storage"
	"fmt"
	_ "github.com/lib/pq"
	"net/http"
	"time"
)

func main() {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		"ilya", "qwerty", "hl_systems")
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		panic(err)
	}
	storage := storage2.PostgresStorage{
		Dbc: db,
	}

	//delete all data
	storage.ClearAllTables()

	m := model.NewModelImpl(&storage, []byte("abacaba"))
	service := api.NewApi(&m)

	server := http.Server{
		Addr:         "localhost:8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      service.Router(),
	}
	println("Starting server")
	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}