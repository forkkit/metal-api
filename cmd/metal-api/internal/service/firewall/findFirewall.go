package firewall

import (
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/datastore"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service/helper"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service/machine"
	v1 "github.com/metal-stack/metal-api/pkg/proto/v1"
	"github.com/metal-stack/metal-api/pkg/util"
	"github.com/metal-stack/metal-lib/httperrors"
	"github.com/metal-stack/metal-lib/zapup"
	"go.uber.org/zap"
	"net/http"
)

func (r *firewallResource) findFirewall(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("id")

	fw, err := r.ds.FindMachineByID(id)
	if service.CheckError(request, response, util.CurrentFuncName(), err) {
		return
	}

	imgs, err := r.ds.ListImages()
	if service.CheckError(request, response, util.CurrentFuncName(), err) {
		return
	}

	if !fw.IsFirewall(imgs.ByID()) {
		service.SendError(util.Logger(request), response, util.CurrentFuncName(), httperrors.NotFound(fmt.Errorf("machine is not a firewall")))
		return
	}

	err = response.WriteHeaderAndEntity(http.StatusOK, makeFirewallResponse(fw, r.ds, util.Logger(request).Sugar()))
	if err != nil {
		zapup.MustRootLogger().Error("Failed to send response", zap.Error(err))
		return
	}
}

func makeFirewallResponse(fw *metal.Machine, ds *datastore.RethinkStore, logger *zap.SugaredLogger) *v1.FirewallResponse {
	return &v1.FirewallResponse{MachineResponse: machine.MakeResponse(fw, ds, logger)}
}

func makeFirewallResponseList(fws metal.Machines, ds *datastore.RethinkStore, logger *zap.SugaredLogger) []*v1.FirewallResponse {
	machineResponseList := helper.MakeMachineResponseList(fws, ds, logger)

	firewallResponseList := []*v1.FirewallResponse{}
	for i := range machineResponseList {
		firewallResponseList = append(firewallResponseList, &v1.FirewallResponse{MachineResponse: machineResponseList[i]})
	}

	return firewallResponseList

}
