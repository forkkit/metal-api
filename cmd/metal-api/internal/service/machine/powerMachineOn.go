package machine

import (
	"github.com/emicklei/go-restful"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
)

func (r machineResource) powerMachineOn(request *restful.Request, response *restful.Response) {
	r.publishMachineCmd("powerMachineOn", metal.MachineOnCmd, request, response)
}
