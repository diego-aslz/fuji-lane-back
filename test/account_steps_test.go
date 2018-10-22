package fujilane

import (
	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/jinzhu/gorm"
	"github.com/nerde/fuji-lane-back/flactions"
	"github.com/nerde/fuji-lane-back/flentities"
	"github.com/nerde/fuji-lane-back/flweb"
)

type accountRow struct {
	flentities.Account
	Country string
}

func requestAccountsCreate(table *gherkin.DataTable) error {
	b, err := assist.ParseMap(table)
	if err != nil {
		return err
	}

	country := &flentities.Country{}
	if err := findByName(country, b["country"]); err != nil {
		return err
	}

	body := flactions.AccountsCreateBody{}
	body.Name = b["name"]
	body.Phone = b["phone"]
	body.UserName = b["userName"]
	body.CountryID = country.ID

	return performPOST(flweb.AccountsPath, body)
}

func tableRowToAccount(r *flentities.Repository, a interface{}) (interface{}, error) {
	row := a.(*accountRow)

	return &row.Account, loadAssociationByName(&row.Account, "Country", row.Country)
}

func accountToTableRow(r *flentities.Repository, a interface{}) (interface{}, error) {
	acc := a.(*flentities.Account)

	acc.Country = &flentities.Country{}
	err := r.Model(acc).Association("Country").Find(acc.Country).Error

	if err != nil && !gorm.IsRecordNotFoundError(err) {
		err = nil
	}

	return &accountRow{Account: *acc, Country: acc.Country.Name}, err
}

func AccountContext(s *godog.Suite) {
	s.Step(`^the following accounts:$`, createFromTableStep(new(accountRow), tableRowToAccount))
	s.Step(`^I create the following account:$`, requestAccountsCreate)
	s.Step(`^I should have the following accounts:$`, assertDatabaseRecordsStep(&[]*flentities.Account{},
		accountToTableRow))
}
