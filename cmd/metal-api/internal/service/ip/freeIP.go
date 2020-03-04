package ip

import (
	"fmt"
	"github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service/helper"
	v1 "github.com/metal-stack/metal-api/cmd/metal-api/internal/service/v1"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/utils"
	"github.com/metal-stack/metal-lib/httperrors"
	"github.com/metal-stack/metal-lib/zapup"
	"go.uber.org/zap"
	"net/http"
)

func (r ipResource) addFreeIPRoute(ws *restful.WebService, tags []string) {
	ws.Route(ws.POST("/free/{id}").
		To(helper.Editor(r.freeIP)).
		Operation("freeIP").
		Doc("frees an ip").
		Param(ws.PathParameter("id", "identifier of the ip").DataType("string")).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Writes(v1.IPResponse{}).
		Returns(http.StatusOK, "OK", v1.IPResponse{}).
		DefaultReturns("Error", httperrors.HTTPErrorResponse{}))
}

func (r ipResource) freeIP(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("id")

	ip, err := r.DS.FindIPByID(id)
	if helper.CheckError(request, response, utils.CurrentFuncName(), err) {
		return
	}

	err = validateIPDelete(ip)
	if helper.CheckError(request, response, utils.CurrentFuncName(), err) {
		return
	}

	err = r.ipamer.ReleaseIP(*ip)
	if helper.CheckError(request, response, utils.CurrentFuncName(), err) {
		return
	}

	err = r.DS.DeleteIP(ip)
	if helper.CheckError(request, response, utils.CurrentFuncName(), err) {
		return
	}
	err = response.WriteHeaderAndEntity(http.StatusOK, v1.NewIPResponse(ip))
	if err != nil {
		zapup.MustRootLogger().Error("Failed to send response", zap.Error(err))
		return
	}
}

func validateIPDelete(ip *metal.IP) error {
	s := ip.GetScope()
	if s != metal.ScopeProject && ip.Type == metal.Static {
		return fmt.Errorf("ip with scope %s can not be deleted", ip.GetScope())
	}
	return nil
}
