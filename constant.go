package utils

import (
	"reflect"
)

const (
	STRING_FORMAT_DATE_LONG  string = "2006-01-02 15:04:05"
	STRING_FORMAT_DATE_SHORT string = "2006-01-02"
	STRING_FORMAT_MONTH_YEAR string = "2006-01"
	STRING_FORMAT_YEAR       string = "2006"
)

var TYPE_TIME reflect.Type = nil
var TYPE_TIME_POINTER reflect.Type = nil
var TYPE_TIMESTAMP reflect.Type = nil
var TYPE_TIMESTAMP_POINTER reflect.Type = nil
var TYPE_GUID reflect.Type = nil
var TYPE_GUID_POINTER reflect.Type = nil
var TYPE_SQL_ERROR reflect.Type = nil

func SetType(time reflect.Type, timePtr reflect.Type, stamp reflect.Type, stampPtr reflect.Type, sqlError reflect.Type, uuid reflect.Type, uuidPtr reflect.Type) {
	TYPE_TIME = time
	TYPE_TIME_POINTER = timePtr
	TYPE_TIMESTAMP = stamp
	TYPE_TIMESTAMP_POINTER = stampPtr
	TYPE_SQL_ERROR = sqlError
	TYPE_GUID = uuid
	TYPE_GUID_POINTER = uuidPtr
}
