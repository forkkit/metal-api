package image

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

func (r *imageResource) updateImage(request *restful.Request, response *restful.Response) {
	var requestPayload v1.ImageUpdateRequest
	err := request.ReadEntity(&requestPayload)
	if helper.CheckError(request, response, util.CurrentFuncName(), err) {
		return
	}

	img := requestPayload.Image

	oldImage, err := r.ds.FindImage(img.Common.Meta.Id)
	if helper.CheckError(request, response, util.CurrentFuncName(), err) {
		return
	}

	newImage := *oldImage
	newImage.Name = img.Common.Name.GetValue()
	newImage.Description = img.Common.Description.GetValue()

	if img.URL != nil {
		newImage.URL = img.URL.GetValue()
	}
	features := make(map[metal.ImageFeatureType]bool)
	for _, f := range img.Features {
		ft, err := metal.ImageFeatureTypeFrom(f.GetValue())
		if helper.CheckError(request, response, util.CurrentFuncName(), err) {
			return
		}
		features[ft] = true
	}
	if len(features) > 0 {
		newImage.Features = features
	}

	err = r.ds.UpdateImage(oldImage, &newImage)
	if helper.CheckError(request, response, util.CurrentFuncName(), err) {
		return
	}
	err = response.WriteHeaderAndEntity(http.StatusOK, NewImageResponse(&newImage))
	if err != nil {
		zapup.MustRootLogger().Error("Failed to send response", zap.Error(err))
		return
	}
}
