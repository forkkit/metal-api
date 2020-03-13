package network

import (
	"github.com/emicklei/go-restful"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service/helper"
	v1 "github.com/metal-stack/metal-api/pkg/proto/v1"
	"github.com/metal-stack/metal-api/pkg/util"
	"github.com/metal-stack/metal-lib/zapup"
	"go.uber.org/zap"
	"net/http"
)

func (r *networkResource) findNetworks(request *restful.Request, response *restful.Response) {
	var requestPayload v1.NetworkSearchQuery
	err := request.ReadEntity(&requestPayload)
	if helper.CheckError(request, response, util.CurrentFuncName(), err) {
		return
	}

	var nws metal.Networks
	err = r.ds.SearchNetworks(&requestPayload, &nws)
	if helper.CheckError(request, response, util.CurrentFuncName(), err) {
		return
	}

	var result []*v1.NetworkResponse
	for _, nw := range nws {
		usage := GetNetworkUsage(&nw, r.ipamer)
		result = append(result, NewNetworkResponse(&nw, usage))
	}
	err = response.WriteHeaderAndEntity(http.StatusOK, result)
	if err != nil {
		zapup.MustRootLogger().Error("Failed to send response", zap.Error(err))
		return
	}
}
