package machine

import (
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service/helper"
	v1 "github.com/metal-stack/metal-api/pkg/proto/v1"
	"github.com/metal-stack/metal-api/pkg/util"
	"github.com/metal-stack/metal-lib/zapup"
	"go.uber.org/zap"
	"net/http"
)

func (r *machineResource) setChassisIdentifyLEDState(request *restful.Request, response *restful.Response) {
	var requestPayload v1.ChassisIdentifyLEDState
	err := request.ReadEntity(&requestPayload)
	if helper.CheckError(request, response, util.CurrentFuncName(), err) {
		return
	}

	ledState, err := metal.LEDStateFrom(requestPayload.Value)
	if helper.CheckError(request, response, util.CurrentFuncName(), err) {
		return
	}

	if ledState == metal.LEDStateOff && requestPayload.Description == "" {
		// we want a cause why this chassis identify LED is off
		if helper.CheckError(request, response, util.CurrentFuncName(), fmt.Errorf("you must supply a description")) {
			return
		}
	}

	id := request.PathParameter("id")
	oldMachine, err := r.ds.FindMachineByID(id)
	if helper.CheckError(request, response, util.CurrentFuncName(), err) {
		return
	}

	newMachine := *oldMachine

	newMachine.LEDState = metal.ChassisIdentifyLEDState{
		Value:       ledState,
		Description: requestPayload.Description,
	}

	err = r.ds.UpdateMachine(oldMachine, &newMachine)
	if helper.CheckError(request, response, util.CurrentFuncName(), err) {
		return
	}

	err = response.WriteHeaderAndEntity(http.StatusOK, MakeResponse(&newMachine, r.ds, util.Logger(request).Sugar()))
	if err != nil {
		zapup.MustRootLogger().Error("Failed to send response", zap.Error(err))
		return
	}
}
