package v1

import (
	r "gopkg.in/rethinkdb/rethinkdb-go.v5"
)

// GenerateTerm generates the machine search query term
func (x *MachineSearchQuery) GenerateTerm(q r.Term) *r.Term {
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
