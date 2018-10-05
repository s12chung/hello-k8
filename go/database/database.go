package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq" // Database driver
)

// DefaultDataBaseConfig returns a String Parameters for the default
// database config used. It follows the postgres docker image.
func DefaultDataBaseConfig() string {
	user := os.Getenv("POSTGRES_USER")
	if user == "" {
		user = "postgres"
	}
	dbname := os.Getenv("POSTGRES_DB")
	if dbname == "" {
		dbname = user
	}
	password := os.Getenv("POSTGRES_PASSWORD")
	host := os.Getenv("POSTGRES_SERVICE_HOST")
	port := os.Getenv("POSTGRES_SERVICE_PORT")

	return fmt.Sprintf("user=%v password=%v dbname=%v host=%v port=%v sslmode=disable", user, password, dbname, host, port)
}

// DefaultDataBase returns the DB given DefaultDataBaseConfig
func DefaultDataBase() (*sql.DB, error) {
	db, err := sql.Open("postgres", DefaultDataBaseConfig())
	if err != nil {
		return nil, err
	}
	return db, nil
}
