package utils

import "reflect"

func ComparableContains[T comparable](finder T, sources ...T) bool {
	for _, e := range sources {
		if e == finder {
			return true
		}
	}

	return false
}

func AnyContains[T any](finder T, sources ...T) bool {
	for _, e := range sources {
		if reflect.ValueOf(e) == reflect.ValueOf(finder) {
			return true
		}
	}

	return false
}
