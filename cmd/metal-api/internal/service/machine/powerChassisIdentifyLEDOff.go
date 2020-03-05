package machine

import (
	"github.com/emicklei/go-restful"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
)

func (r machineResource) powerChassisIdentifyLEDOff(request *restful.Request, response *restful.Response) {
	r.publishMachineCmd("powerChassisIdentifyLEDOff", metal.ChassisIdentifyLEDOffCmd, request, response, request.PathParameter("description"))
}
