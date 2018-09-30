package flactions

import (
	"errors"
	"net/http"

	"github.com/nerde/fuji-lane-back/flentities"
)

// AccountsCreateBody is the expected payload for AccountsCreate
type AccountsCreateBody struct {
	UserName  string `json:"user_name"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	CountryID int    `json:"country_id"`
}

// AccountsCreate is reponsible for creating new accounts
type AccountsCreate struct {
	AccountsCreateBody
}

// Perform executes the action
func (a *AccountsCreate) Perform(c Context) {
	user := c.CurrentUser()

	if user.AccountID != nil {
		c.RespondError(http.StatusUnprocessableEntity, errors.New("You already have an account"))
		return
	}

	account := a.buildAccount()

	c.Repository().Transaction(func(tx *flentities.Repository) {
		err := tx.Create(account).Error
		if err != nil {
			tx.Rollback()
			c.ServerError(err)
			return
		}

		err = tx.Model(user).Updates(map[string]interface{}{"name": a.UserName, "account_id": account.ID}).Error
		if err != nil {
			tx.Rollback()
			c.ServerError(err)
			return
		}

		if err = tx.Commit().Error; err != nil {
			c.ServerError(err)
			return
		}

		c.Respond(http.StatusCreated, account)
	})
}

func (a *AccountsCreate) buildAccount() *flentities.Account {
	account := &flentities.Account{Name: a.Name}
	if a.CountryID > 0 {
		account.CountryID = &a.CountryID
	}

	if a.Phone != "" {
		account.Phone = &a.Phone
	}

	return account
}

// NewAccountsCreate creates new AccountsCreate instances
func NewAccountsCreate() Action {
	return &AccountsCreate{}
}
