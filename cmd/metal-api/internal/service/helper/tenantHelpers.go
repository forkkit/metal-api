package helper

import (
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
	"github.com/metal-stack/metal-api/pkg/helper"
	"github.com/metal-stack/metal-lib/httperrors"
	"github.com/metal-stack/security"
	"net/http"
	"strings"
)

func Viewer(rf restful.RouteFunction) restful.RouteFunction {
	return oneOf(rf, metal.ViewAccess...)
}

func Editor(rf restful.RouteFunction) restful.RouteFunction {
	return oneOf(rf, metal.EditAccess...)
}

func Admin(rf restful.RouteFunction) restful.RouteFunction {
	return oneOf(rf, metal.AdminAccess...)
}

func oneOf(rf restful.RouteFunction, acc ...security.ResourceAccess) restful.RouteFunction {
	return func(request *restful.Request, response *restful.Response) {
		log := helper.Logger(request)
		lg := log.Sugar()
		usr := security.GetUser(request.Request)
		if !usr.HasGroup(acc...) {
			err := fmt.Errorf("you are not member in one of %+v", acc)
			lg.Infow("missing group", "user", usr, "required-group", acc)
			SendError(log, response, helper.CurrentFuncName(), httperrors.NewHTTPError(http.StatusForbidden, err))
			return
		}
		rf(request, response)
	}
}

func tenant(request *restful.Request) string {
	return security.GetUser(request.Request).Tenant
}

type TenantEnsurer struct {
	allowedTenants       map[string]bool
	excludedPathSuffixes []string
}

// NewTenantEnsurer creates a new ensurer with the given tenants.
func NewTenantEnsurer(tenants, excludedPathSuffixes []string) TenantEnsurer {
	result := TenantEnsurer{
		allowedTenants:       make(map[string]bool),
		excludedPathSuffixes: excludedPathSuffixes,
	}
	for _, t := range tenants {
		result.allowedTenants[strings.ToLower(t)] = true
	}
	return result
}

// EnsureAllowedTenantFilter checks if the tenant of the user is allowed.
func (e *TenantEnsurer) EnsureAllowedTenantFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	p := req.Request.URL.Path

	// securing health checks would break monitoring tools
	// preventing liveliness would break status of machines
	for _, suffix := range e.excludedPathSuffixes {
		if strings.HasSuffix(p, suffix) {
			chain.ProcessFilter(req, resp)
			return
		}
	}

	// enforce tenant check otherwise
	tenantID := tenant(req)
	if !e.Allowed(tenantID) {
		err := fmt.Errorf("tenant %s not allowed", tenantID)
		SendError(helper.Logger(req), resp, helper.CurrentFuncName(), httperrors.NewHTTPError(http.StatusForbidden, err))
		return
	}
	chain.ProcessFilter(req, resp)
}

// allowed checks if the given tenant is allowed (case insensitive)
func (e *TenantEnsurer) Allowed(tenant string) bool {
	return e.allowedTenants[strings.ToLower(tenant)]
}
