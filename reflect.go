package utils

import (
	"reflect"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func CloneFields[IN comparable, OUT comparable](input IN, output OUT, ignoreFields ...string) OUT {
	srcType := reflect.TypeOf(input).Elem()
	destType := reflect.TypeOf(output).Elem()
	srcValue := reflect.ValueOf(input).Elem()
	destValue := reflect.ValueOf(output).Elem()

	for i := 0; i < destType.NumField(); i++ {
		destField := destType.Field(i)
		name := destField.Name
		if ComparableContains(name, ignoreFields...) {
			continue
		}

		if srcField, ok := srcType.FieldByName(name); ok {
			srcFieldValue := srcValue.FieldByName(name)
			switch destField.Type {
			case TYPE_TIME:
				value := reflect.ValueOf(TimeStampToTime((srcFieldValue.Interface()).(timestamppb.Timestamp)))
				destValue.FieldByName(name).Set(value)
			case TYPE_TIME_POINTER:
				var value *time.Time = nil
				v := srcFieldValue.Elem()
				if v.IsValid() {
					value = TimeStampToTimePointer(v.Interface().(*timestamppb.Timestamp))
				}
				destValue.FieldByName(name).Set(reflect.ValueOf(value))
			case TYPE_TIMESTAMP:
				value := reflect.ValueOf(TimeToTimeStamp((srcFieldValue.Interface()).(time.Time)))
				destValue.FieldByName(name).Set(value)
			case TYPE_TIMESTAMP_POINTER:
				var value *timestamppb.Timestamp = nil
				v := srcFieldValue.Elem()
				if v.IsValid() {
					value = TimeToTimeStampPointer(v.Interface().(*time.Time))
				}
				destValue.FieldByName(name).Set(reflect.ValueOf(value))
			case srcField.Type:
				destValue.FieldByName(name).Set(srcValue.FieldByName(name))
			}
		}
	}

	return output
}
