package database

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"

	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB, migrationsPath string) error {

	files, err := filepath.Glob(filepath.Join(migrationsPath, "*.sql"))
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	if len(files) == 0 {
		log.Println("No migration files found in %s", migrationsPath)
		return nil
	}

	sort.Strings(files)

	log.Println("running migrations...")

	for _, file := range files {
		fileName := filepath.Base(file)

		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", fileName, err)
		}

		if err = db.Exec(string(content)).Error; err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", fileName, err)
		}

		log.Printf("applied migration: %s", fileName)
	}

	log.Println("all migrations applied successfully")
	return nil
}
