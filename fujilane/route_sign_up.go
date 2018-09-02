package fujilane

import (
	"fmt"
	"net/http"
	"regexp"
)

type signUpBody struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	errors          []string
	usersRepository *usersRepository
}

const (
	emailRegexString = "^(?:(?:(?:(?:[a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(?:\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|(?:(?:\\x22)(?:(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))?(?:\\x20|\\x09)+)?(?:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:\\(?:[\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))?(\\x20|\\x09)+)?(?:\\x22)))@(?:(?:(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])(?:[a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(?:(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])(?:[a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$"
)

func (b *signUpBody) validate() {
	b.errors = []string{}

	passwordReg := regexp.MustCompile(emailRegexString)
	if !passwordReg.Match([]byte(b.Email)) {
		b.errors = append(b.errors, fmt.Sprintf("Invalid email: %s", b.Email))
	}

	passLen := len(b.Password)
	if passLen < 8 || passLen > 30 {
		b.errors = append(b.errors, "Invalid password: length should be between 8 and 30")
	}

	inUse, err := b.usersRepository.isEmailInUse(b.Email)
	if err != nil {
		fmt.Printf("Unable to check email uniqueness: %s\n", err.Error())
		b.errors = append(b.errors, "Unable to validate email")
	} else if inUse {
		b.errors = append(b.errors, fmt.Sprintf("Invalid email: %s is already in use", b.Email))
	}
}

func (a *Application) routeSignUp(c *routeContext) {
	body := &signUpBody{usersRepository: a.usersRepository}
	if !c.parseBodyOrFail(body) {
		return
	}

	body.validate()

	if len(body.errors) > 0 {
		c.respond(http.StatusUnprocessableEntity, c.errorsBody(body.errors))
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

	c.respond(http.StatusCreated, s)
}
