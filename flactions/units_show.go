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
	properties := c.Repository().UserProperties(user).Select("id")

	unit := &flentities.Unit{}
	err = c.Repository().Preload("Amenities").Preload("Images", flentities.Image{Uploaded: true}, imagesDefaultOrder).
		Where(map[string]interface{}{"id": id}).Where("property_id IN (?)", properties.QueryExpr()).Find(unit).Error

	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			c.RespondNotFound()
		} else {
			c.ServerError(err)
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

	c.Respond(http.StatusOK, unit)
}
