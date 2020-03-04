package partition

import (
	"github.com/emicklei/go-restful"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/datastore"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service"
	"github.com/metal-stack/metal-lib/zapup"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

type TopicCreater interface {
	CreateTopic(partitionID, topicFQN string) error
}

type partitionResource struct {
	service.WebResource
	topicCreater TopicCreater
}

// NewPartition returns a webservice for partition specific endpoints.
func NewPartition(ds *datastore.RethinkStore, tc TopicCreater) *restful.WebService {
	r := partitionResource{
		WebResource: service.WebResource{
			DS: ds,
		},
		topicCreater: tc,
	}
	pcc := partitionCapacityCollector{r: &r}
	err := prometheus.Register(pcc)
	if err != nil {
		zapup.MustRootLogger().Error("Failed to register prometheus", zap.Error(err))
	}

	return r.webService()
}
