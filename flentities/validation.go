package flentities

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

const (
	emailRegexString = "^(?:(?:(?:(?:[a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(?:\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|(?:(?:\\x22)(?:(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))?(?:\\x20|\\x09)+)?(?:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:\\(?:[\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))?(\\x20|\\x09)+)?(?:\\x22)))@(?:(?:(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])(?:[a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(?:(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])(?:[a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$"
)

var allowedImageTypes = map[string]interface{}{"image/gif": nil, "image/jpeg": nil, "image/png": nil}

// Validatable can be validated and supports validation errors
type Validatable interface {
	Validate() []error
}

// ValidateFields concatenates the errors resulting from validating several fields
func ValidateFields(field FieldValidation, fields ...FieldValidation) []error {
	errs := field.Errors

	for _, f := range fields {
		errs = append(errs, f.Errors...)
	}

	return errs
}

// ValidateField returns a FieldValidation for the given parameters
func ValidateField(name string, value interface{}) FieldValidation {
	return FieldValidation{name, value, []error{}}
}

// FieldValidation carries the context and errors when validating a specific field
type FieldValidation struct {
	Name   string
	Value  interface{}
	Errors []error
}

// Required adds an error if the value is blank
func (v FieldValidation) Required() FieldValidation {
	valid := true
	switch val := v.Value.(type) {
	case string:
		valid = len(val) != 0
		break
	case uint:
		valid = val != 0
		break
	case int:
		valid = val != 0
		break
	default:
		v.Errors = append(v.Errors, v.unsupportedTypeError("Required"))
	}

	if !valid {
		v.Errors = append(v.Errors, fmt.Errorf("%s is required", v.Name))
	}

	return v
}

// Email adds an error if the value is not a valid email
func (v FieldValidation) Email() FieldValidation {
	switch val := v.Value.(type) {
	case string:
		emailReg := regexp.MustCompile(emailRegexString)
		if !emailReg.Match([]byte(val)) {
			v.Errors = append(v.Errors, fmt.Errorf("Invalid email: %s", val))
		}
		break
	default:
		v.Errors = append(v.Errors, v.unsupportedTypeError("Email"))
	}

	return v
}

// MinLen adds an error if the value does not comply with the minimum size
func (v FieldValidation) MinLen(min int) FieldValidation {
	switch val := v.Value.(type) {
	case string:
		if len(val) < min {
			v.Errors = append(v.Errors, fmt.Errorf("Invalid %s: minimum size is %d", v.Name, min))
		}
		break
	default:
		v.Errors = append(v.Errors, v.unsupportedTypeError("MinLen"))
	}

	return v
}

// MaxLen adds an error if the value does not comply with the maximum size
func (v FieldValidation) MaxLen(max int) FieldValidation {
	switch val := v.Value.(type) {
	case string:
		if len(val) > max {
			v.Errors = append(v.Errors, fmt.Errorf("Invalid %s: maximum size is %d", v.Name, max))
		}
		break
	default:
		v.Errors = append(v.Errors, v.unsupportedTypeError("MaxLen"))
	}

	return v
}

// Min adds an error if the value is less than the given minimum
func (v FieldValidation) Min(min int) FieldValidation {
	switch val := v.Value.(type) {
	case int:
		if val < min {
			v.Errors = append(v.Errors, fmt.Errorf("Invalid %s: minimum is %d", v.Name, min))
		}
		break
	default:
		v.Errors = append(v.Errors, v.unsupportedTypeError("Min"))
	}

	return v
}

// Max adds an error if the value is more than the given maximum
func (v FieldValidation) Max(max int) FieldValidation {
	switch val := v.Value.(type) {
	case int:
		if val > max {
			v.Errors = append(v.Errors, fmt.Errorf("Invalid %s: maximum is %d", v.Name, max))
		}
		break
	default:
		v.Errors = append(v.Errors, v.unsupportedTypeError("Max"))
	}

	return v
}

// Image adds an error if the value is not an image type
func (v FieldValidation) Image() FieldValidation {
	switch val := v.Value.(type) {
	case string:
		if _, ok := allowedImageTypes[val]; !ok {
			v.Errors = append(v.Errors, fmt.Errorf("Invalid %s: needs to be JPEG, PNG or GIF", v.Name))
		}
		break
	default:
		v.Errors = append(v.Errors, v.unsupportedTypeError("Image"))
	}

	return v
}

func (v FieldValidation) unsupportedTypeError(validator string) error {
	return fmt.Errorf("Unsupported type for %s validator: %s", validator, reflect.TypeOf(v.Value))
}

var invalidHTMLTags = []string{"script"}

// HTML adds an error if the value is not an acceptable HTML
func (v FieldValidation) HTML() FieldValidation {
	validator := func(val string) {
		node, err := html.Parse(strings.NewReader(val))
		if err != nil {
			v.Errors = append(v.Errors, fmt.Errorf("%s: %s", v.Name, err.Error()))
			return
		}

		var searchForInvalidTags func(*html.Node)
		searchForInvalidTags = func(n *html.Node) {
			if n.Type == html.ElementNode {
				for _, tagName := range invalidHTMLTags {
					if n.Data == tagName {
						v.Errors = append(v.Errors, fmt.Errorf("%s: %s tags are not allowed", v.Name, tagName))
					}
				}
			}

			for c := n.FirstChild; c != nil; c = c.NextSibling {
				searchForInvalidTags(c)
			}
		}

		searchForInvalidTags(node)
	}

	switch val := v.Value.(type) {
	case *string:
		if val != nil {
			validator(*val)
		}
	case string:
		if val != "" {
			validator(val)
		}
		break
	default:
		v.Errors = append(v.Errors, v.unsupportedTypeError("HTML"))
	}

	return v
}
