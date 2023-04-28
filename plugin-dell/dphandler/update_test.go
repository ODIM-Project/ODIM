// Packahe dphandler ...
package dphandler

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/ODIM-Project/ODIM/plugin-dell/config"
	"github.com/ODIM-Project/ODIM/plugin-dell/dpresponse"
	iris "github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/httptest"
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
		"Targets": []string{"/ODIM/v1/Systems/uuid.1"}}
	attributeByte, _ := json.Marshal(attributes)
	requestBody := map[string]interface{}{
		"ManagerAddress": fmt.Sprintf("%s:%s", deviceHost, devicePort),
		"UserName":       "admin",
		"Password":       []byte("P@$$w0rd"),
		"PostBody":       attributeByte,
	}
	test.POST("/ODIM/v1/UpdateService/Actions.SimpleUpdate").WithJSON(requestBody).Expect().Status(http.StatusOK)
	test.POST("/ODIM/v1/UpdateService/Actions.SimpleUpdate").WithHeader("X-Auth-Token", "Basic YWRtaW46cGFzc3dvcmQ=").
		WithJSON(requestBody).Expect().Status(http.StatusUnauthorized)
	requestBody1 := "request"
	test.POST("/ODIM/v1/UpdateService/Actions.SimpleUpdate").WithHeader("Authorization", "Basic YWRtaW46cGFzc3dvcmQ=").
		// Invalid PostBody
		WithJSON(requestBody1).Expect().Status(http.StatusBadRequest)
	attributes2 := map[string]interface{}{"ImageUri": 1,
		"Targets": []string{"/ODIM/v1/Systems/uuid.1"}}
	attributeByte2, _ := json.Marshal(attributes2)

	requestBody2 := map[string]interface{}{
		"ManagerAddress": fmt.Sprintf("%s:%s", deviceHost, devicePort),
		"UserName":       "admin",
		"Password":       []byte("P@$$w0rd"),
		"PostBody":       attributeByte2,
	}
	test.POST("/ODIM/v1/UpdateService/Actions.SimpleUpdate").WithHeader("Authorization", "Basic YWRtaW46cGFzc3dvcmQ=").
		WithJSON(requestBody2).Expect().Status(http.StatusOK)

	config.Data.KeyCertConf.RootCACertificate = nil
	// invalid client
	test.POST("/ODIM/v1/UpdateService/Actions.SimpleUpdate").WithJSON(requestBody).Expect().Status(http.StatusInternalServerError)
	config.SetUpMockConfig(t)

	//invalid device details
	IoUtilReadAll = func(r io.Reader) ([]byte, error) {
		return nil, fmt.Errorf("fake error ")
	}
	test.POST("/ODIM/v1/UpdateService/Actions.SimpleUpdate").WithJSON(requestBody).Expect().Status(http.StatusOK)
	IoUtilReadAll = func(r io.Reader) ([]byte, error) {
		return ioutil.ReadAll(r)
	}
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
	// invalid credentials
	test.POST("/ODIM/v1/UpdateService/Actions.StartUpdate").WithHeader("X-Auth-Token", "Basic YWRtaW46cGFzc3dvcmQ=").
		WithJSON(requestBody).Expect().Status(http.StatusUnauthorized)

	// invalid request
	requestBody1 := "invalid"
	test.POST("/ODIM/v1/UpdateService/Actions.StartUpdate").WithJSON(requestBody1).Expect().Status(http.StatusBadRequest)

	// Invalid client request
	config.Data.KeyCertConf.RootCACertificate = nil
	test.POST("/ODIM/v1/UpdateService/Actions.StartUpdate").WithJSON(requestBody).Expect().Status(http.StatusInternalServerError)
	config.SetUpMockConfig(t)

	// Invalid device details
	requestBody2 := map[string]interface{}{
		"ManagerAddress": fmt.Sprintf("%s:%s", "deviceHost", "devicePort"),
		"UserName":       "admin",
		"Password":       "password",
	}
	test.POST("/ODIM/v1/UpdateService/Actions.StartUpdate").WithJSON(requestBody2).Expect().Status(http.StatusInternalServerError)

	IoUtilReadAll = func(r io.Reader) ([]byte, error) {
		return nil, fmt.Errorf("fake error ")
	}
	test.POST("/ODIM/v1/UpdateService/Actions.StartUpdate").WithJSON(requestBody).Expect().Status(http.StatusOK)

}
