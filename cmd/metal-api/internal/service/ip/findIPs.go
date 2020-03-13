package ip

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

func (r *ipResource) findIPs(request *restful.Request, response *restful.Response) {
	var requestPayload v1.IPFindRequest
	err := request.ReadEntity(&requestPayload)
	if helper.CheckError(request, response, util.CurrentFuncName(), err) {
		return
	}

	var ips metal.IPs
	err = r.ds.SearchIPs(&requestPayload, &ips)
	if helper.CheckError(request, response, util.CurrentFuncName(), err) {
		return
	}

	var result []*v1.IPResponse
	for i := range ips {
		result = append(result, NewIPResponse(&ips[i]))
	}
	err = response.WriteHeaderAndEntity(http.StatusOK, result)
	if err != nil {
		zapup.MustRootLogger().Error("Failed to send response", zap.Error(err))
		return
	}
}
