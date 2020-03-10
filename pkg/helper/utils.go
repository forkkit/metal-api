package helper

import (
	"reflect"
	"runtime"
	"strconv"
	"strings"

	"github.com/emicklei/go-restful"
	"github.com/metal-stack/metal-lib/zapup"
	"go.uber.org/zap"
)

// Logger returns the request logger from the request.
func Logger(rq *restful.Request) *zap.Logger {
	return zapup.RequestLogger(rq.Request)
}

// GetFunctionName returns the name of the given function
func GetFunctionName(i interface{}) string {
	n := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
	idx := strings.LastIndex(n, ".")
	if idx > 0 {
		n = n[idx+1:]
	}
	idx = strings.Index(n, "-")
	if idx > 0 {
		n = n[:idx]
	}
	return n
}

// CurrentFuncName returns the name of the caller of this function.
func CurrentFuncName() string {
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		return "unknown"
	}
	ffpc := runtime.FuncForPC(pc)
	if ffpc == nil {
		return "unknown"
	}
	pp := strings.Split(ffpc.Name(), ".")
	return pp[len(pp)-1]
}

func SplitCIDR(cidr string) (string, *int) {
	parts := strings.Split(cidr, "/")
	if len(parts) == 2 {
		length, err := strconv.Atoi(parts[1])
		if err != nil {
			return parts[0], nil
		}
		return parts[0], &length
	}

	return cidr, nil
}
