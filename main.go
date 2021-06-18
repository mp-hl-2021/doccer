package main

import (
	"database/sql"
	"doccer/api"
	linter2 "doccer/linter"
	"doccer/model"
	storage2 "doccer/storage"
	"fmt"
	_ "github.com/lib/pq"
	"net/http"
	"time"
)

func main() {
	dbinfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		"db", "5432", "postgres", "qwerty", "postgres")
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		panic(err)
	}
	storage := storage2.PostgresStorage{
		Dbc: db,
	}
	//delete all data
	storage.ClearAllTables()

	linter := linter2.NewGeneralLinter()

	linter.RegisterNewLinter("Text", &linter2.StubLinter{})
	linter.RegisterNewLinter("go", &linter2.GoLinter{})

	m := model.NewModelImpl(&storage, []byte("abacaba"), linter, 10, 10)

	service := api.NewApi(&m)

	server := http.Server {
		Addr:         ":8080",
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