package helper

import (
	"fmt"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/datastore"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
	v12 "github.com/metal-stack/metal-api/cmd/metal-api/internal/service/proto/v1"
	"go.uber.org/zap"
)

func UpdateSwitchNics(oldNics metal.NicMap, newNics metal.NicMap, currentConnections metal.ConnectionMap) (metal.Nics, error) {
	// TODO: Broken switch would change nics, but if this happens we would need to repair broken connections:
	// metal/metal#28

	// To start off we just prevent basic things that can go wrong
	nicsThatGetLost := metal.Nics{}
	for mac, nic := range oldNics {
		_, ok := newNics[mac]
		if !ok {
			nicsThatGetLost = append(nicsThatGetLost, *nic)
		}
	}

	// check if nic gets removed but has a connection
	for _, nicThatGetsLost := range nicsThatGetLost {
		for machineID, connections := range currentConnections {
			for _, c := range connections {
				if c.Nic.MacAddress == nicThatGetsLost.MacAddress {
					return nil, fmt.Errorf("nic with mac address %s gets removed but the machine with id %q is already connected to this nic, which is currently not supported", nicThatGetsLost.MacAddress, machineID)
				}
			}
		}
	}

	nicsThatGetAdded := metal.Nics{}
	nicsThatAlreadyExist := metal.Nics{}
	for mac, nic := range newNics {
		oldNic, ok := oldNics[mac]
		if ok {
			updatedNic := *oldNic

			// check if connection exists and name changes
			for machineID, connections := range currentConnections {
				for _, c := range connections {
					if c.Nic.MacAddress == nic.MacAddress && oldNic.Name != nic.Name {
						return nil, fmt.Errorf("nic with mac address %s wants to be renamed from %q to %q, but already has a connection to machine with id %q, which is currently not supported", nic.MacAddress, oldNic.Name, nic.Name, machineID)
					}
				}
			}

			updatedNic.Name = nic.Name
			nicsThatAlreadyExist = append(nicsThatAlreadyExist, updatedNic)
		} else {
			nicsThatGetAdded = append(nicsThatGetAdded, *nic)
		}
	}

	finalNics := metal.Nics{}
	finalNics = append(finalNics, nicsThatGetAdded...)
	finalNics = append(finalNics, nicsThatAlreadyExist...)

	return finalNics, nil

}

// SetVrfAtSwitches finds the switches connected to the given machine and puts the switch ports into the given vrf.
// Returns the updated switches.
func SetVrfAtSwitches(ds *datastore.RethinkStore, m *metal.Machine, vrf string) ([]metal.Switch, error) {
	switches, err := ds.SearchSwitchesConnectedToMachine(m)
	if err != nil {
		return nil, err
	}
	newSwitches := make([]metal.Switch, 0)
	for _, sw := range switches {
		oldSwitch := sw
		setVrf(&sw, m.ID, vrf)
		err := ds.UpdateSwitch(&oldSwitch, &sw)
		if err != nil {
			return nil, err
		}
		newSwitches = append(newSwitches, sw)
	}
	return newSwitches, nil
}

func setVrf(s *metal.Switch, mid, vrf string) {
	// gather nics within MachineConnections
	changed := metal.Nics{}
	for _, c := range s.MachineConnections[mid] {
		c.Nic.Vrf = vrf
		changed = append(changed, c.Nic)
	}

	if len(changed) == 0 {
		return
	}

	// update sw.Nics
	currentByMac := s.Nics.ByMac()
	changedByMac := changed.ByMac()
	s.Nics = metal.Nics{}
	for mac, old := range currentByMac {
		e := old
		if new, has := changedByMac[mac]; has {
			e = new
		}
		s.Nics = append(s.Nics, *e)
	}
}

func ConnectMachineWithSwitches(ds *datastore.RethinkStore, m *metal.Machine) error {
	switches, err := ds.SearchSwitches(m.RackID, nil)
	if err != nil {
		return err
	}
	for _, sw := range switches {
		oldSwitch := sw
		sw.ConnectMachine(m)
		err := ds.UpdateSwitch(&oldSwitch, &sw)
		if err != nil {
			return err
		}
	}
	return nil
}

func MakeSwitchResponse(s *metal.Switch, ds *datastore.RethinkStore, logger *zap.SugaredLogger) *v12.SwitchResponse {
	p, ips, iMap, machines := findSwitchReferencedEntites(s, ds, logger)
	nics := MakeSwitchNics(s, ips, iMap, machines)
	cons := makeSwitchCons(s)
	return v12.NewSwitchResponse(s, p, nics, cons)
}

func MakeBGPFilterFirewall(m metal.Machine) v12.BGPFilter {
	var vnis, cidrs []string
	for _, net := range m.Allocation.MachineNetworks {
		if net.Underlay {
			for _, ip := range net.IPs {
				cidrs = append(cidrs, fmt.Sprintf("%s/32", ip))
			}
		} else {
			vnis = append(vnis, fmt.Sprintf("%d", net.Vrf))
			// filter for "project" addresses / cidrs is not possible since EVPN Type-5 routes can not be filtered by prefixes
		}
	}
	return v12.NewBGPFilter(vnis, cidrs)
}

func MakeBGPFilterMachine(m metal.Machine, ips metal.IPsMap) v12.BGPFilter {
	var vnis, cidrs []string

	var private *metal.MachineNetwork
	var underlay *metal.MachineNetwork
	for _, net := range m.Allocation.MachineNetworks {
		if net.Private {
			private = net
		} else if net.Underlay {
			underlay = net
		}
	}

	// Allow all prefixes of the private network
	if private != nil {
		cidrs = append(cidrs, private.Prefixes...)
	}
	for _, i := range ips[m.Allocation.ProjectID] {
		// No need to add /32 addresses of the primary network to the whitelist.
		if private != nil && private.ContainsIP(i.IPAddress) {
			continue
		}
		// Do not allow underlay addresses to be announced.
		if underlay != nil && underlay.ContainsIP(i.IPAddress) {
			continue
		}
		// Allow all other ip addresses allocated for the project.
		cidrs = append(cidrs, fmt.Sprintf("%s/32", i.IPAddress))
	}
	return v12.NewBGPFilter(vnis, cidrs)
}

func makeBGPFilter(m metal.Machine, vrf string, ips metal.IPsMap, iMap metal.ImageMap) v12.BGPFilter {
	var filter v12.BGPFilter
	if m.IsFirewall(iMap) {
		// vrf "default" means: the firewall was successfully allocated and the switch port configured
		// otherwise the port is still not configured yet (pxe-setup) and a BGPFilter would break the install routine
		if vrf == "default" {
			filter = MakeBGPFilterFirewall(m)
		}
	} else {
		filter = MakeBGPFilterMachine(m, ips)
	}
	return filter
}

func MakeSwitchNics(s *metal.Switch, ips metal.IPsMap, iMap metal.ImageMap, machines metal.Machines) v12.SwitchNics {
	machinesByID := map[string]*metal.Machine{}
	for i, m := range machines {
		machinesByID[m.ID] = &machines[i]
	}
	machinesBySwp := map[string]*metal.Machine{}
	for mid, metalConnections := range s.MachineConnections {
		for _, mc := range metalConnections {
			if mid == mc.MachineID {
				machinesBySwp[mc.Nic.Name] = machinesByID[mid]
				break
			}
		}
	}
	nics := v12.SwitchNics{}
	for _, n := range s.Nics {
		m := machinesBySwp[n.Name]
		var filter *v12.BGPFilter
		if m != nil && m.Allocation != nil {
			f := makeBGPFilter(*m, n.Vrf, ips, iMap)
			filter = &f
		}
		nic := v12.SwitchNic{
			MacAddress: string(n.MacAddress),
			Name:       n.Name,
			Vrf:        n.Vrf,
			BGPFilter:  filter,
		}
		nics = append(nics, nic)
	}
	return nics
}

func makeSwitchCons(s *metal.Switch) []v12.SwitchConnection {
	var cons []v12.SwitchConnection
	for _, metalConnections := range s.MachineConnections {
		for _, mc := range metalConnections {
			nic := v12.SwitchNic{
				MacAddress: string(mc.Nic.MacAddress),
				Name:       mc.Nic.Name,
				Vrf:        mc.Nic.Vrf,
			}
			con := v12.SwitchConnection{
				Nic:       nic,
				MachineID: mc.MachineID,
			}
			cons = append(cons, con)
		}
	}
	return cons
}

func findSwitchReferencedEntites(s *metal.Switch, ds *datastore.RethinkStore, logger *zap.SugaredLogger) (*metal.Partition, metal.IPsMap, metal.ImageMap, metal.Machines) {
	var err error

	var p *metal.Partition
	var m metal.Machines
	if s.PartitionID != "" {
		p, err = ds.FindPartition(s.PartitionID)
		if err != nil {
			logger.Errorw("switch references partition, but partition cannot be found in database", "switchID", s.ID, "partitionID", s.PartitionID, "error", err)
		}

		err = ds.SearchMachines(&datastore.MachineSearchQuery{PartitionID: &s.PartitionID}, &m)
		if err != nil {
			logger.Errorw("could not search machines of partition", "switchID", s.ID, "partitionID", s.PartitionID, "error", err)
		}
	}

	ips, err := ds.ListIPs()
	if err != nil {
		logger.Errorw("ips could not be listed", "error", err)
	}

	imgs, err := ds.ListImages()
	if err != nil {
		logger.Errorw("images could not be listed", "error", err)
	}

	return p, ips.ByProjectID(), imgs.ByID(), m
}

func MakeSwitchResponseList(ss []metal.Switch, ds *datastore.RethinkStore, logger *zap.SugaredLogger) []*v12.SwitchResponse {
	pMap, ips, iMap := getSwitchReferencedEntityMaps(ds, logger)
	var result []*v12.SwitchResponse
	m, err := ds.ListMachines()
	if err != nil {
		logger.Errorw("could not find machines")
	}
	for _, sw := range ss {
		var p *metal.Partition
		if sw.PartitionID != "" {
			partitionEntity := pMap[sw.PartitionID]
			p = &partitionEntity
		}

		nics := MakeSwitchNics(&sw, ips, iMap, m)
		cons := makeSwitchCons(&sw)
		result = append(result, v12.NewSwitchResponse(&sw, p, nics, cons))
	}

	return result
}

func getSwitchReferencedEntityMaps(ds *datastore.RethinkStore, logger *zap.SugaredLogger) (metal.PartitionMap, metal.IPsMap, metal.ImageMap) {
	p, err := ds.ListPartitions()
	if err != nil {
		logger.Errorw("partitions could not be listed", "error", err)
	}

	ips, err := ds.ListIPs()
	if err != nil {
		logger.Errorw("ips could not be listed", "error", err)
	}

	imgs, err := ds.ListImages()
	if err != nil {
		logger.Errorw("images could not be listed", "error", err)
	}

	return p.ByID(), ips.ByProjectID(), imgs.ByID()
}
