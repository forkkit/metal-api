package ip

import (
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service/helper"
	v1 "github.com/metal-stack/metal-api/pkg/proto/v1"
	"github.com/metal-stack/metal-api/pkg/util"
	"github.com/metal-stack/metal-lib/zapup"
	"go.uber.org/zap"
	"net/http"
)

func (r *ipResource) updateIP(request *restful.Request, response *restful.Response) {
	var requestPayload v1.IPUpdateRequest
	err := request.ReadEntity(&requestPayload)
	if helper.CheckError(request, response, util.CurrentFuncName(), err) {
		return
	}

	oldIP, err := r.ds.FindIPByID(requestPayload.Identifiable.IPAddress)
	if helper.CheckError(request, response, util.CurrentFuncName(), err) {
		return
	}

	newIP := *oldIP
	newIP.Name = requestPayload.Common.Name.GetValue()
	newIP.Description = requestPayload.Common.Description.GetValue()

	tags := make([]string, len(requestPayload.Tags))
	for i, t := range requestPayload.Tags {
		tags[i] = t.GetValue()
	}

	if requestPayload.Type == v1.IP_STATIC {
		newIP.Type = metal.Static
	} else if requestPayload.Type == v1.IP_EPHEMERAL {
		newIP.Type = metal.Ephemeral
	}

	err = r.validateAndUpateIP(oldIP, &newIP)
	if helper.CheckError(request, response, util.CurrentFuncName(), err) {
		return
	}
	err = response.WriteHeaderAndEntity(http.StatusOK, service.NewIPResponse(&newIP))
	if err != nil {
		zapup.MustRootLogger().Error("Failed to send response", zap.Error(err))
		return
	}
}

func (r *ipResource) validateAndUpateIP(oldIP, newIP *metal.IP) error {
	err := validateIPUpdate(oldIP, newIP)
	if err != nil {
		return err
	}
	tags, err := helper.ProcessTags(newIP.Tags)
	if err != nil {
		return err
	}
	newIP.Tags = tags

	err = r.ds.UpdateIP(oldIP, newIP)
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
