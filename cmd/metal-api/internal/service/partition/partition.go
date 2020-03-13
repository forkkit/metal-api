package partition

import (
	"github.com/emicklei/go-restful"
	v12 "github.com/metal-stack/masterdata-api/api/v1"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/datastore"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
	"github.com/metal-stack/metal-api/pkg/proto/v1"
	"github.com/metal-stack/metal-api/pkg/util"
	"github.com/metal-stack/metal-lib/zapup"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

type TopicCreater interface {
	CreateTopic(partitionID, topicFQN string) error
}

type partitionResource struct {
	ds           *datastore.RethinkStore
	topicCreater TopicCreater
}

// NewPartitionService returns a webservice for partition specific endpoints.
func NewPartitionService(ds *datastore.RethinkStore, tc TopicCreater) *restful.WebService {
	r := partitionResource{
		ds:           ds,
		topicCreater: tc,
	}
	pcc := partitionCapacityCollector{r: &r}
	err := prometheus.Register(pcc)
	if err != nil {
		zapup.MustRootLogger().Error("Failed to register prometheus", zap.Error(err))
	}

	return r.webService()
}

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
			Meta: &v12.Meta{
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
