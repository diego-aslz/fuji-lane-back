package fujilane

import (
	"time"

	"github.com/jinzhu/gorm"
)

// User of the system
type User struct {
	gorm.Model
	Name         string
	Email        string
	FacebookID   string
	LastSignedIn time.Time
}
