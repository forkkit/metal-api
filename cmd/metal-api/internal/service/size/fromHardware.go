package size

import (
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service/helper"
	v1 "github.com/metal-stack/metal-api/pkg/proto/v1"
	"github.com/metal-stack/metal-api/pkg/util"
	"github.com/metal-stack/metal-lib/zapup"
	"go.uber.org/zap"
	"net/http"
)

type SizeConstraintMatchingLog struct {
	Constraint v1.SizeConstraint `json:"constraint" description:"the size constraint to which this log relates to"`
	Match      bool              `json:"match" description:"indicates whether the constraint matched or not"`
	Log        string            `json:"log" description:"a string represention of the matching condition"`
}

type SizeMatchingLog struct {
	Name        string                      `json:"name"`
	Log         string                      `json:"log"`
	Match       bool                        `json:"match"`
	Constraints []SizeConstraintMatchingLog `json:"constraints"`
}

func NewSizeMatchingLog(m *metal.SizeMatchingLog) *SizeMatchingLog {
	var constraints []SizeConstraintMatchingLog
	for i := range m.Constraints {
		constraint := SizeConstraintMatchingLog{
			Constraint: v1.SizeConstraint{
				Type: service.ToConstraintType(m.Constraints[i].Constraint.Type),
				Min:  m.Constraints[i].Constraint.Min,
				Max:  m.Constraints[i].Constraint.Max,
			},
			Match: m.Constraints[i].Match,
			Log:   m.Constraints[i].Log,
		}
		constraints = append(constraints, constraint)
	}
	return &SizeMatchingLog{
		Name:        m.Name,
		Match:       m.Match,
		Log:         m.Log,
		Constraints: constraints,
	}
}

func (r *sizeResource) fromHardware(request *restful.Request, response *restful.Response) {
	var requestPayload v1.MachineHardwareExtended
	err := request.ReadEntity(&requestPayload)
	if helper.CheckError(request, response, util.CurrentFuncName(), err) {
		return
	}

	hw := service.NewMetalMachineHardware(&requestPayload)
	_, lg, err := r.ds.FromHardware(hw)
	if helper.CheckError(request, response, util.CurrentFuncName(), err) {
		return
	}

	if len(lg) < 1 {
		if helper.CheckError(request, response, util.CurrentFuncName(), fmt.Errorf("size matching log is empty")) {
			return
		}
	}

	err = response.WriteHeaderAndEntity(http.StatusOK, NewSizeMatchingLog(lg[0]))
	if err != nil {
		zapup.MustRootLogger().Error("Failed to send response", zap.Error(err))
		return
	}
}
