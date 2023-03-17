package utils

import (
	"fmt"
	"reflect"
	"strings"
)

// ComparableContains checks if a given value of type T is contained within the provided slice of values of type T.
// It returns true if the value is found and false otherwise.
// The type parameter T is constrained to the comparable type, which means that values of type T must be comparable using the == operator.
// This function is optimized for performance and type safety by avoiding the use of reflection or any other additional checks.
func ComparableContains[T comparable](finder T, sources ...T) bool {
	// Loop through the slice of values and check if the current element is equal to the given value using the == operator.
	for _, e := range sources {
		if e == finder {
			return true
		}
	}

	// If the value was not found, return false.
	return false
}

// AnyContains checks if the given slice contains the given element.
// It takes a variable number of arguments of any type.
// The first argument is the element to be searched for, and the rest are the elements in the slice.
// It returns true if the element is found in the slice, false otherwise.
func AnyContains[T any](finder T, sources ...T) bool {
	for _, e := range sources {
		if reflect.DeepEqual(e, finder) {
			return true
		}
	}

	return false
}

// StringFormat formats the given string by replacing placeholders with the
// corresponding values of the parameters. Placeholders are in the format of {i},
// where i is the index of the parameter in the list.
func StringFormat(source string, params ...interface{}) string {
	for i, param := range params {
		source = strings.ReplaceAll(source, fmt.Sprintf("{%d}", i), fmt.Sprint(param))
	}
	return source
}

// Select applies the given predicate function to each element in the input slice,
// and returns a new slice containing the results of the predicate function applied to each element.
func Select[T any, V any](sources []T, predicate func(T) V) []V {
	// Create a new slice to hold the result of the predicate function.
	result := make([]V, len(sources))
	// Iterate over each element in the input slice.
	for i, e := range sources {
		// Apply the predicate function to the current element, and store the result in the result slice.
		result[i] = predicate(e)
	}
	// Return the result slice.
	return result
}

// Select applies the given predicate function to each element in the input slice,
// and returns a new slice containing the results of the predicate function applied to each element.
func SelectIndex[T any, V any](sources []T, predicate func(int, T) V) []V {
	// Create a new slice to hold the result of the predicate function.
	result := make([]V, len(sources))
	// Iterate over each element in the input slice.
	for i, e := range sources {
		// Apply the predicate function to the current element, and store the result in the result slice.
		result[i] = predicate(i, e)
	}
	// Return the result slice.
	return result
}

// Where filters the input slice by applying the given conditional function to each element,
// and returns a new slice containing only the elements for which the conditional function returns true.
func Where[T any](sources []T, conditional func(T) bool) []T {
	// Create a new slice to hold the filtered elements.
	result := make([]T, 0)
	// Iterate over each element in the input slice.
	for _, e := range sources {
		// If the conditional function returns true for the current element,
		// append the element to the result slice.
		if conditional(e) {
			result = append(result, e)
		}
	}
	// Return the result slice.
	return result
}

// Find searches the input slice for the first element that matches the given predicate,
// and returns that element. If no element matches the predicate, it returns the zero value of type T.
func Find[T any](sources []T, predicate func(T) bool) T {
	// Iterate over each element in the input slice.
	for _, e := range sources {
		// If the current element matches the predicate, return it.
		if predicate(e) {
			return e
		}
	}
	// If no element matches the predicate, create a new zero value of type T and return it.
	// This ensures that the function always returns a value of type T.
	var zeroT T
	return zeroT
}

// FindLast searches the input slice for the last element that matches the given predicate,
// and returns that element. If no element matches the predicate, it returns the zero value of type T.
func FindLast[T any](sources []T, predicate func(T) bool) T {
	// Iterate over the input slice in reverse order, starting from the last element.
	for i := len(sources) - 1; i >= 0; i-- {
		// Get the current element from the input slice.
		e := sources[i]
		// If the current element matches the predicate, return it.
		if predicate(e) {
			return e
		}
	}
	// If no element matches the predicate, create a new zero value of type T and return it.
	// This ensures that the function always returns a value of type T.
	var zeroT T
	return zeroT
}

// ToInterfaceSlice converts a slice of any type to a slice of interface{}.
// It takes a reflect.Value as input and returns a []interface{}.
// If the input is nil, it returns an empty []interface{}.
// It iterates through the slice and adds each element to the result slice as an interface{}.
// If an element is a pointer, it dereferences the pointer and adds the underlying value as an interface{}.
func ToInterfaceSlice(s reflect.Value) []interface{} {
	result := make([]interface{}, s.Len())
	if s.IsNil() {
		return result
	}

	for i := 0; i < s.Len(); i++ {
		result[i] = handleValuePointer(s.Index(i)).Interface()
	}

	return result
}

func RemoveAt[T any](sources []T, index int) []T {
	var def T
	copy(sources[index:], sources[index+1:])
	sources[len(sources)-1] = def
	return sources[:len(sources)-1]
}
