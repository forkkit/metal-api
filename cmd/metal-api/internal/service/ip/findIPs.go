package ip

import (
	"github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/datastore"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service/helper"
	v1 "github.com/metal-stack/metal-api/cmd/metal-api/internal/service/v1"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/utils"
	"github.com/metal-stack/metal-lib/httperrors"
	"github.com/metal-stack/metal-lib/zapup"
	"go.uber.org/zap"
	"net/http"
)

func (r ipResource) addFindIPsRoute(ws *restful.WebService, tags []string) {
	ws.Route(ws.POST("/find").
		To(helper.Viewer(r.findIPs)).
		Operation("findIPs").
		Doc("get all ips that match given properties").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Reads(v1.IPFindRequest{}).
		Writes([]v1.IPResponse{}).
		Returns(http.StatusOK, "OK", []v1.IPResponse{}).
		DefaultReturns("Error", httperrors.HTTPErrorResponse{}))
}

func (ir ipResource) findIPs(request *restful.Request, response *restful.Response) {
	var requestPayload datastore.IPSearchQuery
	err := request.ReadEntity(&requestPayload)
	if helper.CheckError(request, response, utils.CurrentFuncName(), err) {
		return
	}

	var ips metal.IPs
	err = ir.DS.SearchIPs(&requestPayload, &ips)
	if helper.CheckError(request, response, utils.CurrentFuncName(), err) {
		return
	}

	var result []*v1.IPResponse
	for i := range ips {
		result = append(result, v1.NewIPResponse(&ips[i]))
	}
	err = response.WriteHeaderAndEntity(http.StatusOK, result)
	if err != nil {
		zapup.MustRootLogger().Error("Failed to send response", zap.Error(err))
		return
	}
}
