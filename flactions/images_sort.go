package flactions

import (
	"net/http"

	"github.com/nerde/fuji-lane-back/flentities"

	"github.com/jinzhu/gorm"
)

func imagesDefaultOrder(db *gorm.DB) *gorm.DB {
	return db.Order("images.position, images.id")
}

// ImagesSort marks an image as uploaded
type ImagesSort struct{}

// Perform the action
func (a *ImagesSort) Perform(c Context) {
	ids := []uint{}
	if err := c.BindJSON(&ids); err != nil {
		c.Diagnostics().AddError(err)
		c.Respond(http.StatusBadRequest, nil)
		return
	}

	c.Diagnostics().AddJSON("ids", ids)

	account := c.CurrentAccount()

	c.Repository().Transaction(func(tx *flentities.Repository) {
		for idx, id := range ids {
			image := &flentities.Image{ID: id}
			if err := tx.Preload("Property").Preload("Unit.Property").Find(image).Error; err != nil {
				tx.Rollback()
				c.ServerError(err)
				return
			}

			if !image.BelongsTo(account.ID) {
				continue
			}

			tx.Model(image).Updates(map[string]interface{}{"Position": idx + 1})
		}

		if err := tx.Commit().Error; err != nil {
			c.ServerError(err)
			return
		}
	})
}
