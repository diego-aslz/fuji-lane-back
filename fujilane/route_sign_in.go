package fujilane

import (
	"errors"
	"net/http"
)

type signInBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

const authenticationFailedMessage = "Invalid email or password"

func (a *Application) routeSignIn(c *routeContext) {
	body := &signInBody{}
	if !c.parseBodyOrFail(body) {
		return
	}

	safeBody := *body
	safeBody.Password = "[FILTERED]"
	c.addLogJSON("params", safeBody)

	user, err := a.usersRepository.findByEmail(body.Email)
	if err != nil {
		c.fail(http.StatusInternalServerError, err)
		return
	}

	if user == nil || user.ID == 0 {
		c.addLogQuoted("reason", "User not found")
		c.respond(http.StatusUnauthorized, c.errorsBody([]error{errors.New(authenticationFailedMessage)}))
		return
	}

	if !user.validatePassword(body.Password) {
		c.addLogQuoted("reason", "Password is invalid")
		c.respond(http.StatusUnauthorized, c.errorsBody([]error{errors.New(authenticationFailedMessage)}))
		return
	}

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
