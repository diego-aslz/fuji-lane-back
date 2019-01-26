package flactions

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/nerde/fuji-lane-back/flentities"
	"github.com/nerde/fuji-lane-back/flviews"
)

// Search searches for listings matching the given filters
type Search struct {
	paginatedAction
}

// Perform the action
func (a *Search) Perform() {
	filters := &flentities.ListingsSearchFilters{Page: a.getPage(), PerPage: defaultPageSize}

	a.withIntFilter("cityID", func(i int) { filters.CityID = uint(i) })

	if filters.CityID == 0 {
		a.RespondError(http.StatusBadRequest, errors.New("Please provide a City to filter by"))
		return
	}

	a.withIntFilter("bedrooms", func(i int) { filters.MinBedrooms = i })
	a.withIntFilter("bathrooms", func(i int) { filters.MinBathrooms = i })
	a.withIntFilter("minPriceCents", func(i int) { filters.MinPriceCents = i })
	a.withIntFilter("maxPriceCents", func(i int) { filters.MaxPriceCents = i })
	a.withDateFilter("checkIn", func(d flentities.Date) { filters.CheckIn = &d })
	a.withDateFilter("checkOut", func(d flentities.Date) { filters.CheckOut = &d })

	a.Diagnostics().AddJSON("filters", filters)

	properties, err := flentities.ListingsSearch{Repository: a.Repository(), ListingsSearchFilters: filters}.Search()
	if err != nil {
		a.ServerError(err)
		return
	}

	a.Diagnostics().Add("properties_size", strconv.Itoa(len(properties)))

	a.Respond(http.StatusOK, flviews.NewSearch(properties, filters.Nights()))
}

func (a *Search) withIntFilter(name string, callback func(int)) {
	a.withFilter(name, func(raw string) {
		i, err := strconv.Atoi(raw)

		if err != nil {
			a.addFilterErrorDiagnostic(name, raw, err)

			return
		}

		callback(i)
	})
}

func (a *Search) withDateFilter(name string, callback func(flentities.Date)) {
	a.withFilter(name, func(raw string) {
		d, err := flentities.ParseDate(raw)

		if err != nil {
			a.addFilterErrorDiagnostic(name, raw, err)

			return
		}

		callback(d)
	})
}

func (a *Search) withFilter(name string, callback func(string)) {
	raw := a.Query(name)
	if raw != "" {
		callback(raw)
	}
}

func (a *Search) addFilterErrorDiagnostic(name, raw string, err error) {
	a.Diagnostics().AddQuoted(fmt.Sprintf("%s_filter_error", name), fmt.Sprintf("Unable to parse %s: %s", raw,
		err.Error()))
}

// NewSearch returns a new Search action
func NewSearch(c Context) Action {
	return &Search{paginatedAction{Context: c}}
}
