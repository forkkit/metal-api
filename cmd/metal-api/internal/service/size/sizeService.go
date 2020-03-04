package size

import (
	restful "github.com/emicklei/go-restful"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service"
)

func (r sizeResource) webService() *restful.WebService {
	ws := new(restful.WebService)
	ws.
		Path(service.BasePath + "v1/size").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	tags := []string{"size"}

	r.addFindSizeRoute(ws, tags)
	r.addListSizesRoute(ws, tags)

	r.addCreateSizeRoute(ws, tags)
	r.addUpdateSizeRoute(ws, tags)
	r.addDeleteSizeRoute(ws, tags)

	r.addFromHardwareRoute(ws, tags)

	return ws
}
