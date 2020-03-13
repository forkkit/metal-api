package machine

import (
	"github.com/emicklei/go-restful"
	"github.com/golang/protobuf/ptypes"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service"
	v1 "github.com/metal-stack/metal-api/pkg/proto/v1"
	"github.com/metal-stack/metal-api/pkg/util"
	"github.com/metal-stack/metal-lib/zapup"
	"go.uber.org/zap"
	"net/http"
)

func (r *machineResource) publishMachineCmd(op string, cmd metal.MachineCommand, request *restful.Request, response *restful.Response, params ...string) {
	logger := util.Logger(request).Sugar()
	id := request.PathParameter("id")

	m, err := r.ds.FindMachineByID(id)
	if service.CheckError(request, response, op, err) {
		return
	}

	switch op {
	case "powerResetMachine", "powerMachineOff":
		event := string(metal.ProvisioningEventPlannedReboot)
		_, err = r.provisioningEventForMachine(id, v1.MachineProvisioningEvent{Time: ptypes.TimestampNow(), Event: event, Message: util.StringProto(op)})
		if service.CheckError(request, response, util.CurrentFuncName(), err) {
			return
		}
	}

	err = PublishMachineCmd(logger, m, r, cmd, params...)
	if service.CheckError(request, response, util.CurrentFuncName(), err) {
		return
	}

	err = response.WriteHeaderAndEntity(http.StatusOK, MakeResponse(m, r.ds, util.Logger(request).Sugar()))
	if err != nil {
		zapup.MustRootLogger().Error("Failed to send response", zap.Error(err))
		return
	}
}
