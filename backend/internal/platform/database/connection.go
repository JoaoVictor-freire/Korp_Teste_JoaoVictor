package database

import (
	"errors"
	"os"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func OpenFromEnv() (*gorm.DB, error) {
	dsn := strings.TrimSpace(os.Getenv("DATABASE_URL"))
	if dsn == "" {
		return nil, errors.New("DATABASE_URL is not set")
	}

	return Open(dsn)
}

func Open(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	DB = db
	return db, nil
}

// Connect keeps compatibility with early experiments.
// Prefer Open/OpenFromEnv and handle errors explicitly.
func Connect() {
	_, err := OpenFromEnv()
	if err != nil {
		panic(err)
	}
}
