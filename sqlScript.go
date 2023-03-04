package utils

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// ToSqlScript converts a struct or slice of structs to a SQL script that can be used to declare and insert data
// into a SQL Server table.
//
// The `value` parameter can be either a struct or a slice of structs. If it's a struct, a script to declare the table
// will be created and then a script to insert the data into that table. If it's a slice of structs, a script to declare
// the table will be created only once and then a script to insert each struct in the slice into that table will be created.
//
// The `tableName` parameter specifies the name of the SQL Server table that the data will be inserted into.
//
// The `ignoreFields` parameter is an optional list of struct field names to exclude from the generated script.
func ToSqlScript(value interface{}, tableName string, ignoreFields ...string) string {
	result := &strings.Builder{}
	rfValue, rfType, rfKind := handlePointer(value)
	if rfKind == reflect.Array || rfKind == reflect.Slice {
		objectToScriptDeclare(result, true, &[]string{}, rfType.Elem(), tableName, ignoreFields...)
		arrayToScriptData(result, true, ToInterfaceSlice(rfValue), tableName, ignoreFields...)
	} else {
		objectToScriptDeclare(result, true, &[]string{}, rfType, tableName, ignoreFields...)
		objectToScriptData(result, true, &[]string{}, rfValue, rfType, tableName, ignoreFields...)
	}

	return result.String()
}

// objectToScriptDeclare generates SQL script for declaring a table variable based on a given struct type.
// It iterates over each field of the struct type and generates a SQL column declaration statement based on the
// field's name and type. If the field is a slice type, it recursively calls itself to generate column declarations
// for the slice's element type. Fields can be ignored using the `ignoreFields` parameter.
func objectToScriptDeclare(builder *strings.Builder, write bool, fields *[]string, elemType reflect.Type, tableName string, ignoreFields ...string) {
	result := &strings.Builder{}

	elemType = handleTypePointer(elemType)
	indexFields := findFieldsUsedIndex(elemType, ignoreFields...)
	for _, i := range indexFields {
		field := elemType.Field(i)
		fieldType := handleTypePointer(field.Type)
		if fieldType.Kind() == reflect.Slice {
			objectToScriptDeclare(result, true, &[]string{}, fieldType.Elem(), field.Name, ignoreFields...)
		} else if isStruct(fieldType) {
			objectToScriptDeclare(builder, false, fields, fieldType, field.Name, ignoreFields...)
		} else {
			*fields = append(*fields, fmt.Sprintf("[%s] %s", SafeColumnName(field.Name), findSqlTypeByType(field.Type)))
		}
	}

	if write && len(*fields) > 0 {
		result.WriteString(fmt.Sprintf("declare @$%s table (%s)\n", tableName, strings.Join(*fields, ",")))
	}

	builder.WriteString(result.String())
}

// objectToScriptData converts a single struct instance to an SQL insert statement.
// It takes in the reflect.Value and reflect.Type of the struct instance,
// the table name, and an optional slice of field names to ignore.
// It returns the generated SQL insert statement as a string.
func objectToScriptData(builder *strings.Builder, write bool, values *[]string, elem reflect.Value, elemType reflect.Type, tableName string, ignoreFields ...string) {
	result := &strings.Builder{}

	elem, elemType = handleValueTypePointer(elem, elemType)
	indexFields := findFieldsUsedIndex(elemType, ignoreFields...)
	for _, i := range indexFields {
		fieldField := elemType.Field(i)
		fieldValue := elem.Field(i)
		fieldKind := fieldValue.Kind()

		if ComparableContains(fieldKind, reflect.Array, reflect.Slice) {
			arrayToScriptData(result, true, ToInterfaceSlice(fieldValue.Elem()), fieldField.Name, ignoreFields...)
		} else if isStruct(handleTypePointer(fieldField.Type)) {
			objectToScriptData(builder, false, values, fieldValue, fieldField.Type, fieldField.Name, ignoreFields...)
		} else {
			*values = append(*values, toSqlValue(fieldValue.Kind(), fieldValue))
		}
	}

	if write && len(*values) > 0 {
		result.WriteString(fmt.Sprintf("insert into @$%s select %s\n", tableName, strings.Join(*values, ",")))
	}

	builder.WriteString(result.String())
}

// arrayToScriptData takes a slice of interface{} values, a table name, and an optional list of field names to ignore.
// It returns a string containing a SQL script to insert the data into a table.
func arrayToScriptData(builder *strings.Builder, write bool, values []interface{}, tableName string, ignoreFields ...string) {
	result := &strings.Builder{}

	datas := make([]string, 0)
	for _, item := range values {
		value, valueType, _ := handlePointer(item)
		objectToScriptData(result, write, &datas, value, valueType, tableName, ignoreFields...)
	}

	builder.WriteString(result.String())
}

// handleValueTypePointer checks if the given element type is a pointer type,
// and if it is, returns the value and type of its underlying element type.
// If it is not a pointer type, it returns the given element and its type as is.
func handleValueTypePointer(elem reflect.Value, elemType reflect.Type) (reflect.Value, reflect.Type) {
	// Check if the element type is a pointer type
	if elemType.Kind() == reflect.Ptr {
		// If it is a pointer type, return the value and type of its underlying element type
		return elem.Elem(), elemType.Elem()
	}

	// If it is not a pointer type, return the element and its type as is
	return elem, elemType
}

// handlePointer takes a value and returns its reflect.Value, reflect.Type, and reflect.Kind.
// If the input value is a pointer, it dereferences it first before returning.
// Parameters:
// - value: the input value to handle
// Returns:
// - resultValue: the reflect.Value of the input value, after being dereferenced if it was a pointer
// - rfType: the reflect.Type of the input value, after being dereferenced if it was a pointer
// - rfKind: the reflect.Kind of the input value, after being dereferenced if it was a pointer
func handlePointer(value interface{}) (reflect.Value, reflect.Type, reflect.Kind) {
	rfValue := handleValuePointer(reflect.ValueOf(value))
	rfType := handleTypePointer(reflect.TypeOf(value))
	return rfValue, rfType, rfType.Kind()
}

// handleTypePointer returns the type pointed to by a pointer type, or the original type if it is not a pointer.
func handleTypePointer(pointerType reflect.Type) reflect.Type {
	if ComparableContains(pointerType.Kind(), reflect.Ptr, reflect.Pointer) {
		return pointerType.Elem()
	}
	return pointerType
}

// handleValuePointer returns the value pointed to by a pointer value, or the original value if it is not a pointer.
func handleValuePointer(pointerValue reflect.Value) reflect.Value {
	if ComparableContains(pointerValue.Kind(), reflect.Ptr, reflect.Pointer) {
		return pointerValue.Elem()
	}
	return pointerValue
}

// findFieldsUsedIndex returns the indices of the fields in the struct type that are not ignored.
// The ignoreFields parameter is a variadic argument that specifies the names of the fields to ignore.
func findFieldsUsedIndex(rfType reflect.Type, ignoreFields ...string) []int {
	var result []int

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

// findFieldsUsedIndex returns the type is struct
func isStruct(fieldType reflect.Type) bool {
	if fieldType.Kind() != reflect.Struct {
		return false
	}

	return !ComparableContains(fieldType, TYPE_TIME, TYPE_TIMESTAMP, TYPE_TIME_POINTER, TYPE_TIMESTAMP_POINTER, TYPE_GUID, TYPE_GUID_POINTER)
}

// findSqlTypeByType determines the SQL data type that should be used for a given Go data type.
// It uses reflection to inspect the type of the input and returns the appropriate SQL data type.
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
	case reflect.Ptr:
		return findSqlTypeByType(modelType.Elem())
	default:
		return toOthersSqlType(modelType)
	}
}

// toOthersSqlType returns the corresponding SQL type for a given Go type
func toOthersSqlType(modelType reflect.Type) string {
	if ComparableContains(modelType, TYPE_TIME, TYPE_TIME_POINTER, TYPE_TIMESTAMP, TYPE_TIMESTAMP_POINTER) {
		return "datetime"
	}

	if ComparableContains(modelType, TYPE_GUID, TYPE_GUID_POINTER) {
		return "uniqueidentifier"
	}

	return "nvarchar(1)"
}

// toSqlValue converts the given interface value to a SQL string representation of the corresponding type.
func toSqlValue(kind reflect.Kind, value reflect.Value) string {
	if !value.IsValid() {
		return "null"
	}

	if ComparableContains(kind, reflect.Ptr, reflect.Pointer) {
		if value.IsNil() || !value.Elem().IsValid() {
			return "null"
		}

		value = value.Elem()
		kind = value.Kind()
	}

	switch kind {
	case reflect.Invalid:
		return "null"
	case reflect.Bool:
		return toValueBool(value)
	case reflect.Int8, reflect.Int16, reflect.Uint8, reflect.Uint16:
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
	case reflect.Struct:
		return toValueStruct(value)
	default:
		return toValueText(value)
	}
}

// toValueStruct converts the given interface value to a SQL string representation of a struct type
func toValueStruct(value reflect.Value) string {
	if value.Type() == TYPE_TIME {
		return fmt.Sprintf("'%s'", Safe((value.Interface().(time.Time)).Format("2006-01-02 15:04:05")))
	}

	if value.Type() == TYPE_TIMESTAMP {
		return fmt.Sprintf("'%s'", Safe(TimeStampToTime(value.Interface().(timestamppb.Timestamp)).Format("2006-01-02 15:04:05")))
	}

	return "null"
}

// toValueText converts the given interface value to a SQL string representation of a text type
func toValueText(value reflect.Value) string {
	return fmt.Sprintf("N'%s'", Safe(fmt.Sprintf("%v", value.Interface())))
}

// toValueInt converts the given interface value to a SQL string representation of an integer type
func toValueInt(value reflect.Value) string {
	return fmt.Sprintf("%v", value.Interface().(int))
}

// toValueInt32 converts the given interface value to a SQL string representation of a 32-bit integer type
func toValueInt32(value reflect.Value) string {
	return fmt.Sprintf("%v", value.Interface().(int32))
}

// toValueInt64 converts the given interface value to a SQL string representation of a 64-bit integer type
func toValueInt64(value reflect.Value) string {
	return fmt.Sprintf("%v", value.Interface().(int64))
}

// toValueFloat32 converts the given interface value to a SQL string representation of a 32-bit floating-point type
func toValueFloat32(value reflect.Value) string {
	return fmt.Sprintf("%v", value.Interface().(float32))
}

// toValueFloat64 converts the given interface value to a SQL string representation of a 64-bit floating-point type
func toValueFloat64(value reflect.Value) string {
	return fmt.Sprintf("%v", value.Interface().(float64))
}

// toValueBool converts the given interface value to a SQL string representation of a boolean type
func toValueBool(value reflect.Value) string {
	v := value.Interface()
	return fmt.Sprintf("%v", BoolToInt(v == true || v == "true" || v == 1 || v == "1" || v == float64(1)))
}
