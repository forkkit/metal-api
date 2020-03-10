package service

import (
	mdv1 "github.com/metal-stack/masterdata-api/api/v1"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
	"github.com/metal-stack/metal-api/pkg/helper"
	"github.com/metal-stack/metal-api/pkg/proto/v1"
)

// RecentProvisioningEventsLimit defines how many recent events are added to the MachineRecentProvisioningEvents struct
const RecentProvisioningEventsLimit = 5

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
	machineResponse := NewMachineResponse(m, s, p, i, ec)
	return &v1.MachineIPMIResponse{
		Common:  machineResponse.Common,
		Machine: machineResponse.Machine,
		IPMI:    ToMachineIPMI(m.IPMI),
	}
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
		ChassisPartNumber:   helper.ToStringValue(fru.ChassisPartNumber),
		ChassisPartSerial:   helper.ToStringValue(fru.ChassisPartSerial),
		BoardMfg:            helper.ToStringValue(fru.BoardMfg),
		BoardMfgSerial:      helper.ToStringValue(fru.BoardMfgSerial),
		BoardPartNumber:     helper.ToStringValue(fru.BoardPartNumber),
		ProductManufacturer: helper.ToStringValue(fru.ProductManufacturer),
		ProductPartNumber:   helper.ToStringValue(fru.ProductPartNumber),
		ProductSerial:       helper.ToStringValue(fru.ProductSerial),
	}
}

func NewMachineResponse(m *metal.Machine, s *metal.Size, p *metal.Partition, img *metal.Image, ec *metal.ProvisioningEventContainer) *v1.MachineResponse {
	return &v1.MachineResponse{
		Common: &v1.Common{
			Meta: &mdv1.Meta{
				Id:          m.GetID(),
				Apiversion:  "v1",
				Version:     1,
				CreatedTime: helper.ToTimestamp(m.Created),
				UpdatedTime: helper.ToTimestamp(m.Changed),
			},
			Name:        helper.ToStringValue(m.Name),
			Description: helper.ToStringValue(m.Description),
		},
		Machine: ToMachine(m, s, p, img, ec),
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
		Partition:  NewPartitionResponse(p),
		Size:       NewSizeResponse(s),
		Allocation: ToMachineAllocation(m.Allocation, img),
		RackID:     m.RackID,
		Hardware:   hardware,
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
		Tags:                     helper.ToStringValueSlice(m.Tags...),
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
		Created:         helper.ToTimestamp(alloc.Created),
		Name:            alloc.Name,
		Description:     helper.ToStringValue(alloc.Description),
		Image:           NewImageResponse(img),
		ProjectID:       alloc.ProjectID,
		Hostname:        alloc.Hostname,
		UserData:        helper.ToStringValue(alloc.UserData),
		ConsolePassword: helper.ToStringValue(alloc.ConsolePassword),
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
			Time:    helper.ToTimestamp(machineEvent.Time),
			Event:   string(machineEvent.Event),
			Message: helper.ToStringValue(machineEvent.Message),
		}
		events = append(events, e)
	}
	return &v1.MachineRecentProvisioningEvents{
		Events:                       events,
		IncompleteProvisioningCycles: ec.IncompleteProvisioningCycles,
		LastEventTime:                helper.ToTimestamp(*ec.LastEventTime),
	}
}
