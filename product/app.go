package main

import (
	"database/sql"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *App) Initialize(connectionString string) {
	var err error

	a.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		panic(err)
	}

	a.Router = mux.NewRouter()
}

func (a *App) Run(port string) {}
