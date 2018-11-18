package flactions

import "github.com/jinzhu/gorm"

func imagesDefaultOrder(db *gorm.DB) *gorm.DB {
	return db.Order("images.position, images.id")
}
