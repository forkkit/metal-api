package size

import (
	restful "github.com/emicklei/go-restful"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
	"github.com/metal-stack/metal-api/pkg"
	v1 "github.com/metal-stack/metal-api/pkg/proto/v1"
	"github.com/metal-stack/metal-lib/httperrors"
	"net/http"
)

func (r *sizeResource) webService() *restful.WebService {
	return service.Build(&service.WebService{
		Version: service.V1,
		Path:    "size",
		Routes: []*service.Route{
			{
				Method:  http.MethodGet,
				SubPath: "/",
				Doc:     "get all sizes",
				Writes:  []v1.SizeResponse{},
				Handler: r.listSizes,
			},
			{
				Method:        http.MethodGet,
				SubPath:       "/{id}",
				PathParameter: service.NewPathParameter("id", "identifier of the size"),
				Doc:           "get size by id",
				Writes:        v1.SizeResponse{},
				Handler:       r.findSize,
			},
			{
				Method:  http.MethodPut,
				SubPath: "/",
				Doc:     "create a size. if the given ID already exists a conflict is returned",
				Access:  metal.AdminAccess,
				Reads:   v1.SizeCreateRequest{},
				Writes:  []v1.SizeResponse{},
				Returns: []*service.Return{
					service.NewReturn(http.StatusCreated, "Created", v1.SizeResponse{}),
					service.NewReturn(http.StatusConflict, "Conflict", httperrors.HTTPErrorResponse{}),
				},
				Handler: r.createSize,
			},
			{
				Method:  http.MethodPost,
				SubPath: "/",
				Doc:     "updates a size. if the size was changed since this one was read, a conflict is returned",
				Access:  metal.AdminAccess,
				Reads:   v1.SizeUpdateRequest{},
				Writes:  []v1.SizeResponse{},
				Returns: []*service.Return{
					service.NewReturn(http.StatusOK, "OK", v1.SizeResponse{}),
					service.NewReturn(http.StatusConflict, "Conflict", httperrors.HTTPErrorResponse{}),
				},
				Handler: r.updateSize,
			},
			{
				Method:        http.MethodDelete,
				SubPath:       "/{id}",
				PathParameter: service.NewPathParameter("id", "identifier of the size"),
				Doc:           "deletes an size and returns the deleted entity",
				Access:        metal.AdminAccess,
				Writes:        v1.SizeResponse{},
				Handler:       r.deleteSize,
			},
			{
				Method:  http.MethodPost,
				SubPath: "/from-hardware",
				Doc:     "Searches all sizes for one to match the given hardware specs. If nothing is found, a list of entries is returned, which describe the constraint that did not match",
				Reads:   v1.MachineHardwareExtended{},
				Writes:  []SizeMatchingLog{},
				Handler: r.fromHardware,
			},
		},
	})
}
