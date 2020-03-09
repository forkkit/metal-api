package helper

import (
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/golang/protobuf/ptypes/wrappers"
	"time"
)

func ToTimestamp(t time.Time) *timestamp.Timestamp {
	return &timestamp.Timestamp{
		Seconds: int64(t.Second()),
		Nanos: int32(t.Nanosecond()),
	}
}

func ToBoolValue(value bool) *wrappers.BoolValue {
	return &wrappers.BoolValue{
		Value: value,
	}
}

func ToInt64Value(value int) *wrappers.Int64Value {
	return &wrappers.Int64Value{
		Value: int64(value),
	}
}

func ToStringValue(value string) *wrappers.StringValue {
	return &wrappers.StringValue{
		Value: value,
	}
}

func ToStringValueSlice(values ...string) []*wrappers.StringValue {
	slice := make([]*wrappers.StringValue, 0, len(values))
	for i, value := range values {
		slice[i] = ToStringValue(value)
	}
	return slice
}
