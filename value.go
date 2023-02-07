package utils

func IIF[T any](v bool, trueVal T, falseVal T) T {
	if v {
		return trueVal
	}

	return falseVal
}

func BoolToInt(v bool) int {
	if v {
		return 1
	}
	return 0
}
