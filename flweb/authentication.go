package flweb

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/nerde/fuji-lane-back/flentities"
)

// CurrentAccount returns the account for the currently authenticated user
func (c *Context) CurrentAccount() *flentities.Account {
	if c.currentAccount != nil {
		return c.currentAccount
	}

	user := c.CurrentUser()
	if user == nil || user.AccountID == nil {
		return nil
	}

	c.currentAccount = &flentities.Account{}
	if err := c.Repository().First(c.currentAccount, *user.AccountID).Error; err != nil {
		c.Diagnostics().Add("current_account_load_error",
			fmt.Sprintf("Unable to load Account %d: %s", *user.AccountID, err.Error()))
		c.currentAccount = nil
	}

	return c.currentAccount
}

// CurrentUser returns the currently authenticated user
func (c *Context) CurrentUser() *flentities.User {
	return c.currentUser
}

// CurrentSession returns the session we loaded from authentication token
func (c *Context) CurrentSession() *flentities.Session {
	return c.session
}

func (c *Context) unauthorized(msg ...string) {
	m := "You need to sign in"
	if len(msg) > 0 {
		m = msg[0]
	}

	c.RespondError(http.StatusUnauthorized, errors.New(m))
}

func loadSession(next contextFunc) contextFunc {
	return func(c *Context) {
		auth := c.getHeader("Authorization")

		if auth == "" {
			c.Diagnostics().AddQuoted("session_info", "No authentication token provided")
		} else {
			auth = strings.TrimPrefix(auth, "Bearer ")
			session, err := flentities.LoadSession(auth)
			if err != nil {
				c.Diagnostics().AddJSON("token", auth).AddQuoted("reason", "Unable to load session from token")
				c.ServerError(err)
				return
			}

			if session.ExpiresAt.Before(c.Now()) {
				c.Diagnostics().AddJSON("session", session)
				c.unauthorized("Your session expired")
				return
			}

			user, err := c.repository.FindUserByEmail(session.Email)
			if err != nil {
				c.Diagnostics().AddJSON("session", session).AddQuoted("reason", "Unable to load user")
				c.ServerError(err)
				return
			}

			if user == nil || user.ID == 0 {
				c.Diagnostics().AddJSON("session", session).AddQuoted("session_warn", "User not found")
			} else {
				c.session = session
				c.currentUser = user

				c.Diagnostics().Add("user_email", user.Email).Add("user_id", fmt.Sprint(user.ID))

				accID := ""
				if user.AccountID != nil {
					accID = fmt.Sprint(*user.AccountID)
				}
				c.Diagnostics().Add("account_id", accID)
			}
		}

		next(c)
	}
}

func requireUser(next contextFunc) contextFunc {
	return func(c *Context) {
		if c.currentUser == nil {
			c.unauthorized()
			return
		}

		next(c)
	}
}
