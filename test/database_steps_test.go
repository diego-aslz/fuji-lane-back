package fujilane

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/jinzhu/gorm"
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

func assertDatabaseRecordsStep(slice interface{}, decorators ...entityDecorator) func(table *gherkin.DataTable) error {
	return func(table *gherkin.DataTable) error {
		return assertDatabaseRecords(slice, table, decorators...)
	}
}

func assertDatabaseRecords(slice interface{}, table *gherkin.DataTable, decorators ...entityDecorator) error {
	return flentities.WithRepository(func(r *flentities.Repository) error {
		err := r.Order("id").Find(slice).Error
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

func derefInt(i *int) int {
	if i == nil {
		return 0
	}

	return *i
}

func derefUint(i *uint) uint {
	if i == nil {
		return 0
	}

	return *i
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

func refInt(i int) *int {
	if i == 0 {
		return nil
	}

	return &i
}

func refFloat(f float32) *float32 {
	if f == 0 {
		return nil
	}

	return &f
}

func refTime(t time.Time) *time.Time {
	if t.IsZero() {
		return nil
	}

	return &t
}

func refUint(i uint) *uint {
	if i == 0 {
		return nil
	}

	return &i
}

var allEntities = map[string]interface{}{
	"amenities":  flentities.Amenity{},
	"units":      flentities.Unit{},
	"images":     flentities.Image{},
	"properties": flentities.Property{},
	"users":      flentities.User{},
	"accounts":   flentities.Account{},
	"cities":     flentities.City{},
	"countries":  flentities.Country{},
}

func cleanup(_ interface{}, _ error) {
	err := flentities.WithRepository(func(r *flentities.Repository) error {
		for tableName, model := range allEntities {
			err := r.Unscoped().Delete(model).Error
			if err != nil {
				if strings.Index(err.Error(), "violates foreign key") >= 0 {
					if err = r.Exec(fmt.Sprintf("TRUNCATE %s CASCADE", tableName)).Error; err != nil {
						return err
					}
					continue
				}
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

func loadAssociationByName(record interface{}, nameAndValues ...string) error {
	for i := 0; i <= len(nameAndValues)-2; i += 2 {
		assocName := nameAndValues[i]
		value := nameAndValues[i+1]

		if value == "" {
			continue
		}

		field := reflect.ValueOf(record).Elem().FieldByName(assocName)

		if field.Kind() == reflect.Ptr && field.IsNil() {
			field.Set(reflect.New(field.Type().Elem()))
		} else {
			field = field.Addr()
		}

		if err := findByName(field.Interface(), value); err != nil {
			return err
		}
	}
	return nil
}

func findByName(record interface{}, name string) (err error) {
	return findBy(record, "name", name)
}

func findBy(record interface{}, field, value string) (err error) {
	return flentities.WithRepository(func(r *flentities.Repository) error {
		return r.Find(record, map[string]interface{}{field: value}).Error
	})
}

func loadAssociations(r *flentities.Repository, model interface{}, assocs map[string]interface{}) error {
	for assocName, field := range assocs {
		err := r.Model(model).Association(assocName).Find(field).Error

		if gorm.IsRecordNotFoundError(err) {
			err = nil
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func DatabaseContext(s *godog.Suite) {
	s.Step(`^defaults are loaded$`, databaseDefaultsAreLoaded)
	s.AfterScenario(cleanup)
}
