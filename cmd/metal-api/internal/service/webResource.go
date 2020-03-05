package service

import (
	"fmt"
	"github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/datastore"
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

type WebResource struct {
	DS      *datastore.RethinkStore
	Version Version
	Path    string
	Tags    []string
	Routes  []Route

	ws *restful.WebService
}

func Build(r WebResource) *restful.WebService {
	return r.build()
}

func (r WebResource) build() *restful.WebService {
	if r.ws == nil {
		r.ws = new(restful.WebService)
	}

	path := r.Path
	if strings.HasPrefix(path, "/") {
		path = path[:len(path)-1]
	}

	basePath := BasePath
	if strings.HasPrefix(basePath, "/") {
		basePath = basePath[:len(basePath)-1]
	}

	r.ws.
		Path(fmt.Sprintf("%s/%s/%s", basePath, r.Version, path)).
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	for _, route := range r.Routes {
		r.addRoute(route)
	}

	return r.ws
}

func (r WebResource) addRoute(route Route) {
	if r.ws == nil {
		r.ws = new(restful.WebService)
	}

	if route.Error == nil {
		route.Error = httperrors.HTTPErrorResponse{}
	}

	var rb *restful.RouteBuilder
	switch route.Method {
	case http.MethodPost:
		rb = r.ws.POST(route.SubPath)
	case http.MethodPut:
		rb = r.ws.PUT(route.SubPath)
	case http.MethodDelete:
		rb = r.ws.DELETE(route.SubPath)
	default:
		rb = r.ws.GET(route.SubPath)
	}

	rb.Doc(route.Doc).
		Writes(route.Writes).
		DefaultReturns("Error", route.Error)

	if route.Access == nil {
		rb.To(route.Handler)
	} else if metal.IsAdmin(route.Access) {
		rb.To(helper.Admin(route.Handler))
	} else if metal.IsEdit(route.Access) {
		rb.To(helper.Editor(route.Handler))
	} else {
		rb.To(helper.Viewer(route.Handler))
	}

	tags := make([]string, 0, len(r.Tags))
	copy(tags, r.Tags)
	if len(tags) == 0 {
		tags = []string{r.Path}
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
		rb.Param(r.ws.PathParameter(p.name, p.description).DataType("string"))
	}

	if len(route.Returns) == 0 {
		rb.Returns(http.StatusOK, "OK", route.Writes)
	}

	for _, ret := range route.Returns {
		rb.Returns(ret.Status, ret.Message, ret.Model)
	}

	r.ws.Route(rb)
}
