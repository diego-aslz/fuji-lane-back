package fujilane

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/nerde/fuji-lane-back/flactions"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/nerde/fuji-lane-back/flentities"
	"github.com/nerde/fuji-lane-back/flweb"
)

type imageRow struct {
	flentities.Image
	Property string
}

func (row *imageRow) save(r *flentities.Repository) error {
	if row.Property != "" {
		err := r.Find(&row.Image.Property, flentities.Property{Name: &row.Property}).Error
		if err != nil {
			return err
		}
	}

	return r.Create(&row.Image).Error
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

func createImages(table *gherkin.DataTable) error {
	return createRowEntitiesFromTable(new(imageRow), table)
}

func requestPropertiesImagesNew(table *gherkin.DataTable) error {
	image, err := assist.ParseMap(table)
	if err != nil {
		return err
	}

	propertyName := image["Property"]

	body := flactions.PropertiesImagesCreateBody{}
	body.Name = image["Name"]
	body.Size, err = strconv.Atoi(image["Size"])
	body.Type = image["Type"]

	if err != nil {
		return err
	}

	return flentities.WithRepository(func(r *flentities.Repository) error {
		property := &flentities.Property{}
		if err := r.Find(property, map[string]interface{}{"name": propertyName}).Error; err != nil {
			return err
		}

		path := strings.Replace(flweb.PropertiesImagesPath, ":id", fmt.Sprint(property.ID), 1)

		return performPOST(path, body)
	})
}

func requestPropertiesImagesUploaded(name string) error {
	return flentities.WithRepository(func(r *flentities.Repository) error {
		image := &flentities.Image{}
		if err := r.Find(image, map[string]interface{}{"name": name}).Error; err != nil {
			return err
		}

		path := strings.Replace(flweb.PropertiesImagesUploadedPath, ":property_id", fmt.Sprint(image.PropertyID), 1)
		path = strings.Replace(path, ":id", fmt.Sprint(image.ID), 1)

		return perform("PUT", path, nil)
	})
}

func ImageContext(s *godog.Suite) {
	s.Step(`^I should have the following images:$`, assertImages)
	s.Step(`^the following images:$`, createImages)
	s.Step(`^I request an URL to upload an image called "([^"]*)" for property "([^"]*)"$`, requestPropertiesImagesNew)
	s.Step(`^I mark image "([^"]*)" as uploaded$`, requestPropertiesImagesUploaded)
	s.Step(`^I request an URL to upload the following image:$`, requestPropertiesImagesNew)
}
