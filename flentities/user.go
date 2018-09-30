package flentities

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// User of the system
type User struct {
	gorm.Model
	AccountID         *int
	Account           *Account
	Name              *string
	Email             string
	FacebookID        *string
	EncryptedPassword *string
	LastSignedIn      time.Time
}

// SetPassword calculates the encrypted hash and fills in EncryptedPassword
func (u *User) SetPassword(password string) error {
	passwordBytes, e := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if e != nil {
		return e
	}

	str := string(passwordBytes)
	u.EncryptedPassword = &str

	return nil
}

// ValidatePassword returns true if the parameterized password is correct
func (u *User) ValidatePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(*u.EncryptedPassword), []byte(password)) == nil
}

// Picture returns an URL to the user's profile picture
func (u *User) Picture() string {
	if u.FacebookID != nil {
		return fmt.Sprintf("https://graph.facebook.com/%s/picture?width=64&height=64", *u.FacebookID)
	}

	return ""
}