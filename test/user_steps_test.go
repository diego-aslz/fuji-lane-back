package fujilane

import (
	"time"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/jinzhu/gorm"
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

func tableRowToUser(r *flentities.Repository, a interface{}) (interface{}, error) {
	row := a.(*userRow)

	if row.Password != "" {
		row.User.SetPassword(row.Password)
	}

	if row.Account != "" {
		row.User.Account = &flentities.Account{}
		err := r.Find(row.User.Account, flentities.Account{Name: row.Account}).Error
		if err != nil {
			return nil, err
		}
	}

	row.User.Name = refStr(row.Name)
	row.User.FacebookID = refStr(row.FacebookID)
	row.User.LastSignedIn = refTime(row.LastSignedIn)

	return &row.User, nil
}

func assertUsers(table *gherkin.DataTable) error {
	return assertDatabaseRecords(&[]*flentities.User{}, table, userToTableRow)
}

func userToTableRow(r *flentities.Repository, u interface{}) (interface{}, error) {
	user := u.(*flentities.User)

	user.Account = &flentities.Account{}
	if err := r.Model(u).Association("Account").Find(user.Account).Error; err != nil && !gorm.IsRecordNotFoundError(err) {
		return nil, err
	}

	row := &userRow{
		User:         *user,
		Name:         derefStr(user.Name),
		FacebookID:   derefStr(user.FacebookID),
		LastSignedIn: derefTime(user.LastSignedIn),
	}

	if user.Account != nil {
		row.Account = user.Account.Name
	}

	return row, nil
}

func UserContext(s *godog.Suite) {
	s.Step(`^the following users:$`, createFromTableStep(new(userRow), tableRowToUser))
	s.Step(`^we should have the following users:$`, assertUsers)
	s.Step(`^we should have no users$`, assertNoDatabaseRecordsStep(&flentities.User{}))
}
