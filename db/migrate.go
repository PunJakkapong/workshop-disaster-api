package db

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"sort"
	"strings"
)

// RunMigrations executes all SQL migration files in the migrations directory
func RunMigrations(db *sql.DB) error {
	// Get all migration files
	files, err := ioutil.ReadDir("db/migrations")
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %v", err)
	}

	// Create migrations table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS migrations (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL UNIQUE,
			applied_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %v", err)
	}

	// Sort files by name
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	// Run each migration
	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".sql") {
			continue
		}

		// Check if migration has been applied
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM migrations WHERE name = $1", file.Name()).Scan(&count)
		if err != nil {
			return fmt.Errorf("failed to check migration status: %v", err)
		}

		if count > 0 {
			log.Printf("Skipping migration %s (already applied)", file.Name())
			continue
		}

		// Read and execute migration
		content, err := ioutil.ReadFile(filepath.Join("db/migrations", file.Name()))
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %v", file.Name(), err)
		}

		// Start transaction
		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("failed to begin transaction: %v", err)
		}

		// Execute migration
		_, err = tx.Exec(string(content))
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to execute migration %s: %v", file.Name(), err)
		}

		// Record migration
		_, err = tx.Exec("INSERT INTO migrations (name) VALUES ($1)", file.Name())
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to record migration %s: %v", file.Name(), err)
		}

		// Commit transaction
		if err = tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit migration %s: %v", file.Name(), err)
		}

		log.Printf("Applied migration %s", file.Name())
	}

	return nil
}
