package size

import (
	"bytes"
	"encoding/json"
	mdmv1 "github.com/metal-stack/masterdata-api/api/v1"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service"
	v1 "github.com/metal-stack/metal-api/pkg/proto/v1"
	"github.com/metal-stack/metal-api/pkg/util"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/metal-stack/metal-api/cmd/metal-api/internal/datastore"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/testdata"
	"github.com/metal-stack/metal-lib/httperrors"
	"github.com/stretchr/testify/require"

	restful "github.com/emicklei/go-restful"
)

func TestGetSizes(t *testing.T) {
	ds, mock := datastore.InitMockDB()
	testdata.InitMockDBData(mock)

	sizeService := NewSizeService(ds)
	container := restful.NewContainer().Add(sizeService)
	req := httptest.NewRequest("GET", "/v1/size", nil)
	w := httptest.NewRecorder()
	container.ServeHTTP(w, req)

	resp := w.Result()
	require.Equal(t, http.StatusOK, resp.StatusCode, w.Body.String())
	var result []v1.SizeResponse
	err := json.NewDecoder(resp.Body).Decode(&result)

	require.Nil(t, err)
	require.Len(t, result, 3)
	require.Equal(t, testdata.Sz1.ID, result[0].Size.Common.Meta.Id)
	require.Equal(t, testdata.Sz1.Name, result[0].Size.Common.Name.GetValue())
	require.Equal(t, testdata.Sz1.Description, result[0].Size.Common.Description.GetValue())
	require.Equal(t, testdata.Sz2.ID, result[1].Size.Common.Meta.Id)
	require.Equal(t, testdata.Sz2.Name, result[1].Size.Common.Name.GetValue())
	require.Equal(t, testdata.Sz2.Description, result[1].Size.Common.Description.GetValue())
	require.Equal(t, testdata.Sz3.ID, result[2].Size.Common.Meta.Id)
	require.Equal(t, testdata.Sz3.Name, result[2].Size.Common.Name.GetValue())
	require.Equal(t, testdata.Sz3.Description, result[2].Size.Common.Description.GetValue())
}

func TestGetSize(t *testing.T) {
	ds, mock := datastore.InitMockDB()
	testdata.InitMockDBData(mock)

	sizeService := NewSizeService(ds)
	container := restful.NewContainer().Add(sizeService)
	req := httptest.NewRequest("GET", "/v1/size/1", nil)
	w := httptest.NewRecorder()
	container.ServeHTTP(w, req)

	resp := w.Result()
	require.Equal(t, http.StatusOK, resp.StatusCode, w.Body.String())
	var result v1.SizeResponse
	err := json.NewDecoder(resp.Body).Decode(&result)

	require.Nil(t, err)
	require.Equal(t, testdata.Sz1.ID, result.Size.Common.Meta.Id)
	require.Equal(t, testdata.Sz1.Name, result.Size.Common.Name.GetValue())
	require.Equal(t, testdata.Sz1.Description, result.Size.Common.Description.GetValue())
	require.Equal(t, len(testdata.Sz1.Constraints), len(result.Size.Constraints))
}

func TestGetSizeNotFound(t *testing.T) {
	ds, mock := datastore.InitMockDB()
	testdata.InitMockDBData(mock)

	sizeService := NewSizeService(ds)
	container := restful.NewContainer().Add(sizeService)
	req := httptest.NewRequest("GET", "/v1/size/999", nil)
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

func TestDeleteSize(t *testing.T) {
	ds, mock := datastore.InitMockDB()
	testdata.InitMockDBData(mock)

	sizeService := NewSizeService(ds)
	container := restful.NewContainer().Add(sizeService)
	req := httptest.NewRequest("DELETE", "/v1/size/1", nil)
	container = service.InjectAdmin(container, req)
	w := httptest.NewRecorder()
	container.ServeHTTP(w, req)

	resp := w.Result()
	require.Equal(t, http.StatusOK, resp.StatusCode, w.Body.String())
	var result v1.SizeResponse
	err := json.NewDecoder(resp.Body).Decode(&result)

	require.Nil(t, err)
	require.Equal(t, testdata.Sz1.ID, result.Size.Common.Meta.Id)
	require.Equal(t, testdata.Sz1.Name, result.Size.Common.Name.GetValue())
	require.Equal(t, testdata.Sz1.Description, result.Size.Common.Description.GetValue())
}

func TestCreateSize(t *testing.T) {
	ds, mock := datastore.InitMockDB()
	testdata.InitMockDBData(mock)

	sizeService := NewSizeService(ds)
	container := restful.NewContainer().Add(sizeService)

	createRequest := v1.SizeCreateRequest{
		Size: &v1.Size{
			Common: &v1.Common{
				Meta: &mdmv1.Meta{
					Id: testdata.Sz1.ID,
				},
				Name:        util.StringProto(testdata.Sz1.Name),
				Description: util.StringProto(testdata.Sz1.Description),
			},
		},
	}
	js, _ := json.Marshal(createRequest)
	body := bytes.NewBuffer(js)
	req := httptest.NewRequest("PUT", "/v1/size", body)
	req.Header.Add("Content-Type", "application/json")
	container = service.InjectAdmin(container, req)
	w := httptest.NewRecorder()
	container.ServeHTTP(w, req)

	resp := w.Result()
	require.Equal(t, http.StatusCreated, resp.StatusCode, w.Body.String())
	var result v1.SizeResponse
	err := json.NewDecoder(resp.Body).Decode(&result)

	require.Nil(t, err)
	require.Equal(t, testdata.Sz1.ID, result.Size.Common.Meta.Id)
	require.Equal(t, testdata.Sz1.Name, result.Size.Common.Name.GetValue())
	require.Equal(t, testdata.Sz1.Description, result.Size.Common.Description.GetValue())
}

func TestUpdateSize(t *testing.T) {
	ds, mock := datastore.InitMockDB()
	testdata.InitMockDBData(mock)

	sizeService := NewSizeService(ds)
	container := restful.NewContainer().Add(sizeService)

	minCores := uint64(1)
	maxCores := uint64(4)
	updateRequest := v1.SizeUpdateRequest{
		Size: &v1.Size{
			Common: &v1.Common{
				Meta: &mdmv1.Meta{
					Id: testdata.Sz2.ID,
				},
				Name:        util.StringProto(testdata.Sz2.Name),
				Description: util.StringProto(testdata.Sz2.Description),
			},
			Constraints: []*v1.SizeConstraint{
				{
					Type: v1.SizeConstraint_CORES,
					Min:  minCores,
					Max:  maxCores,
				},
			},
		},
	}
	js, _ := json.Marshal(updateRequest)
	body := bytes.NewBuffer(js)
	req := httptest.NewRequest("POST", "/v1/size", body)
	req.Header.Add("Content-Type", "application/json")
	container = service.InjectAdmin(container, req)
	w := httptest.NewRecorder()
	container.ServeHTTP(w, req)

	resp := w.Result()
	require.Equal(t, http.StatusOK, resp.StatusCode, w.Body.String())
	var result v1.SizeResponse
	err := json.NewDecoder(resp.Body).Decode(&result)

	require.Nil(t, err)
	require.Equal(t, testdata.Sz2.ID, result.Size.Common.Meta.Id)
	require.Equal(t, testdata.Sz2.Name, result.Size.Common.Name.GetValue())
	require.Equal(t, testdata.Sz2.Description, result.Size.Common.Description.GetValue())
	require.Equal(t, metal.CoreConstraint, mapSizeConstraintType(result.Size.Constraints[0].Type))
	require.Equal(t, minCores, result.Size.Constraints[0].Min)
	require.Equal(t, maxCores, result.Size.Constraints[0].Max)
}
