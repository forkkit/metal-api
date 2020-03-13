package machine

import (
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service/sw"
	v1 "github.com/metal-stack/metal-api/pkg/proto/v1"
	"github.com/metal-stack/metal-api/pkg/util"
	"github.com/metal-stack/metal-lib/zapup"
	"go.uber.org/zap"
	"net/http"
)

func (r *machineResource) finalizeAllocation(request *restful.Request, response *restful.Response) {
	var requestPayload v1.MachineFinalizeAllocationRequest
	err := request.ReadEntity(&requestPayload)
	if service.CheckError(request, response, util.CurrentFuncName(), err) {
		return
	}

	id := request.PathParameter("id")
	m, err := r.ds.FindMachineByID(id)
	if service.CheckError(request, response, util.CurrentFuncName(), err) {
		return
	}

	if m.Allocation == nil {
		if service.CheckError(request, response, util.CurrentFuncName(), fmt.Errorf("the machine %q is not allocated", id)) {
			return
		}
	}

	old := *m

	m.Allocation.ConsolePassword = requestPayload.ConsolePassword
	if requestPayload.Setup != nil {
		m.Allocation.PrimaryDisk = requestPayload.Setup.PrimaryDisk
		m.Allocation.OSPartition = requestPayload.Setup.OsPartition
		m.Allocation.Initrd = requestPayload.Setup.Initrd
		m.Allocation.Cmdline = requestPayload.Setup.Cmdline
		m.Allocation.Kernel = requestPayload.Setup.Kernel
		m.Allocation.BootloaderID = requestPayload.Setup.BootloaderID
	}
	m.Allocation.Reinstall = false // just for safety

	err = r.ds.UpdateMachine(&old, m)
	if service.CheckError(request, response, util.CurrentFuncName(), err) {
		return
	}

	var sws []metal.Switch
	var vrf = ""
	imgs, err := r.ds.ListImages()
	if service.CheckError(request, response, util.CurrentFuncName(), err) {
		return
	}

	if m.IsFirewall(imgs.ByID()) {
		// if a machine has multiple networks, it serves as firewall, so it can not be enslaved into the tenant vrf
		vrf = "default"
	} else {
		for _, mn := range m.Allocation.MachineNetworks {
			if mn.Private {
				vrf = fmt.Sprintf("vrf%d", mn.Vrf)
				break
			}
		}
	}

	sws, err = sw.SetVrfAtSwitches(r.ds, m, vrf)
	if err != nil {
		if service.CheckError(request, response, util.CurrentFuncName(), fmt.Errorf("the machine %q could not be enslaved into the vrf %s", id, vrf)) {
			return
		}
	}

	if len(sws) > 0 {
		// Push out events to signal switch configuration change
		evt := metal.SwitchEvent{Type: metal.UPDATE, Machine: *m, Switches: sws}
		err = r.Publish(metal.TopicSwitch.GetFQN(m.PartitionID), evt)
		if err != nil {
			util.Logger(request).Sugar().Infow("switch update event could not be published", "event", evt, "error", err)
		}
	}
	err = response.WriteHeaderAndEntity(http.StatusOK, MakeResponse(m, r.ds, util.Logger(request).Sugar()))
	if err != nil {
		zapup.MustRootLogger().Error("Failed to send response", zap.Error(err))
		return
	}
}
