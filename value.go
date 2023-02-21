package utils

// IIF returns `trueVal` if `v` is true, otherwise returns `falseVal`.
// It is a generic function and works with any data type.
func IIF[T any](v bool, trueVal T, falseVal T) T {
	if v {
		return trueVal
	}

	return falseVal
}

// BoolToInt returns 1 if `v` is true, otherwise returns 0.
func BoolToInt(v bool) int {
	if v {
		return 1
	}
	return 0
}
