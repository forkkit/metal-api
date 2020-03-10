package v1

import (
	mdv1 "github.com/metal-stack/masterdata-api/api/v1"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/datastore"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service/helper"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/utils"
	r "gopkg.in/rethinkdb/rethinkdb-go.v5"
	"strconv"
)

func NewNetworkResponse(network *metal.Network, usage *metal.NetworkUsage) *NetworkResponse {
	if network == nil {
		return nil
	}

	return &NetworkResponse{
		Network:          ToNetwork(network),
		NetworkImmutable: ToNetworkImmutable(network),
		Usage:            ToNetworkUsage(usage),
	}
}

func ToNetwork(network *metal.Network) *Network {
	if network == nil {
		return nil
	}
	return &Network{
		Common: &Common{
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

func ToNetworkImmutable(network *metal.Network) *NetworkImmutable {
	return &NetworkImmutable{
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

func ToNetworkUsage(usage *metal.NetworkUsage) *NetworkUsage {
	return &NetworkUsage{
		AvailableIPs:      usage.AvailableIPs,
		UsedIPs:           usage.UsedIPs,
		AvailablePrefixes: usage.AvailablePrefixes,
		UsedPrefixes:      usage.UsedPrefixes,
	}
}

// GenerateTerm generates the network search query term.
func (x *NetworkSearchQuery) GenerateTerm(rs *datastore.RethinkStore) *r.Term {
	q := *rs.NetworkTable()

	if x.ID != nil {
		q = q.Filter(func(row r.Term) r.Term {
			return row.Field("id").Eq(*x.ID)
		})
	}

	if x.ProjectID != nil {
		q = q.Filter(func(row r.Term) r.Term {
			return row.Field("projectid").Eq(*x.ProjectID)
		})
	}

	if x.PartitionID != nil {
		q = q.Filter(func(row r.Term) r.Term {
			return row.Field("partitionid").Eq(*x.PartitionID)
		})
	}

	if x.ParentNetworkID != nil {
		q = q.Filter(func(row r.Term) r.Term {
			return row.Field("parentnetworkid").Eq(*x.ParentNetworkID)
		})
	}

	if x.Name != nil {
		q = q.Filter(func(row r.Term) r.Term {
			return row.Field("name").Eq(*x.Name)
		})
	}

	if x.Vrf != nil {
		q = q.Filter(func(row r.Term) r.Term {
			return row.Field("vrf").Eq(*x.Vrf)
		})
	}

	if x.Nat != nil {
		q = q.Filter(func(row r.Term) r.Term {
			return row.Field("nat").Eq(*x.Nat)
		})
	}

	if x.PrivateSuper != nil {
		q = q.Filter(func(row r.Term) r.Term {
			return row.Field("privatesuper").Eq(*x.PrivateSuper)
		})
	}

	if x.Underlay != nil {
		q = q.Filter(func(row r.Term) r.Term {
			return row.Field("underlay").Eq(*x.Underlay)
		})
	}

	for k, v := range x.Labels {
		q = q.Filter(func(row r.Term) r.Term {
			return row.Field("labels").Field(k).Eq(v)
		})
	}

	for _, prefix := range x.Prefixes {
		ip, length := utils.SplitCIDR(prefix.GetValue())

		q = q.Filter(func(row r.Term) r.Term {
			return row.Field("prefixes").Map(func(p r.Term) r.Term {
				return p.Field("ip")
			}).Contains(r.Expr(ip))
		})

		if length != nil {
			q = q.Filter(func(row r.Term) r.Term {
				return row.Field("prefixes").Map(func(p r.Term) r.Term {
					return p.Field("length")
				}).Contains(r.Expr(strconv.Itoa(*length)))
			})
		}
	}

	for _, destPrefix := range x.DestinationPrefixes {
		ip, length := utils.SplitCIDR(destPrefix.GetValue())

		q = q.Filter(func(row r.Term) r.Term {
			return row.Field("destinationprefixes").Map(func(dp r.Term) r.Term {
				return dp.Field("ip")
			}).Contains(r.Expr(ip))
		})

		if length != nil {
			q = q.Filter(func(row r.Term) r.Term {
				return row.Field("destinationprefixes").Map(func(dp r.Term) r.Term {
					return dp.Field("length")
				}).Contains(r.Expr(strconv.Itoa(*length)))
			})
		}
	}

	return &q
}
