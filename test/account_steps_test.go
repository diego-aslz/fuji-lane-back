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

type accountBodyTable struct {
	flactions.AccountsCreateBody
	Country string
}

func requestAccountsCreate(table *gherkin.DataTable) error {
	abt, err := assist.CreateInstance(new(accountBodyTable), table)
	if err != nil {
		return err
	}

	body := abt.(*accountBodyTable)

	country := &flentities.Country{}
	if err := findByName(country, body.Country); err != nil {
		return err
	}

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

	if err != nil && gorm.IsRecordNotFoundError(err) {
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
