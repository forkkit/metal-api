package helper

import (
	v12 "github.com/metal-stack/masterdata-api/api/v1"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
	"github.com/metal-stack/metal-api/pkg/proto/v1"
	"github.com/metal-stack/metal-api/pkg/util"
)

func NewImageResponse(img *metal.Image) *v1.ImageResponse {
	if img == nil {
		return nil
	}
	return &v1.ImageResponse{
		Image: ToImage(img),
	}
}

func ToImage(img *metal.Image) *v1.Image {
	var features []string
	for k, v := range img.Features {
		if v {
			features = append(features, string(k))
		}
	}

	return &v1.Image{
		Common: &v1.Common{
			Meta: &v12.Meta{
				Id:          img.ID,
				Apiversion:  "v1",
				Version:     1,
				CreatedTime: util.TimestampProto(img.Created),
				UpdatedTime: util.TimestampProto(img.Changed),
			},
			Name:        util.StringProto(img.Name),
			Description: util.StringProto(img.Description),
		},
		URL:      util.StringProto(img.URL),
		Features: util.StringSliceProto(features...),
	}
}
