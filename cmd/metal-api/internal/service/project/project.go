package project

import (
	"github.com/emicklei/go-restful"
	"github.com/metal-stack/masterdata-api/api/v1"
	mdm "github.com/metal-stack/masterdata-api/pkg/client"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/datastore"
)

type projectResource struct {
	ds  *datastore.RethinkStore
	mdc mdm.Client
}

// NewProjectService returns a webservice for project specific endpoints.
func NewProjectService(ds *datastore.RethinkStore, mdc mdm.Client) *restful.WebService {
	r := projectResource{
		ds:  ds,
		mdc: mdc,
	}
	return r.webService()
}

type ProjectResponse struct {
	v1.Project
}

type ProjectFindRequest struct {
	v1.ProjectFindRequest
}
