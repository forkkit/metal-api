package machine

import (
	"github.com/emicklei/go-restful"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/datastore"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service/helper"
	v1 "github.com/metal-stack/metal-api/pkg/proto/v1"
	"github.com/metal-stack/metal-api/pkg/util"
	"github.com/metal-stack/metal-lib/zapup"
	"go.uber.org/zap"
	"net/http"
)

func (r *machineResource) findMachines(request *restful.Request, response *restful.Response) {
	var requestPayload v1.MachineSearchQuery
	err := request.ReadEntity(&requestPayload)
	if service.CheckError(request, response, util.CurrentFuncName(), err) {
		return
	}

	resp, err := FindMachines(r.ds, &requestPayload)
	if service.CheckError(request, response, util.CurrentFuncName(), err) {
		return
	}
	err = response.WriteHeaderAndEntity(http.StatusOK, resp)
	if err != nil {
		zapup.MustRootLogger().Error("Failed to send response", zap.Error(err))
	}
}

func FindMachines(ds *datastore.RethinkStore, query *v1.MachineSearchQuery) ([]*v1.MachineResponse, error) {
	var ms metal.Machines
	err := ds.SearchMachines(query, &ms)
	if err != nil {
		return nil, err
	}

	resp := helper.MakeMachineResponseList(ms, ds, zapup.MustRootLogger().Sugar())
	return resp, nil
}
