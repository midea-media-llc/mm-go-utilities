package utils

import (
	"encoding/json"
)

// JsonToObject is a generic function that parses a JSON string and unmarshals it
// into a value of type T. The JSON string is passed in as the first argument, and
// the value to unmarshal into is passed in as the second argument. The function
// returns an error if the JSON string cannot be unmarshaled into the value.
func JsonToObject[T any](jsonString string, value T) error {
	if jsonString == "" {
		return nil
	}

	// Convert the JSON string to a byte slice.
	bytes := []byte(jsonString)

	// Unmarshal the JSON string into the value.
	if err := json.Unmarshal(bytes, value); err != nil {
		return err
	}

	return nil
}

// ObjectToJson is a generic function that marshals a value of type T into a
// JSON string. The value to marshal is passed in as the first argument, and a
// pointer to a string variable to store the result is passed in as the second
// argument. The function returns an error if the value cannot be marshaled into
// a JSON string.
func ObjectToJson[T any](object T, value *string) error {
	// Marshal the value into a JSON byte slice.
	result, err := json.Marshal(object)
	if err != nil {
		return err
	}

	// Convert the JSON byte slice to a string and store it in the value pointer.
	*value = string(result)
	return nil
}
