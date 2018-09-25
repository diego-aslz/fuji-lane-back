package fujilane

import (
	"errors"
	"net/http"
)

type signInBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (b signInBody) filterSensitiveInformation() filterableLog {
	b.Password = "[FILTERED]"
	return b
}

const authenticationFailedMessage = "Invalid email or password"

func (a *Application) routeSignIn(c *routeContext) {
	body := &signInBody{}
	if !c.parseBodyOrFail(body) {
		return
	}

	c.addLogFiltered("params", body)

	user, err := a.usersRepository.findByEmail(body.Email)
	if err != nil {
		c.fatal(err)
		return
	}

	if user == nil || user.ID == 0 {
		c.addLogQuoted("reason", "User not found")
		c.respondError(http.StatusUnauthorized, errors.New(authenticationFailedMessage))
		return
	}

	if !user.validatePassword(body.Password) {
		c.addLogQuoted("reason", "Password is invalid")
		c.respondError(http.StatusUnauthorized, errors.New(authenticationFailedMessage))
		return
	}

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
