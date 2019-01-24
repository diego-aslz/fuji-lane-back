package fujilane

import (
	"github.com/DATA-DOG/godog"
	"github.com/jinzhu/gorm"
	"github.com/nerde/fuji-lane-back/flentities"
)

type priceRow struct {
	flentities.Price
	Unit string
}

func tableRowToPrice(r *flentities.Repository, a interface{}) (interface{}, error) {
	row := a.(*priceRow)

	return &row.Price, loadAssociationByName(&row.Price, "Unit", row.Unit)
}

func priceToTableRow(r *flentities.Repository, a interface{}) (interface{}, error) {
	p := a.(*flentities.Price)

	p.Unit = &flentities.Unit{}
	err := r.Model(p).Association("Unit").Find(p.Unit).Error

	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return nil, err
	}

	return &priceRow{Price: *p, Unit: p.Unit.Name}, err
}

func PriceContext(s *godog.Suite) {
	s.Step(`^the following prices:$`, createFromTableStep(new(priceRow), tableRowToPrice))
	s.Step(`^I should have the following prices:$`, assertDatabaseRecordsStep(&[]*flentities.Price{},
		priceToTableRow))
}
