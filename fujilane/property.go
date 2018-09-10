package fujilane

import "github.com/jinzhu/gorm"

// PropertyStatePending means this property is being filled in by the user
const PropertyStatePending = 1

var propertyStates map[int]string

// Property contains address and can have multiple units that can be booked
type Property struct {
	gorm.Model
	Name    string
	StateID int `gorm:"column:state"`
	UserID  int
	User    *User
}

// State returns the state name for the property's state ID
func (p *Property) State() string {
	return propertyStates[p.StateID]
}

func init() {
	propertyStates = make(map[int]string)
	propertyStates[PropertyStatePending] = "Pending"
}
