package v1

import (
	mdv1 "github.com/metal-stack/masterdata-api/api/v1"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service/helper"
	"sort"

	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
)

func NewBGPFilter(vnis, cidrs []string) BGPFilter {
	// Sort VNIs and CIDRs to avoid unnecessary configuration changes on leaf switches
	sort.Strings(vnis)
	sort.Strings(cidrs)
	return BGPFilter{
		VNIs:  helper.ToStringValueSlice(vnis...),
		CIDRs: cidrs,
	}
}

type SwitchNics []*SwitchNic

func (ss SwitchNics) ByMac() map[string]SwitchNic {
	res := make(map[string]SwitchNic)
	for _, s := range ss {
		if s == nil {
			continue
		}
		res[s.MacAddress] = *s
	}
	return res
}

func NewSwitchResponse(s *metal.Switch, p *metal.Partition, nics SwitchNics, cons []*SwitchConnection) *SwitchResponse { //TODO nics unused
	if s == nil {
		return nil
	}

	return &SwitchResponse{
		Switch:      ToSwitch(s),
		Partition:   NewPartitionResponse(p),
		Connections: cons,
	}
}

func ToSwitch(s *metal.Switch) *Switch {
	return &Switch{
		Common: &Common{
			Meta: &mdv1.Meta{
				Id:          s.GetID(),
				Apiversion:  "v1",
				Version:     1,
				CreatedTime: helper.ToTimestamp(s.Created),
				UpdatedTime: helper.ToTimestamp(s.Changed),
			},
			Name:        helper.ToStringValue(s.Name),
			Description: helper.ToStringValue(s.Description),
		},
		RackID: s.RackID,
		Nics:   ToNICs(s.Nics),
	}
}

func ToNICs(nics metal.Nics) SwitchNics {
	nn := make(SwitchNics, 0, len(nics))
	for i, n := range nics {
		nn[i] = ToNIC(n)
	}
	return nn
}

func ToNIC(nic metal.Nic) *SwitchNic {
	return &SwitchNic{
		MacAddress: string(nic.MacAddress),
		Name:       nic.Name,
		Vrf:        helper.ToStringValue(nic.Vrf),
		BGPFilter:  NewBGPFilter(),
	}
}

func NewSwitch(r SwitchRegisterRequest) *metal.Switch {
	nics := make(metal.Nics, 0, len(r.Switch.Nics))
	for i, nic := range r.Switch.Nics {
		nics[i] = metal.Nic{
			MacAddress: metal.MacAddress(nic.MacAddress),
			Name:       nic.Name,
			Vrf:        nic.Vrf.GetValue(),
		}
	}

	return &metal.Switch{
		Base: metal.Base{
			ID:          r.GetSwitch().GetCommon().GetMeta().ID,
			Name:        r.GetSwitch().GetCommon().GetName().GetValue(),
			Description: r.GetSwitch().GetCommon().GetDescription().GetValue(),
		},
		PartitionID:        r.GetPartitionID(),
		RackID:             r.GetSwitch().GetRackID(),
		MachineConnections: make(metal.ConnectionMap),
		Nics:               nics,
	}
}
