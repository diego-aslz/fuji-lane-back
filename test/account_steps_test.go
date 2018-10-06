package fujilane

import (
	"fmt"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/nerde/fuji-lane-back/flactions"
	"github.com/nerde/fuji-lane-back/flentities"
	"github.com/nerde/fuji-lane-back/flweb"
)

type accountRow struct {
	flentities.Account
	Phone   string
	Country string
}

func (row *accountRow) save(r *flentities.Repository) error {
	if row.Country != "" {
		row.Account.Country = &flentities.Country{}
		err := r.Find(row.Account.Country, flentities.Country{Name: row.Country}).Error
		if err != nil {
			return err
		}
	}

	if row.Phone != "" {
		row.Account.Phone = &row.Phone
	}

	return r.Create(&row.Account).Error
}

func theFollowingAccounts(table *gherkin.DataTable) error {
	return createRowEntitiesFromTable(new(accountRow), table)
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
		body.UserName = b["user_name"]
		body.CountryID = int(country.ID)

		return performPOST(flweb.AccountsPath, body)
	})
}

func assertAccounts(table *gherkin.DataTable) error {
	return flentities.WithRepository(func(r *flentities.Repository) error {
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
	s.Step(`^I create the following account:$`, requestAccountsCreate)
	s.Step(`^we should have the following accounts:$`, assertAccounts)
}
