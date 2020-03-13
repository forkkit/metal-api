package helper

import (
	mdmv1 "github.com/metal-stack/masterdata-api/api/v1"
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
			Meta: &mdmv1.Meta{
				Id:          p.GetID(),
				Apiversion:  "v1",
				Version:     1,
				CreatedTime: util.TimestampProto(p.Created),
				UpdatedTime: util.TimestampProto(p.Changed),
			},
			Name:        util.StringProto(p.Name),
			Description: util.StringProto(p.Description),
		},
		ImageURL:                   util.StringProto(p.BootConfiguration.ImageURL),
		KernelURL:                  util.StringProto(p.BootConfiguration.KernelURL),
		CommandLine:                util.StringProto(p.BootConfiguration.CommandLine),
		MgmtServiceAddress:         util.StringProto(p.MgmtServiceAddress),
		PrivateNetworkPrefixLength: util.UInt32Proto(p.PrivateNetworkPrefixLength),
	}
}
