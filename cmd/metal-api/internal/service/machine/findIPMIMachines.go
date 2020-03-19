package machine

import (
	"github.com/emicklei/go-restful"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/datastore"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service/helper"
	v1 "github.com/metal-stack/metal-api/pkg/proto/v1"
	"github.com/metal-stack/metal-api/pkg/util"
	"github.com/metal-stack/metal-lib/zapup"
	"go.uber.org/zap"
	"net/http"
)

func (r *machineResource) findIPMIMachines(request *restful.Request, response *restful.Response) {
	var requestPayload v1.MachineFindRequest
	err := request.ReadEntity(&requestPayload)
	if service.CheckError(request, response, util.CurrentFuncName(), err) {
		return
	}

	resp, err := FindIPMIMachines(r.ds, requestPayload.MachineSearchQuery)
	if service.CheckError(request, response, util.CurrentFuncName(), err) {
		return
	}
	err = response.WriteHeaderAndEntity(http.StatusOK, resp)
	if err != nil {
		zapup.MustRootLogger().Error("Failed to send response", zap.Error(err))
	}
}

func FindIPMIMachines(ds *datastore.RethinkStore, query *v1.MachineSearchQuery) ([]*v1.MachineIPMIResponse, error) {
	var ms metal.Machines
	err := ds.SearchMachines(query, &ms)
	if err != nil {
		return nil, err
	}

	resp := makeMachineIPMIResponseList(ms, ds, zapup.MustRootLogger().Sugar())
	return resp, nil
}

func makeMachineIPMIResponseList(ms metal.Machines, ds *datastore.RethinkStore, logger *zap.SugaredLogger) []*v1.MachineIPMIResponse {
	sMap, pMap, iMap, ecMap := helper.GetReferencedEntityMaps(ds, logger)

	var result []*v1.MachineIPMIResponse

	for _, m := range ms {
		var s *metal.Size
		if m.SizeID != "" {
			sizeEntity := sMap[m.SizeID]
			s = &sizeEntity
		}
		var p *metal.Partition
		if m.PartitionID != "" {
			partitionEntity := pMap[m.PartitionID]
			p = &partitionEntity
		}
		var i *metal.Image
		if m.Allocation != nil {
			if m.Allocation.ImageID != "" {
				imageEntity := iMap[m.Allocation.ImageID]
				i = &imageEntity
			}
		}
		ec := ecMap[m.ID]
		result = append(result, NewMachineIPMIResponse(&m, s, p, i, &ec))
	}

	return result
}
