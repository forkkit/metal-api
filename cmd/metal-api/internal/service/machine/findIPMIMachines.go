package machine

import (
	"github.com/emicklei/go-restful"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/datastore"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service/helper"
	v1 "github.com/metal-stack/metal-api/pkg/proto/v1"
	"github.com/metal-stack/metal-api/pkg/util"
	"github.com/metal-stack/metal-lib/zapup"
	"go.uber.org/zap"
	"net/http"
)

func (r *machineResource) findIPMIMachines(request *restful.Request, response *restful.Response) {
	var requestPayload v1.MachineSearchQuery
	err := request.ReadEntity(&requestPayload)
	if helper.CheckError(request, response, util.CurrentFuncName(), err) {
		return
	}

	ms := metal.Machines{}
	err = r.ds.SearchMachines(&requestPayload, &ms)
	if helper.CheckError(request, response, util.CurrentFuncName(), err) {
		return
	}
	err = response.WriteHeaderAndEntity(http.StatusOK, makeMachineIPMIResponseList(ms, r.ds, util.Logger(request).Sugar()))
	if err != nil {
		zapup.MustRootLogger().Error("Failed to send response", zap.Error(err))
		return
	}
}

func makeMachineIPMIResponseList(ms metal.Machines, ds *datastore.RethinkStore, logger *zap.SugaredLogger) []*v1.MachineIPMIResponse {
	sMap, pMap, iMap, ecMap := getReferencedEntityMaps(ds, logger)

	var result []*v1.MachineIPMIResponse

	for index := range ms {
		var s *metal.Size
		if ms[index].SizeID != "" {
			sizeEntity := sMap[ms[index].SizeID]
			s = &sizeEntity
		}
		var p *metal.Partition
		if ms[index].PartitionID != "" {
			partitionEntity := pMap[ms[index].PartitionID]
			p = &partitionEntity
		}
		var i *metal.Image
		if ms[index].Allocation != nil {
			if ms[index].Allocation.ImageID != "" {
				imageEntity := iMap[ms[index].Allocation.ImageID]
				i = &imageEntity
			}
		}
		ec := ecMap[ms[index].ID]
		result = append(result, NewMachineIPMIResponse(&ms[index], s, p, i, &ec))
	}

	return result
}
