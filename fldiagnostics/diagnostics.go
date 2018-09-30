package fldiagnostics

import "encoding/json"

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
	return d.AddQuoted("error", err.Error())
}

// AddJSON converts the value to JSON and adds it
func (d *Diagnostics) AddJSON(key string, value interface{}) *Diagnostics {
	jsonObj, err := json.Marshal(value)
	if err == nil {
		d.Add(key, string(jsonObj))
	}
	return d
}

// AddSensitive filters out sensitive information before adding as JSON
func (d *Diagnostics) AddSensitive(key string, value SensitivePayload) *Diagnostics {
	return d.AddJSON(key, value.FilterSensitiveInformation())
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
