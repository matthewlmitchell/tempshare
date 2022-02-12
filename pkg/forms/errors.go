package forms

type errors map[string][]string

// Add appends an error message to the specified field in our errors map
func (e errors) Add(field, message string) {
	e[field] = append(e[field], message)
}

// Get retrieves the first error message of a given field from our errors map
func (e errors) Get(field string) string {
	es := e[field]
	if len(es) == 0 {
		return ""
	}

	return es[0]
}