package machine

import (
	"github.com/emicklei/go-restful"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/datastore"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service"
	v1 "github.com/metal-stack/metal-api/pkg/proto/v1"
	"github.com/metal-stack/metal-api/pkg/util"
	"github.com/metal-stack/metal-lib/zapup"
	"go.uber.org/zap"
	"net/http"
)

func (r *machineResource) findMachine(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("id")

	resp, err := FindMachine(r.ds, id)
	if service.CheckError(request, response, util.CurrentFuncName(), err) {
		return
	}
	err = response.WriteHeaderAndEntity(http.StatusOK, resp)
	if err != nil {
		zapup.MustRootLogger().Error("Failed to send response", zap.Error(err))
	}
}

func FindMachine(ds *datastore.RethinkStore, id string) (*v1.MachineResponse, error) {
	m, err := ds.FindMachineByID(id)
	if err != nil {
		return nil, err
	}
	resp := MakeResponse(m, ds, zapup.MustRootLogger().Sugar())
	return resp, nil
}
