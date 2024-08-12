package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/valkey-io/valkey-go"
)

type Connection struct {
	DB     *sql.DB
	Valkey *valkey.Client
}

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

func InitValkey(addresses []string, password string) *valkey.Client {
	client, err := valkey.NewClient(valkey.ClientOption{
		InitAddress: addresses,
		Password:    password,
	})
	if err != nil {
		log.Printf("Error connecting to valkey: %v\n", err)
		return nil
	}

	return &client
}
