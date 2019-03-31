package flactions

import (
	"net/http"

	"github.com/nerde/fuji-lane-back/flentities"
)

// ImagesSort marks an image as uploaded
type ImagesSort struct {
	Context
}

// Perform the action
func (a *ImagesSort) Perform() {
	ids := []uint{}
	if err := a.BindJSON(&ids); err != nil {
		a.Diagnostics().AddError(err)
		a.Respond(http.StatusBadRequest, nil)
		return
	}

	a.Diagnostics().AddJSON("ids", ids)

	account := a.CurrentAccount()

	a.Repository().Transaction(func(tx *flentities.Repository) {
		for idx, id := range ids {
			image := &flentities.Image{ID: id}
			if err := tx.Preload("Property").Preload("Unit.Property").Find(image).Error; err != nil {
				tx.Rollback()
				a.ServerError(err)
				return
			}

			if !image.BelongsTo(account.ID) {
				continue
			}

			if err := tx.UpdatesColVal(image, "Position", idx+1); err != nil {
				tx.Rollback()
				a.ServerError(err)
				return
			}
		}

		if err := tx.Commit().Error; err != nil {
			a.ServerError(err)
			return
		}
	})
}

// NewImagesSort returns a new ImagesSort action
func NewImagesSort(c Context) Action {
	return &ImagesSort{c}
}
