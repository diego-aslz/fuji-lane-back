package flactions

import (
	"fmt"
	"net/http"
	"time"

	"github.com/nerde/fuji-lane-back/flconfig"

	"github.com/snabb/sitemap"
)

// PropertiesSitemap to generate Sitemap XML for properties
type PropertiesSitemap struct {
	Context
}

// Perform executes the action
func (a *PropertiesSitemap) Perform() {
	rows, err := a.Repository().
		Table("properties").
		Joins("INNER JOIN units ON properties.id = units.property_id AND units.published_at IS NOT NULL").
		Where("properties.published_at IS NOT NULL").
		Select("properties.id, properties.slug, MAX(GREATEST(properties.updated_at, units.updated_at))").
		Group("properties.id, properties.slug").
		Rows()

	if err != nil {
		a.ServerError(err)
		return
	}
	defer rows.Close()

	sm := sitemap.New()

	for rows.Next() {
		slug := ""
		id := 0
		var updatedAt time.Time

		if err := rows.Scan(&id, &slug, &updatedAt); err != nil {
			a.ServerError(err)
			return
		}

		sm.Add(&sitemap.URL{
			Loc:        fmt.Sprintf("%s/listings/%s-%d", flconfig.Config.AppURL, slug, id),
			LastMod:    &updatedAt,
			ChangeFreq: sitemap.Daily,
		})
	}

	a.RespondXML(http.StatusOK, sm)
}

// NewPropertiesSitemap returns a new PropertiesSitemap action
func NewPropertiesSitemap(c Context) Action {
	return &PropertiesSitemap{c}
}
