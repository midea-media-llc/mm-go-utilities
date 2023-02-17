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

func Select[T any, V any](sources []T, predicate func(T) V) []V {
	result := make([]V, len(sources))
	for i, e := range sources {
		result[i] = predicate(e)
	}

	return result
}

func Where[T any](sources []T, conditional func(T) bool) []T {
	result := make([]T, 0)
	for _, e := range sources {
		if conditional(e) {
			result = append(result, e)
		}
	}

	return result
}

func Find[T any](sources []T, predicate func(T) bool) T {
	for _, e := range sources {
		if predicate(e) {
			return e
		}
	}

	isPointer := false
	typeT := reflect.TypeOf(sources).Elem()
	for ComparableContains(typeT.Kind(), reflect.Pointer) {
		typeT = typeT.Elem()
		isPointer = true
	}

	v := reflect.New(typeT).Interface()
	if isPointer {
		v = &v
	}
	return v.(T)
}

func FindLast[T any](sources []T, predicate func(T) bool) T {
	for i := len(sources) - 1; i >= 0; i++ {
		e := sources[i]
		if predicate(e) {
			return e
		}

	}

	isPointer := false
	typeT := reflect.TypeOf(sources).Elem()
	for ComparableContains(typeT.Kind(), reflect.Pointer) {
		typeT = typeT.Elem()
		isPointer = true
	}

	v := reflect.New(typeT).Interface()
	if isPointer {
		v = &v
	}
	return v.(T)
}
