package fujilane

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"
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

func assertImages(table *gherkin.DataTable) error {
	return assertDatabaseRecords(&[]*flentities.Image{}, table, imageToTableRow)
}

func imageToTableRow(r *flentities.Repository, i interface{}) (interface{}, error) {
	image := i.(*flentities.Image)

	err := r.Model(i).Association("Property").Find(&image.Property).Error
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return nil, err
	}

	row := &imageRow{Image: *image}

	if image.Property.Name != nil {
		row.Property = *image.Property.Name
	}

	return row, nil
}

func createImages(table *gherkin.DataTable) error {
	return createFromTable(new(imageRow), table, tableRowToImage)
}

func tableRowToImage(r *flentities.Repository, a interface{}) (interface{}, error) {
	row := a.(*imageRow)

	if row.Property != "" {
		err := r.Find(&row.Image.Property, flentities.Property{Name: &row.Property}).Error
		if err != nil {
			return nil, err
		}
	}

	return &row.Image, nil
}

func requestPropertiesImagesCreate(table *gherkin.DataTable) error {
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

func requestPropertiesImagesDestroy(name string) error {
	return flentities.WithRepository(func(r *flentities.Repository) error {
		image := &flentities.Image{}
		if err := r.Find(image, map[string]interface{}{"name": name}).Error; err != nil {
			return err
		}

		path := strings.Replace(flweb.ImagePath, ":id", fmt.Sprint(image.ID), 1)

		return perform("DELETE", path, nil)
	})
}

func ImageContext(s *godog.Suite) {
	s.Step(`^I should have the following images:$`, assertImages)
	s.Step(`^the following images:$`, createImages)
	s.Step(`^I request an URL to upload an image called "([^"]*)" for property "([^"]*)"$`, requestPropertiesImagesCreate)
	s.Step(`^I mark image "([^"]*)" as uploaded$`, requestPropertiesImagesUploaded)
	s.Step(`^I request an URL to upload the following image:$`, requestPropertiesImagesCreate)
	s.Step(`^I remove the image "([^"]*)"$`, requestPropertiesImagesDestroy)
	s.Step(`^I should have no images$`, assertNoDatabaseRecordsStep(&flentities.Image{}))
}
