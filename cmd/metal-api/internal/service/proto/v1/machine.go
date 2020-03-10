package v1

import (
	r "gopkg.in/rethinkdb/rethinkdb-go.v5"

	"github.com/metal-stack/metal-api/cmd/metal-api/internal/datastore"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
)

// RecentProvisioningEventsLimit defines how many recent events are added to the MachineRecentProvisioningEvents struct
const RecentProvisioningEventsLimit = 5

func NewMetalMachineHardware(r *MachineHardwareExtended) metal.MachineHardware {
	var nics metal.Nics
	for i := range r.Nics {
		var neighbors metal.Nics
		for i2 := range r.Nics[i].Neighbors {
			neighbor := metal.Nic{
				MacAddress: metal.MacAddress(r.Nics[i].Neighbors[i2].MacAddress),
				Name:       r.Nics[i].Neighbors[i2].Name,
			}
			neighbors = append(neighbors, neighbor)
		}
		nic := metal.Nic{
			MacAddress: metal.MacAddress(r.Nics[i].MacAddress),
			Name:       r.Nics[i].Name,
			Neighbors:  neighbors,
		}
		nics = append(nics, nic)
	}
	var disks []metal.BlockDevice
	for _, d := range r.Disks {
		disk := metal.BlockDevice{
			Name:    d.Name,
			Size:    d.Size,
			Primary: d.Primary,
		}
		for _, p := range d.Partitions {
			disk.Partitions = append(disk.Partitions, &metal.DiskPartition{
				Label:        p.Label,
				Device:       p.Device,
				Number:       p.Number,
				MountPoint:   p.MountPoint,
				MountOptions: p.MountOptions,
				Size:         p.Size,
				Filesystem:   p.Filesystem,
				GPTType:      p.GPTType,
				GPTGuid:      p.GPTGuid,
				Properties:   p.Properties,
				ContainsOS:   p.ContainsOS,
			})
		}
		disks = append(disks, disk)
	}
	return metal.MachineHardware{
		Memory:   r.Memory,
		CPUCores: r.CPUCores,
		Nics:     nics,
		Disks:    disks,
	}
}

func NewMetalIPMI(r *MachineIPMI) metal.IPMI {
	var chassisPartNumber string
	if r.Fru.ChassisPartNumber != nil {
		chassisPartNumber = *r.Fru.ChassisPartNumber
	}
	var chassisPartSerial string
	if r.Fru.ChassisPartSerial != nil {
		chassisPartSerial = *r.Fru.ChassisPartSerial
	}
	var boardMfg string
	if r.Fru.BoardMfg != nil {
		boardMfg = *r.Fru.BoardMfg
	}
	var boardMfgSerial string
	if r.Fru.BoardMfgSerial != nil {
		boardMfgSerial = *r.Fru.BoardMfgSerial
	}
	var boardPartNumber string
	if r.Fru.BoardPartNumber != nil {
		boardPartNumber = *r.Fru.BoardPartNumber
	}
	var productManufacturer string
	if r.Fru.ProductManufacturer != nil {
		productManufacturer = *r.Fru.ProductManufacturer
	}
	var productPartNumber string
	if r.Fru.ProductPartNumber != nil {
		productPartNumber = *r.Fru.ProductPartNumber
	}
	var productSerial string
	if r.Fru.ProductSerial != nil {
		productSerial = *r.Fru.ProductSerial
	}

	return metal.IPMI{
		Address:    r.Address,
		MacAddress: r.MacAddress,
		User:       r.User,
		Password:   r.Password,
		Interface:  r.Interface,
		BMCVersion: r.BMCVersion,
		Fru: metal.Fru{
			ChassisPartNumber:   chassisPartNumber,
			ChassisPartSerial:   chassisPartSerial,
			BoardMfg:            boardMfg,
			BoardMfgSerial:      boardMfgSerial,
			BoardPartNumber:     boardPartNumber,
			ProductManufacturer: productManufacturer,
			ProductPartNumber:   productPartNumber,
			ProductSerial:       productSerial,
		},
	}
}

func NewMachineIPMIResponse(m *metal.Machine, s *metal.Size, p *metal.Partition, i *metal.Image, ec *metal.ProvisioningEventContainer) *MachineIPMIResponse {
	machineResponse := NewMachineResponse(m, s, p, i, ec)
	return &MachineIPMIResponse{
		Common:      machineResponse.Common,
		MachineBase: machineResponse.MachineBase,
		IPMI: MachineIPMI{
			Address:    m.IPMI.Address,
			MacAddress: m.IPMI.MacAddress,
			User:       m.IPMI.User,
			Password:   m.IPMI.Password,
			Interface:  m.IPMI.Interface,
			BMCVersion: m.IPMI.BMCVersion,
			Fru: MachineFru{
				ChassisPartNumber:   &m.IPMI.Fru.ChassisPartNumber,
				ChassisPartSerial:   &m.IPMI.Fru.ChassisPartSerial,
				BoardMfg:            &m.IPMI.Fru.BoardMfg,
				BoardMfgSerial:      &m.IPMI.Fru.BoardMfgSerial,
				BoardPartNumber:     &m.IPMI.Fru.BoardPartNumber,
				ProductManufacturer: &m.IPMI.Fru.ProductManufacturer,
				ProductPartNumber:   &m.IPMI.Fru.ProductPartNumber,
				ProductSerial:       &m.IPMI.Fru.ProductSerial,
			},
		},
		Timestamps: machineResponse.Timestamps,
	}
}

func NewMachineResponse(m *metal.Machine, s *metal.Size, p *metal.Partition, i *metal.Image, ec *metal.ProvisioningEventContainer) *MachineResponse {
	var hardware MachineHardware
	var nics MachineNics
	for i := range m.Hardware.Nics {
		nic := MachineNic{
			MacAddress: string(m.Hardware.Nics[i].MacAddress),
			Name:       m.Hardware.Nics[i].Name,
		}
		nics = append(nics, nic)

		var disks []MachineBlockDevice
		for i := range m.Hardware.Disks {
			disk := MachineBlockDevice{
				Name: m.Hardware.Disks[i].Name,
				Size: m.Hardware.Disks[i].Size,
			}
			disks = append(disks, disk)
		}
		hardware = MachineHardware{
			MachineHardwareBase: MachineHardwareBase{
				Memory:   m.Hardware.Memory,
				CPUCores: m.Hardware.CPUCores,
				Disks:    disks,
			},
			Nics: nics,
		}
	}

	var allocation *MachineAllocation
	if m.Allocation != nil {
		var networks []MachineNetwork
		for _, nw := range m.Allocation.MachineNetworks {
			ips := append([]string{}, nw.IPs...)
			network := MachineNetwork{
				NetworkID:           nw.NetworkID,
				IPs:                 ips,
				Vrf:                 nw.Vrf,
				ASN:                 nw.ASN,
				Private:             nw.Private,
				Nat:                 nw.Nat,
				Underlay:            nw.Underlay,
				DestinationPrefixes: nw.DestinationPrefixes,
				Prefixes:            nw.Prefixes,
			}
			networks = append(networks, network)
		}

		var consolePassword *string
		if m.Allocation.ConsolePassword != "" {
			consolePassword = &m.Allocation.ConsolePassword
		}

		allocation = &MachineAllocation{
			Created:         m.Allocation.Created,
			Name:            m.Allocation.Name,
			Description:     m.Allocation.Description,
			Image:           NewImageResponse(i),
			Project:         m.Allocation.Project,
			Hostname:        m.Allocation.Hostname,
			SSHPubKeys:      m.Allocation.SSHPubKeys,
			UserData:        m.Allocation.UserData,
			ConsolePassword: consolePassword,
			MachineNetworks: networks,
			Succeeded:       m.Allocation.Succeeded,
		}

		if m.Allocation.Reinstall {
			allocation.Reinstall = &MachineReinstall{
				OldImageID:   m.Allocation.ImageID,
				PrimaryDisk:  m.Allocation.PrimaryDisk,
				OSPartition:  m.Allocation.OSPartition,
				Initrd:       m.Allocation.Initrd,
				Cmdline:      m.Allocation.Cmdline,
				Kernel:       m.Allocation.Kernel,
				BootloaderID: m.Allocation.BootloaderID,
			}
		}
	}

	var tags []string
	if len(m.Tags) > 0 {
		tags = m.Tags
	}
	liveliness := ""
	if ec != nil {
		liveliness = string(ec.Liveliness)
	}

	return &MachineResponse{
		Common: Common{
			Identifiable: Identifiable{
				ID: m.ID,
			},
			Describable: Describable{
				Name:        &m.Name,
				Description: &m.Description,
			},
		},
		MachineBase: MachineBase{
			Partition:  NewPartitionResponse(p),
			Size:       NewSizeResponse(s),
			Allocation: allocation,
			RackID:     m.RackID,
			Hardware:   hardware,
			BIOS: MachineBIOS{
				Version: m.BIOS.Version,
				Vendor:  m.BIOS.Vendor,
				Date:    m.BIOS.Date,
			},
			State: MachineState{
				Value:       string(m.State.Value),
				Description: m.State.Description,
			},
			LEDState: ChassisIdentifyLEDState{
				Value:       string(m.LEDState.Value),
				Description: m.LEDState.Description,
			},
			Liveliness:               liveliness,
			RecentProvisioningEvents: *NewMachineRecentProvisioningEvents(ec),
			Tags:                     tags,
		},
		Timestamps: Timestamps{
			Created: m.Created,
			Changed: m.Changed,
		},
	}
}

func NewMachineRecentProvisioningEvents(ec *metal.ProvisioningEventContainer) *MachineRecentProvisioningEvents {
	var es []MachineProvisioningEvent
	if ec == nil {
		return &MachineRecentProvisioningEvents{
			Events:                       es,
			LastEventTime:                nil,
			IncompleteProvisioningCycles: "0",
		}
	}
	machineEvents := ec.Events
	if len(machineEvents) >= RecentProvisioningEventsLimit {
		machineEvents = machineEvents[:RecentProvisioningEventsLimit]
	}
	for _, machineEvent := range machineEvents {
		e := MachineProvisioningEvent{
			Time:    machineEvent.Time,
			Event:   string(machineEvent.Event),
			Message: machineEvent.Message,
		}
		es = append(es, e)
	}
	return &MachineRecentProvisioningEvents{
		Events:                       es,
		IncompleteProvisioningCycles: ec.IncompleteProvisioningCycles,
		LastEventTime:                ec.LastEventTime,
	}
}

// GenerateTerm generates the project search query term.
func (x *MachineSearchQuery) GenerateTerm(rs *datastore.RethinkStore) *r.Term {
	q := *rs.machineTable()

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
