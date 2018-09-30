package flentities

import (
	"github.com/jinzhu/gorm"
)

type Repository struct {
	*gorm.DB
}

func (r *Repository) FindUserByEmail(email string) (*User, error) {
	user := &User{}
	err := r.Where(User{Email: email}).First(user).Error

	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return nil, err
	}

	return user, nil
}

func (r *Repository) FindUserForFacebookSignIn(facebookID, email string) (*User, error) {
	user := &User{}
	err := r.Where(User{FacebookID: &facebookID}).First(user).Error

	if gorm.IsRecordNotFoundError(err) {
		user, err = r.FindUserByEmail(email)
	}

	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return nil, err
	}

	return user, nil
}

func (r *Repository) SignUp(email, password string) (u *User, err error) {
	u = &User{Email: email}
	if err := u.SetPassword(password); err != nil {
		return nil, err
	}

	return u, r.Create(u).Error
}

func (r *Repository) Transaction(fn func(*Repository)) {
	tx := r.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	fn(&Repository{tx})
}
