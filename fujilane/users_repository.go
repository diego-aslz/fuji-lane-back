package fujilane

import (
	"github.com/jinzhu/gorm"
)

type usersRepository struct{}

func (r *usersRepository) findForFacebookSignIn(facebookID, email string, user *User) error {
	return withDatabase(func(db *gorm.DB) error {
		err := db.Where(User{FacebookID: facebookID}).First(user).Error

		if err == gorm.ErrRecordNotFound {
			err = db.Where(User{Email: email}).First(user).Error
		}

		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}

		user.FacebookID = facebookID

		if user.Email == "" {
			user.Email = email
		}

		return nil
	})
}

func (r *usersRepository) save(user *User) error {
	return withDatabase(func(db *gorm.DB) error {
		return db.Save(user).Error
	})
}
