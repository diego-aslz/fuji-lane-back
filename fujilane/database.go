package fujilane

import (
	"fmt"
	"strings"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file" // File migrations
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // Postgres driver
)

func withDatabase(callback func(*gorm.DB) error) error {
	url := appConfig.databaseURL
	db, err := gorm.Open("postgres", url)
	if err != nil {
		return fmt.Errorf("Unable to connect to %s: %s", url, err.Error())
	}
	defer db.Close()
	db.LogMode(appConfig.databaseLogs)
	return callback(db)
}

// Migrate the database
func Migrate() error {
	return withDatabase(func(db *gorm.DB) error {
		driver, err := postgres.WithInstance(db.DB(), &postgres.Config{})
		if err != nil {
			return err
		}

		m, err := migrate.NewWithDatabaseInstance("file://./migrations", "postgres", driver)
		if err != nil {
			return err
		}

		return m.Up()
	})
}

func isUniqueConstraintViolation(err error) bool {
	return strings.Contains(err.Error(), "violates unique constraint")
}
