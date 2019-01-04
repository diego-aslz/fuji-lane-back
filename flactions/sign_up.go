package flactions

import (
	"errors"
	"net/http"

	"github.com/nerde/fuji-lane-back/fldiagnostics"
	"github.com/nerde/fuji-lane-back/flentities"
)

// SignUpBody is the expected payload for SignUp
type SignUpBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Validate the request body
func (b *SignUpBody) Validate() []error {
	return flentities.ValidateFields(
		flentities.ValidateField("email", b.Email).Required().Email(),
		flentities.ValidateField("password", b.Password).MinLen(8).MaxLen(30),
	)
}

// FilterSensitiveInformation hides the user password
func (b SignUpBody) FilterSensitiveInformation() fldiagnostics.SensitivePayload {
	b.Password = "[FILTERED]"
	return b
}

// SignUp creates properties that can hold units
type SignUp struct {
	SignUpBody
	sessionAction
}

// Perform executes the action
func (a *SignUp) Perform() {
	user, err := a.Repository().SignUp(a.Email, a.Password)
	if err != nil {
		if flentities.IsUniqueConstraintViolation(err) {
			err = errors.New("Invalid email: already in use")
			a.Diagnostics().AddError(err)
			a.RespondError(http.StatusUnprocessableEntity, err)
		} else {
			a.ServerError(err)
		}
		return
	}

	now := a.Now()
	if err = a.Repository().Model(user).Updates(flentities.User{LastSignedIn: &now}).Error; err != nil {
		a.ServerError(err)
		return
	}

	a.createSessionWithStatus(user, http.StatusCreated)
}

// NewSignUp returns a new SignUp action
func NewSignUp(c Context) Action {
	return &SignUp{sessionAction: sessionAction{c}}
}
