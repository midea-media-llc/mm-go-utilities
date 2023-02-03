package utils

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

var Location *time.Location

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

func TimeToTimeStamp(value time.Time) timestamppb.Timestamp {
	return *timestamppb.New(value)
}

func TimeToTimeStampPointer(value *time.Time) *timestamppb.Timestamp {
	if value == nil {
		return nil
	}

	return timestamppb.New(*value)
}
