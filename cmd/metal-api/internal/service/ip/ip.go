package ip

import (
	"github.com/emicklei/go-restful"
	mdm "github.com/metal-stack/masterdata-api/pkg/client"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/datastore"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/ipam"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service"
)

type ipResource struct {
	service.WebResource
	ipamer ipam.IPAMer
	mdc    mdm.Client
}

// NewIP returns a webservice for ip specific endpoints.
func NewIP(ds *datastore.RethinkStore, ipamer ipam.IPAMer, mdc mdm.Client) *restful.WebService {
	ir := ipResource{
		WebResource: service.WebResource{
			DS: ds,
		},
		ipamer: ipamer,
		mdc:    mdc,
	}
	return ir.webService()
}
