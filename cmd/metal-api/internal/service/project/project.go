package project

import (
	"github.com/emicklei/go-restful"
	mdm "github.com/metal-stack/masterdata-api/pkg/client"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/datastore"
)

type projectResource struct {
	ds  *datastore.RethinkStore
	mdc mdm.Client
}

// NewProject returns a webservice for project specific endpoints.
func NewProject(ds *datastore.RethinkStore, mdc mdm.Client) *restful.WebService {
	r := projectResource{
		ds:  ds,
		mdc: mdc,
	}
	return r.webService()
}
