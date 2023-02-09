package utils

import (
	"fmt"
	"reflect"
	"strings"
)

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

func StringFormat[T any](source string, params ...T) string {
	for i, e := range params {
		source = strings.ReplaceAll(source, fmt.Sprintf("{%d}", i), fmt.Sprintf("%v", e))
	}

	return source
}
