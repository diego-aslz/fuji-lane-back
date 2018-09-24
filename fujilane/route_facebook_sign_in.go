package fujilane

import (
	"net/http"
)

type facebookSignInBody struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	Name        string `json:"name"`
	AccessToken string `json:"accessToken"`
}

func (a *Application) routeFacebookSignIn(c *routeContext) {
	body := &facebookSignInBody{}
	if !c.parseBodyOrFail(body) {
		return
	}

	safeBody := *body
	safeBody.AccessToken = "[FILTERED]"
	c.addLogJSON("params", safeBody)

	err := a.facebook.validate(body.AccessToken, body.ID)
	if err != nil {
		c.fail(http.StatusUnauthorized, err)
		return
	}

	user, err := a.usersRepository.findForFacebookSignIn(body.ID, body.Email)
	if err != nil {
		c.fail(http.StatusInternalServerError, err)
		return
	}

	if user.Email == "" {
		user.Email = body.Email
	}
	user.Name = body.Name
	user.FacebookID = body.ID
	user.LastSignedIn = a.timeFunc()

	err = a.usersRepository.save(user)
	if err != nil {
		c.fail(http.StatusInternalServerError, err)
		return
	}

	s := newSession(user, a.timeFunc)
	if err = s.generateToken(); err != nil {
		c.fail(http.StatusInternalServerError, err)
		return
	}

	c.respond(http.StatusOK, s)
}
