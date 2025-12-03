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

// RunMigrations executes all SQL migration files in order
func RunMigrations(db *sql.DB, migrationsPath string) error {
	log.Println("🔄 Running database migrations...")

	// Create migrations tracking table if not exists
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Read all migration files
	files, err := ioutil.ReadDir(migrationsPath)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	// Filter and sort SQL files
	var migrations []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			migrations = append(migrations, file.Name())
		}
	}
	sort.Strings(migrations)

	if len(migrations) == 0 {
		log.Println("⚠️  No migration files found")
		return nil
	}

	// Execute each migration
	appliedCount := 0
	skippedCount := 0

	for _, migration := range migrations {
		// Check if already applied
		var exists bool
		err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version = $1)", migration).Scan(&exists)
		if err != nil {
			return fmt.Errorf("failed to check migration status for %s: %w", migration, err)
		}

		if exists {
			skippedCount++
			continue
		}

		// Read migration file
		filePath := filepath.Join(migrationsPath, migration)
		content, err := ioutil.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", migration, err)
		}

		// Execute migration
		log.Printf("  ▶ Applying migration: %s", migration)
		_, err = db.Exec(string(content))
		if err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", migration, err)
		}

		// Record migration
		_, err = db.Exec("INSERT INTO schema_migrations (version) VALUES ($1)", migration)
		if err != nil {
			return fmt.Errorf("failed to record migration %s: %w", migration, err)
		}

		appliedCount++
		log.Printf("  ✅ Applied: %s", migration)
	}

	log.Printf("✅ Migrations complete: %d applied, %d skipped, %d total", appliedCount, skippedCount, len(migrations))
	return nil
}

// GetAppliedMigrations returns list of applied migrations
func GetAppliedMigrations(db *sql.DB) ([]string, error) {
	rows, err := db.Query("SELECT version FROM schema_migrations ORDER BY version")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var migrations []string
	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		migrations = append(migrations, version)
	}

	return migrations, nil
}

// GetPendingMigrations returns list of migrations that haven't been applied
func GetPendingMigrations(db *sql.DB, migrationsPath string) ([]string, error) {
	// Get all migration files
	files, err := ioutil.ReadDir(migrationsPath)
	if err != nil {
		return nil, err
	}

	var allMigrations []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			allMigrations = append(allMigrations, file.Name())
		}
	}
	sort.Strings(allMigrations)

	// Get applied migrations
	applied, err := GetAppliedMigrations(db)
	if err != nil {
		return nil, err
	}

	appliedMap := make(map[string]bool)
	for _, m := range applied {
		appliedMap[m] = true
	}

	// Find pending
	var pending []string
	for _, m := range allMigrations {
		if !appliedMap[m] {
			pending = append(pending, m)
		}
	}

	return pending, nil
}
