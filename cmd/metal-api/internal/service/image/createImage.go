package image

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

func (r *imageResource) createImage(request *restful.Request, response *restful.Response) {
	var requestPayload v1.ImageCreateRequest
	err := request.ReadEntity(&requestPayload)
	if helper.CheckError(request, response, util.CurrentFuncName(), err) {
		return
	}
	image := requestPayload.Image

	if image.Common.Meta.Id == "" {
		if helper.CheckError(request, response, util.CurrentFuncName(), fmt.Errorf("id should not be empty")) {
			return
		}
	}

	if image.URL.GetValue() == "" {
		if helper.CheckError(request, response, util.CurrentFuncName(), fmt.Errorf("url should not be empty")) {
			return
		}
	}

	features := make(map[metal.ImageFeatureType]bool)
	for _, f := range image.Features {
		ft, err := metal.ImageFeatureTypeFrom(f.GetValue())
		if helper.CheckError(request, response, util.CurrentFuncName(), err) {
			return
		}
		features[ft] = true
	}

	img := &metal.Image{
		Base: metal.Base{
			ID:          image.Common.Meta.Id,
			Name:        image.Common.Name.GetValue(),
			Description: image.Common.Description.GetValue(),
		},
		URL:      image.URL.GetValue(),
		Features: features,
	}

	err = r.ds.CreateImage(img)
	if helper.CheckError(request, response, util.CurrentFuncName(), err) {
		return
	}
	err = response.WriteHeaderAndEntity(http.StatusCreated, service.NewImageResponse(img))
	if err != nil {
		zapup.MustRootLogger().Error("Failed to send response", zap.Error(err))
		return
	}
}
