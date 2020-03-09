package network

import (
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service/helper"
	v12 "github.com/metal-stack/metal-api/cmd/metal-api/internal/service/proto/v1"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/utils"
	"github.com/metal-stack/metal-lib/zapup"
	"go.uber.org/zap"
	"net/http"
)

func (r *networkResource) freeNetwork(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("id")

	nw, err := r.ds.FindNetworkByID(id)
	if helper.CheckError(request, response, utils.CurrentFuncName(), err) {
		return
	}

	for _, prefix := range nw.Prefixes {
		usage, err := r.ipamer.PrefixUsage(prefix.String())
		if helper.CheckError(request, response, utils.CurrentFuncName(), err) {
			return
		}
		if usage.UsedIPs > 2 {
			if helper.CheckError(request, response, utils.CurrentFuncName(), fmt.Errorf("cannot release child prefix %s because IPs in the prefix are still in use: %v", prefix.String(), usage.UsedIPs-2)) {
				return
			}
		}
	}

	for _, prefix := range nw.Prefixes {
		err = r.ipamer.ReleaseChildPrefix(prefix)
		if helper.CheckError(request, response, utils.CurrentFuncName(), err) {
			return
		}
	}

	if nw.Vrf != 0 {
		err = r.ds.ReleaseUniqueInteger(nw.Vrf)
		if err != nil {
			if helper.CheckError(request, response, utils.CurrentFuncName(), fmt.Errorf("could not release vrf: %v", err)) {
				return
			}
		}
	}

	err = r.ds.DeleteNetwork(nw)
	if helper.CheckError(request, response, utils.CurrentFuncName(), err) {
		return
	}
	err = response.WriteHeaderAndEntity(http.StatusOK, v12.NewNetworkResponse(nw, &metal.NetworkUsage{}))
	if err != nil {
		zapup.MustRootLogger().Error("Failed to send response", zap.Error(err))
		return
	}
}
