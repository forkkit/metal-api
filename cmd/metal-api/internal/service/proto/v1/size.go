package v1

import (
	mdv1 "github.com/metal-stack/masterdata-api/api/v1"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service/helper"
)

func NewSizeResponse(s *metal.Size) *SizeResponse {
	if s == nil {
		return nil
	}

	return &SizeResponse{
		Size: ToSize(s),
	}
}

func ToSize(s *metal.Size) *Size {
	return &Size{
		Common: &Common{
			Meta: &mdv1.Meta{
				Id:          s.GetID(),
				Apiversion:  "v1",
				Version:     1,
				CreatedTime: helper.ToTimestamp(s.Created),
				UpdatedTime: helper.ToTimestamp(s.Changed),
			},
			Name:        helper.ToStringValue(s.Name),
			Description: helper.ToStringValue(s.Description),
		},
		Constraints: toConstraintSlice(s.Constraints...),
	}
}

func toConstraintSlice(constraints ...metal.Constraint) []*SizeConstraint {
	var cc []*SizeConstraint
	for _, c := range constraints {
		constraint := &SizeConstraint{
			Type: ToConstraintType(c.Type),
			Min:  c.Min,
			Max:  c.Max,
		}
		cc = append(cc, constraint)
	}
	return cc
}

func ToConstraintType(ct metal.ConstraintType) SizeConstraint_Type {
	switch ct {
	case metal.MemoryConstraint:
		return SizeConstraint_MEMORY
	case metal.StorageConstraint:
		return SizeConstraint_STORAGE
	}
	return SizeConstraint_CORES
}
