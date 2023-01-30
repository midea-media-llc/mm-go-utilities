package utils

import (
	"reflect"
)

func CloneFields[IN comparable, OUT comparable](input IN, output OUT, ignoreFields ...string) OUT {
	srcType := reflect.TypeOf(input).Elem()
	destType := reflect.TypeOf(output).Elem()
	srcValue := reflect.ValueOf(input).Elem()
	destValue := reflect.ValueOf(output).Elem()

	for i := 0; i < destType.NumField(); i++ {
		name := destType.Field(i).Name
		if ArrayStringContains(ignoreFields, name) {
			continue
		}

		if _, ok := srcType.FieldByName(name); ok {
			destValue.FieldByName(name).Set(srcValue.FieldByName(name))
		}
	}

	return output
}
