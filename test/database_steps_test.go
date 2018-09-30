package fujilane

import (
	"log"
	"reflect"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/jinzhu/gorm"
	"github.com/nerde/fuji-lane-back/flconfig"
	"github.com/nerde/fuji-lane-back/flentities"
)

func createFromTable(tp interface{}, table *gherkin.DataTable) error {
	sliceInterface, err := assist.CreateSlice(tp, table)
	if err != nil {
		return err
	}

	records := reflect.ValueOf(sliceInterface)

	return withRepository(func(r *flentities.Repository) error {
		for i := 0; i < records.Len(); i++ {
			err = r.Create(records.Index(i).Interface()).Error
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func cleanup(_ interface{}, _ error) {
	err := withRepository(func(r *flentities.Repository) error {
		for _, model := range flentities.AllEntities() {
			err := r.Unscoped().Delete(model).Error
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		log.Fatal(err.Error())
	}
}

func withRepository(callback func(*flentities.Repository) error) error {
	return flentities.WithDatabase(flconfig.Config, func(db *gorm.DB) error {
		return callback(&flentities.Repository{DB: db})
	})
}

func DatabaseContext(s *godog.Suite) {
	s.AfterScenario(cleanup)
}
