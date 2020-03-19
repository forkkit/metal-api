package machine

import (
	"github.com/emicklei/go-restful"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/datastore"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service"
	v1 "github.com/metal-stack/metal-api/pkg/proto/v1"
	"github.com/metal-stack/metal-api/pkg/util"
	"github.com/metal-stack/metal-lib/zapup"
	"go.uber.org/zap"
	"net/http"
)

func (r *machineResource) findIPMIMachine(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("id")
	resp, err := FindIPMIMachine(r.ds, id)
	if service.CheckError(request, response, util.CurrentFuncName(), err) {
		return
	}
	err = response.WriteHeaderAndEntity(http.StatusOK, resp)
	if err != nil {
		zapup.MustRootLogger().Error("Failed to send response", zap.Error(err))
	}
}

func FindIPMIMachine(ds *datastore.RethinkStore, id string) (*v1.MachineIPMIResponse, error) {
	m, err := ds.FindMachineByID(id)
	if err != nil {
		return nil, err
	}
	resp := makeMachineIPMIResponse(m, ds, zapup.MustRootLogger().Sugar())
	return resp, nil
}

func makeMachineIPMIResponse(m *metal.Machine, ds *datastore.RethinkStore, logger *zap.SugaredLogger) *v1.MachineIPMIResponse {
	s, p, i, ec := FindMachineReferencedEntities(m, ds, logger)
	return NewMachineIPMIResponse(m, s, p, i, ec)
}
