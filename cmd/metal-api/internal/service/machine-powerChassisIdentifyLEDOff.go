package service

import (
	"github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
	v1 "github.com/metal-stack/metal-api/cmd/metal-api/internal/service/v1"
	"github.com/metal-stack/metal-lib/httperrors"
	"net/http"
)

func (r machineResource) addPowerChassisIdentifyLEDOffRoute(ws *restful.WebService, tags []string) {
	ws.Route(ws.POST("/{id}/power/chassis-identify-led-off/{description}").
		To(editor(r.powerChassisIdentifyLEDOff)).
		Operation("chassisIdentifyLEDOff").
		Doc("sends a power-off to the chassis identify LED").
		Param(ws.PathParameter("id", "identifier of the machine").DataType("string")).
		Param(ws.PathParameter("description", "reason why the chassis identify LED has been turned off").DataType("string")).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Reads(v1.EmptyBody{}).
		Returns(http.StatusOK, "OK", v1.MachineResponse{}).
		DefaultReturns("Error", httperrors.HTTPErrorResponse{}))
}

func (r machineResource) powerChassisIdentifyLEDOff(request *restful.Request, response *restful.Response) {
	r.publishMachineCmd("powerChassisIdentifyLEDOff", metal.ChassisIdentifyLEDOffCmd, request, response, request.PathParameter("description"))
}
