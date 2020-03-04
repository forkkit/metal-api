package machine

import (
	"github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service/helper"
	v1 "github.com/metal-stack/metal-api/cmd/metal-api/internal/service/v1"
	"github.com/metal-stack/metal-lib/httperrors"
	"net/http"
)

func (r machineResource) addPowerMachineOnRoute(ws *restful.WebService, tags []string) {
	ws.Route(ws.POST("/{id}/power/on").
		To(helper.Editor(r.powerMachineOn)).
		Operation("machineOn").
		Doc("sends a power-on to the machine").
		Param(ws.PathParameter("id", "identifier of the machine").DataType("string")).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Reads(v1.EmptyBody{}).
		Returns(http.StatusOK, "OK", v1.MachineResponse{}).
		DefaultReturns("Error", httperrors.HTTPErrorResponse{}))
}

func (r machineResource) powerMachineOn(request *restful.Request, response *restful.Response) {
	r.publishMachineCmd("powerMachineOn", metal.MachineOnCmd, request, response)
}
