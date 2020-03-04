package project

import (
	"context"
	"github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
	mdmv1 "github.com/metal-stack/masterdata-api/api/v1"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service/helper"
	v1 "github.com/metal-stack/metal-api/cmd/metal-api/internal/service/v1"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/utils"
	"github.com/metal-stack/metal-lib/httperrors"
	"github.com/metal-stack/metal-lib/zapup"
	"go.uber.org/zap"
	"net/http"
)

func (r projectResource) addFindProjectsRoute(ws *restful.WebService, tags []string) {
	ws.Route(ws.POST("/find").
		To(helper.Viewer(r.findProjects)).
		Operation("findProjects").
		Doc("get all projects that match given properties").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Reads(v1.ProjectFindRequest{}).
		Writes([]v1.ProjectResponse{}).
		Returns(http.StatusOK, "OK", []v1.ProjectResponse{}).
		DefaultReturns("Error", httperrors.HTTPErrorResponse{}))
}

func (r projectResource) findProjects(request *restful.Request, response *restful.Response) {
	var requestPayload mdmv1.ProjectFindRequest
	err := request.ReadEntity(&requestPayload)
	if helper.CheckError(request, response, utils.CurrentFuncName(), err) {
		return
	}

	ps, err := r.mdc.Project().Find(context.Background(), &requestPayload)
	if helper.CheckError(request, response, utils.CurrentFuncName(), err) {
		return
	}

	err = response.WriteHeaderAndEntity(http.StatusOK, ps.Projects)
	if err != nil {
		zapup.MustRootLogger().Error("Failed to send response", zap.Error(err))
		return
	}
}
