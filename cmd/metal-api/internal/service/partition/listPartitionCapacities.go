package partition

import (
	"github.com/emicklei/go-restful"
	mdmv1 "github.com/metal-stack/masterdata-api/api/v1"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service/helper"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service/machine"
	v1 "github.com/metal-stack/metal-api/pkg/proto/v1"
	"github.com/metal-stack/metal-api/pkg/util"
	"github.com/metal-stack/metal-lib/zapup"
	"go.uber.org/zap"
	"net/http"
)

type ServerCapacity struct {
	Size      string `json:"size" description:"the size of the server"`
	Total     int    `json:"total" description:"total amount of servers with this size"`
	Free      int    `json:"free" description:"free servers with this size"`
	Allocated int    `json:"allocated" description:"allocated servers with this size"`
	Faulty    int    `json:"faulty" description:"servers with issues with this size"`
}

type PartitionCapacity struct {
	*v1.Common
	ServerCapacities []ServerCapacity `json:"servers" description:"servers available in this partition"`
}

func (r *partitionResource) listPartitionCapacities(request *restful.Request, response *restful.Response) {
	partitionCapacities, err := r.calcPartitionCapacities()

	if helper.CheckError(request, response, util.CurrentFuncName(), err) {
		return
	}
	err = response.WriteHeaderAndEntity(http.StatusOK, partitionCapacities)
	if err != nil {
		zapup.MustRootLogger().Error("Failed to send response", zap.Error(err))
		return
	}
}

func (r *partitionResource) calcPartitionCapacities() ([]PartitionCapacity, error) {
	// FIXME bad workaround to be able to run make spec
	if r.ds == nil {
		return nil, nil
	}
	ps, err := r.ds.ListPartitions()
	if err != nil {
		return nil, err
	}
	ms, err := r.ds.ListMachines()
	if err != nil {
		return nil, err
	}
	machines := machine.MakeMachineResponseList(ms, r.ds, zapup.MustRootLogger().Sugar())

	var partitionCapacities []PartitionCapacity
	for _, p := range ps {
		capacities := make(map[string]ServerCapacity)
		for _, machineResponse := range machines {
			m := machineResponse.Machine
			if m.PartitionResponse == nil {
				continue
			}
			if m.PartitionResponse.Partition.Common.Meta.Id != p.ID {
				continue
			}
			size := "unknown"
			if m.SizeResponse != nil {
				size = m.SizeResponse.Size.Common.Meta.Id
			}
			available := false
			if len(m.RecentProvisioningEvents.Events) > 0 {
				events := m.RecentProvisioningEvents.Events
				if metal.ProvisioningEventWaiting.Is(events[0].Event) && metal.ProvisioningEventAlive.Is(m.Liveliness) {
					available = true
				}
			}
			oldCap, ok := capacities[size]
			total := 1
			free := 0
			allocated := 0
			faulty := 0
			if ok {
				total = oldCap.Total + 1
			}

			if m.Allocation != nil {
				allocated = 1
			}
			if machine.MachineHasIssues(machineResponse) {
				faulty = 1
			}
			if available && allocated != 1 && faulty != 1 {
				free = 1
			}

			capacities[size] = ServerCapacity{
				Size:      size,
				Total:     total,
				Free:      oldCap.Free + free,
				Allocated: oldCap.Allocated + allocated,
				Faulty:    oldCap.Faulty + faulty,
			}
		}
		var sc []ServerCapacity
		for _, c := range capacities {
			sc = append(sc, c)
		}

		pc := PartitionCapacity{
			Common: &v1.Common{
				Meta: &mdmv1.Meta{
					Id: p.ID,
				},
				Name:        util.StringProto(p.Name),
				Description: util.StringProto(p.Description),
			},
			ServerCapacities: sc,
		}
		partitionCapacities = append(partitionCapacities, pc)
	}
	return partitionCapacities, err
}
