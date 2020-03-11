package machine

import (
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service"
	v1 "github.com/metal-stack/metal-api/pkg/proto/v1"
	"go.uber.org/zap"
	"time"

	mdm "github.com/metal-stack/masterdata-api/pkg/client"

	"github.com/emicklei/go-restful"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/datastore"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/ipam"
	"github.com/metal-stack/metal-lib/bus"
)

const (
	waitForServerTimeout = 30 * time.Second
)

type machineResource struct {
	bus.Publisher
	ds     *datastore.RethinkStore
	ipamer ipam.IPAMer
	mdc    mdm.Client
}

// NewMachine returns a webservice for machine specific endpoints.
func NewMachine(
	ds *datastore.RethinkStore,
	pub bus.Publisher,
	ipamer ipam.IPAMer,
	mdc mdm.Client) *restful.WebService {
	r := machineResource{
		ds:        ds,
		Publisher: pub,
		ipamer:    ipamer,
		mdc:       mdc,
	}
	return r.webService()
}

func PublishMachineCmd(logger *zap.SugaredLogger, m *metal.Machine, publisher bus.Publisher, cmd metal.MachineCommand, params ...string) error {
	var pp []string
	for _, p := range params {
		if len(p) > 0 {
			pp = append(pp, p)
		}
	}
	evt := metal.MachineEvent{
		Type: metal.COMMAND,
		Cmd: &metal.MachineExecCommand{
			Command: cmd,
			Params:  pp,
			Target:  m,
		},
	}

	logger.Infow("publish event", "event", evt, "command", *evt.Cmd)
	err := publisher.Publish(metal.TopicMachine.GetFQN(m.PartitionID), evt)
	if err != nil {
		return err
	}

	return nil
}

func ResurrectMachines(ds *datastore.RethinkStore, publisher bus.Publisher, logger *zap.SugaredLogger) error {
	logger.Info("machine resurrection was requested")

	machines, err := ds.ListMachines()
	if err != nil {
		return err
	}

	for _, m := range machines {
		if m.Allocation != nil {
			continue
		}

		provisioningEvents, err := ds.FindProvisioningEventContainer(m.ID)
		if err != nil {
			// we have no provisioning events... we cannot tell
			logger.Debugw("no provisioningEvents found for resurrection", "machineID", m.ID, "error", err)
			continue
		}

		if provisioningEvents.Liveliness != metal.MachineLivelinessDead {
			continue
		}

		if provisioningEvents.LastEventTime == nil {
			continue
		}

		if time.Since(*provisioningEvents.LastEventTime) < metal.MachineResurrectAfter {
			continue
		}

		logger.Infow("resurrecting dead machine", "machineID", m.ID, "liveliness", provisioningEvents.Liveliness, "since", time.Since(*provisioningEvents.LastEventTime).String())
		err = PublishMachineCmd(logger, &m, publisher, metal.MachineResetCmd)
		if err != nil {
			logger.Errorw("error during machine resurrection when trying to publish machine reset cmd", "machineID", m.ID, "error", err)
		}
	}

	return nil
}

func MachineHasIssues(m *v1.MachineResponse) bool {
	if m.Partition == nil {
		return true
	}
	if !metal.MachineLivelinessAlive.Is(m.Liveliness) {
		return true
	}
	if m.Allocation == nil && len(m.RecentProvisioningEvents.Events) > 0 && metal.ProvisioningEventPhonedHome.Is(m.RecentProvisioningEvents.Events[0].Event) {
		// not allocated, but phones home
		return true
	}
	if m.RecentProvisioningEvents.IncompleteProvisioningCycles != "" && m.RecentProvisioningEvents.IncompleteProvisioningCycles != "0" {
		// Machines with incomplete cycles but in "Waiting" state are considered available
		if len(m.RecentProvisioningEvents.Events) > 0 && !metal.ProvisioningEventWaiting.Is(m.RecentProvisioningEvents.Events[0].Event) {
			return true
		}
	}

	return false
}

func MakeMachineResponse(m *metal.Machine, ds *datastore.RethinkStore, logger *zap.SugaredLogger) *v1.MachineResponse {
	s, p, i, ec := FindMachineReferencedEntities(m, ds, logger)
	return service.NewMachineResponse(m, s, p, i, ec)
}

func MakeMachineResponseList(ms metal.Machines, ds *datastore.RethinkStore, logger *zap.SugaredLogger) []*v1.MachineResponse {
	sMap, pMap, iMap, ecMap := GetMachineReferencedEntityMaps(ds, logger)

	var result []*v1.MachineResponse

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
		result = append(result, service.NewMachineResponse(&ms[index], s, p, i, &ec))
	}

	return result
}

func FindMachineReferencedEntities(m *metal.Machine, ds *datastore.RethinkStore, logger *zap.SugaredLogger) (*metal.Size, *metal.Partition, *metal.Image, *metal.ProvisioningEventContainer) {
	var err error

	var s *metal.Size
	if m.SizeID != "" {
		if m.SizeID == metal.UnknownSize.GetID() {
			s = metal.UnknownSize
		} else {
			s, err = ds.FindSize(m.SizeID)
			if err != nil {
				logger.Errorw("machine references size, but size cannot be found in database", "machineID", m.ID, "sizeID", m.SizeID, "error", err)
			}
		}
	}

	var p *metal.Partition
	if m.PartitionID != "" {
		p, err = ds.FindPartition(m.PartitionID)
		if err != nil {
			logger.Errorw("machine references partition, but partition cannot be found in database", "machineID", m.ID, "partitionID", m.PartitionID, "error", err)
		}
	}

	var i *metal.Image
	if m.Allocation != nil {
		if m.Allocation.ImageID != "" {
			i, err = ds.FindImage(m.Allocation.ImageID)
			if err != nil {
				logger.Errorw("machine references image, but image cannot be found in database", "machineID", m.ID, "imageID", m.Allocation.ImageID, "error", err)
			}
		}
	}

	var ec *metal.ProvisioningEventContainer
	try, err := ds.FindProvisioningEventContainer(m.ID)
	if err != nil {
		logger.Errorw("machine has no provisioning event container in the database", "machineID", m.ID, "error", err)
	} else {
		ec = try
	}

	return s, p, i, ec
}

func GetMachineReferencedEntityMaps(ds *datastore.RethinkStore, logger *zap.SugaredLogger) (metal.SizeMap, metal.PartitionMap, metal.ImageMap, metal.ProvisioningEventContainerMap) {
	s, err := ds.ListSizes()
	if err != nil {
		logger.Errorw("sizes could not be listed", "error", err)
	}

	p, err := ds.ListPartitions()
	if err != nil {
		logger.Errorw("partitions could not be listed", "error", err)
	}

	i, err := ds.ListImages()
	if err != nil {
		logger.Errorw("images could not be listed", "error", err)
	}

	ec, err := ds.ListProvisioningEventContainers()
	if err != nil {
		logger.Errorw("provisioning event containers could not be listed", "error", err)
	}

	return s.ByID(), p.ByID(), i.ByID(), ec.ByID()
}
