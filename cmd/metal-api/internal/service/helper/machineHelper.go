package helper

import (
	v12 "github.com/metal-stack/masterdata-api/api/v1"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/datastore"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
	"github.com/metal-stack/metal-api/pkg/proto/v1"
	"github.com/metal-stack/metal-api/pkg/util"
	"go.uber.org/zap"
)

// RecentProvisioningEventsLimit defines how many recent events are added to the MachineRecentProvisioningEvents struct
const RecentProvisioningEventsLimit = 5

func NewMachineResponse(m *metal.Machine, s *metal.Size, p *metal.Partition, img *metal.Image, ec *metal.ProvisioningEventContainer) *v1.MachineResponse {
	c := &v1.Common{
		Meta: &v12.Meta{
			Id:          m.GetID(),
			Apiversion:  "v1",
			Version:     1,
			CreatedTime: util.TimestampProto(m.Created),
			UpdatedTime: util.TimestampProto(m.Changed),
		},
		Name:        util.StringProto(m.Name),
		Description: util.StringProto(m.Description),
	}
	machine := ToMachine(m, s, p, img, ec)
	machine.Common = c
	return &v1.MachineResponse{
		Common:  c,
		Machine: machine,
	}
}

func GetReferencedEntityMaps(ds *datastore.RethinkStore, logger *zap.SugaredLogger) (metal.SizeMap, metal.PartitionMap, metal.ImageMap, metal.ProvisioningEventContainerMap) {
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

func ToMachineIPMI(ipmi metal.IPMI) *v1.MachineIPMI {
	return &v1.MachineIPMI{
		Address:    ipmi.Address,
		MacAddress: ipmi.MacAddress,
		User:       ipmi.User,
		Password:   ipmi.Password,
		Interface:  ipmi.Interface,
		BmcVersion: ipmi.BMCVersion,
		Fru:        ToMachineFRU(ipmi.Fru),
	}
}

func ToMachineFRU(fru metal.Fru) *v1.MachineFru {
	return &v1.MachineFru{
		ChassisPartNumber:   util.StringProto(fru.ChassisPartNumber),
		ChassisPartSerial:   util.StringProto(fru.ChassisPartSerial),
		BoardMfg:            util.StringProto(fru.BoardMfg),
		BoardMfgSerial:      util.StringProto(fru.BoardMfgSerial),
		BoardPartNumber:     util.StringProto(fru.BoardPartNumber),
		ProductManufacturer: util.StringProto(fru.ProductManufacturer),
		ProductPartNumber:   util.StringProto(fru.ProductPartNumber),
		ProductSerial:       util.StringProto(fru.ProductSerial),
	}
}

func ToMachine(m *metal.Machine, s *metal.Size, p *metal.Partition, img *metal.Image, ec *metal.ProvisioningEventContainer) *v1.Machine {
	var hardware *v1.MachineHardware
	var nics []*v1.MachineNic
	for _, n := range m.Hardware.Nics {
		nic := &v1.MachineNic{
			MacAddress: string(n.MacAddress),
			Name:       n.Name,
		}
		nics = append(nics, nic)

		var disks []*v1.MachineBlockDevice
		for _, d := range m.Hardware.Disks {
			disk := &v1.MachineBlockDevice{
				Name: d.Name,
				Size: d.Size,
			}
			disks = append(disks, disk)
		}
		hardware = &v1.MachineHardware{
			Base: &v1.MachineHardwareBase{
				Memory:   m.Hardware.Memory,
				CpuCores: uint32(m.Hardware.CPUCores),
				Disks:    disks,
			},
			Nics: nics,
		}
	}

	liveliness := ""
	if ec != nil {
		liveliness = string(ec.Liveliness)
	}

	return &v1.Machine{
		PartitionResponse: NewPartitionResponse(p),
		SizeResponse:      NewSizeResponse(s),
		Allocation:        ToMachineAllocation(m.Allocation, img),
		RackID:            m.RackID,
		Hardware:          hardware,
		BIOS: &v1.MachineBIOS{
			Version: m.BIOS.Version,
			Vendor:  m.BIOS.Vendor,
			Date:    m.BIOS.Date,
		},
		State: &v1.MachineState{
			Value:       string(m.State.Value),
			Description: m.State.Description,
		},
		LedState: &v1.ChassisIdentifyLEDState{
			Value:       string(m.LEDState.Value),
			Description: m.LEDState.Description,
		},
		Liveliness:               liveliness,
		RecentProvisioningEvents: NewMachineRecentProvisioningEvents(ec),
		Tags:                     util.StringSliceProto(m.Tags...),
	}
}

func ToMachineAllocation(alloc *metal.MachineAllocation, img *metal.Image) *v1.MachineAllocation {
	if alloc == nil {
		return nil
	}
	var networks []*v1.MachineNetwork
	for _, nw := range alloc.MachineNetworks {
		ips := append([]string{}, nw.IPs...)
		network := &v1.MachineNetwork{
			NetworkID:           nw.NetworkID,
			IPs:                 ips,
			Vrf:                 uint64(nw.Vrf),
			ASN:                 nw.ASN,
			Private:             nw.Private,
			Nat:                 nw.Nat,
			Underlay:            nw.Underlay,
			DestinationPrefixes: nw.DestinationPrefixes,
			Prefixes:            nw.Prefixes,
		}
		networks = append(networks, network)
	}

	ma := &v1.MachineAllocation{
		Created:         util.TimestampProto(alloc.Created),
		Name:            alloc.Name,
		Description:     util.StringProto(alloc.Description),
		ImageResponse:   NewImageResponse(img),
		ProjectID:       alloc.ProjectID,
		Hostname:        alloc.Hostname,
		UserData:        util.StringProto(alloc.UserData),
		ConsolePassword: util.StringProto(alloc.ConsolePassword),
		MachineNetworks: networks,
		Succeeded:       alloc.Succeeded,
		SshPubKeys:      alloc.SSHPubKeys,
	}
	if alloc.Reinstall {
		ma.Reinstall = &v1.MachineReinstall{
			OldImageID: alloc.ImageID,
			Setup:      ToMachineSetup(alloc),
		}
	}
	return ma
}

func ToMachineSetup(alloc *metal.MachineAllocation) *v1.MachineSetup {
	return &v1.MachineSetup{
		PrimaryDisk:  alloc.PrimaryDisk,
		OsPartition:  alloc.OSPartition,
		Initrd:       alloc.Initrd,
		Cmdline:      alloc.Cmdline,
		Kernel:       alloc.Kernel,
		BootloaderID: alloc.BootloaderID,
	}
}

func MakeMachineResponseList(ms metal.Machines, ds *datastore.RethinkStore, logger *zap.SugaredLogger) []*v1.MachineResponse {
	sMap, pMap, iMap, ecMap := GetReferencedEntityMaps(ds, logger)

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
		result = append(result, NewMachineResponse(&ms[index], s, p, i, &ec))
	}

	return result
}

func HasMachineIssues(mr *v1.MachineResponse) bool {
	if mr.Machine == nil || mr.Machine.PartitionResponse == nil || mr.Machine.PartitionResponse.Partition == nil {
		return true
	}
	m := mr.Machine
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

func NewMachineRecentProvisioningEvents(ec *metal.ProvisioningEventContainer) *v1.MachineRecentProvisioningEvents {
	if ec == nil || ec.LastEventTime == nil {
		return &v1.MachineRecentProvisioningEvents{}
	}
	machineEvents := ec.Events
	if len(machineEvents) >= RecentProvisioningEventsLimit {
		machineEvents = machineEvents[:RecentProvisioningEventsLimit]
	}
	var events []*v1.MachineProvisioningEvent
	for _, machineEvent := range machineEvents {
		e := &v1.MachineProvisioningEvent{
			Time:    util.TimestampProto(machineEvent.Time),
			Event:   string(machineEvent.Event),
			Message: util.StringProto(machineEvent.Message),
		}
		events = append(events, e)
	}
	return &v1.MachineRecentProvisioningEvents{
		Events:                       events,
		IncompleteProvisioningCycles: ec.IncompleteProvisioningCycles,
		LastEventTime:                util.TimestampProto(*ec.LastEventTime),
	}
}

func NewMetalMachineHardware(hw *v1.MachineHardwareExtended) metal.MachineHardware {
	var nics metal.Nics
	for _, n := range hw.Nics {
		var neighbors metal.Nics
		for _, neigh := range n.Neighbors {
			neighbor := metal.Nic{
				MacAddress: metal.MacAddress(neigh.MachineNic.MacAddress),
				Name:       neigh.MachineNic.Name,
			}
			neighbors = append(neighbors, neighbor)
		}
		nic := metal.Nic{
			MacAddress: metal.MacAddress(n.MachineNic.MacAddress),
			Name:       n.MachineNic.Name,
			Neighbors:  neighbors,
		}
		nics = append(nics, nic)
	}
	var disks []metal.BlockDevice
	for _, d := range hw.Base.Disks {
		disk := metal.BlockDevice{
			Name:    d.Name,
			Size:    d.Size,
			Primary: d.Primary,
		}
		for _, p := range d.Partitions {
			disk.Partitions = append(disk.Partitions, &metal.DiskPartition{
				Label:        p.Label,
				Device:       p.Device,
				Number:       uint(p.Number),
				MountPoint:   p.MountPoint,
				MountOptions: p.MountOptions,
				Size:         p.Size,
				Filesystem:   p.Filesystem,
				GPTType:      p.GptType,
				GPTGuid:      p.GptGuid,
				Properties:   p.Properties,
				ContainsOS:   p.ContainsOS,
			})
		}
		disks = append(disks, disk)
	}
	return metal.MachineHardware{
		Memory:   hw.Base.Memory,
		CPUCores: uint(hw.Base.CpuCores),
		Nics:     nics,
		Disks:    disks,
	}
}
