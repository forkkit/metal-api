package image

import (
	"github.com/emicklei/go-restful"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service/helper"
	"github.com/metal-stack/metal-api/pkg/helper"
	v1 "github.com/metal-stack/metal-api/pkg/proto"
	"github.com/metal-stack/metal-lib/zapup"
	"go.uber.org/zap"
	"net/http"
)

func (r *imageResource) listImages(request *restful.Request, response *restful.Response) {
	imgs, err := r.ds.ListImages()
	if helper.CheckError(request, response, helper.CurrentFuncName(), err) {
		return
	}

	result := []*v1.ImageResponse{}
	for i := range imgs {
		result = append(result, service.NewImageResponse(&imgs[i]))
	}
	err = response.WriteHeaderAndEntity(http.StatusOK, result)
	if err != nil {
		zapup.MustRootLogger().Error("Failed to send response", zap.Error(err))
		return
	}
}
