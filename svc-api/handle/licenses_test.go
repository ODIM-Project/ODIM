//(C) Copyright [2022] Hewlett Packard Enterprise Development LP
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

//Package handle ...

package handle

import (
	"errors"
	"net/http"
	"testing"

	licenseproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/licenses"
	iris "github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/httptest"
)

func testLicenseService(req licenseproto.GetLicenseServiceRequest) (*licenseproto.GetLicenseResponse, error) {
	var response = &licenseproto.GetLicenseResponse{}
	if req.SessionToken == "ValidToken" {
		response = &licenseproto.GetLicenseResponse{
			StatusCode:    200,
			StatusMessage: "Success",
			Body:          []byte(`{"Response":"Success"}`),
		}
	} else if req.SessionToken == "InvalidToken" {
		response = &licenseproto.GetLicenseResponse{
			StatusCode:    401,
			StatusMessage: "Unauthorized", Body: []byte(`{"Response":"Unauthorized"}`),
		}
	} else if req.SessionToken == "token" {
		return &licenseproto.GetLicenseResponse{}, errors.New("Unable to RPC Call")
	}
	return response, nil
}

func testLicenseCollection(req licenseproto.GetLicenseRequest) (*licenseproto.GetLicenseResponse, error) {
	var response = &licenseproto.GetLicenseResponse{}
	if req.SessionToken == "ValidToken" {
		response = &licenseproto.GetLicenseResponse{
			StatusCode:    200,
			StatusMessage: "Success",
			Body:          []byte(`{"Response":"Success"}`),
		}
	} else if req.SessionToken == "InvalidToken" {
		response = &licenseproto.GetLicenseResponse{
			StatusCode:    401,
			StatusMessage: "Unauthorized", Body: []byte(`{"Response":"Unauthorized"}`),
		}
	} else if req.SessionToken == "token" {
		return &licenseproto.GetLicenseResponse{}, errors.New("Unable to RPC Call")
	}
	return response, nil
}

func testLicenseResource(req licenseproto.GetLicenseResourceRequest) (*licenseproto.GetLicenseResponse, error) {
	var response = &licenseproto.GetLicenseResponse{}
	if req.SessionToken == "ValidToken" {
		response = &licenseproto.GetLicenseResponse{
			StatusCode:    200,
			StatusMessage: "Success",
			Body:          []byte(`{"Response":"Success"}`),
		}
	} else if req.SessionToken == "InvalidToken" {
		response = &licenseproto.GetLicenseResponse{
			StatusCode:    401,
			StatusMessage: "Unauthorized", Body: []byte(`{"Response":"Unauthorized"}`),
		}
	} else if req.SessionToken == "token" {
		return &licenseproto.GetLicenseResponse{}, errors.New("Unable to RPC Call")
	}
	return response, nil
}

func testInstallLicenseService(req licenseproto.InstallLicenseRequest) (*licenseproto.GetLicenseResponse, error) {
	var response = &licenseproto.GetLicenseResponse{}
	if req.SessionToken == "ValidToken" {
		response = &licenseproto.GetLicenseResponse{
			StatusCode: 204,
		}
	} else if req.SessionToken == "InvalidToken" {
		response = &licenseproto.GetLicenseResponse{
			StatusCode:    401,
			StatusMessage: "Unauthorized", Body: []byte(`{"Response":"Unauthorized"}`),
		}
	} else if req.SessionToken == "token" {
		return &licenseproto.GetLicenseResponse{}, errors.New("Unable to RPC Call")
	}
	return response, nil
}

func TestGetLicenseService(t *testing.T) {
	header["Allow"] = []string{"GET"}
	defer delete(header, "Allow")
	var a LicenseRPCs
	a.GetLicenseServiceRPC = testLicenseService
	testApp := iris.New()
	redfishRoutes := testApp.Party("/redfish/v1/LicenseService")
	redfishRoutes.Get("/", a.GetLicenseService)
	test := httptest.New(t, testApp)
	test.GET(
		"/redfish/v1/LicenseService",
	).WithHeader("X-Auth-Token", "ValidToken").Expect().Status(http.StatusOK).Headers().Equal(header)
	test.GET(
		"/redfish/v1/LicenseService",
	).WithHeader("X-Auth-Token", "").Expect().Status(http.StatusUnauthorized)
	test.GET(
		"/redfish/v1/LicenseService",
	).WithHeader("X-Auth-Token", "token").Expect().Status(http.StatusInternalServerError)
}

func TestGetLicenseCollection(t *testing.T) {
	header["Allow"] = []string{"GET, POST"}
	defer delete(header, "Allow")
	var a LicenseRPCs
	a.GetLicenseCollectionRPC = testLicenseCollection
	testApp := iris.New()
	redfishRoutes := testApp.Party("/redfish/v1/LicenseService")
	redfishRoutes.Get("/Licenses", a.GetLicenseCollection)
	test := httptest.New(t, testApp)
	test.GET(
		"/redfish/v1/LicenseService/Licenses",
	).WithHeader("X-Auth-Token", "ValidToken").Expect().Status(http.StatusOK).Headers().Equal(header)
	test.GET(
		"/redfish/v1/LicenseService/Licenses",
	).WithHeader("X-Auth-Token", "").Expect().Status(http.StatusUnauthorized)
	test.GET(
		"/redfish/v1/LicenseService/Licenses",
	).WithHeader("X-Auth-Token", "token").Expect().Status(http.StatusInternalServerError)
}

func TestGetLicenseResource(t *testing.T) {
	header["Allow"] = []string{"GET"}
	defer delete(header, "Allow")
	var a LicenseRPCs
	a.GetLicenseResourceRPC = testLicenseResource
	testApp := iris.New()
	redfishRoutes := testApp.Party("/redfish/v1/LicenseService")
	redfishRoutes.Get("/Licenses/{id}", a.GetLicenseResource)
	test := httptest.New(t, testApp)
	test.GET(
		"/redfish/v1/LicenseService/Licenses/1",
	).WithHeader("X-Auth-Token", "ValidToken").Expect().Status(http.StatusOK).Headers().Equal(header)
	test.GET(
		"/redfish/v1/LicenseService/Licenses/1",
	).WithHeader("X-Auth-Token", "").Expect().Status(http.StatusUnauthorized)
	test.GET(
		"/redfish/v1/LicenseService/Licenses/1",
	).WithHeader("X-Auth-Token", "token").Expect().Status(http.StatusInternalServerError)
}

func TestInstallLicenseService(t *testing.T) {
	var a LicenseRPCs
	a.InstallLicenseServiceRPC = testInstallLicenseService
	testApp := iris.New()
	redfishRoutes := testApp.Party("/redfish/v1/LicenseService")
	redfishRoutes.Post("/Licenses", a.InstallLicenseService)
	test := httptest.New(t, testApp)
	test.POST(
		"/redfish/v1/LicenseService/Licenses",
	).WithHeader("X-Auth-Token", "ValidToken").WithJSON(map[string]interface{}{
		"AuthorizedDevices": "/redfish/v1/Managers/{id}",
		"LicenseString":     "MzMzSzItOFFMVjQtWThSM0ctTEpRUVgtN0JLNk0=",
	}).Expect().Status(http.StatusNoContent).Headers().Equal(header)
	test.POST(
		"/redfish/v1/LicenseService/Licenses",
	).WithHeader("X-Auth-Token", "").WithJSON(map[string]interface{}{
		"AuthorizedDevices": "/redfish/v1/Managers/{id}",
		"LicenseString":     "MzMzSzItOFFMVjQtWThSM0ctTEpRUVgtN0JLNk0=",
	}).Expect().Status(http.StatusUnauthorized)
	test.POST(
		"/redfish/v1/LicenseService/Licenses",
	).WithHeader("X-Auth-Token", "token").WithJSON(map[string]interface{}{
		"AuthorizedDevices": "/redfish/v1/Managers/{id}",
		"LicenseString":     "MzMzSzItOFFMVjQtWThSM0ctTEpRUVgtN0JLNk0=",
	}).Expect().Status(http.StatusInternalServerError)
}
