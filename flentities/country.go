package flentities

import "github.com/jinzhu/gorm"

// Country we support
type Country struct {
	gorm.Model
	Name string
}
