package flactions

import (
	"net/http"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/nerde/fuji-lane-back/flentities"
)

// UnitsPublish marks a property as published, allowing it to appear in search results
type UnitsPublish struct {
	Context
}

// Perform executes the action
func (a *UnitsPublish) Perform() {
	user := a.CurrentUser()

	id, err := strconv.Atoi(a.Param("id"))
	if err != nil {
		a.Diagnostics().AddError(err)
		a.RespondNotFound()
		return
	}

	properties := a.Repository().UserProperties(user).Select("id")

	unit := &flentities.Unit{}
	err = a.Repository().Preload("Amenities").Preload("Prices").Preload("Images", flentities.Image{Uploaded: true}).
		Where(map[string]interface{}{"id": id}).Where("property_id IN (?)", properties.QueryExpr()).Find(unit).Error

	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			a.RespondNotFound()
		} else {
			a.ServerError(err)
		}
		return
	}

	errs := unit.CanBePublished()
	if len(errs) > 0 {
		a.RespondError(http.StatusUnprocessableEntity, errs...)
		return
	}

	updates := map[string]interface{}{"PublishedAt": a.Now(), "EverPublished": true}
	if err = a.Repository().Model(unit).Updates(updates).Error; err != nil {
		a.ServerError(err)
		return
	}

	a.Respond(http.StatusOK, nil)
}

// NewUnitsPublish returns a new UnitsPublish action
func NewUnitsPublish(c Context) Action {
	return &UnitsPublish{c}
}
