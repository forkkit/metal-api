package machine

import (
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/datastore"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service/helper"
	v1 "github.com/metal-stack/metal-api/cmd/metal-api/internal/service/v1"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/utils"
	"github.com/metal-stack/metal-lib/zapup"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

func (r machineResource) ipmiReport(request *restful.Request, response *restful.Response) {
	var requestPayload v1.MachineIpmiReport
	log := utils.Logger(request)
	logger := log.Sugar()
	err := request.ReadEntity(&requestPayload)
	if helper.CheckError(request, response, utils.CurrentFuncName(), err) {
		return
	}
	if requestPayload.PartitionID == "" {
		err := fmt.Errorf("given partition id was not found")
		helper.CheckError(request, response, utils.CurrentFuncName(), err)
		return
	}

	var ms metal.Machines
	err = r.DS.SearchMachines(&datastore.MachineSearchQuery{}, &ms)
	if helper.CheckError(request, response, utils.CurrentFuncName(), err) {
		return
	}
	known := v1.Leases{}
	for _, m := range ms {
		uuid := m.ID
		if uuid == "" {
			continue
		}
		known[uuid] = m.IPMI.Address
	}
	resp := v1.MachineIpmiReportResponse{
		Updated: v1.Leases{},
		Created: v1.Leases{},
	}
	// create empty machines for uuids that are not yet known to the metal-api
	const defaultIPMIPort = "623"
	for uuid, ip := range requestPayload.Leases {
		if uuid == "" {
			continue
		}
		if _, ok := known[uuid]; ok {
			continue
		}
		m := &metal.Machine{
			Base: metal.Base{
				ID: uuid,
			},
			PartitionID: requestPayload.PartitionID,
			IPMI: metal.IPMI{
				Address: ip + ":" + defaultIPMIPort,
			},
		}
		err = r.DS.CreateMachine(m)
		if err != nil {
			logger.Errorf("could not create machine", "id", uuid, "ipmi-ip", ip, "m", m, "err", err)
			continue
		}
		resp.Created[uuid] = ip
	}
	// update machine ipmi data if ipmi ip changed
	for _, oldMachine := range ms {
		uuid := oldMachine.ID
		if uuid == "" {
			continue
		}
		// if oldmachine.uuid is not part of this update cycle skip it
		ip, ok := requestPayload.Leases[uuid]
		if !ok {
			continue
		}
		newMachine := oldMachine

		// Replace host part of ipmi address with the ip from the ipmicatcher
		hostAndPort := strings.Split(oldMachine.IPMI.Address, ":")
		if len(hostAndPort) == 2 {
			newMachine.IPMI.Address = ip + ":" + hostAndPort[1]
		} else if len(hostAndPort) < 2 {
			newMachine.IPMI.Address = ip + ":" + defaultIPMIPort
		} else {
			logger.Errorf("not updating ipmi, address is garbage", "id", uuid, "ip", ip, "machine", newMachine, "address", newMachine.IPMI.Address)
			continue
		}

		if newMachine.IPMI.Address == oldMachine.IPMI.Address {
			continue
		}
		// machine was created by a PXE boot event and has no partition set.
		if oldMachine.PartitionID == "" {
			newMachine.PartitionID = requestPayload.PartitionID
		}

		if newMachine.PartitionID != requestPayload.PartitionID {
			logger.Errorf("could not update machine because overlapping id found", "id", uuid, "machine", newMachine, "partition", requestPayload.PartitionID)
			continue
		}

		err = r.DS.UpdateMachine(&oldMachine, &newMachine)
		if err != nil {
			logger.Errorf("could not update machine", "id", uuid, "ip", ip, "machine", newMachine, "err", err)
			continue
		}
		resp.Updated[uuid] = ip
	}
	err = response.WriteHeaderAndEntity(http.StatusOK, resp)
	if err != nil {
		zapup.MustRootLogger().Error("Failed to send response", zap.Error(err))
		return
	}
}
