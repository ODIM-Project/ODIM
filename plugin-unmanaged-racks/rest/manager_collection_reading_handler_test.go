package rest

import (
	"net/http"
	"testing"

	"github.com/kataras/iris/v12/httptest"
)

func Test_get_manager_collection(t *testing.T) {
	testApp, _ := createTestApplication()
	httptest.New(t, testApp).
		GET("/ODIM/v1/Managers").
		WithBasicAuth("admin", "Od!m12$4").
		Expect().
		Status(http.StatusOK).
		ContentType("application/json", "UTF-8").
		JSON().Object().
		ValueEqual("Members@odata.count", 1).
		Path("$.Members").Array().Length().Equal(1)
}
