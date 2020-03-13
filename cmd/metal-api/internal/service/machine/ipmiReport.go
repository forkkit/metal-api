package machine

import (
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service"
	v1 "github.com/metal-stack/metal-api/pkg/proto/v1"
	"github.com/metal-stack/metal-api/pkg/util"
	"github.com/metal-stack/metal-lib/zapup"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

func (r *machineResource) ipmiReport(request *restful.Request, response *restful.Response) {
	var requestPayload v1.MachineIpmiReport
	log := util.Logger(request)
	logger := log.Sugar()
	err := request.ReadEntity(&requestPayload)
	if service.CheckError(request, response, util.CurrentFuncName(), err) {
		return
	}
	if requestPayload.PartitionID == "" {
		err := fmt.Errorf("given partition id was not found")
		service.CheckError(request, response, util.CurrentFuncName(), err)
		return
	}

	var ms metal.Machines
	err = r.ds.SearchMachines(&v1.MachineSearchQuery{}, &ms)
	if service.CheckError(request, response, util.CurrentFuncName(), err) {
		return
	}
	known := make(map[string]string, len(ms))
	for _, m := range ms {
		uuid := m.ID
		if uuid == "" {
			continue
		}
		known[uuid] = m.IPMI.Address
	}
	resp := v1.MachineIpmiReportResponse{
		UpdatedLeases: make(map[string]string),
		CreatedLeases: make(map[string]string),
	}
	// create empty machines for uuids that are not yet known to the metal-api
	const defaultIPMIPort = "623"
	for uuid, ip := range requestPayload.ActiveLeases {
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
		err = r.ds.CreateMachine(m)
		if err != nil {
			logger.Errorf("could not create machine", "id", uuid, "ipmi-ip", ip, "m", m, "err", err)
			continue
		}
		resp.CreatedLeases[uuid] = ip
	}
	// update machine ipmi data if ipmi ip changed
	for _, oldMachine := range ms {
		uuid := oldMachine.ID
		if uuid == "" {
			continue
		}
		// if oldmachine.uuid is not part of this update cycle skip it
		ip, ok := requestPayload.ActiveLeases[uuid]
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

		err = r.ds.UpdateMachine(&oldMachine, &newMachine)
		if err != nil {
			logger.Errorf("could not update machine", "id", uuid, "ip", ip, "machine", newMachine, "err", err)
			continue
		}
		resp.UpdatedLeases[uuid] = ip
	}
	err = response.WriteHeaderAndEntity(http.StatusOK, resp)
	if err != nil {
		zapup.MustRootLogger().Error("Failed to send response", zap.Error(err))
		return
	}
}
