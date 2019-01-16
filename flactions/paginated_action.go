package flactions

import (
	"strconv"

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
	return db.Limit(pageSize).Offset((page - 1) * pageSize)
}

func (a paginatedAction) addPageDiagnostic() {
	a.Diagnostics().Add("page", strconv.Itoa(a.getPage()))
}
