package flactions

import (
	"errors"
	"net/http"

	"github.com/nerde/fuji-lane-back/flentities"
)

// AccountsCreateBody is the expected payload for AccountsCreate
type AccountsCreateBody struct {
	UserName  string `json:"userName"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	CountryID uint   `json:"countryID"`
}

// AccountsCreate is reponsible for creating new accounts
type AccountsCreate struct {
	AccountsCreateBody
	Context
}

// Perform executes the action
func (a *AccountsCreate) Perform() {
	user := a.CurrentUser()

	if user.AccountID != nil {
		a.RespondError(http.StatusUnprocessableEntity, errors.New("You already have an account"))
		return
	}

	account := a.buildAccount()

	a.Repository().Transaction(func(tx *flentities.Repository) {
		err := tx.Create(account).Error
		if err != nil {
			tx.Rollback()
			a.ServerError(err)
			return
		}

		err = tx.Model(user).Updates(map[string]interface{}{"name": a.UserName, "account_id": account.ID}).Error
		if err != nil {
			tx.Rollback()
			a.ServerError(err)
			return
		}

		if err = tx.Commit().Error; err != nil {
			a.ServerError(err)
			return
		}

		a.Respond(http.StatusCreated, account)
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

// NewAccountsCreate returns a new AccountsCreate action
func NewAccountsCreate(c Context) Action {
	return &AccountsCreate{Context: c}
}
