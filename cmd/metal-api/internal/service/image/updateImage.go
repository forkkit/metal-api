package image

import (
	"github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service/helper"
	v1 "github.com/metal-stack/metal-api/cmd/metal-api/internal/service/v1"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/utils"
	"github.com/metal-stack/metal-lib/httperrors"
	"github.com/metal-stack/metal-lib/zapup"
	"go.uber.org/zap"
	"net/http"
)

func (r imageResource) addUpdateImageRoute(ws *restful.WebService, tags []string) {
	ws.Route(ws.POST("/").
		To(helper.Admin(r.updateImage)).
		Operation("updateImage").
		Doc("updates an image. if the image was changed since this one was read, a conflict is returned").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Reads(v1.ImageUpdateRequest{}).
		Returns(http.StatusOK, "OK", v1.ImageResponse{}).
		Returns(http.StatusConflict, "Conflict", httperrors.HTTPErrorResponse{}).
		DefaultReturns("Error", httperrors.HTTPErrorResponse{}))
}

func (r imageResource) updateImage(request *restful.Request, response *restful.Response) {
	var requestPayload v1.ImageUpdateRequest
	err := request.ReadEntity(&requestPayload)
	if helper.CheckError(request, response, utils.CurrentFuncName(), err) {
		return
	}

	oldImage, err := r.DS.FindImage(requestPayload.ID)
	if helper.CheckError(request, response, utils.CurrentFuncName(), err) {
		return
	}

	newImage := *oldImage

	if requestPayload.Name != nil {
		newImage.Name = *requestPayload.Name
	}
	if requestPayload.Description != nil {
		newImage.Description = *requestPayload.Description
	}
	if requestPayload.URL != nil {
		newImage.URL = *requestPayload.URL
	}
	features := make(map[metal.ImageFeatureType]bool)
	for _, f := range requestPayload.Features {
		ft, err := metal.ImageFeatureTypeFrom(f)
		if helper.CheckError(request, response, utils.CurrentFuncName(), err) {
			return
		}
		features[ft] = true
	}
	if len(features) > 0 {
		newImage.Features = features
	}

	err = r.DS.UpdateImage(oldImage, &newImage)
	if helper.CheckError(request, response, utils.CurrentFuncName(), err) {
		return
	}
	err = response.WriteHeaderAndEntity(http.StatusOK, v1.NewImageResponse(&newImage))
	if err != nil {
		zapup.MustRootLogger().Error("Failed to send response", zap.Error(err))
		return
	}
}
