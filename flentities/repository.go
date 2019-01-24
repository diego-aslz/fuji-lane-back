package flentities

import (
	"github.com/jinzhu/gorm"
)

// Repository provides useful methods for common database operations. It also allows bypassing directly to the
// persistence framework for simpler queries
type Repository struct {
	*gorm.DB
}

// FindUserByEmail tries to find a User by its Email. Returns `nil` if it cannot find it
func (r *Repository) FindUserByEmail(email string) (*User, error) {
	user := &User{}
	err := r.Where(User{Email: email}).First(user).Error

	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return nil, err
	}

	return user, nil
}

// FindUserForFacebookSignIn tries to find the User by its FacebookID or Email. Returns `nil` if it cannot find it
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

// SignUp registers a new user for the given credentials
func (r *Repository) SignUp(email, password string) (u *User, err error) {
	u = &User{Email: email}
	if err := u.SetPassword(password); err != nil {
		return nil, err
	}

	return u, r.Create(u).Error
}

// RemoveImage removes an image and nullifies references
func (r *Repository) RemoveImage(image *Image) (err error) {
	r.Transaction(func(t *Repository) {
		err = t.Table("units").Where(map[string]interface{}{"floor_plan_image_id": image.ID}).
			Updates(map[string]interface{}{"floor_plan_image_id": nil}).Error
		if err != nil {
			t.Rollback()
			return
		}

		err = t.Delete(image).Error
		if err != nil {
			t.Rollback()
			return
		}

		err = t.Commit().Error
	})
	return
}

// UserProperties returns a *gorm.DB which selects all Properties the given user owns
func (r *Repository) UserProperties(u *User) *gorm.DB {
	return r.Table("properties").Where(map[string]interface{}{"account_id": *u.AccountID})
}

// Paginate to paginate query
func (r Repository) Paginate(page, perPage int) *gorm.DB {
	return r.Limit(perPage).Offset((page - 1) * perPage)
}

// Transaction calls the callback function with a transactional Repository. Any panics will be rolled back
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
