package image

import (
	"github.com/emicklei/go-restful"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/datastore"
)

type imageResource struct {
	ds *datastore.RethinkStore
}

// NewImageService returns a webservice for image specific endpoints.
func NewImageService(ds *datastore.RethinkStore) *restful.WebService {
	r := imageResource{
		ds: ds,
	}
	return r.webService()
}
