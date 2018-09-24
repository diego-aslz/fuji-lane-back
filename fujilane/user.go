package fujilane

import (
	"time"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// User of the system
type User struct {
	gorm.Model
	Name              string
	Email             string
	FacebookID        string
	EncryptedPassword string
	LastSignedIn      time.Time
}

func (u *User) setPassword(password string) error {
	passwordBytes, e := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if e != nil {
		return e
	}

	u.EncryptedPassword = string(passwordBytes)

	return nil
}

func (u *User) validatePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.EncryptedPassword), []byte(password)) == nil
}
