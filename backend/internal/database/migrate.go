package database

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

// RunMigrations executes SQL migration files on startup
func RunMigrations(migrationsPath string) error {
	log.Info().Str("path", migrationsPath).Msg("Running database migrations")

	// Read all migration files
	files, err := filepath.Glob(filepath.Join(migrationsPath, "*.sql"))
	if err != nil {
		return fmt.Errorf("failed to read migration files: %w", err)
	}

	if len(files) == 0 {
		log.Warn().Msg("No migration files found")
		return nil
	}

	// Execute each migration file in order
	for _, file := range files {
		log.Info().Str("file", filepath.Base(file)).Msg("Executing migration")

		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", file, err)
		}

		_, err = DB.Exec(string(content))
		if err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", file, err)
		}

		log.Info().Str("file", filepath.Base(file)).Msg("Migration executed successfully")
	}

	log.Info().Int("count", len(files)).Msg("All migrations completed successfully")
	return nil
}
