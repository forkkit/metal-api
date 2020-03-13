package partition

import (
	"github.com/emicklei/go-restful"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service/helper"
	v1 "github.com/metal-stack/metal-api/pkg/proto/v1"
	"github.com/metal-stack/metal-api/pkg/util"
	"github.com/metal-stack/metal-lib/zapup"
	"go.uber.org/zap"
	"net/http"
)

func (r *partitionResource) updatePartition(request *restful.Request, response *restful.Response) {
	var requestPayload v1.PartitionUpdateRequest
	err := request.ReadEntity(&requestPayload)
	if service.CheckError(request, response, util.CurrentFuncName(), err) {
		return
	}

	partition := requestPayload.Partition

	oldPartition, err := r.ds.FindPartition(partition.Common.Meta.Id)
	if service.CheckError(request, response, util.CurrentFuncName(), err) {
		return
	}

	newPartition := *oldPartition

	newPartition.Name = partition.Common.Name.GetValue()
	newPartition.Description = partition.Common.Description.GetValue()
	newPartition.MgmtServiceAddress = partition.MgmtServiceAddress.GetValue()
	newPartition.BootConfiguration.ImageURL = partition.ImageURL.GetValue()
	newPartition.BootConfiguration.KernelURL = partition.KernelURL.GetValue()
	newPartition.BootConfiguration.CommandLine = partition.CommandLine.GetValue()

	err = r.ds.UpdatePartition(oldPartition, &newPartition)
	if service.CheckError(request, response, util.CurrentFuncName(), err) {
		return
	}
	err = response.WriteHeaderAndEntity(http.StatusOK, helper.NewPartitionResponse(&newPartition))
	if err != nil {
		zapup.MustRootLogger().Error("Failed to send response", zap.Error(err))
		return
	}
}
