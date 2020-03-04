package service

import (
	"github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
	v1 "github.com/metal-stack/metal-api/cmd/metal-api/internal/service/v1"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/utils"
	"github.com/metal-stack/metal-lib/httperrors"
	"github.com/metal-stack/metal-lib/zapup"
	"net/http"
)

func (r machineResource) addReinstallMachineRoute(ws *restful.WebService, tags []string) {
	ws.Route(ws.POST("/{id}/reinstall").
		To(editor(r.reinstallMachine)).
		Operation("reinstallMachine").
		Doc("reinstall this machine").
		Param(ws.PathParameter("id", "identifier of the machine").DataType("string")).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Reads(v1.MachineReinstallRequest{}).
		Returns(http.StatusOK, "OK", v1.MachineResponse{}).
		Returns(http.StatusGatewayTimeout, "Timeout", httperrors.HTTPErrorResponse{}).
		DefaultReturns("Error", httperrors.HTTPErrorResponse{}))
}

func (r machineResource) reinstallMachine(request *restful.Request, response *restful.Response) {
	log := utils.Logger(request).Sugar()
	var requestPayload v1.MachineReinstallRequest
	err := request.ReadEntity(&requestPayload)
	if checkError(request, response, utils.CurrentFuncName(), err) {
		return
	}

	err = r.reinstallOrDeleteMachine(request, response, &requestPayload.ImageID)
	if err != nil {
		sendError(log.Desugar(), response, utils.CurrentFuncName(), httperrors.InternalServerError(err))
	}
}

func (r machineResource) abortReinstall(machineID string) {
	log := zapup.MustRootLogger().Sugar()

	m, err := r.ds.FindMachineByID(machineID)
	if err != nil {
		log.Errorw("unable to find machine", "machineID", machineID, "error", err)
		return
	}

	if m.Allocation != nil && m.Allocation.Reinstall {
		old := *m
		m.Allocation.Reinstall = false
		err = r.ds.UpdateMachine(&old, m)
		if err != nil {
			log.Errorw("unable to find machine", "machineID", machineID, "error", err)
		}
	}

	err = publishMachineCmd(log, m, r, metal.MachineAbortReinstall)
	if err != nil {
		log.Errorw("unable to publish ’Abort Reinstall' command", "machineID", machineID, "error", err)
	}
}
