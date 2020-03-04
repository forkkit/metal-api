package network

import (
	"github.com/emicklei/go-restful"
	mdm "github.com/metal-stack/masterdata-api/pkg/client"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/datastore"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/ipam"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service"
	"github.com/metal-stack/metal-lib/zapup"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

type networkResource struct {
	service.WebResource
	ipamer ipam.IPAMer
	mdc    mdm.Client
}

// NewNetwork returns a webservice for network specific endpoints.
func NewNetwork(ds *datastore.RethinkStore, ipamer ipam.IPAMer, mdc mdm.Client) *restful.WebService {
	r := networkResource{
		WebResource: service.WebResource{
			DS: ds,
		},
		ipamer: ipamer,
		mdc:    mdc,
	}
	nuc := networkUsageCollector{r: &r}
	err := prometheus.Register(nuc)
	if err != nil {
		zapup.MustRootLogger().Error("Failed to register prometheus", zap.Error(err))
	}
	return r.webService()
}
