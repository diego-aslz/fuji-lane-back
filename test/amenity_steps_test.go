package fujilane

import (
	"strings"

	"github.com/DATA-DOG/godog"
	"github.com/nerde/fuji-lane-back/flentities"
	"github.com/nerde/fuji-lane-back/flweb"
)

func requestAmenities(target string) error {
	return perform("GET", strings.Replace(flweb.AmenityTypesPath, ":target", target, 1), nil)
}

func AmenityContext(s *godog.Suite) {
	s.Step(`^I list amenity types for "([^"]*)"$`, requestAmenities)
	s.Step(`^the system should respond with "([^"]*)" and the following amenity types:$`,
		assertResponseStatusAndListStep(&[]*flentities.AmenityType{}))
}
