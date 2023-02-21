package utils

import (
	"reflect"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// CloneFields clones values from input to output based on the field names of output type
// ignores fields listed in ignoreFields.
// It returns the cloned output.
func CloneFields[IN comparable, OUT comparable](input IN, output OUT, ignoreFields ...string) OUT {
	srcType := reflect.TypeOf(input).Elem()     // Get the type of the input
	destType := reflect.TypeOf(output).Elem()   // Get the type of the output
	srcValue := reflect.ValueOf(input).Elem()   // Get the value of the input
	destValue := reflect.ValueOf(output).Elem() // Get the value of the output

	// Loop through all fields of the output
	for i := 0; i < destType.NumField(); i++ {
		destField := destType.Field(i)                 // Get the field of the output
		name := destField.Name                         // Get the name of the field
		if ComparableContains(name, ignoreFields...) { // If the name is in ignoreFields, skip it
			continue
		}

		// Check if the field with the same name exists in the input
		if srcField, ok := srcType.FieldByName(name); ok {
			srcFieldValue := srcValue.FieldByName(name) // Get the value of the field from the input
			switch destField.Type {
			case TYPE_TIME:
				// If the type is time.Time, convert the input value to time.Time and set it to the output field
				value := reflect.ValueOf(TimeStampToTime((srcFieldValue.Interface()).(timestamppb.Timestamp)))
				destValue.FieldByName(name).Set(value)
			case TYPE_TIME_POINTER:
				// If the type is *time.Time, convert the input value to *time.Time and set it to the output field
				var value *time.Time = nil
				v := srcFieldValue.Elem()
				if v.IsValid() {
					i := v.Interface().(timestamppb.Timestamp)
					value = TimeStampToTimePointer(&i)
				}
				destValue.FieldByName(name).Set(reflect.ValueOf(value))
			case TYPE_TIMESTAMP:
				// If the type is timestamppb.Timestamp, convert the input value to timestamppb.Timestamp and set it to the output field
				value := reflect.ValueOf(TimeToTimeStamp((srcFieldValue.Interface()).(time.Time)))
				destValue.FieldByName(name).Set(value)
			case TYPE_TIMESTAMP_POINTER:
				// If the type is *timestamppb.Timestamp, convert the input value to *timestamppb.Timestamp and set it to the output field
				var value *timestamppb.Timestamp = nil
				v := srcFieldValue.Elem()
				if v.IsValid() {
					i := v.Interface().(time.Time)
					value = TimeToTimeStampPointer(&i)
				}
				destValue.FieldByName(name).Set(reflect.ValueOf(value))
			case srcField.Type:
				// If the types are the same, set the input value to the output field
				destValue.FieldByName(name).Set(srcValue.FieldByName(name))
			}
		}
	}

	return output // Return the cloned output
}
