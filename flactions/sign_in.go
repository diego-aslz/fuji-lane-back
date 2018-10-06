package flactions

import (
	"errors"
	"net/http"

	"github.com/nerde/fuji-lane-back/fldiagnostics"
	"github.com/nerde/fuji-lane-back/flentities"
)

// SignInBody is the expected payload for SignIn
type SignInBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// FilterSensitiveInformation hides the user password
func (b SignInBody) FilterSensitiveInformation() fldiagnostics.SensitivePayload {
	b.Password = "[FILTERED]"
	return b
}

// SignIn creates properties that can hold units
type SignIn struct {
	SignInBody
}

const authenticationFailedMessage = "Invalid email or password"

// Perform executes the action
func (a *SignIn) Perform(c Context) {
	c.Diagnostics().AddSensitive("params", a)

	user, err := c.Repository().FindUserByEmail(a.Email)
	if err != nil {
		c.ServerError(err)
		return
	}

	if user == nil || user.ID == 0 {
		c.Diagnostics().AddQuoted("reason", "User not found")
		c.RespondError(http.StatusUnauthorized, errors.New(authenticationFailedMessage))
		return
	}

	if !user.ValidatePassword(a.Password) {
		c.Diagnostics().AddQuoted("reason", "Password is invalid")
		c.RespondError(http.StatusUnauthorized, errors.New(authenticationFailedMessage))
		return
	}

	now := c.Now()
	if err = c.Repository().Model(user).Updates(flentities.User{LastSignedIn: &now}).Error; err != nil {
		c.ServerError(err)
		return
	}

	createSession(c, user)
}
