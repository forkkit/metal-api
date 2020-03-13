package network

import (
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service/helper"
	v1 "github.com/metal-stack/metal-api/pkg/proto/v1"
	"github.com/metal-stack/metal-api/pkg/util"
	"github.com/metal-stack/metal-lib/zapup"
	"go.uber.org/zap"
	"net/http"
)

func (r *networkResource) updateNetwork(request *restful.Request, response *restful.Response) {
	var requestPayload v1.NetworkUpdateRequest
	err := request.ReadEntity(&requestPayload)
	if helper.CheckError(request, response, util.CurrentFuncName(), err) {
		return
	}

	oldNetwork, err := r.ds.FindNetworkByID(requestPayload.Common.Meta.Id)
	if helper.CheckError(request, response, util.CurrentFuncName(), err) {
		return
	}

	newNetwork := *oldNetwork

	newNetwork.Name = requestPayload.Common.Name.GetValue()
	newNetwork.Description = requestPayload.Common.Description.GetValue()

	var prefixesToBeRemoved metal.Prefixes
	var prefixesToBeAdded metal.Prefixes

	if len(requestPayload.Prefixes) > 0 {
		var prefixesFromRequest metal.Prefixes
		for _, prefixCidr := range requestPayload.Prefixes {
			requestPrefix, err := metal.NewPrefixFromCIDR(prefixCidr)
			if err != nil {
				if helper.CheckError(request, response, util.CurrentFuncName(), err) {
					return
				}
			}
			prefixesFromRequest = append(prefixesFromRequest, *requestPrefix)
		}
		newNetwork.Prefixes = prefixesFromRequest

		prefixesToBeRemoved = oldNetwork.SubstractPrefixes(prefixesFromRequest...)

		// now validate if there are ips which have a prefix to be removed as a parent
		allIPs, err := r.ds.ListIPs()
		if helper.CheckError(request, response, util.CurrentFuncName(), err) {
			return
		}
		err = CheckAnyIPOfPrefixesInUse(allIPs, prefixesToBeRemoved)
		if err != nil {
			if helper.CheckError(request, response, util.CurrentFuncName(), fmt.Errorf("unable to update Network: %v", err)) {
				return
			}
		}

		prefixesToBeAdded = newNetwork.SubstractPrefixes(oldNetwork.Prefixes...)
	}

	for _, p := range prefixesToBeRemoved {
		err := r.ipamer.DeletePrefix(p)
		if helper.CheckError(request, response, util.CurrentFuncName(), err) {
			return
		}
	}

	for _, p := range prefixesToBeAdded {
		err := r.ipamer.CreatePrefix(p)
		if helper.CheckError(request, response, util.CurrentFuncName(), err) {
			return
		}
	}

	err = r.ds.UpdateNetwork(oldNetwork, &newNetwork)
	if helper.CheckError(request, response, util.CurrentFuncName(), err) {
		return
	}

	usage := GetNetworkUsage(&newNetwork, r.ipamer)
	err = response.WriteHeaderAndEntity(http.StatusOK, NewNetworkResponse(&newNetwork, usage))
	if err != nil {
		zapup.MustRootLogger().Error("Failed to send response", zap.Error(err))
		return
	}
}
