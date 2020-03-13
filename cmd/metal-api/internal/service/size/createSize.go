package size

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

func (r *sizeResource) createSize(request *restful.Request, response *restful.Response) {
	var requestPayload v1.SizeCreateRequest
	err := request.ReadEntity(&requestPayload)
	if helper.CheckError(request, response, util.CurrentFuncName(), err) {
		return
	}

	size := requestPayload.Size

	if size.Common.Meta.Id == "" {
		if helper.CheckError(request, response, util.CurrentFuncName(), fmt.Errorf("id should not be empty")) {
			return
		}
	}

	if size.Common.Meta.Id == metal.UnknownSize.GetID() {
		if helper.CheckError(request, response, util.CurrentFuncName(), fmt.Errorf("id cannot be %q", metal.UnknownSize.GetID())) {
			return
		}
	}

	var constraints []metal.Constraint
	for _, c := range size.Constraints {
		constraint := metal.Constraint{
			Type: mapSizeConstraintType(c.Type),
			Min:  c.Min,
			Max:  c.Max,
		}
		constraints = append(constraints, constraint)
	}

	s := &metal.Size{
		Base: metal.Base{
			ID:          size.Common.Meta.Id,
			Name:        size.Common.Name.GetValue(),
			Description: size.Common.Description.GetValue(),
		},
		Constraints: constraints,
	}

	err = r.ds.CreateSize(s)
	if helper.CheckError(request, response, util.CurrentFuncName(), err) {
		return
	}
	err = response.WriteHeaderAndEntity(http.StatusCreated, NewSizeResponse(s))
	if err != nil {
		zapup.MustRootLogger().Error("Failed to send response", zap.Error(err))
		return
	}
}
