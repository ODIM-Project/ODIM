package rest

import (
	"testing"

	"github.com/kataras/iris/v12/httptest"
	"github.com/stretchr/testify/require"
)

func Test_get_not_empty_chassis_collection(t *testing.T) {
	testApp, testRedis := createTestApplication()
	//should be returned
	require.NoError(t, testRedis.Set("Chassis:/ODIM/v1/Chassis/1", ""))
	require.NoError(t, testRedis.Set("Chassis:/ODIM/v1/Chassis/2", ""))
	//should not be returned
	require.NoError(t, testRedis.Set("CONTAINS:Chassis:/ODIM/v1/Chassis/2", ""))
	require.NoError(t, testRedis.Set("CONTAINEDIN:Chassis:/ODIM/v1/Chassis/2", ""))

	httptest.New(t, testApp).
		GET("/ODIM/v1/Chassis/").
		WithBasicAuth("admin", "Od!m12$4").
		Expect().Status(httptest.StatusOK).
		ContentType("application/json", "UTF-8").
		JSON().Object().
		Path(`$.Members..["@odata.id"]`).Array().ContainsOnly("/ODIM/v1/Chassis/1", "/ODIM/v1/Chassis/2")
}
