package helper

import (
	mdmv1 "github.com/metal-stack/masterdata-api/api/v1"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
	"github.com/metal-stack/metal-api/pkg/proto/v1"
	"github.com/metal-stack/metal-api/pkg/util"
)

func NewSwitchResponse(s *metal.Switch, p *metal.Partition, nics SwitchNics, cons []*v1.SwitchConnection) *v1.SwitchResponse { //TODO nics unused
	if s == nil {
		return nil
	}

	return &v1.SwitchResponse{
		Switch:            ToSwitch(s),
		PartitionResponse: NewPartitionResponse(p),
		Connections:       cons,
	}
}

func ToSwitch(s *metal.Switch) *v1.Switch {
	return &v1.Switch{
		Common: &v1.Common{
			Meta: &mdmv1.Meta{
				Id:          s.ID,
				Apiversion:  "v1",
				Version:     1,
				CreatedTime: util.TimestampProto(s.Created),
				UpdatedTime: util.TimestampProto(s.Changed),
			},
			Name:        util.StringProto(s.Name),
			Description: util.StringProto(s.Description),
		},
		RackID: s.RackID,
		Nics:   ToNICs(s.Nics),
	}
}

func FromNICs(nics SwitchNics) metal.Nics {
	nn := make(metal.Nics, len(nics))
	for i, n := range nics {
		nn[i] = metal.Nic{
			MacAddress: metal.MacAddress(n.MacAddress),
			Name:       n.Name,
			Vrf:        n.Vrf.GetValue(),
			//Neighbors:  FromNICs(), //TODO
		}
	}
	return nn
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
		Vrf:        util.StringProto(nic.Vrf),
		//BGPFilter:  NewBGPFilter(), //TODO
	}
}

func FromSwitch(r v1.SwitchRegisterRequest) *metal.Switch {
	return &metal.Switch{
		Base: metal.Base{
			ID:          r.Switch.Common.Meta.Id,
			Name:        r.Switch.Common.Name.GetValue(),
			Description: r.Switch.Common.Description.GetValue(),
			Created:     util.Time(r.Switch.Common.Meta.CreatedTime),
			Changed:     util.Time(r.Switch.Common.Meta.UpdatedTime),
		},
		PartitionID:        r.PartitionID,
		RackID:             r.Switch.RackID,
		MachineConnections: make(metal.ConnectionMap),
		Nics:               FromNICs(r.Switch.Nics),
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
