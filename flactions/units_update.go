package flactions

import (
	"errors"
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/nerde/fuji-lane-back/flentities"
)

// UnitsUpdateBody is the representation of the payload for creating a Unit
type UnitsUpdateBody struct {
	Name                   string  `json:"name"`
	Overview               *string `json:"overview"`
	Bedrooms               int     `json:"bedrooms"`
	Bathrooms              int     `json:"bathrooms"`
	SizeM2                 int     `json:"sizeM2"`
	MaxOccupancy           int     `json:"maxOccupancy"`
	Count                  int     `json:"count"`
	BasePriceCents         int     `json:"basePriceCents"`
	OneNightPriceCents     int     `json:"oneNightPriceCents"`
	OneWeekPriceCents      int     `json:"oneWeekPriceCents"`
	ThreeMonthsPriceCents  int     `json:"threeMonthsPriceCents"`
	SixMonthsPriceCents    int     `json:"sixMonthsPriceCents"`
	TwelveMonthsPriceCents int     `json:"twelveMonthsPriceCents"`
	FloorPlanImageID       uint    `json:"floorPlanImageID"`
	bodyWithAmenities
}

// Validate the request body
func (b *UnitsUpdateBody) Validate() []error {
	return flentities.ValidateFields(
		flentities.ValidateField("overview", b.Overview).HTML(),
	)
}

func (b *UnitsUpdateBody) toMap() (updates map[string]interface{}) {
	updates = map[string]interface{}{}

	if b.Name != "" {
		updates["name"] = b.Name
	}

	if b.Overview != nil {
		updates["Overview"] = *b.Overview
	}

	optionals := map[string]uint{
		"Bedrooms":               uint(b.Bedrooms),
		"Bathrooms":              uint(b.Bathrooms),
		"SizeM2":                 uint(b.SizeM2),
		"Count":                  uint(b.Count),
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
	Context
}

// Perform executes the action
func (a *UnitsUpdate) Perform() {
	unit := &flentities.Unit{}

	conditions := map[string]interface{}{"id": a.Param("id")}
	err := a.Repository().Preload("Property").Preload("Amenities").Find(unit, conditions).Error
	if gorm.IsRecordNotFoundError(err) {
		a.Diagnostics().AddQuoted("reason", "Could not find unit")
		a.RespondNotFound()
		return
	}
	if err != nil {
		a.ServerError(err)
		return
	}

	if unit.Property.AccountID != a.CurrentAccount().ID {
		a.Diagnostics().AddQuoted("reason", "Unit belongs to another account")
		a.RespondNotFound()
		return
	}

	updates := a.toMap()
	if imageID, ok := updates["FloorPlanImageID"]; ok {
		image := &flentities.Image{}
		image.ID = imageID.(uint)

		if err := a.Repository().Preload("Unit.Property").Find(image).Error; err != nil {
			a.ServerError(err)
			return
		}

		if image.Unit == nil || image.Unit.Property.AccountID != a.CurrentAccount().ID {
			a.RespondError(http.StatusUnprocessableEntity, errors.New("floor plan image does not exist"))
			return
		}
	}

	a.Repository().Transaction(func(tx *flentities.Repository) {
		amenitiesToDelete, amenitiesToCreate := a.amenitiesDiff(unit.Amenities)

		for _, am := range amenitiesToDelete {
			if err := tx.Delete(am).Error; err != nil {
				tx.Rollback()
				a.ServerError(err)
				return
			}
		}

		for _, am := range amenitiesToCreate {
			am.UnitID = &unit.ID
			if err := tx.Create(am).Error; err != nil {
				tx.Rollback()
				a.ServerError(err)
				return
			}
		}

		if err := tx.Model(unit).Updates(updates).Error; err != nil {
			tx.Rollback()
			a.ServerError(err)
			return
		}

		if err := tx.Commit().Error; err != nil {
			a.ServerError(err)
			return
		}

		a.Respond(http.StatusOK, unit)
	})
}

// NewUnitsUpdate returns a new UnitsUpdate action
func NewUnitsUpdate(c Context) Action {
	return &UnitsUpdate{Context: c}
}
