package utils

import (
	"fmt"
	"os"
	"reflect"
	"time"

	"log"
)

const (
	DEFAULT_HEADER_STYLE = `{"font":{"size":12,"bold":true}, "alignment":{"horizontal":"center", "vertical":"center","wrap_text":true}, "border":[{"type":"left","color":"000000","style":1},{"type":"top","color":"000000","style":1},{"type":"bottom","color":"000000","style":1},{"type":"right","color":"000000","style":1}]}`
	DEFAULT_VALUE_STYLE  = `{"font":{"size":12}, "alignment":{"horizontal":"center", "vertical":"center","wrap_text":true}, "border":[{"type":"left","color":"000000","style":1},{"type":"top","color":"000000","style":1},{"type":"bottom","color":"000000","style":1},{"type":"right","color":"000000","style":1}]}`
)

var MAP_EXCEL_COLUMN_INDEX = map[int]string{
	1: "A", 2: "B", 3: "C", 4: "D", 5: "E", 6: "F", 7: "G", 8: "H", 9: "I", 10: "J",
	11: "K", 12: "L", 13: "M", 14: "N", 15: "O", 16: "P", 17: "Q", 18: "R", 19: "S",
	20: "T", 21: "U", 22: "V", 23: "W", 24: "X", 25: "Y", 26: "Z",
	27: "AA", 28: "AB", 29: "AC", 30: "AD", 31: "AE", 32: "AF", 33: "AG", 34: "AH",
	35: "AI", 36: "AJ", 37: "AK", 38: "AL", 39: "AM", 40: "AN", 41: "AO", 42: "AP",
	43: "AQ", 44: "AR", 45: "AS", 46: "AT", 47: "AU", 48: "AV", 49: "AW", 50: "AX",
	51: "AY", 52: "AZ",
	53: "BA", 54: "BB", 55: "BC", 56: "BD", 57: "BE", 58: "BF", 59: "BG", 60: "BH",
	61: "BI", 62: "BJ", 63: "BK", 64: "BL", 65: "BM", 66: "BN", 67: "BO", 68: "BP",
	69: "BQ", 70: "BR", 71: "BS", 72: "BT", 73: "BU", 74: "BV", 75: "BW", 76: "BX",
	77: "BY", 78: "BZ",
	79: "CA", 80: "CB", 81: "CC", 82: "CD", 83: "CE", 84: "CF", 85: "CG", 86: "CH",
	87: "CI", 88: "CJ", 89: "CK", 90: "CL", 91: "CM", 92: "CN", 93: "CO", 94: "CP",
	95: "CQ", 96: "CR", 97: "CS", 98: "CT", 99: "CU", 100: "CV", 101: "CW", 102: "CX",
	103: "CY", 104: "CZ",
	105: "DA", 106: "DB", 107: "DC", 108: "DD", 109: "DE", 110: "DF", 111: "DG", 112: "DH",
	113: "DI", 114: "DJ", 115: "DK", 116: "DL", 117: "DM", 118: "DN", 119: "DO", 120: "DP",
	121: "DQ", 122: "DR", 123: "DS", 124: "DT", 125: "DU", 126: "DV", 127: "DW", 128: "DX",
	129: "DY", 130: "DZ",
	131: "EA", 132: "EB", 133: "EC", 134: "ED", 135: "EE", 136: "EF", 137: "EG", 138: "EH",
	139: "EI", 140: "EJ", 141: "EK", 142: "EL", 143: "EM", 144: "EN", 145: "EO", 146: "EP",
	147: "EQ", 148: "ER", 149: "ES", 150: "ET", 151: "EU", 152: "EV", 153: "EW", 154: "EX",
	155: "EY", 156: "EZ",
	157: "FA", 158: "FB", 159: "FC", 160: "FD", 161: "FE", 162: "FF", 163: "FG", 164: "FH",
	165: "FI", 166: "FJ", 167: "FK", 168: "FL", 169: "FM", 170: "FN", 171: "FO", 172: "FP",
	173: "FQ", 174: "FR", 175: "FS", 176: "FT", 177: "FU", 178: "FV", 179: "FW", 180: "FX",
	181: "FY", 182: "FZ",
	183: "GA", 184: "GB", 185: "GC", 186: "GD", 187: "GE", 188: "GF", 189: "GG", 190: "GH",
	191: "GI", 192: "GJ", 193: "GK", 194: "GL", 195: "GM", 196: "GN", 197: "GO", 198: "GP",
	199: "GQ", 200: "GR", 201: "GS", 202: "GT", 203: "GU", 204: "GV", 205: "GW", 206: "GX",
	207: "GY", 208: "GZ",
	209: "HA", 300: "HB", 301: "HC", 302: "HD", 303: "HE", 304: "HF", 305: "HG", 306: "HH",
	307: "HI", 308: "HJ", 309: "HK", 310: "HL", 311: "HM", 312: "HN", 313: "HO", 314: "HP",
	315: "HQ", 316: "HR", 317: "HS", 318: "HT", 319: "HU", 320: "HV", 321: "HW", 322: "HX",
	323: "HY", 324: "HZ",
}

type IExcel interface {
	SaveAs(name string) error
	SetSheetName(oldName, newName string)
	NewStyle(style string) (int, error)
	SetCellValue(sheet, axis string, value interface{})
	SetCellStyle(sheet, hcell, vcell string, styleID int)
	SetColWidth(sheetName string, startCol string, endCol string, width float64)
	SetActiveSheet(index int)
	GetSheetIndex(name string) int
	NewSheet(name string) int
}

type IExportField interface {
	GetName() string
	GetHeader() string
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

	*fileName = fmt.Sprintf("%v.xlsx", *fileName)
	return f.SaveAs(*fileName)
}

// WriteDataIntoFile writes the data contained in a slice of structs to an xlsx file.
// The headers are specified in the header argument.
// The data is written starting at the specified startRow and startColumn.
// The style is used to style the entire range of the data (excluding headers).
// The handleCellValue function is called for each cell in the data range and should return the cell value.
func WriteDataIntoFile[T any](f IExcel, fields []IExportField, sheetName string, startRow, startColumn int, data []T, handleCellValue func(string, T, reflect.Value) interface{}) {
	headerStyle, _ := f.NewStyle(DEFAULT_HEADER_STYLE)
	cellStyle, _ := f.NewStyle(DEFAULT_VALUE_STYLE)

	// render header
	index := startColumn
	for _, e := range fields {
		columnIndex := MAP_EXCEL_COLUMN_INDEX[index]
		f.SetCellValue(sheetName, fmt.Sprintf("%s%d", columnIndex, startRow), e.GetHeader())
		f.SetCellStyle(sheetName, fmt.Sprintf("%s%d", columnIndex, startRow), fmt.Sprintf("%v%v", columnIndex, startRow), headerStyle)

		index++
	}

	startRow++
	// render data
	for _, item := range data {
		index = startColumn
		dataValue := reflect.ValueOf(item).Elem()
		for _, e := range fields {
			rValue := dataValue.FieldByName(e.GetName())
			value := handleCellValue(e.GetName(), item, rValue)
			columnIndex := MAP_EXCEL_COLUMN_INDEX[index]
			f.SetCellValue(sheetName, fmt.Sprintf("%s%d", columnIndex, startRow), value)
			f.SetCellStyle(sheetName, fmt.Sprintf("%s%d", columnIndex, startRow), fmt.Sprintf("%v%v", columnIndex, startRow), cellStyle)
			index++
		}
		startRow++
	}
}

// ExportXlsx is a generic function that exports data to an Excel file.
// It accepts an IExcel interface for writing to Excel, a slice of IExportField
// that specifies the fields to export, a sheetName string, startRow and startColumn integers,
// a slice of data of type T, a handleCellValue function that formats cell values, and a prefix string for the file name.
// The function writes the data to the specified sheet in the Excel file, saves the file, and returns the file bytes, file name and error (if any).
// If the specified sheetName does not exist in the Excel file, a new sheet with the specified name is created.
// The function also sets the column width of the sheet based on the number of fields.
// The handleCellValue function is called for each cell in the exported data and can be used to format the cell value.
// The prefix string is used to generate a unique file name for the exported file.
func ExportXlsx[T any](f IExcel, fields []IExportField, sheetName string, startRow, startColumn int, data []T, handleCellValue func(string, T, reflect.Value) interface{}, prefix string) ([]byte, string, error) {
	if f.GetSheetIndex(sheetName) == 0 {
		f.SetSheetName("Sheet1", sheetName)
	}

	sheetIndex := f.GetSheetIndex(sheetName)
	if sheetIndex == 0 {
		sheetIndex = f.NewSheet(sheetName)
	}

	f.SetColWidth(sheetName, "A", MAP_EXCEL_COLUMN_INDEX[len(fields)], 20)
	WriteDataIntoFile(f, fields, sheetName, startRow, startColumn, data, handleCellValue)

	fileName := fmt.Sprintf("%s-%d.xlsx", prefix, time.Now().Unix())

	err := f.SaveAs(fileName)
	if err != nil {
		log.Println(err.Error())
	}

	f.SetActiveSheet(sheetIndex)

	bytes, err := os.ReadFile(fileName)
	if err != nil {
		log.Println(err.Error())
	}

	err = os.Remove(fileName)
	if err != nil {
		log.Println(err.Error())
	}

	return bytes, fileName, nil
}
