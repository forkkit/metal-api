package helper

import (
	"bytes"
	"github.com/emicklei/go-restful"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
	"github.com/metal-stack/metal-lib/jwt/sec"
	"github.com/metal-stack/metal-lib/rest"
	"github.com/metal-stack/security"
	"io/ioutil"
	"net/http"
	"time"
)

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
		EMail:  "metal-view@metal-pod.io",
		Name:   "Metal-View",
		Groups: sec.MergeResourceAccess(metal.ViewGroups),
		Tenant: providerTenant,
	}
	ud.edit = security.User{
		EMail:  "metal-edit@metal-pod.io",
		Name:   "Metal-Edit",
		Groups: sec.MergeResourceAccess(metal.EditGroups),
		Tenant: providerTenant,
	}
	ud.admin = security.User{
		EMail:  "metal-admin@metal-pod.io",
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
	keys := make([]string, 0, len(ud.metalUsers))
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
