package firewall

import (
	"github.com/emicklei/go-restful"
	mdm "github.com/metal-stack/masterdata-api/pkg/client"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/datastore"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/ipam"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service"
	"github.com/metal-stack/metal-lib/bus"
)

type firewallResource struct {
	service.WebResource
	bus.Publisher
	ipamer ipam.IPAMer
	mdc    mdm.Client
}

// NewFirewall returns a webservice for firewall specific endpoints.
func NewFirewall(
	ds *datastore.RethinkStore,
	ipamer ipam.IPAMer,
	mdc mdm.Client) *restful.WebService {
	r := firewallResource{
		WebResource: service.WebResource{
			DS: ds,
		},
		ipamer: ipamer,
		mdc:    mdc,
	}
	return r.webService()
}
