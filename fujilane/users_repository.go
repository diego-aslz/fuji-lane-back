package fujilane

import (
	"github.com/jinzhu/gorm"
	"github.com/nerde/fuji-lane-back/flentities"
)

type usersRepository struct{}

func (r *usersRepository) findForFacebookSignIn(facebookID, email string) (*flentities.User, error) {
	user := &flentities.User{}
	return user, withDatabase(func(db *gorm.DB) error {
		err := db.Where(flentities.User{FacebookID: &facebookID}).First(user).Error

		if gorm.IsRecordNotFoundError(err) {
			user, err = r.findByEmail(email)
		}

		if err != nil && !gorm.IsRecordNotFoundError(err) {
			return err
		}

		return nil
	})
}

func (r *usersRepository) findByEmail(email string) (*flentities.User, error) {
	user := &flentities.User{}
	return user, withDatabase(func(db *gorm.DB) error {
		err := db.Where(flentities.User{Email: email}).First(user).Error

		if err != nil && !gorm.IsRecordNotFoundError(err) {
			return err
		}

		return nil
	})
}

func (r *usersRepository) save(user *flentities.User) error {
	return withDatabase(func(db *gorm.DB) error {
		return db.Save(user).Error
	})
}

func (r *usersRepository) signUp(email, password string) (u *flentities.User, err error) {
	err = withDatabase(func(db *gorm.DB) error {
		u = &flentities.User{Email: email}
		if err := u.SetPassword(password); err != nil {
			return err
		}

		return db.Create(u).Error
	})
	return
}
