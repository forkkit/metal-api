package v1

import (
	mdv1 "github.com/metal-stack/masterdata-api/api/v1"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/datastore"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service/helper"
	r "gopkg.in/rethinkdb/rethinkdb-go.v5"
	"strings"
)

func NewIPResponse(ip *metal.IP) *IPResponse {
	return &IPResponse{
		IP: ToIP(ip),
		Identifiable: &IPIdentifiable{
			IPAddress: ip.IPAddress,
		},
	}
}

func ToIP(ip *metal.IP) *IP {
	return &IP{
		Common: &Common{
			Meta: &mdv1.Meta{
				Id:                   ip.GetID(),
				Apiversion:           "v1",
				Version:              1,
				CreatedTime:          helper.ToTimestamp(ip.Created),
				UpdatedTime:          helper.ToTimestamp(ip.Changed),
			},
			Name:        helper.ToStringValue(ip.Name),
			Description: helper.ToStringValue(ip.Description),
		},
		NetworkID: ip.NetworkID,
		ProjectID: ip.ProjectID,
		Type:      toIPType(ip.Type),
		Tags:      helper.ToStringValueSlice(ip.Tags...),
	}
}

// GenerateTerm generates the project search query term.
func (x *IPFindRequest) GenerateTerm(rs *datastore.RethinkStore) *r.Term {
	q := *rs.IPTable()

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
		x.Tags = append(x.Tags, helper.ToStringValue(metal.IpTag(metal.TagIPMachineID, x.MachineID.GetValue())))
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

func toIPType(ipType metal.IPType) IP_Type {
	if strings.EqualFold(string(ipType), "ephemeral") {
		return IP_EPHEMERAL
	}
	return IP_STATIC
}