package fujilane

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/nerde/fuji-lane-back/flentities"
)

// CurrentUser returns the currently authenticated user
func (c *Context) CurrentUser() *flentities.User {
	v, ok := c.context.Get("current-user")
	if !ok {
		return nil
	}

	return v.(*flentities.User)
}

func (c *Context) unauthorized() {
	c.RespondError(http.StatusUnauthorized, errors.New("You need to sign in"))
}

func authenticateUser(next func(*Context)) func(*Context) {
	return func(c *Context) {
		auth := c.getHeader("Authorization")
		if auth == "" {
			c.Diagnostics().AddQuoted("reason", "Missing authentication token")
			c.unauthorized()
			return
		}

		auth = strings.TrimPrefix(auth, "Bearer ")
		session, err := flentities.LoadSession(auth)
		if err != nil {
			c.Diagnostics().AddJSON("token", auth).AddQuoted("reason", "Unable to load session from token")
			c.ServerError(err)
			return
		}

		if session.ExpiresAt.Before(c.Now()) {
			c.Diagnostics().AddSensitive("session", session)
			c.RespondError(http.StatusUnauthorized, errors.New("Your session expired"))
			return
		}

		user, err := c.repository.FindUserByEmail(session.Email)
		if err != nil {
			c.Diagnostics().AddSensitive("session", session).AddQuoted("reason", "Unable to load user")
			c.ServerError(err)
			return
		}

		if user == nil || user.ID == 0 {
			c.Diagnostics().AddSensitive("session", session).AddQuoted("reason", "User not found")
			c.unauthorized()
			return
		}

		c.Diagnostics().Add("user", user.Email).Add("user_id", fmt.Sprint(user.ID))
		c.set("current-user", user)

		next(c)
	}
}
