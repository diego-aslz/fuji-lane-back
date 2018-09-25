package fujilane

import (
	"log"
	"reflect"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/jinzhu/gorm"
)

func createFromTable(tp interface{}, table *gherkin.DataTable) error {
	sliceInterface, err := assist.CreateSlice(tp, table)
	if err != nil {
		return err
	}

	records := reflect.ValueOf(sliceInterface)

	return withDatabase(func(db *gorm.DB) error {
		for i := 0; i < records.Len(); i++ {
			err = db.Create(records.Index(i).Interface()).Error
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func cleanup(_ interface{}, _ error) {
	err := withDatabase(func(db *gorm.DB) error {
		for _, model := range []interface{}{Property{}, User{}, Account{}, Country{}} {
			err := db.Unscoped().Delete(model).Error
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

func DatabaseContext(s *godog.Suite) {
	s.AfterScenario(cleanup)
}
