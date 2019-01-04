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
	sessionAction
}

const authenticationFailedMessage = "Invalid email or password"

// Perform executes the action
func (a *SignIn) Perform() {
	user, err := a.Repository().FindUserByEmail(a.Email)
	if err != nil {
		a.ServerError(err)
		return
	}

	if user == nil || user.ID == 0 {
		a.Diagnostics().AddQuoted("reason", "User not found")
		a.RespondError(http.StatusUnauthorized, errors.New(authenticationFailedMessage))
		return
	}

	if !user.ValidatePassword(a.Password) {
		a.Diagnostics().AddQuoted("reason", "Password is invalid")
		a.RespondError(http.StatusUnauthorized, errors.New(authenticationFailedMessage))
		return
	}

	now := a.Now()
	if err = a.Repository().Model(user).Updates(flentities.User{LastSignedIn: &now}).Error; err != nil {
		a.ServerError(err)
		return
	}

	a.createSession(user)
}

// NewSignIn returns a new SignIn action
func NewSignIn(c Context) Action {
	return &SignIn{sessionAction: sessionAction{c}}
}
