package partition

import (
	restful "github.com/emicklei/go-restful"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service"
)

func (r partitionResource) webService() *restful.WebService {
	ws := new(restful.WebService)
	ws.
		Path(service.BasePath + "v1/partition").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	tags := []string{"Partition"}

	r.addFindPartitionRoute(ws, tags)
	r.addListPartitionsRoute(ws, tags)

	r.addCreatePartitionRoute(ws, tags)
	r.addUpdatePartitionRoute(ws, tags)
	r.addDeletePartitionRoute(ws, tags)

	r.addListPartitionCapacitiesRoute(ws, tags)

	return ws
}
