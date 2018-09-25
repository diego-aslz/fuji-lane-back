package fujilane

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

func (c *routeContext) currentUser() *User {
	v, _ := c.context.Get("current-user")
	return v.(*User)
}

func (c *routeContext) unauthorized() {
	c.respondError(http.StatusUnauthorized, errors.New("You need to sign in"))
}

func (a *Application) authenticateUser(next func(*routeContext)) func(*routeContext) {
	return func(c *routeContext) {
		auth := c.getHeader("Authorization")
		if auth == "" {
			c.addLogQuoted("reason", "Missing authentication token")
			c.unauthorized()
			return
		}

		auth = strings.TrimPrefix(auth, "Bearer ")
		session, err := loadSession(auth)
		if err != nil {
			c.addLogJSON("token", auth)
			c.addLogQuoted("reason", "Unable to load session from token")
			c.fatal(err)
			return
		}

		if session.ExpiresAt.Before(a.timeFunc()) {
			c.addLogFiltered("session", session)
			c.respondError(http.StatusUnauthorized, errors.New("Your session expired"))
			return
		}

		user, err := a.usersRepository.findByEmail(session.Email)
		if err != nil {
			c.addLogFiltered("session", session)
			c.addLogQuoted("reason", "Unable to load user")
			c.fatal(err)
			return
		}

		if user == nil || user.ID == 0 {
			c.addLogFiltered("session", session)
			c.addLogQuoted("reason", "User not found")
			c.unauthorized()
			return
		}

		c.addLog("user", user.Email)
		c.addLog("user_id", fmt.Sprint(user.ID))
		c.set("current-user", user)

		next(c)
	}
}
