package fujilane

import (
	"github.com/jinzhu/gorm"
	"github.com/nerde/fuji-lane-back/flentities"
)

type propertiesRepository struct{}

func (r *propertiesRepository) create(user *flentities.User) (p *flentities.Property, err error) {
	err = withDatabase(func(db *gorm.DB) error {
		p = &flentities.Property{AccountID: *user.AccountID, StateID: flentities.PropertyStateDraft}

		return db.Create(p).Error
	})
	return
}
