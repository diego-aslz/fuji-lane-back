package fujilane

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (a *Application) routeFacebookSignIn(c *gin.Context) {
	body := &facebookSignInBody{}
	err := c.BindJSON(body)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	err = a.facebook.validate(body.AccessToken, body.ID)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	user := &User{}
	err = a.usersRepository.findForFacebookSignIn(body.ID, body.Email, user)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
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
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	s := newSession(user.Email, a.timeFunc)
	if err = s.generateToken(); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, s)
}
