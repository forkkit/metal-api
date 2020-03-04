package sw

import (
	restful "github.com/emicklei/go-restful"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service"
)

func (r switchResource) webService() *restful.WebService {
	ws := new(restful.WebService)
	ws.
		Path(service.BasePath + "v1/switch").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	tags := []string{"switch"}

	r.addFindSwitchRoute(ws, tags)
	r.addListSwitchesRoute(ws, tags)

	r.addRegisterSwitchRoute(ws, tags)
	r.addDeleteSwitchRoute(ws, tags)

	return ws
}
