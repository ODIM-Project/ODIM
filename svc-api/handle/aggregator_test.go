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

func testUpdateAggregationSourceRPCCall(req aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error) {
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

func testDeleteAggregationSourceRPCCall(req aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error) {
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

func testAggregateRPCCall(req aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error) {
	var response = &aggregatorproto.AggregatorResponse{}
	if req.SessionToken == "ValidToken" {
		response = &aggregatorproto.AggregatorResponse{
			StatusCode:    http.StatusCreated,
			StatusMessage: "Success",
			Body:          []byte(`{"Response":"Success"}`),
		}
	} else if req.SessionToken == "InvalidToken" {
		response = &aggregatorproto.AggregatorResponse{
			StatusCode:    http.StatusUnauthorized,
			StatusMessage: "Unauthorized", Body: []byte(`{"Response":"Unauthorized"}`),
		}
	} else if req.SessionToken == "" {
		response = &aggregatorproto.AggregatorResponse{
			StatusCode: http.StatusUnauthorized,
		}
	} else if req.SessionToken == "token" {
		return &aggregatorproto.AggregatorResponse{}, errors.New("Unable to RPC Call")
	}
	return response, nil
}

func testGetAggregateRPCCall(req aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error) {
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
	} else if req.SessionToken == "" {
		response = &aggregatorproto.AggregatorResponse{
			StatusCode: http.StatusUnauthorized,
		}
	} else if req.SessionToken == "token" {
		return &aggregatorproto.AggregatorResponse{}, errors.New("Unable to RPC Call")
	}
	return response, nil
}

func testDeleteAggregateRPCCall(req aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error) {
	var response = &aggregatorproto.AggregatorResponse{}
	if req.SessionToken == "ValidToken" {
		response = &aggregatorproto.AggregatorResponse{
			StatusCode: http.StatusNoContent,
		}
	} else if req.SessionToken == "InvalidToken" {
		response = &aggregatorproto.AggregatorResponse{
			StatusCode:    http.StatusUnauthorized,
			StatusMessage: "Unauthorized", Body: []byte(`{"Response":"Unauthorized"}`),
		}
	} else if req.SessionToken == "" {
		response = &aggregatorproto.AggregatorResponse{
			StatusCode: http.StatusUnauthorized,
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
	"Host":     "9.9.9.0",
	"UserName": "admin",
	"Password": "Password1234",
	"Links":    links,
}

var updateAggregationSourceRequest = map[string]interface{}{
	"Host":     "9.9.9.0",
	"UserName": "admin",
	"Password": "Password1234",
}

func TestAddAggregationSource(t *testing.T) {
	var a AggregatorRPCs
	a.AddAggregationSourceRPC = testAddAggregationSourceRPCCall
	testApp := iris.New()
	redfishRoutes := testApp.Party("/redfish/v1/AggregationService/AggregationSources")
	redfishRoutes.Post("/", a.AddAggregationSource)
	test := httptest.New(t, testApp)
	test.POST("/redfish/v1/AggregationService/AggregationSources").WithHeader("X-Auth-Token", "ValidToken").WithJSON(addAggregationSourceRequest).Expect().Status(http.StatusAccepted)
	test.POST("/redfish/v1/AggregationService/AggregationSources").WithHeader("X-Auth-Token", "InvalidToken").WithJSON(addAggregationSourceRequest).Expect().Status(http.StatusUnauthorized)
	test.POST("/redfish/v1/AggregationService/AggregationSources").WithHeader("X-Auth-Token", "token").WithJSON(addAggregationSourceRequest).Expect().Status(http.StatusInternalServerError)
}

func TestGetAllAggregationSource(t *testing.T) {
	var a AggregatorRPCs
	a.GetAllAggregationSourceRPC = testGetAllAggregationSourceRPC
	testApp := iris.New()
	redfishRoutes := testApp.Party("/redfish/v1/AggregationService/AggregationSources")
	redfishRoutes.Get("/", a.GetAllAggregationSource)
	test := httptest.New(t, testApp)
	test.GET(
		"/redfish/v1/AggregationService/AggregationSources",
	).WithHeader("X-Auth-Token", "ValidToken").Expect().Status(http.StatusOK)
	test.GET(
		"/redfish/v1/AggregationService/AggregationSources",
	).WithHeader("X-Auth-Token", "").Expect().Status(http.StatusUnauthorized)
	test.GET(
		"/redfish/v1/AggregationService/AggregationSources",
	).WithHeader("X-Auth-Token", "token").Expect().Status(http.StatusInternalServerError)
}

func TestGetAggregationSource(t *testing.T) {
	var a AggregatorRPCs
	a.GetAggregationSourceRPC = testGetAggregationSourceRPC
	testApp := iris.New()
	redfishRoutes := testApp.Party("/redfish/v1/AggregationService/AggregationSources")
	redfishRoutes.Get("/{id}", a.GetAggregationSource)
	test := httptest.New(t, testApp)
	test.GET(
		"/redfish/v1/AggregationService/AggregationSources/someid",
	).WithHeader("X-Auth-Token", "ValidToken").Expect().Status(http.StatusOK)
	test.GET(
		"/redfish/v1/AggregationService/AggregationSources/someid",
	).WithHeader("X-Auth-Token", "").Expect().Status(http.StatusUnauthorized)
	test.GET(
		"/redfish/v1/AggregationService/AggregationSources/someid",
	).WithHeader("X-Auth-Token", "token").Expect().Status(http.StatusInternalServerError)
}

func TestUpdateAggregationSource(t *testing.T) {
	var a AggregatorRPCs
	a.UpdateAggregationSourceRPC = testUpdateAggregationSourceRPCCall
	testApp := iris.New()
	redfishRoutes := testApp.Party("/redfish/v1/AggregationService/AggregationSources")
	redfishRoutes.Patch("/{id}", a.UpdateAggregationSource)
	test := httptest.New(t, testApp)
	test.PATCH("/redfish/v1/AggregationService/AggregationSources/someid").WithHeader("X-Auth-Token", "ValidToken").WithJSON(updateAggregationSourceRequest).Expect().Status(http.StatusOK)
	test.PATCH("/redfish/v1/AggregationService/AggregationSources/someid").WithHeader("X-Auth-Token", "InvalidToken").WithJSON(updateAggregationSourceRequest).Expect().Status(http.StatusUnauthorized)
	test.PATCH("/redfish/v1/AggregationService/AggregationSources/someid").WithHeader("X-Auth-Token", "token").WithJSON(updateAggregationSourceRequest).Expect().Status(http.StatusInternalServerError)
}

func TestDeleteAggregationSource(t *testing.T) {
	var a AggregatorRPCs
	a.DeleteAggregationSourceRPC = testDeleteAggregationSourceRPCCall
	testApp := iris.New()
	redfishRoutes := testApp.Party("/redfish/v1/AggregationService/AggregationSources")
	redfishRoutes.Delete("/{id}", a.DeleteAggregationSource)
	test := httptest.New(t, testApp)
	test.DELETE("/redfish/v1/AggregationService/AggregationSources/someid").WithHeader("X-Auth-Token", "ValidToken").Expect().Status(http.StatusAccepted)
	test.DELETE("/redfish/v1/AggregationService/AggregationSources/someid").WithHeader("X-Auth-Token", "InvalidToken").Expect().Status(http.StatusUnauthorized)
	test.DELETE("/redfish/v1/AggregationService/AggregationSources/someid").WithHeader("X-Auth-Token", "token").Expect().Status(http.StatusInternalServerError)
}

var aggregateRequest = map[string]interface{}{
	"Elements": []string{
		"/redfish/v1/Systems/423e8254-e3ef-42bd-a130-f096c93a4wq2:1",
		"/redfish/v1/Systems/c14d91b5-3333-48bb-a7b7-75f74a137d48:1",
	},
}

func TestCreateAggregate(t *testing.T) {
	var a AggregatorRPCs
	a.CreateAggregateRPC = testAggregateRPCCall
	testApp := iris.New()
	redfishRoutes := testApp.Party("/redfish/v1/AggregationService/Aggregates")
	redfishRoutes.Post("/", a.CreateAggregate)
	test := httptest.New(t, testApp)
	//  update status code after the code is added
	// test with valid token
	test.POST(
		"/redfish/v1/AggregationService/Aggregates",
	).WithHeader("X-Auth-Token", "ValidToken").WithJSON(aggregateRequest).Expect().Status(http.StatusCreated)

	// test with Invalid token
	test.POST(
		"/redfish/v1/AggregationService/Aggregates",
	).WithHeader("X-Auth-Token", "InvalidToken").WithJSON(aggregateRequest).Expect().Status(http.StatusUnauthorized)

	// test without token
	test.POST(
		"/redfish/v1/AggregationService/Aggregates",
	).WithHeader("X-Auth-Token", "").WithJSON(aggregateRequest).Expect().Status(http.StatusUnauthorized)

	// test without RequestBody
	test.POST(
		"/redfish/v1/AggregationService/Aggregates",
	).WithHeader("X-Auth-Token", "ValidToken").Expect().Status(http.StatusBadRequest)

	// test for RPC Error
	test.POST(
		"/redfish/v1/AggregationService/Aggregates",
	).WithHeader("X-Auth-Token", "token").WithJSON(aggregateRequest).Expect().Status(http.StatusInternalServerError)
}

func TestGetAggregateCollection(t *testing.T) {
	var a AggregatorRPCs
	a.GetAggregateCollectionRPC = testGetAggregateRPCCall
	testApp := iris.New()
	redfishRoutes := testApp.Party("/redfish/v1/AggregationService/Aggregates")
	redfishRoutes.Get("/", a.GetAggregateCollection)
	test := httptest.New(t, testApp)
	// test with valid token
	test.GET(
		"/redfish/v1/AggregationService/Aggregates",
	).WithHeader("X-Auth-Token", "ValidToken").Expect().Status(http.StatusOK)

	// test with Invalid token
	test.GET(
		"/redfish/v1/AggregationService/Aggregates",
	).WithHeader("X-Auth-Token", "InvalidToken").Expect().Status(http.StatusUnauthorized)

	// test without token
	test.GET(
		"/redfish/v1/AggregationService/Aggregates",
	).WithHeader("X-Auth-Token", "").Expect().Status(http.StatusUnauthorized)

	// test for RPC Error
	test.GET(
		"/redfish/v1/AggregationService/Aggregates",
	).WithHeader("X-Auth-Token", "token").Expect().Status(http.StatusInternalServerError)
}

func TestGetAggregate(t *testing.T) {
	var a AggregatorRPCs
	a.GetAggregateRPC = testGetAggregateRPCCall
	testApp := iris.New()
	redfishRoutes := testApp.Party("/redfish/v1/AggregationService/Aggregates/{id}")
	redfishRoutes.Get("/", a.GetAggregate)
	test := httptest.New(t, testApp)
	// test with valid token
	test.GET(
		"/redfish/v1/AggregationService/Aggregates/7ff3bd97-c41c-5de0-937d-85d390691b73",
	).WithHeader("X-Auth-Token", "ValidToken").Expect().Status(http.StatusOK)

	// test with Invalid token
	test.GET(
		"/redfish/v1/AggregationService/Aggregates/7ff3bd97-c41c-5de0-937d-85d390691b73",
	).WithHeader("X-Auth-Token", "InvalidToken").Expect().Status(http.StatusUnauthorized)

	// test without token
	test.GET(
		"/redfish/v1/AggregationService/Aggregates/7ff3bd97-c41c-5de0-937d-85d390691b73",
	).WithHeader("X-Auth-Token", "").Expect().Status(http.StatusUnauthorized)

	// test for RPC Error
	test.GET(
		"/redfish/v1/AggregationService/Aggregates/7ff3bd97-c41c-5de0-937d-85d390691b73",
	).WithHeader("X-Auth-Token", "token").Expect().Status(http.StatusInternalServerError)
}

func TestDeleteAggregate(t *testing.T) {
	var a AggregatorRPCs
	a.DeleteAggregateRPC = testDeleteAggregateRPCCall
	testApp := iris.New()
	redfishRoutes := testApp.Party("/redfish/v1/AggregationService/Aggregates/{id}")
	redfishRoutes.Delete("/", a.DeleteAggregate)
	test := httptest.New(t, testApp)
	// test with valid token
	test.DELETE(
		"/redfish/v1/AggregationService/Aggregates/7ff3bd97-c41c-5de0-937d-85d390691b73",
	).WithHeader("X-Auth-Token", "ValidToken").Expect().Status(http.StatusNoContent)

	// test with Invalid token
	test.DELETE(
		"/redfish/v1/AggregationService/Aggregates/7ff3bd97-c41c-5de0-937d-85d390691b73",
	).WithHeader("X-Auth-Token", "InvalidToken").Expect().Status(http.StatusUnauthorized)

	// test without token
	test.DELETE(
		"/redfish/v1/AggregationService/Aggregates/7ff3bd97-c41c-5de0-937d-85d390691b73",
	).WithHeader("X-Auth-Token", "").Expect().Status(http.StatusUnauthorized)

	// test for RPC Error
	test.DELETE(
		"/redfish/v1/AggregationService/Aggregates/7ff3bd97-c41c-5de0-937d-85d390691b73",
	).WithHeader("X-Auth-Token", "token").Expect().Status(http.StatusInternalServerError)
}

func TestAddElementsToAggregate(t *testing.T) {
	var a AggregatorRPCs
	a.AddElementsToAggregateRPC = testGetAggregateRPCCall
	testApp := iris.New()
	redfishRoutes := testApp.Party("/redfish/v1/AggregationService/Aggregates/{id}/Actions/Aggregate.AddElements")
	redfishRoutes.Post("/", a.AddElementsToAggregate)
	test := httptest.New(t, testApp)
	// test with valid token
	test.POST(
		"/redfish/v1/AggregationService/Aggregates/7ff3bd97-c41c-5de0-937d-85d390691b73/Actions/Aggregate.AddElements",
	).WithHeader("X-Auth-Token", "ValidToken").WithJSON(aggregateRequest).Expect().Status(http.StatusOK)

	// test with Invalid token
	test.POST(
		"/redfish/v1/AggregationService/Aggregates/7ff3bd97-c41c-5de0-937d-85d390691b73/Actions/Aggregate.AddElements",
	).WithHeader("X-Auth-Token", "InvalidToken").WithJSON(aggregateRequest).Expect().Status(http.StatusUnauthorized)

	// test without token
	test.POST(
		"/redfish/v1/AggregationService/Aggregates/7ff3bd97-c41c-5de0-937d-85d390691b73/Actions/Aggregate.AddElements",
	).WithHeader("X-Auth-Token", "").WithJSON(aggregateRequest).Expect().Status(http.StatusUnauthorized)

	// test without request body
	test.POST(
		"/redfish/v1/AggregationService/Aggregates/7ff3bd97-c41c-5de0-937d-85d390691b73/Actions/Aggregate.AddElements",
	).WithHeader("X-Auth-Token", "ValidToken").Expect().Status(http.StatusBadRequest)

	// test for RPC Error
	test.POST(
		"/redfish/v1/AggregationService/Aggregates/7ff3bd97-c41c-5de0-937d-85d390691b73/Actions/Aggregate.AddElements",
	).WithHeader("X-Auth-Token", "token").WithJSON(aggregateRequest).Expect().Status(http.StatusInternalServerError)

}

func TestRemoveElementsFromAggregate(t *testing.T) {
	var a AggregatorRPCs
	a.RemoveElementsFromAggregateRPC = testGetAggregateRPCCall
	testApp := iris.New()
	redfishRoutes := testApp.Party("/redfish/v1/AggregationService/Aggregates/{id}/Actions/Aggregate.RemoveElements")
	redfishRoutes.Post("/", a.RemoveElementsFromAggregate)
	test := httptest.New(t, testApp)
	// test with valid token
	test.POST(
		"/redfish/v1/AggregationService/Aggregates/7ff3bd97-c41c-5de0-937d-85d390691b73/Actions/Aggregate.RemoveElements",
	).WithHeader("X-Auth-Token", "ValidToken").WithJSON(aggregateRequest).Expect().Status(http.StatusOK)

	// test with Invalid token
	test.POST(
		"/redfish/v1/AggregationService/Aggregates/7ff3bd97-c41c-5de0-937d-85d390691b73/Actions/Aggregate.RemoveElements",
	).WithHeader("X-Auth-Token", "InvalidToken").WithJSON(aggregateRequest).Expect().Status(http.StatusUnauthorized)

	// test without token
	test.POST(
		"/redfish/v1/AggregationService/Aggregates/7ff3bd97-c41c-5de0-937d-85d390691b73/Actions/Aggregate.RemoveElements",
	).WithHeader("X-Auth-Token", "").WithJSON(aggregateRequest).Expect().Status(http.StatusUnauthorized)

	// test without request body
	test.POST(
		"/redfish/v1/AggregationService/Aggregates/7ff3bd97-c41c-5de0-937d-85d390691b73/Actions/Aggregate.RemoveElements",
	).WithHeader("X-Auth-Token", "ValidToken").Expect().Status(http.StatusBadRequest)

	// test for RPC Error
	test.POST(
		"/redfish/v1/AggregationService/Aggregates/7ff3bd97-c41c-5de0-937d-85d390691b73/Actions/Aggregate.RemoveElements",
	).WithHeader("X-Auth-Token", "token").WithJSON(aggregateRequest).Expect().Status(http.StatusInternalServerError)
}

func TestResetAggregateElements(t *testing.T) {
	var a AggregatorRPCs
	a.ResetAggregateElementsRPC = testGetAggregateRPCCall
	var aggregateRequest = map[string]interface{}{
		"BatchSize":                    2,
		"DelayBetweenBatchesInSeconds": 2,
		"ResetType":                    "ForceOff",
	}
	testApp := iris.New()
	redfishRoutes := testApp.Party("/redfish/v1/AggregationService/Aggregates/{id}/Actions/Aggregate.Reset")
	redfishRoutes.Post("/", a.ResetAggregateElements)
	test := httptest.New(t, testApp)
	// test with valid token
	test.POST(
		"/redfish/v1/AggregationService/Aggregates/7ff3bd97-c41c-5de0-937d-85d390691b73/Actions/Aggregate.Reset",
	).WithHeader("X-Auth-Token", "ValidToken").WithJSON(aggregateRequest).Expect().Status(http.StatusOK)

	// test with Invalid token
	test.POST(
		"/redfish/v1/AggregationService/Aggregates/7ff3bd97-c41c-5de0-937d-85d390691b73/Actions/Aggregate.Reset",
	).WithHeader("X-Auth-Token", "InvalidToken").WithJSON(aggregateRequest).Expect().Status(http.StatusUnauthorized)

	// test without token
	test.POST(
		"/redfish/v1/AggregationService/Aggregates/7ff3bd97-c41c-5de0-937d-85d390691b73/Actions/Aggregate.Reset",
	).WithHeader("X-Auth-Token", "").WithJSON(aggregateRequest).Expect().Status(http.StatusUnauthorized)

	// test without request body
	test.POST(
		"/redfish/v1/AggregationService/Aggregates/7ff3bd97-c41c-5de0-937d-85d390691b73/Actions/Aggregate.Reset",
	).WithHeader("X-Auth-Token", "ValidToken").Expect().Status(http.StatusBadRequest)

	// test for RPC Error
	test.POST(
		"/redfish/v1/AggregationService/Aggregates/7ff3bd97-c41c-5de0-937d-85d390691b73/Actions/Aggregate.Reset",
	).WithHeader("X-Auth-Token", "token").WithJSON(aggregateRequest).Expect().Status(http.StatusInternalServerError)
}

func TestSetDefaultBootOrderAggregateElements(t *testing.T) {
	var a AggregatorRPCs
	a.SetDefaultBootOrderAggregateElementsRPC = testGetAggregateRPCCall
	var aggregateRequest = map[string]interface{}{
		"BatchSize":                    2,
		"DelayBetweenBatchesInSeconds": 2,
		"ResetType":                    "ForceOff",
	}
	testApp := iris.New()
	redfishRoutes := testApp.Party("/redfish/v1/AggregationService/Aggregates/{id}/Actions/Aggregate.SetDefaultBootOrder")
	redfishRoutes.Post("/", a.SetDefaultBootOrderAggregateElements)
	test := httptest.New(t, testApp)
	// test with valid token
	test.POST(
		"/redfish/v1/AggregationService/Aggregates/7ff3bd97-c41c-5de0-937d-85d390691b73/Actions/Aggregate.SetDefaultBootOrder",
	).WithHeader("X-Auth-Token", "ValidToken").WithJSON(aggregateRequest).Expect().Status(http.StatusOK)

	// test with Invalid token
	test.POST(
		"/redfish/v1/AggregationService/Aggregates/7ff3bd97-c41c-5de0-937d-85d390691b73/Actions/Aggregate.SetDefaultBootOrder",
	).WithHeader("X-Auth-Token", "InvalidToken").WithJSON(aggregateRequest).Expect().Status(http.StatusUnauthorized)

	// test without token
	test.POST(
		"/redfish/v1/AggregationService/Aggregates/7ff3bd97-c41c-5de0-937d-85d390691b73/Actions/Aggregate.SetDefaultBootOrder",
	).WithHeader("X-Auth-Token", "").WithJSON(aggregateRequest).Expect().Status(http.StatusUnauthorized)

	// test for RPC Error
	test.POST(
		"/redfish/v1/AggregationService/Aggregates/7ff3bd97-c41c-5de0-937d-85d390691b73/Actions/Aggregate.SetDefaultBootOrder",
	).WithHeader("X-Auth-Token", "token").WithJSON(aggregateRequest).Expect().Status(http.StatusInternalServerError)
}

func TestGetAllConnectionMethods(t *testing.T) {
	var a AggregatorRPCs
	a.GetAllConnectionMethodsRPC = testGetAggregateRPCCall

	testApp := iris.New()
	redfishRoutes := testApp.Party("/redfish/v1/AggregationService/ConnectionMethods")
	redfishRoutes.Get("/", a.GetAllConnectionMethods)
	test := httptest.New(t, testApp)
	// test with valid token
	test.GET(
		"/redfish/v1/AggregationService/ConnectionMethods",
	).WithHeader("X-Auth-Token", "ValidToken").Expect().Status(http.StatusOK)

	// test with Invalid token
	test.GET(
		"/redfish/v1/AggregationService/ConnectionMethods",
	).WithHeader("X-Auth-Token", "InvalidToken").Expect().Status(http.StatusUnauthorized)

	// test without token
	test.GET(
		"/redfish/v1/AggregationService/ConnectionMethods",
	).WithHeader("X-Auth-Token", "").Expect().Status(http.StatusUnauthorized)

	// test for RPC Error
	test.GET(
		"/redfish/v1/AggregationService/ConnectionMethods",
	).WithHeader("X-Auth-Token", "token").Expect().Status(http.StatusInternalServerError)
}

func TestGetConnectionMethod(t *testing.T) {
	var a AggregatorRPCs
	a.GetConnectionMethodRPC = testGetAggregateRPCCall

	testApp := iris.New()
	redfishRoutes := testApp.Party("/redfish/v1/AggregationService/ConnectionMethods")
	redfishRoutes.Get("/{id}", a.GetConnectionMethod)
	test := httptest.New(t, testApp)
	// test with valid token
	test.GET(
		"/redfish/v1/AggregationService/ConnectionMethods/74116e00-0a4a-53e6-a959-e6a7465d6358",
	).WithHeader("X-Auth-Token", "ValidToken").Expect().Status(http.StatusOK) //TODO : replace with http.StatusOK

	// test with Invalid token
	test.GET(
		"/redfish/v1/AggregationService/ConnectionMethods/74116e00-0a4a-53e6-a959-e6a7465d6358",
	).WithHeader("X-Auth-Token", "InvalidToken").Expect().Status(http.StatusUnauthorized) //TODO : replace with http.StatusUnauthorized

	// test without token
	test.GET(
		"/redfish/v1/AggregationService/ConnectionMethods/74116e00-0a4a-53e6-a959-e6a7465d6358",
	).WithHeader("X-Auth-Token", "").Expect().Status(http.StatusUnauthorized) //TODO : replace with http.StatusUnauthorized

	// test for RPC Error
	test.GET(
		"/redfish/v1/AggregationService/ConnectionMethods/74116e00-0a4a-53e6-a959-e6a7465d6358",
	).WithHeader("X-Auth-Token", "token").Expect().Status(http.StatusInternalServerError) //TODO : replace with http.StatusInternalServerError
}
