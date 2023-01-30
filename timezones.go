package utils

import (
	"os"
	"strings"
	"time"

	"github.com/midea-media-llc/mm-go-utilities/logs"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var Location *time.Location

func LoadTimeZones() {
	defaultTimeZones := os.Getenv("TIME_ZONES")
	if strings.Trim(defaultTimeZones, " ") == "" {
		defaultTimeZones = "Asia/Ho_Chi_Minh"
	}
	location, err := time.LoadLocation(os.Getenv("TIME_ZONES"))
	if err != nil {
		logs.Errorf("Cannot load Time zones", err)
	}
	Location = location
	logs.Infof("Time zones loaded " + os.Getenv("TIME_ZONES"))
}

func ToLongTime(value *time.Time) string {
	if value == nil {
		return ""
	}

	return value.Format("2006-01-02T15:04:05.999Z")
}

func ToFormat(value *time.Time, format string) string {
	if value == nil {
		return ""
	}

	return value.Format(format)
}

func TimeStampToTime(timestamp timestamppb.Timestamp) time.Time {
	return timestamp.AsTime()
}

func TimeStampToTimePointer(timestamp *timestamppb.Timestamp) *time.Time {
	time := timestamp.AsTime()
	return &time
}

func TimeToTimeStamp(value *time.Time) *timestamppb.Timestamp {
	if value == nil {
		return nil
	}

	return timestamppb.New(*value)
}
