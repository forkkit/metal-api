package size

import (
	"github.com/emicklei/go-restful"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/datastore"
)

type sizeResource struct {
	ds *datastore.RethinkStore
}

// NewSize returns a webservice for size specific endpoints.
func NewSize(ds *datastore.RethinkStore) *restful.WebService {
	r := sizeResource{
		ds: ds,
	}
	return r.webService()
}
