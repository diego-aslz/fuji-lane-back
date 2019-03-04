package optional

import (
	"encoding/json"
	"reflect"
)

// String represents a string that may or may not be set
type String struct {
	Value *string
	Set   bool
}

// UnmarshalJSON marks as Set and unmarshals the Value
func (o *String) UnmarshalJSON(data []byte) error {
	o.Set = true
	return json.Unmarshal(data, &o.Value)
}

// MarshalJSON returns the JSON encoding of o.Value
func (o String) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.Value)
}

// Update changes dst only if Set is true
func (o String) Update(dst interface{}, zeroMeansNil bool) {
	if !o.Set {
		return
	}

	var v interface{}

	switch dst.(type) {
	case *string:
		if o.Value == nil {
			v = ""
		} else {
			v = *o.Value
		}
		break
	case **string:
		if zeroMeansNil && o.Value != nil && *o.Value == "" {
			v = nil
		} else {
			v = o.Value
		}

		break
	}

	reflect.ValueOf(dst).Elem().Set(reflect.ValueOf(v))
}
