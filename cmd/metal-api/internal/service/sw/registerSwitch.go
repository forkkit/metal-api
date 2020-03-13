package sw

import (
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service/helper"
	v1 "github.com/metal-stack/metal-api/pkg/proto/v1"
	"github.com/metal-stack/metal-api/pkg/util"
	"github.com/metal-stack/metal-lib/zapup"
	"go.uber.org/zap"
	"net/http"
)

func (r *switchResource) registerSwitch(request *restful.Request, response *restful.Response) {
	var requestPayload v1.SwitchRegisterRequest
	err := request.ReadEntity(&requestPayload)
	if service.CheckError(request, response, util.CurrentFuncName(), err) {
		return
	}

	sw := requestPayload.Switch

	if sw.Common.Meta.Id == "" {
		if service.CheckError(request, response, util.CurrentFuncName(), fmt.Errorf("uuid cannot be empty")) {
			return
		}
	}

	_, err = r.ds.FindPartition(requestPayload.PartitionID)
	if service.CheckError(request, response, util.CurrentFuncName(), err) {
		return
	}

	s, err := r.ds.FindSwitch(sw.Common.Meta.Id)
	if err != nil && !metal.IsNotFound(err) {
		if service.CheckError(request, response, util.CurrentFuncName(), err) {
			return
		}
	}

	returnCode := http.StatusOK

	spec := helper.FromSwitch(requestPayload)

	if len(sw.Nics) != len(spec.Nics.ByMac()) {
		if service.CheckError(request, response, util.CurrentFuncName(), fmt.Errorf("duplicate mac addresses found in nics")) {
			return
		}
	}

	if s == nil {
		s = spec

		err = r.ds.CreateSwitch(s)
		if service.CheckError(request, response, util.CurrentFuncName(), err) {
			return
		}

		// TODO: Broken switch: A machine was registered before this new switch is getting registered
		// It needs to take over the existing connections from the broken switch or something?
		// metal/metal#28

		returnCode = http.StatusCreated
	} else {
		old := *s

		nics, err := UpdateSwitchNics(old.Nics.ByMac(), spec.Nics.ByMac(), old.MachineConnections)
		if service.CheckError(request, response, util.CurrentFuncName(), err) {
			return
		}

		s.Name = sw.Common.Name.GetValue()
		s.Description = sw.Common.Description.GetValue()
		s.RackID = spec.RackID
		s.PartitionID = spec.PartitionID

		s.Nics = nics
		// Do not replace connections here: We do not want to loose them!

		err = r.ds.UpdateSwitch(&old, s)

		if service.CheckError(request, response, util.CurrentFuncName(), err) {
			return
		}
	}

	err = response.WriteHeaderAndEntity(returnCode, MakeSwitchResponse(s, r.ds, util.Logger(request).Sugar()))
	if err != nil {
		zapup.MustRootLogger().Error("Failed to send response", zap.Error(err))
		return
	}
}
