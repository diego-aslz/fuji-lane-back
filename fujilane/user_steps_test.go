package fujilane

import (
	"fmt"
	"reflect"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/jinzhu/gorm"
)

type testUser struct {
	User
	Password string
	Account  string
}

func theFollowingUsers(table *gherkin.DataTable) error {
	sliceInterface, err := assist.CreateSlice(new(testUser), table)
	if err != nil {
		return err
	}

	users := reflect.ValueOf(sliceInterface)

	return withDatabase(func(db *gorm.DB) error {
		for i := 0; i < users.Len(); i++ {
			tu, _ := users.Index(i).Interface().(*testUser)

			if tu.Password != "" {
				tu.User.setPassword(tu.Password)
			}

			if tu.Account != "" {
				tu.User.Account = &Account{}
				err = db.Find(tu.User.Account, Account{Name: tu.Account}).Error
				if err != nil {
					return err
				}
			}

			err = db.Create(&tu.User).Error
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
		db.Find(&users)
		return assist.CompareToSlice(users, table)
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
