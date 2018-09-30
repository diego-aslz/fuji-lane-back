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
	facebook *flservices.Facebook
}

// Perform the action
func (a *FacebookSignIn) Perform(c Context) {
	c.Diagnostics().AddSensitive("params", a)

	err := a.facebook.Validate(a.AccessToken, a.ID)
	if err != nil {
		c.Diagnostics().AddError(err)
		c.RespondError(http.StatusUnauthorized, errors.New("You could not be authenticated"))
		return
	}

	user, err := c.Repository().FindUserForFacebookSignIn(a.ID, a.Email)
	if err != nil {
		c.ServerError(err)
		return
	}

	if user.Email == "" {
		user.Email = a.Email
	}
	user.Name = &a.Name
	user.FacebookID = &a.ID
	user.LastSignedIn = c.Now()

	err = c.Repository().Save(user).Error
	if err != nil {
		c.ServerError(err)
		return
	}

	s := flentities.NewSession(user, c.Now)
	if err = s.GenerateToken(); err != nil {
		c.ServerError(err)
		return
	}

	c.Respond(http.StatusOK, s)
}

// NewFacebookSignIn creates a new FacebookSignIn instance
func NewFacebookSignIn(client flservices.FacebookClient) Action {
	return &FacebookSignIn{facebook: flservices.NewFacebook(client)}
}
