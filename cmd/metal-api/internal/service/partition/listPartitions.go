package partition

import (
	"github.com/emicklei/go-restful"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service/helper"
	v1 "github.com/metal-stack/metal-api/pkg/proto/v1"
	"github.com/metal-stack/metal-api/pkg/util"
	"github.com/metal-stack/metal-lib/zapup"
	"go.uber.org/zap"
	"net/http"
)

func (r *partitionResource) listPartitions(request *restful.Request, response *restful.Response) {
	ps, err := r.ds.ListPartitions()
	if helper.CheckError(request, response, util.CurrentFuncName(), err) {
		return
	}

	var result []*v1.PartitionResponse
	for i := range ps {
		result = append(result, service.NewPartitionResponse(&ps[i]))
	}
	err = response.WriteHeaderAndEntity(http.StatusOK, result)
	if err != nil {
		zapup.MustRootLogger().Error("Failed to send response", zap.Error(err))
		return
	}
}
