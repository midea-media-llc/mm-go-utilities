package utils

import (
	"encoding/json"
)

func JsonToObject[T any](jsonString string, value T) error {
	if jsonString == "" {
		return nil
	}

	bytes := []byte(jsonString)

	if err := json.Unmarshal(bytes, value); err != nil {
		return err
	}

	return nil
}

func ObjectToJson[T any](object T, value *string) error {
	result, err := json.Marshal(object)
	if err != nil {
		return err
	}

	*value = string(result)
	return nil
}
