package firewall

import (
	"github.com/emicklei/go-restful"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/datastore"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service/helper"
	"github.com/metal-stack/metal-api/pkg/helper"
	"github.com/metal-stack/metal-lib/zapup"
	"go.uber.org/zap"
	"net/http"
)

func (r *firewallResource) findFirewalls(request *restful.Request, response *restful.Response) {
	var requestPayload datastore.MachineSearchQuery
	err := request.ReadEntity(&requestPayload)
	if helper.CheckError(request, response, helper.CurrentFuncName(), err) {
		return
	}

	var possibleFws metal.Machines
	err = r.ds.SearchMachines(&requestPayload, &possibleFws)
	if helper.CheckError(request, response, helper.CurrentFuncName(), err) {
		return
	}

	imgs, err := r.ds.ListImages()
	if helper.CheckError(request, response, helper.CurrentFuncName(), err) {
		return
	}

	fws := metal.Machines{}
	imageMap := imgs.ByID()
	for i := range possibleFws {
		if possibleFws[i].IsFirewall(imageMap) {
			fws = append(fws, possibleFws[i])
		}
	}

	err = response.WriteHeaderAndEntity(http.StatusOK, makeFirewallResponseList(fws, r.ds, helper.Logger(request).Sugar()))
	if err != nil {
		zapup.MustRootLogger().Error("Failed to send response", zap.Error(err))
		return
	}
}