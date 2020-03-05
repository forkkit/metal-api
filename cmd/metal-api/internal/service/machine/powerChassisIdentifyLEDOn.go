package machine

import (
	"github.com/emicklei/go-restful"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
)

func (r machineResource) powerChassisIdentifyLEDOn(request *restful.Request, response *restful.Response) {
	r.powerChassisIdentifyLEDOnWithDescription(request, response)
}

func (r machineResource) powerChassisIdentifyLEDOnWithDescription(request *restful.Request, response *restful.Response) {
	description := request.PathParameter("description")
	r.publishMachineCmd("powerChassisIdentifyLEDOn", metal.ChassisIdentifyLEDOnCmd, request, response, description)
}
