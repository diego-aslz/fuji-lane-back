package flactions

import (
	"net/http"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/nerde/fuji-lane-back/flentities"
)

// UnitsUnpublish unpublishes a unit to hide it from search results
type UnitsUnpublish struct {
	Context
}

// Perform executes the action
func (a *UnitsUnpublish) Perform() {
	user := a.CurrentUser()

	id, err := strconv.Atoi(a.Param("id"))
	if err != nil {
		a.Diagnostics().AddError(err)
		a.RespondNotFound()
		return
	}

	properties := a.Repository().UserProperties(user).Select("id")

	unit := &flentities.Unit{}
	err = a.Repository().Where(map[string]interface{}{"id": id}).Where("property_id IN (?)", properties.QueryExpr()).
		Find(unit).Error

	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			a.RespondNotFound()
		} else {
			a.ServerError(err)
		}
		return
	}

	if unit.PublishedAt != nil {
		updates := map[string]interface{}{"PublishedAt": nil}
		if err = a.Repository().Model(unit).Updates(updates).Error; err != nil {
			a.ServerError(err)
			return
		}
	}

	a.Respond(http.StatusOK, nil)
}

// NewUnitsUnpublish returns a new UnitsUnpublish action
func NewUnitsUnpublish(c Context) Action {
	return &UnitsUnpublish{c}
}
