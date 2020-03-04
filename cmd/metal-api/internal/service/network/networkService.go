package network

import (
	restful "github.com/emicklei/go-restful"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service"
)

func (r networkResource) webService() *restful.WebService {
	ws := new(restful.WebService)
	ws.
		Path(service.BasePath + "v1/network").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	tags := []string{"network"}

	r.addListNetworksRoute(ws, tags)

	r.addFindNetworkRoute(ws, tags)
	r.addFindNetworksRoute(ws, tags)

	r.addCreateNetworkRoute(ws, tags)
	r.addUpdateNetworkRoute(ws, tags)
	r.addDeleteNetworkRoute(ws, tags)

	r.addAllocateNetworkRoute(ws, tags)
	r.addFreeNetworkRoute(ws, tags)

	return ws
}
