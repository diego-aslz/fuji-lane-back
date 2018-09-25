package fujilane

import (
	"errors"
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/nerde/fuji-lane-back/flentities"
)

type accountCreateBody struct {
	UserName  string `json:"user_name"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	CountryID int    `json:"country_id"`
}

type accountCreateAction struct {
	accountCreateBody
}

func (a *accountCreateAction) perform(c *routeContext, db *gorm.DB) error {
	user := c.currentUser()

	if user.AccountID != nil {
		c.respondError(http.StatusUnprocessableEntity, errors.New("You already have an account"))
		return nil
	}

	account := a.buildAccount()

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	err := tx.Create(account).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Model(user).Updates(map[string]interface{}{"name": a.UserName, "account_id": account.ID}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	c.respond(http.StatusCreated, account)

	return tx.Commit().Error
}

func (a *accountCreateAction) buildAccount() *flentities.Account {
	account := &flentities.Account{Name: a.Name}
	if a.CountryID > 0 {
		account.CountryID = &a.CountryID
	}

	if a.Phone != "" {
		account.Phone = &a.Phone
	}

	return account
}

func (a *Application) routeAccountsCreate(c *routeContext) {
	action := &accountCreateAction{}
	if !c.parseBodyOrFail(action) {
		return
	}

	err := withDatabase(func(db *gorm.DB) error {
		return action.perform(c, db)
	})

	if err != nil {
		c.fatal(err)
	}
}
