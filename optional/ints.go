package optional

import (
	"encoding/json"
	"reflect"
)

// Int represents an int that may or may not be set
type Int struct {
	Value *int
	Set   bool
}

// UnmarshalJSON marks as Set and unmarshals the Value
func (o *Int) UnmarshalJSON(data []byte) error {
	o.Set = true
	return json.Unmarshal(data, &o.Value)
}

// MarshalJSON returns the JSON encoding of o.Value
func (o Int) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.Value)
}

// Update changes dst only if Set is true
func (o Int) Update(dst interface{}, zeroMeansNil bool) {
	if !o.Set {
		return
	}

	var v interface{}

	switch dst.(type) {
	case *int:
		if o.Value == nil {
			v = 0
		} else {
			v = *o.Value
		}
		break
	case **int:
		if zeroMeansNil && o.Value != nil && *o.Value == 0 {
			var it *int
			v = it
		} else {
			v = o.Value
		}

		break
	}

	reflect.ValueOf(dst).Elem().Set(reflect.ValueOf(v))
}

// Uint represents an uint that may or may not be set
type Uint struct {
	Value *uint
	Set   bool
}

// UnmarshalJSON marks as Set and unmarshals the Value
func (o *Uint) UnmarshalJSON(data []byte) error {
	o.Set = true
	return json.Unmarshal(data, &o.Value)
}

// MarshalJSON returns the JSON encoding of o.Value
func (o Uint) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.Value)
}

// Update changes dst only if Set is true
func (o Uint) Update(dst interface{}, zeroMeansNil bool) {
	if !o.Set {
		return
	}

	var v interface{}

	switch dst.(type) {
	case *uint:
		if o.Value == nil {
			v = uint(0)
		} else {
			v = *o.Value
		}
		break
	case **uint:
		if zeroMeansNil && o.Value != nil && *o.Value == 0 {
			var it *uint
			v = it
		} else {
			v = o.Value
		}

		break
	}

	reflect.ValueOf(dst).Elem().Set(reflect.ValueOf(v))
}
