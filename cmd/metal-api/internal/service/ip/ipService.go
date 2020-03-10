package ip

import (
	restful "github.com/emicklei/go-restful"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
	"github.com/metal-stack/metal-api/pkg"
	v1 "github.com/metal-stack/metal-api/pkg/proto/v1"
	"github.com/metal-stack/metal-lib/httperrors"
	"net/http"
)

func (r *ipResource) webService() *restful.WebService {
	return service.Build(&service.WebService{
		Version: service.V1,
		Path:    "ip",
		Routes: []*service.Route{
			{
				Method:  http.MethodGet,
				SubPath: "/",
				Doc:     "get all ips",
				Access:  metal.ViewAccess,
				Writes:  []v1.IPResponse{},
				Handler: r.listIPs,
			},
			{
				Method:        http.MethodGet,
				SubPath:       "/{id}",
				PathParameter: service.NewPathParameter("id", "identifier of the ip"),
				Doc:           "get ip by id",
				Access:        metal.ViewAccess,
				Writes:        v1.IPResponse{},
				Handler:       r.findIP,
			},
			{
				Method:  http.MethodPost,
				SubPath: "/find",
				Doc:     "get all ips that match given properties",
				Access:  metal.ViewAccess,
				Reads:   v1.IPFindRequest{},
				Writes:  []v1.IPResponse{},
				Handler: r.findIPs,
			},
			{
				Method:  http.MethodPost,
				SubPath: "/allocate",
				Doc:     "allocate an ip in the given network",
				Access:  metal.EditAccess,
				Reads:   v1.IPAllocateRequest{},
				Writes:  v1.IPResponse{},
				Returns: []*service.Return{
					service.NewReturn(http.StatusCreated, "Created", v1.IPResponse{}),
				},
				Handler: r.allocateIP,
			},
			{
				Method:        http.MethodPost,
				SubPath:       "/allocate/{ip}",
				PathParameter: service.NewPathParameter("ip", "ip to try to allocate"),
				Doc:           "allocate a specific ip in the given network",
				Access:        metal.EditAccess,
				Reads:         v1.IPAllocateRequest{},
				Writes:        v1.IPResponse{},
				Returns: []*service.Return{
					service.NewReturn(http.StatusCreated, "Created", v1.IPResponse{}),
				},
				Handler: r.allocateSpecificIP,
			},
			{
				Method:  http.MethodPost,
				SubPath: "/",
				Doc:     "updates an ip. if the ip was changed since this one was read, a conflict is returned",
				Access:  metal.EditAccess,
				Reads:   v1.IPUpdateRequest{},
				Writes:  v1.IPResponse{},
				Returns: []*service.Return{
					service.NewReturn(http.StatusOK, "OK", v1.IPResponse{}),
					service.NewReturn(http.StatusConflict, "Conflict", httperrors.HTTPErrorResponse{}),
				},
				Handler: r.updateIP,
			},
			{
				Method:        http.MethodPost,
				SubPath:       "/free/{id}",
				PathParameter: service.NewPathParameter("id", "identifier of the ip"),
				Doc:           "frees an ip",
				Access:        metal.EditAccess,
				Writes:        v1.IPResponse{},
				Handler:       r.freeIP,
			},
		},
	})
}
