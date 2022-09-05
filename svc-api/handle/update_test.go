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

	updateproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/update"
	iris "github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/httptest"
)

func testGetUpdateService(req updateproto.UpdateRequest) (*updateproto.UpdateResponse, error) {
	var response = &updateproto.UpdateResponse{}
	if req.SessionToken == "ValidToken" {
		response = &updateproto.UpdateResponse{
			StatusCode:    200,
			StatusMessage: "Success",
			Body:          []byte(`{"Response":"Success"}`),
		}
	} else if req.SessionToken == "InvalidToken" {
		response = &updateproto.UpdateResponse{
			StatusCode:    401,
			StatusMessage: "Unauthorized", Body: []byte(`{"Response":"Unauthorized"}`),
		}
	} else if req.SessionToken == "token" {
		return &updateproto.UpdateResponse{}, errors.New("Unable to RPC Call")
	}
	return response, nil
}

func mockGetInventory(updateproto.UpdateRequest) (*updateproto.UpdateResponse, error) {
	return &updateproto.UpdateResponse{
		StatusCode: http.StatusOK,
	}, nil
}

func mockSimpleUpdate(req updateproto.UpdateRequest) (*updateproto.UpdateResponse, error) {
	var response = &updateproto.UpdateResponse{}
	if req.SessionToken == "" {
		response = &updateproto.UpdateResponse{
			StatusCode:    401,
			StatusMessage: "Unauthorized",
			Body:          []byte(`{"Response":"Unauthorized"}`),
		}
	} else if req.SessionToken == "InvalidToken" {
		response = &updateproto.UpdateResponse{
			StatusCode:    401,
			StatusMessage: "Unauthorized",
			Body:          []byte(`{"Response":"Unauthorized"}`),
		}
	} else if req.SessionToken == "TokenRPC" {
		return &updateproto.UpdateResponse{}, errors.New("Unable to RPC Call")
	}
	response = &updateproto.UpdateResponse{
		StatusCode:    200,
		StatusMessage: "Success",
		Body:          []byte(`{"Response":"Success"}`),
	}
	return response, nil
}

func mockStartUpdate(req updateproto.UpdateRequest) (*updateproto.UpdateResponse, error) {
	var response = &updateproto.UpdateResponse{}
	if req.SessionToken == "" {
		response = &updateproto.UpdateResponse{
			StatusCode:    401,
			StatusMessage: "Unauthorized",
			Body:          []byte(`{"Response":"Unauthorized"}`),
		}
	} else if req.SessionToken == "InvalidToken" {
		response = &updateproto.UpdateResponse{
			StatusCode:    401,
			StatusMessage: "Unauthorized",
			Body:          []byte(`{"Response":"Unauthorized"}`),
		}
	} else if req.SessionToken == "TokenRPC" {
		return &updateproto.UpdateResponse{}, errors.New("Unable to RPC Call")
	}
	response = &updateproto.UpdateResponse{
		StatusCode:    200,
		StatusMessage: "Success",
		Body:          []byte(`{"Response":"Success"}`),
	}
	return response, nil
}

func TestGetUpdateService(t *testing.T) {
	header["Allow"] = []string{"GET"}
	defer delete(header, "Allow")
	var a UpdateRPCs
	a.GetUpdateServiceRPC = testGetUpdateService
	testApp := iris.New()
	redfishRoutes := testApp.Party("/redfish/v1/UpdateService")
	redfishRoutes.Get("/", a.GetUpdateService)
	test := httptest.New(t, testApp)
	test.GET(
		"/redfish/v1/UpdateService",
	).WithHeader("X-Auth-Token", "ValidToken").Expect().Status(http.StatusOK).Headers().Equal(header)
	test.GET(
		"/redfish/v1/UpdateService",
	).WithHeader("X-Auth-Token", "").Expect().Status(http.StatusUnauthorized)
	test.GET(
		"/redfish/v1/UpdateService",
	).WithHeader("X-Auth-Token", "token").Expect().Status(http.StatusInternalServerError)
}

func TestGetFirmwareInventoryCollection(t *testing.T) {
	var a UpdateRPCs
	a.GetFirmwareInventoryCollectionRPC = testGetUpdateService
	testApp := iris.New()
	redfishRoutes := testApp.Party("/redfish/v1/UpdateService/FirmwareInventory")
	redfishRoutes.Get("/", a.GetFirmwareInventoryCollection)
	test := httptest.New(t, testApp)
	test.GET(
		"/redfish/v1/UpdateService/FirmwareInventory",
	).WithHeader("X-Auth-Token", "ValidToken").Expect().Status(http.StatusOK)
	test.GET(
		"/redfish/v1/UpdateService/FirmwareInventory",
	).WithHeader("X-Auth-Token", "").Expect().Status(http.StatusUnauthorized)
	test.GET(
		"/redfish/v1/UpdateService/FirmwareInventory",
	).WithHeader("X-Auth-Token", "token").Expect().Status(http.StatusInternalServerError).Headers().Equal(header)
}

func TestGetSoftwareInventoryCollection(t *testing.T) {
	header["Allow"] = []string{"GET"}
	defer delete(header, "Allow")
	var a UpdateRPCs
	a.GetSoftwareInventoryCollectionRPC = testGetUpdateService
	testApp := iris.New()
	redfishRoutes := testApp.Party("/redfish/v1/UpdateService/SoftwareInventory")
	redfishRoutes.Get("/", a.GetSoftwareInventoryCollection)
	test := httptest.New(t, testApp)
	test.GET(
		"/redfish/v1/UpdateService/SoftwareInventory",
	).WithHeader("X-Auth-Token", "ValidToken").Expect().Status(http.StatusOK).Headers().Equal(header)
	test.GET(
		"/redfish/v1/UpdateService/SoftwareInventory",
	).WithHeader("X-Auth-Token", "").Expect().Status(http.StatusUnauthorized)
	test.GET(
		"/redfish/v1/UpdateService/SoftwareInventory",
	).WithHeader("X-Auth-Token", "token").Expect().Status(http.StatusInternalServerError)
}

func TestGetFirmwareInventory(t *testing.T) {
	header["Allow"] = []string{"GET"}
	defer delete(header, "Allow")
	var a UpdateRPCs
	a.GetFirmwareInventoryRPC = mockGetInventory
	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1/UpdateService/FirmwareInventory")
	redfishRoutes.Get("/{id}", a.GetFirmwareInventory)

	e := httptest.New(t, mockApp)
	e.GET(
		"/redfish/v1/UpdateService/FirmwareInventory/6d4a0a66-7efa-578e-83cf-44dc68d2874e.1",
	).WithHeader("X-Auth-Token", "token").Expect().Status(http.StatusOK).Headers().Equal(header)
}

func TestGetSoftwareInventory(t *testing.T) {
	header["Allow"] = []string{"GET"}
	defer delete(header, "Allow")
	var a UpdateRPCs
	a.GetSoftwareInventoryRPC = mockGetInventory
	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1/UpdateService/SoftwareInventory")
	redfishRoutes.Get("/{id}", a.GetSoftwareInventory)

	e := httptest.New(t, mockApp)
	e.GET(
		"/redfish/v1/UpdateService/SoftwareInventory/6d4a0a66-7efa-578e-83cf-44dc68d2874e.1",
	).WithHeader("X-Auth-Token", "token").Expect().Status(http.StatusOK).Headers().Equal(header)
}

func TestSimpleUpdateWithValidToken(t *testing.T) {
	var a UpdateRPCs
	a.SimpleUpdateRPC = mockSimpleUpdate
	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1/UpdateService")
	redfishRoutes.Post("/Actions/UpdateService.SimpleUpdate", a.SimpleUpdate)

	e := httptest.New(t, mockApp)
	e.POST(
		"/redfish/v1/UpdateService/Actions/UpdateService.SimpleUpdate",
	).WithJSON(map[string]interface{}{
		"ImageURI": "/abc/abc",
		"Targets": []interface{}{
			"/redfish/v1/Systems/0356b6f0-5a20-4614-b04a-809c956fe751:1",
		},
		"@Redfish.OperationApplyTime": "OnStartUpdateRequest",
	}).WithHeader("X-Auth-Token", "ValidToken").Expect().Status(http.StatusOK)
}

func TestSimpleUpdateWithoutToken(t *testing.T) {
	var a UpdateRPCs
	a.SimpleUpdateRPC = mockSimpleUpdate
	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1/UpdateService")
	redfishRoutes.Post("/Actions/UpdateService.SimpleUpdate", a.SimpleUpdate)

	e := httptest.New(t, mockApp)
	e.POST(
		"/redfish/v1/UpdateService/Actions/UpdateService.SimpleUpdate",
	).WithJSON(map[string]interface{}{
		"ImageURI": "/abc/abc",
		"Targets": []interface{}{
			"/redfish/v1/Systems/0356b6f0-5a20-4614-b04a-809c956fe751:1",
		},
		"@Redfish.OperationApplyTime": "OnStartUpdateRequest",
	}).WithHeader("X-Auth-Token", "ValidToken").Expect().Status(http.StatusOK)
}

func TestSimpleUpdateWithInvalidToken(t *testing.T) {
	var a UpdateRPCs
	a.SimpleUpdateRPC = mockSimpleUpdate
	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1/UpdateService")
	redfishRoutes.Post("/Actions/UpdateService.SimpleUpdate", a.SimpleUpdate)

	e := httptest.New(t, mockApp)
	e.POST(
		"/redfish/v1/UpdateService/Actions/UpdateService.SimpleUpdate",
	).WithJSON(map[string]interface{}{
		"ImageURI": "/abc/abc",
		"Targets": []interface{}{
			"/redfish/v1/Systems/0356b6f0-5a20-4614-b04a-809c956fe751:1",
		},
		"@Redfish.OperationApplyTime": "OnStartUpdateRequest",
	}).WithHeader("X-Auth-Token", "ValidToken").Expect().Status(http.StatusOK).Headers().Equal(header)
}

func TestSimpleUpdateNegativeTestCases(t *testing.T) {
	//ToDo
}

func TestStartUpdateWithValidToken(t *testing.T) {
	var a UpdateRPCs
	a.StartUpdateRPC = mockStartUpdate
	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1/UpdateService")
	redfishRoutes.Post("/Actions/UpdateService.StartUpdate", a.StartUpdate)

	e := httptest.New(t, mockApp)
	e.POST(
		"/redfish/v1/UpdateService/Actions/UpdateService.StartUpdate",
	).WithHeader("X-Auth-Token", "ValidToken").Expect().Status(http.StatusOK).Headers().Equal(header)
}

func TestStartUpdateWithoutToken(t *testing.T) {
	var a UpdateRPCs
	a.StartUpdateRPC = mockStartUpdate
	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1/UpdateService")
	redfishRoutes.Post("/Actions/UpdateService.StartUpdate", a.StartUpdate)

	e := httptest.New(t, mockApp)
	e.POST(
		"/redfish/v1/UpdateService/Actions/UpdateService.StartUpdate",
	).WithHeader("X-Auth-Token", "ValidToken").Expect().Status(http.StatusOK).Headers().Equal(header)
}

func TestStartUpdateWithInvalidToken(t *testing.T) {
	var a UpdateRPCs
	a.StartUpdateRPC = mockStartUpdate
	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1/UpdateService")
	redfishRoutes.Post("/Actions/UpdateService.StartUpdate", a.StartUpdate)

	e := httptest.New(t, mockApp)
	e.POST(
		"/redfish/v1/UpdateService/Actions/UpdateService.StartUpdate",
	).WithHeader("X-Auth-Token", "ValidToken").Expect().Status(http.StatusOK)
}
