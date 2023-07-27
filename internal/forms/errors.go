package forms

// package for error validation server-side

type errors map[string][]string

// Add adds an error message for a given form field
func (e errors) Add(field, message string) {
	e[field] = append(e[field], message) // add message to the slice of strings having under field key in e map
}

// Get returns the first error message
func (e errors) Get(field string) string {
	es := e[field]

	if len(es) == 0 {
		return ""
	}

	return es[0]
}
