package helper

import (
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/emicklei/go-restful"
)

func TestTenantEnsurer(t *testing.T) {
	e := service.NewTenantEnsurer([]string{"pvdr", "Pv", "pv-DR"}, nil)
	require.True(t, e.Allowed("pvdr"))
	require.True(t, e.Allowed("Pv"))
	require.True(t, e.Allowed("pv"))
	require.True(t, e.Allowed("pv-DR"))
	require.True(t, e.Allowed("PV-DR"))
	require.True(t, e.Allowed("PV-dr"))
	require.False(t, e.Allowed(""))
	require.False(t, e.Allowed("abc"))
}

func foo(req *restful.Request, resp *restful.Response) {
	_, _ = io.WriteString(resp.ResponseWriter, "foo")
}

func TestAllowedPathSuffixes(t *testing.T) {
	e := service.NewTenantEnsurer([]string{"a", "b", "c"}, []string{"/health", "/liveliness"})
	ws := new(restful.WebService).Path("")
	ws.Filter(e.EnsureAllowedTenantFilter)
	health := ws.GET("/health").To(foo)
	liveliness := ws.GET("/liveliness").To(foo)
	machine := ws.GET("/machine").To(foo)
	ws.Route(health)
	ws.Route(liveliness)
	ws.Route(machine)
	restful.DefaultContainer.Add(ws)

	// health must be allowed without tenant check
	httpRequest, _ := http.NewRequest("GET", "http://localhost/health", nil)
	httpRequest.Header.Set("Accept", "*/*")
	httpWriter := httptest.NewRecorder()

	restful.DefaultContainer.Dispatch(httpWriter, httpRequest)

	require.Equal(t, http.StatusOK, httpWriter.Code)

	// liveliness must be allowed without tenant check
	httpRequest, _ = http.NewRequest("GET", "http://localhost/liveliness", nil)
	httpRequest.Header.Set("Accept", "*/*")
	httpWriter = httptest.NewRecorder()

	restful.DefaultContainer.Dispatch(httpWriter, httpRequest)

	require.Equal(t, http.StatusOK, httpWriter.Code)

	// machine must not be allowed without tenant check
	httpRequest, _ = http.NewRequest("GET", "http://localhost/machine", nil)
	httpRequest.Header.Set("Accept", "*/*")
	httpWriter = httptest.NewRecorder()

	restful.DefaultContainer.Dispatch(httpWriter, httpRequest)

	require.Equal(t, http.StatusForbidden, httpWriter.Code)
}
