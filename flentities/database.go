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

// WithRepository gets a database connection and calls the callback with a connected Repository, returning any
// connection errors
func WithRepository(callback func(*Repository) error) error {
	db, err := gorm.Open("postgres", flconfig.Config.DatabaseURL)
	if err != nil {
		return fmt.Errorf("Unable to connect to %s: %s", flconfig.Config.DatabaseURL, err.Error())
	}

	defer db.Close()

	return callback(&Repository{db.
		LogMode(flconfig.Config.DatabaseLogs).
		Set("gorm:association_autocreate", false).
		Set("gorm:association_autoupdate", false)})
}

// Migrate the database
func Migrate() error {
	return withMigrations(func(r *Repository, m *migrate.Migrate) error {
		if err := m.Up(); err != migrate.ErrNoChange {
			return err
		}

		return nil
	})
}

// Reset the database, redoing all migrations
func Reset() error {
	err := withMigrations(func(r *Repository, m *migrate.Migrate) error {
		if err := m.Down(); err != migrate.ErrNoChange {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	// We have to call withMigrations because otherwise Up() won't take effect because migrate.NewWithDatabaseInstance
	// loads the current database version once
	return withMigrations(func(r *Repository, m *migrate.Migrate) error {
		if err := m.Up(); err != migrate.ErrNoChange {
			return err
		}

		return nil
	})
}

func withMigrations(fn func(*Repository, *migrate.Migrate) error) error {
	return WithRepository(func(r *Repository) error {
		driver, err := postgres.WithInstance(r.DB.DB(), &postgres.Config{})
		if err != nil {
			return err
		}

		m, err := migrate.NewWithDatabaseInstance("file://./migrations", "postgres", driver)
		if err != nil {
			return err
		}

		return fn(r, m)
	})
}

// IsUniqueConstraintViolation returns true if the error is a unique constraint violation
func IsUniqueConstraintViolation(err error) bool {
	return strings.Contains(err.Error(), "violates unique constraint")
}
