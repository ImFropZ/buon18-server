package database

import (
	"database/sql"
	"log"
	"server/config"

	_ "github.com/lib/pq"
)

func InitSQL() *sql.DB {
	connectionString := config.Env("DB_CONNECTION_STRING")

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
