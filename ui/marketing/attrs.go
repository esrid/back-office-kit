package marketing

import "strconv"

// intAttr renders a dimension, or nothing when it is unset. An explicit
// width="0" would be worse than no attribute at all.
func intAttr(v int) string {
	if v <= 0 {
		return ""
	}
	return strconv.Itoa(v)
}
