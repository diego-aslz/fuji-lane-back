package flweb

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/nerde/fuji-lane-back/flentities"
)

// CurrentAccount returns the account for the currently authenticated user
func (c *Context) CurrentAccount() *flentities.Account {
	user := c.CurrentUser()
	if user == nil || user.AccountID == nil {
		return nil
	}

	if user.Account == nil {
		user.Account = &flentities.Account{}
		if err := c.Repository().First(user.Account, *user.AccountID).Error; err != nil {
			c.Diagnostics().Add("current_account_load_error",
				fmt.Sprintf("Unable to load Account %d: %s", *user.AccountID, err.Error()))
			user.Account = nil
		}
	}

	return user.Account
}

// CurrentUser returns the currently authenticated user
func (c *Context) CurrentUser() *flentities.User {
	v, ok := c.Get("current-user")
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
		if user.AccountID != nil {
			c.Diagnostics().Add("account_id", strconv.Itoa(*user.AccountID))
		} else {
			c.Diagnostics().Add("account_id", "")
		}

		c.set("current-user", user)

		next(c)
	}
}
