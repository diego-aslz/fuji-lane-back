package fujilane

import (
	"time"

	"github.com/DATA-DOG/godog"
	"github.com/jinzhu/gorm"
	"github.com/nerde/fuji-lane-back/flentities"
	"github.com/nerde/fuji-lane-back/flweb"
)

type userRow struct {
	flentities.User
	Password     string
	Account      string
	LastSignedIn time.Time
}

func tableRowToUser(r *flentities.Repository, a interface{}) (interface{}, error) {
	row := a.(*userRow)
	row.User.LastSignedIn = refTime(row.LastSignedIn)

	if row.Password != "" {
		row.User.SetPassword(row.Password)
	}

	return &row.User, loadAssociationByName(&row.User, "Account", row.Account)
}

func userToTableRow(r *flentities.Repository, u interface{}) (interface{}, error) {
	user := u.(*flentities.User)

	user.Account = &flentities.Account{}
	if err := r.Model(u).Association("Account").Find(user.Account).Error; err != nil && !gorm.IsRecordNotFoundError(err) {
		return nil, err
	}

	row := &userRow{
		User:         *user,
		LastSignedIn: derefTime(user.LastSignedIn),
	}

	if user.Account != nil {
		row.Account = user.Account.Name
	}

	return row, nil
}

func UserContext(s *godog.Suite) {
	s.Step(`^the following users:$`, createFromTableStep(new(userRow), tableRowToUser))
	s.Step(`^I should have the following users:$`, assertDatabaseRecordsStep(&[]*flentities.User{}, userToTableRow))
	s.Step(`^I should have no users$`, assertNoDatabaseRecordsStep(&flentities.User{}))
	s.Step(`^I get my user details$`, performGETStep(flweb.UsersMePath))
}
