package partition

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

	"github.com/stretchr/testify/assert"

	restful "github.com/emicklei/go-restful"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/datastore"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/testdata"
	"github.com/metal-stack/metal-lib/httperrors"
	"github.com/stretchr/testify/require"
)

type nopTopicCreator struct {
}

func (n nopTopicCreator) CreateTopic(partitionID, topicFQN string) error {
	return nil
}

type expectingTopicCreator struct {
	t              *testing.T
	expectedTopics []string
}

func (n expectingTopicCreator) CreateTopic(partitionID, topicFQN string) error {
	assert := assert.New(n.t)
	assert.NotEmpty(topicFQN)
	assert.Contains(n.expectedTopics, topicFQN, "Expectation %v contains %s failed.", n.expectedTopics, topicFQN)
	return nil
}

func TestGetPartitions(t *testing.T) {
	ds, mock := datastore.InitMockDB()
	testdata.InitMockDBData(mock)

	service := NewPartitionService(ds, &nopTopicCreator{})
	container := restful.NewContainer().Add(service)
	req := httptest.NewRequest("GET", "/v1/partition", nil)
	w := httptest.NewRecorder()
	container.ServeHTTP(w, req)

	resp := w.Result()
	require.Equal(t, http.StatusOK, resp.StatusCode, w.Body.String())
	var result []v1.PartitionResponse
	err := json.NewDecoder(resp.Body).Decode(&result)

	require.Nil(t, err)
	require.Len(t, result, 3)
	require.Equal(t, testdata.Partition1.ID, result[0].Partition.Common.Meta.Id)
	require.Equal(t, testdata.Partition1.Name, result[0].Partition.Common.Name.GetValue())
	require.Equal(t, testdata.Partition1.Description, result[0].Partition.Common.Description.GetValue())
	require.Equal(t, testdata.Partition2.ID, result[1].Partition.Common.Meta.Id)
	require.Equal(t, testdata.Partition2.Name, result[1].Partition.Common.Name.GetValue())
	require.Equal(t, testdata.Partition2.Description, result[1].Partition.Common.Description.GetValue())
	require.Equal(t, testdata.Partition3.ID, result[2].Partition.Common.Meta.Id)
	require.Equal(t, testdata.Partition3.Name, result[2].Partition.Common.Name.GetValue())
	require.Equal(t, testdata.Partition3.Description, result[2].Partition.Common.Description.GetValue())
}

func TestGetPartition(t *testing.T) {
	ds, mock := datastore.InitMockDB()
	testdata.InitMockDBData(mock)

	service := NewPartitionService(ds, &nopTopicCreator{})
	container := restful.NewContainer().Add(service)
	req := httptest.NewRequest("GET", "/v1/partition/1", nil)
	w := httptest.NewRecorder()
	container.ServeHTTP(w, req)

	resp := w.Result()
	require.Equal(t, http.StatusOK, resp.StatusCode, w.Body.String())
	var result v1.PartitionResponse
	err := json.NewDecoder(resp.Body).Decode(&result)

	require.Nil(t, err)
	require.Equal(t, testdata.Partition1.ID, result.Partition.Common.Meta.Id)
	require.Equal(t, testdata.Partition1.Name, result.Partition.Common.Name.GetValue())
	require.Equal(t, testdata.Partition1.Description, result.Partition.Common.Description.GetValue())
}

func TestGetPartitionNotFound(t *testing.T) {
	ds, mock := datastore.InitMockDB()
	testdata.InitMockDBData(mock)

	service := NewPartitionService(ds, &nopTopicCreator{})
	container := restful.NewContainer().Add(service)
	req := httptest.NewRequest("GET", "/v1/partition/999", nil)
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

func TestDeletePartition(t *testing.T) {
	ds, mock := datastore.InitMockDB()
	testdata.InitMockDBData(mock)

	service := NewPartitionService(ds, &nopTopicCreator{})
	container := restful.NewContainer().Add(service)
	req := httptest.NewRequest("DELETE", "/v1/partition/1", nil)
	container = helper.InjectAdmin(container, req)
	w := httptest.NewRecorder()
	container.ServeHTTP(w, req)

	resp := w.Result()
	require.Equal(t, http.StatusOK, resp.StatusCode, w.Body.String())
	var result v1.PartitionResponse
	err := json.NewDecoder(resp.Body).Decode(&result)

	require.Nil(t, err)
	require.Equal(t, testdata.Partition1.ID, result.Partition.Common.Meta.Id)
	require.Equal(t, testdata.Partition1.Name, result.Partition.Common.Name.GetValue())
	require.Equal(t, testdata.Partition1.Description, result.Partition.Common.Description.GetValue())
}

func TestCreatePartition(t *testing.T) {
	ds, mock := datastore.InitMockDB()
	testdata.InitMockDBData(mock)

	topicCreator := expectingTopicCreator{
		t:              t,
		expectedTopics: []string{"1-switch", "1-machine"},
	}
	service := NewPartitionService(ds, topicCreator)
	container := restful.NewContainer().Add(service)

	createRequest := v1.PartitionCreateRequest{
		Partition: &v1.Partition{
			Common: &v1.Common{
				Meta: &mdmv1.Meta{
					Id: testdata.Partition1.ID,
				},
				Name:        util.StringProto(testdata.Partition1.Name),
				Description: util.StringProto(testdata.Partition1.Description),
			},
		},
	}
	js, _ := json.Marshal(createRequest)
	body := bytes.NewBuffer(js)
	req := httptest.NewRequest("PUT", "/v1/partition", body)
	req.Header.Add("Content-Type", "application/json")
	container = helper.InjectAdmin(container, req)
	w := httptest.NewRecorder()
	container.ServeHTTP(w, req)

	resp := w.Result()
	require.Equal(t, http.StatusCreated, resp.StatusCode, w.Body.String())
	var result v1.PartitionResponse
	err := json.NewDecoder(resp.Body).Decode(&result)

	require.Nil(t, err)
	require.Equal(t, testdata.Partition1.ID, result.Partition.Common.Meta.Id)
	require.Equal(t, testdata.Partition1.Name, result.Partition.Common.Name.GetValue())
	require.Equal(t, testdata.Partition1.Description, result.Partition.Common.Description.GetValue())
}

func TestUpdatePartition(t *testing.T) {
	ds, mock := datastore.InitMockDB()
	testdata.InitMockDBData(mock)

	service := NewPartitionService(ds, &nopTopicCreator{})
	container := restful.NewContainer().Add(service)

	mgmtService := "mgmt"
	imageURL := "http://somewhere/image1.zip"
	kernelURL := "http://somewhere/kernel1.zip"
	cmdLine := "cmdline"
	updateRequest := v1.PartitionUpdateRequest{
		Partition: &v1.Partition{
			Common: &v1.Common{
				Meta: &mdmv1.Meta{
					Id: testdata.Partition2.ID,
				},
				Name:        util.StringProto(testdata.Partition2.Name),
				Description: util.StringProto(testdata.Partition2.Description),
			},
			MgmtServiceAddress: util.StringProto(mgmtService),
			ImageURL:           util.StringProto(imageURL),
			KernelURL:          util.StringProto(kernelURL),
			CommandLine:        util.StringProto(cmdLine),
		},
	}
	js, _ := json.Marshal(updateRequest)
	body := bytes.NewBuffer(js)
	req := httptest.NewRequest("POST", "/v1/partition", body)
	req.Header.Add("Content-Type", "application/json")
	container = helper.InjectAdmin(container, req)
	w := httptest.NewRecorder()
	container.ServeHTTP(w, req)

	resp := w.Result()
	require.Equal(t, http.StatusOK, resp.StatusCode, w.Body.String())
	var result v1.PartitionResponse
	err := json.NewDecoder(resp.Body).Decode(&result)

	require.Nil(t, err)
	require.Equal(t, testdata.Partition2.ID, result.Partition.Common.Meta.Id)
	require.Equal(t, testdata.Partition2.Name, result.Partition.Common.Name.GetValue())
	require.Equal(t, testdata.Partition2.Description, result.Partition.Common.Description.GetValue())
	require.Equal(t, mgmtService, result.Partition.MgmtServiceAddress.GetValue())
	require.Equal(t, imageURL, result.Partition.ImageURL.GetValue())
	require.Equal(t, kernelURL, result.Partition.KernelURL.GetValue())
	require.Equal(t, cmdLine, result.Partition.CommandLine.GetValue())
}

func TestPartitionCapacity(t *testing.T) {
	ds, mock := datastore.InitMockDB()
	testdata.InitMockDBData(mock)

	service := NewPartitionService(ds, &nopTopicCreator{})
	container := restful.NewContainer().Add(service)

	req := httptest.NewRequest("GET", "/v1/partition/capacity", nil)
	req.Header.Add("Content-Type", "application/json")
	container = helper.InjectAdmin(container, req)
	w := httptest.NewRecorder()
	container.ServeHTTP(w, req)

	resp := w.Result()
	require.Equal(t, http.StatusOK, resp.StatusCode, w.Body.String())
	var result []*v1.PartitionCapacity
	err := json.NewDecoder(resp.Body).Decode(&result)

	require.Nil(t, err)
	require.Equal(t, testdata.Partition1.ID, result[0].Common.Meta.Id)
	require.NotNil(t, result[0].ServerCapacities)
	require.Equal(t, 1, len(result[0].ServerCapacities))
	capacity := result[0].ServerCapacities[0]
	require.Equal(t, "1", capacity.Size)
	require.Equal(t, 5, capacity.Total)
	require.Equal(t, 0, capacity.Free)
}
