package ip

import (
	restful "github.com/emicklei/go-restful"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service"
)

func (r ipResource) webService() *restful.WebService {
	ws := new(restful.WebService)
	ws.
		Path(service.BasePath + "v1/ip").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	tags := []string{"ip"}

	r.addListIPsRoute(ws, tags)
	r.addFindIPRoute(ws, tags)
	r.addFindIPsRoute(ws, tags)

	r.addAllocateIPRoute(ws, tags)
	r.addUpdateIPRoute(ws, tags)
	r.addFreeIPRoute(ws, tags)

	return ws
}
