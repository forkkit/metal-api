package service

import (
	"github.com/emicklei/go-restful"
	"github.com/metal-stack/security"
)

type Return struct {
	Status  int
	Message string
	Model   interface{}
}

func NewReturn(httpStatus int, message string, model interface{}) *Return {
	return &Return{
		Status:  httpStatus,
		Message: message,
		Model:   model,
	}
}

type Route struct {
	Method         string
	SubPath        string
	PathParameter  *PathParameter
	PathParameters []*PathParameter
	Doc            string
	Access         []security.ResourceAccess
	Reads          interface{}
	Writes         interface{}
	Handler        restful.RouteFunction
	Returns        []*Return
	Ok             interface{}
	Error          interface{}
}
