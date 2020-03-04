package partition

import (
	"github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service/helper"
	v1 "github.com/metal-stack/metal-api/cmd/metal-api/internal/service/v1"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/utils"
	"github.com/metal-stack/metal-lib/httperrors"
	"github.com/metal-stack/metal-lib/zapup"
	"go.uber.org/zap"
	"net/http"
)

func (r partitionResource) addUpdatePartitionRoute(ws *restful.WebService, tags []string) {
	ws.Route(ws.POST("/").
		To(helper.Admin(r.updatePartition)).
		Operation("updatePartition").
		Doc("updates a Partition. if the Partition was changed since this one was read, a conflict is returned").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Reads(v1.PartitionUpdateRequest{}).
		Returns(http.StatusOK, "OK", v1.PartitionResponse{}).
		Returns(http.StatusConflict, "Conflict", httperrors.HTTPErrorResponse{}).
		DefaultReturns("Error", httperrors.HTTPErrorResponse{}))
}

func (r partitionResource) updatePartition(request *restful.Request, response *restful.Response) {
	var requestPayload v1.PartitionUpdateRequest
	err := request.ReadEntity(&requestPayload)
	if helper.CheckError(request, response, utils.CurrentFuncName(), err) {
		return
	}

	oldPartition, err := r.DS.FindPartition(requestPayload.ID)
	if helper.CheckError(request, response, utils.CurrentFuncName(), err) {
		return
	}

	newPartition := *oldPartition

	if requestPayload.Name != nil {
		newPartition.Name = *requestPayload.Name
	}
	if requestPayload.Description != nil {
		newPartition.Description = *requestPayload.Description
	}
	if requestPayload.MgmtServiceAddress != nil {
		newPartition.MgmtServiceAddress = *requestPayload.MgmtServiceAddress
	}
	if requestPayload.PartitionBootConfiguration.ImageURL != nil {
		newPartition.BootConfiguration.ImageURL = *requestPayload.PartitionBootConfiguration.ImageURL
	}
	if requestPayload.PartitionBootConfiguration.KernelURL != nil {
		newPartition.BootConfiguration.KernelURL = *requestPayload.PartitionBootConfiguration.KernelURL
	}
	if requestPayload.PartitionBootConfiguration.CommandLine != nil {
		newPartition.BootConfiguration.CommandLine = *requestPayload.PartitionBootConfiguration.CommandLine
	}

	err = r.DS.UpdatePartition(oldPartition, &newPartition)
	if helper.CheckError(request, response, utils.CurrentFuncName(), err) {
		return
	}
	err = response.WriteHeaderAndEntity(http.StatusOK, v1.NewPartitionResponse(&newPartition))
	if err != nil {
		zapup.MustRootLogger().Error("Failed to send response", zap.Error(err))
		return
	}
}
