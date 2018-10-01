package fujilane

import (
	"log"
	"reflect"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/nerde/fuji-lane-back/flentities"
)

type rowEntity interface {
	save(*flentities.Repository) error
}

func createEntitiesFromTable(tp interface{}, table *gherkin.DataTable) error {
	return createFromTable(tp, table, func(obj interface{}, r *flentities.Repository) error {
		return r.Create(obj).Error
	})
}

func createRowEntitiesFromTable(tp interface{}, table *gherkin.DataTable) error {
	return createFromTable(tp, table, func(obj interface{}, r *flentities.Repository) error {
		return obj.(rowEntity).save(r)
	})
}

func createFromTable(tp interface{}, table *gherkin.DataTable, onSave func(interface{}, *flentities.Repository) error) error {
	sliceInterface, err := assist.CreateSlice(tp, table)
	if err != nil {
		return err
	}

	records := reflect.ValueOf(sliceInterface)

	return flentities.WithRepository(func(r *flentities.Repository) error {
		for i := 0; i < records.Len(); i++ {
			if err = onSave(records.Index(i).Interface(), r); err != nil {
				return err
			}
		}

		return nil
	})
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

func DatabaseContext(s *godog.Suite) {
	s.AfterScenario(cleanup)
}
