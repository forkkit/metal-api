package sw

import (
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/datastore"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service/helper"
	v1 "github.com/metal-stack/metal-api/pkg/proto/v1"
	"github.com/metal-stack/metal-api/pkg/util"
	"go.uber.org/zap"
	"sort"
)

type switchResource struct {
	ds *datastore.RethinkStore
}

// NewSwitchService returns a webservice for switch specific endpoints.
func NewSwitchService(ds *datastore.RethinkStore) *restful.WebService {
	r := switchResource{
		ds: ds,
	}
	return r.webService()
}

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
		if nic, has := changedByMac[mac]; has {
			e = nic
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

func MakeSwitchResponse(s *metal.Switch, ds *datastore.RethinkStore, logger *zap.SugaredLogger) *v1.SwitchResponse {
	p, ips, iMap, machines := findSwitchReferencedEntities(s, ds, logger)
	nics := MakeSwitchNics(s, ips, iMap, machines)
	cons := makeSwitchCons(s)
	return helper.NewSwitchResponse(s, p, nics, cons)
}

func MakeBGPFilterFirewall(m metal.Machine) v1.BGPFilter {
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
	return NewBGPFilter(vnis, cidrs)
}

func MakeBGPFilterMachine(m metal.Machine, ips metal.IPsMap) v1.BGPFilter {
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
	return NewBGPFilter(vnis, cidrs)
}

func makeBGPFilter(m metal.Machine, vrf string, ips metal.IPsMap, iMap metal.ImageMap) v1.BGPFilter {
	var filter v1.BGPFilter
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

func MakeSwitchNics(s *metal.Switch, ips metal.IPsMap, iMap metal.ImageMap, machines metal.Machines) helper.SwitchNics {
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
	nics := helper.SwitchNics{}
	for _, n := range s.Nics {
		m := machinesBySwp[n.Name]
		var filter *v1.BGPFilter
		if m != nil && m.Allocation != nil {
			f := makeBGPFilter(*m, n.Vrf, ips, iMap)
			filter = &f
		}
		nic := &v1.SwitchNic{
			MacAddress: string(n.MacAddress),
			Name:       n.Name,
			Vrf:        util.StringProto(n.Vrf),
			BGPFilter:  filter,
		}
		nics = append(nics, nic)
	}
	return nics
}

func makeSwitchCons(s *metal.Switch) []*v1.SwitchConnection {
	var cons []*v1.SwitchConnection
	for _, metalConnections := range s.MachineConnections {
		for _, mc := range metalConnections {
			nic := &v1.SwitchNic{
				MacAddress: string(mc.Nic.MacAddress),
				Name:       mc.Nic.Name,
				Vrf:        util.StringProto(mc.Nic.Vrf),
			}
			con := &v1.SwitchConnection{
				Nic:       nic,
				MachineID: util.StringProto(mc.MachineID),
			}
			cons = append(cons, con)
		}
	}
	return cons
}

func findSwitchReferencedEntities(s *metal.Switch, ds *datastore.RethinkStore, logger *zap.SugaredLogger) (*metal.Partition, metal.IPsMap, metal.ImageMap, metal.Machines) {
	var err error

	var p *metal.Partition
	var m metal.Machines
	if s.PartitionID != "" {
		p, err = ds.FindPartition(s.PartitionID)
		if err != nil {
			logger.Errorw("switch references partition, but partition cannot be found in database", "switchID", s.ID, "partitionID", s.PartitionID, "error", err)
		}

		err = ds.SearchMachines(&v1.MachineSearchQuery{PartitionID: util.StringProto(s.PartitionID)}, &m)
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

func MakeSwitchResponseList(ss []metal.Switch, ds *datastore.RethinkStore, logger *zap.SugaredLogger) []*v1.SwitchResponse {
	pMap, ips, iMap := getSwitchReferencedEntityMaps(ds, logger)
	var result []*v1.SwitchResponse
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
		result = append(result, helper.NewSwitchResponse(&sw, p, nics, cons))
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

func NewBGPFilter(vnis, cidrs []string) v1.BGPFilter {
	// Sort VNIs and CIDRs to avoid unnecessary configuration changes on leaf switches
	sort.Strings(vnis)
	sort.Strings(cidrs)
	return v1.BGPFilter{
		VNIs:  util.StringSliceProto(vnis...),
		CIDRs: cidrs,
	}
}
