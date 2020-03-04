package network

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

func (r networkResource) addFindNetworkRoute(ws *restful.WebService, tags []string) {
	ws.Route(ws.GET("/{id}").
		To(helper.Viewer(r.findNetwork)).
		Operation("findNetwork").
		Doc("get network by id").
		Param(ws.PathParameter("id", "identifier of the network").DataType("string")).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Writes(v1.NetworkResponse{}).
		Returns(http.StatusOK, "OK", v1.NetworkResponse{}).
		DefaultReturns("Error", httperrors.HTTPErrorResponse{}))
}

func (r networkResource) findNetwork(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("id")

	nw, err := r.DS.FindNetworkByID(id)
	if helper.CheckError(request, response, utils.CurrentFuncName(), err) {
		return
	}
	usage := helper.GetNetworkUsage(nw, r.ipamer)
	err = response.WriteHeaderAndEntity(http.StatusOK, v1.NewNetworkResponse(nw, usage))
	if err != nil {
		zapup.MustRootLogger().Error("Failed to send response", zap.Error(err))
		return
	}
}
