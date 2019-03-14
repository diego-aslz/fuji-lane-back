package flactions

import (
	"net/http"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/nerde/fuji-lane-back/flentities"
)

// PropertiesUnpublish hides the property so it won't show up on search results anymore.
type PropertiesUnpublish struct {
	Context
}

// Perform executes the action
func (a *PropertiesUnpublish) Perform() {
	user := a.CurrentUser()

	id, err := strconv.Atoi(a.Param("id"))
	if err != nil {
		a.Diagnostics().AddError(err)
		a.RespondNotFound()
		return
	}

	conditions := map[string]interface{}{
		"id":         id,
		"account_id": *user.AccountID,
	}

	property := &flentities.Property{}
	err = a.Repository().Find(property, conditions).Error

	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			a.RespondNotFound()
		} else {
			a.ServerError(err)
		}
		return
	}

	if property.PublishedAt != nil {
		updates := map[string]interface{}{"PublishedAt": nil}
		if err = a.Repository().Model(property).Updates(updates).Error; err != nil {
			a.ServerError(err)
			return
		}
	}

	a.Respond(http.StatusOK, nil)
}

// NewPropertiesUnpublish returns a new PropertiesUnpublish action
func NewPropertiesUnpublish(c Context) Action {
	return &PropertiesUnpublish{c}
}
