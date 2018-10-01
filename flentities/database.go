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
	return WithRepository(func(r *Repository) error {
		driver, err := postgres.WithInstance(r.DB.DB(), &postgres.Config{})
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

// Seed the database
func Seed() error {
	return WithRepository(func(r *Repository) error {
		findOrCreate := [][]interface{}{
			[]interface{}{
				Country{Model: gorm.Model{ID: 1}},
				&Country{Model: gorm.Model{ID: 1}, Name: "China"},
			},
			[]interface{}{
				Country{Model: gorm.Model{ID: 2}},
				&Country{Model: gorm.Model{ID: 2}, Name: "Hong Kong"},
			},
			[]interface{}{
				Country{Model: gorm.Model{ID: 3}},
				&Country{Model: gorm.Model{ID: 3}, Name: "Japan"},
			},
			[]interface{}{
				Country{Model: gorm.Model{ID: 4}},
				&Country{Model: gorm.Model{ID: 4}, Name: "Singapore"},
			},
			[]interface{}{
				Country{Model: gorm.Model{ID: 5}},
				&Country{Model: gorm.Model{ID: 5}, Name: "Vietnam"},
			},
		}

		for _, pairs := range findOrCreate {
			if err := r.Where(pairs[0]).FirstOrCreate(pairs[1]).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// IsUniqueConstraintViolation returns true if the error is a unique constraint violation
func IsUniqueConstraintViolation(err error) bool {
	return strings.Contains(err.Error(), "violates unique constraint")
}
