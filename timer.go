package utils

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

const longTimeFormat = "2006-01-02T15:04:05.999Z"

// ToLongTime converts a time.Time value to a string in long time format
func ToLongTime(value *time.Time) string {
	if value == nil {
		return ""
	}

	return value.Format(longTimeFormat)
}

// ToFormat converts a time.Time value to a string using the given format
func ToFormat(value *time.Time, format string) string {
	if value == nil {
		return ""
	}

	return value.Format(format)
}

// TimeStampToTime converts a timestamppb.Timestamp value to a time.Time value
func TimeStampToTime(timestamp timestamppb.Timestamp) time.Time {
	return timestamp.AsTime()
}

// TimeStampToTimePointer converts a pointer to timestamppb.Timestamp value to a pointer to time.Time value
func TimeStampToTimePointer(timestamp *timestamppb.Timestamp) *time.Time {
	if timestamp == nil {
		return nil
	}

	t := timestamp.AsTime()
	return &t
}

// TimeToTimeStamp converts a time.Time value to a timestamppb.Timestamp value
func TimeToTimeStamp(value time.Time) timestamppb.Timestamp {
	return *timestamppb.New(value)
}

// TimeToTimeStampPointer converts a pointer to time.Time value to a pointer to timestamppb.Timestamp value
func TimeToTimeStampPointer(value *time.Time) *timestamppb.Timestamp {
	if value == nil {
		return nil
	}

	result := TimeToTimeStamp(*value)
	return &result
}
