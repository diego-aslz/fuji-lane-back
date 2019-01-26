package fujilane

import (
	"encoding/json"
	"fmt"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/nerde/fuji-lane-back/flweb"
)

type searchUnit struct {
	Name               string `json:"name"`
	PerNightPriceCents int    `json:"perNightPriceCents"`
	TotalCents         int    `json:"totalCents"`
}

type searchProperty struct {
	Name  string       `json:"name"`
	Units []searchUnit `json:"units"`
}

type searchResult struct {
	PropertyName string
	searchUnit
}

func assertSearchResults(status string, table *gherkin.DataTable) error {
	if err := assertResponseStatus(status); err != nil {
		return err
	}

	actualBody := []searchProperty{}
	if err := json.Unmarshal([]byte(response.Body.String()), &actualBody); err != nil {
		return fmt.Errorf("Unable to unmarshal %s: %s", response.Body.String(), err.Error())
	}

	searchResults := []*searchResult{}
	for _, prop := range actualBody {
		for _, unit := range prop.Units {
			result := &searchResult{}
			result.PropertyName = prop.Name
			result.searchUnit = unit

			searchResults = append(searchResults, result)
		}
	}

	return assist.CompareToSlice(searchResults, table)
}

func SearchContext(s *godog.Suite) {
	s.Step(`^I search for units with the following filters:$`, performGETWithQueryStep(flweb.SearchPath))
	s.Step(`^the system should respond with "([^"]*)" and the following search results:$`, assertSearchResults)
}
