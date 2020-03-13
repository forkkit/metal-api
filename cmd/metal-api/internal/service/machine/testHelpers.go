package machine

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/emicklei/go-restful"
	v12 "github.com/metal-stack/masterdata-api/api/v1"
	"github.com/metal-stack/masterdata-api/pkg/client"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/datastore"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/eventbus"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/ipam"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service/helper"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service/image"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service/ip"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service/network"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service/partition"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service/size"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service/sw"
	"github.com/metal-stack/metal-api/pkg/proto/v1"
	"github.com/metal-stack/metal-api/pkg/util"
	"github.com/metal-stack/metal-lib/jwt/sec"
	"github.com/metal-stack/metal-lib/rest"
	"github.com/metal-stack/security"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"go.uber.org/zap"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
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

type testEnv struct {
	imageService        *restful.WebService
	switchService       *restful.WebService
	sizeService         *restful.WebService
	networkService      *restful.WebService
	partitionService    *restful.WebService
	machineService      *restful.WebService
	ipService           *restful.WebService
	PrivateSuperNetwork *v1.NetworkResponse
	PrivateNetwork      *v1.NetworkResponse
	rethinkContainer    testcontainers.Container
	ctx                 context.Context
}

func (te *testEnv) Teardown() {
	_ = te.rethinkContainer.Terminate(te.ctx)
}

func CreateTestEnvironment(t *testing.T) testEnv {
	require := require.New(t)
	log, err := zap.NewDevelopment()
	require.NoError(err)

	ipamer := ipam.InitTestIpam(t)
	nsq := eventbus.InitTestPublisher(t)
	ds, rc, ctx := datastore.InitTestDB(t)
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	mdc, err := client.NewClient(timeoutCtx, "localhost", 50051, "certs/client.pem", "certs/client-key.pem", "certs/ca.pem", "hmac", log)
	require.NoError(err)

	machineService := NewMachineService(ds, nsq.Publisher, ipamer, mdc)
	imageService := image.NewImageService(ds)
	switchService := sw.NewSwitchService(ds)
	sizeService := size.NewSizeService(ds)
	networkService := network.NewNetworkService(ds, ipamer, mdc)
	partitionService := partition.NewPartitionService(ds, nsq)
	ipService := ip.NewIPService(ds, ipamer, mdc)

	te := testEnv{
		imageService:     imageService,
		switchService:    switchService,
		sizeService:      sizeService,
		networkService:   networkService,
		partitionService: partitionService,
		machineService:   machineService,
		ipService:        ipService,
		rethinkContainer: rc,
		ctx:              ctx,
	}

	imageID := "test-image"
	imageName := "testimage"
	imageDesc := "Test Image"
	img := v1.ImageCreateRequest{
		Image: &v1.Image{
			Common: &v1.Common{
				Meta: &v12.Meta{
					Id: imageID,
				},
				Name:        util.StringProto(imageName),
				Description: util.StringProto(imageDesc),
			},
			URL:      util.StringProto("https://blobstore/image"),
			Features: util.StringSliceProto(string(metal.ImageFeatureMachine)),
		},
	}
	var createdImage v1.ImageResponse

	status := te.ImageCreate(t, img, &createdImage)
	require.Equal(http.StatusCreated, status)
	require.NotNil(createdImage)
	require.Equal(img.Image.Common.Meta.Id, createdImage.Image.Common.Meta.Id)

	sizeName := "testsize"
	sizeDesc := "Test Size"
	s := v1.SizeCreateRequest{
		Size: &v1.Size{
			Common: &v1.Common{
				Meta: &v12.Meta{
					Id: "test-size",
				},
				Name:        util.StringProto(sizeName),
				Description: util.StringProto(sizeDesc),
			},
			Constraints: []*v1.SizeConstraint{
				{
					Type: v1.SizeConstraint_CORES,
					Min:  8,
					Max:  8,
				},
				{
					Type: v1.SizeConstraint_MEMORY,
					Min:  1000,
					Max:  2000,
				},
				{
					Type: v1.SizeConstraint_STORAGE,
					Min:  2000,
					Max:  3000,
				},
			},
		},
	}
	var createdSize v1.SizeResponse
	status = te.SizeCreate(t, s, &createdSize)
	require.Equal(http.StatusCreated, status)
	require.NotNil(createdSize)
	require.Equal(s.Size.Common.Meta.Id, createdSize.Size.Common.Meta.Id)

	partName := "test-partition"
	partDesc := "Test Partition"
	part := v1.PartitionCreateRequest{
		Partition: &v1.Partition{
			Common: &v1.Common{
				Meta: &v12.Meta{
					Id: "test-partition",
				},
				Name:        util.StringProto(partName),
				Description: util.StringProto(partDesc),
			},
		},
	}
	var createdPartition v1.PartitionResponse
	status = te.PartitionCreate(t, part, &createdPartition)
	require.Equal(http.StatusCreated, status)
	require.NotNil(createdPartition)
	require.Equal(part.Partition.Common.Name.GetValue(), createdPartition.Partition.Common.Name.GetValue())
	require.NotEmpty(createdPartition.Partition.Common.Meta.Id)

	switchID := "test-switch01"
	sw := v1.SwitchRegisterRequest{
		Switch: &v1.Switch{
			Common: &v1.Common{
				Meta: &v12.Meta{
					Id: switchID,
				},
			},
			RackID: "test-rack",
			Nics: helper.SwitchNics{
				{
					MacAddress: "bb:aa:aa:aa:aa:aa",
					Name:       "swp1",
				},
			},
		},
		PartitionID: "test-partition",
	}
	var createdSwitch v1.SwitchResponse

	status = te.SwitchRegister(t, sw, &createdSwitch)
	require.Equal(http.StatusCreated, status)
	require.NotNil(createdSwitch)
	require.Equal(sw.Switch.Common.Meta.Id, createdSwitch.Switch.Common.Meta.Id)
	require.Len(sw.Switch.Nics, 1)
	require.Equal(sw.Switch.Nics[0].Name, createdSwitch.Switch.Nics[0].Name)
	require.Equal(sw.Switch.Nics[0].MacAddress, createdSwitch.Switch.Nics[0].MacAddress)

	var createdNetwork v1.NetworkResponse
	networkID := "test-private-super"
	networkName := "test-private-super-network"
	networkDesc := "Test Private Super Network"
	testPrivateSuperCidr := "10.0.0.0/16"
	ncr := v1.NetworkCreateRequest{
		Network: &v1.Network{
			Common: &v1.Common{
				Meta: &v12.Meta{
					Id: networkID,
				},
				Name:        util.StringProto(networkName),
				Description: util.StringProto(networkDesc),
			},
			PartitionID: util.StringProto(part.Partition.Common.Meta.Id),
		},
		NetworkImmutable: &v1.NetworkImmutable{
			Prefixes:     []string{testPrivateSuperCidr},
			PrivateSuper: true,
		},
	}
	status = te.NetworkCreate(t, ncr, &createdNetwork)
	require.Equal(http.StatusCreated, status)
	require.NotNil(createdNetwork)
	require.Equal(ncr.Network.Common.Meta.Id, createdNetwork.Network.Common.Meta.Id)

	te.PrivateSuperNetwork = &createdNetwork

	var acquiredPrivateNetwork v1.NetworkResponse
	privateNetworkName := "test-private-network"
	privateNetworkDesc := "Test Private Network"
	projectID := "test-project-1"
	nar := v1.NetworkAllocateRequest{
		Network: &v1.Network{
			Common: &v1.Common{
				Name:        util.StringProto(privateNetworkName),
				Description: util.StringProto(privateNetworkDesc),
			},
			ProjectID:   util.StringProto(projectID),
			PartitionID: util.StringProto(part.Partition.Common.Meta.Id),
		},
	}
	status = te.NetworkAcquire(t, nar, &acquiredPrivateNetwork)
	require.Equal(http.StatusCreated, status)
	require.NotNil(acquiredPrivateNetwork)
	require.Equal(ncr.Network.Common.Meta.Id, acquiredPrivateNetwork.NetworkImmutable.ParentNetworkID)
	require.Len(acquiredPrivateNetwork.NetworkImmutable.Prefixes, 1)
	_, ipnet, _ := net.ParseCIDR(testPrivateSuperCidr)
	_, privateNet, _ := net.ParseCIDR(acquiredPrivateNetwork.NetworkImmutable.Prefixes[0])
	require.True(ipnet.Contains(privateNet.IP), "%s must be within %s", privateNet, ipnet)
	te.PrivateNetwork = &acquiredPrivateNetwork

	return te
}

func (te *testEnv) SizeCreate(t *testing.T, icr v1.SizeCreateRequest, response interface{}) int {
	return webRequestPut(t, te.sizeService, icr, "/v1/size/", response)
}

func (te *testEnv) PartitionCreate(t *testing.T, icr v1.PartitionCreateRequest, response interface{}) int {
	return webRequestPut(t, te.partitionService, icr, "/v1/partition/", response)
}

func (te *testEnv) SwitchRegister(t *testing.T, srr v1.SwitchRegisterRequest, response interface{}) int {
	return webRequestPost(t, te.switchService, srr, "/v1/switch/register", response)
}

func (te *testEnv) SwitchGet(t *testing.T, swid string, response interface{}) int {
	return webRequestGet(t, te.switchService, emptyBody{}, "/v1/switch/"+swid, response)
}

func (te *testEnv) ImageCreate(t *testing.T, icr v1.ImageCreateRequest, response interface{}) int {
	return webRequestPut(t, te.imageService, icr, "/v1/image/", response)
}

func (te *testEnv) NetworkCreate(t *testing.T, icr v1.NetworkCreateRequest, response interface{}) int {
	return webRequestPut(t, te.networkService, icr, "/v1/network/", response)
}

func (te *testEnv) NetworkAcquire(t *testing.T, nar v1.NetworkAllocateRequest, response interface{}) int {
	return webRequestPost(t, te.networkService, nar, "/v1/network/allocate", response)
}

func (te *testEnv) MachineAllocate(t *testing.T, mar v1.MachineAllocateRequest, response interface{}) int {
	return webRequestPost(t, te.machineService, mar, "/v1/machine/allocate", response)
}

func (te *testEnv) MachineFree(t *testing.T, uuid string, response interface{}) int {
	return webRequestDelete(t, te.machineService, &emptyBody{}, "/v1/machine/"+uuid+"/free", response)
}

func (te *testEnv) MachineRegister(t *testing.T, mrr v1.MachineRegisterRequest, response interface{}) int {
	return webRequestPost(t, te.machineService, mrr, "/v1/machine/register", response)
}

func (te *testEnv) MachineWait(uuid string) {
	container := restful.NewContainer().Add(te.machineService)
	createReq := httptest.NewRequest(http.MethodGet, "/v1/machine/"+uuid+"/wait", nil)
	container = InjectAdmin(container, createReq)
	w := httptest.NewRecorder()
	for {
		container.ServeHTTP(w, createReq)
		resp := w.Result()
		var response map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&response)
		if err != nil {
			panic(err)
		}
		if resp.StatusCode == http.StatusOK {
			break
		}
		if resp.StatusCode == http.StatusInternalServerError {
			break
		}
	}
}

//nolint:golint,unused
type emptyBody struct{}

//nolint:golint,unused
func webRequestPut(t *testing.T, service *restful.WebService, request interface{}, path string, response interface{}) int {
	return webRequest(t, http.MethodPut, service, request, path, response)
}

//nolint:golint,unused
func webRequestPost(t *testing.T, service *restful.WebService, request interface{}, path string, response interface{}) int {
	return webRequest(t, http.MethodPost, service, request, path, response)
}

//nolint:golint,unused
func webRequestDelete(t *testing.T, service *restful.WebService, request interface{}, path string, response interface{}) int {
	return webRequest(t, http.MethodDelete, service, request, path, response)
}

//nolint:golint,unused
func webRequestGet(t *testing.T, service *restful.WebService, request interface{}, path string, response interface{}) int {
	return webRequest(t, http.MethodGet, service, request, path, response)
}

//nolint:golint,unused
func webRequest(t *testing.T, method string, service *restful.WebService, request interface{}, path string, response interface{}) int {
	container := restful.NewContainer().Add(service)

	jsonBody, err := json.Marshal(request)
	require.NoError(t, err)
	body := ioutil.NopCloser(strings.NewReader(string(jsonBody)))
	createReq := httptest.NewRequest(method, path, body)
	createReq.Header.Set("Content-Type", "application/json")

	container = InjectAdmin(container, createReq)
	w := httptest.NewRecorder()
	container.ServeHTTP(w, createReq)

	resp := w.Result()
	err = json.NewDecoder(resp.Body).Decode(response)
	require.NoError(t, err)
	return resp.StatusCode
}
