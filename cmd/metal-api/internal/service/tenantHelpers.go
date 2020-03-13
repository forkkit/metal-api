package service

import (
	"bytes"
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
	"github.com/metal-stack/metal-api/pkg/util"
	"github.com/metal-stack/metal-lib/httperrors"
	"github.com/metal-stack/metal-lib/jwt/sec"
	"github.com/metal-stack/metal-lib/rest"
	"github.com/metal-stack/security"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
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
		log := util.Logger(request)
		lg := log.Sugar()
		usr := security.GetUser(request.Request)
		if !usr.HasGroup(acc...) {
			err := fmt.Errorf("you are not member in one of %+v", acc)
			lg.Infow("missing group", "user", usr, "required-group", acc)
			SendError(log, response, util.CurrentFuncName(), httperrors.NewHTTPError(http.StatusForbidden, err))
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
		SendError(util.Logger(req), resp, util.CurrentFuncName(), httperrors.NewHTTPError(http.StatusForbidden, err))
		return
	}
	chain.ProcessFilter(req, resp)
}

// allowed checks if the given tenant is allowed (case insensitive)
func (e *TenantEnsurer) Allowed(tenant string) bool {
	return e.allowedTenants[strings.ToLower(tenant)]
}

type UserDirectory struct {
	viewer security.User
	edit   security.User
	admin  security.User

	metalUsers map[string]security.User
}

func NewUserDirectory(providerTenant string) *UserDirectory {
	ud := &UserDirectory{}

	// User.Name is used as AuthType for HMAC
	ud.viewer = security.User{
		EMail:  "metal-view@metal-stack.io",
		Name:   "Metal-View",
		Groups: sec.MergeResourceAccess(metal.ViewGroups),
		Tenant: providerTenant,
	}
	ud.edit = security.User{
		EMail:  "metal-edit@metal-stack.io",
		Name:   "Metal-Edit",
		Groups: sec.MergeResourceAccess(metal.EditGroups),
		Tenant: providerTenant,
	}
	ud.admin = security.User{
		EMail:  "metal-admin@metal-stack.io",
		Name:   "Metal-Admin",
		Groups: sec.MergeResourceAccess(metal.AdminGroups),
		Tenant: providerTenant,
	}
	ud.metalUsers = map[string]security.User{
		"view":  ud.viewer,
		"edit":  ud.edit,
		"admin": ud.admin,
	}

	return ud
}

func (ud *UserDirectory) UserNames() []string {
	keys := make([]string, len(ud.metalUsers))
	for k := range ud.metalUsers {
		keys = append(keys, k)
	}
	return keys
}

func (ud *UserDirectory) Get(user string) security.User {
	return ud.metalUsers[user]
}

var testUserDirectory = NewUserDirectory("")

func InjectViewer(container *restful.Container, rq *http.Request) *restful.Container {
	return injectUser(testUserDirectory.viewer, container, rq)
}

func InjectEditor(container *restful.Container, rq *http.Request) *restful.Container {
	return injectUser(testUserDirectory.edit, container, rq)
}
func InjectAdmin(container *restful.Container, rq *http.Request) *restful.Container {
	return injectUser(testUserDirectory.admin, container, rq)
}

func injectUser(u security.User, container *restful.Container, rq *http.Request) *restful.Container {
	hma := security.NewHMACAuth(u.Name, []byte{1, 2, 3}, security.WithUser(u))
	usergetter := security.NewCreds(security.WithHMAC(hma))
	container.Filter(rest.UserAuth(usergetter))
	var body []byte
	if rq.Body != nil {
		data, _ := ioutil.ReadAll(rq.Body)
		body = data
		rq.Body.Close()
		rq.Body = ioutil.NopCloser(bytes.NewReader(data))
	}
	hma.AddAuth(rq, time.Now(), body)
	return container
}
