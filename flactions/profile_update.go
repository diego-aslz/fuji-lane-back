package flactions

import (
	"errors"
	"net/http"

	"github.com/nerde/fuji-lane-back/optional"

	"github.com/nerde/fuji-lane-back/fldiagnostics"
	"github.com/nerde/fuji-lane-back/flentities"
)

// ProfileUpdateBody is the expected payload for UsersUpdate
type ProfileUpdateBody struct {
	Name                     optional.String `json:"name"`
	Email                    optional.String `json:"email"`
	Password                 string          `json:"password"`
	ResetUnreadBookingsCount bool            `json:"resetUnreadBookingsCount"`
}

// Validate the request body
func (b *ProfileUpdateBody) Validate() []error {
	validations := []flentities.FieldValidation{}

	if b.Email.Set {
		validations = append(validations, flentities.ValidateField("email", b.Email.Value).Required().Email())
	}

	return flentities.ValidateFields(validations...)
}

// FilterSensitiveInformation hides password
func (b ProfileUpdateBody) FilterSensitiveInformation() fldiagnostics.SensitivePayload {
	b.Password = "[FILTERED]"
	return b
}

// ProfileUpdate is reponsible for creating new accounts
type ProfileUpdate struct {
	ProfileUpdateBody
	sessionAction
}

// Perform executes the action
func (a *ProfileUpdate) Perform() {
	user := a.CurrentUser()

	if a.Name.Set || a.Email.Set {
		if !user.ValidatePassword(a.Password) {
			a.RespondError(http.StatusUnauthorized, errors.New("Password is incorrect"))
			return
		}

		optional.Update(a.ProfileUpdateBody, user)

		if err := a.Repository().Save(user).Error; err != nil {
			a.ServerError(err)
			return
		}
	} else if a.ResetUnreadBookingsCount {
		user.UnreadBookingsCount = 0
		if err := a.Repository().UpdatesColVal(user, "UnreadBookingsCount", 0); err != nil {
			a.ServerError(err)
			return
		}
	}

	a.createSession(user)
}

// NewProfileUpdate returns a new ProfileUpdate action
func NewProfileUpdate(c Context) Action {
	return &ProfileUpdate{sessionAction: sessionAction{c}}
}
