// Packahe dphandler ...
package dphandler

import (
	"encoding/json"
	"fmt"
	"github.com/ODIM-Project/ODIM/plugin-dell/config"
	"github.com/ODIM-Project/ODIM/plugin-dell/dpresponse"
	iris "github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/httptest"
	"net/http"
	"testing"
)

func mockSimpleUpdate(username, password, url string, w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
}

func TestSimpleUpdate(t *testing.T) {
	config.SetUpMockConfig(t)

	deviceHost := "localhost"
	devicePort := "1234"
	ts := startTestServer(mockSimpleUpdate)
	// Start the server.
	ts.StartTLS()
	defer ts.Close()
	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/ODIM/v1")

	redfishRoutes.Post("/UpdateService/Actions.SimpleUpdate", SimpleUpdate)
	dpresponse.PluginToken = "token"
	test := httptest.New(t, mockApp)
	attributes := map[string]interface{}{"ImageUri": "abc",
		"Targets": []string{"/ODIM/v1/Systems/uuid:1"}}
	attributeByte, _ := json.Marshal(attributes)
	requestBody := map[string]interface{}{
		"ManagerAddress": fmt.Sprintf("%s:%s", deviceHost, devicePort),
		"UserName":       "admin",
		"Password":       []byte("P@$$w0rd"),
		"PostBody":       attributeByte,
	}
	test.POST("/ODIM/v1/UpdateService/Actions.SimpleUpdate").WithJSON(requestBody).Expect().Status(http.StatusOK)
}

func TestStartUpdate(t *testing.T) {
	config.SetUpMockConfig(t)

	deviceHost := "localhost"
	devicePort := "1234"
	ts := startTestServer(mockSimpleUpdate)
	// Start the server.
	ts.StartTLS()
	defer ts.Close()
	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/ODIM/v1")

	redfishRoutes.Post("/UpdateService/Actions.StartUpdate", StartUpdate)
	dpresponse.PluginToken = "token"
	test := httptest.New(t, mockApp)
	requestBody := map[string]interface{}{
		"ManagerAddress": fmt.Sprintf("%s:%s", deviceHost, devicePort),
		"UserName":       "admin",
		"Password":       []byte("P@$$w0rd"),
	}
	test.POST("/ODIM/v1/UpdateService/Actions.StartUpdate").WithJSON(requestBody).Expect().Status(http.StatusOK)
}
