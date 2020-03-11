package partition

import (
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service/helper"
	v1 "github.com/metal-stack/metal-api/pkg/proto/v1"
	"github.com/metal-stack/metal-api/pkg/util"
	"github.com/metal-stack/metal-lib/zapup"
	"go.uber.org/zap"
	"net/http"
)

func (r *partitionResource) createPartition(request *restful.Request, response *restful.Response) {
	var requestPayload v1.PartitionCreateRequest
	err := request.ReadEntity(&requestPayload)
	if helper.CheckError(request, response, util.CurrentFuncName(), err) {
		return
	}

	if requestPayload.ID == "" {
		if helper.CheckError(request, response, util.CurrentFuncName(), fmt.Errorf("id should not be empty")) {
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
	var mgmtServiceAddress string
	if requestPayload.MgmtServiceAddress != nil {
		mgmtServiceAddress = *requestPayload.MgmtServiceAddress
	}
	prefixLength := 22
	if requestPayload.PrivateNetworkPrefixLength != nil {
		prefixLength = *requestPayload.PrivateNetworkPrefixLength
		if prefixLength < 16 || prefixLength > 30 {
			if helper.CheckError(request, response, util.CurrentFuncName(), fmt.Errorf("private network prefix length is out of range")) {
				return
			}
		}
	}
	var imageURL string
	if requestPayload.PartitionBootConfiguration.ImageURL != nil {
		imageURL = *requestPayload.PartitionBootConfiguration.ImageURL
	}
	var kernelURL string
	if requestPayload.PartitionBootConfiguration.KernelURL != nil {
		kernelURL = *requestPayload.PartitionBootConfiguration.KernelURL
	}
	var commandLine string
	if requestPayload.PartitionBootConfiguration.CommandLine != nil {
		commandLine = *requestPayload.PartitionBootConfiguration.CommandLine
	}

	p := &metal.Partition{
		Base: metal.Base{
			ID:          requestPayload.ID,
			Name:        name,
			Description: description,
		},
		MgmtServiceAddress:         mgmtServiceAddress,
		PrivateNetworkPrefixLength: prefixLength,
		BootConfiguration: metal.BootConfiguration{
			ImageURL:    imageURL,
			KernelURL:   kernelURL,
			CommandLine: commandLine,
		},
	}

	fqns := []string{metal.TopicMachine.GetFQN(p.GetID()), metal.TopicSwitch.GetFQN(p.GetID())}
	for _, fqn := range fqns {
		if err := r.topicCreater.CreateTopic(p.GetID(), fqn); err != nil {
			if helper.CheckError(request, response, util.CurrentFuncName(), err) {
				return
			}
		}
	}

	err = r.ds.CreatePartition(p)
	if helper.CheckError(request, response, util.CurrentFuncName(), err) {
		return
	}

	err = response.WriteHeaderAndEntity(http.StatusCreated, service.NewPartitionResponse(p))
	if err != nil {
		zapup.MustRootLogger().Error("Failed to send response", zap.Error(err))
		return
	}
}
