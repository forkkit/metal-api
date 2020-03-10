package v1

import (
	"fmt"
	"github.com/metal-stack/metal-api/pkg/helper"
	"github.com/metal-stack/metal-lib/pkg/tag"
	r "gopkg.in/rethinkdb/rethinkdb-go.v5"
)

func IpTag(key, value string) string {
	return fmt.Sprintf("%s=%s", key, value)
}

// GenerateTerm generates the IP search query term.
func (x *IPFindRequest) GenerateTerm(q r.Term) *r.Term {
	if x.IPAddress != nil {
		q = q.Filter(func(row r.Term) r.Term {
			return row.Field("id").Eq(*x.IPAddress)
		})
	}

	if x.ProjectID != nil {
		q = q.Filter(func(row r.Term) r.Term {
			return row.Field("projectid").Eq(*x.ProjectID)
		})
	}

	if x.NetworkID != nil {
		q = q.Filter(func(row r.Term) r.Term {
			return row.Field("networkid").Eq(*x.NetworkID)
		})
	}

	if x.ParentPrefixCidr != nil {
		q = q.Filter(func(row r.Term) r.Term {
			return row.Field("networkprefix").Eq(*x.ParentPrefixCidr)
		})
	}

	if x.MachineID != nil {
		x.Tags = append(x.Tags, helper.ToStringValue(IpTag(tag.MachineID, x.MachineID.GetValue())))
	}

	for _, tag := range x.Tags {
		q = q.Filter(func(row r.Term) r.Term {
			return row.Field("tags").Contains(r.Expr(tag))
		})
	}

	if x.Type != nil {
		q = q.Filter(func(row r.Term) r.Term {
			return row.Field("type").Eq(*x.Type)
		})
	}

	return &q
}
