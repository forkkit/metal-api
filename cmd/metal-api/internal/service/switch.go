package service

import (
	mdv1 "github.com/metal-stack/masterdata-api/api/v1"
	"github.com/metal-stack/metal-api/pkg/proto/v1"
	"github.com/metal-stack/metal-api/pkg/util"
	"sort"

	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
)

func NewBGPFilter(vnis, cidrs []string) v1.BGPFilter {
	// Sort VNIs and CIDRs to avoid unnecessary configuration changes on leaf switches
	sort.Strings(vnis)
	sort.Strings(cidrs)
	return v1.BGPFilter{
		VNIs:  util.ToStringValueSlice(vnis...),
		CIDRs: cidrs,
	}
}

type SwitchNics []*v1.SwitchNic

func (ss SwitchNics) ByMac() map[string]v1.SwitchNic {
	res := make(map[string]v1.SwitchNic)
	for _, s := range ss {
		if s == nil {
			continue
		}
		res[s.MacAddress] = *s
	}
	return res
}

func NewSwitchResponse(s *metal.Switch, p *metal.Partition, nics SwitchNics, cons []*v1.SwitchConnection) *v1.SwitchResponse { //TODO nics unused
	if s == nil {
		return nil
	}

	return &v1.SwitchResponse{
		Switch:      ToSwitch(s),
		Partition:   NewPartitionResponse(p),
		Connections: cons,
	}
}

func ToSwitch(s *metal.Switch) *v1.Switch {
	return &v1.Switch{
		Common: &v1.Common{
			Meta: &mdv1.Meta{
				Id:          s.GetID(),
				Apiversion:  "v1",
				Version:     1,
				CreatedTime: util.ToTimestamp(s.Created),
				UpdatedTime: util.ToTimestamp(s.Changed),
			},
			Name:        util.ToStringValue(s.Name),
			Description: util.ToStringValue(s.Description),
		},
		RackID: s.RackID,
		Nics:   ToNICs(s.Nics),
	}
}

func ToNICs(nics metal.Nics) SwitchNics {
	nn := make(SwitchNics, len(nics))
	for i, n := range nics {
		nn[i] = ToNIC(n)
	}
	return nn
}

func ToNIC(nic metal.Nic) *v1.SwitchNic {
	return &v1.SwitchNic{
		MacAddress: string(nic.MacAddress),
		Name:       nic.Name,
		Vrf:        util.ToStringValue(nic.Vrf),
		//BGPFilter:  NewBGPFilter(), //TODO
	}
}

func NewSwitch(r v1.SwitchRegisterRequest) *metal.Switch {
	nics := make(metal.Nics, len(r.Switch.Nics))
	for i, nic := range r.Switch.Nics {
		nics[i] = metal.Nic{
			MacAddress: metal.MacAddress(nic.MacAddress),
			Name:       nic.Name,
			Vrf:        nic.Vrf.GetValue(),
		}
	}

	return &metal.Switch{
		Base: metal.Base{
			ID:          r.Switch.Common.Meta.Id,
			Name:        r.Switch.Common.Name.GetValue(),
			Description: r.Switch.Common.Description.GetValue(),
		},
		PartitionID:        r.GetPartitionID(),
		RackID:             r.Switch.GetRackID(),
		MachineConnections: make(metal.ConnectionMap),
		Nics:               nics,
	}
}
