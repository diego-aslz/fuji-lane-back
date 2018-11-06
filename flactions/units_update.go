package flactions

import (
	"errors"
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/nerde/fuji-lane-back/flentities"
)

// UnitsUpdateBody is the representation of the payload for creating a Unit
type UnitsUpdateBody struct {
	Name                   string `json:"name"`
	Bedrooms               int    `json:"bedrooms"`
	SizeM2                 int    `json:"sizeM2"`
	MaxOccupancy           int    `json:"maxOccupancy"`
	Count                  int    `json:"count"`
	BasePriceCents         int    `json:"basePriceCents"`
	OneNightPriceCents     int    `json:"oneNightPriceCents"`
	OneWeekPriceCents      int    `json:"oneWeekPriceCents"`
	ThreeMonthsPriceCents  int    `json:"threeMonthsPriceCents"`
	SixMonthsPriceCents    int    `json:"sixMonthsPriceCents"`
	TwelveMonthsPriceCents int    `json:"twelveMonthsPriceCents"`
	FloorPlanImageID       uint   `json:"floorPlainImageID"`
	bodyWithAmenities
}

func (b *UnitsUpdateBody) toMap() (updates map[string]interface{}) {
	updates = map[string]interface{}{}

	if b.Name != "" {
		updates["name"] = b.Name
	}

	optionals := map[string]uint{
		"bedrooms":               uint(b.Bedrooms),
		"sizeM2":                 uint(b.SizeM2),
		"count":                  uint(b.Count),
		"MaxOccupancy":           uint(b.MaxOccupancy),
		"BasePriceCents":         uint(b.BasePriceCents),
		"OneNightPriceCents":     uint(b.OneNightPriceCents),
		"OneWeekPriceCents":      uint(b.OneWeekPriceCents),
		"ThreeMonthsPriceCents":  uint(b.ThreeMonthsPriceCents),
		"SixMonthsPriceCents":    uint(b.SixMonthsPriceCents),
		"TwelveMonthsPriceCents": uint(b.TwelveMonthsPriceCents),
		"FloorPlanImageID":       b.FloorPlanImageID,
	}

	for field, value := range optionals {
		if value == 0 {
			continue
		}
		updates[field] = value
	}

	return
}

// UnitsUpdate creates a new Unit
type UnitsUpdate struct {
	UnitsUpdateBody
}

// Perform executes the action
func (a *UnitsUpdate) Perform(c Context) {
	unit := &flentities.Unit{}

	conditions := map[string]interface{}{"id": c.Param("id")}
	err := c.Repository().Preload("Property").Preload("Amenities").Find(unit, conditions).Error
	if gorm.IsRecordNotFoundError(err) {
		c.Diagnostics().AddQuoted("reason", "Could not find unit")
		c.RespondNotFound()
		return
	}
	if err != nil {
		c.ServerError(err)
		return
	}

	if unit.Property.AccountID != c.CurrentAccount().ID {
		c.Diagnostics().AddQuoted("reason", "Unit belongs to another account")
		c.RespondNotFound()
		return
	}

	updates := a.toMap()
	if imageID, ok := updates["FloorPlanImageID"]; ok {
		image := &flentities.Image{}
		image.ID = imageID.(uint)

		if err := c.Repository().Preload("Unit.Property").Find(image).Error; err != nil {
			c.ServerError(err)
			return
		}

		if image.Unit == nil || image.Unit.Property.AccountID != c.CurrentAccount().ID {
			c.RespondError(http.StatusUnprocessableEntity, errors.New("floor plan image does not exist"))
			return
		}
	}

	c.Repository().Transaction(func(tx *flentities.Repository) {
		amenitiesToDelete, amenitiesToCreate := a.amenitiesDiff(unit.Amenities)

		for _, am := range amenitiesToDelete {
			if err := tx.Delete(am).Error; err != nil {
				tx.Rollback()
				c.ServerError(err)
				return
			}
		}

		for _, am := range amenitiesToCreate {
			am.UnitID = &unit.ID
			if err := tx.Create(am).Error; err != nil {
				tx.Rollback()
				c.ServerError(err)
				return
			}
		}

		if err := tx.Model(unit).Updates(updates).Error; err != nil {
			tx.Rollback()
			c.ServerError(err)
			return
		}

		if err := tx.Commit().Error; err != nil {
			c.ServerError(err)
			return
		}

		c.Respond(http.StatusOK, unit)
	})
}
