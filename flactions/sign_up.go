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
}

// Perform executes the action
func (a *SignUp) Perform(c Context) {
	c.Diagnostics().AddSensitive("params", a)

	user, err := c.Repository().SignUp(a.Email, a.Password)
	if err != nil {
		if flentities.IsUniqueConstraintViolation(err) {
			err = errors.New("Invalid email: already in use")
			c.Diagnostics().AddError(err)
			c.RespondError(http.StatusUnprocessableEntity, err)
		} else {
			c.ServerError(err)
		}
		return
	}

	now := c.Now()
	if err = c.Repository().Model(user).Updates(flentities.User{LastSignedIn: &now}).Error; err != nil {
		c.ServerError(err)
		return
	}

	s := flentities.NewSession(user, c.Now)
	if err = s.GenerateToken(); err != nil {
		c.ServerError(err)
		return
	}

	c.Respond(http.StatusCreated, s)
}
