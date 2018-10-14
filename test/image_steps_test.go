package fujilane

import (
	"encoding/json"
	"fmt"
	"regexp"
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

func tableRowToImage(r *flentities.Repository, a interface{}) (interface{}, error) {
	row := a.(*imageRow)
	return &row.Image, loadAssociationByName(&row.Image, "Property", row.Property)
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

	property := &flentities.Property{}
	if err := findByName(property, propertyName); err != nil {
		return err
	}

	path := strings.Replace(flweb.PropertiesImagesPath, ":id", fmt.Sprint(property.ID), 1)

	return performPOST(path, body)
}

func requestImagesUploaded(name string) error {
	image := &flentities.Image{}
	if err := findByName(image, name); err != nil {
		return err
	}

	path := strings.Replace(flweb.ImagesUploadedPath, ":id", fmt.Sprint(image.ID), 1)

	return perform("PUT", path, nil)
}

func requestPropertiesImagesDestroy(name string) error {
	image := &flentities.Image{}
	if err := findByName(image, name); err != nil {
		return err
	}

	path := strings.Replace(flweb.ImagePath, ":id", fmt.Sprint(image.ID), 1)

	return perform("DELETE", path, nil)
}

func assertResponseStatusAndImage(status string, table *gherkin.DataTable) error {
	if err := assertResponseStatus(status); err != nil {
		return err
	}

	image := &flentities.Image{}
	if err := json.Unmarshal([]byte(response.Body.String()), image); err != nil {
		return fmt.Errorf("Unable to unmarshal %s: %s", response.Body.String(), err.Error())
	}

	// Discarding fields that cannot be asserted because they are dynamic
	reps := map[string]string{"Amz-Signature": "SIGNATURE", "Amz-Date": "DATE", "Amz-Credential": "CREDENTIAL"}
	for key, replacement := range reps {
		reg := regexp.MustCompile(key + "=([\\w\\-]|%2F)+")
		image.URL = string(reg.ReplaceAll([]byte(image.URL), []byte(key+"="+replacement)))
	}

	return assist.CompareToInstance(image, table)
}

func ImageContext(s *godog.Suite) {
	s.Step(`^I should have the following images:$`, assertDatabaseRecordsStep(&[]*flentities.Image{}, imageToTableRow))
	s.Step(`^the following images:$`, createFromTableStep(new(imageRow), tableRowToImage))
	s.Step(`^I mark image "([^"]*)" as uploaded$`, requestImagesUploaded)
	s.Step(`^I request an URL to upload the following image:$`, requestPropertiesImagesCreate)
	s.Step(`^I remove the image "([^"]*)"$`, requestPropertiesImagesDestroy)
	s.Step(`^I should have no images$`, assertNoDatabaseRecordsStep(&flentities.Image{}))
	s.Step(`^the system should respond with "([^"]*)" and the following image:$`, assertResponseStatusAndImage)
}
