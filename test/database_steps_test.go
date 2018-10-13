package fujilane

import (
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/nerde/fuji-lane-back/flentities"
)

type entityDecorator func(*flentities.Repository, interface{}) (interface{}, error)

func createFromTableStep(tp interface{}, decorators ...entityDecorator) func(table *gherkin.DataTable) error {
	return func(table *gherkin.DataTable) error {
		return createFromTable(tp, table, decorators...)
	}
}

func createFromTable(tp interface{}, table *gherkin.DataTable, decorators ...entityDecorator) error {
	sliceInterface, err := assist.CreateSlice(tp, table)
	if err != nil {
		return err
	}

	records := reflect.ValueOf(sliceInterface)

	return flentities.WithRepository(func(r *flentities.Repository) error {
		var err error

		for i := 0; i < records.Len(); i++ {
			record := records.Index(i).Interface()

			for _, dec := range decorators {
				record, err = dec(r, record)
				if err != nil {
					return err
				}
			}

			if err = r.Create(record).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func assertDatabaseRecords(slice interface{}, table *gherkin.DataTable, decorators ...entityDecorator) error {
	return flentities.WithRepository(func(r *flentities.Repository) error {
		err := r.Find(slice).Error
		if err != nil {
			return err
		}

		elem := reflect.ValueOf(slice).Elem()
		count := elem.Len()
		expectedCount := len(table.Rows) - 1
		if count != expectedCount {
			return fmt.Errorf("Expected to have %d records in the database, got %d", expectedCount, count)
		}

		records := []interface{}{}

		for i := 0; i < elem.Len(); i++ {
			records = append(records, elem.Index(i).Interface())

			for _, dec := range decorators {
				records[i], err = dec(r, records[i])
				if err != nil {
					return err
				}
			}
		}

		return assist.CompareToSlice(records, table)
	})
}

func assertNoDatabaseRecordsStep(model interface{}) func() error {
	return func() error {
		return assertNoDatabaseRecords(model)
	}
}

func assertNoDatabaseRecords(model interface{}) error {
	return flentities.WithRepository(func(r *flentities.Repository) error {
		count := 0
		err := r.Model(model).Count(&count).Error
		if err != nil {
			return err
		}

		if count != 0 {
			return fmt.Errorf("Expected to have %d records in the database, got %d", 0, count)
		}

		return nil
	})
}

func derefStr(str *string) string {
	if str == nil {
		return ""
	}

	return *str
}

func derefTime(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}

	return *t
}

func refStr(str string) *string {
	if str == "" {
		return nil
	}

	return &str
}

func refTime(t time.Time) *time.Time {
	if t.IsZero() {
		return nil
	}

	return &t
}

func cleanup(_ interface{}, _ error) {
	err := flentities.WithRepository(func(r *flentities.Repository) error {
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

func databaseDefaultsAreLoaded() error {
	return flentities.Seed()
}

func DatabaseContext(s *godog.Suite) {
	s.Step(`^defaults are loaded$`, databaseDefaultsAreLoaded)
	s.AfterScenario(cleanup)
}
