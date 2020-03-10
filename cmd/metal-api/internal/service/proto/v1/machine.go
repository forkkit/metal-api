package v1

import (
	mdv1 "github.com/metal-stack/masterdata-api/api/v1"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service/helper"
	r "gopkg.in/rethinkdb/rethinkdb-go.v5"

	"github.com/metal-stack/metal-api/cmd/metal-api/internal/datastore"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
)

// RecentProvisioningEventsLimit defines how many recent events are added to the MachineRecentProvisioningEvents struct
const RecentProvisioningEventsLimit = 5

func NewMetalMachineHardware(hw *MachineHardwareExtended) metal.MachineHardware {
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

func NewMetalIPMI(ipmi *MachineIPMI) metal.IPMI {
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

func NewMachineIPMIResponse(m *metal.Machine, s *metal.Size, p *metal.Partition, i *metal.Image, ec *metal.ProvisioningEventContainer) *MachineIPMIResponse {
	machineResponse := NewMachineResponse(m, s, p, i, ec)
	return &MachineIPMIResponse{
		Common:  machineResponse.Common,
		Machine: machineResponse.Machine,
		IPMI:    ToMachineIPMI(m.IPMI),
	}
}

func ToMachineIPMI(ipmi metal.IPMI) *MachineIPMI {
	return &MachineIPMI{
		Address:    ipmi.Address,
		MacAddress: ipmi.MacAddress,
		User:       ipmi.User,
		Password:   ipmi.Password,
		Interface:  ipmi.Interface,
		BmcVersion: ipmi.BMCVersion,
		Fru:        ToMachineFRU(ipmi.Fru),
	}
}

func ToMachineFRU(fru metal.Fru) *MachineFru {
	return &MachineFru{
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

func NewMachineResponse(m *metal.Machine, s *metal.Size, p *metal.Partition, img *metal.Image, ec *metal.ProvisioningEventContainer) *MachineResponse {
	return &MachineResponse{
		Common: &Common{
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

func ToMachine(m *metal.Machine, s *metal.Size, p *metal.Partition, img *metal.Image, ec *metal.ProvisioningEventContainer) *Machine {
	var hardware *MachineHardware
	var nics []*MachineNic
	for _, n := range m.Hardware.Nics {
		nic := &MachineNic{
			MacAddress: string(n.MacAddress),
			Name:       n.Name,
		}
		nics = append(nics, nic)

		var disks []*MachineBlockDevice
		for _, d := range m.Hardware.Disks {
			disk := &MachineBlockDevice{
				Name: d.Name,
				Size: d.Size,
			}
			disks = append(disks, disk)
		}
		hardware = &MachineHardware{
			Base: &MachineHardwareBase{
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

	return &Machine{
		Partition:  NewPartitionResponse(p),
		Size:       NewSizeResponse(s),
		Allocation: ToMachineAllocation(m.Allocation, img),
		RackID:     m.RackID,
		Hardware:   hardware,
		BIOS: &MachineBIOS{
			Version: m.BIOS.Version,
			Vendor:  m.BIOS.Vendor,
			Date:    m.BIOS.Date,
		},
		State: &MachineState{
			Value:       string(m.State.Value),
			Description: m.State.Description,
		},
		LedState: &ChassisIdentifyLEDState{
			Value:       string(m.LEDState.Value),
			Description: m.LEDState.Description,
		},
		Liveliness:               liveliness,
		RecentProvisioningEvents: NewMachineRecentProvisioningEvents(ec),
		Tags:                     helper.ToStringValueSlice(m.Tags...),
	}
}

func ToMachineAllocation(alloc *metal.MachineAllocation, img *metal.Image) *MachineAllocation {
	if alloc == nil {
		return nil
	}
	var networks []*MachineNetwork
	for _, nw := range alloc.MachineNetworks {
		ips := append([]string{}, nw.IPs...)
		network := &MachineNetwork{
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

	ma := &MachineAllocation{
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
		ma.Reinstall = &MachineReinstall{
			OldImageID: alloc.ImageID,
			Setup:      ToMachineSetup(alloc),
		}
	}
	return ma
}

func ToMachineSetup(alloc *metal.MachineAllocation) *MachineSetup {
	return &MachineSetup{
		PrimaryDisk:  alloc.PrimaryDisk,
		OsPartition:  alloc.OSPartition,
		Initrd:       alloc.Initrd,
		Cmdline:      alloc.Cmdline,
		Kernel:       alloc.Kernel,
		BootloaderID: alloc.BootloaderID,
	}
}

func NewMachineRecentProvisioningEvents(ec *metal.ProvisioningEventContainer) *MachineRecentProvisioningEvents {
	if ec == nil || ec.LastEventTime == nil {
		return &MachineRecentProvisioningEvents{}
	}
	machineEvents := ec.Events
	if len(machineEvents) >= RecentProvisioningEventsLimit {
		machineEvents = machineEvents[:RecentProvisioningEventsLimit]
	}
	var events []*MachineProvisioningEvent
	for _, machineEvent := range machineEvents {
		e := &MachineProvisioningEvent{
			Time:    helper.ToTimestamp(machineEvent.Time),
			Event:   string(machineEvent.Event),
			Message: helper.ToStringValue(machineEvent.Message),
		}
		events = append(events, e)
	}
	return &MachineRecentProvisioningEvents{
		Events:                       events,
		IncompleteProvisioningCycles: ec.IncompleteProvisioningCycles,
		LastEventTime:                helper.ToTimestamp(*ec.LastEventTime),
	}
}

// GenerateTerm generates the machine search query term
func (x *MachineSearchQuery) GenerateTerm(rs *datastore.RethinkStore) *r.Term {
	q := *rs.MachineTable()

	if x.ID != nil {
		q = q.Filter(func(row r.Term) r.Term {
			return row.Field("id").Eq(*x.ID)
		})
	}

	if x.Name != nil {
		q = q.Filter(func(row r.Term) r.Term {
			return row.Field("name").Eq(*x.Name)
		})
	}

	if x.PartitionID != nil {
		q = q.Filter(func(row r.Term) r.Term {
			return row.Field("partitionid").Eq(*x.PartitionID)
		})
	}

	if x.SizeID != nil {
		q = q.Filter(func(row r.Term) r.Term {
			return row.Field("sizeid").Eq(*x.SizeID)
		})
	}

	if x.RackID != nil {
		q = q.Filter(func(row r.Term) r.Term {
			return row.Field("rackid").Eq(*x.RackID)
		})
	}

	if x.Liveliness != nil {
		q = q.Filter(func(row r.Term) r.Term {
			return row.Field("liveliness").Eq(*x.Liveliness)
		})
	}

	for _, tag := range x.Tags {
		q = q.Filter(func(row r.Term) r.Term {
			return row.Field("tags").Contains(r.Expr(tag))
		})
	}

	if x.AllocationName != nil {
		q = q.Filter(func(row r.Term) r.Term {
			return row.Field("allocation").Field("name").Eq(*x.AllocationName)
		})
	}

	if x.AllocationProject != nil {
		q = q.Filter(func(row r.Term) r.Term {
			return row.Field("allocation").Field("project").Eq(*x.AllocationProject)
		})
	}

	if x.AllocationImageID != nil {
		q = q.Filter(func(row r.Term) r.Term {
			return row.Field("allocation").Field("imageid").Eq(*x.AllocationImageID)
		})
	}

	if x.AllocationHostname != nil {
		q = q.Filter(func(row r.Term) r.Term {
			return row.Field("allocation").Field("hostname").Eq(*x.AllocationHostname)
		})
	}

	if x.AllocationSucceeded != nil {
		q = q.Filter(func(row r.Term) r.Term {
			return row.Field("allocation").Field("succeeded").Eq(*x.AllocationSucceeded)
		})
	}

	for _, id := range x.NetworkIDs {
		q = q.Filter(func(row r.Term) r.Term {
			return row.Field("allocation").Field("networks").Map(func(nw r.Term) r.Term {
				return nw.Field("networkid")
			}).Contains(r.Expr(id))
		})
	}

	for _, prefix := range x.NetworkPrefixes {
		q = q.Filter(func(row r.Term) r.Term {
			return row.Field("allocation").Field("networks").Map(func(nw r.Term) r.Term {
				return nw.Field("prefixes")
			}).Contains(r.Expr(prefix))
		})
	}

	for _, ip := range x.NetworkIPs {
		q = q.Filter(func(row r.Term) r.Term {
			return row.Field("allocation").Field("networks").Map(func(nw r.Term) r.Term {
				return nw.Field("ips")
			}).Contains(r.Expr(ip))
		})
	}

	for _, destPrefix := range x.NetworkDestinationPrefixes {
		q = q.Filter(func(row r.Term) r.Term {
			return row.Field("allocation").Field("networks").Map(func(nw r.Term) r.Term {
				return nw.Field("destinationprefixes")
			}).Contains(r.Expr(destPrefix))
		})
	}

	for _, vrf := range x.NetworkVrfs {
		q = q.Filter(func(row r.Term) r.Term {
			return row.Field("allocation").Field("networks").Map(func(nw r.Term) r.Term {
				return nw.Field("vrf")
			}).Contains(r.Expr(vrf))
		})
	}

	if x.NetworkPrivate != nil {
		q = q.Filter(func(row r.Term) r.Term {
			return row.Field("allocation").Field("networks").Map(func(nw r.Term) r.Term {
				return nw.Field("private")
			}).Contains(*x.NetworkPrivate)
		})
	}

	for _, asn := range x.NetworkASNs {
		q = q.Filter(func(row r.Term) r.Term {
			return row.Field("allocation").Field("networks").Map(func(nw r.Term) r.Term {
				return nw.Field("asn")
			}).Contains(r.Expr(asn))
		})
	}

	if x.NetworkNat != nil {
		q = q.Filter(func(row r.Term) r.Term {
			return row.Field("allocation").Field("networks").Map(func(nw r.Term) r.Term {
				return nw.Field("nat")
			}).Contains(*x.NetworkNat)
		})
	}

	if x.NetworkUnderlay != nil {
		q = q.Filter(func(row r.Term) r.Term {
			return row.Field("allocation").Field("networks").Map(func(nw r.Term) r.Term {
				return nw.Field("underlay")
			}).Contains(*x.NetworkUnderlay)
		})
	}

	if x.HardwareMemory != nil {
		q = q.Filter(func(row r.Term) r.Term {
			return row.Field("hardware").Field("memory").Eq(*x.HardwareMemory)
		})
	}

	if x.HardwareCPUCores != nil {
		q = q.Filter(func(row r.Term) r.Term {
			return row.Field("hardware").Field("cpu_cores").Eq(*x.HardwareCPUCores)
		})
	}

	for _, mac := range x.NicsMacAddresses {
		q = q.Filter(func(row r.Term) r.Term {
			return row.Field("hardware").Field("network_interfaces").Map(func(nic r.Term) r.Term {
				return nic.Field("macAddress")
			}).Contains(r.Expr(mac))
		})
	}

	for _, name := range x.NicsNames {
		q = q.Filter(func(row r.Term) r.Term {
			return row.Field("hardware").Field("network_interfaces").Map(func(nic r.Term) r.Term {
				return nic.Field("name")
			}).Contains(r.Expr(name))
		})
	}

	for _, vrf := range x.NicsVrfs {
		q = q.Filter(func(row r.Term) r.Term {
			return row.Field("hardware").Field("network_interfaces").Map(func(nic r.Term) r.Term {
				return nic.Field("vrf")
			}).Contains(r.Expr(vrf))
		})
	}

	for _, mac := range x.NicsNeighborMacAddresses {
		q = q.Filter(func(row r.Term) r.Term {
			return row.Field("hardware").Field("network_interfaces").Map(func(nic r.Term) r.Term {
				return nic.Field("neighbors").Map(func(neigh r.Term) r.Term {
					return neigh.Field("macAddress")
				})
			}).Contains(r.Expr(mac))
		})
	}

	for _, name := range x.NicsNames {
		q = q.Filter(func(row r.Term) r.Term {
			return row.Field("hardware").Field("network_interfaces").Map(func(nic r.Term) r.Term {
				return nic.Field("neighbors").Map(func(neigh r.Term) r.Term {
					return neigh.Field("name")
				})
			}).Contains(r.Expr(name))
		})
	}

	for _, vrf := range x.NicsVrfs {
		q = q.Filter(func(row r.Term) r.Term {
			return row.Field("hardware").Field("network_interfaces").Map(func(nic r.Term) r.Term {
				return nic.Field("neighbors").Map(func(neigh r.Term) r.Term {
					return neigh.Field("vrf")
				})
			}).Contains(r.Expr(vrf))
		})
	}

	for _, name := range x.DiskNames {
		q = q.Filter(func(row r.Term) r.Term {
			return row.Field("block_devices").Map(func(bd r.Term) r.Term {
				return bd.Field("name")
			}).Contains(r.Expr(name))
		})
	}

	for _, size := range x.DiskSizes {
		q = q.Filter(func(row r.Term) r.Term {
			return row.Field("block_devices").Map(func(bd r.Term) r.Term {
				return bd.Field("size")
			}).Contains(r.Expr(size))
		})
	}

	if x.StateValue != nil {
		q = q.Filter(func(row r.Term) r.Term {
			return row.Field("state_value").Eq(*x.StateValue)
		})
	}

	if x.IpmiAddress != nil {
		q = q.Filter(func(row r.Term) r.Term {
			return row.Field("ipmi").Field("address").Eq(*x.IpmiAddress)
		})
	}

	if x.IpmiMacAddress != nil {
		q = q.Filter(func(row r.Term) r.Term {
			return row.Field("ipmi").Field("mac").Eq(*x.IpmiMacAddress)
		})
	}

	if x.IpmiUser != nil {
		q = q.Filter(func(row r.Term) r.Term {
			return row.Field("ipmi").Field("user").Eq(*x.IpmiUser)
		})
	}

	if x.IpmiInterface != nil {
		q = q.Filter(func(row r.Term) r.Term {
			return row.Field("ipmi").Field("interface").Eq(*x.IpmiInterface)
		})
	}

	if x.FruChassisPartNumber != nil {
		q = q.Filter(func(row r.Term) r.Term {
			return row.Field("ipmi").Field("fru").Field("chassis_part_number").Eq(*x.FruChassisPartNumber)
		})
	}

	if x.FruChassisPartSerial != nil {
		q = q.Filter(func(row r.Term) r.Term {
			return row.Field("ipmi").Field("fru").Field("chassis_part_serial").Eq(*x.FruChassisPartSerial)
		})
	}

	if x.FruBoardMfg != nil {
		q = q.Filter(func(row r.Term) r.Term {
			return row.Field("ipmi").Field("fru").Field("board_mfg").Eq(*x.FruBoardMfg)
		})
	}

	if x.FruBoardMfgSerial != nil {
		q = q.Filter(func(row r.Term) r.Term {
			return row.Field("ipmi").Field("fru").Field("board_mfg_serial").Eq(*x.FruBoardMfgSerial)
		})
	}

	if x.FruBoardPartNumber != nil {
		q = q.Filter(func(row r.Term) r.Term {
			return row.Field("ipmi").Field("fru").Field("board_part_number").Eq(*x.FruBoardPartNumber)
		})
	}

	if x.FruProductManufacturer != nil {
		q = q.Filter(func(row r.Term) r.Term {
			return row.Field("ipmi").Field("fru").Field("product_manufacturer").Eq(*x.FruProductManufacturer)
		})
	}

	if x.FruProductPartNumber != nil {
		q = q.Filter(func(row r.Term) r.Term {
			return row.Field("ipmi").Field("fru").Field("product_part_number").Eq(*x.FruProductPartNumber)
		})
	}

	if x.FruProductSerial != nil {
		q = q.Filter(func(row r.Term) r.Term {
			return row.Field("ipmi").Field("fru").Field("product_serial").Eq(*x.FruProductSerial)
		})
	}

	return &q
}
