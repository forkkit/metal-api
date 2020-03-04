package helper

import (
	"encoding/json"
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/go-stack/stack"
	mdmv1 "github.com/metal-stack/masterdata-api/api/v1"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/utils"
	"github.com/metal-stack/metal-lib/httperrors"
	"go.uber.org/zap"
	"net/http"
)

func CheckError(rq *restful.Request, rsp *restful.Response, opname string, err error) bool {
	log := utils.Logger(rq)
	if err != nil {
		if metal.IsNotFound(err) {
			sendErrorImpl(log, rsp, opname, httperrors.NotFound(err), 2)
			return true
		}
		if metal.IsConflict(err) {
			sendErrorImpl(log, rsp, opname, httperrors.Conflict(err), 2)
			return true
		}
		if metal.IsInternal(err) {
			sendErrorImpl(log, rsp, opname, httperrors.InternalServerError(err), 2)
			return true
		}
		if mdmv1.IsNotFound(err) {
			sendErrorImpl(log, rsp, opname, httperrors.NotFound(err), 2)
			return true
		}
		if mdmv1.IsConflict(err) {
			sendErrorImpl(log, rsp, opname, httperrors.Conflict(err), 2)
			return true
		}
		if mdmv1.IsInternal(err) {
			sendErrorImpl(log, rsp, opname, httperrors.InternalServerError(err), 2)
			return true
		}
		sendErrorImpl(log, rsp, opname, httperrors.NewHTTPError(http.StatusUnprocessableEntity, err), 2)
		return true
	}
	return false
}

func SendError(log *zap.Logger, rsp *restful.Response, opname string, errRsp *httperrors.HTTPErrorResponse) {
	sendErrorImpl(log, rsp, opname, errRsp, 1)
}

func sendErrorImpl(log *zap.Logger, rsp *restful.Response, opname string, errRsp *httperrors.HTTPErrorResponse, stackup int) {
	s := stack.Caller(stackup)
	response, merr := json.Marshal(errRsp)
	log.Error("service error", zap.String("operation", opname), zap.Int("status", errRsp.StatusCode), zap.String("error", errRsp.Message), zap.Stringer("service-caller", s), zap.String("resp", string(response)))
	if merr != nil {
		err := rsp.WriteError(http.StatusInternalServerError, fmt.Errorf("unable to format error string: %v", merr))
		if err != nil {
			log.Error("Failed to send response", zap.Error(err))
			return
		}
		return
	}
	err := rsp.WriteErrorString(errRsp.StatusCode, string(response))
	if err != nil {
		log.Error("Failed to send response", zap.Error(err))
		return
	}
}
