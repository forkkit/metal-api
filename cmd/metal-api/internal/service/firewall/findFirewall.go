package firewall

import (
	"fmt"
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

func (r firewallResource) addFindFirewallRoute(ws *restful.WebService, tags []string) {
	ws.Route(ws.GET("/{id}").
		To(helper.Viewer(r.findFirewall)).
		Operation("findFirewall").
		Doc("get firewall by id").
		Param(ws.PathParameter("id", "identifier of the firewall").DataType("string")).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Writes(v1.FirewallResponse{}).
		Returns(http.StatusOK, "OK", v1.FirewallResponse{}).
		DefaultReturns("Error", httperrors.HTTPErrorResponse{}))
}

func (r firewallResource) findFirewall(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("id")

	fw, err := r.DS.FindMachineByID(id)
	if helper.CheckError(request, response, utils.CurrentFuncName(), err) {
		return
	}

	imgs, err := r.DS.ListImages()
	if helper.CheckError(request, response, utils.CurrentFuncName(), err) {
		return
	}

	if !fw.IsFirewall(imgs.ByID()) {
		helper.SendError(utils.Logger(request), response, utils.CurrentFuncName(), httperrors.NotFound(fmt.Errorf("machine is not a firewall")))
		return
	}

	err = response.WriteHeaderAndEntity(http.StatusOK, makeFirewallResponse(fw, r.DS, utils.Logger(request).Sugar()))
	if err != nil {
		zapup.MustRootLogger().Error("Failed to send response", zap.Error(err))
		return
	}
}

func makeFirewallResponse(fw *metal.Machine, ds *datastore.RethinkStore, logger *zap.SugaredLogger) *v1.FirewallResponse {
	return &v1.FirewallResponse{MachineResponse: *helper.MakeMachineResponse(fw, ds, logger)}
}

func makeFirewallResponseList(fws metal.Machines, ds *datastore.RethinkStore, logger *zap.SugaredLogger) []*v1.FirewallResponse {
	machineResponseList := helper.MakeMachineResponseList(fws, ds, logger)

	firewallResponseList := []*v1.FirewallResponse{}
	for i := range machineResponseList {
		firewallResponseList = append(firewallResponseList, &v1.FirewallResponse{MachineResponse: *machineResponseList[i]})
	}

	return firewallResponseList

}
