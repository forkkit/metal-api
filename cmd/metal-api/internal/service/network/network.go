package network

import (
	"fmt"
	"github.com/emicklei/go-restful"
	v12 "github.com/metal-stack/masterdata-api/api/v1"
	mdm "github.com/metal-stack/masterdata-api/pkg/client"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/datastore"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/ipam"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
	"github.com/metal-stack/metal-api/pkg/proto/v1"
	"github.com/metal-stack/metal-api/pkg/util"
	"github.com/metal-stack/metal-lib/zapup"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

type networkResource struct {
	ds     *datastore.RethinkStore
	ipamer ipam.IPAMer
	mdc    mdm.Client
}

// NewNetworkService returns a webservice for network specific endpoints.
func NewNetworkService(ds *datastore.RethinkStore, ipamer ipam.IPAMer, mdc mdm.Client) *restful.WebService {
	r := networkResource{
		ds:     ds,
		ipamer: ipamer,
		mdc:    mdc,
	}
	nuc := networkUsageCollector{r: &r}
	err := prometheus.Register(nuc)
	if err != nil {
		zapup.MustRootLogger().Error("Failed to register prometheus", zap.Error(err))
	}
	return r.webService()
}

func GetNetworkUsage(nw *metal.Network, ipamer ipam.IPAMer) *metal.NetworkUsage {
	usage := &metal.NetworkUsage{}
	if nw == nil {
		return usage
	}
	for _, prefix := range nw.Prefixes {
		u, err := ipamer.PrefixUsage(prefix.String())
		if err != nil {
			continue
		}
		usage.AvailableIPs = usage.AvailableIPs + u.AvailableIPs
		usage.UsedIPs = usage.UsedIPs + u.UsedIPs
		usage.AvailablePrefixes = usage.AvailablePrefixes + u.AvailablePrefixes
		usage.UsedPrefixes = usage.UsedPrefixes + u.UsedPrefixes
	}
	return usage
}

func CheckAnyIPOfPrefixesInUse(ips []metal.IP, prefixes metal.Prefixes) error {
	for _, ip := range ips {
		for _, p := range prefixes {
			if ip.ParentPrefixCidr == p.String() {
				return fmt.Errorf("prefix %s has ip %s in use", p.String(), ip.IPAddress)
			}
		}
	}
	return nil
}

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
