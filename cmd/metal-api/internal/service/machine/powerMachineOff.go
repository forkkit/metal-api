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

func (r machineResource) addPowerMachineOffRoute(ws *restful.WebService, tags []string) {
	ws.Route(ws.POST("/{id}/power/off").
		To(helper.Editor(r.powerMachineOff)).
		Operation("machineOff").
		Doc("sends a power-off to the machine").
		Param(ws.PathParameter("id", "identifier of the machine").DataType("string")).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Reads(v1.EmptyBody{}).
		Returns(http.StatusOK, "OK", v1.MachineResponse{}).
		DefaultReturns("Error", httperrors.HTTPErrorResponse{}))
}

func (r machineResource) powerMachineOff(request *restful.Request, response *restful.Response) {
	r.publishMachineCmd("powerMachineOff", metal.MachineOffCmd, request, response)
}
