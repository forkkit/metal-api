package machine

import (
	"github.com/emicklei/go-restful"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/datastore"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service/helper"
	v1 "github.com/metal-stack/metal-api/pkg/proto/v1"
	"github.com/metal-stack/metal-api/pkg/util"
	"github.com/metal-stack/metal-lib/zapup"
	"go.uber.org/zap"
	"net/http"
)

func (r *machineResource) listMachines(request *restful.Request, response *restful.Response) {
	resp, err := ListMachines(r.ds)
	if service.CheckError(request, response, util.CurrentFuncName(), err) {
		return
	}
	err = response.WriteHeaderAndEntity(http.StatusOK, resp)
	if err != nil {
		zapup.MustRootLogger().Error("Failed to send response", zap.Error(err))
	}
}

func ListMachines(ds *datastore.RethinkStore) ([]*v1.MachineResponse, error) {
	ms, err := ds.ListMachines()
	if err != nil {
		return nil, err
	}
	resp := helper.MakeMachineResponseList(ms, ds, zapup.MustRootLogger().Sugar())
	return resp, nil
}
