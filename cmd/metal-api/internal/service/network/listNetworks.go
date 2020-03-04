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

func (r networkResource) addListNetworksRoute(ws *restful.WebService, tags []string) {
	ws.Route(ws.GET("/").
		To(helper.Viewer(r.listNetworks)).
		Operation("listNetworks").
		Doc("get all networks").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Writes([]v1.NetworkResponse{}).
		Returns(http.StatusOK, "OK", []v1.NetworkResponse{}).
		DefaultReturns("Error", httperrors.HTTPErrorResponse{}))
}

func (r networkResource) listNetworks(request *restful.Request, response *restful.Response) {
	nws, err := r.DS.ListNetworks()
	if helper.CheckError(request, response, utils.CurrentFuncName(), err) {
		return
	}

	var result []*v1.NetworkResponse
	for i := range nws {
		usage := helper.GetNetworkUsage(&nws[i], r.ipamer)
		result = append(result, v1.NewNetworkResponse(&nws[i], usage))
	}
	err = response.WriteHeaderAndEntity(http.StatusOK, result)
	if err != nil {
		zapup.MustRootLogger().Error("Failed to send response", zap.Error(err))
		return
	}
}
