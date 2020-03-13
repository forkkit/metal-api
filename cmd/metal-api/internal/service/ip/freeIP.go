package ip

import (
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service/helper"
	"github.com/metal-stack/metal-api/pkg/util"
	"github.com/metal-stack/metal-lib/zapup"
	"go.uber.org/zap"
	"net/http"
)

func (r *ipResource) freeIP(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("id")

	ip, err := r.ds.FindIPByID(id)
	if helper.CheckError(request, response, util.CurrentFuncName(), err) {
		return
	}

	err = validateIPDelete(ip)
	if helper.CheckError(request, response, util.CurrentFuncName(), err) {
		return
	}

	err = r.ipamer.ReleaseIP(*ip)
	if helper.CheckError(request, response, util.CurrentFuncName(), err) {
		return
	}

	err = r.ds.DeleteIP(ip)
	if helper.CheckError(request, response, util.CurrentFuncName(), err) {
		return
	}
	err = response.WriteHeaderAndEntity(http.StatusOK, NewIPResponse(ip))
	if err != nil {
		zapup.MustRootLogger().Error("Failed to send response", zap.Error(err))
		return
	}
}

func validateIPDelete(ip *metal.IP) error {
	s := ip.GetScope()
	if s != metal.ScopeProject && ip.Type == metal.Static {
		return fmt.Errorf("ip with scope %s can not be deleted", ip.GetScope())
	}
	return nil
}
