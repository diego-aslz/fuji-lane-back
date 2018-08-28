package fujilane

import (
	"fmt"
	"os"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file" // File migrations
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // Postgres driver
)

func withDatabase(callback func(*gorm.DB) error) error {
	url := os.Getenv("DATABASE_URL")
	db, err := gorm.Open("postgres", url)
	if err != nil {
		return fmt.Errorf("Unable to connect to %s: %s", url, err.Error())
	}
	defer db.Close()
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
