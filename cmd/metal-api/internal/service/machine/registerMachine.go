package machine

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/emicklei/go-restful"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service/helper"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service/sw"
	v1 "github.com/metal-stack/metal-api/pkg/proto/v1"
	"github.com/metal-stack/metal-api/pkg/util"
	"github.com/metal-stack/metal-lib/zapup"
	"go.uber.org/zap"
	"net/http"
)

func (r *machineResource) registerMachine(request *restful.Request, response *restful.Response) {
	var requestPayload v1.MachineRegisterRequest
	err := request.ReadEntity(&requestPayload)
	if service.CheckError(request, response, util.CurrentFuncName(), err) {
		return
	}

	if requestPayload.UUID == "" {
		if service.CheckError(request, response, util.CurrentFuncName(), fmt.Errorf("uuid cannot be empty")) {
			return
		}
	}

	partition, err := r.ds.FindPartition(requestPayload.PartitionID)
	if service.CheckError(request, response, util.CurrentFuncName(), err) {
		return
	}

	machineHardware := helper.NewMetalMachineHardware(requestPayload.Hardware)
	size, _, err := r.ds.FromHardware(machineHardware)
	if err != nil {
		size = metal.UnknownSize
		util.Logger(request).Sugar().Errorw("no size found for hardware, defaulting to unknown size", "hardware", machineHardware, "error", err)
	}

	m, err := r.ds.FindMachineByID(requestPayload.UUID)
	if err != nil && !metal.IsNotFound(err) {
		if service.CheckError(request, response, util.CurrentFuncName(), err) {
			return
		}
	}

	returnCode := http.StatusOK

	if m == nil {
		// machine is not in the database, create it
		name := fmt.Sprintf("%d-core/%s", machineHardware.CPUCores, humanize.Bytes(machineHardware.Memory))
		descr := fmt.Sprintf("a machine with %d core(s) and %s of RAM", machineHardware.CPUCores, humanize.Bytes(machineHardware.Memory))
		m = &metal.Machine{
			Base: metal.Base{
				ID:          requestPayload.UUID,
				Name:        name,
				Description: descr,
			},
			Allocation:  nil,
			SizeID:      size.ID,
			PartitionID: partition.ID,
			RackID:      requestPayload.RackID,
			Hardware:    machineHardware,
			State: metal.MachineState{
				Value: v1.MachineState_AVAILABLE,
				Description: "Machine just registered",
			},
			LEDState: metal.ChassisIdentifyLEDState{
				Value:       v1.ChassisIdentifyLEDState_LED_OFF,
				Description: "Machine just registered",
			},
			Tags: util.StringSlice(requestPayload.Tags),
			IPMI: NewMetalIPMI(requestPayload.IPMI),
		}

		if requestPayload.BIOS != nil {
			m.BIOS = metal.BIOS{
				Version: requestPayload.BIOS.Version,
				Vendor:  requestPayload.BIOS.Vendor,
				Date:    requestPayload.BIOS.Date,
			}
		}

		err = r.ds.CreateMachine(m)
		if service.CheckError(request, response, util.CurrentFuncName(), err) {
			return
		}

		returnCode = http.StatusCreated
	} else {
		// machine has already registered, update it
		old := *m

		m.SizeID = size.ID
		m.PartitionID = partition.ID
		m.RackID = requestPayload.RackID
		m.Hardware = machineHardware
		if requestPayload.BIOS != nil {
			m.BIOS.Version = requestPayload.BIOS.Version
			m.BIOS.Vendor = requestPayload.BIOS.Vendor
			m.BIOS.Date = requestPayload.BIOS.Date
		}
		m.IPMI = NewMetalIPMI(requestPayload.IPMI)

		err = r.ds.UpdateMachine(&old, m)
		if service.CheckError(request, response, util.CurrentFuncName(), err) {
			return
		}
	}

	ec, err := r.ds.FindProvisioningEventContainer(m.ID)
	if err != nil && !metal.IsNotFound(err) {
		if service.CheckError(request, response, util.CurrentFuncName(), err) {
			return
		}
	}

	if ec == nil {
		err = r.ds.CreateProvisioningEventContainer(&metal.ProvisioningEventContainer{
			Base:                         metal.Base{ID: m.ID},
			Liveliness:                   metal.MachineLivelinessAlive,
			IncompleteProvisioningCycles: "0"},
		)
		if service.CheckError(request, response, util.CurrentFuncName(), err) {
			return
		}
	}

	err = sw.ConnectMachineWithSwitches(r.ds, m)
	if service.CheckError(request, response, util.CurrentFuncName(), err) {
		return
	}
	err = response.WriteHeaderAndEntity(returnCode, MakeResponse(m, r.ds, util.Logger(request).Sugar()))
	if err != nil {
		zapup.MustRootLogger().Error("Failed to send response", zap.Error(err))
		return
	}
}
