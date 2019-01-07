package flentities

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
	"time"
)

// Layout for parsing and formatting Date objects
const Layout = "2006-01-02"

// Date with no time
type Date struct {
	time.Time
}

func (d Date) String() string {
	return d.Time.Format(Layout)
}

// MarshalJSON marshals this date to a String
func (d Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

// UnmarshalJSON parses a String from a JSON into this date
func (d *Date) UnmarshalJSON(data []byte) (err error) {
	raw := ""
	err = json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}

	*d, err = ParseDate(raw)
	return
}

// ParseDate converts a raw string input into a *Date
func ParseDate(raw string) (Date, error) {
	t, err := time.Parse(Layout, raw)
	return Date{t}, err
}

// Value returns a String representation for d
func (d Date) Value() (driver.Value, error) {
	return d.String(), nil
}

// Scan loads a scanned value into d
func (d *Date) Scan(value interface{}) (err error) {
	switch v := value.(type) {
	case string:
		*d, err = ParseDate(v)
	case time.Time:
		*d = Date{v}
	default:
		err = fmt.Errorf("Cannot Scan Date from a %s", reflect.TypeOf(value))
	}
	return
}
