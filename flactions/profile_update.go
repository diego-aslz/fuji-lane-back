package flactions

import (
	"errors"
	"net/http"

	"github.com/nerde/fuji-lane-back/flentities"
)

// ProfileUpdateBody is the expected payload for UsersUpdate
type ProfileUpdateBody struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Validate the request body
func (b *ProfileUpdateBody) Validate() []error {
	return flentities.ValidateFields(
		flentities.ValidateField("email", b.Email).Required().Email(),
	)
}

// ProfileUpdate is reponsible for creating new accounts
type ProfileUpdate struct {
	ProfileUpdateBody
}

// Perform executes the action
func (a *ProfileUpdate) Perform(c Context) {
	user := c.CurrentUser()

	if !user.ValidatePassword(a.Password) {
		c.RespondError(http.StatusUnauthorized, errors.New("Password does not match"))
		return
	}

	if a.Name == "" {
		user.Name = nil
	} else {
		user.Name = &a.Name
	}

	user.Email = a.Email

	if err := c.Repository().Save(user).Error; err != nil {
		c.ServerError(err)
		return
	}

	c.Respond(http.StatusOK, user)
}
