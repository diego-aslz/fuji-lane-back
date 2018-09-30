package fujilane

import (
	"fmt"
	"time"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/nerde/fuji-lane-back/flentities"
)

type userRow struct {
	flentities.User
	Name         string
	Password     string
	Account      string
	FacebookID   string
	LastSignedIn time.Time
}

func (row *userRow) save(r *flentities.Repository) error {
	if row.Password != "" {
		row.User.SetPassword(row.Password)
	}

	if row.Account != "" {
		row.User.Account = &flentities.Account{}
		err := r.Find(row.User.Account, flentities.Account{Name: row.Account}).Error
		if err != nil {
			return err
		}
	}

	if row.Name != "" {
		row.User.Name = &row.Name
	}

	if row.FacebookID != "" {
		row.User.FacebookID = &row.FacebookID
	}

	if !row.LastSignedIn.IsZero() {
		row.User.LastSignedIn = &row.LastSignedIn
	}

	return r.Create(&row.User).Error
}

func newUserRow(u *flentities.User) (row *userRow) {
	row = &userRow{}
	row.User = *u

	if u.Name != nil {
		row.Name = *u.Name
	}

	if u.FacebookID != nil {
		row.FacebookID = *u.FacebookID
	}

	if u.Account != nil {
		row.Account = u.Account.Name
	}

	if u.LastSignedIn != nil {
		row.LastSignedIn = *u.LastSignedIn
	}

	return
}

func theFollowingUsers(table *gherkin.DataTable) error {
	return createRowEntitiesFromTable(new(userRow), table)
}

func weShouldHaveTheFollowingUsers(table *gherkin.DataTable) error {
	return withRepository(func(r *flentities.Repository) error {
		count := 0
		err := r.Model(&flentities.User{}).Count(&count).Error
		if err != nil {
			return err
		}

		rowsCount := len(table.Rows) - 1
		if count != rowsCount {
			return fmt.Errorf("Expected to have %d users in the DB, got %d", rowsCount, count)
		}

		users := []*flentities.User{}
		err = r.Preload("Account").Find(&users).Error
		if err != nil {
			return err
		}

		rows := []*userRow{}
		for _, user := range users {
			rows = append(rows, newUserRow(user))
		}

		return assist.CompareToSlice(rows, table)
	})
}

func weShouldHaveNoUsers() error {
	return withRepository(func(r *flentities.Repository) error {
		count := 0
		err := r.Model(&flentities.User{}).Count(&count).Error
		if err != nil {
			return err
		}

		if count != 0 {
			return fmt.Errorf("Expected to have %d users in the DB, got %d", 0, count)
		}

		return nil
	})
}

func UserContext(s *godog.Suite) {
	s.Step(`^the following users:$`, theFollowingUsers)
	s.Step(`^we should have the following users:$`, weShouldHaveTheFollowingUsers)
	s.Step(`^we should have no users$`, weShouldHaveNoUsers)
}
