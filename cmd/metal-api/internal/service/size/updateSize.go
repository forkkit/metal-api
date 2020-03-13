package size

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

func (r *sizeResource) updateSize(request *restful.Request, response *restful.Response) {
	var requestPayload v1.SizeUpdateRequest
	err := request.ReadEntity(&requestPayload)
	if helper.CheckError(request, response, util.CurrentFuncName(), err) {
		return
	}

	size := requestPayload.Size

	oldSize, err := r.ds.FindSize(size.Common.Meta.Id)
	if helper.CheckError(request, response, util.CurrentFuncName(), err) {
		return
	}

	newSize := *oldSize
	newSize.Name = size.Common.Name.GetValue()
	newSize.Description = size.Common.Description.GetValue()
	var constraints []metal.Constraint
	if size.Constraints != nil {
		for _, c := range size.Constraints {
			constraint := metal.Constraint{
				Type: mapSizeConstraintType(c.Type),
				Min:  c.Min,
				Max:  c.Max,
			}
			constraints = append(constraints, constraint)
		}
		newSize.Constraints = constraints
	}

	err = r.ds.UpdateSize(oldSize, &newSize)
	if helper.CheckError(request, response, util.CurrentFuncName(), err) {
		return
	}
	err = response.WriteHeaderAndEntity(http.StatusOK, NewSizeResponse(&newSize))
	if err != nil {
		zapup.MustRootLogger().Error("Failed to send response", zap.Error(err))
		return
	}
}
