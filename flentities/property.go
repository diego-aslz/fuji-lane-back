package flentities

import (
	"encoding/json"

	"github.com/jinzhu/gorm"
)

// PropertyStateDraft means this property is being filled in by the user
const PropertyStateDraft = 1

var propertyStates map[int]string

// Property contains address and can have multiple units that can be booked
type Property struct {
	gorm.Model
	Name      *string  `json:"name"`
	StateID   int      `gorm:"column:state" json:"-"`
	AccountID int      `json:"-"`
	Account   *Account `json:"-"`
}

// State returns the state name for the property's state ID
func (p *Property) State() string {
	return propertyStates[p.StateID]
}

func init() {
	propertyStates = make(map[int]string)
	propertyStates[PropertyStateDraft] = "draft"
}

type propertyAlias Property

type propertyUI struct {
	ID    uint   `json:"id"`
	State string `json:"state"`
	*propertyAlias
}

func (p *propertyUI) stateID() int {
	for id, state := range propertyStates {
		if state == p.State {
			return id
		}
	}
	return 0
}

// MarshalJSON returns JSON bytes for a Property
func (p *Property) MarshalJSON() ([]byte, error) {
	return json.Marshal(propertyUI{
		ID:            p.ID,
		State:         p.State(),
		propertyAlias: (*propertyAlias)(p),
	})
}

// UnmarshalJSON loads a Property from JSON bytes
func (p *Property) UnmarshalJSON(data []byte) error {
	aux := &propertyUI{propertyAlias: (*propertyAlias)(p)}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	p.ID = aux.ID
	p.StateID = aux.stateID()
	return nil
}
