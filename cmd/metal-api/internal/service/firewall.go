package service

import (
	mdv1 "github.com/metal-stack/masterdata-api/api/v1"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
	"github.com/metal-stack/metal-api/pkg/helper"
	"github.com/metal-stack/metal-api/pkg/proto/v1"
)

func NewFirewallResponse(fw *metal.Image) *v1.FirewallResponse { //TODO
	if fw == nil {
		return nil
	}
	return &v1.FirewallResponse{
		Firewall: ToFirewall(fw),
	}
}

func ToFirewall(fw *metal.Image) *v1.Firewall { //TODO
	return &v1.Firewall{
		Common: &v1.Common{
			Meta: &mdv1.Meta{
				Id:          fw.ID,
				Apiversion:  "v1",
				Version:     1,
				CreatedTime: helper.ToTimestamp(fw.Created),
				UpdatedTime: helper.ToTimestamp(fw.Changed),
			},
			Name:        helper.ToStringValue(fw.Name),
			Description: helper.ToStringValue(fw.Description),
		},
		//Ha: helper.ToBoolValue(f.Ha), //TODO
	}
}
