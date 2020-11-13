package rest

import (
	"github.com/kataras/iris/v12/httptest"
	"net/http"
	"testing"
)

func Test_get_urp_manager(t *testing.T) {
	testApp, _ := createTestApplication()
	httptest.New(t, testApp).
		GET("/ODIM/v1/Managers/"+TEST_CONFIG.RootServiceUUID).
		WithBasicAuth("admin", "Od!m12$4").
		Expect().
		Status(http.StatusOK).
		ContentType("application/json", "UTF-8").
		JSON().Object().
		ValueEqual("@odata.id", "/ODIM/v1/Managers/"+TEST_CONFIG.RootServiceUUID).
		ValueEqual("Name", TEST_CONFIG.ID).
		ValueEqual("UUID", TEST_CONFIG.RootServiceUUID).
		ValueEqual("Id", TEST_CONFIG.RootServiceUUID).
		ValueEqual("FirmwareVersion", TEST_CONFIG.FirmwareVersion)
}
