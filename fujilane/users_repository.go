package fujilane

import (
	"github.com/jinzhu/gorm"
)

type usersRepository struct{}

func (r *usersRepository) findForFacebookSignIn(facebookID, email string, user *User) error {
	return withDatabase(func(db *gorm.DB) error {
		err := db.Where(User{FacebookID: facebookID}).First(user).Error

		if gorm.IsRecordNotFoundError(err) {
			err = db.Where(User{Email: email}).First(user).Error
		}

		if err != nil && !gorm.IsRecordNotFoundError(err) {
			return err
		}

		return nil
	})
}

func (r *usersRepository) save(user *User) error {
	return withDatabase(func(db *gorm.DB) error {
		return db.Save(user).Error
	})
}

func (r *usersRepository) signUp(email, password string) (u *User, err error) {
	err = withDatabase(func(db *gorm.DB) error {
		u = &User{Email: email}
		if err := u.setPassword(password); err != nil {
			return err
		}

		return db.Create(u).Error
	})
	return
}
