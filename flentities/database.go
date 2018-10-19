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
			[]interface{}{
				City{Model: gorm.Model{ID: 101}},
				&City{Model: gorm.Model{ID: 101}, CountryID: 1, Name: "Beijing"},
			},
			[]interface{}{
				City{Model: gorm.Model{ID: 102}},
				&City{Model: gorm.Model{ID: 102}, CountryID: 1, Name: "Chengdu"},
			},
			[]interface{}{
				City{Model: gorm.Model{ID: 103}},
				&City{Model: gorm.Model{ID: 103}, CountryID: 1, Name: "Chongqing"},
			},
			[]interface{}{
				City{Model: gorm.Model{ID: 104}},
				&City{Model: gorm.Model{ID: 104}, CountryID: 1, Name: "Dongguan"},
			},
			[]interface{}{
				City{Model: gorm.Model{ID: 105}},
				&City{Model: gorm.Model{ID: 105}, CountryID: 1, Name: "Guangzhou"},
			},
			[]interface{}{
				City{Model: gorm.Model{ID: 106}},
				&City{Model: gorm.Model{ID: 106}, CountryID: 1, Name: "Shanghai"},
			},
			[]interface{}{
				City{Model: gorm.Model{ID: 107}},
				&City{Model: gorm.Model{ID: 107}, CountryID: 1, Name: "Shenyang"},
			},
			[]interface{}{
				City{Model: gorm.Model{ID: 108}},
				&City{Model: gorm.Model{ID: 108}, CountryID: 1, Name: "Shenzhen"},
			},
			[]interface{}{
				City{Model: gorm.Model{ID: 109}},
				&City{Model: gorm.Model{ID: 109}, CountryID: 1, Name: "Tianjin"},
			},
			[]interface{}{
				City{Model: gorm.Model{ID: 110}},
				&City{Model: gorm.Model{ID: 110}, CountryID: 1, Name: "Wuhan"},
			},
			[]interface{}{
				City{Model: gorm.Model{ID: 201}},
				&City{Model: gorm.Model{ID: 201}, CountryID: 2, Name: "Hong Kong"},
			},
			[]interface{}{
				City{Model: gorm.Model{ID: 301}},
				&City{Model: gorm.Model{ID: 301}, CountryID: 3, Name: "Fukuoka"},
			},
			[]interface{}{
				City{Model: gorm.Model{ID: 302}},
				&City{Model: gorm.Model{ID: 302}, CountryID: 3, Name: "Kawasaki"},
			},
			[]interface{}{
				City{Model: gorm.Model{ID: 303}},
				&City{Model: gorm.Model{ID: 303}, CountryID: 3, Name: "Kobe"},
			},
			[]interface{}{
				City{Model: gorm.Model{ID: 304}},
				&City{Model: gorm.Model{ID: 304}, CountryID: 3, Name: "Kyoto"},
			},
			[]interface{}{
				City{Model: gorm.Model{ID: 305}},
				&City{Model: gorm.Model{ID: 305}, CountryID: 3, Name: "Nagoya"},
			},
			[]interface{}{
				City{Model: gorm.Model{ID: 306}},
				&City{Model: gorm.Model{ID: 306}, CountryID: 3, Name: "Osaka"},
			},
			[]interface{}{
				City{Model: gorm.Model{ID: 307}},
				&City{Model: gorm.Model{ID: 307}, CountryID: 3, Name: "Saitama"},
			},
			[]interface{}{
				City{Model: gorm.Model{ID: 308}},
				&City{Model: gorm.Model{ID: 308}, CountryID: 3, Name: "Sapporo"},
			},
			[]interface{}{
				City{Model: gorm.Model{ID: 309}},
				&City{Model: gorm.Model{ID: 309}, CountryID: 3, Name: "Tokyo"},
			},
			[]interface{}{
				City{Model: gorm.Model{ID: 310}},
				&City{Model: gorm.Model{ID: 310}, CountryID: 3, Name: "Yokohama"},
			},
			[]interface{}{
				City{Model: gorm.Model{ID: 401}},
				&City{Model: gorm.Model{ID: 401}, CountryID: 4, Name: "Singapore"},
			},
			[]interface{}{
				City{Model: gorm.Model{ID: 501}},
				&City{Model: gorm.Model{ID: 501}, CountryID: 5, Name: "Ho Chi Minh"},
			},
			[]interface{}{
				City{Model: gorm.Model{ID: 502}},
				&City{Model: gorm.Model{ID: 502}, CountryID: 5, Name: "Hanoi"},
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
