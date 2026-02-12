package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func NewPgDatabase() (*sql.DB, error) {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, fmt.Errorf("DB connection failed with error: %s" + err.Error())
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("DB Ping failed with error: %s" + err.Error())
	}

	fmt.Println("Connected To DB")
	return db, nil
}
