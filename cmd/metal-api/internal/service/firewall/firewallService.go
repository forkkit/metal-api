package firewall

import (
	restful "github.com/emicklei/go-restful"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service"
)

// webService creates the webservice endpoint
func (r firewallResource) webService() *restful.WebService {
	ws := new(restful.WebService)
	ws.
		Path(service.BasePath + "v1/firewall").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	tags := []string{"firewall"}

	r.addListFirewallsRoute(ws, tags)
	r.addFindFirewallRoute(ws, tags)
	r.addFindFirewallsRoute(ws, tags)

	r.addAllocateFirewallRoute(ws, tags)

	return ws
}
