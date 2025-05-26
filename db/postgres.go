package db

import (
	"database/sql"
	"fmt"
	"os"

	// Import postgres driver
	_ "github.com/lib/pq"
)

// createDatabase creates the database if it doesn't exist
func createDatabase() error {
	// Connect to default postgres database
	postgresURL := fmt.Sprintf("postgres://%s:%s@%s:%s/postgres?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"))

	db, err := sql.Open("postgres", postgresURL)
	if err != nil {
		return fmt.Errorf("failed to connect to postgres: %v", err)
	}
	defer db.Close()

	// Check if database exists
	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT datname FROM pg_catalog.pg_database WHERE datname = $1)",
		os.Getenv("POSTGRES_DB")).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check if database exists: %v", err)
	}

	// Create database if it doesn't exist
	if !exists {
		_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", os.Getenv("POSTGRES_DB")))
		if err != nil {
			return fmt.Errorf("failed to create database: %v", err)
		}
		fmt.Printf("Created database %s\n", os.Getenv("POSTGRES_DB"))
	}

	return nil
}

// ConnectPostgres establishes a connection to PostgreSQL database
func ConnectPostgres() (*sql.DB, error) {
	// Create database if it doesn't exist
	if err := createDatabase(); err != nil {
		return nil, err
	}

	postgresURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"))

	db, err := sql.Open("postgres", postgresURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %v", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping PostgreSQL: %v", err)
	}

	return db, nil
}
