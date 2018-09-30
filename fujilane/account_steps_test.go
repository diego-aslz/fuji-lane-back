package fujilane

import (
	"fmt"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/nerde/fuji-lane-back/flactions"
	"github.com/nerde/fuji-lane-back/flentities"
)

type accountRow struct {
	flentities.Account
	Phone   string
	Country string
}

func theFollowingAccounts(table *gherkin.DataTable) error {
	return createFromTable(new(flentities.Account), table)
}

func iCreateTheFollowingAccount(table *gherkin.DataTable) error {
	b, err := assist.ParseMap(table)
	if err != nil {
		return err
	}

	return withRepository(func(r *flentities.Repository) error {
		country := &flentities.Country{}

		if err := r.Find(country, flentities.Country{Name: b["country"]}).Error; err != nil {
			return err
		}

		body := flactions.AccountsCreateBody{}
		body.Name = b["name"]
		body.Phone = b["phone"]
		body.UserName = b["user_name"]
		body.CountryID = int(country.ID)

		return performPOST(accountsPath, body)
	})
}

func iShouldHaveTheFollowingAccounts(table *gherkin.DataTable) error {
	return withRepository(func(r *flentities.Repository) error {
		accounts := []*flentities.Account{}
		err := r.Preload("Country").Find(&accounts).Error
		if err != nil {
			return err
		}

		rowsCount := len(table.Rows) - 1
		if len(accounts) != rowsCount {
			return fmt.Errorf("Expected to have %d accounts in the DB, got %d", rowsCount, len(accounts))
		}

		rows := []*accountRow{}
		for _, acc := range accounts {
			row := &accountRow{Account: *acc}

			if acc.Country != nil {
				row.Country = acc.Country.Name
			}

			if acc.Phone != nil {
				row.Phone = *acc.Phone
			}

			rows = append(rows, row)
		}

		return assist.CompareToSlice(rows, table)
	})
}

func AccountContext(s *godog.Suite) {
	s.Step(`^the following accounts:$`, theFollowingAccounts)
	s.Step(`^I create the following account:$`, iCreateTheFollowingAccount)
	s.Step(`^we should have the following accounts:$`, iShouldHaveTheFollowingAccounts)
}