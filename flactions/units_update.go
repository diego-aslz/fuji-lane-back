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
}

func (b *UnitsUpdateBody) toMap() (updates map[string]interface{}) {
	updates = map[string]interface{}{}

	updates["name"] = b.Name
	updates["bedrooms"] = b.Bedrooms
	updates["sizeM2"] = b.SizeM2
	updates["count"] = b.Count

	optionals := map[string]uint{
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

// Validate the request body
func (b *UnitsUpdateBody) Validate() []error {
	return flentities.ValidateFields(
		flentities.ValidateField("name", b.Name).Required(),
		flentities.ValidateField("bedrooms", b.Bedrooms).Required(),
		flentities.ValidateField("size", b.SizeM2).Required(),
		flentities.ValidateField("number of unit type", b.Count).Required(),
	)
}

// UnitsUpdate creates a new Unit
type UnitsUpdate struct {
	UnitsUpdateBody
}

// Perform executes the action
func (a *UnitsUpdate) Perform(c Context) {
	unit := &flentities.Unit{}

	conditions := map[string]interface{}{"id": c.Param("id")}
	err := c.Repository().Preload("Property").Find(unit, conditions).Error
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

	if err := c.Repository().Model(unit).Updates(updates).Error; err != nil {
		c.ServerError(err)
		return
	}

	c.Respond(http.StatusOK, unit)
}
