package image

import (
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service/helper"
	"github.com/metal-stack/metal-api/pkg/util"
	"github.com/metal-stack/metal-lib/zapup"
	"go.uber.org/zap"
	"net/http"
)

func (r *imageResource) deleteImage(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("id")

	img, err := r.ds.FindImage(id)
	if service.CheckError(request, response, util.CurrentFuncName(), err) {
		return
	}

	machines, err := r.ds.ListMachines()
	if service.CheckError(request, response, util.CurrentFuncName(), err) {
		return
	}
	for _, m := range machines {
		if m.Allocation == nil {
			continue
		}
		if m.Allocation.ImageID == img.ID {
			if service.CheckError(request, response, util.CurrentFuncName(), fmt.Errorf("image %s is in use by machine:%s", img.ID, m.ID)) {
				return
			}
		}
	}

	err = r.ds.DeleteImage(img)
	if service.CheckError(request, response, util.CurrentFuncName(), err) {
		return
	}
	err = response.WriteHeaderAndEntity(http.StatusOK, helper.NewImageResponse(img))
	if err != nil {
		zapup.MustRootLogger().Error("Failed to send response", zap.Error(err))
		return
	}
}
