package network

import (
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/datastore"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service/helper"
	v1 "github.com/metal-stack/metal-api/cmd/metal-api/internal/service/v1"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/utils"
	"github.com/metal-stack/metal-lib/zapup"
	"go.uber.org/zap"
	"net/http"
)

func (r networkResource) deleteNetwork(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("id")

	nw, err := r.DS.FindNetworkByID(id)
	if helper.CheckError(request, response, utils.CurrentFuncName(), err) {
		return
	}

	var children metal.Networks
	err = r.DS.SearchNetworks(&datastore.NetworkSearchQuery{ParentNetworkID: &nw.ID}, &children)
	if helper.CheckError(request, response, utils.CurrentFuncName(), err) {
		return
	}

	if len(children) != 0 {
		if helper.CheckError(request, response, utils.CurrentFuncName(), fmt.Errorf("network cannot be deleted because there are children of this network")) {
			return
		}
	}

	allIPs, err := r.DS.ListIPs()
	if helper.CheckError(request, response, utils.CurrentFuncName(), err) {
		return
	}

	err = helper.CheckAnyIPOfPrefixesInUse(allIPs, nw.Prefixes)
	if err != nil {
		if helper.CheckError(request, response, utils.CurrentFuncName(), fmt.Errorf("unable to delete Network: %v", err)) {
			return
		}
	}

	for _, p := range nw.Prefixes {
		err := r.ipamer.DeletePrefix(p)
		if helper.CheckError(request, response, utils.CurrentFuncName(), err) {
			return
		}
	}

	if nw.Vrf != 0 {
		err = r.DS.ReleaseUniqueInteger(nw.Vrf)
		if err != nil {
			if helper.CheckError(request, response, utils.CurrentFuncName(), fmt.Errorf("could not release vrf: %v", err)) {
				return
			}
		}
	}

	err = r.DS.DeleteNetwork(nw)
	if helper.CheckError(request, response, utils.CurrentFuncName(), err) {
		return
	}
	err = response.WriteHeaderAndEntity(http.StatusOK, v1.NewNetworkResponse(nw, &metal.NetworkUsage{}))
	if err != nil {
		zapup.MustRootLogger().Error("Failed to send response", zap.Error(err))
		return
	}
}
