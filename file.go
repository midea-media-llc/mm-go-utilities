package utils

import (
	"fmt"
	"reflect"
)

var MAP_EXCEL_COLUMN_INDEX = map[int]string{1: "A", 2: "B", 3: "C", 4: "D", 5: "E", 6: "F", 7: "G", 8: "H", 9: "I", 10: "J", 11: "K", 12: "L", 13: "M", 14: "N", 15: "O", 16: "P", 17: "Q", 18: "R", 19: "S", 20: "T", 21: "U", 22: "V", 23: "W", 24: "X", 25: "Y", 26: "Z"}

type IExcel interface {
	SaveAs(name string) error
	SetSheetName(oldName, newName string)
	SetCellValue(sheet, axis string, value interface{})
	SetCellStyle(sheet, hcell, vcell string, styleID int)
}

type ExportFieldConfig struct {
	Name   string
	Header string
}

// WriteArrayToFileXlsx writes the data contained in a slice to an xlsx file.
// The name of the file and the name of the sheet are specified in the arguments.
// If headers are specified, they are added as the first row of the sheet.
// The data is written starting at the second row of the sheet.
// The data is expected to be a slice of structs.
func WriteArrayToFileXlsx(f IExcel, value interface{}, fileName *string, sheetName string, headers ...string) error {
	f.SetSheetName("Sheet1", sheetName)
	arrs := ToInterfaceSlice(reflect.ValueOf(value))

	countMore := 1
	if len(headers) > 0 {
		for i, item := range headers {
			f.SetCellValue(sheetName, fmt.Sprintf("%s%v", MAP_EXCEL_COLUMN_INDEX[i+1], 1), item)
		}
		countMore = 2
	}

	for i, item := range arrs {
		v := reflect.ValueOf(item)

		for j := 0; j < v.NumField(); j++ {
			f.SetCellValue(sheetName, fmt.Sprintf("%s%v", MAP_EXCEL_COLUMN_INDEX[j+1], i+countMore), v.Field(j))
		}
	}

	return f.SaveAs(fmt.Sprintf("%v.xlsx", *fileName))
}

// WriteDataIntoFile writes the data contained in a slice of structs to an xlsx file.
// The headers are specified in the header argument.
// The data is written starting at the specified startRow and startColumn.
// The style is used to style the entire range of the data (excluding headers).
// The handleCellValue function is called for each cell in the data range and should return the cell value.
func WriteDataIntoFile[T any](header []ExportFieldConfig, f IExcel, sheetName string, startRow, startColumn int, data []T, style int, handleCellValue func(string, T, reflect.Value) interface{}) {
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
