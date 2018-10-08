package fujilane

import (
	"fmt"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/nerde/fuji-lane-back/flentities"
)

type imageRow struct {
	flentities.Image
	Property string
}

func assertImages(table *gherkin.DataTable) error {
	return flentities.WithRepository(func(r *flentities.Repository) error {
		images := []*flentities.Image{}
		err := r.Preload("Property").Find(&images).Error
		if err != nil {
			return err
		}

		rowsCount := len(table.Rows) - 1
		if len(images) != rowsCount {
			return fmt.Errorf("Expected to have %d images in the DB, got %d", rowsCount, len(images))
		}

		rows := []*imageRow{}
		for _, acc := range images {
			row := &imageRow{Image: *acc}

			if acc.Property.Name != nil {
				row.Property = *acc.Property.Name
			}

			rows = append(rows, row)
		}

		return assist.CompareToSlice(rows, table)
	})
}

func ImageContext(s *godog.Suite) {
	s.Step(`^I should have the following images:$`, assertImages)
}
