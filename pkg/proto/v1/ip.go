package v1

import (
	"fmt"
	"github.com/metal-stack/metal-api/pkg/util"
	"github.com/metal-stack/metal-lib/pkg/tag"
	r "gopkg.in/rethinkdb/rethinkdb-go.v5"
)

func IpTag(key, value string) string {
	return fmt.Sprintf("%s=%s", key, value)
}

// GenerateTerm generates the IP search query term.
func (ip *IPFindRequest) GenerateTerm(q r.Term) *r.Term {
	if ip.IPAddress != nil {
		q = q.Filter(func(row r.Term) r.Term {
			return row.Field("id").Eq(ip.IPAddress.GetValue())
		})
	}

	if ip.ProjectID != nil {
		q = q.Filter(func(row r.Term) r.Term {
			return row.Field("projectid").Eq(ip.ProjectID.GetValue())
		})
	}

	if ip.NetworkID != nil {
		q = q.Filter(func(row r.Term) r.Term {
			return row.Field("networkid").Eq(ip.NetworkID.GetValue())
		})
	}

	if ip.ParentPrefixCidr != nil {
		q = q.Filter(func(row r.Term) r.Term {
			return row.Field("networkprefix").Eq(ip.ParentPrefixCidr.GetValue())
		})
	}

	if ip.MachineID != nil {
		ip.Tags = append(ip.Tags, util.ToStringValue(IpTag(tag.MachineID, ip.MachineID.GetValue())))
	}

	for _, tag := range ip.Tags {
		q = q.Filter(func(row r.Term) r.Term {
			return row.Field("tags").Contains(r.Expr(tag.GetValue()))
		})
	}

	if ip.Type != nil {
		q = q.Filter(func(row r.Term) r.Term {
			return row.Field("type").Eq(ip.Type.GetValue())
		})
	}

	return &q
}
