package ip

import (
	"context"
	"fmt"
	"github.com/emicklei/go-restful"
	mdmv1 "github.com/metal-stack/masterdata-api/api/v1"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service/helper"
	v1 "github.com/metal-stack/metal-api/pkg/proto/v1"
	"github.com/metal-stack/metal-api/pkg/util"
	"github.com/metal-stack/metal-lib/pkg/tag"
	"github.com/metal-stack/metal-lib/zapup"
	"go.uber.org/zap"
	"net/http"
)

func (r *ipResource) allocateIP(request *restful.Request, response *restful.Response) {
	r.allocateSpecificIP(request, response)
}

func (r *ipResource) allocateSpecificIP(request *restful.Request, response *restful.Response) {
	specificIP := request.PathParameter("ip")
	var requestPayload v1.IPAllocateRequest
	err := request.ReadEntity(&requestPayload)
	if helper.CheckError(request, response, util.CurrentFuncName(), err) {
		return
	}

	ip := requestPayload.IP

	if ip.NetworkID == "" {
		if helper.CheckError(request, response, util.CurrentFuncName(), fmt.Errorf("networkid should not be empty")) {
			return
		}
	}
	if ip.ProjectID == "" {
		if helper.CheckError(request, response, util.CurrentFuncName(), fmt.Errorf("projectid should not be empty")) {
			return
		}
	}

	nw, err := r.ds.FindNetworkByID(ip.NetworkID)
	if helper.CheckError(request, response, util.CurrentFuncName(), err) {
		return
	}

	p, err := r.mdc.Project().Get(context.Background(), &mdmv1.ProjectGetRequest{Id: ip.ProjectID})
	if helper.CheckError(request, response, util.CurrentFuncName(), err) {
		return
	}

	tags := make([]string, len(ip.Tags))
	for i, t := range ip.Tags {
		tags[i] = t.GetValue()
	}
	if requestPayload.MachineID != nil {
		tags = append(tags, v1.IpTag(tag.MachineID, requestPayload.MachineID.GetValue()))
	}

	tags, err = helper.ProcessTags(tags)
	if helper.CheckError(request, response, util.CurrentFuncName(), err) {
		return
	}

	// TODO: Following operations should span a database transaction if possible

	ipAddress, ipParentCidr, err := AllocateIP(nw, specificIP, r.ipamer)
	if helper.CheckError(request, response, util.CurrentFuncName(), err) {
		return
	}
	util.Logger(request).Sugar().Debugw("found an ip to allocate", "ip", ipAddress, "network", nw.ID)

	ipType := metal.Ephemeral
	if ip.Type == v1.IP_STATIC {
		ipType = metal.Static
	}

	metalIP := &metal.IP{
		IPAddress:        ipAddress,
		ParentPrefixCidr: ipParentCidr,
		Name:             ip.Common.Name.GetValue(),
		Description:      ip.Common.Description.GetValue(),
		NetworkID:        nw.ID,
		ProjectID:        p.GetProject().GetMeta().GetId(),
		Type:             ipType,
		Tags:             tags,
	}

	err = r.ds.CreateIP(metalIP)
	if helper.CheckError(request, response, util.CurrentFuncName(), err) {
		return
	}
	err = response.WriteHeaderAndEntity(http.StatusCreated, NewIPResponse(metalIP))
	if err != nil {
		zapup.MustRootLogger().Error("Failed to send response", zap.Error(err))
		return
	}
}
