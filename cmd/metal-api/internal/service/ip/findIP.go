package ip

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

func (r ipResource) addFindIPRoute(ws *restful.WebService, tags []string) {
	ws.Route(ws.GET("/{id}").
		To(helper.Viewer(r.findIP)).
		Operation("findIP").
		Doc("get ip by id").
		Param(ws.PathParameter("id", "identifier of the ip").DataType("string")).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Writes(v1.IPResponse{}).
		Returns(http.StatusOK, "OK", v1.IPResponse{}).
		DefaultReturns("Error", httperrors.HTTPErrorResponse{}))
}

func (ir ipResource) findIP(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("id")

	ip, err := ir.DS.FindIPByID(id)
	if helper.CheckError(request, response, utils.CurrentFuncName(), err) {
		return
	}
	err = response.WriteHeaderAndEntity(http.StatusOK, v1.NewIPResponse(ip))
	if err != nil {
		zapup.MustRootLogger().Error("Failed to send response", zap.Error(err))
		return
	}
}
