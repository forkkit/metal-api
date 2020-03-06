package network

import (
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service/helper"
	v1 "github.com/metal-stack/metal-api/cmd/metal-api/internal/service/v1"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/utils"
	"github.com/metal-stack/metal-lib/zapup"
	"go.uber.org/zap"
	"net/http"
)

func (r networkResource) updateNetwork(request *restful.Request, response *restful.Response) {
	var requestPayload v1.NetworkUpdateRequest
	err := request.ReadEntity(&requestPayload)
	if helper.CheckError(request, response, utils.CurrentFuncName(), err) {
		return
	}

	oldNetwork, err := r.DS.FindNetworkByID(requestPayload.ID)
	if helper.CheckError(request, response, utils.CurrentFuncName(), err) {
		return
	}

	newNetwork := *oldNetwork

	if requestPayload.Name != nil {
		newNetwork.Name = *requestPayload.Name
	}
	if requestPayload.Description != nil {
		newNetwork.Description = *requestPayload.Description
	}

	var prefixesToBeRemoved metal.Prefixes
	var prefixesToBeAdded metal.Prefixes

	if len(requestPayload.Prefixes) > 0 {
		var prefixesFromRequest metal.Prefixes
		for _, prefixCidr := range requestPayload.Prefixes {
			requestPrefix, err := metal.NewPrefixFromCIDR(prefixCidr)
			if err != nil {
				if helper.CheckError(request, response, utils.CurrentFuncName(), err) {
					return
				}
			}
			prefixesFromRequest = append(prefixesFromRequest, *requestPrefix)
		}
		newNetwork.Prefixes = prefixesFromRequest

		prefixesToBeRemoved = oldNetwork.SubstractPrefixes(prefixesFromRequest...)

		// now validate if there are ips which have a prefix to be removed as a parent
		allIPs, err := r.DS.ListIPs()
		if helper.CheckError(request, response, utils.CurrentFuncName(), err) {
			return
		}
		err = helper.CheckAnyIPOfPrefixesInUse(allIPs, prefixesToBeRemoved)
		if err != nil {
			if helper.CheckError(request, response, utils.CurrentFuncName(), fmt.Errorf("unable to update Network: %v", err)) {
				return
			}
		}

		prefixesToBeAdded = newNetwork.SubstractPrefixes(oldNetwork.Prefixes...)
	}

	for _, p := range prefixesToBeRemoved {
		err := r.ipamer.DeletePrefix(p)
		if helper.CheckError(request, response, utils.CurrentFuncName(), err) {
			return
		}
	}

	for _, p := range prefixesToBeAdded {
		err := r.ipamer.CreatePrefix(p)
		if helper.CheckError(request, response, utils.CurrentFuncName(), err) {
			return
		}
	}

	err = r.DS.UpdateNetwork(oldNetwork, &newNetwork)
	if helper.CheckError(request, response, utils.CurrentFuncName(), err) {
		return
	}

	usage := helper.GetNetworkUsage(&newNetwork, r.ipamer)
	err = response.WriteHeaderAndEntity(http.StatusOK, v1.NewNetworkResponse(&newNetwork, usage))
	if err != nil {
		zapup.MustRootLogger().Error("Failed to send response", zap.Error(err))
		return
	}
}
