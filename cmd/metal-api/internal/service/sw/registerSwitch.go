package sw

import (
	"fmt"
	"github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service/helper"
	v1 "github.com/metal-stack/metal-api/cmd/metal-api/internal/service/v1"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/utils"
	"github.com/metal-stack/metal-lib/httperrors"
	"github.com/metal-stack/metal-lib/zapup"
	"go.uber.org/zap"
	"net/http"
)

func (r switchResource) addRegisterSwitchRoute(ws *restful.WebService, tags []string) {
	ws.Route(ws.POST("/register").
		To(helper.Editor(r.registerSwitch)).
		Doc("register a switch").
		Operation("registerSwitch").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Reads(v1.SwitchRegisterRequest{}).
		Returns(http.StatusOK, "OK", v1.SwitchResponse{}).
		Returns(http.StatusCreated, "Created", v1.SwitchResponse{}).
		DefaultReturns("Error", httperrors.HTTPErrorResponse{}))
}

func (r switchResource) registerSwitch(request *restful.Request, response *restful.Response) {
	var requestPayload v1.SwitchRegisterRequest
	err := request.ReadEntity(&requestPayload)
	if helper.CheckError(request, response, utils.CurrentFuncName(), err) {
		return
	}

	if requestPayload.ID == "" {
		if helper.CheckError(request, response, utils.CurrentFuncName(), fmt.Errorf("uuid cannot be empty")) {
			return
		}
	}

	_, err = r.DS.FindPartition(requestPayload.PartitionID)
	if helper.CheckError(request, response, utils.CurrentFuncName(), err) {
		return
	}

	s, err := r.DS.FindSwitch(requestPayload.ID)
	if err != nil && !metal.IsNotFound(err) {
		if helper.CheckError(request, response, utils.CurrentFuncName(), err) {
			return
		}
	}

	returnCode := http.StatusOK

	if s == nil {
		s = v1.NewSwitch(requestPayload)

		if len(requestPayload.Nics) != len(s.Nics.ByMac()) {
			if helper.CheckError(request, response, utils.CurrentFuncName(), fmt.Errorf("duplicate mac addresses found in nics")) {
				return
			}
		}

		err = r.DS.CreateSwitch(s)
		if helper.CheckError(request, response, utils.CurrentFuncName(), err) {
			return
		}

		// TODO: Broken switch: A machine was registered before this new switch is getting registered
		// It needs to take over the existing connections from the broken switch or something?
		// metal/metal#28

		returnCode = http.StatusCreated
	} else {
		old := *s

		spec := v1.NewSwitch(requestPayload)

		if len(requestPayload.Nics) != len(spec.Nics.ByMac()) {
			if helper.CheckError(request, response, utils.CurrentFuncName(), fmt.Errorf("duplicate mac addresses found in nics")) {
				return
			}
		}

		nics, err := helper.UpdateSwitchNics(old.Nics.ByMac(), spec.Nics.ByMac(), old.MachineConnections)
		if helper.CheckError(request, response, utils.CurrentFuncName(), err) {
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

		err = r.DS.UpdateSwitch(&old, s)

		if helper.CheckError(request, response, utils.CurrentFuncName(), err) {
			return
		}
	}
	err = response.WriteHeaderAndEntity(returnCode, helper.MakeSwitchResponse(s, r.DS, utils.Logger(request).Sugar()))
	if err != nil {
		zapup.MustRootLogger().Error("Failed to send response", zap.Error(err))
		return
	}
}
