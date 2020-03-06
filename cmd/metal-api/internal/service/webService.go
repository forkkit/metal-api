package service

import (
	"fmt"
	"github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service/helper"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/utils"
	"github.com/metal-stack/metal-lib/httperrors"
	"net/http"
	"strings"
)

var (
	BasePath = "/"
)

type WebService struct {
	Version Version
	Path    string
	Tags    []string
	Routes  []*Route
}

func Build(webService *WebService) *restful.WebService {
	return webService.build()
}

func (s *WebService) build() *restful.WebService {
	ws := new(restful.WebService)

	basePath := BasePath
	if !strings.HasSuffix(basePath, "/") {
		basePath = basePath + "/"
	}

	path := s.Path
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	ws.
		Path(fmt.Sprintf("%s%s%s", basePath, s.Version, path)).
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	for _, route := range s.Routes {
		s.addRoute(route, ws)
	}

	return ws
}

func (s *WebService) addRoute(route *Route, ws *restful.WebService) {
	var rb *restful.RouteBuilder

	subPath := strings.TrimPrefix(route.SubPath, "/")

	switch route.Method {
	case http.MethodPost:
		rb = ws.POST(subPath)
	case http.MethodPut:
		rb = ws.PUT(subPath)
	case http.MethodDelete:
		rb = ws.DELETE(subPath)
	default:
		rb = ws.GET(subPath)
	}

	rb.
		Doc(route.Doc).
		Writes(route.Writes)

	if route.Error == nil {
		route.Error = httperrors.HTTPErrorResponse{}
	}
	rb.DefaultReturns("Error", route.Error)

	if route.Access == nil {
		rb.To(route.Handler)
	} else if metal.IsAdmin(route.Access) {
		rb.To(helper.Admin(route.Handler))
	} else if metal.IsEdit(route.Access) {
		rb.To(helper.Editor(route.Handler))
	} else {
		rb.To(helper.Viewer(route.Handler))
	}

	tags := make([]string, 0, len(s.Tags))
	copy(tags, s.Tags)
	if len(tags) == 0 {
		tags = []string{s.Path}
	}
	rb.Metadata(restfulspec.KeyOpenAPITags, tags)

	op := utils.GetFunctionName(route.Handler)
	rb.Operation(op)

	if route.Reads != nil {
		rb.Reads(route.Reads)
	}

	pp := append(route.PathParameters, route.PathParameter)
	for _, p := range pp {
		if p == nil {
			continue
		}
		rb.Param(ws.PathParameter(p.name, p.description).DataType("string"))
	}

	if len(route.Returns) == 0 {
		rb.Returns(http.StatusOK, "OK", route.Writes)
	}

	for _, ret := range route.Returns {
		rb.Returns(ret.Status, ret.Message, ret.Model)
	}

	ws.Route(rb)
}
