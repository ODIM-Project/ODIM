//(C) Copyright [2020] Hewlett Packard Enterprise Development LP
//
//Licensed under the Apache License, Version 2.0 (the "License"); you may
//not use this file except in compliance with the License. You may obtain
//a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//License for the specific language governing permissions and limitations
// under the License.
package handle

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-api/models"
	iris "github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/httptest"
)

//TestGetVersion is unittest method for GetVersion func.
func TestGetVersion(t *testing.T) {
	router := iris.New()
	redfishRoutes := router.Party("/redfish")
	redfishRoutes.Get("/", GetVersion)
	e := httptest.New(t, router)

	//Expected reponse body decalration and initilaization to string
	expectedBody := "{\n  \"v1\": \"/redfish/v1/\"\n}\n"

	//Check for status code 200 which is StatusOK
	e.GET("/redfish").Expect().Status(http.StatusOK)

	//Check for the response body which should be equal to the expextecBody
	e.GET("/redfish").Expect().Status(http.StatusOK).Body().Equal(expectedBody)
}

func mockGetService(a []string, b string) models.ServiceRoot {
	return models.ServiceRoot{}
}

//TestGetServiceRoot is unittest method for GetServiceRoot func.
func TestGetServiceRoot(t *testing.T) {
	s := ServiceRoot{getService: mockGetService}

	router := iris.New()
	redfishRoutes := router.Party("/redfish")

	redfishRoutes.Get("/v1", s.GetServiceRoot)
	e := httptest.New(t, router)

	//Check for status code 200 which is StatusOK
	e.GET("/redfish/v1").Expect().Status(http.StatusOK)
}

//TestGetOdata is unittest method for GetOdata func.
func TestGetOdata(t *testing.T) {
	router := iris.New()
	redfishRoutes := router.Party("/redfish")
	redfishRoutes.Get("/v1/odata", GetOdata)
	e := httptest.New(t, router)

	//Check for status code 200 which is StatusOK
	e.GET("/redfish/v1/odata").Expect().Status(http.StatusOK)

	list := [4]string{"@odata.context", "value", "@Redfish.Copyright", "Session"}

	//Check if body contains the fileds mentioned in list.
	for _, field := range list {
		e.GET("/redfish/v1/odata").Expect().Status(http.StatusOK).Body().Contains(field)
	}

}

//TestGetMetadata is unittest method for GetOdata func.
func TestGetMetadata(t *testing.T) {
	router := iris.New()
	redfishRoutes := router.Party("/redfish")
	redfishRoutes.Get("/v1/$metadata", GetMetadata)
	e := httptest.New(t, router)

	//Check for status code 200 which is StatusOK
	e.GET("/redfish/v1/$metadata").Expect().Status(http.StatusOK)

	list := [4]string{"Reference", "Uri", "Namespace", "Include"}

	//Check if body contains the fileds mentioned in list.
	for _, field := range list {
		e.GET("/redfish/v1/$metadata").Expect().Status(http.StatusOK).Body().Contains(field)
	}

}

//TestAsMethodNotAllowed is unittest method for AsMethodNotAllowed func.
func TestAsMethodNotAllowed(t *testing.T) {
	router := iris.New()
	redfishRoutes := router.Party("/redfish")
	redfishRoutes.Any("/v1/AccountService", AsMethodNotAllowed)
	e := httptest.New(t, router)

	//Check for status code 405 for http methods which are not allowed on Account service URL
	e.POST("/redfish/v1/AccountService").Expect().Status(http.StatusMethodNotAllowed)
	e.PUT("/redfish/v1/AccountService").Expect().Status(http.StatusMethodNotAllowed)
	e.DELETE("/redfish/v1/AccountService").Expect().Status(http.StatusMethodNotAllowed)
}

//TestSsMethodNotAllowed is unittest method for SsMethodNotAllowed func.
func TestSsMethodNotAllowed(t *testing.T) {
	router := iris.New()
	redfishRoutes := router.Party("/redfish")
	redfishRoutes.Any("/v1/SessionService", SsMethodNotAllowed)
	e := httptest.New(t, router)

	//Check for status code 405 for http methods which are not allowed on Account service URL
	e.POST("/redfish/v1/SessionService").Expect().Status(http.StatusMethodNotAllowed)
	e.PUT("/redfish/v1/SessionService").Expect().Status(http.StatusMethodNotAllowed)
	e.DELETE("/redfish/v1/SessionService").Expect().Status(http.StatusMethodNotAllowed)
}

//TestSystemsMethodNotAllowed is unittest method for SystemsMethodNotAllowed func.
func TestSystemsMethodNotAllowed(t *testing.T) {
	router := iris.New()
	redfishRoutes := router.Party("/redfish")
	redfishRoutes.Any("/v1/Systems", SystemsMethodNotAllowed)
	redfishRoutes.Any("/v1/Systems/{id}", SystemsMethodNotAllowed)
	redfishRoutes.Any("/v1/Systems/{id}/EthernetInterfaces", SystemsMethodNotAllowed)
	redfishRoutes.Any("/v1/Systems/{id}/EthernetInterfaces/{rid}", SystemsMethodNotAllowed)
	redfishRoutes.Any("/v1/Systems/{id}/Memory", SystemsMethodNotAllowed)
	redfishRoutes.Any("/v1/Systems/{id}/Processors", SystemsMethodNotAllowed)
	redfishRoutes.Any("/v1/Systems/{id}/Storage", SystemsMethodNotAllowed)
	redfishRoutes.Any("/v1/Systems/{id}/Storage/{rid}/Drives/{rid2}", SystemsMethodNotAllowed)
	redfishRoutes.Any("/v1/Systems/{id}/Storage/{rid}", SystemsMethodNotAllowed)
	redfishRoutes.Any("/v1/Systems/{id}/Storage/{rid}/Volumes", SystemsMethodNotAllowed)
	redfishRoutes.Any("/v1/Systems/{id}/Processors/{rid}", SystemsMethodNotAllowed)
	redfishRoutes.Any("/v1/Systems/{id}/Storage/{rid}/Volumes/{rid2}", SystemsMethodNotAllowed)

	e := httptest.New(t, router)
	systemID := "74116e00-0a4a-53e6-a959-e6a7465d6358:1"
	rID := "1"

	//Check for status code 405 for http methods which are not allowed on systems URLs
	e.POST("/redfish/v1/Systems").Expect().Status(http.StatusMethodNotAllowed)
	e.PUT("/redfish/v1/Systems").Expect().Status(http.StatusMethodNotAllowed)
	e.PATCH("/redfish/v1/Systems").Expect().Status(http.StatusMethodNotAllowed)
	e.DELETE("/redfish/v1/Systems").Expect().Status(http.StatusMethodNotAllowed)

	e.POST("/redfish/v1/Systems/" + systemID).Expect().Status(http.StatusMethodNotAllowed)
	e.PUT("/redfish/v1/Systems/" + systemID).Expect().Status(http.StatusMethodNotAllowed)
	e.DELETE("/redfish/v1/Systems/" + systemID).Expect().Status(http.StatusMethodNotAllowed)

	e.POST("/redfish/v1/Systems/" + systemID + "/EthernetInterfaces").Expect().Status(http.StatusMethodNotAllowed)
	e.PUT("/redfish/v1/Systems/" + systemID + "/EthernetInterfaces").Expect().Status(http.StatusMethodNotAllowed)
	e.PATCH("/redfish/v1/Systems/" + systemID + "/EthernetInterfaces").Expect().Status(http.StatusMethodNotAllowed)
	e.DELETE("/redfish/v1/Systems/" + systemID + "/EthernetInterfaces").Expect().Status(http.StatusMethodNotAllowed)

	e.POST("/redfish/v1/Systems/" + systemID + "/EthernetInterfaces/" + rID).Expect().Status(http.StatusMethodNotAllowed)
	e.PUT("/redfish/v1/Systems/" + systemID + "/EthernetInterfaces/" + rID).Expect().Status(http.StatusMethodNotAllowed)
	e.PATCH("/redfish/v1/Systems/" + systemID + "/EthernetInterfaces/" + rID).Expect().Status(http.StatusMethodNotAllowed)
	e.DELETE("/redfish/v1/Systems/" + systemID + "/EthernetInterfaces/" + rID).Expect().Status(http.StatusMethodNotAllowed)

	e.POST("/redfish/v1/Systems/" + systemID + "/Memory").Expect().Status(http.StatusMethodNotAllowed)
	e.PUT("/redfish/v1/Systems/" + systemID + "/Memory").Expect().Status(http.StatusMethodNotAllowed)
	e.PATCH("/redfish/v1/Systems/" + systemID + "/Memory").Expect().Status(http.StatusMethodNotAllowed)
	e.DELETE("/redfish/v1/Systems/" + systemID + "/Memory").Expect().Status(http.StatusMethodNotAllowed)

	e.POST("/redfish/v1/Systems/" + systemID + "/Processors").Expect().Status(http.StatusMethodNotAllowed)
	e.PUT("/redfish/v1/Systems/" + systemID + "/Processors").Expect().Status(http.StatusMethodNotAllowed)
	e.PATCH("/redfish/v1/Systems/" + systemID + "/Processors").Expect().Status(http.StatusMethodNotAllowed)
	e.DELETE("/redfish/v1/Systems/" + systemID + "/Processors").Expect().Status(http.StatusMethodNotAllowed)

	e.POST("/redfish/v1/Systems/" + systemID + "/Storage").Expect().Status(http.StatusMethodNotAllowed)
	e.PUT("/redfish/v1/Systems/" + systemID + "/Storage").Expect().Status(http.StatusMethodNotAllowed)
	e.PATCH("/redfish/v1/Systems/" + systemID + "/Storage").Expect().Status(http.StatusMethodNotAllowed)
	e.DELETE("/redfish/v1/Systems/" + systemID + "/Storage").Expect().Status(http.StatusMethodNotAllowed)

	e.POST("/redfish/v1/Systems/" + systemID + "/Storage/{rid}/Drives/{rid2}").Expect().Status(http.StatusMethodNotAllowed)
	e.PUT("/redfish/v1/Systems/" + systemID + "/Storage/{rid}/Drives/{rid2}").Expect().Status(http.StatusMethodNotAllowed)
	e.PATCH("/redfish/v1/Systems/" + systemID + "/Storage/{rid}/Drives/{rid2}").Expect().Status(http.StatusMethodNotAllowed)
	e.DELETE("/redfish/v1/Systems/" + systemID + "/Storage/{rid}/Drives/{rid2}").Expect().Status(http.StatusMethodNotAllowed)

	e.POST("/redfish/v1/Systems/" + systemID + "/Storage/{rid}").Expect().Status(http.StatusMethodNotAllowed)
	e.PUT("/redfish/v1/Systems/" + systemID + "/Storage/{rid}").Expect().Status(http.StatusMethodNotAllowed)
	e.PATCH("/redfish/v1/Systems/" + systemID + "/Storage/{rid}").Expect().Status(http.StatusMethodNotAllowed)
	e.DELETE("/redfish/v1/Systems/" + systemID + "/Storage/{rid}").Expect().Status(http.StatusMethodNotAllowed)

	e.POST("/redfish/v1/Systems/" + systemID + "/Storage/{rid}/Volumes/{rid2}").Expect().Status(http.StatusMethodNotAllowed)
	e.PUT("/redfish/v1/Systems/" + systemID + "/Storage/{rid}/Volumes/{rid2}").Expect().Status(http.StatusMethodNotAllowed)
	e.PATCH("/redfish/v1/Systems/" + systemID + "/Storage/{rid}/Volumes/{rid2}").Expect().Status(http.StatusMethodNotAllowed)

	e.PUT("/redfish/v1/Systems/" + systemID + "/Storage/{rid}/Volumes").Expect().Status(http.StatusMethodNotAllowed)
	e.PATCH("/redfish/v1/Systems/" + systemID + "/Storage/{rid}/Volumes").Expect().Status(http.StatusMethodNotAllowed)
	e.DELETE("/redfish/v1/Systems/" + systemID + "/Storage/{rid}/Volumes").Expect().Status(http.StatusMethodNotAllowed)

	e.POST("/redfish/v1/Systems/" + systemID + "/Processors/{rid}").Expect().Status(http.StatusMethodNotAllowed)
	e.PUT("/redfish/v1/Systems/" + systemID + "/Processors/{rid}").Expect().Status(http.StatusMethodNotAllowed)
	e.PATCH("/redfish/v1/Systems/" + systemID + "/Processors/{rid}").Expect().Status(http.StatusMethodNotAllowed)
	e.DELETE("/redfish/v1/Systems/" + systemID + "/Processors/{rid}").Expect().Status(http.StatusMethodNotAllowed)
}

//TestMethodNotAllowedForLogServices is unit test method for
//LogService path in ManagersMethodNotAllowed and SystemsMethodNotAllowed funcs.
func TestMethodNotAllowedForLogServices(t *testing.T) {
	logServicesURI := "{id}/LogServices/{rID}"
	entriesURI := logServicesURI + "/Entries"
	subEntriesURI := logServicesURI + "/Entries/{rID2}"
	actionsURI := logServicesURI + "/Actions"
	clearLogURI := logServicesURI + "/Actions/LogService.ClearLog"

	router := iris.New()
	systemsRoutes := router.Party("/redfish/v1/Systems")
	systemsRoutes.Any("{id}/LogServices", SystemsMethodNotAllowed)
	systemsRoutes.Any(logServicesURI, SystemsMethodNotAllowed)
	systemsRoutes.Any(entriesURI, SystemsMethodNotAllowed)
	systemsRoutes.Any(subEntriesURI, SystemsMethodNotAllowed)
	systemsRoutes.Any(actionsURI, SystemsMethodNotAllowed)
	systemsRoutes.Any(clearLogURI, SystemsMethodNotAllowed)
	managersRoutes := router.Party("/redfish/v1/Managers")
	managersRoutes.Any("{id}/LogServices", ManagersMethodNotAllowed)
	managersRoutes.Any(logServicesURI, ManagersMethodNotAllowed)
	managersRoutes.Any(entriesURI, ManagersMethodNotAllowed)
	managersRoutes.Any(subEntriesURI, ManagersMethodNotAllowed)
	managersRoutes.Any(actionsURI, ManagersMethodNotAllowed)
	managersRoutes.Any(clearLogURI, ManagersMethodNotAllowed)

	e := httptest.New(t, router)

	for _, module := range []string{"/redfish/v1/Systems", "/redfish/v1/Managers"} {
		uri := module + "/23256e00-0a4a-53e6-a959-e6a7465d2325:1/LogServices"
		func(uri string) {
			uriForRid := uri + "/1"
			uriForEntries := uriForRid + "/Entries"
			uriForSubEntries := uriForRid + "/Entries/1"
			uriForActions := uriForRid + "/Actions"
			uriForClearLog := uriForRid + "/Actions/LogService.ClearLog"

			e.GET(uriForActions).Expect().Status(http.StatusMethodNotAllowed)
			e.GET(uriForClearLog).Expect().Status(http.StatusMethodNotAllowed)

			e.PUT(uri).Expect().Status(http.StatusMethodNotAllowed)
			e.PUT(uriForRid).Expect().Status(http.StatusMethodNotAllowed)
			e.PUT(uriForEntries).Expect().Status(http.StatusMethodNotAllowed)
			e.PUT(uriForSubEntries).Expect().Status(http.StatusMethodNotAllowed)
			e.PUT(uriForActions).Expect().Status(http.StatusMethodNotAllowed)
			e.PUT(uriForClearLog).Expect().Status(http.StatusMethodNotAllowed)

			e.POST(uri).Expect().Status(http.StatusMethodNotAllowed)
			e.POST(uriForRid).Expect().Status(http.StatusMethodNotAllowed)
			e.POST(uriForEntries).Expect().Status(http.StatusMethodNotAllowed)
			e.POST(uriForSubEntries).Expect().Status(http.StatusMethodNotAllowed)
			e.POST(uriForActions).Expect().Status(http.StatusMethodNotAllowed)

			e.PATCH(uri).Expect().Status(http.StatusMethodNotAllowed)
			e.PATCH(uriForRid).Expect().Status(http.StatusMethodNotAllowed)
			e.PATCH(uriForEntries).Expect().Status(http.StatusMethodNotAllowed)
			e.PATCH(uriForSubEntries).Expect().Status(http.StatusMethodNotAllowed)
			e.PATCH(uriForActions).Expect().Status(http.StatusMethodNotAllowed)
			e.PATCH(uriForClearLog).Expect().Status(http.StatusMethodNotAllowed)

			e.DELETE(uri).Expect().Status(http.StatusMethodNotAllowed)
			e.DELETE(uriForRid).Expect().Status(http.StatusMethodNotAllowed)
			e.DELETE(uriForEntries).Expect().Status(http.StatusMethodNotAllowed)
			e.DELETE(uriForSubEntries).Expect().Status(http.StatusMethodNotAllowed)
			e.DELETE(uriForActions).Expect().Status(http.StatusMethodNotAllowed)
			e.DELETE(uriForClearLog).Expect().Status(http.StatusMethodNotAllowed)
		}(uri)
	}
}
func authMock(token string, b []string, c []string) response.RPC {
	if token == "invalidToken" {
		return common.GeneralError(http.StatusUnauthorized, response.NoValidSession, "", nil, nil)
	}
	return common.GeneralError(http.StatusOK, response.Success, "", nil, nil)
}

func TestGetRegistryFileCollection(t *testing.T) {
	err := common.SetUpMockConfig()
	if err != nil {
		t.Fatalf("fatal: error while trying to collect mock db config: %v", err)
		return
	}
	r := Registry{
		Auth: authMock,
	}
	router := iris.New()
	redfishRoutes := router.Party("/redfish/v1")
	redfishRoutes.Get("/Registries", r.GetRegistryFileCollection)
	test := httptest.New(t, router)
	test.GET("/redfish/v1/Registries").WithHeader("X-Auth-Token", "validToken").Expect().Status(http.StatusOK)
	test.GET("/redfish/v1/Registries").Expect().Status(http.StatusUnauthorized)
	test.GET("/redfish/v1/Registries").WithHeader("X-Auth-Token", "invalidToken").Expect().Status(http.StatusUnauthorized)
}
func TestGetMessageRegistryFileID(t *testing.T) {
	err := common.SetUpMockConfig()
	if err != nil {
		t.Fatalf("fatal: error while trying to collect mock db config: %v", err)
		return
	}
	r := Registry{
		Auth: authMock,
	}
	message := []byte("Just Testing")
	err = ioutil.WriteFile("/tmp/Base.1.10.0.json", message, 0644)
	if err != nil {
		t.Fatalf(err.Error())
	}
	router := iris.New()
	redfishRoutes := router.Party("/redfish/v1")
	redfishRoutes.Get("/Registries/{id}", r.GetMessageRegistryFileID)
	test := httptest.New(t, router)
	test.GET("/redfish/v1/Registries/UnknownID").WithHeader("X-Auth-Token", "validToken").Expect().Status(http.StatusNotFound)
	test.GET("/redfish/v1/Registries/Base.1.10.0").WithHeader("X-Auth-Token", "validToken").Expect().Status(http.StatusOK)
	test.GET("/redfish/v1/Registries/Base.1.10.0").Expect().Status(http.StatusUnauthorized)
	test.GET("/redfish/v1/Registries/Base.1.10.0").WithHeader("X-Auth-Token", "invalidToken").Expect().Status(http.StatusUnauthorized)
}
func TestGetMessageRegistryFile(t *testing.T) {
	err := common.SetUpMockConfig()
	if err != nil {
		t.Fatalf("fatal: error while trying to collect mock db config: %v", err)
		return
	}
	r := Registry{
		Auth: authMock,
	}
	message := []byte("Just Testing")
	err = ioutil.WriteFile("/tmp/Base.1.10.0.json", message, 0644)
	if err != nil {
		t.Fatalf(err.Error())
	}
	router := iris.New()
	redfishRoutes := router.Party("/redfish/v1")
	redfishRoutes.Get("/registries/{id}", r.GetMessageRegistryFile)
	test := httptest.New(t, router)
	test.GET("/redfish/v1/registries/UnknownID").WithHeader("X-Auth-Token", "validToken").Expect().Status(http.StatusNotFound)
	test.GET("/redfish/v1/registries/Base.1.10.0.json").WithHeader("X-Auth-Token", "validToken").Expect().Status(http.StatusOK)
	test.GET("/redfish/v1/registries/Base.1.10.0.json").Expect().Status(http.StatusUnauthorized)
	test.GET("/redfish/v1/registries/Base.1.10.0.json").WithHeader("X-Auth-Token", "invalidToken").Expect().Status(http.StatusUnauthorized)
}

//TestTsMethodNotAllowed is unittest method for TsMethodNotAllowed func.
func TestTsMethodNotAllowed(t *testing.T) {
	router := iris.New()
	redfishRoutes := router.Party("/redfish/v1")
	redfishRoutes.Any("/TaskService", TsMethodNotAllowed)
	redfishRoutes.Any("/TaskService/Tasks", TsMethodNotAllowed)
	redfishRoutes.Any("/TaskService/Tasks/{TaskID}", TsMethodNotAllowed)
	e := httptest.New(t, router)

	//Check for status code 405 for http methods which are not allowed on Task service URLs
	e.POST("/redfish/v1/TaskService").Expect().Status(http.StatusMethodNotAllowed)
	e.PUT("/redfish/v1/TaskService").Expect().Status(http.StatusMethodNotAllowed)
	e.DELETE("/redfish/v1/TaskService").Expect().Status(http.StatusMethodNotAllowed)

	e.POST("/redfish/v1/TaskService/Tasks").Expect().Status(http.StatusMethodNotAllowed)
	e.PUT("/redfish/v1/TaskService/Tasks").Expect().Status(http.StatusMethodNotAllowed)
	e.DELETE("/redfish/v1/TaskService/Tasks").Expect().Status(http.StatusMethodNotAllowed)

	e.POST("/redfish/v1/TaskService/Tasks/{TaskID}").Expect().Status(http.StatusMethodNotAllowed)
	e.PUT("/redfish/v1/TaskService/Tasks/{TaskID}").Expect().Status(http.StatusMethodNotAllowed)
}

//TestEvtMethodNotAllowed is unittest method for EvtMethodNotAllowed func.
func TestEvtMethodNotAllowed(t *testing.T) {
	router := iris.New()
	redfishRoutes := router.Party("/redfish/v1")
	redfishRoutes.Any("/EventService", EvtMethodNotAllowed)
	redfishRoutes.Any("/EventService/Actions", EvtMethodNotAllowed)
	redfishRoutes.Any("/EventService/Actions/EventService.SubmitTestEvent", EvtMethodNotAllowed)
	redfishRoutes.Any("/EventService/Subscriptions/", EvtMethodNotAllowed)
	e := httptest.New(t, router)

	//Check for status code 405 for http methods which are not allowed on Task service URLs
	e.POST("/redfish/v1/EventService").Expect().Status(http.StatusMethodNotAllowed)
	e.PUT("/redfish/v1/EventService").Expect().Status(http.StatusMethodNotAllowed)
	e.DELETE("/redfish/v1/EventService").Expect().Status(http.StatusMethodNotAllowed)

	e.GET("/redfish/v1/EventService/Actions/EventService.SubmitTestEvent").Expect().Status(http.StatusMethodNotAllowed)
	e.PUT("/redfish/v1/EventService/Actions/EventService.SubmitTestEvent").Expect().Status(http.StatusMethodNotAllowed)
	e.DELETE("/redfish/v1/EventService/Actions/EventService.SubmitTestEvent").Expect().Status(http.StatusMethodNotAllowed)

	e.DELETE("/redfish/v1/EventService/Subscriptions").Expect().Status(http.StatusMethodNotAllowed)
	e.PUT("/redfish/v1/EventService/Subscriptions").Expect().Status(http.StatusMethodNotAllowed)
	e.PATCH("/redfish/v1/EventService/Subscriptions").Expect().Status(http.StatusMethodNotAllowed)
}

//TestAggMethodNotAllowed is unittest method for AggMethodNotAllowed func.
func TestAggMethodNotAllowed(t *testing.T) {
	router := iris.New()
	redfishRoutes := router.Party("/redfish/v1")
	redfishRoutes.Any("/AggregationService", AggMethodNotAllowed)
	redfishRoutes.Any("/AggregationService/ConnectionMethods", AggMethodNotAllowed)
	redfishRoutes.Any("/AggregationService/ConnectionMethods/{id}", AggMethodNotAllowed)
	e := httptest.New(t, router)

	//Check for status code 405 for http methods which are not allowed on aggregation servicee URLs
	e.POST("/redfish/v1/AggregationService").Expect().Status(http.StatusMethodNotAllowed)
	e.PUT("/redfish/v1/AggregationService").Expect().Status(http.StatusMethodNotAllowed)
	e.DELETE("/redfish/v1/AggregationService").Expect().Status(http.StatusMethodNotAllowed)

	//Check for status code 405 for http methods which are not allowed on aggregation service connection methods URLs
	e.POST("/redfish/v1/AggregationService/ConnectionMethods").Expect().Status(http.StatusMethodNotAllowed)
	e.PUT("/redfish/v1/AggregationService/ConnectionMethods").Expect().Status(http.StatusMethodNotAllowed)
	e.DELETE("/redfish/v1/AggregationService/ConnectionMethods").Expect().Status(http.StatusMethodNotAllowed)

	connMethodID := "74116e00-0a4a-53e6-a959-e6a7465d6358"
	//Check for status code 405 for http methods which are not allowed on aggregation service connection method URLs
	e.POST("/redfish/v1/AggregationService/ConnectionMethods/" + connMethodID).Expect().Status(http.StatusMethodNotAllowed)
	e.PUT("/redfish/v1/AggregationService/ConnectionMethods/" + connMethodID).Expect().Status(http.StatusMethodNotAllowed)
	e.DELETE("/redfish/v1/AggregationService/ConnectionMethods/" + connMethodID).Expect().Status(http.StatusMethodNotAllowed)
}

//TestFabricsMethodNotAllowed is unittest method for FabricsMethodNotAllowed func.
func TestFabricsMethodNotAllowed(t *testing.T) {
	router := iris.New()
	redfishRoutes := router.Party("/redfish/v1")
	redfishRoutes.Any("/Fabrics", FabricsMethodNotAllowed)
	e := httptest.New(t, router)

	//Check for status code 405 for http methods which are not allowed on Task service URLs
	e.POST("/redfish/v1/Fabrics").Expect().Status(http.StatusMethodNotAllowed)
	e.PUT("/redfish/v1/Fabrics").Expect().Status(http.StatusMethodNotAllowed)
	e.DELETE("/redfish/v1/Fabrics").Expect().Status(http.StatusMethodNotAllowed)
}

//TestChassisMethodNotAllowed is unittest method for ChassisMethodNotAllowed func.
func TestChassisMethodNotAllowed(t *testing.T) {
	router := iris.New()
	redfishRoutes := router.Party("/redfish")
	redfishRoutes.Any("/v1/Chassis", ChassisMethodNotAllowed)
	redfishRoutes.Any("/v1/Chassis/{id}/NetworkAdapters", ChassisMethodNotAllowed)
	redfishRoutes.Any("/v1/Chassis/{id}/Power#PowerControl/{rid}", ChassisMethodNotAllowed)
	redfishRoutes.Any("/v1/Chassis/{id}/Power#PowerSupplies/{rid}", ChassisMethodNotAllowed)
	redfishRoutes.Any("/v1/Chassis/{id}/Power#Redundancy/{rid}", ChassisMethodNotAllowed)
	redfishRoutes.Any("/v1/Chassis/{id}/Thermal#Fans/{rid}", ChassisMethodNotAllowed)
	redfishRoutes.Any("/v1/Chassis/{id}/Thermal#Temperatures/{rid}", ChassisMethodNotAllowed)
	redfishRoutes.Any("/v1/Chassis/{id}/Assembly", ChassisMethodNotAllowed)
	redfishRoutes.Any("/v1/Chassis/{id}/PCIeSlots", ChassisMethodNotAllowed)
	redfishRoutes.Any("/v1/Chassis/{id}/PCIeSlots/{rid}", ChassisMethodNotAllowed)
	redfishRoutes.Any("/v1/Chassis/{id}/PCIeDevices", ChassisMethodNotAllowed)
	redfishRoutes.Any("/v1/Chassis/{id}/PCIeDevices/{rid}", ChassisMethodNotAllowed)
	redfishRoutes.Any("/v1/Chassis/{id}/Sensors", ChassisMethodNotAllowed)
	redfishRoutes.Any("/v1/Chassis/{id}/Sensors/{rid}", ChassisMethodNotAllowed)

	e := httptest.New(t, router)
	chassisID := "74116e00-0a4a-53e6-a959-e6a7465d6358:1"
	rID := "1"
	//Check for status code 405 for http methods which are not allowed on systems URLs
	e.POST("/redfish/v1/Chassis").Expect().Status(http.StatusMethodNotAllowed)
	e.PUT("/redfish/v1/Chassis").Expect().Status(http.StatusMethodNotAllowed)
	e.PATCH("/redfish/v1/Chassis").Expect().Status(http.StatusMethodNotAllowed)
	e.DELETE("/redfish/v1/Chassis").Expect().Status(http.StatusMethodNotAllowed)

	e.POST("/redfish/v1/Chassis/" + chassisID + "/NetworkAdapters").Expect().Status(http.StatusMethodNotAllowed)
	e.PUT("/redfish/v1/Chassis/" + chassisID + "/NetworkAdapters").Expect().Status(http.StatusMethodNotAllowed)
	e.PATCH("/redfish/v1/Chassis/" + chassisID + "/NetworkAdapters").Expect().Status(http.StatusMethodNotAllowed)
	e.DELETE("/redfish/v1/Chassis/" + chassisID + "/NetworkAdapters").Expect().Status(http.StatusMethodNotAllowed)

	e.POST("/redfish/v1/Chassis/" + chassisID + "/Power#PowerControl/" + rID).Expect().Status(http.StatusMethodNotAllowed)
	e.PUT("/redfish/v1/Chassis/" + chassisID + "/Power#PowerControl/" + rID).Expect().Status(http.StatusMethodNotAllowed)
	e.PATCH("/redfish/v1/Chassis/" + chassisID + "/Power#PowerControl/" + rID).Expect().Status(http.StatusMethodNotAllowed)
	e.DELETE("/redfish/v1/Chassis/" + chassisID + "/Power#PowerControl/" + rID).Expect().Status(http.StatusMethodNotAllowed)

	e.POST("/redfish/v1/Chassis/" + chassisID + "/Power#PowerSupplies/" + rID).Expect().Status(http.StatusMethodNotAllowed)
	e.PUT("/redfish/v1/Chassis/" + chassisID + "/Power#PowerSupplies/" + rID).Expect().Status(http.StatusMethodNotAllowed)
	e.PATCH("/redfish/v1/Chassis/" + chassisID + "/Power#PowerSupplies/" + rID).Expect().Status(http.StatusMethodNotAllowed)
	e.DELETE("/redfish/v1/Chassis/" + chassisID + "/Power#PowerSupplies/" + rID).Expect().Status(http.StatusMethodNotAllowed)

	e.POST("/redfish/v1/Chassis/" + chassisID + "/Power#Redundancy/" + rID).Expect().Status(http.StatusMethodNotAllowed)
	e.PUT("/redfish/v1/Chassis/" + chassisID + "/Power#Redundancy/" + rID).Expect().Status(http.StatusMethodNotAllowed)
	e.PATCH("/redfish/v1/Chassis/" + chassisID + "/Power#Redundancy/" + rID).Expect().Status(http.StatusMethodNotAllowed)
	e.DELETE("/redfish/v1/Chassis/" + chassisID + "/Power#Redundancy/" + rID).Expect().Status(http.StatusMethodNotAllowed)

	e.POST("/redfish/v1/Chassis/" + chassisID + "/Thermal#Fans/" + rID).Expect().Status(http.StatusMethodNotAllowed)
	e.PUT("/redfish/v1/Chassis/" + chassisID + "/Thermal#Fans/" + rID).Expect().Status(http.StatusMethodNotAllowed)
	e.PATCH("/redfish/v1/Chassis/" + chassisID + "/Thermal#Fans/" + rID).Expect().Status(http.StatusMethodNotAllowed)
	e.DELETE("/redfish/v1/Chassis/" + chassisID + "/Thermal#Fans/" + rID).Expect().Status(http.StatusMethodNotAllowed)

	e.POST("/redfish/v1/Chassis/" + chassisID + "/Thermal#Temperatures/" + rID).Expect().Status(http.StatusMethodNotAllowed)
	e.PUT("/redfish/v1/Chassis/" + chassisID + "/Thermal#Temperatures/" + rID).Expect().Status(http.StatusMethodNotAllowed)
	e.PATCH("/redfish/v1/Chassis/" + chassisID + "/Thermal#Temperatures/" + rID).Expect().Status(http.StatusMethodNotAllowed)
	e.DELETE("/redfish/v1/Chassis/" + chassisID + "/Thermal#Temperatures/" + rID).Expect().Status(http.StatusMethodNotAllowed)

	e.POST("/redfish/v1/Chassis/" + chassisID + "/Assembly").Expect().Status(http.StatusMethodNotAllowed)
	e.PUT("/redfish/v1/Chassis/" + chassisID + "/Assembly").Expect().Status(http.StatusMethodNotAllowed)
	e.PATCH("/redfish/v1/Chassis/" + chassisID + "/Assembly").Expect().Status(http.StatusMethodNotAllowed)
	e.DELETE("/redfish/v1/Chassis/" + chassisID + "/Assembly").Expect().Status(http.StatusMethodNotAllowed)

	e.POST("/redfish/v1/Chassis/" + chassisID + "/PCIeSlots").Expect().Status(http.StatusMethodNotAllowed)
	e.PUT("/redfish/v1/Chassis/" + chassisID + "/PCIeSlots").Expect().Status(http.StatusMethodNotAllowed)
	e.PATCH("/redfish/v1/Chassis/" + chassisID + "/PCIeSlots").Expect().Status(http.StatusMethodNotAllowed)
	e.DELETE("/redfish/v1/Chassis/" + chassisID + "/PCIeSlots").Expect().Status(http.StatusMethodNotAllowed)

	e.POST("/redfish/v1/Chassis/" + chassisID + "/PCIeSlots/" + rID).Expect().Status(http.StatusMethodNotAllowed)
	e.PUT("/redfish/v1/Chassis/" + chassisID + "/PCIeSlots/" + rID).Expect().Status(http.StatusMethodNotAllowed)
	e.PATCH("/redfish/v1/Chassis/" + chassisID + "/PCIeSlots/" + rID).Expect().Status(http.StatusMethodNotAllowed)
	e.DELETE("/redfish/v1/Chassis/" + chassisID + "/PCIeSlots/" + rID).Expect().Status(http.StatusMethodNotAllowed)

	e.POST("/redfish/v1/Chassis/" + chassisID + "/PCIeDevices").Expect().Status(http.StatusMethodNotAllowed)
	e.PUT("/redfish/v1/Chassis/" + chassisID + "/PCIeDevices").Expect().Status(http.StatusMethodNotAllowed)
	e.PATCH("/redfish/v1/Chassis/" + chassisID + "/PCIeDevices").Expect().Status(http.StatusMethodNotAllowed)
	e.DELETE("/redfish/v1/Chassis/" + chassisID + "/PCIeDevices").Expect().Status(http.StatusMethodNotAllowed)

	e.POST("/redfish/v1/Chassis/" + chassisID + "/PCIeDevices/" + rID).Expect().Status(http.StatusMethodNotAllowed)
	e.PUT("/redfish/v1/Chassis/" + chassisID + "/PCIeDevices/" + rID).Expect().Status(http.StatusMethodNotAllowed)
	e.PATCH("/redfish/v1/Chassis/" + chassisID + "/PCIeDevices/" + rID).Expect().Status(http.StatusMethodNotAllowed)
	e.DELETE("/redfish/v1/Chassis/" + chassisID + "/PCIeDevices/" + rID).Expect().Status(http.StatusMethodNotAllowed)

	e.POST("/redfish/v1/Chassis/" + chassisID + "/Sensors").Expect().Status(http.StatusMethodNotAllowed)
	e.PUT("/redfish/v1/Chassis/" + chassisID + "/Sensors").Expect().Status(http.StatusMethodNotAllowed)
	e.PATCH("/redfish/v1/Chassis/" + chassisID + "/Sensors").Expect().Status(http.StatusMethodNotAllowed)
	e.DELETE("/redfish/v1/Chassis/" + chassisID + "/Sensors").Expect().Status(http.StatusMethodNotAllowed)

	e.POST("/redfish/v1/Chassis/" + chassisID + "/Sensors/" + rID).Expect().Status(http.StatusMethodNotAllowed)
	e.PUT("/redfish/v1/Chassis/" + chassisID + "/Sensors/" + rID).Expect().Status(http.StatusMethodNotAllowed)
	e.PATCH("/redfish/v1/Chassis/" + chassisID + "/Sensors/" + rID).Expect().Status(http.StatusMethodNotAllowed)
	e.DELETE("/redfish/v1/Chassis/" + chassisID + "/Sensors/" + rID).Expect().Status(http.StatusMethodNotAllowed)
}

// TestRegMethodNotAllowed is the unit test method for RegMethodNotAllowed func.
func TestRegMethodNotAllowed(t *testing.T) {
	router := iris.New()
	redfishRoutes := router.Party("/redfish")
	redfishRoutes.Any("/v1/Registries", RegMethodNotAllowed)
	redfishRoutes.Any("/v1/Registries/{id}", RegMethodNotAllowed)
	redfishRoutes.Any("/v1/registries", RegMethodNotAllowed)
	redfishRoutes.Any("/v1/registries/{id}", RegMethodNotAllowed)

	e := httptest.New(t, router)
	id := "Base.1.6.0"
	file := "Base.1.6.0.json"

	//Check for status code 405 for http methods which are not allowed on registry URLs
	e.POST("/redfish/v1/Registries").Expect().Status(http.StatusMethodNotAllowed)
	e.PUT("/redfish/v1/Registries").Expect().Status(http.StatusMethodNotAllowed)
	e.PATCH("/redfish/v1/Registries").Expect().Status(http.StatusMethodNotAllowed)
	e.DELETE("/redfish/v1/Registries").Expect().Status(http.StatusMethodNotAllowed)

	e.POST("/redfish/v1/Registries/" + id).Expect().Status(http.StatusMethodNotAllowed)
	e.PUT("/redfish/v1/Registries/" + id).Expect().Status(http.StatusMethodNotAllowed)
	e.PATCH("/redfish/v1/Registries/" + id).Expect().Status(http.StatusMethodNotAllowed)
	e.DELETE("/redfish/v1/Registries/" + id).Expect().Status(http.StatusMethodNotAllowed)

	e.POST("/redfish/v1/registries").Expect().Status(http.StatusMethodNotAllowed)
	e.PUT("/redfish/v1/registries").Expect().Status(http.StatusMethodNotAllowed)
	e.PATCH("/redfish/v1/registries").Expect().Status(http.StatusMethodNotAllowed)
	e.DELETE("/redfish/v1/registries").Expect().Status(http.StatusMethodNotAllowed)

	e.POST("/redfish/v1/registries/" + file).Expect().Status(http.StatusMethodNotAllowed)
	e.PUT("/redfish/v1/registries/" + file).Expect().Status(http.StatusMethodNotAllowed)
	e.PATCH("/redfish/v1/registries/" + file).Expect().Status(http.StatusMethodNotAllowed)
	e.DELETE("/redfish/v1/registries/" + file).Expect().Status(http.StatusMethodNotAllowed)
}

// TestManagersMethodNotAllowed is the unit test method for ManagerMethodNotAllowed func.
func TestManagersMethodNotAllowed(t *testing.T) {
	router := iris.New()
	redfishRoutes := router.Party("/redfish")
	redfishRoutes.Any("/v1/Managers", ManagersMethodNotAllowed)
	redfishRoutes.Any("/v1/Managers/{id}", ManagersMethodNotAllowed)
	e := httptest.New(t, router)

	e.PUT("/redfish/v1/Managers").Expect().Status(http.StatusMethodNotAllowed)
	e.PATCH("/redfish/v1/Managers").Expect().Status(http.StatusMethodNotAllowed)
	e.DELETE("/redfish/v1/Managers").Expect().Status(http.StatusMethodNotAllowed)

	e.PUT("/redfish/v1/Managers/{id}").Expect().Status(http.StatusMethodNotAllowed)
	e.PATCH("/redfish/v1/Managers/{id}").Expect().Status(http.StatusMethodNotAllowed)
	e.DELETE("/redfish/v1/Managers/{id}").Expect().Status(http.StatusMethodNotAllowed)
}

//TestAggregateMethodNotAllowed is unittest method for AggregateMethodNotAllowed func.
func TestAggregateMethodNotAllowed(t *testing.T) {
	router := iris.New()
	redfishRoutes := router.Party("/redfish/v1/AggregationService/Aggregates")
	redfishRoutes.Any("/", AggregateMethodNotAllowed)
	redfishRoutes.Any("/{id}", AggregateMethodNotAllowed)
	redfishRoutes.Any("/{id}/Actions/Aggregate.AddElements/", AggregateMethodNotAllowed)
	redfishRoutes.Any("/{id}/Actions/Aggregate.RemoveElements/", AggregateMethodNotAllowed)
	redfishRoutes.Any("/{id}/Actions/Aggregate.Reset/", AggregateMethodNotAllowed)
	redfishRoutes.Any("/{id}/Actions/Aggregate.SetDefaultBootOrder/", AggregateMethodNotAllowed)

	e := httptest.New(t, router)
	id := "74116e00-0a4a-53e6-a959-e6a7465d6358"
	//Check for status code 405 for http methods which are not allowed
	e.PUT("/redfish/v1/AggregationService/Aggregates").Expect().Status(http.StatusMethodNotAllowed)
	e.PATCH("/redfish/v1/AggregationService/Aggregates").Expect().Status(http.StatusMethodNotAllowed)
	e.DELETE("/redfish/v1/AggregationService/Aggregates").Expect().Status(http.StatusMethodNotAllowed)

	e.POST("/redfish/v1/AggregationService/Aggregates/" + id).Expect().Status(http.StatusMethodNotAllowed)
	e.PUT("/redfish/v1/AggregationService/Aggregates/" + id).Expect().Status(http.StatusMethodNotAllowed)
	e.PATCH("/redfish/v1/AggregationService/Aggregates/" + id).Expect().Status(http.StatusMethodNotAllowed)

	e.GET("/redfish/v1/AggregationService/Aggregates/" + id + "/Actions/Aggregate.AddElements").Expect().Status(http.StatusMethodNotAllowed)
	e.PUT("/redfish/v1/AggregationService/Aggregates/" + id + "/Actions/Aggregate.AddElements").Expect().Status(http.StatusMethodNotAllowed)
	e.PATCH("/redfish/v1/AggregationService/Aggregates/" + id + "/Actions/Aggregate.AddElements").Expect().Status(http.StatusMethodNotAllowed)
	e.DELETE("/redfish/v1/AggregationService/Aggregates/" + id + "/Actions/Aggregate.AddElements").Expect().Status(http.StatusMethodNotAllowed)

	e.GET("/redfish/v1/AggregationService/Aggregates/" + id + "/Actions/Aggregate.RemoveElements").Expect().Status(http.StatusMethodNotAllowed)
	e.PUT("/redfish/v1/AggregationService/Aggregates/" + id + "/Actions/Aggregate.RemoveElements").Expect().Status(http.StatusMethodNotAllowed)
	e.PATCH("/redfish/v1/AggregationService/Aggregates/" + id + "/Actions/Aggregate.RemoveElements").Expect().Status(http.StatusMethodNotAllowed)
	e.DELETE("/redfish/v1/AggregationService/Aggregates/" + id + "/Actions/Aggregate.RemoveElements").Expect().Status(http.StatusMethodNotAllowed)

	e.GET("/redfish/v1/AggregationService/Aggregates/" + id + "/Actions/Aggregate.Reset").Expect().Status(http.StatusMethodNotAllowed)
	e.PUT("/redfish/v1/AggregationService/Aggregates/" + id + "/Actions/Aggregate.Reset").Expect().Status(http.StatusMethodNotAllowed)
	e.PATCH("/redfish/v1/AggregationService/Aggregates/" + id + "/Actions/Aggregate.Reset").Expect().Status(http.StatusMethodNotAllowed)
	e.DELETE("/redfish/v1/AggregationService/Aggregates/" + id + "/Actions/Aggregate.Reset").Expect().Status(http.StatusMethodNotAllowed)

	e.GET("/redfish/v1/AggregationService/Aggregates/" + id + "/Actions/Aggregate.SetDefaultBootOrder").Expect().Status(http.StatusMethodNotAllowed)
	e.PUT("/redfish/v1/AggregationService/Aggregates/" + id + "/Actions/Aggregate.SetDefaultBootOrder").Expect().Status(http.StatusMethodNotAllowed)
	e.PATCH("/redfish/v1/AggregationService/Aggregates/" + id + "/Actions/Aggregate.SetDefaultBootOrder").Expect().Status(http.StatusMethodNotAllowed)
	e.DELETE("/redfish/v1/AggregationService/Aggregates/" + id + "/Actions/Aggregate.SetDefaultBootOrder").Expect().Status(http.StatusMethodNotAllowed)
}
