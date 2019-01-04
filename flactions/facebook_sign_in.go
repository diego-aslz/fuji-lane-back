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
	Context
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

	user, err := a.Repository().FindUserForFacebookSignIn(a.ID, a.Email)
	if err != nil {
		a.ServerError(err)
		return
	}

	now := a.Now()
	if user.ID > 0 {
		updates := flentities.User{Name: &a.Name, FacebookID: &a.ID, LastSignedIn: &now}
		err = a.Repository().Model(user).Updates(updates).Error
	} else {
		user.Email = a.Email
		user.Name = &a.Name
		user.FacebookID = &a.ID
		user.LastSignedIn = &now
		err = a.Repository().Create(user).Error
	}

	if err != nil {
		a.ServerError(err)
		return
	}

	createSession(a.Context, user)
}

// NewFacebookSignIn creates a new FacebookSignIn instance
func NewFacebookSignIn(client flservices.FacebookClient, c Context) Action {
	return &FacebookSignIn{facebook: flservices.NewFacebook(client), Context: c}
}
