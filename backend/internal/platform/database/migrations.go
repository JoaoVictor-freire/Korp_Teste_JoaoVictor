package database

import (
	"os"
	"strings"

	"gorm.io/gorm"
)

func AutoMigrateEnabled() bool {
	value := strings.ToLower(strings.TrimSpace(os.Getenv("AUTO_MIGRATE")))
	return value == "1" || value == "true" || value == "yes"
}

func Migrate(db *gorm.DB, models ...any) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// Serialize migrations across services sharing the same database.
		if err := tx.Exec(`SELECT pg_advisory_xact_lock(20260409);`).Error; err != nil {
			return err
		}

		if err := tx.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`).Error; err != nil {
			return err
		}

		if len(models) == 0 || !AutoMigrateEnabled() {
			return nil
		}

		return tx.AutoMigrate(models...)
	})
}
