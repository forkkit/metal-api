package partition

import (
	"github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service/helper"
	v1 "github.com/metal-stack/metal-api/cmd/metal-api/internal/service/v1"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/utils"
	"github.com/metal-stack/metal-lib/httperrors"
	"github.com/metal-stack/metal-lib/zapup"
	"go.uber.org/zap"
	"net/http"
)

func (r partitionResource) addFindPartitionRoute(ws *restful.WebService, tags []string) {
	ws.Route(ws.GET("/{id}").
		To(r.findPartition).
		Operation("findPartition").
		Doc("get Partition by id").
		Param(ws.PathParameter("id", "identifier of the Partition").DataType("string")).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Writes(v1.PartitionResponse{}).
		Returns(http.StatusOK, "OK", v1.PartitionResponse{}).
		DefaultReturns("Error", httperrors.HTTPErrorResponse{}))
}

func (r partitionResource) findPartition(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("id")

	p, err := r.DS.FindPartition(id)
	if helper.CheckError(request, response, utils.CurrentFuncName(), err) {
		return
	}
	err = response.WriteHeaderAndEntity(http.StatusOK, v1.NewPartitionResponse(p))
	if err != nil {
		zapup.MustRootLogger().Error("Failed to send response", zap.Error(err))
		return
	}
}
