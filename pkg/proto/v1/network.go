package v1

import (
	"github.com/metal-stack/metal-api/pkg/helper"
	r "gopkg.in/rethinkdb/rethinkdb-go.v5"
	"strconv"
)

// GenerateTerm generates the network search query term.
func (x *NetworkSearchQuery) GenerateTerm(q r.Term) *r.Term {
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
		ip, length := helper.SplitCIDR(prefix.GetValue())

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
		ip, length := helper.SplitCIDR(destPrefix.GetValue())

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
