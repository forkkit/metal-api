package machine

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
	"time"
)

func (r machineResource) addAddProvisioningEventRoute(ws *restful.WebService, tags []string) {
	ws.Route(ws.POST("/{id}/event").
		To(helper.Editor(r.addProvisioningEvent)).
		Operation("addProvisioningEvent").
		Doc("adds a machine provisioning event").
		Param(ws.PathParameter("id", "identifier of the machine").DataType("string")).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Reads(v1.MachineProvisioningEvent{}).
		Returns(http.StatusOK, "OK", v1.MachineRecentProvisioningEvents{}).
		DefaultReturns("Error", httperrors.HTTPErrorResponse{}))
}

func (r machineResource) addProvisioningEvent(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("id")
	m, err := r.DS.FindMachineByID(id)
	if err != nil && !metal.IsNotFound(err) {
		if helper.CheckError(request, response, utils.CurrentFuncName(), err) {
			return
		}
	}

	// an event can actually create an empty machine. This enables us to also catch the very first PXE Booting event
	// in a machine lifecycle
	if m == nil {
		m = &metal.Machine{
			Base: metal.Base{
				ID: id,
			},
		}
		err = r.DS.CreateMachine(m)
		if helper.CheckError(request, response, utils.CurrentFuncName(), err) {
			return
		}
	}

	var requestPayload v1.MachineProvisioningEvent
	err = request.ReadEntity(&requestPayload)
	if helper.CheckError(request, response, utils.CurrentFuncName(), err) {
		return
	}
	ok := metal.AllProvisioningEventTypes[metal.ProvisioningEventType(requestPayload.Event)]
	if !ok {
		if helper.CheckError(request, response, utils.CurrentFuncName(), fmt.Errorf("unknown provisioning event")) {
			return
		}
	}

	ec, err := r.provisioningEventForMachine(id, requestPayload)
	if helper.CheckError(request, response, utils.CurrentFuncName(), err) {
		return
	}

	err = response.WriteHeaderAndEntity(http.StatusOK, v1.NewMachineRecentProvisioningEvents(ec))
	if err != nil {
		zapup.MustRootLogger().Error("Failed to send response", zap.Error(err))
		return
	}
}

func (r machineResource) provisioningEventForMachine(machineID string, e v1.MachineProvisioningEvent) (*metal.ProvisioningEventContainer, error) {
	ec, err := r.DS.FindProvisioningEventContainer(machineID)
	if err != nil && !metal.IsNotFound(err) {
		return nil, err
	}

	if ec == nil {
		ec = &metal.ProvisioningEventContainer{
			Base: metal.Base{
				ID: machineID,
			},
			Liveliness: metal.MachineLivelinessAlive,
		}
	}
	now := time.Now()
	ec.LastEventTime = &now

	event := metal.ProvisioningEvent{
		Time:    now,
		Event:   metal.ProvisioningEventType(e.Event),
		Message: e.Message,
	}
	if event.Event == metal.ProvisioningEventAlive {
		zapup.MustRootLogger().Sugar().Debugw("received provisioning alive event", "id", ec.ID)
		ec.Liveliness = metal.MachineLivelinessAlive
	} else if event.Event == metal.ProvisioningEventPhonedHome && len(ec.Events) > 0 && ec.Events[0].Event == metal.ProvisioningEventPhonedHome {
		zapup.MustRootLogger().Sugar().Debugw("swallowing repeated phone home event", "id", ec.ID)
		ec.Liveliness = metal.MachineLivelinessAlive
	} else if event.Event == metal.ProvisioningEventReinstallAborted {
		r.abortReinstall(machineID)
	} else {
		ec.Events = append([]metal.ProvisioningEvent{event}, ec.Events...)
		ec.IncompleteProvisioningCycles = ec.CalculateIncompleteCycles(zapup.MustRootLogger().Sugar())
		ec.Liveliness = metal.MachineLivelinessAlive
	}
	ec.TrimEvents(metal.ProvisioningEventsInspectionLimit)

	err = r.DS.UpsertProvisioningEventContainer(ec)
	return ec, err
}
