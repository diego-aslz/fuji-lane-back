package flactions

import (
	"net/http"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/nerde/fuji-lane-back/flentities"
)

// UnitsShow exposes details for a unit
type UnitsShow struct {
	Context
}

// Perform executes the action
func (a *UnitsShow) Perform() {
	id, err := strconv.Atoi(a.Param("id"))
	if err != nil {
		a.Diagnostics().AddError(err)
		a.RespondNotFound()
		return
	}

	user := a.CurrentUser()
	properties := a.Repository().UserProperties(user).Select("id")

	unit := &flentities.Unit{}
	err = a.Repository().Preload("Amenities").Preload("Images", flentities.Image{Uploaded: true}, imagesDefaultOrder).
		Where(map[string]interface{}{"id": id}).Where("property_id IN (?)", properties.QueryExpr()).Find(unit).Error

	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			a.RespondNotFound()
		} else {
			a.ServerError(err)
		}
		return
	}

	if unit.FloorPlanImageID != nil {
		filteredImages := unit.Images[:0]

		for _, img := range unit.Images {
			if img.ID == *unit.FloorPlanImageID {
				unit.FloorPlanImage = img
			} else {
				filteredImages = append(filteredImages, img)
			}
		}

		unit.Images = filteredImages
	}

	a.Respond(http.StatusOK, unit)
}

// NewUnitsShow returns a new UnitsShow action
func NewUnitsShow(c Context) Action {
	return &UnitsShow{c}
}
