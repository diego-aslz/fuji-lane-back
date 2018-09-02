package fujilane

import (
	"net/http"
)

type signUpBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (a *Application) routeSignUp(c *routeContext) {
	body := &signUpBody{}
	if !c.parseBodyOrFail(body) {
		return
	}

	user, err := a.usersRepository.signUp(body.Email, body.Password)
	if err != nil {
		c.fail(http.StatusInternalServerError, err)
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

	c.success(http.StatusCreated, s)
}
