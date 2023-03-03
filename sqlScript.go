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
		// If `value` is a slice, create a script to declare the table only once
		objectToScriptDeclare(result, true, &[]string{}, rfType.Elem(), tableName, ignoreFields...)
		// Loop through each struct in the slice and create a script to insert its data into the table
		arrayToScriptData(result, true, ToInterfaceSlice(rfValue), tableName, ignoreFields...)
	} else {
		// If `value` is a single struct, create a script to declare the table and insert its data into the table
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
	// Create a new strings.Builder to hold the resulting SQL script.
	result := &strings.Builder{}

	// Create a slice to hold the SQL column declarations for each field.
	// Handle pointer types by dereferencing them.
	elemType = handleTypePointer(elemType)

	// Get the indexes of the fields to include.
	indexFields := findFieldsUsedIndex(elemType, ignoreFields...)

	// Iterate over each field in the struct.
	for _, i := range indexFields {
		// Get the field information.
		field := elemType.Field(i)

		// If the field is a slice, recursively call this function to generate SQL column declarations for its element type.
		if field.Type.Kind() == reflect.Slice {
			objectToScriptDeclare(result, true, &[]string{}, field.Type.Elem(), field.Name, ignoreFields...)
		} else if isStruct(field.Type) {
			objectToScriptDeclare(builder, false, fields, field.Type.Elem(), field.Name, ignoreFields...)
		} else {
			// Otherwise, generate a SQL column declaration for the field based on its name and type.
			fieldDeclaration := fmt.Sprintf("[%s] %s", SafeColumnName(field.Name), findSqlTypeByType(field.Type))
			*fields = append(*fields, fieldDeclaration)
		}
	}

	// If there are any fields, generate a SQL declare statement for the table variable using the given tableName parameter.
	if write && len(*fields) > 0 {
		result.WriteString(fmt.Sprintf("declare @$%s table (%s)\n", tableName, strings.Join(*fields, ",")))
	}

	// Return the resulting SQL script as a string.
	builder.WriteString(result.String())
}

// objectToScriptData converts a single struct instance to an SQL insert statement.
// It takes in the reflect.Value and reflect.Type of the struct instance,
// the table name, and an optional slice of field names to ignore.
// It returns the generated SQL insert statement as a string.
func objectToScriptData(builder *strings.Builder, write bool, values *[]string, elem reflect.Value, elemType reflect.Type, tableName string, ignoreFields ...string) {
	// Create a strings.Builder to store the generated SQL statement.
	result := &strings.Builder{}
	// Create a slice to store the values of each field in the struct instance.

	// Dereference any pointer type that may be present.
	elem, elemType = handleValueTypePointer(elem, elemType)

	// Get the index of the fields that should be included in the SQL insert statement.
	indexFields := findFieldsUsedIndex(elemType, ignoreFields...)
	for _, i := range indexFields {
		// Get the reflect.StructField and reflect.Value of the current field.
		fieldField := elemType.Field(i)
		fieldValue := elem.Field(i)
		// Get the kind of the field.
		fieldKind := fieldField.Type.Kind()

		// If the field is an array or slice, generate a separate SQL insert statement for each element.
		if fieldKind == reflect.Array || fieldKind == reflect.Slice {
			arrayToScriptData(result, true, ToInterfaceSlice(fieldValue), fieldField.Name, ignoreFields...)
		} else if isStruct(fieldField.Type) {
			objectToScriptData(builder, false, values, fieldValue.Elem(), fieldValue.Elem().Type(), fieldField.Name, ignoreFields...)
		} else {
			// Convert the field value to a SQL string.
			*values = append(*values, toSqlValue(fieldKind, fieldValue.Interface()))
		}
	}
	// If there are any values to insert, generate the SQL insert statement.
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
	// Iterate over each value in the slice
	for _, item := range values {
		// Get the reflect.Value, reflect.Type, and reflect.Kind of the value using handlePointer
		value, valueType, _ := handlePointer(item)

		// Call objectToScriptData to convert the value to a SQL script and append the result to the builder
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
func handlePointer(value interface{}) (resultValue reflect.Value, rfType reflect.Type, rfKind reflect.Kind) {
	// Get the reflect.Value and reflect.Type of the input value
	resultValue = reflect.ValueOf(value)
	rfType = reflect.TypeOf(value)

	// Get the reflect.Kind of the input value
	rfKind = rfType.Kind()

	// If the input value is a pointer, dereference it and update the reflect.Type and reflect.Kind accordingly
	if rfKind == reflect.Ptr {
		resultValue = resultValue.Elem()
		rfType = rfType.Elem()
		rfKind = rfType.Kind()
	}

	return resultValue, rfType, rfKind
}

// handleTypePointer returns the type pointed to by a pointer type, or the original type if it is not a pointer.
func handleTypePointer(pointerType reflect.Type) reflect.Type {
	// Check if the given type is a pointer type.
	if pointerType.Kind() == reflect.Pointer {
		// If it is, return the type that the pointer points to.
		return pointerType.Elem()
	}
	// If it's not a pointer, return the original type.
	return pointerType
}

// findFieldsUsedIndex returns the indices of the fields in the struct type that are not ignored.
// The ignoreFields parameter is a variadic argument that specifies the names of the fields to ignore.
func findFieldsUsedIndex(rfType reflect.Type, ignoreFields ...string) []int {
	var result []int

	// Loop over all the fields of the struct type.
	numFields := rfType.NumField()
	for i := 0; i < numFields; i++ {
		field := rfType.Field(i)

		// If the field is in the ignoreFields list, skip it.
		if ComparableContains(field.Name, ignoreFields...) {
			continue
		}

		// Otherwise, add the index of the field to the result slice.
		result = append(result, i)
	}

	return result
}

// findSqlTypeByType determines the SQL data type that should be used for a given Go data type.
// It uses reflection to inspect the type of the input and returns the appropriate SQL data type.
func findSqlTypeByType(modelType reflect.Type) string {
	// Use a switch statement to handle different Go data types.
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
		// If the input is a pointer, recursively call this function on the underlying type.
		return findSqlTypeByType(modelType.Elem())
	default:
		// If the input is not a supported data type, use the toOthersSqlType function to determine the appropriate SQL data type.
		return toOthersSqlType(modelType)
	}
}

// toOthersSqlType returns the corresponding SQL type for a given Go type
func toOthersSqlType(modelType reflect.Type) string {
	// If the modelType is a time.Time or *time.Time, return 'datetime'
	if modelType == reflect.TypeOf(time.Time{}) || modelType == reflect.TypeOf(&time.Time{}) {
		return "datetime"
	}

	// If the modelType is a timestamppb.Timestamp or *timestamppb.Timestamp, return 'datetime'
	if modelType == reflect.TypeOf(timestamppb.Timestamp{}) || modelType == reflect.TypeOf(&timestamppb.Timestamp{}) {
		return "datetime"
	}

	// For any other type, return 'nvarchar(max)'
	return "nvarchar(max)"
}

// toValuePointer converts the given interface value to a SQL string representation of a pointer type
func toValuePointer(value interface{}) string {
	// if the value is nil, return 'null'
	if value == nil {
		return "null"
	}

	// check if the value is of type *time.Time
	var timeNil *time.Time = nil
	if reflect.TypeOf(value) == reflect.TypeOf(timeNil) {
		// if the value is nil, return 'null'
		if value == timeNil {
			return "null"
		}
		// if the value is not nil, format it to 'YYYY-MM-DD HH:MM:SS' and return it within single quotes
		return fmt.Sprintf("'%s'", Safe((value.(*time.Time)).Format("2006-01-02 15:04:05")))
	}

	// check if the value is of type *timestamppb.Timestamp
	var timeNil2 *timestamppb.Timestamp = nil
	if reflect.TypeOf(value) == reflect.TypeOf(timeNil2) {
		// if the value is nil, return 'null'
		if value == timeNil2 {
			return "null"
		}
		// if the value is not nil, convert it to time.Time type and format it to 'YYYY-MM-DD HH:MM:SS' and return it within single quotes
		return fmt.Sprintf("'%s'", Safe(TimeStampToTimePointer(value.(*timestamppb.Timestamp)).Format("2006-01-02 15:04:05")))
	}

	// if the value is not of type *time.Time or *timestamppb.Timestamp, get its element value
	elemValue := reflect.ValueOf(value).Elem()

	// if the element value is invalid, return 'null'
	if elemValue.Kind() == reflect.Invalid {
		return "null"
	}

	// return the SQL string representation of the element value using the toSqlValue function
	return toSqlValue(reflect.TypeOf(value).Elem().Kind(), elemValue.Interface())
}

// toSqlValue converts the given interface value to a SQL string representation of the corresponding type.
func toSqlValue(kind reflect.Kind, value interface{}) string {
	// if the value is nil, return 'null'
	if value == nil {
		return "null"
	}

	// switch case statement based on the type of the interface value
	switch kind {
	case reflect.Invalid: // if the kind is Invalid, return 'null'
		return "null"
	case reflect.Bool: // if the kind is Bool, convert the value to a SQL string representation of a boolean type
		return toValueBool(value)
	case reflect.Int8, reflect.Int16, reflect.Uint8, reflect.Uint16: // if the kind is an integer type, convert the value to a SQL string representation of an integer type
		return toValueInt(value)
	case reflect.Int32, reflect.Uint32, reflect.Uint:
		return toValueInt32(value)
	case reflect.Int64, reflect.Uint64: // if the kind is a 64-bit integer type, convert the value to a SQL string representation of a 64-bit integer type
		return toValueInt64(value)
	case reflect.Float32: // if the kind is a 32-bit floating-point type, convert the value to a SQL string representation of a 32-bit floating-point type
		return toValueFloat32(value)
	case reflect.Float64: // if the kind is a 64-bit floating-point type, convert the value to a SQL string representation of a 64-bit floating-point type
		return toValueFloat64(value)
	case reflect.String: // if the kind is a string type, convert the value to a SQL string representation of a text type
		return toValueText(value)
	case reflect.Pointer: // if the kind is a pointer type, convert the value to a SQL string representation of the corresponding type
		return toValuePointer(value)
	case reflect.Struct: // if the kind is a struct type, convert the value to a SQL string representation of the corresponding type
		return toValueStruct(value)
	default: // if none of the above cases match, convert the value to a SQL string representation of a text type
		return toValueText(value)
	}
}

// toValueStruct converts the given interface value to a SQL string representation of a struct type
func toValueStruct(value interface{}) string {
	// if the value is nil, return 'null'
	if value == nil {
		return "null"
	}

	// if the value is of type time.Time, format it to 'YYYY-MM-DD HH:MM:SS' and return it within single quotes
	if reflect.TypeOf(value) == reflect.TypeOf(time.Time{}) {
		return fmt.Sprintf("'%s'", Safe((value.(time.Time)).Format("2006-01-02 15:04:05")))
	}

	// if the value is of type timestamppb.Timestamp, convert it to time.Time type and format it to 'YYYY-MM-DD HH:MM:SS'
	// and return it within single quotes
	if reflect.TypeOf(value) == reflect.TypeOf(timestamppb.Timestamp{}) {
		return fmt.Sprintf("'%s'", Safe(TimeStampToTime(value.(timestamppb.Timestamp)).Format("2006-01-02 15:04:05")))
	}

	// if none of the above cases match, return 'null'
	return "null"
}

// toValueText converts the given interface value to a SQL string representation of a text type
func toValueText(value interface{}) string {
	// format the value to a string representation and return it within an N'...' string with proper escaping
	return fmt.Sprintf("N'%s'", Safe(fmt.Sprintf("%v", value)))
}

// toValueInt converts the given interface value to a SQL string representation of an integer type
func toValueInt(value interface{}) string {
	// format the value to a string representation and return it
	return fmt.Sprintf("%v", value.(int))
}

// toValueInt32 converts the given interface value to a SQL string representation of a 32-bit integer type
func toValueInt32(value interface{}) string {
	// format the value to a string representation and return it
	return fmt.Sprintf("%v", value.(int32))
}

// toValueInt64 converts the given interface value to a SQL string representation of a 64-bit integer type
func toValueInt64(value interface{}) string {
	// format the value to a string representation and return it
	return fmt.Sprintf("%v", value.(int64))
}

// toValueFloat32 converts the given interface value to a SQL string representation of a 32-bit floating-point type
func toValueFloat32(value interface{}) string {
	// format the value to a string representation and return it
	return fmt.Sprintf("%v", value.(float32))
}

// toValueFloat64 converts the given interface value to a SQL string representation of a 64-bit floating-point type
func toValueFloat64(value interface{}) string {
	// format the value to a string representation and return it
	return fmt.Sprintf("%v", value.(float64))
}

// toValueBool converts the given interface value to a SQL string representation of a boolean type
func toValueBool(value interface{}) string {
	// check if the value is true or 'true' or 1 or '1' or float64(1), if so return '1', otherwise return '0'
	if value == true || value == "true" || value == 1 || value == "1" || (value == float64(1)) {
		return "1"
	}
	return "0"
}

func isStruct(fieldType reflect.Type) bool {
	if fieldType.Kind() != reflect.Ptr {
		return false
	}

	if fieldType.Elem().Kind() != reflect.Struct {
		return false
	}

	return !AnyContains(fieldType, TYPE_TIME_POINTER, TYPE_TIMESTAMP_POINTER)
}
