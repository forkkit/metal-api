package firewall

import (
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service/helper"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service/machine"
	v1 "github.com/metal-stack/metal-api/pkg/proto/v1"
	"github.com/metal-stack/metal-api/pkg/util"
	"github.com/metal-stack/metal-lib/zapup"
	"go.uber.org/zap"
	"net/http"
)

func (r *firewallResource) allocateFirewall(request *restful.Request, response *restful.Response) {
	var requestPayload v1.FirewallCreateRequest
	err := request.ReadEntity(&requestPayload)
	if helper.CheckError(request, response, util.CurrentFuncName(), err) {
		return
	}
	allocReq := requestPayload.MachineAllocateRequest

	hostname := "metal"
	if allocReq.Hostname != nil {
		hostname = allocReq.Hostname.GetValue()
	}

	if allocReq.Networks != nil && len(allocReq.Networks) <= 0 {
		if helper.CheckError(request, response, util.CurrentFuncName(), fmt.Errorf("network ids cannot be empty")) {
			return
		}
	}
	ha := false
	if requestPayload.HA != nil {
		ha = requestPayload.HA.GetValue()
	}
	if ha {
		if helper.CheckError(request, response, util.CurrentFuncName(), fmt.Errorf("highly-available firewall not supported for the time being")) {
			return
		}
	}

	image, err := r.ds.FindImage(allocReq.ImageID)
	if helper.CheckError(request, response, util.CurrentFuncName(), err) {
		return
	}

	if !image.HasFeature(metal.ImageFeatureFirewall) {
		if helper.CheckError(request, response, util.CurrentFuncName(), fmt.Errorf("given image is not usable for a firewall, features: %s", image.ImageFeatureString())) {
			return
		}
	}

	sshPubKeys := make([]string, len(allocReq.SSHPubKeys))
	for i, pubKey := range allocReq.SSHPubKeys {
		sshPubKeys[i] = pubKey.GetValue()
	}

	tags := make([]string, len(allocReq.Tags))
	for i, tag := range allocReq.Tags {
		tags[i] = tag.GetValue()
	}

	ips := make([]string, len(allocReq.IPs))
	for i, ipAddress := range allocReq.IPs {
		ips[i] = ipAddress.GetValue()
	}

	networks := make([]v1.MachineAllocationNetwork, len(allocReq.Networks))
	//for i, nw := range allocReq.Networks { //TODO
	//	networks[i] = service.FromNetwork(nw)
	//}

	spec := machine.AllocationSpec{
		UUID:        allocReq.Common.Meta.GetId(),
		Name:        allocReq.Common.Name.GetValue(),
		Description: allocReq.Common.Description.GetValue(),
		Hostname:    hostname,
		ProjectID:   allocReq.ProjectID,
		PartitionID: allocReq.PartitionID,
		SizeID:      allocReq.SizeID,
		Image:       image,
		SSHPubKeys:  sshPubKeys,
		UserData:    allocReq.UserData.GetValue(),
		Tags:        tags,
		Networks:    networks,
		IPs:         ips,
		HA:          ha,
		IsFirewall:  true,
	}

	m, err := machine.Allocate(r.ds, r.ipamer, &spec, r.mdc)
	if helper.CheckError(request, response, util.CurrentFuncName(), err) {
		return
	}
	err = response.WriteHeaderAndEntity(http.StatusOK, machine.MakeResponse(m, r.ds, util.Logger(request).Sugar()))
	if err != nil {
		zapup.MustRootLogger().Error("Failed to send response", zap.Error(err))
		return
	}
}
