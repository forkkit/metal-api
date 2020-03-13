package helper

import (
	v12 "github.com/metal-stack/masterdata-api/api/v1"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
	"github.com/metal-stack/metal-api/pkg/proto/v1"
	"github.com/metal-stack/metal-api/pkg/util"
)

func NewNetworkResponse(network *metal.Network, usage *metal.NetworkUsage) *v1.NetworkResponse {
	if network == nil {
		return nil
	}

	return &v1.NetworkResponse{
		Network:          ToNetwork(network),
		NetworkImmutable: ToNetworkImmutable(network),
		Usage:            ToNetworkUsage(usage),
	}
}

func FromNetwork(network *v1.Network) *metal.Network {
	if network == nil {
		return nil
	}
	return &metal.Network{
		Base: metal.Base{
			ID:          network.Common.Meta.Id,
			Name:        network.Common.Name.GetValue(),
			Description: network.Common.Description.GetValue(),
			Created:     util.Time(network.Common.Meta.CreatedTime),
			Changed:     util.Time(network.Common.Meta.UpdatedTime),
		},
		PartitionID: network.PartitionID.GetValue(),
		ProjectID:   network.ProjectID.GetValue(),
		Labels:      network.Labels,
	}
}

func ToNetwork(network *metal.Network) *v1.Network {
	if network == nil {
		return nil
	}
	return &v1.Network{
		Common: &v1.Common{
			Meta: &v12.Meta{
				Id:          network.GetID(),
				Apiversion:  "v1",
				Version:     1,
				CreatedTime: util.TimestampProto(network.Created),
				UpdatedTime: util.TimestampProto(network.Changed),
			},
			Name:        util.StringProto(network.Name),
			Description: util.StringProto(network.Description),
		},
		PartitionID: util.StringProto(network.PartitionID),
		ProjectID:   util.StringProto(network.ProjectID),
		Labels:      network.Labels,
	}
}

func ToNetworkImmutable(network *metal.Network) *v1.NetworkImmutable {
	return &v1.NetworkImmutable{
		Prefixes:            network.Prefixes.String(),
		DestinationPrefixes: network.DestinationPrefixes.String(),
		Nat:                 network.Nat,
		PrivateSuper:        network.PrivateSuper,
		Underlay:            network.Underlay,
		Vrf:                 util.UInt64Proto(network.Vrf),
		//VrfShared:           helper.BoolProto(network.VrfShared), //TODO network.VrfShared is not defined
		ParentNetworkID: util.StringProto(network.ParentNetworkID),
	}
}

func ToNetworkUsage(usage *metal.NetworkUsage) *v1.NetworkUsage {
	return &v1.NetworkUsage{
		AvailableIPs:      usage.AvailableIPs,
		UsedIPs:           usage.UsedIPs,
		AvailablePrefixes: usage.AvailablePrefixes,
		UsedPrefixes:      usage.UsedPrefixes,
	}
}
