package sw

import (
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service/helper"
	"github.com/metal-stack/metal-api/pkg/helper"
	v1 "github.com/metal-stack/metal-api/pkg/proto/v1"
	"github.com/metal-stack/metal-lib/zapup"
	"go.uber.org/zap"
	"net/http"
)

func (r *switchResource) registerSwitch(request *restful.Request, response *restful.Response) {
	var requestPayload v1.SwitchRegisterRequest
	err := request.ReadEntity(&requestPayload)
	if helper.CheckError(request, response, helper.CurrentFuncName(), err) {
		return
	}

	if requestPayload.ID == "" {
		if helper.CheckError(request, response, helper.CurrentFuncName(), fmt.Errorf("uuid cannot be empty")) {
			return
		}
	}

	_, err = r.ds.FindPartition(requestPayload.PartitionID)
	if helper.CheckError(request, response, helper.CurrentFuncName(), err) {
		return
	}

	s, err := r.ds.FindSwitch(requestPayload.ID)
	if err != nil && !metal.IsNotFound(err) {
		if helper.CheckError(request, response, helper.CurrentFuncName(), err) {
			return
		}
	}

	returnCode := http.StatusOK

	if s == nil {
		s = service.ToSwitch(requestPayload)

		if len(requestPayload.Nics) != len(s.Nics.ByMac()) {
			if helper.CheckError(request, response, helper.CurrentFuncName(), fmt.Errorf("duplicate mac addresses found in nics")) {
				return
			}
		}

		err = r.ds.CreateSwitch(s)
		if helper.CheckError(request, response, helper.CurrentFuncName(), err) {
			return
		}

		// TODO: Broken switch: A machine was registered before this new switch is getting registered
		// It needs to take over the existing connections from the broken switch or something?
		// metal/metal#28

		returnCode = http.StatusCreated
	} else {
		old := *s

		spec := service.ToSwitch(requestPayload)

		if len(requestPayload.Nics) != len(spec.Nics.ByMac()) {
			if helper.CheckError(request, response, helper.CurrentFuncName(), fmt.Errorf("duplicate mac addresses found in nics")) {
				return
			}
		}

		nics, err := helper.UpdateSwitchNics(old.Nics.ByMac(), spec.Nics.ByMac(), old.MachineConnections)
		if helper.CheckError(request, response, helper.CurrentFuncName(), err) {
			return
		}

		if requestPayload.Name != nil {
			s.Name = *requestPayload.Name
		}
		if requestPayload.Description != nil {
			s.Description = *requestPayload.Description
		}
		s.RackID = spec.RackID
		s.PartitionID = spec.PartitionID

		s.Nics = nics
		// Do not replace connections here: We do not want to loose them!

		err = r.ds.UpdateSwitch(&old, s)

		if helper.CheckError(request, response, helper.CurrentFuncName(), err) {
			return
		}
	}
	err = response.WriteHeaderAndEntity(returnCode, helper.MakeSwitchResponse(s, r.ds, helper.Logger(request).Sugar()))
	if err != nil {
		zapup.MustRootLogger().Error("Failed to send response", zap.Error(err))
		return
	}
}
