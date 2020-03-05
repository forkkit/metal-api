package image

import (
	restful "github.com/emicklei/go-restful"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service"
	v1 "github.com/metal-stack/metal-api/cmd/metal-api/internal/service/v1"
	"github.com/metal-stack/metal-lib/httperrors"
	"net/http"
)

func (r imageResource) webService() *restful.WebService {
	return service.Build(service.WebResource{
		Version: service.V1,
		Path:    "image",
		Routes: []service.Route{
			{
				Method:  http.MethodGet,
				SubPath: "/",
				Doc:     "get all images",
				Writes:  []v1.ImageResponse{},
				Handler: r.listImages,
			},
			{
				Method:        http.MethodGet,
				SubPath:       "/{id}",
				PathParameter: service.NewPathParameter("id", "identifier of the image"),
				Doc:           "get image by id",
				Writes:        v1.ImageResponse{},
				Handler:       r.findImage,
			},
			{
				Method:  http.MethodPut,
				SubPath: "/",
				Doc:     "create an image. if the given ID already exists a conflict is returned",
				Access:  metal.AdminAccess,
				Reads:   v1.ImageCreateRequest{},
				Writes:  []v1.ImageResponse{},
				Returns: []*service.Return{
					service.NewReturn(http.StatusCreated, "Created", v1.ImageResponse{}),
					service.NewReturn(http.StatusConflict, "Conflict", httperrors.HTTPErrorResponse{}),
				},
				Handler: r.createImage,
			},
			{
				Method:  http.MethodPost,
				SubPath: "/",
				Doc:     "updates an image. if the image was changed since this one was read, a conflict is returned",
				Access:  metal.AdminAccess,
				Reads:   v1.ImageUpdateRequest{},
				Writes:  v1.ImageResponse{},
				Returns: []*service.Return{
					service.NewReturn(http.StatusOK, "OK", v1.ImageResponse{}),
					service.NewReturn(http.StatusConflict, "Conflict", httperrors.HTTPErrorResponse{}),
				},
				Handler: r.updateImage,
			},
			{
				Method:        http.MethodDelete,
				SubPath:       "/{id}",
				PathParameter: service.NewPathParameter("id", "identifier of the image"),
				Doc:           "deletes an image and returns the deleted entity",
				Access:        metal.AdminAccess,
				Writes:        v1.ImageResponse{},
				Handler:       r.deleteImage,
			},
		},
	})
}
