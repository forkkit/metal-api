package size

import (
	"fmt"
	"github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service/helper"
	v1 "github.com/metal-stack/metal-api/cmd/metal-api/internal/service/v1"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/utils"
	"github.com/metal-stack/metal-lib/httperrors"
	"github.com/metal-stack/metal-lib/zapup"
	"go.uber.org/zap"
	"net/http"
)

func (r sizeResource) addFromHardwareRoute(ws *restful.WebService, tags []string) {
	ws.Route(ws.POST("/from-hardware").
		To(r.fromHardware).
		Operation("fromHardware").
		Doc("Searches all sizes for one to match the given hardwarespecs. If nothing is found, a list of entries is returned which describe the constraint which did not match").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Reads(v1.MachineHardwareExtended{}).
		Returns(http.StatusOK, "OK", v1.SizeMatchingLog{}).
		DefaultReturns("Error", httperrors.HTTPErrorResponse{}))
}

func (r sizeResource) fromHardware(request *restful.Request, response *restful.Response) {
	var requestPayload v1.MachineHardwareExtended
	err := request.ReadEntity(&requestPayload)
	if helper.CheckError(request, response, utils.CurrentFuncName(), err) {
		return
	}

	hw := v1.NewMetalMachineHardware(&requestPayload)
	_, lg, err := r.DS.FromHardware(hw)
	if helper.CheckError(request, response, utils.CurrentFuncName(), err) {
		return
	}

	if len(lg) < 1 {
		if helper.CheckError(request, response, utils.CurrentFuncName(), fmt.Errorf("size matching log is empty")) {
			return
		}
	}

	err = response.WriteHeaderAndEntity(http.StatusOK, v1.NewSizeMatchingLog(lg[0]))
	if err != nil {
		zapup.MustRootLogger().Error("Failed to send response", zap.Error(err))
		return
	}
}
