package utils

import (
	"reflect"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
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

func init() {
	TYPE_TIME = reflect.TypeOf(time.Time{})
	TYPE_TIME_POINTER = reflect.TypeOf(&time.Time{})
	TYPE_TIMESTAMP = reflect.TypeOf(timestamp.Timestamp{})
	TYPE_TIMESTAMP_POINTER = reflect.TypeOf(&timestamp.Timestamp{})
}
