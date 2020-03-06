package sw

import (
	"github.com/emicklei/go-restful"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/datastore"
)

type switchResource struct {
	ds *datastore.RethinkStore
}

// NewSwitch returns a webservice for switch specific endpoints.
func NewSwitch(ds *datastore.RethinkStore) *restful.WebService {
	r := switchResource{
		ds: ds,
	}
	return r.webService()
}
