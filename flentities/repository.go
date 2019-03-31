package flentities

import (
	"time"

	"github.com/jinzhu/gorm"
)

// Repository provides useful methods for common database operations. It also allows bypassing directly to the
// persistence framework for simpler queries
type Repository struct {
	*gorm.DB
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

// FindBy tries to find the first record matching the conditions and handles "not found" errors. Returns `true`
// if the record was found.
func (r *Repository) FindBy(model, conditions interface{}) (bool, error) {
	if err := r.Where(conditions).First(model).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
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

// UsersRepository handles User specific queries
type UsersRepository struct {
	*Repository
}

// FacebookSignIn signs an user in from Facebook Authentication request. If user does not exist, it will be created
func (r *UsersRepository) FacebookSignIn(claims map[string]string, now time.Time) (*User, error) {
	email := claims["email"]
	name := claims["name"]
	facebookID := claims["facebookID"]

	user, err := r.FindUserForFacebookSignIn(facebookID, email)
	if err != nil {
		return nil, err
	}

	if user.ID > 0 {
		return user, r.TrackFacebookAuth(user, name, facebookID, now)
	}
	user.Email = email
	user.Name = &name
	user.FacebookID = &facebookID
	user.LastSignedIn = &now

	return user, r.Create(user).Error
}

// FindUserForFacebookSignIn tries to find the User by its FacebookID or Email. Returns `nil` if it cannot find it
func (r *UsersRepository) FindUserForFacebookSignIn(facebookID, email string) (*User, error) {
	user := &User{}
	found, err := r.FindBy(user, User{FacebookID: &facebookID})
	if err != nil {
		return nil, err
	}

	if !found {
		found, err = r.FindBy(user, map[string]interface{}{"email": email})
		if err != nil {
			return nil, err
		}
	}

	return user, nil
}

// TrackFacebookAuth updates the given user with Facebook Auth details
func (r *UsersRepository) TrackFacebookAuth(user *User, name, facebookID string, now time.Time) error {
	updates := map[string]interface{}{"Name": name, "FacebookID": facebookID, "LastSignedIn": &now}
	return r.Model(user).Updates(updates).Error
}

// GoogleSignIn signs an user in from Google Authentication request. If user does not exist, it will be created
func (r *UsersRepository) GoogleSignIn(claims map[string]string, now time.Time) (*User, error) {
	user := &User{}
	found, err := r.FindBy(user, map[string]interface{}{"email": claims["email"]})
	if err != nil {
		return nil, err
	}

	if !found {
		return r.CreateFromGoogleAuth(claims, now)
	}

	return user, r.UpdateFromGoogleAuth(user, claims, now)
}

// CreateFromGoogleAuth creates a new user from Google Auth claims
func (r *UsersRepository) CreateFromGoogleAuth(claims map[string]string, now time.Time) (*User, error) {
	name := claims["name"]
	picture := claims["picture"]
	googleID := claims["googleID"]

	user := &User{
		Email:        claims["email"],
		Name:         &name,
		PictureURL:   &picture,
		GoogleID:     &googleID,
		LastSignedIn: &now,
	}

	return user, r.Save(user).Error
}

// UpdateFromGoogleAuth updates user attributes from a google auth
func (r *UsersRepository) UpdateFromGoogleAuth(user *User, claims map[string]string, now time.Time) error {
	return r.Model(user).Updates(map[string]interface{}{
		"Name":         claims["name"],
		"PictureURL":   claims["picture"],
		"GoogleID":     claims["googleID"],
		"LastSignedIn": &now,
	}).Error
}
