package fujilane

import "github.com/jinzhu/gorm"

type propertiesRepository struct{}

func (r *propertiesRepository) create(user *User) (p *Property, err error) {
	err = withDatabase(func(db *gorm.DB) error {
		p = &Property{AccountID: *user.AccountID, StateID: PropertyStateDraft}

		return db.Create(p).Error
	})
	return
}
