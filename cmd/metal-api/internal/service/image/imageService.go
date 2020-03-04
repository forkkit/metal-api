package image

import (
	restful "github.com/emicklei/go-restful"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service"
)

func (r imageResource) webService() *restful.WebService {
	ws := new(restful.WebService)
	ws.
		Path(service.BasePath + "v1/image").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	tags := []string{"image"}

	r.addListImagesRoute(ws, tags)
	r.addFindImageRoute(ws, tags)

	r.addCreateImageRoute(ws, tags)
	r.addUpdateImageRoute(ws, tags)
	r.addDeleteImageRoute(ws, tags)

	return ws
}
