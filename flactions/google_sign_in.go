package flactions

import (
	"errors"
	"net/http"

	"github.com/nerde/fuji-lane-back/flentities"

	"github.com/nerde/fuji-lane-back/flservices"
)

// GoogleSignIn signs the user in via Google authentication
type GoogleSignIn struct {
	sessionAction
	googleAuth flservices.GoogleAuth
}

// Perform the action
func (a *GoogleSignIn) Perform() {
	rawToken := a.Context.GetHeader("Authorization")

	err := a.googleAuth.Verify(rawToken)
	if err != nil {
		a.Diagnostics().AddError(err)
		a.RespondError(http.StatusBadRequest, errors.New("Could not verify token"))
		return
	}

	claims, err := a.googleAuth.Decode(rawToken)
	if err != nil {
		a.ServerError(err)
		return
	}

	users := flentities.UsersRepository{Repository: a.Repository()}
	user, err := users.GoogleSignIn(claims, a.Now())
	if err != nil {
		a.ServerError(err)
		return
	}

	a.createSession(user)
}

// NewGoogleSignIn creates a new GoogleSignIn instance
func NewGoogleSignIn(googleAuth flservices.GoogleAuth, c Context) Action {
	return &GoogleSignIn{googleAuth: googleAuth, sessionAction: sessionAction{c}}
}
