package flactions

import (
	"net/http"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/nerde/fuji-lane-back/flentities"
)

// UnitsShow exposes details for a unit
type UnitsShow struct{}

// Perform executes the action
func (a *UnitsShow) Perform(c Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Diagnostics().AddError(err)
		c.RespondNotFound()
		return
	}

	user := c.CurrentUser()
	propertyConditions := map[string]interface{}{"account_id": *user.AccountID}
	userProperties := c.Repository().Table("properties").Where(propertyConditions)
	conditions := map[string]interface{}{"id": id}

	unit := &flentities.Unit{}
	err = c.Repository().Preload("FloorPlanImage").Where(conditions).Where(
		"property_id IN (?)", userProperties.Select("id").QueryExpr()).Find(unit).Error

	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			c.RespondNotFound()
		} else {
			c.ServerError(err)
		}
		return
	}

	c.Respond(http.StatusOK, unit)
}
