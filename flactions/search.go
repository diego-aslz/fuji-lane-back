package flactions

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/nerde/fuji-lane-back/flentities"
	"github.com/nerde/fuji-lane-back/flviews"
)

// Search searches for listings matching the given filters
type Search struct {
	paginatedAction
}

// Perform the action
func (a *Search) Perform() {
	cityID := a.Query("cityID")
	if cityID == "" {
		a.RespondError(http.StatusBadRequest, errors.New("Please provide a City to filter by"))
		return
	}

	a.addPageDiagnostic()
	a.Diagnostics().Add("cityID", cityID)

	publishedNull := map[string]interface{}{"published_at": nil}
	unitConditions := a.Repository().Not(publishedNull).Where(map[string]interface{}{"deleted_at": nil})
	unitRawConditions := []string{"units.published_at IS NOT NULL", "units.deleted_at IS NULL"}
	unitJoinArgs := []interface{}{}

	rawBedrooms := a.Query("bedrooms")
	if rawBedrooms != "" {
		bedrooms, err := strconv.Atoi(rawBedrooms)

		if err == nil {
			condition := "bedrooms >= ?"
			unitConditions = unitConditions.Where(condition, bedrooms)
			unitRawConditions = append(unitRawConditions, condition)
			unitJoinArgs = append(unitJoinArgs, bedrooms)
		}
	}

	builder := a.Repository().
		Preload("Images", flentities.Image{Uploaded: true}, imagesDefaultOrder).
		Preload("Amenities").
		Preload("Units", func(_ *gorm.DB) *gorm.DB { return unitConditions.Order("base_price_cents") }).
		Preload("Units.Images", flentities.Image{Uploaded: true}, imagesDefaultOrder).
		Preload("Units.Amenities").
		Where("city_id = ?", cityID).
		Joins(fmt.Sprintf("INNER JOIN units ON properties.id = units.property_id AND %s",
			strings.Join(unitRawConditions, " AND ")), unitJoinArgs...).
		Select("DISTINCT(properties.*)")

	properties := []*flentities.Property{}
	if err := a.paginate(builder, a.getPage(), defaultPageSize).Find(&properties).Error; err != nil {
		a.ServerError(err)
		return
	}

	a.Diagnostics().Add("properties_size", strconv.Itoa(len(properties)))

	a.Respond(http.StatusOK, flviews.NewSearch(properties))
}

// NewSearch returns a new Search action
func NewSearch(c Context) Action {
	return &Search{paginatedAction{Context: c}}
}
