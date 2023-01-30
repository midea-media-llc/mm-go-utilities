package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"time"
)

// Map fills out the fields in dest with values from source. All fields in the
// destination object must exist in the source object.
//
// Object hierarchies with nested structs and slices are supported, as long as
// type types of nested structs/slices follow the same rules, i.e. all fields
// in destination structs must be found on the source struct.
//
// Embedded/anonymous structs are supported
//
// Values that are not exported/not public will not be mapped.
//
// It is a design decision to panic when a field cannot be mapped in the
// destination to ensure that a renamed field in either the source or
// destination does not result in subtle silent bug.
func Map(source, dest interface{}) {
	var destType = reflect.TypeOf(dest)
	if destType.Kind() != reflect.Ptr {
		panic("Dest must be a pointer type")
	}
	var sourceVal = reflect.ValueOf(source)
	var destVal = reflect.ValueOf(dest).Elem()
	mapValues(sourceVal, destVal, false)
}

// MapLoose works just like Map, except it doesn't fail when the destination
// type contains fields not supplied by the source.
//
// This function is meant to be a temporary solution - the general idea is
// that the Map function should take a number of options that can modify its
// behavior - but I'd rather not add that functionality before I have a better
// idea what is a good options format.
func MapLoose(source, dest interface{}) {
	var destType = reflect.TypeOf(dest)
	if destType.Kind() != reflect.Ptr {
		panic("Dest must be a pointer type")
	}
	var sourceVal = reflect.ValueOf(source)
	var destVal = reflect.ValueOf(dest).Elem()
	mapValues(sourceVal, destVal, true)
}

func MapLoosePro(source, dest interface{}) error {
	var destType = reflect.TypeOf(dest)
	if destType.Kind() != reflect.Ptr {
		return errors.New("dest must be a pointer type")
	}
	var sourceVal = reflect.ValueOf(source)
	var destVal = reflect.ValueOf(dest).Elem()
	result := mapValues(sourceVal, destVal, true)

	if !result {
		return fmt.Errorf("cannot map %T to %T", source, dest)
	}

	return nil
}

func mapValues(sourceVal, destVal reflect.Value, loose bool) bool {
	destType := destVal.Type()

	if destType == reflect.TypeOf(time.Time{}) {
		switch sourceVal.Kind() {
		case reflect.String:
			if sourceVal.String() == "" {
				return false
			}

			parsedTime, err := time.Parse(time.RFC3339, sourceVal.String())
			if err != nil {
				panic("Cannot convert time from string using RFC3339 layout")
			}
			val := reflect.ValueOf(parsedTime)
			destVal.Set(val)
		case reflect.Float64:
			if sourceVal.Float() == 0 {
				return false
			}

			val := reflect.ValueOf(time.Unix(0, int64(sourceVal.Float())*int64(time.Millisecond)))
			destVal.Set(val)
		case reflect.Int64:
			if sourceVal.Int() == 0 {
				return false
			}

			val := reflect.ValueOf(time.Unix(0, sourceVal.Int()*int64(time.Millisecond)))
			destVal.Set(val)
		}
	} else if destType.Kind() == reflect.Struct {
		if sourceVal.Type().Kind() == reflect.Ptr {
			if sourceVal.IsNil() {
				// If source is nil, it maps to an empty struct
				sourceVal = reflect.New(sourceVal.Type().Elem())
			}
			sourceVal = sourceVal.Elem()
		}
		for i := 0; i < destVal.NumField(); i++ {
			mapField(sourceVal, destVal, i, loose)
		}
	} else if destType == sourceVal.Type() {
		destVal.Set(sourceVal)
	} else if destType.Kind() == reflect.Ptr {
		if valueIsNil(sourceVal) {
			return false
		}
		val := reflect.New(destType.Elem())
		shouldContinue := mapValues(sourceVal, val.Elem(), loose)

		if shouldContinue {
			destVal.Set(val)
		}
	} else if destType.Kind() == reflect.Slice {
		if sourceVal.Index(0).Interface() != nil && sourceVal.Index(0).Kind() == reflect.Struct {
			convertedSource, err := json.Marshal(sourceVal.Interface())

			if err != nil {
				convertedSource = []byte{}
			}

			newSourceVal := reflect.ValueOf(convertedSource)
			mapSlice(newSourceVal, destVal, loose)
		} else {
			mapSlice(sourceVal, destVal, loose)
		}
	} else {
		panic("Currently not supported")
	}

	return true
}

func mapSlice(sourceVal, destVal reflect.Value, loose bool) {
	destType := destVal.Type()
	length := sourceVal.Len()
	target := reflect.MakeSlice(destType, length, length)

	for j := 0; j < length; j++ {
		val := reflect.New(destType.Elem()).Elem()
		mapValues(sourceVal.Index(j), val, loose)
		target.Index(j).Set(val)
	}

	if length == 0 {
		verifyArrayTypesAreCompatible(sourceVal, destVal, loose)
	}

	destVal.Set(target)
}

func verifyArrayTypesAreCompatible(sourceVal, destVal reflect.Value, loose bool) {
	dummyDest := reflect.New(reflect.PtrTo(destVal.Type()))
	dummySource := reflect.MakeSlice(sourceVal.Type(), 1, 1)
	mapValues(dummySource, dummyDest.Elem(), loose)
}

func mapField(source, destVal reflect.Value, i int, loose bool) {
	destType := destVal.Type()
	fieldName := destType.Field(i).Name
	defer func() {
		if r := recover(); r != nil {
			panic(fmt.Sprintf("Error mapping field: %s. DestType: %v. SourceType: %v. Error: %v", fieldName, destType, source.Type(), r))
		}
	}()

	destField := destVal.Field(i)
	if destType.Field(i).Anonymous {
		mapValues(source, destField, loose)
	} else {
		if valueIsContainedInNilEmbeddedType(source, fieldName) {
			return
		}
		sourceField := source.FieldByName(fieldName)
		if (sourceField == reflect.Value{}) {
			if loose {
				return
			}
			if destField.Kind() == reflect.Struct {
				mapValues(source, destField, loose)
				return
			} else {
				for i := 0; i < source.NumField(); i++ {
					if source.Field(i).Kind() != reflect.Struct {
						continue
					}
					if sourceField = source.Field(i).FieldByName(fieldName); (sourceField != reflect.Value{}) {
						break
					}
				}
			}
		}
		mapValues(sourceField, destField, loose)
	}
}

func valueIsNil(value reflect.Value) bool {
	return value.Type().Kind() == reflect.Ptr && value.IsNil()
}

func valueIsContainedInNilEmbeddedType(source reflect.Value, fieldName string) bool {
	structField, _ := source.Type().FieldByName(fieldName)
	ix := structField.Index
	if len(structField.Index) > 1 {
		parentField := source.FieldByIndex(ix[:len(ix)-1])
		if valueIsNil(parentField) {
			return true
		}
	}
	return false
}

// StructToMap map struct to map and remove nil value, empty string
func StructToMap(source interface{}) (map[string]interface{}, error) {
	m := make(map[string]interface{})

	if reflect.TypeOf(source).Kind() != reflect.Ptr {
		return nil, errors.New("Source must be a pointer type")
	}

	elem := reflect.ValueOf(source).Elem()
	relType := elem.Type()

	for i := 0; i < relType.NumField(); i++ {
		field := elem.Field(i)
		if !field.IsNil() &&
			(reflect.TypeOf(field.Interface()).String() != "string" || !field.IsZero()) {
			m[relType.Field(i).Tag.Get("structmap")] = field.Interface()
		}
	}

	return m, nil
}
