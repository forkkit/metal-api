package ip

import (
	"fmt"
	"github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service/helper"
	v1 "github.com/metal-stack/metal-api/cmd/metal-api/internal/service/v1"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/utils"
	"github.com/metal-stack/metal-lib/httperrors"
	"github.com/metal-stack/metal-lib/zapup"
	"go.uber.org/zap"
	"net/http"
)

func (r ipResource) addUpdateIPRoute(ws *restful.WebService, tags []string) {
	ws.Route(ws.POST("/").
		To(helper.Editor(r.updateIP)).
		Operation("updateIP").
		Doc("updates an ip. if the ip was changed since this one was read, a conflict is returned").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Reads(v1.IPUpdateRequest{}).
		Writes(v1.IPResponse{}).
		Returns(http.StatusOK, "OK", v1.IPResponse{}).
		Returns(http.StatusConflict, "Conflict", httperrors.HTTPErrorResponse{}).
		DefaultReturns("Error", httperrors.HTTPErrorResponse{}))
}

func (r ipResource) updateIP(request *restful.Request, response *restful.Response) {
	var requestPayload v1.IPUpdateRequest
	err := request.ReadEntity(&requestPayload)
	if helper.CheckError(request, response, utils.CurrentFuncName(), err) {
		return
	}

	oldIP, err := r.DS.FindIPByID(requestPayload.IPAddress)
	if helper.CheckError(request, response, utils.CurrentFuncName(), err) {
		return
	}

	newIP := *oldIP
	if requestPayload.Name != nil {
		newIP.Name = *requestPayload.Name
	}
	if requestPayload.Description != nil {
		newIP.Description = *requestPayload.Description
	}
	if requestPayload.Tags != nil {
		newIP.Tags = requestPayload.Tags
	}
	if requestPayload.Type == metal.Static || requestPayload.Type == metal.Ephemeral {
		newIP.Type = requestPayload.Type
	}

	err = r.validateAndUpateIP(oldIP, &newIP)
	if helper.CheckError(request, response, utils.CurrentFuncName(), err) {
		return
	}
	err = response.WriteHeaderAndEntity(http.StatusOK, v1.NewIPResponse(&newIP))
	if err != nil {
		zapup.MustRootLogger().Error("Failed to send response", zap.Error(err))
		return
	}
}

func (r ipResource) validateAndUpateIP(oldIP, newIP *metal.IP) error {
	err := validateIPUpdate(oldIP, newIP)
	if err != nil {
		return err
	}
	tags, err := helper.ProcessTags(newIP.Tags)
	if err != nil {
		return err
	}
	newIP.Tags = tags

	err = r.DS.UpdateIP(oldIP, newIP)
	if err != nil {
		return err
	}
	return nil
}

// Checks whether an ip update is allowed:
// (1) allow update of ephemeral to static
// (2) allow update within a scope
// (3) allow update from and to scope project
// (4) deny all other updates
func validateIPUpdate(old *metal.IP, new *metal.IP) error {
	// constraint 1
	if old.Type == metal.Static && new.Type == metal.Ephemeral {
		return fmt.Errorf("cannot change type of ip address from static to ephemeral")
	}
	os := old.GetScope()
	ns := new.GetScope()
	// constraint 2
	if os == ns {
		return nil
	}
	// constraint 3
	if os == metal.ScopeProject || ns == metal.ScopeProject {
		return nil
	}
	return fmt.Errorf("can not use ip of scope %v with scope %v", os, ns)
}
