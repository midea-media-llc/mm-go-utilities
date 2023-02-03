package utils

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToInterfaceSlice(s reflect.Value) []interface{} {
	result := make([]interface{}, s.Len())
	if s.IsNil() {
		return result
	}

	for i := 0; i < s.Len(); i++ {
		item := s.Index(i)
		if item.Kind() == reflect.Pointer {
			result[i] = item.Elem().Interface()
		} else {
			result[i] = item.Interface()
		}
	}

	return result
}

func ToSqlScript(value interface{}, tableName string, ignoreFields ...string) string {
	result := strings.Builder{}
	rfValue, rfType, rfKind := handlePointer(value)

	if rfKind == reflect.Array || rfKind == reflect.Slice {
		result.WriteString(objectToScriptDeclare(rfType.Elem(), tableName, ignoreFields...))
		result.WriteString(arrayToScriptData(ToInterfaceSlice(rfValue), tableName, ignoreFields...))
	} else {
		result.WriteString(objectToScriptDeclare(rfType, tableName, ignoreFields...))
		result.WriteString(objectToScriptData(rfValue, rfType, tableName, ignoreFields...))
	}

	return result.String()
}

func objectToScriptDeclare(elemType reflect.Type, tableName string, ignoreFields ...string) string {
	result := strings.Builder{}
	fields := make([]string, 0)
	elemType = handleTypePointer(elemType)

	indexFields := findFieldsUsedIndex(elemType, ignoreFields...)
	for _, i := range indexFields {
		field := elemType.Field(i)
		if field.Type.Kind() == reflect.Slice {
			result.WriteString(objectToScriptDeclare(field.Type.Elem(), field.Name, ignoreFields...))
		} else {
			fields = append(fields, fmt.Sprintf("[%s] %s", SafeColumnName(field.Name), findSqlTypeByType(field.Type)))
		}
	}
	if len(fields) > 0 {
		result.WriteString(fmt.Sprintf("declare @$%s table (%s)\n", tableName, strings.Join(fields, ",")))
	}
	return result.String()
}

func objectToScriptData(elem reflect.Value, elemType reflect.Type, tableName string, ignoreFields ...string) string {
	result := strings.Builder{}
	values := make([]string, 0)
	elem, elemType = handleValueTypePointer(elem, elemType)

	indexFields := findFieldsUsedIndex(elemType, ignoreFields...)
	for _, i := range indexFields {
		fieldField := elemType.Field(i)
		fieldValue := elem.Field(i)
		fieldKind := fieldField.Type.Kind()

		if fieldKind == reflect.Array || fieldKind == reflect.Slice {
			result.WriteString(arrayToScriptData(ToInterfaceSlice(fieldValue), fieldField.Name, ignoreFields...))
		} else {
			values = append(values, toSqlValue(fieldKind, fieldValue.Interface()))
		}
	}
	if len(values) > 0 {
		result.WriteString(fmt.Sprintf("insert into @$%s select %s\n", tableName, strings.Join(values, ",")))
	}
	return result.String()
}

func arrayToScriptData(values []interface{}, tableName string, ignoreFields ...string) string {
	result := strings.Builder{}
	for _, item := range values {
		value, valueType, _ := handlePointer(item)
		result.WriteString(objectToScriptData(value, valueType, tableName, ignoreFields...))
	}
	return result.String()
}

func handleValueTypePointer(elem reflect.Value, elemType reflect.Type) (reflect.Value, reflect.Type) {
	if elemType.Kind() == reflect.Pointer {
		return elem.Elem(), elemType.Elem()
	}
	return elem, elemType
}

func handlePointer(value interface{}) (reflect.Value, reflect.Type, reflect.Kind) {
	resultValue := reflect.ValueOf(value)
	rfType := reflect.TypeOf(value)
	rfKind := rfType.Kind()
	if rfKind == reflect.Ptr {
		resultValue = resultValue.Elem()
		rfType = rfType.Elem()
		rfKind = rfType.Kind()
	}

	return resultValue, rfType, rfKind
}

func handleTypePointer(pointerType reflect.Type) reflect.Type {
	if pointerType.Kind() == reflect.Pointer {
		return pointerType.Elem()
	}
	return pointerType
}

func findFieldsUsedIndex(rfType reflect.Type, ignoreFields ...string) []int {
	result := make([]int, 0)
	numFields := rfType.NumField()
	for i := 0; i < numFields; i++ {
		field := rfType.Field(i)
		if ComparableContains(field.Name, ignoreFields...) {
			continue
		}
		result = append(result, i)
	}
	return result
}

func findSqlTypeByType(modelType reflect.Type) string {
	switch modelType.Kind() {
	case reflect.Bool:
		return "bit"
	case reflect.Int8, reflect.Int16, reflect.Uint8, reflect.Uint16:
		return "smallint"
	case reflect.Int32, reflect.Int, reflect.Uint32, reflect.Uint:
		return "int"
	case reflect.Int64, reflect.Uint64:
		return "bigint"
	case reflect.Float32, reflect.Float64:
		return "numeric(38,12)"
	case reflect.String:
		return "nvarchar(max)"
	case reflect.Pointer:
		return findSqlTypeByType(modelType.Elem())
	default:
		return toOthersSqlType(modelType)
	}
}

func toOthersSqlType(modelType reflect.Type) string {
	if modelType == reflect.TypeOf(time.Time{}) || modelType == reflect.TypeOf(&time.Time{}) {
		return "datetime"
	}

	if modelType == reflect.TypeOf(timestamppb.Timestamp{}) || modelType == reflect.TypeOf(&timestamppb.Timestamp{}) {
		return "datetime"
	}

	return "nvarchar(max)"
}

func toSqlValue(kind reflect.Kind, value interface{}) string {
	if value == nil {
		return "null"
	}

	switch kind {
	case reflect.Invalid:
		return "null"
	case reflect.Bool:
		return toValueBool(value)
	case reflect.Int8, reflect.Int16, reflect.Uint8, reflect.Uint16:
		return toValueInt(value)
	case reflect.Int:
		return toValueInt(value)
	case reflect.Int32, reflect.Uint32, reflect.Uint:
		return toValueInt32(value)
	case reflect.Int64, reflect.Uint64:
		return toValueInt64(value)
	case reflect.Float32:
		return toValueFloat32(value)
	case reflect.Float64:
		return toValueFloat64(value)
	case reflect.String:
		return toValueText(value)
	case reflect.Pointer:
		return toValuePointer(value)
	case reflect.Struct:
		return toValueStruct(value)
	default:
		return toValueText(value)
	}
}

func toValuePointer(value interface{}) string {
	if value == nil {
		return "null"
	}

	var timeNil *time.Time = nil
	if reflect.TypeOf(value) == reflect.TypeOf(timeNil) {
		if value == timeNil {
			return "null"
		}
		return fmt.Sprintf("'%s'", Safe((value.(*time.Time)).Format("2006-01-02 15:04:05")))
	}

	var timeNil2 *timestamppb.Timestamp = nil
	if reflect.TypeOf(value) == reflect.TypeOf(timeNil2) {
		if value == timeNil2 {
			return "null"
		}
		return fmt.Sprintf("'%s'", Safe(TimeStampToTimePointer(value.(*timestamppb.Timestamp)).Format("2006-01-02 15:04:05")))
	}

	elemValue := reflect.ValueOf(value).Elem()
	if elemValue.Kind() == reflect.Invalid {
		return "null"
	}
	return toSqlValue(reflect.TypeOf(value).Elem().Kind(), elemValue.Interface())
}

func toValueStruct(value interface{}) string {
	if value == nil {
		return "null"
	}

	if reflect.TypeOf(value) == reflect.TypeOf(time.Time{}) {
		return fmt.Sprintf("'%s'", Safe((value.(time.Time)).Format("2006-01-02 15:04:05")))
	}

	if reflect.TypeOf(value) == reflect.TypeOf(timestamppb.Timestamp{}) {
		return fmt.Sprintf("'%s'", Safe(TimeStampToTime(value.(timestamppb.Timestamp)).Format("2006-01-02 15:04:05")))
	}

	return "null"
}

func toValueText(value interface{}) string {
	return fmt.Sprintf("N'%s'", Safe(fmt.Sprintf("%v", value)))
}

func toValueInt(value interface{}) string {
	return fmt.Sprintf("%v", value.(int))
}

func toValueInt32(value interface{}) string {
	return fmt.Sprintf("%v", value.(int32))
}

func toValueInt64(value interface{}) string {
	return fmt.Sprintf("%v", value.(int64))
}

func toValueFloat32(value interface{}) string {
	return fmt.Sprintf("%v", value.(float32))
}

func toValueFloat64(value interface{}) string {
	return fmt.Sprintf("%v", value.(float64))
}

func toValueBool(value interface{}) string {
	if value == true || value == "true" || value == 1 || value == "1" {
		return "1"
	}
	return "0"
}
