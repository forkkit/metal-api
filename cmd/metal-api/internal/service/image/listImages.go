package image

import (
	"github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service/helper"
	v1 "github.com/metal-stack/metal-api/cmd/metal-api/internal/service/v1"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/utils"
	"github.com/metal-stack/metal-lib/httperrors"
	"github.com/metal-stack/metal-lib/zapup"
	"go.uber.org/zap"
	"net/http"
)

func (r imageResource) addListImagesRoute(ws *restful.WebService, tags []string) {
	ws.Route(ws.GET("/").
		To(r.listImages).
		Operation("listImages").
		Doc("get all images").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Writes([]v1.ImageResponse{}).
		Returns(http.StatusOK, "OK", []v1.ImageResponse{}).
		DefaultReturns("Error", httperrors.HTTPErrorResponse{}))
}

func (r imageResource) listImages(request *restful.Request, response *restful.Response) {
	imgs, err := r.DS.ListImages()
	if helper.CheckError(request, response, utils.CurrentFuncName(), err) {
		return
	}

	result := []*v1.ImageResponse{}
	for i := range imgs {
		result = append(result, v1.NewImageResponse(&imgs[i]))
	}
	err = response.WriteHeaderAndEntity(http.StatusOK, result)
	if err != nil {
		zapup.MustRootLogger().Error("Failed to send response", zap.Error(err))
		return
	}
}
