package ip

import (
	"github.com/emicklei/go-restful"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service/helper"
	"github.com/metal-stack/metal-api/pkg/util"
	"github.com/metal-stack/metal-lib/zapup"
	"go.uber.org/zap"
	"net/http"
)

func (r *ipResource) findIP(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("id")

	ip, err := r.ds.FindIPByID(id)
	if service.CheckError(request, response, util.CurrentFuncName(), err) {
		return
	}
	err = response.WriteHeaderAndEntity(http.StatusOK, helper.NewIPResponse(ip))
	if err != nil {
		zapup.MustRootLogger().Error("Failed to send response", zap.Error(err))
		return
	}
}