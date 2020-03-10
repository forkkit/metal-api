package v1

import (
	mdv1 "github.com/metal-stack/masterdata-api/api/v1"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service/helper"
)

func NewPartitionResponse(p *metal.Partition) *PartitionResponse {
	if p == nil {
		return nil
	}
	return &PartitionResponse{
		Partition: ToPartition(p),
	}
}

func ToPartition(p *metal.Partition) *Partition {
	return &Partition{
		Common: &Common{
			Meta: &mdv1.Meta{
				Id:          p.GetID(),
				Apiversion:  "v1",
				Version:     1,
				CreatedTime: helper.ToTimestamp(p.Created),
				UpdatedTime: helper.ToTimestamp(p.Changed),
			},
			Name:        helper.ToStringValue(p.Name),
			Description: helper.ToStringValue(p.Description),
		},
		ImageURL:                   helper.ToStringValue(p.BootConfiguration.ImageURL),
		KernelURL:                  helper.ToStringValue(p.BootConfiguration.KernelURL),
		CommandLine:                helper.ToStringValue(p.BootConfiguration.CommandLine),
		MgmtServiceAddress:         helper.ToStringValue(p.MgmtServiceAddress),
		PrivateNetworkPrefixLength: helper.ToInt64Value(p.PrivateNetworkPrefixLength),
	}
}
