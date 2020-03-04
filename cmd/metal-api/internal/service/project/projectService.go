package project

import (
	restful "github.com/emicklei/go-restful"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service"
)

func (r projectResource) webService() *restful.WebService {
	ws := new(restful.WebService)
	ws.
		Path(service.BasePath + "v1/project").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	tags := []string{"project"}

	r.addListProjectsRoute(ws, tags)
	r.addFindProjectRoute(ws, tags)
	r.addFindProjectsRoute(ws, tags)

	return ws
}
