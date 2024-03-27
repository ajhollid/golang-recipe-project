package forms

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"net/url"
	"strings"
)

// Form creates a custom form struct and embeds a url.Values object
type Form struct {
	url.Values
	Errors errors
}

// New initializes a form struct
func New(data url.Values) *Form {
	formErrors := errors(map[string][]string{})
	newForm := Form{
		data,
		formErrors,
	}
	return &newForm
}

// Validators

// Required checks for required fields
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		input := f.Get(field)
		if strings.TrimSpace(input) == "" {
			f.Errors.Add(field, "This field cannot be empty")
		}
	}
}

// Has checks if a form field exists and is not empty
func (f *Form) Has(field string) bool {
	input := f.Get("field")
	if input == "" {
		return false
	}
	return true
}

// MinLength checks for min length of a field
func (f *Form) MinLength(field string, length int) bool {
	input := f.Get(field)
	if len(input) < length {
		f.Errors.Add(field, fmt.Sprintf("This field must be at least %d characters long", length))
		return false
	}
	return true
}

// IsEmail Checks for valid Email address
func (f *Form) IsEmail(field string) {
	if !govalidator.IsEmail(f.Get(field)) {
		f.Errors.Add(field, "Invalid email address")
	}
}

// Valid returns true if there are no errors, otherwise false
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}
