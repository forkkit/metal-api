package partition

import (
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
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

	partition := requestPayload.Partition

	if partition.Common.Meta.Id == "" {
		if helper.CheckError(request, response, util.CurrentFuncName(), fmt.Errorf("id should not be empty")) {
			return
		}
	}

	prefixLength := uint32(22)
	if partition.PrivateNetworkPrefixLength != nil {
		prefixLength = partition.PrivateNetworkPrefixLength.GetValue()
		if prefixLength < 16 || prefixLength > 30 {
			if helper.CheckError(request, response, util.CurrentFuncName(), fmt.Errorf("private network prefix length is out of range")) {
				return
			}
		}
	}

	p := &metal.Partition{
		Base: metal.Base{
			ID:          partition.Common.Meta.Id,
			Name:        partition.Common.Name.GetValue(),
			Description: partition.Common.Description.GetValue(),
		},
		MgmtServiceAddress:         partition.MgmtServiceAddress.GetValue(),
		PrivateNetworkPrefixLength: uint(prefixLength),
		BootConfiguration: metal.BootConfiguration{
			ImageURL:    partition.ImageURL.GetValue(),
			KernelURL:   partition.KernelURL.GetValue(),
			CommandLine: partition.CommandLine.GetValue(),
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

	err = response.WriteHeaderAndEntity(http.StatusCreated, NewPartitionResponse(p))
	if err != nil {
		zapup.MustRootLogger().Error("Failed to send response", zap.Error(err))
		return
	}
}
