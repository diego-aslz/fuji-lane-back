package fujilane

import (
	"errors"
	"net/http"
)

type facebookSignInBody struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	Name        string `json:"name"`
	AccessToken string `json:"accessToken"`
}

func (b facebookSignInBody) filterSensitiveInformation() filterableLog {
	b.AccessToken = "[FILTERED]"
	return b
}

func (a *Application) routeFacebookSignIn(c *routeContext) {
	body := &facebookSignInBody{}
	if !c.parseBodyOrFail(body) {
		return
	}

	c.addLogFiltered("params", body)

	err := a.facebook.validate(body.AccessToken, body.ID)
	if err != nil {
		c.addLogError(err)
		c.respondError(http.StatusUnauthorized, errors.New("You could not be authenticated"))
		return
	}

	user, err := a.usersRepository.findForFacebookSignIn(body.ID, body.Email)
	if err != nil {
		c.fatal(err)
		return
	}

	if user.Email == "" {
		user.Email = body.Email
	}
	user.Name = &body.Name
	user.FacebookID = &body.ID
	user.LastSignedIn = a.timeFunc()

	err = a.usersRepository.save(user)
	if err != nil {
		c.fatal(err)
		return
	}

	s := newSession(user, a.timeFunc)
	if err = s.generateToken(); err != nil {
		c.fatal(err)
		return
	}

	c.respond(http.StatusOK, s)
}
