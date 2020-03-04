package sw

import (
	"github.com/emicklei/go-restful"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/datastore"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service"
)

type switchResource struct {
	service.WebResource
}

// NewSwitch returns a webservice for switch specific endpoints.
func NewSwitch(ds *datastore.RethinkStore) *restful.WebService {
	r := switchResource{
		WebResource: service.WebResource{
			DS: ds,
		},
	}
	return r.webService()
}
