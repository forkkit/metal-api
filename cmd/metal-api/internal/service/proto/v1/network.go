package v1

import (
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/datastore"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/utils"
	r "gopkg.in/rethinkdb/rethinkdb-go.v5"
	"strconv"
)

func NewNetworkResponse(network *metal.Network, usage *metal.NetworkUsage) *NetworkResponse {
	if network == nil {
		return nil
	}

	var parentNetworkID *string
	if network.ParentNetworkID != "" {
		parentNetworkID = &network.ParentNetworkID
	}

	return &NetworkResponse{
		Common: Common{
			Identifiable: Identifiable{
				ID: network.ID,
			},
			Describable: Describable{
				Name:        &network.Name,
				Description: &network.Description,
			},
		},
		NetworkBase: NetworkBase{
			PartitionID: &network.PartitionID,
			ProjectID:   &network.ProjectID,
			Labels:      network.Labels,
		},
		NetworkImmutable: NetworkImmutable{
			Prefixes:            network.Prefixes.String(),
			DestinationPrefixes: network.DestinationPrefixes.String(),
			Nat:                 network.Nat,
			PrivateSuper:        network.PrivateSuper,
			Underlay:            network.Underlay,
			Vrf:                 &network.Vrf,
			ParentNetworkID:     parentNetworkID,
		},
		Usage: NetworkUsage{
			AvailableIPs:      usage.AvailableIPs,
			UsedIPs:           usage.UsedIPs,
			AvailablePrefixes: usage.AvailablePrefixes,
			UsedPrefixes:      usage.UsedPrefixes,
		},
		Timestamps: Timestamps{
			Created: network.Created,
			Changed: network.Changed,
		},
	}
}

// GenerateTerm generates the project search query term.
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
		ip, length := utils.SplitCIDR(prefix)

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
		ip, length := utils.SplitCIDR(destPrefix)

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
