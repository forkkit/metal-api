package network

import (
	restful "github.com/emicklei/go-restful"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service"
	v12 "github.com/metal-stack/metal-api/cmd/metal-api/internal/service/proto/v1"
	"github.com/metal-stack/metal-lib/httperrors"
	"net/http"
)

func (r *networkResource) webService() *restful.WebService {
	return service.Build(&service.WebService{
		Version: service.V1,
		Path:    "network",
		Routes: []*service.Route{
			{
				Method:  http.MethodGet,
				SubPath: "/",
				Doc:     "get all networks",
				Access:  metal.ViewAccess,
				Writes:  []v12.NetworkResponse{},
				Handler: r.listNetworks,
			},
			{
				Method:        http.MethodGet,
				SubPath:       "/{id}",
				PathParameter: service.NewPathParameter("id", "identifier of the network"),
				Doc:           "get network by id",
				Access:        metal.ViewAccess,
				Writes:        v12.NetworkResponse{},
				Handler:       r.findNetwork,
			},
			{
				Method:  http.MethodPost,
				SubPath: "/find",
				Doc:     "get all networks that match given properties",
				Access:  metal.ViewAccess,
				Reads:   v12.NetworkFindRequest{},
				Writes:  []v12.NetworkResponse{},
				Handler: r.findNetworks,
			},
			{
				Method:  http.MethodPut,
				SubPath: "/",
				Doc:     "create a network. if the given ID already exists a conflict is returned",
				Access:  metal.AdminAccess,
				Reads:   v12.NetworkCreateRequest{},
				Writes:  v12.NetworkResponse{},
				Returns: []*service.Return{
					service.NewReturn(http.StatusCreated, "Created", v12.NetworkResponse{}),
					service.NewReturn(http.StatusConflict, "Conflict", httperrors.HTTPErrorResponse{}),
				},
				Handler: r.createNetwork,
			},
			{
				Method:  http.MethodPost,
				SubPath: "/",
				Doc:     "updates a network. if the network was changed since this one was read, a conflict is returned",
				Access:  metal.AdminAccess,
				Reads:   v12.NetworkUpdateRequest{},
				Writes:  v12.NetworkResponse{},
				Returns: []*service.Return{
					service.NewReturn(http.StatusOK, "OK", v12.NetworkResponse{}),
					service.NewReturn(http.StatusConflict, "Conflict", httperrors.HTTPErrorResponse{}),
				},
				Handler: r.updateNetwork,
			},
			{
				Method:        http.MethodDelete,
				SubPath:       "/{id}",
				PathParameter: service.NewPathParameter("id", "identifier of the network"),
				Doc:           "deletes a network and returns the deleted entity",
				Access:        metal.AdminAccess,
				Writes:        v12.NetworkResponse{},
				Handler:       r.deleteNetwork,
			},
			{
				Method:        http.MethodPost,
				SubPath:       "/free/{id}",
				PathParameter: service.NewPathParameter("id", "identifier of the network"),
				Doc:           "frees a network",
				Access:        metal.EditAccess,
				Writes:        v12.NetworkResponse{},
				Returns: []*service.Return{
					service.NewReturn(http.StatusOK, "OK", v12.NetworkResponse{}),
					service.NewReturn(http.StatusConflict, "Conflict", httperrors.HTTPErrorResponse{}),
				},
				Handler: r.freeNetwork,
			},
			{
				Method:  http.MethodPost,
				SubPath: "/allocate",
				Doc:     "allocates a child network from a partition's private super network",
				Access:  metal.EditAccess,
				Reads:   v12.NetworkAllocateRequest{},
				Writes:  v12.NetworkResponse{},
				Returns: []*service.Return{
					service.NewReturn(http.StatusCreated, "Created", v12.NetworkResponse{}),
					service.NewReturn(http.StatusConflict, "Conflict", httperrors.HTTPErrorResponse{}),
				},
				Handler: r.allocateNetwork,
			},
		},
	})
}
