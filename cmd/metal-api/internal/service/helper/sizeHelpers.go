package helper

import (
	v12 "github.com/metal-stack/masterdata-api/api/v1"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
	"github.com/metal-stack/metal-api/pkg/proto/v1"
	"github.com/metal-stack/metal-api/pkg/util"
)

func ToConstraintType(ct metal.ConstraintType) v1.SizeConstraint_Type {
	switch ct {
	case metal.MemoryConstraint:
		return v1.SizeConstraint_MEMORY
	case metal.StorageConstraint:
		return v1.SizeConstraint_STORAGE
	}
	return v1.SizeConstraint_CORES
}

func NewSizeResponse(s *metal.Size) *v1.SizeResponse {
	if s == nil {
		return nil
	}

	return &v1.SizeResponse{
		Size: ToSize(s),
	}
}

func ToSize(s *metal.Size) *v1.Size {
	return &v1.Size{
		Common: &v1.Common{
			Meta: &v12.Meta{
				Id:          s.GetID(),
				Apiversion:  "v1",
				Version:     1,
				CreatedTime: util.TimestampProto(s.Created),
				UpdatedTime: util.TimestampProto(s.Changed),
			},
			Name:        util.StringProto(s.Name),
			Description: util.StringProto(s.Description),
		},
		Constraints: toConstraintSlice(s.Constraints...),
	}
}

func toConstraintSlice(constraints ...metal.Constraint) []*v1.SizeConstraint {
	var cc []*v1.SizeConstraint
	for _, c := range constraints {
		constraint := &v1.SizeConstraint{
			Type: ToConstraintType(c.Type),
			Min:  c.Min,
			Max:  c.Max,
		}
		cc = append(cc, constraint)
	}
	return cc
}
