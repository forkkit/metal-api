package image

import (
	"bytes"
	"encoding/json"
	mdv1 "github.com/metal-stack/masterdata-api/api/v1"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service/machine"
	v1 "github.com/metal-stack/metal-api/pkg/proto/v1"
	"github.com/metal-stack/metal-api/pkg/util"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/metal-stack/metal-api/cmd/metal-api/internal/datastore"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/testdata"

	restful "github.com/emicklei/go-restful"
	"github.com/metal-stack/metal-lib/httperrors"
	"github.com/stretchr/testify/require"
)

func TestGetImages(t *testing.T) {
	ds, mock := datastore.InitMockDB()
	testdata.InitMockDBData(mock)

	imageService := NewImageService(ds)
	container := restful.NewContainer().Add(imageService)
	req := httptest.NewRequest("GET", "/v1/image", nil)
	w := httptest.NewRecorder()
	container.ServeHTTP(w, req)

	resp := w.Result()
	require.Equal(t, http.StatusOK, resp.StatusCode, w.Body.String())
	var result []v1.ImageResponse
	err := json.NewDecoder(resp.Body).Decode(&result)

	require.Nil(t, err)
	require.Len(t, result, 3)
	require.Equal(t, testdata.Img1.ID, result[0].Image.Common.Meta.Id)
	require.Equal(t, testdata.Img1.Name, result[0].Image.Common.Name.GetValue())
	require.Equal(t, testdata.Img2.ID, result[1].Image.Common.Meta.Id)
	require.Equal(t, testdata.Img2.Name, result[1].Image.Common.Name.GetValue())
	require.Equal(t, testdata.Img3.ID, result[2].Image.Common.Meta.Id)
	require.Equal(t, testdata.Img3.Name, result[2].Image.Common.Name.GetValue())
}

func TestGetImage(t *testing.T) {
	ds, mock := datastore.InitMockDB()
	testdata.InitMockDBData(mock)

	imageService := NewImageService(ds)
	container := restful.NewContainer().Add(imageService)
	req := httptest.NewRequest("GET", "/v1/image/1", nil)
	w := httptest.NewRecorder()
	container.ServeHTTP(w, req)

	resp := w.Result()
	require.Equal(t, http.StatusOK, resp.StatusCode, w.Body.String())
	var result v1.ImageResponse
	err := json.NewDecoder(resp.Body).Decode(&result)

	require.Nil(t, err)
	require.Equal(t, testdata.Img1.ID, result.Image.Common.Meta.Id)
	require.Equal(t, testdata.Img1.Name, result.Image.Common.Name.GetValue())
}

func TestGetImageNotFound(t *testing.T) {
	ds, mock := datastore.InitMockDB()
	testdata.InitMockDBData(mock)

	imageService := NewImageService(ds)
	container := restful.NewContainer().Add(imageService)
	req := httptest.NewRequest("GET", "/v1/image/999", nil)
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

func TestDeleteImage(t *testing.T) {
	ds, mock := datastore.InitMockDB()
	testdata.InitMockDBData(mock)

	imageService := NewImageService(ds)
	container := restful.NewContainer().Add(imageService)
	req := httptest.NewRequest("DELETE", "/v1/image/3", nil)
	container = machine.InjectAdmin(container, req)
	w := httptest.NewRecorder()
	container.ServeHTTP(w, req)

	resp := w.Result()
	require.Equal(t, http.StatusOK, resp.StatusCode, w.Body.String())
	var result v1.ImageResponse
	err := json.NewDecoder(resp.Body).Decode(&result)

	require.Nil(t, err)
	require.Equal(t, testdata.Img3.ID, result.Image.Common.Meta.Id)
	require.Equal(t, testdata.Img3.Name, result.Image.Common.Name.GetValue())
}

func TestCreateImage(t *testing.T) {
	ds, mock := datastore.InitMockDB()
	testdata.InitMockDBData(mock)

	createRequest := v1.ImageUpdateRequest{
		Image: &v1.Image{
			Common: &v1.Common{
				Meta: &mdv1.Meta{
					Id: testdata.Img1.ID,
				},
				Name:        util.StringProto(testdata.Img1.Name),
				Description: util.StringProto(testdata.Img1.Description),
			},
			URL: util.StringProto(testdata.Img1.URL),
		},
	}
	js, _ := json.Marshal(createRequest)
	body := bytes.NewBuffer(js)
	req := httptest.NewRequest("PUT", "/v1/image", body)
	container := machine.InjectAdmin(restful.NewContainer().Add(NewImageService(ds)), req)
	req.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	container.ServeHTTP(w, req)

	resp := w.Result()
	require.Equal(t, http.StatusCreated, resp.StatusCode, w.Body.String())
	var result v1.ImageResponse
	err := json.NewDecoder(resp.Body).Decode(&result)

	require.Nil(t, err)
	require.Equal(t, testdata.Img1.ID, result.Image.Common.Meta.Id)
	require.Equal(t, testdata.Img1.Name, result.Image.Common.Name.GetValue())
	require.Equal(t, testdata.Img1.Description, result.Image.Common.Description.GetValue())
	require.Equal(t, testdata.Img1.URL, result.Image.URL.GetValue())
}

func TestUpdateImage(t *testing.T) {
	ds, mock := datastore.InitMockDB()
	testdata.InitMockDBData(mock)

	imageService := NewImageService(ds)
	container := restful.NewContainer().Add(imageService)

	updateRequest := v1.ImageUpdateRequest{
		Image: &v1.Image{
			Common: &v1.Common{
				Meta: &mdv1.Meta{
					Id: testdata.Img1.ID,
				},
				Name:        util.StringProto(testdata.Img2.Name),
				Description: util.StringProto(testdata.Img2.Description),
			},
			URL: util.StringProto(testdata.Img2.URL),
		},
	}
	js, _ := json.Marshal(updateRequest)
	body := bytes.NewBuffer(js)
	req := httptest.NewRequest("POST", "/v1/image", body)
	container = machine.InjectAdmin(container, req)
	req.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	container.ServeHTTP(w, req)

	resp := w.Result()
	require.Equal(t, http.StatusOK, resp.StatusCode, w.Body.String())
	var result v1.ImageResponse
	err := json.NewDecoder(resp.Body).Decode(&result)

	require.Nil(t, err)
	require.Equal(t, testdata.Img1.ID, result.Image.Common.Meta.Id)
	require.Equal(t, testdata.Img2.Name, result.Image.Common.Name.GetValue())
	require.Equal(t, testdata.Img2.Description, result.Image.Common.Description.GetValue())
	require.Equal(t, testdata.Img2.URL, result.Image.URL.GetValue())
}
