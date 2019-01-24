package flactions

import (
	"strconv"

	"github.com/nerde/fuji-lane-back/flentities"

	"github.com/jinzhu/gorm"
)

type paginatedAction struct {
	Context
	page int
}

func (a paginatedAction) getPage() int {
	if a.page > 0 {
		return a.page
	}

	rawPage := a.Query("page")
	if rawPage == "" {
		return 1
	}

	a.page, _ = strconv.Atoi(rawPage)
	if a.page < 1 {
		a.page = 1
	}

	return a.page
}

func (a paginatedAction) paginate(db *gorm.DB, page, pageSize int) *gorm.DB {
	return flentities.Repository{DB: db}.Paginate(page, pageSize)
}

func (a paginatedAction) addPageDiagnostic() {
	a.Diagnostics().Add("page", strconv.Itoa(a.getPage()))
}
