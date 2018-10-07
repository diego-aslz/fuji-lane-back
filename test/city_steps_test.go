package fujilane

import (
	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/nerde/fuji-lane-back/flentities"
)

type cityRow struct {
	flentities.City
	Country string
}

func (row *cityRow) save(r *flentities.Repository) error {
	if row.Country != "" {
		err := r.Find(&row.City.Country, flentities.Country{Name: row.Country}).Error
		if err != nil {
			return err
		}
	}
	return r.Create(&row.City).Error
}

func theFollowingCities(table *gherkin.DataTable) error {
	return createRowEntitiesFromTable(new(cityRow), table)
}

func CityContext(s *godog.Suite) {
	s.Step(`^the following cities:$`, theFollowingCities)
}
