package optional

import (
	"encoding/json"
	"reflect"
)

// Bool represents a bool that may or may not be set
type Bool struct {
	Value *bool
	Set   bool
}

// UnmarshalJSON marks as Set and unmarshals the Value
func (o *Bool) UnmarshalJSON(data []byte) error {
	o.Set = true
	return json.Unmarshal(data, &o.Value)
}

// MarshalJSON returns the JSON encoding of o.Value
func (o Bool) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.Value)
}

// Update changes dst only if Set is true
func (o Bool) Update(dst interface{}) {
	if !o.Set {
		return
	}

	var v interface{}

	switch dst.(type) {
	case *bool:
		if o.Value == nil {
			v = false
		} else {
			v = *o.Value
		}
		break
	case **bool:
		v = o.Value

		break
	}

	reflect.ValueOf(dst).Elem().Set(reflect.ValueOf(v))
}
