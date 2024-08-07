package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func InitSQL(connectionString string) *sql.DB {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatalf("Error connecting to database: %v\n", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Error pinging database: %v\n", err)
	} else {
		log.Println("Database connected")
	}

	return db
}
