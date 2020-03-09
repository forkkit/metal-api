package v1

import (
	mdv1 "github.com/metal-stack/masterdata-api/api/v1"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service/helper"
)

func NewFirewallResponse(fw *metal.Image) *FirewallResponse { //TODO
	if fw == nil {
		return nil
	}
	return &FirewallResponse{
		Firewall: ToFirewall(fw),
	}
}

func ToFirewall(fw *metal.Image) *Firewall { //TODO
	return &Firewall{
		Common: &Common{
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
