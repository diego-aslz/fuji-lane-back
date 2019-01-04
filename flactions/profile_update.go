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
	Context
}

// Perform executes the action
func (a *ProfileUpdate) Perform() {
	user := a.CurrentUser()

	if !user.ValidatePassword(a.Password) {
		a.RespondError(http.StatusUnauthorized, errors.New("Password does not match"))
		return
	}

	if a.Name == "" {
		user.Name = nil
	} else {
		user.Name = &a.Name
	}

	user.Email = a.Email

	if err := a.Repository().Save(user).Error; err != nil {
		a.ServerError(err)
		return
	}

	a.Respond(http.StatusOK, user)
}

// NewProfileUpdate returns a new ProfileUpdate action
func NewProfileUpdate(c Context) Action {
	return &ProfileUpdate{Context: c}
}
