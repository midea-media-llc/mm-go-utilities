package utils

import (
	"fmt"
	"reflect"

	"github.com/360EntSecGroup-Skylar/excelize"
)

var MAP_EXCEL_COLUMN_INDEX = map[int]string{1: "A", 2: "B", 3: "C", 4: "D", 5: "E", 6: "F", 7: "G", 8: "H", 9: "I", 10: "J", 11: "K", 12: "L", 13: "M", 14: "N", 15: "O", 16: "P", 17: "Q", 18: "R", 19: "S", 20: "T", 21: "U", 22: "V", 23: "W", 24: "X", 25: "Y", 26: "Z"}

type ExportFieldConfig struct {
	Name   string
	Header string
}

func WriteArrayToFileXlsx(value interface{}, fileName *string, sheetName string, headers ...string) error {
	fileExcel := excelize.NewFile()
	fileExcel.SetSheetName("Sheet1", sheetName)

	arrs := ToInterfaceSlice(reflect.ValueOf(value))

	countMore := 1
	if len(headers) > 0 {
		for i, item := range headers {
			fileExcel.SetCellValue(sheetName, fmt.Sprintf("%s%v", MAP_EXCEL_COLUMN_INDEX[i+1], 1), item)
		}
		countMore = 2
	}

	for i, item := range arrs {
		v := reflect.ValueOf(item)

		for j := 0; j < v.NumField(); j++ {
			fileExcel.SetCellValue(sheetName, fmt.Sprintf("%s%v", MAP_EXCEL_COLUMN_INDEX[j+1], i+countMore), v.Field(j))
		}
	}

	*fileName = fmt.Sprintf("%v.xlsx", *fileName)
	return fileExcel.SaveAs(*fileName)
}

func WriteDataIntoFile[T any](header []ExportFieldConfig, f *excelize.File, sheetName string, startRow, startColumn int, data []T, style int, handleCellValue func(string, T, reflect.Value) interface{}) {
	// render header
	index := startColumn
	for _, e := range header {
		columnIndex := MAP_EXCEL_COLUMN_INDEX[index]
		f.SetCellValue(sheetName, fmt.Sprintf("%v%d", columnIndex, startRow), e.Header)
		f.SetCellStyle(sheetName, fmt.Sprintf("%v%d", columnIndex, startRow), fmt.Sprintf("%v%v", columnIndex, len(data)+startRow), style)
		index++
	}
	startRow++

	// render data
	for _, item := range data {
		index = startColumn
		dataValue := reflect.ValueOf(item).Elem()
		for _, e := range header {
			valueByName := dataValue.FieldByName(e.Name)
			if !valueByName.IsValid() && valueByName.IsZero() {
				continue
			}
			value := handleCellValue(e.Name, item, valueByName)
			columnIndex := MAP_EXCEL_COLUMN_INDEX[index]
			f.SetCellValue(sheetName, fmt.Sprintf("%s%d", columnIndex, startRow), value)
			index++
		}
		startRow++
	}
}
