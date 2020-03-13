package network

import (
	"bytes"
	"encoding/json"
	mdmv1 "github.com/metal-stack/masterdata-api/api/v1"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service/helper"
	v1 "github.com/metal-stack/metal-api/pkg/proto/v1"
	"github.com/metal-stack/metal-api/pkg/util"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/metal-stack/metal-lib/httperrors"
	r "gopkg.in/rethinkdb/rethinkdb-go.v5"

	"github.com/metal-stack/metal-api/cmd/metal-api/internal/datastore"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/ipam"

	restful "github.com/emicklei/go-restful"
	goipam "github.com/metal-stack/go-ipam"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/testdata"
	"github.com/stretchr/testify/require"
)

func TestGetNetworks(t *testing.T) {
	ds, mock := datastore.InitMockDB()
	testdata.InitMockDBData(mock)

	networkservice := NewNetworkService(ds, ipam.New(goipam.New()), nil)
	container := restful.NewContainer().Add(networkservice)
	req := httptest.NewRequest("GET", "/v1/network", nil)
	container = helper.InjectViewer(container, req)
	w := httptest.NewRecorder()
	container.ServeHTTP(w, req)

	resp := w.Result()
	require.Equal(t, http.StatusOK, resp.StatusCode, w.Body.String())
	var result []v1.NetworkResponse
	err := json.NewDecoder(resp.Body).Decode(&result)

	require.Nil(t, err)
	require.Len(t, result, 4)
	require.Equal(t, testdata.Nw1.ID, result[0].Network.Common.Meta.Id)
	require.Equal(t, testdata.Nw1.Name, result[0].Network.Common.Name.GetValue())
	require.Equal(t, testdata.Nw1.Description, result[0].Network.Common.Description.GetValue())
	require.Equal(t, testdata.Nw2.ID, result[1].Network.Common.Meta.Id)
	require.Equal(t, testdata.Nw2.Name, result[1].Network.Common.Name.GetValue())
	require.Equal(t, testdata.Nw2.Description, result[1].Network.Common.Description.GetValue())
	require.Equal(t, testdata.Nw3.ID, result[2].Network.Common.Meta.Id)
	require.Equal(t, testdata.Nw3.Name, result[2].Network.Common.Name.GetValue())
	require.Equal(t, testdata.Nw3.Description, result[2].Network.Common.Description.GetValue())
}

func TestGetNetwork(t *testing.T) {
	ds, mock := datastore.InitMockDB()
	testdata.InitMockDBData(mock)

	networkservice := NewNetworkService(ds, ipam.New(goipam.New()), nil)
	container := restful.NewContainer().Add(networkservice)
	req := httptest.NewRequest("GET", "/v1/network/1", nil)
	container = helper.InjectViewer(container, req)
	w := httptest.NewRecorder()
	container.ServeHTTP(w, req)

	resp := w.Result()
	require.Equal(t, http.StatusOK, resp.StatusCode, w.Body.String())
	var result v1.NetworkResponse
	err := json.NewDecoder(resp.Body).Decode(&result)

	require.Nil(t, err)
	require.Equal(t, testdata.Nw1.ID, result.Network.Common.Meta.Id)
	require.Equal(t, testdata.Nw1.Name, result.Network.Common.Name.GetValue())
	require.Equal(t, testdata.Nw1.Description, result.Network.Common.Description.GetValue())
}

func TestGetNetworkNotFound(t *testing.T) {
	ds, mock := datastore.InitMockDB()
	testdata.InitMockDBData(mock)

	networkservice := NewNetworkService(ds, ipam.New(goipam.New()), nil)
	container := restful.NewContainer().Add(networkservice)
	req := httptest.NewRequest("GET", "/v1/network/999", nil)
	container = helper.InjectViewer(container, req)
	w := httptest.NewRecorder()
	container.ServeHTTP(w, req)

	resp := w.Result()
	require.Equal(t, http.StatusNotFound, resp.StatusCode, w.Body.String())
	var result httperrors.HTTPErrorResponse
	err := json.NewDecoder(resp.Body).Decode(&result)

	require.Nil(t, err)
	require.Contains(t, result.Message, "999")
	require.Equal(t, 404, result.StatusCode)
}

func TestDeleteNetwork(t *testing.T) {
	ds, mock := datastore.InitMockDB()
	mock.On(r.DB("mockdb").Table("network").Filter(r.MockAnything())).Return([]interface{}{}, nil)
	ipamer, err := testdata.InitMockIpamData(mock, false)
	require.Nil(t, err)
	testdata.InitMockDBData(mock)

	networkservice := NewNetworkService(ds, ipamer, nil)
	container := restful.NewContainer().Add(networkservice)
	req := httptest.NewRequest("DELETE", "/v1/network/"+testdata.NwIPAM.ID, nil)
	container = helper.InjectAdmin(container, req)
	w := httptest.NewRecorder()
	container.ServeHTTP(w, req)

	resp := w.Result()
	require.Equal(t, http.StatusOK, resp.StatusCode, w.Body.String())
	var result v1.NetworkResponse
	err = json.NewDecoder(resp.Body).Decode(&result)

	require.Nil(t, err)
	require.Equal(t, testdata.NwIPAM.ID, result.Network.Common.Meta.Id)
	require.Equal(t, testdata.NwIPAM.Name, result.Network.Common.Name.GetValue())
	require.Equal(t, testdata.NwIPAM.Description, result.Network.Common.Description.GetValue())
}

func TestDeleteNetworkIPInUse(t *testing.T) {
	ds, mock := datastore.InitMockDB()
	mock.On(r.DB("mockdb").Table("network").Filter(r.MockAnything())).Return([]interface{}{}, nil)
	ipamer, err := testdata.InitMockIpamData(mock, true)
	require.Nil(t, err)
	testdata.InitMockDBData(mock)

	networkservice := NewNetworkService(ds, ipamer, nil)
	container := restful.NewContainer().Add(networkservice)
	req := httptest.NewRequest("DELETE", "/v1/network/"+testdata.NwIPAM.ID, nil)
	container = helper.InjectAdmin(container, req)
	w := httptest.NewRecorder()
	container.ServeHTTP(w, req)

	resp := w.Result()
	require.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode, w.Body.String())
	var result httperrors.HTTPErrorResponse
	err = json.NewDecoder(resp.Body).Decode(&result)

	require.Nil(t, err)
	require.Equal(t, 422, result.StatusCode)
	require.Contains(t, result.Message, "unable to delete Network: prefix 10.0.0.0/16 has ip 10.0.0.1 in use")
}

func TestCreateNetwork(t *testing.T) {
	ds, mock := datastore.InitMockDB()
	ipamer, err := testdata.InitMockIpamData(mock, false)
	require.Nil(t, err)
	testdata.InitMockDBData(mock)

	networkservice := NewNetworkService(ds, ipamer, nil)
	container := restful.NewContainer().Add(networkservice)

	prefixes := []string{"172.0.0.0/24"}
	destPrefixes := []string{"0.0.0.0/0"}
	vrf := uint(10000)
	createRequest := &v1.NetworkCreateRequest{
		Network: &v1.Network{
			Common: &v1.Common{
				Meta: &mdmv1.Meta{
					Id: testdata.Nw1.ID,
				},
				Name:        util.StringProto(testdata.Nw1.Name),
				Description: util.StringProto(testdata.Nw1.Description),
			},
			ProjectID:   util.StringProto(testdata.Nw1.ProjectID),
			PartitionID: util.StringProto(testdata.Nw1.PartitionID),
		},
		NetworkImmutable: &v1.NetworkImmutable{
			Prefixes:            prefixes,
			DestinationPrefixes: destPrefixes,
			Vrf:                 util.UInt64Proto(vrf),
		},
	}
	js, _ := json.Marshal(createRequest)
	body := bytes.NewBuffer(js)
	req := httptest.NewRequest("PUT", "/v1/network", body)
	req.Header.Add("Content-Type", "application/json")
	container = helper.InjectAdmin(container, req)
	w := httptest.NewRecorder()
	container.ServeHTTP(w, req)

	resp := w.Result()
	require.Equal(t, http.StatusCreated, resp.StatusCode, w.Body.String())
	var result v1.NetworkResponse
	err = json.NewDecoder(resp.Body).Decode(&result)

	require.Nil(t, err)
	require.Equal(t, testdata.Nw1.Name, result.Network.Common.Name.GetValue())
	require.Equal(t, testdata.Nw1.Description, result.Network.Common.Description.GetValue())
	require.Equal(t, testdata.Nw1.PartitionID, result.Network.PartitionID.GetValue())
	require.Equal(t, testdata.Nw1.ProjectID, result.Network.ProjectID.GetValue())
	require.Equal(t, destPrefixes, result.NetworkImmutable.DestinationPrefixes)
}

func TestUpdateNetwork(t *testing.T) {
	ds, mock := datastore.InitMockDB()
	testdata.InitMockDBData(mock)

	networkservice := NewNetworkService(ds, ipam.New(goipam.New()), nil)
	container := restful.NewContainer().Add(networkservice)

	newName := "new"
	updateRequest := &v1.NetworkUpdateRequest{
		Common: &v1.Common{
			Meta: &mdmv1.Meta{
				Id: testdata.Nw1.ID,
			},
			Name:        util.StringProto(newName),
			Description: util.StringProto(testdata.Nw1.Description),
		},
	}
	js, _ := json.Marshal(updateRequest)
	body := bytes.NewBuffer(js)
	req := httptest.NewRequest("POST", "/v1/network", body)
	req.Header.Add("Content-Type", "application/json")
	container = helper.InjectAdmin(container, req)
	w := httptest.NewRecorder()
	container.ServeHTTP(w, req)

	resp := w.Result()
	require.Equal(t, http.StatusOK, resp.StatusCode, w.Body.String())
	var result v1.NetworkResponse
	err := json.NewDecoder(resp.Body).Decode(&result)

	require.Nil(t, err)
	require.Equal(t, testdata.Nw1.ID, result.Network.Common.Meta.Id)
	require.Equal(t, newName, result.Network.Common.Name.GetValue())
}

func TestSearchNetwork(t *testing.T) {
	ds, mock := datastore.InitMockDB()
	mock.On(r.DB("mockdb").Table("network").Filter(r.MockAnything())).Return([]interface{}{testdata.Nw1}, nil)
	testdata.InitMockDBData(mock)

	networkService := NewNetworkService(ds, ipam.New(goipam.New()), nil)
	container := restful.NewContainer().Add(networkService)
	findReq := &v1.NetworkFindRequest{
		NetworkSearchQuery: &v1.NetworkSearchQuery{
			PartitionID: util.StringProto("1"),
		},
	}
	requestJSON, err := json.Marshal(findReq)
	require.Nil(t, err)
	req := httptest.NewRequest("POST", "/v1/network/find", bytes.NewBuffer(requestJSON))
	req.Header.Add("Content-Type", "application/json")
	container = helper.InjectViewer(container, req)
	w := httptest.NewRecorder()
	container.ServeHTTP(w, req)

	resp := w.Result()
	require.Equal(t, http.StatusOK, resp.StatusCode, w.Body.String())
	var results []*v1.NetworkResponse
	err = json.NewDecoder(resp.Body).Decode(&results)

	require.Nil(t, err)
	require.Len(t, results, len(testdata.TestNetworks))
	result := results[0]
	require.Equal(t, testdata.Nw1.ID, result.Network.Common.Meta.Id)
	require.Equal(t, testdata.Nw1.Name, result.Network.Common.Name.GetValue())
	require.Equal(t, testdata.Nw1.Description, result.Network.Common.Description.GetValue())
	require.Equal(t, testdata.Nw1.ProjectID, result.Network.ProjectID.GetValue())
	require.Equal(t, testdata.Nw1.PartitionID, result.Network.PartitionID.GetValue())
}
