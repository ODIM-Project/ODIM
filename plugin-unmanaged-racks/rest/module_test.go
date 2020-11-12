package rest

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/config"
	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/db"

	"github.com/alicebob/miniredis/v2"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/httptest"
	"github.com/stretchr/testify/require"
)

var TEST_CONFIG = config.PluginConfig{
	RootServiceUUID: "99999999-9999-9999-9999-999999999999",
	UserName:        "admin",
	Password:        "O01bKrP7Tzs7YoO3YvQt4pRa2J_R6HI34ZfP4MxbqNIYAVQVt2ewGXmhjvBfzMifM7bHFccXKGmdHvj3hY44Hw==",
	EventConf: &config.EventConf{
		DestURI: "/eventHandler",
	},
}

func createTestApplication() (*iris.Application, *miniredis.Miniredis) {
	r, err := miniredis.Run()
	if err != nil {
		panic(err)
	}

	cm := db.NewConnectionManager("tcp", r.Host(), r.Port())
	return createApplication(&TEST_CONFIG, cm), r
}

func Test_secured_endpoints_return_401_when_unauthorized(t *testing.T) {
	tests := []struct {
		method string
		uri    string
	}{
		{http.MethodPost, "/ODIM/v1/Startup"},

		{http.MethodGet, "/ODIM/v1/Chassis"},
		{http.MethodPost, "/ODIM/v1/Chassis"},
		{http.MethodGet, "/ODIM/v1/Chassis/1"},
		{http.MethodDelete, "/ODIM/v1/Chassis/1"},
		{http.MethodPatch, "/ODIM/v1/Chassis/1"},

		{http.MethodGet, "/ODIM/v1/Managers"},
		{http.MethodGet, "/ODIM/v1/Managers/1"},
	}

	testApp, _ := createTestApplication()
	for _, test := range tests {
		t.Run(fmt.Sprintf("%s %s", test.method, test.uri), func(t *testing.T) {
			he := httptest.New(t, testApp)
			resp := he.Request(test.method, test.uri).Expect()
			resp.Status(http.StatusUnauthorized)
		})
	}
}

func Test_unsecured_endpoints_return_401_when_unauthorized(t *testing.T) {
	tests := []struct {
		method string
		uri    string
	}{
		{http.MethodGet, "/ODIM/v1/Status"},
		{http.MethodPost, TEST_CONFIG.EventConf.DestURI},
	}

	testApp, _ := createTestApplication()
	for _, test := range tests {
		t.Run(fmt.Sprintf("%s %s", test.method, test.uri), func(t *testing.T) {
			he := httptest.New(t, testApp)
			resp := he.Request(test.method, test.uri).Expect()
			require.NotEqual(t, http.StatusUnauthorized, resp.Raw().Status)
		})
	}
}
