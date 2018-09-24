package fujilane

import (
	"fmt"
	"regexp"
)

const (
	emailRegexString = "^(?:(?:(?:(?:[a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(?:\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|(?:(?:\\x22)(?:(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))?(?:\\x20|\\x09)+)?(?:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:\\(?:[\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))?(\\x20|\\x09)+)?(?:\\x22)))@(?:(?:(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])(?:[a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(?:(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])(?:[a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$"
)

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
func ValidateField(name, value string) FieldValidation {
	return FieldValidation{name, value, []error{}}
}

// FieldValidation carries the context and errors when validating a specific field
type FieldValidation struct {
	Name   string
	Value  string
	Errors []error
}

// Required adds an error if the value is blank
func (v FieldValidation) Required() FieldValidation {
	if len(v.Value) == 0 {
		v.Errors = append(v.Errors, fmt.Errorf("Invalid %s: cannot be blank", v.Name))
	}

	return v
}

// Email adds an error if the value is not a valid email
func (v FieldValidation) Email() FieldValidation {
	emailReg := regexp.MustCompile(emailRegexString)
	if !emailReg.Match([]byte(v.Value)) {
		v.Errors = append(v.Errors, fmt.Errorf("Invalid email: %s", v.Value))
	}

	return v
}

// MinLen adds an error if the value does not comply with the minimum size
func (v FieldValidation) MinLen(min int) FieldValidation {
	if len(v.Value) < min {
		v.Errors = append(v.Errors, fmt.Errorf("Invalid %s: minimum size is %d", v.Name, min))
	}

	return v
}

// MaxLen adds an error if the value does not comply with the maximum size
func (v FieldValidation) MaxLen(max int) FieldValidation {
	if len(v.Value) > max {
		v.Errors = append(v.Errors, fmt.Errorf("Invalid %s: maximum size is %d", v.Name, max))
	}

	return v
}
