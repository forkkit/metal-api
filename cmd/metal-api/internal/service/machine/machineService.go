package machine

import (
	"github.com/emicklei/go-restful"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service"
)

// webService creates the webservice endpoint
func (r machineResource) webService() *restful.WebService {
	ws := new(restful.WebService)
	ws.
		Path(service.BasePath + "v1/machine").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	tags := []string{"machine"}

	r.addFindMachineRoute(ws, tags)
	r.addFindMachinesRoute(ws, tags)
	r.addListMachinesRoute(ws, tags)

	r.addIPMIReportRoute(ws, tags)
	r.addFindIPMIMachineRoute(ws, tags)
	r.addFindIPMIMachinesRoute(ws, tags)

	r.addWaitForAllocationRoute(ws, tags)
	r.addAllocateMachineRoute(ws, tags)
	r.addFinalizeAllocationRoute(ws, tags)

	r.addRegisterMachineRoute(ws, tags)

	r.addReinstallMachineRoute(ws, tags)
	r.addFreeMachineRoute(ws, tags)

	r.addGetProvisioningEventContainerRoute(ws, tags)
	r.addAddProvisioningEventRoute(ws, tags)

	r.addCheckMachineLivelinessRoute(ws, tags)
	r.addSetMachineStateRoute(ws, tags)

	r.addPowerMachineOnRoute(ws, tags)
	r.addPowerMachineOffRoute(ws, tags)
	r.addPowerResetMachineRoute(ws, tags)
	r.addBootMachineBIOSRoute(ws, tags)

	r.addSetChassisIdentifyLEDStateRoute(ws, tags)
	r.addPowerChassisIdentifyLEDOnRoute(ws, tags)
	r.addPowerChassisIdentifyLEDOffRoute(ws, tags)

	return ws
}
