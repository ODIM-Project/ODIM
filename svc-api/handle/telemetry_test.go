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
	"errors"
	"net/http"
	"testing"

	teleproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/telemetry"
	iris "github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/httptest"
)

func testTelemetryService(req teleproto.TelemetryRequest) (*teleproto.TelemetryResponse, error) {
	var response = &teleproto.TelemetryResponse{}
	if req.SessionToken == "ValidToken" {
		response = &teleproto.TelemetryResponse{
			StatusCode:    200,
			StatusMessage: "Success",
			Body:          []byte(`{"Response":"Success"}`),
		}
	} else if req.SessionToken == "InvalidToken" {
		response = &teleproto.TelemetryResponse{
			StatusCode:    401,
			StatusMessage: "Unauthorized", Body: []byte(`{"Response":"Unauthorized"}`),
		}
	} else if req.SessionToken == "token" {
		return &teleproto.TelemetryResponse{}, errors.New("Unable to RPC Call")
	}
	return response, nil
}

func TestGetTelemetryService(t *testing.T) {
	var a TelemetryRPCs
	a.GetTelemetryServiceRPC = testTelemetryService
	testApp := iris.New()
	redfishRoutes := testApp.Party("/redfish/v1/TelemetryService")
	redfishRoutes.Get("/", a.GetTelemetryService)
	test := httptest.New(t, testApp)
	test.GET(
		"/redfish/v1/TelemetryService",
	).WithHeader("X-Auth-Token", "ValidToken").Expect().Status(http.StatusOK)
	test.GET(
		"/redfish/v1/TelemetryService",
	).WithHeader("X-Auth-Token", "").Expect().Status(http.StatusUnauthorized)
	test.GET(
		"/redfish/v1/TelemetryService",
	).WithHeader("X-Auth-Token", "token").Expect().Status(http.StatusInternalServerError)
}

func TestGetMetricDefinitionCollection(t *testing.T) {
	var a TelemetryRPCs
	a.GetMetricDefinitionCollectionRPC = testTelemetryService
	testApp := iris.New()
	redfishRoutes := testApp.Party("/redfish/v1/TelemetryService")
	redfishRoutes.Get("/MetricDefinitions", a.GetMetricDefinitionCollection)
	test := httptest.New(t, testApp)
	test.GET(
		"/redfish/v1/TelemetryService/MetricDefinitions",
	).WithHeader("X-Auth-Token", "ValidToken").Expect().Status(http.StatusOK)
	test.GET(
		"/redfish/v1/TelemetryService/MetricDefinitions",
	).WithHeader("X-Auth-Token", "").Expect().Status(http.StatusUnauthorized)
	test.GET(
		"/redfish/v1/TelemetryService/MetricDefinitions",
	).WithHeader("X-Auth-Token", "token").Expect().Status(http.StatusInternalServerError)
}

func TestGetMetricReportDefinitionCollection(t *testing.T) {
	var a TelemetryRPCs
	a.GetMetricReportDefinitionCollectionRPC = testTelemetryService
	testApp := iris.New()
	redfishRoutes := testApp.Party("/redfish/v1/TelemetryService")
	redfishRoutes.Get("/MetricReportDefinitions", a.GetMetricReportDefinitionCollection)
	test := httptest.New(t, testApp)
	test.GET(
		"/redfish/v1/TelemetryService/MetricReportDefinitions",
	).WithHeader("X-Auth-Token", "ValidToken").Expect().Status(http.StatusOK)
	test.GET(
		"/redfish/v1/TelemetryService/MetricReportDefinitions",
	).WithHeader("X-Auth-Token", "").Expect().Status(http.StatusUnauthorized)
	test.GET(
		"/redfish/v1/TelemetryService/MetricReportDefinitions",
	).WithHeader("X-Auth-Token", "token").Expect().Status(http.StatusInternalServerError)
}

func TestGetMetricReportCollection(t *testing.T) {
	var a TelemetryRPCs
	a.GetMetricReportCollectionRPC = testTelemetryService
	testApp := iris.New()
	redfishRoutes := testApp.Party("/redfish/v1/TelemetryService")
	redfishRoutes.Get("/MetricReports", a.GetMetricReportCollection)
	test := httptest.New(t, testApp)
	test.GET(
		"/redfish/v1/TelemetryService/MetricReports",
	).WithHeader("X-Auth-Token", "ValidToken").Expect().Status(http.StatusOK)
	test.GET(
		"/redfish/v1/TelemetryService/MetricReports",
	).WithHeader("X-Auth-Token", "").Expect().Status(http.StatusUnauthorized)
	test.GET(
		"/redfish/v1/TelemetryService/MetricReports",
	).WithHeader("X-Auth-Token", "token").Expect().Status(http.StatusInternalServerError)
}

func TestGetTriggerCollection(t *testing.T) {
	var a TelemetryRPCs
	a.GetTriggerCollectionRPC = testTelemetryService
	testApp := iris.New()
	redfishRoutes := testApp.Party("/redfish/v1/TelemetryService")
	redfishRoutes.Get("/Triggers", a.GetTriggerCollection)
	test := httptest.New(t, testApp)
	test.GET(
		"/redfish/v1/TelemetryService/Triggers",
	).WithHeader("X-Auth-Token", "ValidToken").Expect().Status(http.StatusOK)
	test.GET(
		"/redfish/v1/TelemetryService/Triggers",
	).WithHeader("X-Auth-Token", "").Expect().Status(http.StatusUnauthorized)
	test.GET(
		"/redfish/v1/TelemetryService/Triggers",
	).WithHeader("X-Auth-Token", "token").Expect().Status(http.StatusInternalServerError)
}

func TestGetMetricDefinition(t *testing.T) {
	var a TelemetryRPCs
	a.GetMetricDefinitionRPC = testTelemetryService
	testApp := iris.New()
	redfishRoutes := testApp.Party("/redfish/v1/TelemetryService")
	redfishRoutes.Get("/MetricDefinitions/{id}", a.GetMetricDefinition)
	test := httptest.New(t, testApp)
	test.GET(
		"/redfish/v1/TelemetryService/MetricDefinitions/1",
	).WithHeader("X-Auth-Token", "ValidToken").Expect().Status(http.StatusOK)
	test.GET(
		"/redfish/v1/TelemetryService/MetricDefinitions/1",
	).WithHeader("X-Auth-Token", "").Expect().Status(http.StatusUnauthorized)
	test.GET(
		"/redfish/v1/TelemetryService/MetricDefinitions/1",
	).WithHeader("X-Auth-Token", "token").Expect().Status(http.StatusInternalServerError)
}

func TestGetMetricReportDefinition(t *testing.T) {
	var a TelemetryRPCs
	a.GetMetricReportDefinitionRPC = testTelemetryService
	testApp := iris.New()
	redfishRoutes := testApp.Party("/redfish/v1/TelemetryService")
	redfishRoutes.Get("/MetricReportDefinitions/{id}", a.GetMetricReportDefinition)
	test := httptest.New(t, testApp)
	test.GET(
		"/redfish/v1/TelemetryService/MetricReportDefinitions/1",
	).WithHeader("X-Auth-Token", "ValidToken").Expect().Status(http.StatusOK)
	test.GET(
		"/redfish/v1/TelemetryService/MetricReportDefinitions/1",
	).WithHeader("X-Auth-Token", "").Expect().Status(http.StatusUnauthorized)
	test.GET(
		"/redfish/v1/TelemetryService/MetricReportDefinitions/1",
	).WithHeader("X-Auth-Token", "token").Expect().Status(http.StatusInternalServerError)
}

func TestGetMetricReport(t *testing.T) {
	var a TelemetryRPCs
	a.GetMetricReportRPC = testTelemetryService
	testApp := iris.New()
	redfishRoutes := testApp.Party("/redfish/v1/TelemetryService")
	redfishRoutes.Get("/MetricReports/{id}", a.GetMetricReport)
	test := httptest.New(t, testApp)
	test.GET(
		"/redfish/v1/TelemetryService/MetricReports/1",
	).WithHeader("X-Auth-Token", "ValidToken").Expect().Status(http.StatusOK)
	test.GET(
		"/redfish/v1/TelemetryService/MetricReports/1",
	).WithHeader("X-Auth-Token", "").Expect().Status(http.StatusUnauthorized)
	test.GET(
		"/redfish/v1/TelemetryService/MetricReports/1",
	).WithHeader("X-Auth-Token", "token").Expect().Status(http.StatusInternalServerError)
}

func TestGetTrigger(t *testing.T) {
	var a TelemetryRPCs
	a.GetTriggerRPC = testTelemetryService
	testApp := iris.New()
	redfishRoutes := testApp.Party("/redfish/v1/TelemetryService")
	redfishRoutes.Get("/Triggers/{id}", a.GetTrigger)
	test := httptest.New(t, testApp)
	test.GET(
		"/redfish/v1/TelemetryService/Triggers/1",
	).WithHeader("X-Auth-Token", "ValidToken").Expect().Status(http.StatusOK)
	test.GET(
		"/redfish/v1/TelemetryService/Triggers/1",
	).WithHeader("X-Auth-Token", "").Expect().Status(http.StatusUnauthorized)
	test.GET(
		"/redfish/v1/TelemetryService/Triggers/1",
	).WithHeader("X-Auth-Token", "token").Expect().Status(http.StatusInternalServerError)
}

func TestUpdateTrigger(t *testing.T) {
	var a TelemetryRPCs
	a.UpdateTriggerRPC = testTelemetryService
	testApp := iris.New()
	redfishRoutes := testApp.Party("/redfish/v1/TelemetryService")
	redfishRoutes.Patch("/Triggers/{id}", a.UpdateTrigger)
	test := httptest.New(t, testApp)
	test.PATCH(
		"/redfish/v1/TelemetryService/Triggers/1",
	).WithHeader("X-Auth-Token", "ValidToken").Expect().Status(http.StatusOK)
	test.PATCH(
		"/redfish/v1/TelemetryService/Triggers/1",
	).WithHeader("X-Auth-Token", "").Expect().Status(http.StatusUnauthorized)
	test.PATCH(
		"/redfish/v1/TelemetryService/Triggers/1",
	).WithHeader("X-Auth-Token", "token").Expect().Status(http.StatusInternalServerError)
}
