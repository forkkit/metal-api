package machine

import (
	"fmt"
	"github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metrics"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service/helper"
	v1 "github.com/metal-stack/metal-api/cmd/metal-api/internal/service/v1"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/utils"
	"github.com/metal-stack/metal-lib/httperrors"
	"github.com/metal-stack/metal-lib/zapup"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func (r machineResource) addCheckMachineLivelinessRoute(ws *restful.WebService, tags []string) {
	ws.Route(ws.POST("/liveliness").
		To(r.checkMachineLiveliness).
		Operation("checkMachineLiveliness").
		Doc("external trigger for evaluating machine liveliness").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Reads(v1.EmptyBody{}).
		Returns(http.StatusOK, "OK", v1.MachineLivelinessReport{}).
		DefaultReturns("Error", httperrors.HTTPErrorResponse{}))
}

func (r machineResource) checkMachineLiveliness(request *restful.Request, response *restful.Response) {
	logger := utils.Logger(request).Sugar()
	logger.Info("liveliness report was requested")

	machines, err := r.DS.ListMachines()
	if helper.CheckError(request, response, utils.CurrentFuncName(), err) {
		return
	}

	liveliness := make(metrics.PartitionLiveliness)

	unknown := 0
	alive := 0
	dead := 0
	for _, m := range machines {
		p := liveliness[m.PartitionID]
		lvlness, err := r.evaluateMachineLiveliness(m)
		if err != nil {
			logger.Errorw("cannot update liveliness", "error", err, "machine", m)
			// fall through, so the caller should get the evaulated state, although it is not persistet
		}
		switch lvlness {
		case metal.MachineLivelinessAlive:
			alive++
			p.Alive++
		case metal.MachineLivelinessDead:
			dead++
			p.Dead++
		default:
			unknown++
			p.Unknown++
		}
		liveliness[m.PartitionID] = p
	}

	report := v1.MachineLivelinessReport{
		AliveCount:   alive,
		DeadCount:    dead,
		UnknownCount: unknown,
	}

	metrics.ProvideLiveliness(liveliness)
	err = response.WriteHeaderAndEntity(http.StatusOK, report)
	if err != nil {
		zapup.MustRootLogger().Error("Failed to send response", zap.Error(err))
		return
	}
}

// EvaluateMachineLiveliness evaluates the liveliness of a given machine
func (r machineResource) evaluateMachineLiveliness(m metal.Machine) (metal.MachineLiveliness, error) {
	provisioningEvents, err := r.DS.FindProvisioningEventContainer(m.ID)
	if err != nil {
		// we have no provisioning events... we cannot tell
		return metal.MachineLivelinessUnknown, fmt.Errorf("no provisioningEvents found for ID: %s", m.ID)
	}

	old := *provisioningEvents

	if provisioningEvents.LastEventTime != nil {
		if time.Since(*provisioningEvents.LastEventTime) > metal.MachineDeadAfter {
			if m.Allocation != nil {
				// the machine is either dead or the customer did turn off the phone home service
				provisioningEvents.Liveliness = metal.MachineLivelinessUnknown
			} else {
				// the machine is just dead
				provisioningEvents.Liveliness = metal.MachineLivelinessDead
			}
		} else {
			provisioningEvents.Liveliness = metal.MachineLivelinessAlive
		}
		err = r.DS.UpdateProvisioningEventContainer(&old, provisioningEvents)
		if err != nil {
			return provisioningEvents.Liveliness, err
		}
	}

	return provisioningEvents.Liveliness, nil
}
