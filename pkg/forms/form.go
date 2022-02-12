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

func (f *Form) PermittedValues(field string, options ...string) {
	value := f.Get(field)
	if value == "" {
		return
	}

	for _, option := range options {
		if value == option {
			return
		}
	}

	f.Errors.Add(field, "This field is invalid.")
}

func (f *Form) MinLength(field string, minLength int) {
	value := f.Get(field)
	if value == "" {
		return
	}

	strLength := utf8.RuneCountInString(value)
	if strLength < minLength {
		f.Errors.Add(field, fmt.Sprintf("This field must contain more than %d characters", minLength))
	}
}

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

// Valid will return false when f.Errors has non-zero length, true otherwise
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}