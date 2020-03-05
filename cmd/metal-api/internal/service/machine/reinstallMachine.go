package machine

import (
	"github.com/emicklei/go-restful"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service/helper"
	v1 "github.com/metal-stack/metal-api/cmd/metal-api/internal/service/v1"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/utils"
	"github.com/metal-stack/metal-lib/httperrors"
	"github.com/metal-stack/metal-lib/zapup"
)

func (r machineResource) reinstallMachine(request *restful.Request, response *restful.Response) {
	log := utils.Logger(request).Sugar()
	var requestPayload v1.MachineReinstallRequest
	err := request.ReadEntity(&requestPayload)
	if helper.CheckError(request, response, utils.CurrentFuncName(), err) {
		return
	}

	err = r.reinstallOrDeleteMachine(request, response, &requestPayload.ImageID)
	if err != nil {
		helper.SendError(log.Desugar(), response, utils.CurrentFuncName(), httperrors.InternalServerError(err))
	}
}

func (r machineResource) abortReinstall(machineID string) {
	log := zapup.MustRootLogger().Sugar()

	m, err := r.DS.FindMachineByID(machineID)
	if err != nil {
		log.Errorw("unable to find machine", "machineID", machineID, "error", err)
		return
	}

	if m.Allocation != nil && m.Allocation.Reinstall {
		old := *m
		m.Allocation.Reinstall = false
		err = r.DS.UpdateMachine(&old, m)
		if err != nil {
			log.Errorw("unable to find machine", "machineID", machineID, "error", err)
		}
	}

	err = helper.PublishMachineCmd(log, m, r, metal.MachineAbortReinstall)
	if err != nil {
		log.Errorw("unable to publish ’Abort Reinstall' command", "machineID", machineID, "error", err)
	}
}
