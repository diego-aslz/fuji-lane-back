package fujilane

import (
	"fmt"
	"reflect"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/jinzhu/gorm"
)

func theFollowingUsers(table *gherkin.DataTable) error {
	sliceInterface, err := assist.CreateSlice(new(User), table)
	if err != nil {
		return err
	}

	users := reflect.ValueOf(sliceInterface)

	return withDatabase(func(db *gorm.DB) error {
		for i := 0; i < users.Len(); i++ {
			user, _ := users.Index(i).Interface().(*User)
			err = db.Create(user).Error
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

func UserContext(s *godog.Suite) {
	s.Step(`^the following users:$`, theFollowingUsers)
	s.Step(`^we should have the following users:$`, weShouldHaveTheFollowingUsers)
}
