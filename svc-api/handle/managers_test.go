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
	"fmt"
	"net/http"
	"testing"

	managersproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/managers"
	iris "github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/httptest"
)

func mockGetManagersRequest(req managersproto.ManagerRequest) (*managersproto.ManagerResponse, error) {
	var response = &managersproto.ManagerResponse{}
	if req.ManagerID == "1A" && req.SessionToken == "ValidToken" {
		response = &managersproto.ManagerResponse{
			StatusCode:    200,
			StatusMessage: "Success",
			Body:          []byte(`{"Response":"Success"}`),
		}
	} else if req.ManagerID == "1A" && req.SessionToken == "InvalidToken" {
		response = &managersproto.ManagerResponse{
			StatusCode:    401,
			StatusMessage: "Unauthorized",
			Body:          []byte(`{"Response":"Unauthorized"}`),
		}
	} else if req.ManagerID == "2A" && req.SessionToken == "ValidToken" {
		response = &managersproto.ManagerResponse{
			StatusCode:    403,
			StatusMessage: "Forbidden",
			Body:          []byte(`{"Response":"Forbidden"}`),
		}
	} else if req.ManagerID == "3A" {
		return &managersproto.ManagerResponse{}, fmt.Errorf("RPC Error")
	}
	return response, nil
}

func mockGetManagersResourceRequest(req managersproto.ManagerRequest) (*managersproto.ManagerResponse, error) {
	var response = &managersproto.ManagerResponse{}
	if req.ManagerID == "1A" && req.ResourceID == "1B" && req.SessionToken == "ValidToken" {
		response = &managersproto.ManagerResponse{
			StatusCode:    200,
			StatusMessage: "Success",
			Body:          []byte(`{"Response":"Success"}`),
		}
	} else if req.ManagerID == "1A" && req.ResourceID == "1B" && req.SessionToken == "InvalidToken" {
		response = &managersproto.ManagerResponse{
			StatusCode:    401,
			StatusMessage: "Unauthorized",
			Body:          []byte(`{"Response":"Unauthorized"}`),
		}
	} else if req.ManagerID == "2A" && req.ResourceID == "1B" && req.SessionToken == "ValidToken" {
		response = &managersproto.ManagerResponse{
			StatusCode:    403,
			StatusMessage: "Forbidden",
			Body:          []byte(`{"Response":"Forbidden"}`),
		}
	} else if req.ManagerID == "1A" && req.ResourceID == "2B" && req.SessionToken == "ValidToken" {
		response = &managersproto.ManagerResponse{
			StatusCode:    403,
			StatusMessage: "Forbidden",
			Body:          []byte(`{"Response":"Forbidden"}`),
		}
	} else if req.ManagerID == "1A" && req.ResourceID == "2B" && req.SessionToken == "" {
		response = &managersproto.ManagerResponse{
			StatusCode:    401,
			StatusMessage: "Unauthorized",
			Body:          []byte(`{"Response":"Unauthorized"}`),
		}
	} else if req.ManagerID == "3A" {
		return &managersproto.ManagerResponse{}, fmt.Errorf("RPC Error")
	}
	return response, nil
}
func mockGetManagersCollectionRequest(req managersproto.ManagerRequest) (*managersproto.ManagerResponse, error) {
	var response = &managersproto.ManagerResponse{}
	if req.SessionToken == "ValidToken" {
		response = &managersproto.ManagerResponse{
			StatusCode:    200,
			StatusMessage: "Success",
			Body:          []byte(`{"Response":"Success"}`),
		}
	} else if req.SessionToken == "InvalidToken" {
		response = &managersproto.ManagerResponse{
			StatusCode:    401,
			StatusMessage: "Unauthorized",
			Body:          []byte(`{"Response":"Unauthorized"}`),
		}
	} else if req.SessionToken == "" {
		response = &managersproto.ManagerResponse{
			StatusCode:    401,
			StatusMessage: "Unauthorized",
			Body:          []byte(`{"Response":"Unauthorized"}`),
		}
	} else if req.SessionToken == "token" {
		return &managersproto.ManagerResponse{}, fmt.Errorf("RPC Error")
	}
	return response, nil
}
func mockVirtualMediaInsertRequest(req managersproto.ManagerRequest) (*managersproto.ManagerResponse, error) {
	var response = &managersproto.ManagerResponse{}
	if req.ManagerID == "1A" && req.ResourceID == "1B" && req.SessionToken == "ValidToken" {
		response = &managersproto.ManagerResponse{
			StatusCode:    200,
			StatusMessage: "Success",
			Body:          []byte(`{"Response":"Success"}`),
		}
	} else if req.ManagerID == "1A" && req.ResourceID == "1B" && req.SessionToken == "InvalidToken" {
		response = &managersproto.ManagerResponse{
			StatusCode:    401,
			StatusMessage: "Unauthorized",
			Body:          []byte(`{"Response":"Unauthorized"}`),
		}
	} else if req.ManagerID == "2A" && req.ResourceID == "1B" && req.SessionToken == "ValidToken" {
		response = &managersproto.ManagerResponse{
			StatusCode:    403,
			StatusMessage: "Forbidden",
			Body:          []byte(`{"Response":"Forbidden"}`),
		}
	} else if req.ManagerID == "1A" && req.ResourceID == "2B" && req.SessionToken == "ValidToken" {
		response = &managersproto.ManagerResponse{
			StatusCode:    403,
			StatusMessage: "Forbidden",
			Body:          []byte(`{"Response":"Forbidden"}`),
		}
	} else if req.ManagerID == "1A" && req.ResourceID == "2B" && req.SessionToken == "" {
		response = &managersproto.ManagerResponse{
			StatusCode:    401,
			StatusMessage: "Unauthorized",
			Body:          []byte(`{"Response":"Unauthorized"}`),
		}
	} else if req.ManagerID == "3A" {
		return &managersproto.ManagerResponse{}, fmt.Errorf("RPC Error")
	}
	return response, nil
}

func mockVirtualMediaEjectRequest(req managersproto.ManagerRequest) (*managersproto.ManagerResponse, error) {
	var response = &managersproto.ManagerResponse{}
	if req.ManagerID == "1A" && req.ResourceID == "1B" && req.SessionToken == "ValidToken" {
		response = &managersproto.ManagerResponse{
			StatusCode:    200,
			StatusMessage: "Success",
			Body:          []byte(`{"Response":"Success"}`),
		}
	} else if req.ManagerID == "1A" && req.ResourceID == "1B" && req.SessionToken == "InvalidToken" {
		response = &managersproto.ManagerResponse{
			StatusCode:    401,
			StatusMessage: "Unauthorized",
			Body:          []byte(`{"Response":"Unauthorized"}`),
		}
	} else if req.ManagerID == "2A" && req.ResourceID == "1B" && req.SessionToken == "ValidToken" {
		response = &managersproto.ManagerResponse{
			StatusCode:    403,
			StatusMessage: "Forbidden",
			Body:          []byte(`{"Response":"Forbidden"}`),
		}
	} else if req.ManagerID == "1A" && req.ResourceID == "2B" && req.SessionToken == "ValidToken" {
		response = &managersproto.ManagerResponse{
			StatusCode:    403,
			StatusMessage: "Forbidden",
			Body:          []byte(`{"Response":"Forbidden"}`),
		}
	} else if req.ManagerID == "1A" && req.ResourceID == "2B" && req.SessionToken == "" {
		response = &managersproto.ManagerResponse{
			StatusCode:    401,
			StatusMessage: "Unauthorized",
			Body:          []byte(`{"Response":"Unauthorized"}`),
		}
	} else if req.ManagerID == "3A" {
		return &managersproto.ManagerResponse{}, fmt.Errorf("RPC Error")
	}
	return response, nil
}
func TestGetManager_ValidManagerID(t *testing.T) {
	var mgr ManagersRPCs
	mgr.GetManagersRPC = mockGetManagersRequest
	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1/Managers")
	redfishRoutes.Get("/{id}", mgr.GetManager)
	test := httptest.New(t, mockApp)
	test.GET(
		"/redfish/v1/Managers/1A",
	).WithHeader("X-Auth-Token", "ValidToken").Expect().Status(http.StatusOK)
}

func TestGetManager_InvalidManagerID(t *testing.T) {
	var mgr ManagersRPCs
	mgr.GetManagersRPC = mockGetManagersRequest
	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1/Managers")
	redfishRoutes.Get("/{id}", mgr.GetManager)
	test := httptest.New(t, mockApp)
	test.GET(
		"/redfish/v1/Managers/2A",
	).WithHeader("X-Auth-Token", "ValidToken").Expect().Status(http.StatusForbidden)
}

func TestGetManager_InvalidToken(t *testing.T) {
	var mgr ManagersRPCs
	mgr.GetManagersRPC = mockGetManagersRequest
	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1/Managers")
	redfishRoutes.Get("/{id}", mgr.GetManager)
	test := httptest.New(t, mockApp)
	test.GET(
		"/redfish/v1/Managers/1A",
	).WithHeader("X-Auth-Token", "").Expect().Status(http.StatusUnauthorized)
}

func TestGetManager_RPCError(t *testing.T) {
	var mgr ManagersRPCs
	mgr.GetManagersRPC = mockGetManagersRequest
	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1/Managers")
	redfishRoutes.Get("/{id}", mgr.GetManager)
	test := httptest.New(t, mockApp)
	test.GET(
		"/redfish/v1/Managers/3A",
	).WithHeader("X-Auth-Token", "InvalidToken").Expect().Status(http.StatusInternalServerError)
}

func TestGetManagersCollection_ValidManagerID(t *testing.T) {
	var mgr ManagersRPCs
	mgr.GetManagersCollectionRPC = mockGetManagersCollectionRequest
	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1/Managers")
	redfishRoutes.Get("/", mgr.GetManagersCollection)
	test := httptest.New(t, mockApp)
	test.GET(
		"/redfish/v1/Managers",
	).WithHeader("X-Auth-Token", "ValidToken").Expect().Status(http.StatusOK)
	test.GET(
		"/redfish/v1/Managers",
	).WithHeader("X-Auth-Token", "InvalidToken").Expect().Status(http.StatusUnauthorized)
	test.GET(
		"/redfish/v1/Managers",
	).WithHeader("X-Auth-Token", "").Expect().Status(http.StatusUnauthorized)
	test.GET(
		"/redfish/v1/Managers",
	).WithHeader("X-Auth-Token", "token").Expect().Status(http.StatusInternalServerError)
}

func TestGetManagersResource_ValidManagerID(t *testing.T) {
	var mgr ManagersRPCs
	mgr.GetManagersResourceRPC = mockGetManagersResourceRequest
	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1/Managers")
	redfishRoutes.Get("/{id}/EthernetInterfaces/{rid}", mgr.GetManagersResource)
	test := httptest.New(t, mockApp)
	test.GET(
		"/redfish/v1/Managers/1A/EthernetInterfaces/1B",
	).WithHeader("X-Auth-Token", "ValidToken").Expect().Status(http.StatusOK)
	test.GET(
		"/redfish/v1/Managers/1A/EthernetInterfaces/1B",
	).WithHeader("X-Auth-Token", "InvalidToken").Expect().Status(http.StatusUnauthorized)
	test.GET(
		"/redfish/v1/Managers/2A/EthernetInterfaces/1B",
	).WithHeader("X-Auth-Token", "ValidToken").Expect().Status(http.StatusForbidden)
	test.GET(
		"/redfish/v1/Managers/2A/EthernetInterfaces/1B",
	).WithHeader("X-Auth-Token", "").Expect().Status(http.StatusUnauthorized)
	test.GET(
		"/redfish/v1/Managers/3A/EthernetInterfaces/1B",
	).WithHeader("X-Auth-Token", "InvalidToken").Expect().Status(http.StatusInternalServerError)
}

func TestGetManagerResource_RPCError(t *testing.T) {
	var mgr ManagersRPCs
	mgr.GetManagersRPC = mockGetManagersResourceRequest
	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1/Managers")
	redfishRoutes.Get("/{id}/NetworkInterfaces/{rid}", mgr.GetManager)
	test := httptest.New(t, mockApp)
	test.GET(
		"/redfish/v1/Managers/3A/NetworkInterfaces/1B",
	).WithHeader("X-Auth-Token", "InvalidToken").Expect().Status(http.StatusInternalServerError)
}

func TestVirtualMediaInsert(t *testing.T) {
	var mgr ManagersRPCs
	mgr.VirtualMediaInsertRPC = mockVirtualMediaInsertRequest
	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1/Managers")
	redfishRoutes.Post("/{id}/VirtualMedia/{rid}/VirtualMedia.InsertMedia", mgr.VirtualMediaInsert)
	test := httptest.New(t, mockApp)

	test.POST(
		"/redfish/v1/Managers/1A/VirtualMedia/1B/VirtualMedia.InsertMedia",
	).WithHeader("X-Auth-Token", "ValidToken").WithJSON(map[string]string{"Image": "Body"}).Expect().Status(http.StatusOK)
	test.POST(
		"/redfish/v1/Managers/1A/VirtualMedia/1B/VirtualMedia.InsertMedia",
	).WithHeader("X-Auth-Token", "InvalidToken").WithJSON(map[string]string{"Image": "Body"}).Expect().Status(http.StatusUnauthorized)
	test.POST(
		"/redfish/v1/Managers/2A/VirtualMedia/1B/VirtualMedia.InsertMedia",
	).WithHeader("X-Auth-Token", "ValidToken").WithJSON(map[string]string{"Image": "Body"}).Expect().Status(http.StatusForbidden)
	test.POST(
		"/redfish/v1/Managers/2A/VirtualMedia/1B/VirtualMedia.InsertMedia",
	).WithHeader("X-Auth-Token", "").WithJSON(map[string]string{"Image": "Body"}).Expect().Status(http.StatusUnauthorized)
	test.POST(
		"/redfish/v1/Managers/3A/VirtualMedia/1B/VirtualMedia.InsertMedia",
	).WithHeader("X-Auth-Token", "InvalidToken").WithJSON(map[string]string{"Image": "Body"}).Expect().Status(http.StatusInternalServerError)
}

func TestVirtualMediaEject(t *testing.T) {
	var mgr ManagersRPCs
	mgr.VirtualMediaEjectRPC = mockVirtualMediaInsertRequest
	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1/Managers")
	redfishRoutes.Post("/{id}/VirtualMedia/{rid}/VirtualMedia.EjectMedia", mgr.VirtualMediaEject)
	test := httptest.New(t, mockApp)

	test.POST(
		"/redfish/v1/Managers/1A/VirtualMedia/1B/VirtualMedia.EjectMedia",
	).WithHeader("X-Auth-Token", "ValidToken").Expect().Status(http.StatusOK)
	test.POST(
		"/redfish/v1/Managers/1A/VirtualMedia/1B/VirtualMedia.EjectMedia",
	).WithHeader("X-Auth-Token", "InvalidToken").Expect().Status(http.StatusUnauthorized)
	test.POST(
		"/redfish/v1/Managers/2A/VirtualMedia/1B/VirtualMedia.EjectMedia",
	).WithHeader("X-Auth-Token", "ValidToken").Expect().Status(http.StatusForbidden)
	test.POST(
		"/redfish/v1/Managers/2A/VirtualMedia/1B/VirtualMedia.EjectMedia",
	).WithHeader("X-Auth-Token", "").Expect().Status(http.StatusUnauthorized)
	test.POST(
		"/redfish/v1/Managers/3A/VirtualMedia/1B/VirtualMedia.EjectMedia",
	).WithHeader("X-Auth-Token", "InvalidToken").Expect().Status(http.StatusInternalServerError)
}
