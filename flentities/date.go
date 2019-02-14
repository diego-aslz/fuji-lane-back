package flentities

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"math"
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

// NightsFrom returns the number of nights between two dates
func (d Date) NightsFrom(other Date) int {
	return int(math.Ceil(d.Sub(other.Time).Hours() / 24))
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

// Human returns the date in a human friendly format
func (d Date) Human() string {
	return d.Time.Format("Mon, 02 Jan 2006")
}
