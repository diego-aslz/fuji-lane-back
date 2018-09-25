package fujilane

import "github.com/jinzhu/gorm"

type accountsRepository struct{}

func (r *accountsRepository) create(account *Account) error {
	return withDatabase(func(db *gorm.DB) error {
		return db.Create(account).Error
	})
}
