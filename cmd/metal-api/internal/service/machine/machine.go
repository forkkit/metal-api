package machine

import (
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service/helper"
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

// NewMachineService returns a webservice for machine specific endpoints.
func NewMachineService(
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
	// TODO Find better place
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

func MakeResponse(m *metal.Machine, ds *datastore.RethinkStore, logger *zap.SugaredLogger) *v1.MachineResponse {
	s, p, i, ec := FindMachineReferencedEntities(m, ds, logger)
	return helper.NewMachineResponse(m, s, p, i, ec)
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

func NewMetalIPMI(ipmi *v1.MachineIPMI) metal.IPMI {
	return metal.IPMI{
		Address:    ipmi.Address,
		MacAddress: ipmi.MacAddress,
		User:       ipmi.User,
		Password:   ipmi.Password,
		Interface:  ipmi.Interface,
		BMCVersion: ipmi.BmcVersion,
		Fru: metal.Fru{
			ChassisPartNumber:   ipmi.Fru.ChassisPartNumber.GetValue(),
			ChassisPartSerial:   ipmi.Fru.ChassisPartSerial.GetValue(),
			BoardMfg:            ipmi.Fru.BoardMfg.GetValue(),
			BoardMfgSerial:      ipmi.Fru.BoardMfgSerial.GetValue(),
			BoardPartNumber:     ipmi.Fru.BoardPartNumber.GetValue(),
			ProductManufacturer: ipmi.Fru.ProductManufacturer.GetValue(),
			ProductPartNumber:   ipmi.Fru.ProductPartNumber.GetValue(),
			ProductSerial:       ipmi.Fru.ProductSerial.GetValue(),
		},
	}
}

func NewMachineIPMIResponse(m *metal.Machine, s *metal.Size, p *metal.Partition, i *metal.Image, ec *metal.ProvisioningEventContainer) *v1.MachineIPMIResponse {
	machineResponse := helper.NewMachineResponse(m, s, p, i, ec)
	return &v1.MachineIPMIResponse{
		Common:  machineResponse.Common,
		Machine: machineResponse.Machine,
		IPMI:    helper.ToMachineIPMI(m.IPMI),
	}
}
