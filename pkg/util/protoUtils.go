package util

import (
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/golang/protobuf/ptypes/wrappers"
	"time"
)

func Time(t *timestamp.Timestamp) time.Time {
	if t == nil {
		return time.Unix(0, int64(0))
	}
	return time.Unix(t.Seconds, int64(t.Nanos))
}

func StringSlice(slice []*wrappers.StringValue) []string {
	var ss []string
	for _, s := range slice {
		ss = append(ss, s.GetValue())
	}
	return ss
}

func TimestampProto(t time.Time) *timestamp.Timestamp {
	return &timestamp.Timestamp{
		Seconds: int64(t.Second()),
		Nanos:   int32(t.Nanosecond()),
	}
}

func BoolProto(value bool) *wrappers.BoolValue {
	return &wrappers.BoolValue{
		Value: value,
	}
}

func Int64Proto(value uint) *wrappers.Int64Value {
	return &wrappers.Int64Value{
		Value: int64(value),
	}
}

func Int64SliceProto(values ...int64) []*wrappers.Int64Value {
	var slice []*wrappers.Int64Value
	for _, value := range values {
		slice = append(slice, Int64Proto(uint(value)))
	}
	return slice
}

func UInt32Proto(value uint) *wrappers.UInt32Value {
	return &wrappers.UInt32Value{
		Value: uint32(value),
	}
}

func UInt32SliceProto(values ...uint) []*wrappers.UInt32Value {
	var slice []*wrappers.UInt32Value
	for _, value := range values {
		slice = append(slice, UInt32Proto(value))
	}
	return slice
}

func UInt64Proto(value uint) *wrappers.UInt64Value {
	return &wrappers.UInt64Value{
		Value: uint64(value),
	}
}

func UInt64SliceProto(values ...uint) []*wrappers.UInt64Value {
	var slice []*wrappers.UInt64Value
	for _, value := range values {
		slice = append(slice, UInt64Proto(value))
	}
	return slice
}

func StringProto(value string) *wrappers.StringValue {
	return &wrappers.StringValue{
		Value: value,
	}
}

func StringSliceProto(values ...string) []*wrappers.StringValue {
	var slice []*wrappers.StringValue
	for _, value := range values {
		slice = append(slice, StringProto(value))
	}
	return slice
}
