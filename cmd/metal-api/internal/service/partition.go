package service

import (
	mdv1 "github.com/metal-stack/masterdata-api/api/v1"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
	"github.com/metal-stack/metal-api/pkg/proto/v1"
	"github.com/metal-stack/metal-api/pkg/util"
)

func NewPartitionResponse(p *metal.Partition) *v1.PartitionResponse {
	if p == nil {
		return nil
	}
	return &v1.PartitionResponse{
		Partition: ToPartition(p),
	}
}

func ToPartition(p *metal.Partition) *v1.Partition {
	return &v1.Partition{
		Common: &v1.Common{
			Meta: &mdv1.Meta{
				Id:          p.GetID(),
				Apiversion:  "v1",
				Version:     1,
				CreatedTime: util.ToTimestamp(p.Created),
				UpdatedTime: util.ToTimestamp(p.Changed),
			},
			Name:        util.ToStringValue(p.Name),
			Description: util.ToStringValue(p.Description),
		},
		ImageURL:                   util.ToStringValue(p.BootConfiguration.ImageURL),
		KernelURL:                  util.ToStringValue(p.BootConfiguration.KernelURL),
		CommandLine:                util.ToStringValue(p.BootConfiguration.CommandLine),
		MgmtServiceAddress:         util.ToStringValue(p.MgmtServiceAddress),
		PrivateNetworkPrefixLength: util.ToInt64Value(p.PrivateNetworkPrefixLength),
	}
}
