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
	Phone   string
	Country string
}

func requestAccountsCreate(table *gherkin.DataTable) error {
	b, err := assist.ParseMap(table)
	if err != nil {
		return err
	}

	return flentities.WithRepository(func(r *flentities.Repository) error {
		country := &flentities.Country{}

		if err := r.Find(country, flentities.Country{Name: b["country"]}).Error; err != nil {
			return err
		}

		body := flactions.AccountsCreateBody{}
		body.Name = b["name"]
		body.Phone = b["phone"]
		body.UserName = b["userName"]
		body.CountryID = int(country.ID)

		return performPOST(flweb.AccountsPath, body)
	})
}

func assertAccounts(table *gherkin.DataTable) error {
	return assertDatabaseRecords(&[]*flentities.Account{}, table, accountToTableRow)
}

func tableRowToAccount(r *flentities.Repository, a interface{}) (interface{}, error) {
	row := a.(*accountRow)

	if row.Country != "" {
		row.Account.Country = &flentities.Country{}
		err := r.Find(row.Account.Country, flentities.Country{Name: row.Country}).Error
		if err != nil {
			return nil, err
		}
	}

	if row.Phone != "" {
		row.Account.Phone = &row.Phone
	}

	return &row.Account, nil
}

func accountToTableRow(r *flentities.Repository, a interface{}) (interface{}, error) {
	acc := a.(*flentities.Account)

	acc.Country = &flentities.Country{}
	err := r.Model(acc).Association("Country").Find(acc.Country).Error

	if err != nil && !gorm.IsRecordNotFoundError(err) {
		err = nil
	}

	return &accountRow{Account: *acc, Country: acc.Country.Name, Phone: derefStr(acc.Phone)}, err
}

func AccountContext(s *godog.Suite) {
	s.Step(`^the following accounts:$`, createFromTableStep(new(accountRow), tableRowToAccount))
	s.Step(`^I create the following account:$`, requestAccountsCreate)
	s.Step(`^we should have the following accounts:$`, assertAccounts)
}
