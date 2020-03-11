package util

import (
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/golang/protobuf/ptypes/wrappers"
	"time"
)

func FromTimestamp(t *timestamp.Timestamp) time.Time {
	return time.Unix(t.Seconds, int64(t.Nanos))
}

func ToTimestamp(t time.Time) *timestamp.Timestamp {
	return &timestamp.Timestamp{
		Seconds: int64(t.Second()),
		Nanos:   int32(t.Nanosecond()),
	}
}

func ToBoolValue(value bool) *wrappers.BoolValue {
	return &wrappers.BoolValue{
		Value: value,
	}
}

func ToInt64Value(value uint) *wrappers.Int64Value {
	return &wrappers.Int64Value{
		Value: int64(value),
	}
}

func ToInt64ValueSlice(values ...int64) []*wrappers.Int64Value {
	slice := make([]*wrappers.Int64Value, len(values))
	for i, value := range values {
		slice[i] = ToInt64Value(uint(value))
	}
	return slice
}

func ToUInt64Value(value uint) *wrappers.UInt64Value {
	return &wrappers.UInt64Value{
		Value: uint64(value),
	}
}

func ToUInt64ValueSlice(values ...uint) []*wrappers.UInt64Value {
	slice := make([]*wrappers.UInt64Value, len(values))
	for i, value := range values {
		slice[i] = ToUInt64Value(value)
	}
	return slice
}

func ToStringValue(value string) *wrappers.StringValue {
	return &wrappers.StringValue{
		Value: value,
	}
}

func ToStringValueSlice(values ...string) []*wrappers.StringValue {
	slice := make([]*wrappers.StringValue, len(values))
	for i, value := range values {
		slice[i] = ToStringValue(value)
	}
	return slice
}
