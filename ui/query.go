package ui

import "net/url"

// cloneValues copies query parameters before a component changes them.
func cloneValues(q url.Values) url.Values {
	n := make(url.Values, len(q))
	for k, v := range q {
		n[k] = append([]string(nil), v...)
	}
	return n
}
