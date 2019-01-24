package flactions

import (
	"errors"
	"net/http"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/nerde/fuji-lane-back/flentities"
)

// UnitPriceBody is the payload for a single unit price
type UnitPriceBody struct {
	MinNights int `json:"minNights"`
	Cents     int `json:"cents"`
}

// UnitsUpdateBody is the representation of the payload for creating a Unit
type UnitsUpdateBody struct {
	Name             string          `json:"name"`
	Overview         *string         `json:"overview"`
	Bedrooms         int             `json:"bedrooms"`
	Bathrooms        int             `json:"bathrooms"`
	SizeM2           int             `json:"sizeM2"`
	MaxOccupancy     *int            `json:"maxOccupancy"`
	Count            int             `json:"count"`
	FloorPlanImageID *uint           `json:"floorPlanImageID"`
	Prices           []UnitPriceBody `json:"prices"`
	bodyWithAmenities
}

// Validate the request body
func (b *UnitsUpdateBody) Validate() []error {
	return flentities.ValidateFields(
		flentities.ValidateField("overview", b.Overview).HTML(),
	)
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
	err := a.Repository().Preload("Property").Preload("Amenities").Preload("Prices").Find(unit, conditions).Error
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

	if a.Name != "" {
		unit.Name = a.Name
	}

	if a.Overview != nil {
		unit.Overview = a.Overview
	}

	if a.Bedrooms > 0 {
		unit.Bedrooms = a.Bedrooms
	}

	if a.Bathrooms > 0 {
		unit.Bathrooms = a.Bathrooms
	}

	if a.SizeM2 > 0 {
		unit.SizeM2 = a.SizeM2
	}

	if a.Count > 0 {
		unit.Count = a.Count
	}

	if a.MaxOccupancy != nil {
		unit.MaxOccupancy = a.MaxOccupancy
	}

	if a.FloorPlanImageID != nil {
		image := &flentities.Image{}
		image.ID = *a.FloorPlanImageID

		if err := a.Repository().Preload("Unit.Property").Find(image).Error; err != nil {
			a.ServerError(err)
			return
		}

		if image.Unit == nil || image.Unit.Property.AccountID != a.CurrentAccount().ID {
			a.RespondError(http.StatusUnprocessableEntity, errors.New("floor plan image does not exist"))
			return
		}

		unit.FloorPlanImageID = a.FloorPlanImageID
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

		if a.Prices != nil && len(a.Prices) > 0 {
			minNightsFound := []int{}

			for _, pBody := range a.Prices {
				if pBody.Cents == 0 {
					continue
				}

				var price *flentities.Price

				for _, p := range unit.Prices {
					if p.MinNights == pBody.MinNights {
						price = p
						break
					}
				}

				if price == nil {
					price = &flentities.Price{Unit: unit}
				}

				price.MinNights = pBody.MinNights
				price.Cents = pBody.Cents

				if err := tx.Save(price).Error; err != nil {
					tx.Rollback()
					a.ServerError(err)
					return
				}

				minNightsFound = append(minNightsFound, price.MinNights)
			}

			tx.Where("unit_id = ? AND min_nights NOT IN (?)", unit.ID, minNightsFound).Delete(&flentities.Price{})
		}

		if err := tx.Save(unit).Error; err != nil {
			tx.Rollback()

			if flentities.IsUniqueConstraintViolation(err) && strings.Index(err.Error(), "_slug") > -1 {
				a.RespondError(http.StatusUnprocessableEntity, errors.New("Name is already in use"))
				return
			}

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
