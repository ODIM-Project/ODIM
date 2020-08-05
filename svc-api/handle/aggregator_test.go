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

	aggregatorproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/aggregator"
	iris "github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/httptest"
)

func testDeleteComputeRPC(req aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error) {
	return &aggregatorproto.AggregatorResponse{
		StatusCode: http.StatusOK,
	}, nil
}
func testAddComputeRPC(req aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error) {
	return &aggregatorproto.AggregatorResponse{
		StatusCode: http.StatusOK,
	}, nil
}
func testAddComputeRPCWithRPCError(req aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error) {
	return &aggregatorproto.AggregatorResponse{}, errors.New("Unable to RPC Call")
}
func testDeleteComputeRPCWIthRPCError(req aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error) {
	return &aggregatorproto.AggregatorResponse{}, errors.New("Unable to RPC Call")
}
func testGetAggregationService(req aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error) {
	var response = &aggregatorproto.AggregatorResponse{}
	if req.SessionToken == "ValidToken" {
		response = &aggregatorproto.AggregatorResponse{
			StatusCode:    200,
			StatusMessage: "Success",
			Body:          []byte(`{"Response":"Success"}`),
		}
	} else if req.SessionToken == "InvalidToken" {
		response = &aggregatorproto.AggregatorResponse{
			StatusCode:    401,
			StatusMessage: "Unauthorized", Body: []byte(`{"Response":"Unauthorized"}`),
		}
	} else if req.SessionToken == "token" {
		return &aggregatorproto.AggregatorResponse{}, errors.New("Unable to RPC Call")
	}
	return response, nil
}

func testAddAggregationSourceRPCCall(req aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error) {
	var response = &aggregatorproto.AggregatorResponse{}
	if req.SessionToken == "ValidToken" {
		response = &aggregatorproto.AggregatorResponse{
			StatusCode:    http.StatusAccepted,
			StatusMessage: "Success",
			Body:          []byte(`{"Response":"Success"}`),
		}
	} else if req.SessionToken == "InvalidToken" {
		response = &aggregatorproto.AggregatorResponse{
			StatusCode:    http.StatusUnauthorized,
			StatusMessage: "Unauthorized", Body: []byte(`{"Response":"Unauthorized"}`),
		}
	} else if req.SessionToken == "token" {
		return &aggregatorproto.AggregatorResponse{}, errors.New("Unable to RPC Call")
	}
	return response, nil
}

func testGetAllAggregationSourceRPC(req aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error) {
	var response = &aggregatorproto.AggregatorResponse{}
	if req.SessionToken == "ValidToken" {
		response = &aggregatorproto.AggregatorResponse{
			StatusCode:    http.StatusOK,
			StatusMessage: "Success",
			Body:          []byte(`{"Response":"Success"}`),
		}
	} else if req.SessionToken == "InvalidToken" {
		response = &aggregatorproto.AggregatorResponse{
			StatusCode:    http.StatusUnauthorized,
			StatusMessage: "Unauthorized", Body: []byte(`{"Response":"Unauthorized"}`),
		}
	} else if req.SessionToken == "token" {
		return &aggregatorproto.AggregatorResponse{}, errors.New("Unable to RPC Call")
	}
	return response, nil
}

func testGetAggregationSourceRPC(req aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error) {
	var response = &aggregatorproto.AggregatorResponse{}
	if req.SessionToken == "ValidToken" {
		response = &aggregatorproto.AggregatorResponse{
			StatusCode:    http.StatusOK,
			StatusMessage: "Success",
			Body:          []byte(`{"Response":"Success"}`),
		}
	} else if req.SessionToken == "InvalidToken" {
		response = &aggregatorproto.AggregatorResponse{
			StatusCode:    http.StatusUnauthorized,
			StatusMessage: "Unauthorized", Body: []byte(`{"Response":"Unauthorized"}`),
		}
	} else if req.SessionToken == "token" {
		return &aggregatorproto.AggregatorResponse{}, errors.New("Unable to RPC Call")
	}
	return response, nil
}

type params struct {
	Name string
}

var request = map[string]interface{}{
	"@odata.context": "/redfish/v1/$metadata#ActionInfo.ActionInfo",
	"@odata.id":      "/redfish/v1/AggregationService/RemoveActionInfo",
	"@odata.type":    "#ActionInfo.v1_0_3.ActionInfo",
	"Id":             "RemoveActionInfo",
	"Name":           "Remove Action Info",
	"Oem":            "",
	"Parameters":     []params{{Name: "uri"}},
}

func TestDeleteCompute(t *testing.T) {
	var a AggregatorRPCs
	a.DeleteComputeRPC = testDeleteComputeRPC

	testApp := iris.New()
	redfishRoutes := testApp.Party("/redfish/v1")
	redfishRoutes.Delete("/AggregationService#RemoveActionInfo", a.DeleteCompute)

	e := httptest.New(t, testApp)
	e.DELETE(
		"/redfish/v1/AggregationService#RemoveActionInfo",
	).WithHeader("X-Auth-Token", "token").WithJSON(request).Expect().Status(http.StatusOK)
}
func TestAddCompute(t *testing.T) {
	var a AggregatorRPCs
	a.AddComputeRPC = testAddComputeRPC

	testApp := iris.New()
	redfishRoutes := testApp.Party("/redfish/v1")
	redfishRoutes.Post("/AggregationService#Add", a.AddCompute)

	e := httptest.New(t, testApp)
	e.POST(
		"/redfish/v1/AggregationService#Add",
	).WithHeader("X-Auth-Token", "token").WithJSON(request).Expect().Status(http.StatusOK)
	e.POST(
		"/redfish/v1/AggregationService#Add",
	).WithHeader("X-Auth-Token", "token").Expect().Status(http.StatusBadRequest)
	e.POST(
		"/redfish/v1/AggregationService#Add",
	).WithHeader("X-Auth-Token", "").WithJSON(request).Expect().Status(http.StatusUnauthorized)
}
func TestAddComputeWithRPCError(t *testing.T) {
	var a AggregatorRPCs
	a.AddComputeRPC = testAddComputeRPCWithRPCError

	testApp := iris.New()
	redfishRoutes := testApp.Party("/redfish/v1")
	redfishRoutes.Post("/AggregationService#Add", a.AddCompute)

	e := httptest.New(t, testApp)
	e.POST(
		"/redfish/v1/AggregationService#Add",
	).WithHeader("X-Auth-Token", "token").WithJSON(request).Expect().Status(http.StatusInternalServerError)
}
func TestResetCompute(t *testing.T) {
	var a AggregatorRPCs
	a.ResetRPC = testAddComputeRPC

	testApp := iris.New()
	redfishRoutes := testApp.Party("/redfish/v1")
	redfishRoutes.Post("/AggregationService#Reset", a.Reset)

	e := httptest.New(t, testApp)
	e.POST(
		"/redfish/v1/AggregationService#Reset",
	).WithHeader("X-Auth-Token", "token").WithJSON(request).Expect().Status(http.StatusOK)
	e.POST(
		"/redfish/v1/AggregationService#Reset",
	).WithHeader("X-Auth-Token", "token").Expect().Status(http.StatusBadRequest)
	e.POST(
		"/redfish/v1/AggregationService#Reset",
	).WithHeader("X-Auth-Token", "").WithJSON(request).Expect().Status(http.StatusUnauthorized)
}

func TestResetComputeWithRPCError(t *testing.T) {
	var a AggregatorRPCs
	a.ResetRPC = testAddComputeRPCWithRPCError

	testApp := iris.New()
	redfishRoutes := testApp.Party("/redfish/v1")
	redfishRoutes.Post("/AggregationService#Reset", a.Reset)

	e := httptest.New(t, testApp)
	e.POST(
		"/redfish/v1/AggregationService#Reset",
	).WithHeader("X-Auth-Token", "token").WithJSON(request).Expect().Status(http.StatusInternalServerError)
}
func TestSetDefaultBootOrderCompute(t *testing.T) {
	var a AggregatorRPCs
	a.SetDefaultBootOrderRPC = testAddComputeRPC

	testApp := iris.New()
	redfishRoutes := testApp.Party("/redfish/v1")
	redfishRoutes.Post("/AggregationService#SetDefaultBootOrder", a.SetDefaultBootOrder)

	e := httptest.New(t, testApp)
	e.POST(
		"/redfish/v1/AggregationService#SetDefaultBootOrder",
	).WithHeader("X-Auth-Token", "token").WithJSON(request).Expect().Status(http.StatusOK)
	e.POST(
		"/redfish/v1/AggregationService#SetDefaultBootOrder",
	).WithHeader("X-Auth-Token", "token").Expect().Status(http.StatusBadRequest)
	e.POST(
		"/redfish/v1/AggregationService#SetDefaultBootOrder",
	).WithHeader("X-Auth-Token", "").WithJSON(request).Expect().Status(http.StatusUnauthorized)
}

func TestSetDefaultBootOrderComputeWithRPCError(t *testing.T) {
	var a AggregatorRPCs
	a.SetDefaultBootOrderRPC = testAddComputeRPCWithRPCError

	testApp := iris.New()
	redfishRoutes := testApp.Party("/redfish/v1")
	redfishRoutes.Post("/AggregationService#SetDefaultBootOrder", a.SetDefaultBootOrder)

	e := httptest.New(t, testApp)
	e.POST(
		"/redfish/v1/AggregationService#SetDefaultBootOrder",
	).WithHeader("X-Auth-Token", "token").WithJSON(request).Expect().Status(http.StatusInternalServerError)
}
func TestDeleteComputeWithoutToken(t *testing.T) {
	var a AggregatorRPCs
	a.DeleteComputeRPC = testDeleteComputeRPC

	testApp := iris.New()
	redfishRoutes := testApp.Party("/redfish/v1")
	redfishRoutes.Delete("/AggregationService#RemoveActionInfo", a.DeleteCompute)

	e := httptest.New(t, testApp)
	e.DELETE(
		"/redfish/v1/AggregationService#RemoveActionInfo",
	).WithJSON(request).Expect().Status(http.StatusUnauthorized)
}

func TestDeleteComputeWithoutRequestbody(t *testing.T) {
	var a AggregatorRPCs
	a.DeleteComputeRPC = testDeleteComputeRPC

	testApp := iris.New()
	redfishRoutes := testApp.Party("/redfish/v1")
	redfishRoutes.Delete("/AggregationService#RemoveActionInfo", a.DeleteCompute)

	e := httptest.New(t, testApp)
	e.DELETE(
		"/redfish/v1/AggregationService#RemoveActionInfo",
	).WithHeader("X-Auth-Token", "token").Expect().Status(http.StatusBadRequest)
}

func TestDeleteComputeWithoutRPCCall(t *testing.T) {
	var a AggregatorRPCs
	a.DeleteComputeRPC = testDeleteComputeRPCWIthRPCError

	testApp := iris.New()
	redfishRoutes := testApp.Party("/redfish/v1")
	redfishRoutes.Delete("/AggregationService#RemoveActionInfo", a.DeleteCompute)

	e := httptest.New(t, testApp)
	e.DELETE(
		"/redfish/v1/AggregationService#RemoveActionInfo",
	).WithHeader("X-Auth-Token", "token").WithJSON(request).Expect().Status(http.StatusInternalServerError)
}

func TestGetAggregationService(t *testing.T) {
	var a AggregatorRPCs
	a.GetAggregationServiceRPC = testGetAggregationService
	testApp := iris.New()
	redfishRoutes := testApp.Party("/redfish/v1/AggregationService")
	redfishRoutes.Get("/", a.GetAggregationService)
	test := httptest.New(t, testApp)
	test.GET(
		"/redfish/v1/AggregationService",
	).WithHeader("X-Auth-Token", "ValidToken").Expect().Status(http.StatusOK)
	test.GET(
		"/redfish/v1/AggregationService",
	).WithHeader("X-Auth-Token", "").Expect().Status(http.StatusUnauthorized)
	test.GET(
		"/redfish/v1/AggregationService",
	).WithHeader("X-Auth-Token", "token").Expect().Status(http.StatusInternalServerError)
}

var oem = map[string]interface{}{
	"PluginID": "ILO",
}
var links = map[string]interface{}{
	"Oem": oem,
}
var addAggregationSourceRequest = map[string]interface{}{
	"Host":     "10.24.0.14",
	"UserName": "admin",
	"Password": "Password1234",
	"Links":    links,
}

func TestAddAggregationSource(t *testing.T) {
	var a AggregatorRPCs
	a.AddAggregationSourceRPC = testAddAggregationSourceRPCCall
	testApp := iris.New()
	redfishRoutes := testApp.Party("/redfish/v1/AggregationService/AggregationSource")
	redfishRoutes.Post("/", a.AddAggregationSource)
	test := httptest.New(t, testApp)
	//  update status code after the code is added
	test.POST("/redfish/v1/AggregationService/AggregationSource").WithHeader("X-Auth-Token", "ValidToken").WithJSON(addAggregationSourceRequest).Expect().Status(http.StatusAccepted)
	test.POST("/redfish/v1/AggregationService/AggregationSource").WithHeader("X-Auth-Token", "InvalidToken").WithJSON(addAggregationSourceRequest).Expect().Status(http.StatusUnauthorized)
	test.POST("/redfish/v1/AggregationService/AggregationSource").WithHeader("X-Auth-Token", "token").WithJSON(addAggregationSourceRequest).Expect().Status(http.StatusInternalServerError)
}

func TestGetAllAggregationSource(t *testing.T) {
	var a AggregatorRPCs
	a.GetAllAggregationSourceRPC = testGetAllAggregationSourceRPC
	testApp := iris.New()
	redfishRoutes := testApp.Party("/redfish/v1/AggregationService/AggregationSource")
	redfishRoutes.Get("/", a.GetAllAggregationSource)
	test := httptest.New(t, testApp)
	// change the  status code after the code is added
	test.GET(
		"/redfish/v1/AggregationService/AggregationSource",
	).WithHeader("X-Auth-Token", "ValidToken").Expect().Status(http.StatusOK)
	test.GET(
		"/redfish/v1/AggregationService/AggregationSource",
	).WithHeader("X-Auth-Token", "").Expect().Status(http.StatusUnauthorized)
	test.GET(
		"/redfish/v1/AggregationService/AggregationSource",
	).WithHeader("X-Auth-Token", "token").Expect().Status(http.StatusInternalServerError)
}

func TestGetAggregationSource(t *testing.T) {
	var a AggregatorRPCs
	a.GetAggregationSourceRPC = testGetAggregationSourceRPC
	testApp := iris.New()
	redfishRoutes := testApp.Party("/redfish/v1/AggregationService/AggregationSource")
	redfishRoutes.Get("/{id}", a.GetAggregationSource)
	test := httptest.New(t, testApp)
	//  change the  status code after the code is added
	test.GET(
		"/redfish/v1/AggregationService/AggregationSource/someid",
	).WithHeader("X-Auth-Token", "ValidToken").Expect().Status(http.StatusOK)
	test.GET(
		"/redfish/v1/AggregationService/AggregationSource/someid",
	).WithHeader("X-Auth-Token", "").Expect().Status(http.StatusUnauthorized)
	test.GET(
		"/redfish/v1/AggregationService/AggregationSource/someid",
	).WithHeader("X-Auth-Token", "token").Expect().Status(http.StatusInternalServerError)
}
