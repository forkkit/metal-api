package partition

import (
	"github.com/emicklei/go-restful"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/datastore"
	"github.com/metal-stack/metal-lib/zapup"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

type TopicCreator interface {
	CreateTopic(partitionID, topicFQN string) error
}

type partitionResource struct {
	ds           *datastore.RethinkStore
	topicCreator TopicCreator
}

// NewPartitionService returns a webservice for partition specific endpoints.
func NewPartitionService(ds *datastore.RethinkStore, tc TopicCreator) *restful.WebService {
	r := partitionResource{
		ds:           ds,
		topicCreator: tc,
	}
	pcc := partitionCapacityCollector{r: &r}
	err := prometheus.Register(pcc)
	if err != nil {
		zapup.MustRootLogger().Error("Failed to register prometheus", zap.Error(err))
	}

	return r.webService()
}
