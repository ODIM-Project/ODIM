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

	systemsproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/systems"
	iris "github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/httptest"
)

func mockGetSystemRequest(req systemsproto.GetSystemsRequest) (*systemsproto.SystemsResponse, error) {
	var response = &systemsproto.SystemsResponse{}
	if req.URL == "/redfish/v1/Systems/1A" && req.SessionToken == "ValidToken" {
		response = &systemsproto.SystemsResponse{
			StatusCode:    200,
			StatusMessage: "Success",
			Body:          []byte(`{"Response":"Success"}`),
		}
	} else if req.URL == "/redfish/v1/Systems%2f1A" && req.SessionToken == "ValidToken" {
		response = &systemsproto.SystemsResponse{
			StatusCode:    200,
			StatusMessage: "Success",
			Body:          []byte(`{"Response":"Success"}`),
		}
	} else if req.URL == "/redfish/v1/Systems/1A" && req.SessionToken == "InvalidToken" {
		response = &systemsproto.SystemsResponse{
			StatusCode:    401,
			StatusMessage: "Unauthorized",
			Body:          []byte(`{"Response":"Unauthorized"}`),
		}
	} else if req.URL == "/redfish/v1/Systems/1A" && req.SessionToken == "" {
		response = &systemsproto.SystemsResponse{
			StatusCode:    401,
			StatusMessage: "Unauthorized",
			Body:          []byte(`{"Response":"Unauthorized"}`),
		}
	} else if req.URL == "/redfish/v1/Systems/2A" && req.SessionToken == "ValidToken" {
		response = &systemsproto.SystemsResponse{
			StatusCode:    403,
			StatusMessage: "Forbidden",
			Body:          []byte(`{"Response":"Forbidden"}`),
		}
	} else if req.SessionToken == "TokenRPC" {
		return &systemsproto.SystemsResponse{}, errors.New("Unable to RPC Call")
	}
	return response, nil
}

func mockGetSystemsCollection(req systemsproto.GetSystemsRequest) (*systemsproto.SystemsResponse, error) {
	var response = &systemsproto.SystemsResponse{}
	if req.URL == "/redfish/v1/Systems" && req.SessionToken == "ValidToken" {
		response = &systemsproto.SystemsResponse{
			StatusCode:    200,
			StatusMessage: "Success",
			Body:          []byte(`{"Response":"Success"}`),
		}
	} else if req.URL == "/redfish/v1/Systems" && req.SessionToken == "InvalidToken" {
		response = &systemsproto.SystemsResponse{
			StatusCode:    401,
			StatusMessage: "Unauthorized",
			Body:          []byte(`{"Response":"Unauthorized"}`),
		}
	} else if req.URL == "/redfish/v1/Systems" && req.SessionToken == "" {
		response = &systemsproto.SystemsResponse{
			StatusCode:    401,
			StatusMessage: "Unauthorized",
			Body:          []byte(`{"Response":"Unauthorized"}`),
		}
	} else if req.SessionToken == "TokenRPC" {
		return &systemsproto.SystemsResponse{}, errors.New("Unable to RPC Call")
	}
	return response, nil
}

func TestGetSystemsCollection_ValidToken(t *testing.T) {
	var sys SystemRPCs
	sys.GetSystemsCollectionRPC = mockGetSystemsCollection
	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1")
	redfishRoutes.Get("/Systems", sys.GetSystemsCollection)
	test := httptest.New(t, mockApp)
	test.GET(
		"/redfish/v1/Systems",
	).WithHeader("X-Auth-Token", "ValidToken").Expect().Status(http.StatusOK)
	test.GET(
		"/redfish/v1/Systems",
	).WithHeader("X-Auth-Token", "TokenRPC").Expect().Status(http.StatusInternalServerError)
}

func TestGetSystemsCollection_NoValidToken(t *testing.T) {
	var sys SystemRPCs
	sys.GetSystemsCollectionRPC = mockGetSystemsCollection
	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1")
	redfishRoutes.Get("/Systems", sys.GetSystemsCollection)
	test := httptest.New(t, mockApp)
	test.GET(
		"/redfish/v1/Systems",
	).WithHeader("X-Auth-Token", "").Expect().Status(http.StatusUnauthorized)
}

func TestGetSystem_ValidUUid(t *testing.T) {
	var sys SystemRPCs
	sys.GetSystemRPC = mockGetSystemRequest
	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1/Systems")
	redfishRoutes.Get("/{id}", sys.GetSystem)
	test := httptest.New(t, mockApp)
	test.GET(
		"/redfish/v1/Systems/1A",
	).WithHeader("X-Auth-Token", "ValidToken").Expect().Status(http.StatusOK)
	test.GET(
		"/redfish/v1/Systems/1A",
	).WithHeader("X-Auth-Token", "TokenRPC").Expect().Status(http.StatusInternalServerError)
}

func TestGetSystem_InvalidUUid(t *testing.T) {
	var sys SystemRPCs
	sys.GetSystemRPC = mockGetSystemRequest
	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1/Systems")
	redfishRoutes.Get("/{id}", sys.GetSystem)
	test := httptest.New(t, mockApp)
	test.GET(
		"/redfish/v1/Systems/2A",
	).WithHeader("X-Auth-Token", "ValidToken").Expect().Status(http.StatusForbidden)
}

func TestGetSystem_NoToken(t *testing.T) {
	var sys SystemRPCs
	sys.GetSystemRPC = mockGetSystemRequest
	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1/Systems")
	redfishRoutes.Get("/{id}", sys.GetSystem)
	test := httptest.New(t, mockApp)
	test.GET(
		"/redfish/v1/Systems/1A",
	).WithHeader("X-Auth-Token", "").Expect().Status(http.StatusUnauthorized)
}

func TestGetSystem_InvalidToken(t *testing.T) {
	var sys SystemRPCs
	sys.GetSystemRPC = mockGetSystemRequest
	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1/Systems")
	redfishRoutes.Get("/{id}", sys.GetSystem)
	test := httptest.New(t, mockApp)
	test.GET(
		"/redfish/v1/Systems/1A",
	).WithHeader("X-Auth-Token", "InvalidToken").Expect().Status(http.StatusUnauthorized)
}

func mockGetSystemResource(systemsproto.GetSystemsRequest) (*systemsproto.SystemsResponse, error) {
	return &systemsproto.SystemsResponse{
		StatusCode: http.StatusOK,
	}, nil
}

func mockGetSystemResourceRPCError(systemsproto.GetSystemsRequest) (*systemsproto.SystemsResponse, error) {
	return &systemsproto.SystemsResponse{}, errors.New("RPC Error")
}

func TestSystemRPCs_GetSystemResource(t *testing.T) {
	var sys SystemRPCs
	sys.GetSystemResourceRPC = mockGetSystemResource
	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1/Systems")
	redfishRoutes.Get("/{id}/SecureBoot", sys.GetSystemResource)

	e := httptest.New(t, mockApp)
	e.GET(
		"/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1/SecureBoot",
	).WithHeader("X-Auth-Token", "token").Expect().Status(http.StatusOK)
}

func TestSystemRPCs_GetSystemResourceRPCError(t *testing.T) {
	var sys SystemRPCs
	sys.GetSystemResourceRPC = mockGetSystemResourceRPCError
	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1/Systems")
	redfishRoutes.Get("/{id}/SecureBoot", sys.GetSystemResource)

	e := httptest.New(t, mockApp)
	e.GET(
		"/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1/SecureBoot",
	).WithHeader("X-Auth-Token", "token").Expect().Status(http.StatusInternalServerError)
}

func TestSystemRPCs_GetSystemResourceWithoutToken(t *testing.T) {
	var sys SystemRPCs
	sys.GetSystemResourceRPC = mockGetSystemResource
	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1/Systems")
	redfishRoutes.Get("/{id}/SecureBoot", sys.GetSystemResource)

	e := httptest.New(t, mockApp)
	e.GET(
		"/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1/SecureBoot",
	).Expect().Status(http.StatusUnauthorized)
}

func mockChangeBiosSettings(req systemsproto.BiosSettingsRequest) (*systemsproto.SystemsResponse, error) {
	var response = &systemsproto.SystemsResponse{}
	if req.SessionToken == "" {
		response = &systemsproto.SystemsResponse{
			StatusCode:    401,
			StatusMessage: "Unauthorized",
			Body:          []byte(`{"Response":"Unauthorized"}`),
		}
	} else if req.SessionToken == "InvalidToken" {
		response = &systemsproto.SystemsResponse{
			StatusCode:    401,
			StatusMessage: "Unauthorized",
			Body:          []byte(`{"Response":"Unauthorized"}`),
		}
	} else if req.SessionToken == "ValidToken" && req.SystemID == "" {
		response = &systemsproto.SystemsResponse{
			StatusCode:    400,
			StatusMessage: "BadRequest",
			Body:          []byte(`{"Response":"BadRequest"}`),
		}
	} else if req.SessionToken == "TokenRPC" {
		return &systemsproto.SystemsResponse{}, errors.New("Unable to RPC Call")
	} else {
		response = &systemsproto.SystemsResponse{
			StatusCode:    200,
			StatusMessage: "Success",
			Body:          []byte(`{"Response":"Success"}`),
		}
	}
	return response, nil
}

func TestChangeBiosSettingsWithValidToken(t *testing.T) {
	var sys SystemRPCs
	sys.ChangeBiosSettingsRPC = mockChangeBiosSettings
	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1/Systems")
	redfishRoutes.Patch("/{id}/Bios/Settings", sys.ChangeBiosSettings)

	e := httptest.New(t, mockApp)
	e.PATCH(
		"/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1/Bios/Settings",
	).WithJSON(map[string]string{"Sample": "Body"}).WithHeader("X-Auth-Token", "ValidToken").Expect().Status(http.StatusOK)
}

func TestChangeBiosSettingsWithoutToken(t *testing.T) {
	var sys SystemRPCs
	sys.ChangeBiosSettingsRPC = mockChangeBiosSettings
	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1/Systems")
	redfishRoutes.Patch("/{id}/Bios/Settings", sys.ChangeBiosSettings)

	e := httptest.New(t, mockApp)
	e.PATCH(
		"/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1/Bios/Settings",
	).WithJSON(map[string]string{"Sample": "Body"}).WithHeader("X-Auth-Token", "").Expect().Status(http.StatusUnauthorized)
}

func TestChangeBiosSettingsWithInvalidToken(t *testing.T) {
	var sys SystemRPCs
	sys.ChangeBiosSettingsRPC = mockChangeBiosSettings
	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1/Systems")
	redfishRoutes.Patch("/{id}/Bios/Settings", sys.ChangeBiosSettings)

	e := httptest.New(t, mockApp)
	e.PATCH(
		"/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1/Bios/Settings",
	).WithJSON(map[string]string{"Sample": "Body"}).WithHeader("X-Auth-Token", "InvalidToken").Expect().Status(http.StatusUnauthorized)
}

func TestChangeBiosSettingsNegativeTestCases(t *testing.T) {
	var sys SystemRPCs
	sys.ChangeBiosSettingsRPC = mockChangeBiosSettings
	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1/Systems")
	redfishRoutes.Patch("/{id}/Bios/Settings", sys.ChangeBiosSettings)

	e := httptest.New(t, mockApp)
	e.PATCH(
		"/redfish/v1/Systems//Bios/Settings",
	).WithJSON(map[string]string{"Sample": "Body"}).WithHeader("X-Auth-Token", "ValidToken").Expect().Status(http.StatusBadRequest)
	e.PATCH(
		"/redfish/v1/Systems//Bios/Settings",
	).WithHeader("X-Auth-Token", "ValidToken").Expect().Status(http.StatusBadRequest)
	e.PATCH(
		"/redfish/v1/Systems//Bios/Settings",
	).WithJSON(map[string]string{"Sample": "Body"}).WithHeader("X-Auth-Token", "TokenRPC").Expect().Status(http.StatusInternalServerError)
}

func mockChangeBootOrderSettings(req systemsproto.BootOrderSettingsRequest) (*systemsproto.SystemsResponse, error) {
	var response = &systemsproto.SystemsResponse{}
	if req.SessionToken == "" {
		response = &systemsproto.SystemsResponse{
			StatusCode:    401,
			StatusMessage: "Unauthorized",
			Body:          []byte(`{"Response":"Unauthorized"}`),
		}
	} else if req.SessionToken == "InvalidToken" {
		response = &systemsproto.SystemsResponse{
			StatusCode:    401,
			StatusMessage: "Unauthorized",
			Body:          []byte(`{"Response":"Unauthorized"}`),
		}
	} else if req.SessionToken == "TokenRPC" {
		return &systemsproto.SystemsResponse{}, errors.New("Unable to RPC Call")
	} else if req.SessionToken == "ValidToken" {
		response = &systemsproto.SystemsResponse{
			StatusCode:    200,
			StatusMessage: "Success",
			Body:          []byte(`{"Response":"Success"}`),
		}
	}
	return response, nil
}

func TestChangeBootOrderSettingsWithValidToken(t *testing.T) {
	var sys SystemRPCs
	sys.ChangeBootOrderSettingsRPC = mockChangeBootOrderSettings
	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1/Systems")
	redfishRoutes.Patch("/{id}", sys.ChangeBootOrderSettings)

	e := httptest.New(t, mockApp)
	e.PATCH(
		"/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1",
	).WithJSON(map[string]string{"Sample": "Body"}).WithHeader("X-Auth-Token", "ValidToken").Expect().Status(http.StatusOK)
}

func TestChangeBootOrderSettingsWithoutToken(t *testing.T) {
	var sys SystemRPCs
	sys.ChangeBootOrderSettingsRPC = mockChangeBootOrderSettings
	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1/Systems")
	redfishRoutes.Patch("/{id}", sys.ChangeBootOrderSettings)

	e := httptest.New(t, mockApp)
	e.PATCH(
		"/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1",
	).WithJSON(map[string]string{"Sample": "Body"}).WithHeader("X-Auth-Token", "").Expect().Status(http.StatusUnauthorized)
}

func TestChangeBootOrderSettingsNegativeTestCases(t *testing.T) {
	var sys SystemRPCs
	sys.ChangeBootOrderSettingsRPC = mockChangeBootOrderSettings
	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1/Systems")
	redfishRoutes.Patch("/{id}", sys.ChangeBootOrderSettings)

	e := httptest.New(t, mockApp)
	e.PATCH(
		"/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1",
	).WithJSON(map[string]string{"Sample": "Body"}).WithHeader("X-Auth-Token", "InvalidToken").Expect().Status(http.StatusUnauthorized)
	e.PATCH(
		"/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1",
	).WithJSON(map[string]string{"Sample": "Body"}).WithHeader("X-Auth-Token", "TokenRPC").Expect().Status(http.StatusInternalServerError)
	e.PATCH(
		"/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1",
	).WithHeader("X-Auth-Token", "Token").Expect().Status(http.StatusBadRequest)
}

func mockComputerSystemReset(req systemsproto.ComputerSystemResetRequest) (*systemsproto.SystemsResponse, error) {
	var response = &systemsproto.SystemsResponse{}
	if req.SessionToken == "" {
		response = &systemsproto.SystemsResponse{
			StatusCode:    401,
			StatusMessage: "Unauthorized",
			Body:          []byte(`{"Response":"Unauthorized"}`),
		}
	} else if req.SessionToken == "TokenRPC" {
		return &systemsproto.SystemsResponse{}, errors.New("Unable to RPC Call")
	} else {
		response = &systemsproto.SystemsResponse{
			StatusCode:    200,
			StatusMessage: "Success",
			Body:          []byte(`{"Response":"Success"}`),
		}
	}
	return response, nil
}

func TestComputerSystemResetNegativeTestCases(t *testing.T) {
	var sys SystemRPCs
	sys.SystemResetRPC = mockComputerSystemReset
	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1/Systems")
	redfishRoutes.Post("/{id}/Actions/ComputerSystem.Reset", sys.ComputerSystemReset)

	e := httptest.New(t, mockApp)
	e.POST(
		"/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1/Actions/ComputerSystem.Reset",
	).WithJSON(map[string]string{"Sample": "Body"}).WithHeader("X-Auth-Token", "").Expect().Status(http.StatusUnauthorized)
	e.POST(
		"/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1/Actions/ComputerSystem.Reset",
	).WithJSON("Body").WithHeader("X-Auth-Token", "").Expect().Status(http.StatusUnauthorized)
	e.POST(
		"/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1/Actions/ComputerSystem.Reset",
	).WithJSON(map[string]string{"Sample": "Body"}).WithHeader("X-Auth-Token", "TokenRPC").Expect().Status(http.StatusInternalServerError)
}

func TestComputerSystemResetWithValidData(t *testing.T) {
	var sys SystemRPCs
	sys.SystemResetRPC = mockComputerSystemReset
	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1/Systems")
	redfishRoutes.Post("/{id}/Actions/ComputerSystem.Reset", sys.ComputerSystemReset)

	e := httptest.New(t, mockApp)
	e.POST(
		"/redfish/v1/Systems/123/Actions/ComputerSystem.Reset",
	).WithJSON(map[string]string{"Sample": "Body"}).WithHeader("X-Auth-Token", "ValidToken").Expect().Status(http.StatusOK)
}

func mockSetDefaultBootOrder(req systemsproto.DefaultBootOrderRequest) (*systemsproto.SystemsResponse, error) {
	var response = &systemsproto.SystemsResponse{}
	if req.SessionToken == "" {
		response = &systemsproto.SystemsResponse{
			StatusCode:    401,
			StatusMessage: "Unauthorized",
			Body:          []byte(`{"Response":"Unauthorized"}`),
		}
	} else if req.SessionToken == "ValidToken" {
		response = &systemsproto.SystemsResponse{
			StatusCode:    200,
			StatusMessage: "Success",
			Body:          []byte(`{"ResetType":"ForceRestart"}`),
		}
	} else if req.SessionToken == "TokenRPC" {
		return &systemsproto.SystemsResponse{}, errors.New("Unable to RPC Call")
	}
	return response, nil
}

func TestMockSetDefaultBootOrderWithoutToken(t *testing.T) {
	var sys SystemRPCs
	sys.SetDefaultBootOrderRPC = mockSetDefaultBootOrder
	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1/Systems")
	redfishRoutes.Patch("/{id}/Actions/ComputerSystem.SetDefaultBootOrder", sys.SetDefaultBootOrder)

	e := httptest.New(t, mockApp)
	e.PATCH(
		"/redfish/v1/Systems/6d4a0a66:1/Actions/ComputerSystem.SetDefaultBootOrder",
	).WithJSON(map[string]string{"Sample": "Body"}).WithHeader("X-Auth-Token", "").Expect().Status(http.StatusUnauthorized)
	e.PATCH(
		"/redfish/v1/Systems/6d4a0a66:1/Actions/ComputerSystem.SetDefaultBootOrder",
	).WithJSON(map[string]string{"Sample": "Body"}).WithHeader("X-Auth-Token", "TokenRPC").Expect().Status(http.StatusInternalServerError)
}

func TestSetDefaultBootOrderWithValidData(t *testing.T) {
	var sys SystemRPCs
	sys.SetDefaultBootOrderRPC = mockSetDefaultBootOrder
	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1/Systems")
	redfishRoutes.Patch("/{id}/Actions/ComputerSystem.SetDefaultBootOrder", sys.SetDefaultBootOrder)

	e := httptest.New(t, mockApp)
	e.PATCH(
		"/redfish/v1/Systems/123:1/Actions/ComputerSystem.SetDefaultBootOrder",
	).WithJSON(map[string]string{"Sample": "Body"}).WithHeader("X-Auth-Token", "ValidToken").Expect().Status(http.StatusOK)
}

// Create volume unit tests
func mockCreateVolume(req systemsproto.VolumeRequest) (*systemsproto.SystemsResponse, error) {
	var response = &systemsproto.SystemsResponse{}
	if req.SessionToken == "" {
		response = &systemsproto.SystemsResponse{
			StatusCode:    http.StatusUnauthorized,
			StatusMessage: "Unauthorized",
			Body:          []byte(`{"Response":"Unauthorized"}`),
		}
	} else if req.SessionToken == "InvalidToken" {
		response = &systemsproto.SystemsResponse{
			StatusCode:    http.StatusUnauthorized,
			StatusMessage: "Unauthorized",
			Body:          []byte(`{"Response":"Unauthorized"}`),
		}
	} else if req.SessionToken == "ValidToken" && req.StorageInstance == "" {
		response = &systemsproto.SystemsResponse{
			StatusCode:    http.StatusBadRequest,
			StatusMessage: "BadRequest",
			Body:          []byte(`{"Response":"BadRequest"}`),
		}
	} else if req.SessionToken == "TokenRPC" {
		return &systemsproto.SystemsResponse{}, errors.New("Unable to RPC Call")
	} else {
		response = &systemsproto.SystemsResponse{
			StatusCode:    http.StatusOK,
			StatusMessage: "Success",
			Body:          []byte(`{"Response":"Success"}`),
		}
	}
	return response, nil
}

func TestCreateVolume(t *testing.T) {
	var sys SystemRPCs
	sys.CreateVolumeRPC = mockCreateVolume
	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1/Systems/{id}/Storage")
	redfishRoutes.Post("/{id2}/Volumes", sys.CreateVolume)

	e := httptest.New(t, mockApp)
	e.POST(
		"/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1/Storage/ArrayControllers-0/Volumes",
	).WithJSON(map[string]string{"Sample": "Body"}).WithHeader("X-Auth-Token", "ValidToken").Expect().Status(http.StatusOK)
}

func TestCreateVolumeWithoutToken(t *testing.T) {
	var sys SystemRPCs
	sys.CreateVolumeRPC = mockCreateVolume
	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1/Systems/{id}/Storage")
	redfishRoutes.Patch("/{id2}/Volumes", sys.CreateVolume)

	e := httptest.New(t, mockApp)
	e.POST(
		"/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1/Storage/ArrayControllers-0/Volumes",
	).WithJSON(map[string]string{"Sample": "Body"}).WithHeader("X-Auth-Token", "").Expect().Status(http.StatusNotFound)
}

func TestCreateVolumeWithInvalidToken(t *testing.T) {
	var sys SystemRPCs
	sys.CreateVolumeRPC = mockCreateVolume
	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1/Systems/{id}/Storage")
	redfishRoutes.Post("/{id2}/Volumes", sys.CreateVolume)

	e := httptest.New(t, mockApp)
	e.POST(
		"/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1/Storage/ArrayControllers-0/Volumes",
	).WithJSON(map[string]string{"Sample": "Body"}).WithHeader("X-Auth-Token", "InvalidToken").Expect().Status(http.StatusUnauthorized)
}

func TestCreateVolumeNegativeTestCases(t *testing.T) {
	var sys SystemRPCs
	sys.CreateVolumeRPC = mockCreateVolume
	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1/Systems/{id}/Storage")
	redfishRoutes.Post("/{id2}/Volumes", sys.CreateVolume)

	e := httptest.New(t, mockApp)
	e.POST(
		"/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1/Storage//Volumes",
	).WithJSON(map[string]string{"Sample": "Body"}).WithHeader("X-Auth-Token", "ValidToken").Expect().Status(http.StatusBadRequest)
	e.POST(
		"/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1/Storage/ArrayControllers-0/Volumes",
	).WithHeader("X-Auth-Token", "ValidToken").Expect().Status(http.StatusBadRequest)
	e.POST(
		"/redfish/v1/Systems//Storage/ArrayControllers-0/Volumes",
	).WithJSON(map[string]string{"Sample": "Body"}).WithHeader("X-Auth-Token", "TokenRPC").Expect().Status(http.StatusInternalServerError)
}

func mockDeleteVolume(req systemsproto.VolumeRequest) (*systemsproto.SystemsResponse, error) {
	var response = &systemsproto.SystemsResponse{}
	if req.SessionToken == "" {
		response = &systemsproto.SystemsResponse{
			StatusCode:    http.StatusUnauthorized,
			StatusMessage: "Unauthorized",
			Body:          []byte(`{"Response":"Unauthorized"}`),
		}
	} else if req.SessionToken == "InvalidToken" {
		response = &systemsproto.SystemsResponse{
			StatusCode:    http.StatusUnauthorized,
			StatusMessage: "Unauthorized",
			Body:          []byte(`{"Response":"Unauthorized"}`),
		}
	} else if req.SessionToken == "ValidToken" && req.VolumeID == "2" {
		response = &systemsproto.SystemsResponse{
			StatusCode:    http.StatusNotFound,
			StatusMessage: "NotFound",
			Body:          []byte(`{"Response":"NotFound"}`),
		}
	} else if req.SessionToken == "TokenRPC" {
		return &systemsproto.SystemsResponse{}, errors.New("Unable to RPC Call")
	} else {
		response = &systemsproto.SystemsResponse{
			StatusCode:    http.StatusOK,
			StatusMessage: "Success",
			Body:          []byte(`{"Response":"Success"}`),
		}
	}
	return response, nil
}

func TestDeleteVolume(t *testing.T) {
	var sys SystemRPCs
	sys.DeleteVolumeRPC = mockDeleteVolume
	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1/Systems/{id}/Storage/{id2}")
	redfishRoutes.Delete("/Volumes/{rid}", sys.DeleteVolume)

	e := httptest.New(t, mockApp)
	// test with valid token
	e.DELETE(
		"/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1/Storage/ArrayControllers-0/Volumes/1",
	).WithJSON(map[string]string{"Sample": "Body"}).WithHeader("X-Auth-Token", "ValidToken").Expect().Status(http.StatusOK)

	// test with Invalid token
	e.DELETE(
		"/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1/Storage/ArrayControllers-0/Volumes/1",
	).WithJSON(map[string]string{"Sample": "Body"}).WithHeader("X-Auth-Token", "InvalidToken").Expect().Status(http.StatusUnauthorized)

	// test without token
	e.DELETE(
		"/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1/Storage/ArrayControllers-0/Volumes/1",
	).WithJSON(map[string]string{"Sample": "Body"}).WithHeader("X-Auth-Token", "").Expect().Status(http.StatusUnauthorized)

	// test with invalid volume id
	e.DELETE(
		"/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1/Storage/ArrayControllers-0/Volumes/2",
	).WithJSON(map[string]string{"Sample": "Body"}).WithHeader("X-Auth-Token", "ValidToken").Expect().Status(http.StatusNotFound)

	// test with rpc error
	e.DELETE(
		"/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1/Storage/ArrayControllers-0/Volumes/2",
	).WithJSON(map[string]string{"Sample": "Body"}).WithHeader("X-Auth-Token", "TokenRPC").Expect().Status(http.StatusInternalServerError)
}
