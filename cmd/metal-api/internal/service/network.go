package service

import (
	mdv1 "github.com/metal-stack/masterdata-api/api/v1"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
	"github.com/metal-stack/metal-api/pkg/helper"
	"github.com/metal-stack/metal-api/pkg/proto/v1"
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

func ToNetwork(network *metal.Network) *v1.Network {
	if network == nil {
		return nil
	}
	return &v1.Network{
		Common: &v1.Common{
			Meta: &mdv1.Meta{
				Id:          network.GetID(),
				Apiversion:  "v1",
				Version:     1,
				CreatedTime: helper.ToTimestamp(network.Created),
				UpdatedTime: helper.ToTimestamp(network.Changed),
			},
			Name:        helper.ToStringValue(network.Name),
			Description: helper.ToStringValue(network.Description),
		},
		PartitionID: helper.ToStringValue(network.PartitionID),
		ProjectID:   helper.ToStringValue(network.ProjectID),
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
		Vrf:                 helper.ToUInt64Value(network.Vrf),
		//VrfShared:           helper.ToBoolValue(network.VrfShared), //TODO network.VrfShared is not defined
		ParentNetworkID: helper.ToStringValue(network.ParentNetworkID),
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
