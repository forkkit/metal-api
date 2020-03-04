package image

import (
	"github.com/emicklei/go-restful"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/datastore"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service"
)

type imageResource struct {
	service.WebResource
}

// NewImage returns a webservice for image specific endpoints.
func NewImage(ds *datastore.RethinkStore) *restful.WebService {
	r := imageResource{
		WebResource: service.WebResource{
			DS: ds,
		},
	}
	return r.webService()
}
