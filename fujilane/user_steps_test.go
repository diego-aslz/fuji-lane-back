package fujilane

import (
	"fmt"
	"reflect"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/jinzhu/gorm"
)

type userRow struct {
	User
	Name       string
	Password   string
	Account    string
	FacebookID string
}

func (row *userRow) save(db *gorm.DB) error {
	if row.Password != "" {
		row.User.setPassword(row.Password)
	}

	if row.Account != "" {
		row.User.Account = &Account{}
		err := db.Find(row.User.Account, Account{Name: row.Account}).Error
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

	return db.Create(&row.User).Error
}

func newUserRow(u *User) (row *userRow) {
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

	return
}

func theFollowingUsers(table *gherkin.DataTable) error {
	sliceInterface, err := assist.CreateSlice(new(userRow), table)
	if err != nil {
		return err
	}

	users := reflect.ValueOf(sliceInterface)

	return withDatabase(func(db *gorm.DB) error {
		for i := 0; i < users.Len(); i++ {
			row, _ := users.Index(i).Interface().(*userRow)

			err = row.save(db)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func weShouldHaveTheFollowingUsers(table *gherkin.DataTable) error {
	return withDatabase(func(db *gorm.DB) error {
		count := 0
		err := db.Model(&User{}).Count(&count).Error
		if err != nil {
			return err
		}

		rowsCount := len(table.Rows) - 1
		if count != rowsCount {
			return fmt.Errorf("Expected to have %d users in the DB, got %d", rowsCount, count)
		}

		users := []*User{}
		err = db.Preload("Account").Find(&users).Error
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
	return withDatabase(func(db *gorm.DB) error {
		count := 0
		err := db.Model(&User{}).Count(&count).Error
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
