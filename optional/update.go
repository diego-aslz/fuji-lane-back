package optional

import (
	"reflect"
)

// Field represents an optional value in a payload
type Field interface {
	Update(interface{}, bool)
}

// Update copies all optional field values that got set in src to the dst struct matching fields by name
func Update(src, dst interface{}) {
	v := reflect.ValueOf(src)
	d := reflect.ValueOf(dst).Elem()

	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)

		df := d.FieldByName(field.Name)
		if !df.IsValid() {
			continue
		}

		if of, ok := v.Field(i).Interface().(Field); ok {
			of.Update(df.Addr().Interface(), true)
		}
	}
}
