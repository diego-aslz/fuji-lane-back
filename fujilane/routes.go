package fujilane

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	statusPath         = "/status"
	facebookSignInPath = "/sign_in/facebook"
)

// AddRoutes to a Gin Engine
func (a *Application) AddRoutes(e *gin.Engine) {
	e.GET(statusPath, a.statusRoute)
	e.POST(facebookSignInPath, a.facebookSignInRoute)
}

func (a *Application) statusRoute(c *gin.Context) {
	c.JSON(200, gin.H{"status": "active"})
}

type facebookSignInBody struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	Name        string `json:"name"`
	AccessToken string `json:"accessToken"`
}

func (a *Application) facebookSignInRoute(c *gin.Context) {
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

	user.Name = body.Name
	user.LastSignedIn = time.Now()

	err = a.usersRepository.save(user)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "active"})
}
