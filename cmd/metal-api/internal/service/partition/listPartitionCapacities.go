package partition

import (
	"github.com/emicklei/go-restful"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service/helper"
	v1 "github.com/metal-stack/metal-api/cmd/metal-api/internal/service/v1"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/utils"
	"github.com/metal-stack/metal-lib/zapup"
	"go.uber.org/zap"
	"net/http"
)

func (r partitionResource) listPartitionCapacities(request *restful.Request, response *restful.Response) {
	partitionCapacities, err := r.calcPartitionCapacities()

	if helper.CheckError(request, response, utils.CurrentFuncName(), err) {
		return
	}
	err = response.WriteHeaderAndEntity(http.StatusOK, partitionCapacities)
	if err != nil {
		zapup.MustRootLogger().Error("Failed to send response", zap.Error(err))
		return
	}
}

func (r partitionResource) calcPartitionCapacities() ([]v1.PartitionCapacity, error) {
	// FIXME bad workaround to be able to run make spec
	if r.DS == nil {
		return nil, nil
	}
	ps, err := r.DS.ListPartitions()
	if err != nil {
		return nil, err
	}
	ms, err := r.DS.ListMachines()
	if err != nil {
		return nil, err
	}
	machines := helper.MakeMachineResponseList(ms, r.DS, zapup.MustRootLogger().Sugar())

	var partitionCapacities []v1.PartitionCapacity
	for _, p := range ps {
		capacities := make(map[string]v1.ServerCapacity)
		for _, m := range machines {
			if m.Partition == nil {
				continue
			}
			if m.Partition.ID != p.ID {
				continue
			}
			size := "unknown"
			if m.Size != nil {
				size = m.Size.ID
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
			if helper.MachineHasIssues(m) {
				faulty = 1
			}
			if available && allocated != 1 && faulty != 1 {
				free = 1
			}

			cap := v1.ServerCapacity{
				Size:      size,
				Total:     total,
				Free:      oldCap.Free + free,
				Allocated: oldCap.Allocated + allocated,
				Faulty:    oldCap.Faulty + faulty,
			}
			capacities[size] = cap
		}
		var sc []v1.ServerCapacity
		for _, c := range capacities {
			sc = append(sc, c)
		}

		pc := v1.PartitionCapacity{
			Common: v1.Common{
				Identifiable: v1.Identifiable{
					ID: p.ID,
				},
				Describable: v1.Describable{
					Name:        &p.Name,
					Description: &p.Description,
				},
			},
			ServerCapacities: sc,
		}
		partitionCapacities = append(partitionCapacities, pc)
	}
	return partitionCapacities, err
}
