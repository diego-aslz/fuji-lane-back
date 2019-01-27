package flactions

import (
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/nerde/fuji-lane-back/flentities"
)

// PropertiesList lists user properties
type PropertiesList struct {
	Context
}

// Perform executes the action
func (a *PropertiesList) Perform() {
	user := a.CurrentUser()

	properties := []*flentities.Property{}
	err := a.Repository().Order("name").
		Preload("Images", flentities.Image{Uploaded: true}, flentities.ImagesDefaultOrder).
		Preload("Units", func(db *gorm.DB) *gorm.DB {
			return db.Joins("LEFT JOIN prices ON prices.unit_id = units.id AND prices.min_nights = 1").
				Order(flentities.PerNightPriceSQL)
		}).
		Preload("Units.Prices").
		Find(&properties, map[string]interface{}{"account_id": user.AccountID}).Error
	if err != nil {
		a.ServerError(err)
		return
	}

	a.Respond(http.StatusOK, properties)
}

// NewPropertiesList returns a new PropertiesList action
func NewPropertiesList(c Context) Action {
	return &PropertiesList{c}
}
