package rest

import (
	"github.com/kataras/iris/v12/httptest"
	"testing"
)

func Test_get_not_empty_chassis_collection(t *testing.T) {
	testApp, testRedis := createTestApplication()
	//should be returned
	testRedis.Set("Chassis:/ODIM/v1/Chassis/1", "")
	testRedis.Set("Chassis:/ODIM/v1/Chassis/2", "")
	//should not be returned
	testRedis.Set("CONTAINS:Chassis:/ODIM/v1/Chassis/2", "")
	testRedis.Set("CONTAINEDIN:Chassis:/ODIM/v1/Chassis/2", "")

	httptest.New(t, testApp).
		GET("/ODIM/v1/Chassis/").
		WithBasicAuth("admin", "Od!m12$4").
		Expect().Status(httptest.StatusOK).
		ContentType("application/json", "UTF-8").
		JSON().Object().
		Path(`$.Members..["@odata.id"]`).Array().ContainsOnly("/ODIM/v1/Chassis/1", "/ODIM/v1/Chassis/2")
}
