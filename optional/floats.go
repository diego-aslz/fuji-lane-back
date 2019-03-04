package optional

import (
	"encoding/json"
	"reflect"
)

// Float32 represents a float32 that may or may not be set
type Float32 struct {
	Value *float32
	Set   bool
}

// UnmarshalJSON marks as Set and unmarshals the Value
func (o *Float32) UnmarshalJSON(data []byte) error {
	o.Set = true
	return json.Unmarshal(data, &o.Value)
}

// MarshalJSON returns the JSON encoding of o.Value
func (o Float32) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.Value)
}

// Update changes dst only if Set is true
func (o Float32) Update(dst interface{}, zeroMeansNil bool) {
	if !o.Set {
		return
	}

	var v interface{}

	switch dst.(type) {
	case *float32:
		if o.Value == nil {
			v = float32(0)
		} else {
			v = *o.Value
		}
		break
	case **float32:
		if zeroMeansNil && o.Value != nil && *o.Value == 0 {
			v = nil
		} else {
			v = o.Value
		}

		break
	}

	reflect.ValueOf(dst).Elem().Set(reflect.ValueOf(v))
}
