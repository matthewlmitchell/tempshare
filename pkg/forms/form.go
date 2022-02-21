package forms

import (
	"fmt"
	"net/url"
	"strings"
	"unicode/utf8"
)

type Form struct {
	url.Values
	Errors errors
}

// New initializes a new Form struct given a set of values (url.Values)
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// Required ensures that necessary form fields are non-empty
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		// If the field is an empty string after trimming white-space, append an error
		// to the field
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field must not be blank")
		}
	}
}

// PermittedValues asserts that the given form field's inputs
// are valid and acceptable. If a non-permitted value is given for a field,
// an error is appended to the field: "This field is invalid"
func (f *Form) PermittedValues(field string, options ...string) {
	// Extract the value of a given field from the form
	value := f.Get(field)
	if value == "" {
		return
	}

	// Iterate through the possible valid inputs, and compare them to what
	// the client provided
	for _, option := range options {
		if value == option {
			return
		}
	}

	f.Errors.Add(field, "This field is invalid.")
}

// MinLength asserts that the value in a given field contains at least minLength
// number of runes in the string
func (f *Form) MinLength(field string, minLength int) {
	value := f.Get(field)
	if value == "" {
		return
	}

	// Count the number of runes in the string to determine if it is of appropriate length
	// Note: if we simply did len(value) this check would return an incorrect value
	//		 for characters that take up more than one byte, e.g. high-ansi
	strLength := utf8.RuneCountInString(value)
	if strLength < minLength {
		f.Errors.Add(field, fmt.Sprintf("This field must contain more than %d characters", minLength))
	}
}

// MaxLength asserts that the value in a given field contains at most maxLength
// number of runes in the string
func (f *Form) MaxLength(field string, maxLength int) {
	value := f.Get(field)
	if value == "" {
		return
	}

	strLength := utf8.RuneCountInString(value)
	if strLength > maxLength {
		f.Errors.Add(field, fmt.Sprintf("This field must contain less than %d characters", maxLength))
	}
}

// Valid is used for determining if a form's inputs are invalid or not.
// If there are any errors, the length of f.Errors will be non-zero, and return false.
// If f.Errors is of length zero, return true (the form inputs are valid).
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}
