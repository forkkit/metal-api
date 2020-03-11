package v1

import (
	"github.com/metal-stack/metal-api/pkg/util"
	r "gopkg.in/rethinkdb/rethinkdb-go.v5"
	"strconv"
)

// GenerateTerm generates the network search query term.
func (nw *NetworkSearchQuery) GenerateTerm(q r.Term) *r.Term {
	if nw.ID != nil {
		q = q.Filter(func(row r.Term) r.Term {
			return row.Field("id").Eq(nw.ID.GetValue())
		})
	}

	if nw.ProjectID != nil {
		q = q.Filter(func(row r.Term) r.Term {
			return row.Field("projectid").Eq(nw.ProjectID.GetValue())
		})
	}

	if nw.PartitionID != nil {
		q = q.Filter(func(row r.Term) r.Term {
			return row.Field("partitionid").Eq(nw.PartitionID.GetValue())
		})
	}

	if nw.ParentNetworkID != nil {
		q = q.Filter(func(row r.Term) r.Term {
			return row.Field("parentnetworkid").Eq(nw.ParentNetworkID.GetValue())
		})
	}

	if nw.Name != nil {
		q = q.Filter(func(row r.Term) r.Term {
			return row.Field("name").Eq(nw.Name.GetValue())
		})
	}

	if nw.Vrf != nil {
		q = q.Filter(func(row r.Term) r.Term {
			return row.Field("vrf").Eq(nw.Vrf.GetValue())
		})
	}

	if nw.Nat != nil {
		q = q.Filter(func(row r.Term) r.Term {
			return row.Field("nat").Eq(nw.Nat.GetValue())
		})
	}

	if nw.PrivateSuper != nil {
		q = q.Filter(func(row r.Term) r.Term {
			return row.Field("privatesuper").Eq(nw.PrivateSuper.GetValue())
		})
	}

	if nw.Underlay != nil {
		q = q.Filter(func(row r.Term) r.Term {
			return row.Field("underlay").Eq(nw.Underlay.GetValue())
		})
	}

	for k, v := range nw.Labels {
		q = q.Filter(func(row r.Term) r.Term {
			return row.Field("labels").Field(k).Eq(v)
		})
	}

	for _, prefix := range nw.Prefixes {
		ip, length := util.SplitCIDR(prefix.GetValue())

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

	for _, destPrefix := range nw.DestinationPrefixes {
		ip, length := util.SplitCIDR(destPrefix.GetValue())

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
