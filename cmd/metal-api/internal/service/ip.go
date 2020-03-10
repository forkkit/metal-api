package service

import (
	mdv1 "github.com/metal-stack/masterdata-api/api/v1"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
	"github.com/metal-stack/metal-api/pkg/helper"
	"github.com/metal-stack/metal-api/pkg/proto/v1"
	"strings"
)

func NewIPResponse(ip *metal.IP) *v1.IPResponse {
	return &v1.IPResponse{
		IP: ToIP(ip),
		Identifiable: &v1.IPIdentifiable{
			IPAddress: ip.IPAddress,
		},
	}
}

func ToIP(ip *metal.IP) *v1.IP {
	return &v1.IP{
		Common: &v1.Common{
			Meta: &mdv1.Meta{
				Id:          ip.GetID(),
				Apiversion:  "v1",
				Version:     1,
				CreatedTime: helper.ToTimestamp(ip.Created),
				UpdatedTime: helper.ToTimestamp(ip.Changed),
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

func toIPType(ipType metal.IPType) v1.IP_Type {
	if strings.EqualFold(string(ipType), "ephemeral") {
		return v1.IP_EPHEMERAL
	}
	return v1.IP_STATIC
}
