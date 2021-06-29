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

	fabricsproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/fabrics"
	iris "github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/httptest"
)

func mockGetFabricResource(req fabricsproto.FabricRequest) (*fabricsproto.FabricResponse, error) {
	var response = &fabricsproto.FabricResponse{}
	if req.SessionToken == "ValidToken" {
		response = &fabricsproto.FabricResponse{
			StatusCode:    200,
			StatusMessage: "Success",
			Body:          []byte(`{"Response":"Success"}`),
		}
	} else if req.SessionToken == "" {
		response = &fabricsproto.FabricResponse{
			StatusCode:    401,
			StatusMessage: "Unauthorized", Body: []byte(`{"Response":"Unauthorized"}`),
		}
	} else if req.SessionToken == "TokenRPC" {
		return nil, errors.New("RPC Error")
	}
	return response, nil
}

func mockDeleteFabricResource(req fabricsproto.FabricRequest) (*fabricsproto.FabricResponse, error) {
	var response = &fabricsproto.FabricResponse{}
	if req.SessionToken == "ValidToken" {
		response = &fabricsproto.FabricResponse{
			StatusCode:    200,
			StatusMessage: "Success",
			Body:          []byte(`{"Response":"Success"}`),
		}
	} else if req.SessionToken == "" {
		response = &fabricsproto.FabricResponse{
			StatusCode:    401,
			StatusMessage: "Unauthorized", Body: []byte(`{"Response":"Unauthorized"}`),
		}
	} else if req.SessionToken == "TokenRPC" {
		return nil, errors.New("RPC Error")
	}
	return response, nil
}

func mockUpdateFabricResource(req fabricsproto.FabricRequest) (*fabricsproto.FabricResponse, error) {
	var response = &fabricsproto.FabricResponse{}
	if req.SessionToken == "ValidToken" {
		response = &fabricsproto.FabricResponse{
			StatusCode:    200,
			StatusMessage: "Success",
			Body:          []byte(`{"Response":"Success"}`),
		}
	} else if req.SessionToken == "" {
		response = &fabricsproto.FabricResponse{
			StatusCode:    401,
			StatusMessage: "Unauthorized", Body: []byte(`{"Response":"Unauthorized"}`),
		}
	}
	return response, nil
}

func TestGetFabricResource(t *testing.T) {
	var fabrics FabricRPCs
	fabrics.GetFabricResourceRPC = mockGetFabricResource
	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1/Fabrics")
	redfishRoutes.Get("/", fabrics.GetFabricCollection)
	redfishRoutes.Get("/{id}", fabrics.GetFabric)
	redfishRoutes.Get("/{id}/Switches/{switchID}", fabrics.GetFabricCollection)
	redfishRoutes.Get("/{id}/Switches/{switchID}/Ports", fabrics.GetFabricCollection)
	redfishRoutes.Get("/{id}/Switches", fabrics.GetFabricCollection)
	redfishRoutes.Get("/{id}/Switches/{switchID}/Ports/{port_uuid}", fabrics.GetFabricCollection)
	redfishRoutes.Get("/{id}/Zones/", fabrics.GetFabricCollection)
	redfishRoutes.Get("/{id}/Endpoints/", fabrics.GetFabricCollection)
	redfishRoutes.Get("/{id}/AddressPools/", fabrics.GetFabricCollection)
	redfishRoutes.Get("/{id}/Zones/{zone_uuid}", fabrics.GetFabricCollection)
	redfishRoutes.Get("/{id}/Endpoints/{endpoint_uuid}", fabrics.GetFabricCollection)
	redfishRoutes.Get("/{id}/AddressPools/{addresspool_uuid}", fabrics.GetFabricCollection)
	test := httptest.New(t, mockApp)
	test.GET(
		"/redfish/v1/Fabrics/",
	).WithHeader("X-Auth-Token", "ValidToken").Expect().Status(http.StatusOK)
	test.GET(
		"/redfish/v1/Fabrics/",
	).WithHeader("X-Auth-Token", "").Expect().Status(http.StatusUnauthorized)
	test.GET(
		"/redfish/v1/Fabrics/",
	).WithHeader("X-Auth-Token", "TokenRPC").Expect().Status(http.StatusInternalServerError)
	test.GET(
		"/redfish/v1/Fabrics/1",
	).WithHeader("X-Auth-Token", "ValidToken").Expect().Status(http.StatusOK)
	test.GET(
		"/redfish/v1/Fabrics/1",
	).WithHeader("X-Auth-Token", "").Expect().Status(http.StatusUnauthorized)
	test.GET(
		"/redfish/v1/Fabrics/1/",
	).WithHeader("X-Auth-Token", "TokenRPC").Expect().Status(http.StatusInternalServerError)
	test.GET(
		"/redfish/v1/Fabrics/1/Switches",
	).WithHeader("X-Auth-Token", "ValidToken").Expect().Status(http.StatusOK)
	test.GET(
		"/redfish/v1/Fabrics/1/Switches",
	).WithHeader("X-Auth-Token", "").Expect().Status(http.StatusUnauthorized)
	test.GET(
		"/redfish/v1/Fabrics/1/Switches",
	).WithHeader("X-Auth-Token", "TokenRPC").Expect().Status(http.StatusInternalServerError)
	test.GET(
		"/redfish/v1/Fabrics/1/Switches/2",
	).WithHeader("X-Auth-Token", "ValidToken").Expect().Status(http.StatusOK)
	test.GET(
		"/redfish/v1/Fabrics/1/Switches/2",
	).WithHeader("X-Auth-Token", "").Expect().Status(http.StatusUnauthorized)
	test.GET(
		"/redfish/v1/Fabrics/1/Switches/2",
	).WithHeader("X-Auth-Token", "TokenRPC").Expect().Status(http.StatusInternalServerError)
	test.GET(
		"/redfish/v1/Fabrics/1/Switches/2/Ports",
	).WithHeader("X-Auth-Token", "ValidToken").Expect().Status(http.StatusOK)
	test.GET(
		"/redfish/v1/Fabrics/1/Switches/2/Ports",
	).WithHeader("X-Auth-Token", "").Expect().Status(http.StatusUnauthorized)
	test.GET(
		"/redfish/v1/Fabrics/1/Switches/2/Ports",
	).WithHeader("X-Auth-Token", "TokenRPC").Expect().Status(http.StatusInternalServerError)
	test.GET(
		"/redfish/v1/Fabrics/1/Switches/2/Ports/3",
	).WithHeader("X-Auth-Token", "ValidToken").Expect().Status(http.StatusOK)
	test.GET(
		"/redfish/v1/Fabrics/1/Switches/2/Ports/3",
	).WithHeader("X-Auth-Token", "").Expect().Status(http.StatusUnauthorized)
	test.GET(
		"/redfish/v1/Fabrics/1/Switches/2/Ports/3",
	).WithHeader("X-Auth-Token", "TokenRPC").Expect().Status(http.StatusInternalServerError)
	test.GET(
		"/redfish/v1/Fabrics/1/Zones/",
	).WithHeader("X-Auth-Token", "ValidToken").Expect().Status(http.StatusOK)
	test.GET(
		"/redfish/v1/Fabrics/1/Zones/",
	).WithHeader("X-Auth-Token", "").Expect().Status(http.StatusUnauthorized)
	test.GET(
		"/redfish/v1/Fabrics/1/Zones/",
	).WithHeader("X-Auth-Token", "TokenRPC").Expect().Status(http.StatusInternalServerError)
	test.GET(
		"/redfish/v1/Fabrics/1/Endpoints/",
	).WithHeader("X-Auth-Token", "ValidToken").Expect().Status(http.StatusOK)
	test.GET(
		"/redfish/v1/Fabrics/1/Endpoints/",
	).WithHeader("X-Auth-Token", "").Expect().Status(http.StatusUnauthorized)
	test.GET(
		"/redfish/v1/Fabrics/1/Endpoints/",
	).WithHeader("X-Auth-Token", "TokenRPC").Expect().Status(http.StatusInternalServerError)
	test.GET(
		"/redfish/v1/Fabrics/1/AddressPools/",
	).WithHeader("X-Auth-Token", "ValidToken").Expect().Status(http.StatusOK)
	test.GET(
		"/redfish/v1/Fabrics/1/AddressPools/",
	).WithHeader("X-Auth-Token", "").Expect().Status(http.StatusUnauthorized)
	test.GET(
		"/redfish/v1/Fabrics/1/AddressPools/",
	).WithHeader("X-Auth-Token", "TokenRPC").Expect().Status(http.StatusInternalServerError)
	test.GET(
		"/redfish/v1/Fabrics/1/Zones/2",
	).WithHeader("X-Auth-Token", "ValidToken").Expect().Status(http.StatusOK)
	test.GET(
		"/redfish/v1/Fabrics/1/Zones/2",
	).WithHeader("X-Auth-Token", "").Expect().Status(http.StatusUnauthorized)
	test.GET(
		"/redfish/v1/Fabrics/1/Zones/2",
	).WithHeader("X-Auth-Token", "TokenRPC").Expect().Status(http.StatusInternalServerError)
	test.GET(
		"/redfish/v1/Fabrics/1/Endpoints/2",
	).WithHeader("X-Auth-Token", "ValidToken").Expect().Status(http.StatusOK)
	test.GET(
		"/redfish/v1/Fabrics/1/Endpoints/2",
	).WithHeader("X-Auth-Token", "").Expect().Status(http.StatusUnauthorized)
	test.GET(
		"/redfish/v1/Fabrics/1/Endpoints/2",
	).WithHeader("X-Auth-Token", "TokenRPC").Expect().Status(http.StatusInternalServerError)
	test.GET(
		"/redfish/v1/Fabrics/1/AddressPools/2",
	).WithHeader("X-Auth-Token", "ValidToken").Expect().Status(http.StatusOK)
	test.GET(
		"/redfish/v1/Fabrics/1/AddressPools/2",
	).WithHeader("X-Auth-Token", "").Expect().Status(http.StatusUnauthorized)
	test.GET(
		"/redfish/v1/Fabrics/1/AddressPools/2",
	).WithHeader("X-Auth-Token", "TokenRPC").Expect().Status(http.StatusInternalServerError)
}

func TestDeleteFabricResource(t *testing.T) {
	var fabrics FabricRPCs
	fabrics.DeleteFabricResourceRPC = mockDeleteFabricResource
	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1/Fabrics")
	redfishRoutes.Get("/", fabrics.DeleteFabricResource)
	test := httptest.New(t, mockApp)
	test.GET(
		"/redfish/v1/Fabrics/",
	).WithHeader("X-Auth-Token", "ValidToken").Expect().Status(http.StatusOK)
	test.GET(
		"/redfish/v1/Fabrics/",
	).WithHeader("X-Auth-Token", "").Expect().Status(http.StatusUnauthorized)
	test.GET(
		"/redfish/v1/Fabrics/",
	).WithHeader("X-Auth-Token", "TokenRPC").Expect().Status(http.StatusInternalServerError)
}

func TestUpdateFabricResource(t *testing.T) {
	var fabrics FabricRPCs
	fabrics.UpdateFabricResourceRPC = mockUpdateFabricResource
	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1/Fabrics")
	redfishRoutes.Get("/", fabrics.UpdateFabricResource)
	test := httptest.New(t, mockApp)
	test.PATCH(
		"/redfish/v1/Fabrics/",
	).WithHeader("X-Auth-Token", "ValidToken").WithJSON("").Expect().Status(http.StatusBadRequest)
	test.PATCH(
		"/redfish/v1/Fabrics/",
	).WithHeader("X-Auth-Token", "").Expect().Status(http.StatusUnauthorized)
}
