package v1

import (
	mdv1 "github.com/metal-stack/masterdata-api/api/v1"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service/helper"
)

func NewImageResponse(img *metal.Image) *ImageResponse {
	if img == nil {
		return nil
	}
	return &ImageResponse{
		Image: ToImage(img),
	}
}

func ToImage(img *metal.Image) *Image {
	var features []string
	for k, v := range img.Features {
		if v {
			features = append(features, string(k))
		}
	}

	return &Image{
		Common: &Common{
			Meta: &mdv1.Meta{
				Id:          img.ID,
				Apiversion:  "v1",
				Version:     1,
				CreatedTime: helper.ToTimestamp(img.Created),
				UpdatedTime: helper.ToTimestamp(img.Changed),
			},
			Name:        helper.ToStringValue(img.Name),
			Description: helper.ToStringValue(img.Description),
		},
		URL:      helper.ToStringValue(img.URL),
		Features: helper.ToStringValueSlice(features...),
	}
}
