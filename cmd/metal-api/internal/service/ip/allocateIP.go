package ip

import (
	"context"
	"fmt"
	"github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
	mdmv1 "github.com/metal-stack/masterdata-api/api/v1"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service/helper"
	v1 "github.com/metal-stack/metal-api/cmd/metal-api/internal/service/v1"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/utils"
	"github.com/metal-stack/metal-lib/httperrors"
	"github.com/metal-stack/metal-lib/zapup"
	"go.uber.org/zap"
	"net/http"
)

func (r ipResource) addAllocateIPRoute(ws *restful.WebService, tags []string) {
	ws.Route(ws.POST("/allocate").
		To(helper.Editor(r.allocateIP)).
		Operation("allocateIP").
		Doc("allocate an ip in the given network.").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Reads(v1.IPAllocateRequest{}).
		Writes(v1.IPResponse{}).
		Returns(http.StatusCreated, "Created", v1.IPResponse{}).
		DefaultReturns("Error", httperrors.HTTPErrorResponse{}))

	ws.Route(ws.POST("/allocate/{ip}").
		To(helper.Editor(r.allocateIP)).
		Operation("allocateSpecificIP").
		Param(ws.PathParameter("ip", "ip to try to allocate").DataType("string")).
		Doc("allocate a specific ip in the given network.").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Reads(v1.IPAllocateRequest{}).
		Writes(v1.IPResponse{}).
		Returns(http.StatusCreated, "Created", v1.IPResponse{}).
		DefaultReturns("Error", httperrors.HTTPErrorResponse{}))
}

func (ir ipResource) allocateIP(request *restful.Request, response *restful.Response) {
	specificIP := request.PathParameter("ip")
	var requestPayload v1.IPAllocateRequest
	err := request.ReadEntity(&requestPayload)
	if helper.CheckError(request, response, utils.CurrentFuncName(), err) {
		return
	}

	if requestPayload.NetworkID == "" {
		if helper.CheckError(request, response, utils.CurrentFuncName(), fmt.Errorf("networkid should not be empty")) {
			return
		}
	}
	if requestPayload.ProjectID == "" {
		if helper.CheckError(request, response, utils.CurrentFuncName(), fmt.Errorf("projectid should not be empty")) {
			return
		}
	}

	var name string
	if requestPayload.Name != nil {
		name = *requestPayload.Name
	}
	var description string
	if requestPayload.Description != nil {
		description = *requestPayload.Description
	}

	nw, err := ir.DS.FindNetworkByID(requestPayload.NetworkID)
	if helper.CheckError(request, response, utils.CurrentFuncName(), err) {
		return
	}

	p, err := ir.mdc.Project().Get(context.Background(), &mdmv1.ProjectGetRequest{Id: requestPayload.ProjectID})
	if helper.CheckError(request, response, utils.CurrentFuncName(), err) {
		return
	}

	tags := requestPayload.Tags
	if requestPayload.MachineID != nil {
		tags = append(tags, metal.IpTag(metal.TagIPMachineID, *requestPayload.MachineID))
	}

	tags, err = helper.ProcessTags(tags)
	if helper.CheckError(request, response, utils.CurrentFuncName(), err) {
		return
	}

	// TODO: Following operations should span a database transaction if possible

	ipAddress, ipParentCidr, err := helper.AllocateIP(nw, specificIP, ir.ipamer)
	if helper.CheckError(request, response, utils.CurrentFuncName(), err) {
		return
	}
	utils.Logger(request).Sugar().Debugw("found an ip to allocate", "ip", ipAddress, "network", nw.ID)

	ipType := metal.Ephemeral
	if requestPayload.Type == metal.Static {
		ipType = metal.Static
	}

	ip := &metal.IP{
		IPAddress:        ipAddress,
		ParentPrefixCidr: ipParentCidr,
		Name:             name,
		Description:      description,
		NetworkID:        nw.ID,
		ProjectID:        p.GetProject().GetMeta().GetId(),
		Type:             ipType,
		Tags:             tags,
	}

	err = ir.DS.CreateIP(ip)
	if helper.CheckError(request, response, utils.CurrentFuncName(), err) {
		return
	}
	err = response.WriteHeaderAndEntity(http.StatusCreated, v1.NewIPResponse(ip))
	if err != nil {
		zapup.MustRootLogger().Error("Failed to send response", zap.Error(err))
		return
	}
}
