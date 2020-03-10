package machine

import (
	"github.com/emicklei/go-restful"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service/helper"
	"github.com/metal-stack/metal-api/pkg/helper"
	v1 "github.com/metal-stack/metal-api/pkg/proto/v1"
	"github.com/metal-stack/metal-lib/zapup"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func (r *machineResource) publishMachineCmd(op string, cmd metal.MachineCommand, request *restful.Request, response *restful.Response, params ...string) {
	logger := helper.Logger(request).Sugar()
	id := request.PathParameter("id")

	m, err := r.ds.FindMachineByID(id)
	if helper.CheckError(request, response, op, err) {
		return
	}

	switch op {
	case "powerResetMachine", "powerMachineOff":
		event := string(metal.ProvisioningEventPlannedReboot)
		_, err = r.provisioningEventForMachine(id, v1.MachineProvisioningEvent{Time: time.Now(), Event: event, Message: op})
		if helper.CheckError(request, response, helper.CurrentFuncName(), err) {
			return
		}
	}

	err = helper.PublishMachineCmd(logger, m, r, cmd, params...)
	if helper.CheckError(request, response, helper.CurrentFuncName(), err) {
		return
	}

	err = response.WriteHeaderAndEntity(http.StatusOK, helper.MakeMachineResponse(m, r.ds, helper.Logger(request).Sugar()))
	if err != nil {
		zapup.MustRootLogger().Error("Failed to send response", zap.Error(err))
		return
	}
}
