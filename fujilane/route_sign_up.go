package fujilane

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type signUpBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (a *Application) routeSignUp(c *gin.Context) {
	body := &signUpBody{}
	err := c.BindJSON(body)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	user, err := a.usersRepository.signUp(body.Email, body.Password)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	user.LastSignedIn = a.timeFunc()

	err = a.usersRepository.save(user)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	s := newSession(user, a.timeFunc)
	if err = s.generateToken(); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, s)
}
