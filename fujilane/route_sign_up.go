package fujilane

import (
	"errors"
	"net/http"
)

type signUpBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (b *signUpBody) Validate() []error {
	return ValidateFields(
		ValidateField("email", b.Email).Required().Email(),
		ValidateField("password", b.Password).MinLen(8).MaxLen(30),
	)
}

func (a *Application) routeSignUp(c *routeContext) {
	body := &signUpBody{}
	if !c.parseBodyAndValidate(body) {
		return
	}

	safeBody := *body
	safeBody.Password = "[FILTERED]"
	c.addLogJSON("params", safeBody)

	user, err := a.usersRepository.signUp(body.Email, body.Password)
	if err != nil {
		if isUniqueConstraintViolation(err) {
			err = errors.New("Invalid email: already in use")
			c.addLogError(err)
			c.respond(http.StatusUnprocessableEntity, c.errorsBody([]error{err}))
		} else {
			c.fail(http.StatusInternalServerError, err)
		}
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

	c.respond(http.StatusCreated, s)
}
