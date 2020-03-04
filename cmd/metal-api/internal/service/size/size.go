package size

import (
	"github.com/emicklei/go-restful"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/datastore"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service"
)

type sizeResource struct {
	service.WebResource
}

// NewSize returns a webservice for size specific endpoints.
func NewSize(ds *datastore.RethinkStore) *restful.WebService {
	r := sizeResource{
		WebResource: service.WebResource{
			DS: ds,
		},
	}
	return r.webService()
}
