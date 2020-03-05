package firewall

import (
	restful "github.com/emicklei/go-restful"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service"
	v1 "github.com/metal-stack/metal-api/cmd/metal-api/internal/service/v1"
	"net/http"
)

// webService creates the webservice endpoint
func (r firewallResource) webService() *restful.WebService {
	return service.Build(service.WebResource{
		Version: service.V1,
		Path:    "firewall",
		Routes: []service.Route{
			{
				Method:  http.MethodGet,
				SubPath: "/",
				Doc:     "get all known firewalls",
				Access:  metal.ViewAccess,
				Writes:  []v1.FirewallResponse{},
				Handler: r.listFirewalls,
			},
			{
				Method:        http.MethodGet,
				SubPath:       "/{id}",
				PathParameter: service.NewPathParameter("id", "identifier of the firewall"),
				Doc:           "get firewall by id",
				Access:        metal.ViewAccess,
				Writes:        v1.FirewallResponse{},
				Handler:       r.findFirewall,
			},
			{
				Method:  http.MethodGet,
				SubPath: "/find",
				Doc:     "find firewalls by multiple criteria",
				Access:  metal.ViewAccess,
				Reads:   v1.FirewallFindRequest{},
				Writes:  []v1.FirewallResponse{},
				Handler: r.findFirewalls,
			},
			{
				Method:  http.MethodPost,
				SubPath: "/allocate",
				Doc:     "allocate a firewall",
				Access:  metal.EditAccess,
				Reads:   v1.FirewallCreateRequest{},
				Writes:  []v1.FirewallResponse{},
				Handler: r.allocateFirewall,
			},
		},
	})
}
