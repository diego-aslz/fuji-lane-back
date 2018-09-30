package flentities

import (
	"fmt"
	"strings"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file" // File migrations
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // Postgres driver
	"github.com/nerde/fuji-lane-back/flconfig"
)

// WithDatabase gets a database connection and calls the callback with it, taking care of connection errors
func WithDatabase(config *flconfig.Configuration, callback func(*gorm.DB) error) error {
	db, err := gorm.Open("postgres", config.DatabaseURL)
	if err != nil {
		return fmt.Errorf("Unable to connect to %s: %s", config.DatabaseURL, err.Error())
	}

	defer db.Close()

	return callback(db.
		LogMode(config.DatabaseLogs).
		Set("gorm:association_autocreate", false).
		Set("gorm:association_autoupdate", false))
}

// Migrate the database
func Migrate(config *flconfig.Configuration) error {
	return WithDatabase(config, func(db *gorm.DB) error {
		driver, err := postgres.WithInstance(db.DB(), &postgres.Config{})
		if err != nil {
			return err
		}

		m, err := migrate.NewWithDatabaseInstance("file://./migrations", "postgres", driver)
		if err != nil {
			return err
		}

		if err = m.Up(); err == migrate.ErrNoChange {
			err = nil
		}

		return err
	})
}

// IsUniqueConstraintViolation returns true if the error is a unique constraint violation
func IsUniqueConstraintViolation(err error) bool {
	return strings.Contains(err.Error(), "violates unique constraint")
}
