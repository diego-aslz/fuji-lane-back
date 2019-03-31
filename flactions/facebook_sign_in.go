package flactions

import (
	"errors"
	"net/http"

	"github.com/nerde/fuji-lane-back/fldiagnostics"
	"github.com/nerde/fuji-lane-back/flservices"

	"github.com/nerde/fuji-lane-back/flentities"
)

// FacebookSignInBody is the expected request body for FacebookSignIn
type FacebookSignInBody struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	Name        string `json:"name"`
	AccessToken string `json:"accessToken"`
}

// FilterSensitiveInformation hides the access token
func (b FacebookSignInBody) FilterSensitiveInformation() fldiagnostics.SensitivePayload {
	b.AccessToken = "[FILTERED]"
	return b
}

// FacebookSignIn signs the user in via Facebook authentication
type FacebookSignIn struct {
	FacebookSignInBody
	sessionAction
	facebook *flservices.Facebook
}

// Perform the action
func (a *FacebookSignIn) Perform() {
	err := a.facebook.Validate(a.AccessToken, a.ID)
	if err != nil {
		a.Diagnostics().AddError(err)
		a.RespondError(http.StatusUnauthorized, errors.New("You could not be authenticated"))
		return
	}

	users := flentities.UsersRepository{Repository: a.Repository()}
	user, err := users.FacebookSignIn(a.claims(), a.Now())
	if err != nil {
		a.ServerError(err)
		return
	}

	a.createSession(user)
}

func (a *FacebookSignIn) claims() map[string]string {
	return map[string]string{
		"email":      a.Email,
		"name":       a.Name,
		"facebookID": a.ID,
	}
}

// NewFacebookSignIn creates a new FacebookSignIn instance
func NewFacebookSignIn(client flservices.FacebookClient, c Context) Action {
	return &FacebookSignIn{facebook: flservices.NewFacebook(client), sessionAction: sessionAction{c}}
}
