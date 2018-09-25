package fujilane

import (
	"reflect"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/jinzhu/gorm"
)

func theFollowingAccounts(table *gherkin.DataTable) error {
	sliceInterface, err := assist.CreateSlice(new(Account), table)
	if err != nil {
		return err
	}

	accounts := reflect.ValueOf(sliceInterface)

	return withDatabase(func(db *gorm.DB) error {
		for i := 0; i < accounts.Len(); i++ {
			acc, _ := accounts.Index(i).Interface().(*Account)

			err = db.Create(&acc).Error
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func AccountContext(s *godog.Suite) {
	s.Step(`^the following accounts:$`, theFollowingAccounts)
}
