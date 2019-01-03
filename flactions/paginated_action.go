package flactions

import (
	"strconv"

	"github.com/jinzhu/gorm"
)

type paginatedAction struct{}

func (paginatedAction) page(c Context) int {
	page := c.Query("page")
	if page == "" {
		return 1
	}

	p, _ := strconv.Atoi(page)
	if p < 1 {
		return 1
	}

	return p
}

func (a paginatedAction) paginate(db *gorm.DB, page, pageSize int) *gorm.DB {
	return db.Limit(pageSize).Offset((page - 1) * bookingsPageSize)
}
