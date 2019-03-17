package fldiagnostics

import (
	"encoding/json"
)

// Diagnostics is used to store information relevant for logging and debugging
type Diagnostics struct {
	keys   []string
	values []string
}

// Add new information to this Diagnostics
func (d *Diagnostics) Add(key, value string) *Diagnostics {
	d.keys = append(d.keys, key)
	d.values = append(d.values, value)
	return d
}

// AddQuoted adds a quoted string
func (d *Diagnostics) AddQuoted(key, value string) *Diagnostics {
	return d.Add(key, "\""+value+"\"")
}

// AddError adds an error message
func (d *Diagnostics) AddError(err error) *Diagnostics {
	return d.AddErrorAs("error", err)
}

// AddErrorAs adds an error with a custom key
func (d *Diagnostics) AddErrorAs(key string, err error) *Diagnostics {
	return d.AddQuoted(key, err.Error())
}

// AddJSON converts the value to JSON and adds it
func (d *Diagnostics) AddJSON(key string, value interface{}) *Diagnostics {
	if sensitive, ok := value.(SensitivePayload); ok {
		value = sensitive.FilterSensitiveInformation()
	}

	jsonObj, err := json.Marshal(value)
	if err == nil {
		d.Add(key, string(jsonObj))
	} else {
		d.AddQuoted("json_parsing_error", err.Error())
	}

	return d
}

// Concat adds all keys and values from other Diagnostics into this one
func (d *Diagnostics) Concat(other *Diagnostics) *Diagnostics {
	for i, key := range other.keys {
		d.Add(key, other.values[i])
	}
	return d
}

func (d Diagnostics) String() string {
	result := ""

	for i, key := range d.keys {
		if i > 0 {
			result += " "
		}

		result += key + "=" + d.values[i]
	}

	return result
}

// ToMap returns a map representation of this diagnostics
func (d Diagnostics) ToMap() map[string]string {
	res := map[string]string{}

	for i, key := range d.keys {
		res[key] = d.values[i]
	}

	return res
}
