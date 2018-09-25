package fujilane

import (
	"github.com/jinzhu/gorm"
	"github.com/nerde/fuji-lane-back/flentities"
)

type accountsRepository struct{}

func (r *accountsRepository) create(account *flentities.Account) error {
	return withDatabase(func(db *gorm.DB) error {
		return db.Create(account).Error
	})
}
