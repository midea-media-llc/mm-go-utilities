package utils

import (
	"fmt"
	"reflect"
)

func SplitArrayInt64(arrData []int64, limitLen int) [][]int64 {
	arrRes := [][]int64{}
	for i := 0; i < len(arrData); i += limitLen {
		batch := arrData[i:min(i+limitLen, len(arrData))]
		arrRes = append(arrRes, batch)
	}
	return arrRes
}

func RemoveAt(arrData []interface{}, index int) []interface{} {
	return append(arrData[:index], arrData[index+1:]...)
}

func StructToStringArray(source interface{}) []string {
	v := reflect.ValueOf(source)

	var res []string
	for j := 0; j < v.NumField(); j++ {
		strValue := fmt.Sprintf("%v", v.Field(j).Interface())
		res = append(res, strValue)
	}
	return res
}

func HandleFilterArrayInt(filterValue string) string {
	if filterValue == "" {
		filterValue = "-1"
	}
	if filterValue != "-1" {
		arrs, err := StringToArrayInt(filterValue, ",")
		if err != nil {
			return "-1"
		}

		filterValue = JoinArrayIntToString(arrs, ",")
	}

	return filterValue
}

func HandleFilterArrayIntPointer(filterValue *string) *string {
	if filterValue == nil {
		result := "-1"
		return &result
	}
	result := HandleFilterArrayInt(*filterValue)
	return &result
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}
