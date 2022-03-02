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

	accountproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/account"
	iris "github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/httptest"
)

var header = map[string][]string{
	"Connection":             {"keep-alive"},
	"Odata-Version":          {"4.0"},
	"X-Frame-Options":        {"sameorigin"},
	"Content-Type":           {"application/json; charset=utf-8"},
	"X-Content-Type-Options": {"nosniff"},
	"Cache-Control":          {"no-cache, no-store, must-revalidate"},
	"Transfer-Encoding":      {"chunked"},
}

func mockGetAccountServiceRPC(req accountproto.AccountRequest) (*accountproto.AccountResponse, error) {
	if req.SessionToken == "TokenRPC" {
		return nil, errors.New("RPC Error")
	}
	return &accountproto.AccountResponse{
		StatusCode: http.StatusOK,
	}, nil
}

func mockCreateAccountRPC(req accountproto.CreateAccountRequest) (*accountproto.AccountResponse, error) {
	if req.SessionToken == "TokenRPC" {
		return nil, errors.New("RPC Error")
	}
	return &accountproto.AccountResponse{
		StatusCode: http.StatusCreated,
	}, nil
}

func mockGetAllAccountsRPC(req accountproto.AccountRequest) (*accountproto.AccountResponse, error) {
	if req.SessionToken == "TokenRPC" {
		return nil, errors.New("RPC Error")
	}
	return &accountproto.AccountResponse{
		StatusCode: http.StatusOK,
	}, nil
}

func mockGetAccountRPC(req accountproto.GetAccountRequest) (*accountproto.AccountResponse, error) {
	if req.SessionToken == "TokenRPC" {
		return nil, errors.New("RPC Error")
	}
	return &accountproto.AccountResponse{
		StatusCode: http.StatusOK,
	}, nil
}

func mockUpdateAccountRPC(req accountproto.UpdateAccountRequest) (*accountproto.AccountResponse, error) {
	if req.SessionToken == "TokenRPC" {
		return nil, errors.New("RPC Error")
	}
	return &accountproto.AccountResponse{
		StatusCode: http.StatusOK,
	}, nil
}

func mockDeleteAccountRPC(req accountproto.DeleteAccountRequest) (*accountproto.AccountResponse, error) {
	if req.SessionToken == "TokenRPC" {
		return nil, errors.New("RPC Error")
	}
	return &accountproto.AccountResponse{
		StatusCode: http.StatusOK,
	}, nil
}

func TestAccountRPCs_GetAccountService(t *testing.T) {
	header["Allow"] = []string{"GET"}
	defer delete(header, "Allow")
	var a AccountRPCs
	a.GetServiceRPC = mockGetAccountServiceRPC
	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1")
	redfishRoutes.Get("/AccountService", a.GetAccountService)

	e := httptest.New(t, mockApp)
	e.GET(
		"/redfish/v1/AccountService",
	).WithHeader("X-Auth-Token", "token").Expect().Status(http.StatusOK).Headers().Equal(header)

	e.GET(
		"/redfish/v1/AccountService",
	).WithHeader("X-Auth-Token", "TokenRPC").Expect().Status(http.StatusInternalServerError)
}

func TestAccountRPCs_GetAccountServiceWithOutToken(t *testing.T) {
	var a AccountRPCs
	a.GetServiceRPC = mockGetAccountServiceRPC

	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1")
	redfishRoutes.Get("/AccountService", a.GetAccountService)

	e := httptest.New(t, mockApp)
	e.GET(
		"/redfish/v1/AccountService",
	).Expect().Status(http.StatusUnauthorized).Headers().Equal(header)
}

func TestAccountRPCs_CreateAccount(t *testing.T) {
	var a AccountRPCs
	a.CreateRPC = mockCreateAccountRPC

	body := map[string]interface{}{
		"UserName": "someUser",
		"Password": "somePassword",
		"RoleId":   "someRole",
	}
	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1")
	redfishRoutes.Post("/AccountService/Accounts", a.CreateAccount)

	e := httptest.New(t, mockApp)
	e.POST(
		"/redfish/v1/AccountService/Accounts",
	).WithHeader("X-Auth-Token", "Token").WithJSON(body).Expect().Status(http.StatusCreated).Headers().Equal(header)
	e.POST(
		"/redfish/v1/AccountService/Accounts",
	).WithHeader("X-Auth-Token", "").WithJSON(body).Expect().Status(http.StatusUnauthorized)
	e.POST(
		"/redfish/v1/AccountService/Accounts",
	).WithHeader("X-Auth-Token", "Token").Expect().Status(http.StatusBadRequest)
	e.POST(
		"/redfish/v1/AccountService/Accounts",
	).WithHeader("X-Auth-Token", "TokenRPC").WithJSON(body).Expect().Status(http.StatusInternalServerError)

}

func TestAccountRPCs_GetAllAccounts(t *testing.T) {
	header["Allow"] = []string{"GET, POST"}
	defer delete(header, "Allow")
	var a AccountRPCs
	a.GetAllAccountsRPC = mockGetAllAccountsRPC

	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1")
	redfishRoutes.Get("/AccountService/Accounts", a.GetAllAccounts)

	e := httptest.New(t, mockApp)
	e.GET(
		"/redfish/v1/AccountService/Accounts",
	).WithHeader("X-Auth-Token", "token").Expect().Status(http.StatusOK).Headers().Equal(header)
	e.GET(
		"/redfish/v1/AccountService/Accounts",
	).WithHeader("X-Auth-Token", "").Expect().Status(http.StatusUnauthorized)
	e.GET(
		"/redfish/v1/AccountService/Accounts",
	).WithHeader("X-Auth-Token", "TokenRPC").Expect().Status(http.StatusInternalServerError)
}

func TestAccountRPCs_GetAccount(t *testing.T) {
	header["Allow"] = []string{"GET, PATCH, DELETE"}
	defer delete(header, "Allow")
	var a AccountRPCs
	a.GetAccountRPC = mockGetAccountRPC

	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1")
	redfishRoutes.Get("/AccountService/Accounts/{id}", a.GetAccount)

	e := httptest.New(t, mockApp)
	e.GET(
		"/redfish/v1/AccountService/Accounts/someID",
	).WithHeader("X-Auth-Token", "token").Expect().Status(http.StatusOK).Headers().Equal(header)
	e.GET(
		"/redfish/v1/AccountService/Accounts/someID",
	).WithHeader("X-Auth-Token", "TokenRPC").Expect().Status(http.StatusInternalServerError)
}

func TestAccountRPCs_GetAccountWithOutToken(t *testing.T) {
	var a AccountRPCs
	a.GetAccountRPC = mockGetAccountRPC

	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1")
	redfishRoutes.Get("/AccountService/Accounts/{id}", a.GetAccount)

	e := httptest.New(t, mockApp)
	e.GET(
		"/redfish/v1/AccountService/Accounts/someID",
	).Expect().Status(http.StatusUnauthorized)
}

func TestAccountRPCs_UpdateAccount(t *testing.T) {
	var a AccountRPCs
	a.UpdateRPC = mockUpdateAccountRPC

	body := map[string]interface{}{
		"Password": "somePassword",
		"RoleId":   "someRole",
	}
	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1")
	redfishRoutes.Patch("/AccountService/Accounts/{id}", a.UpdateAccount)

	e := httptest.New(t, mockApp)
	e.PATCH(
		"/redfish/v1/AccountService/Accounts/someID",
	).WithHeader("X-Auth-Token", "token").WithJSON(body).Expect().Status(http.StatusOK)
	e.PATCH(
		"/redfish/v1/AccountService/Accounts/someID",
	).WithHeader("X-Auth-Token", "").WithJSON(body).Expect().Status(http.StatusUnauthorized)
	e.PATCH(
		"/redfish/v1/AccountService/Accounts/someID",
	).WithHeader("X-Auth-Token", "").Expect().Status(http.StatusBadRequest)
	e.PATCH(
		"/redfish/v1/AccountService/Accounts/someID",
	).WithHeader("X-Auth-Token", "TokenRPC").WithJSON(body).Expect().Status(http.StatusInternalServerError)
}

func TestAccountRPCs_DeleteAccount(t *testing.T) {
	var a AccountRPCs
	a.DeleteRPC = mockDeleteAccountRPC

	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1")
	redfishRoutes.Get("/AccountService/Accounts/{id}", a.DeleteAccount)

	e := httptest.New(t, mockApp)
	e.GET(
		"/redfish/v1/AccountService/Accounts/someID",
	).WithHeader("X-Auth-Token", "token").Expect().Status(http.StatusOK)
	e.GET(
		"/redfish/v1/AccountService/Accounts/someID",
	).WithHeader("X-Auth-Token", "TokenRPC").Expect().Status(http.StatusInternalServerError)
}

func TestAccountRPCs_DeleteAccountWithOutToken(t *testing.T) {
	var a AccountRPCs
	a.DeleteRPC = mockDeleteAccountRPC

	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1")
	redfishRoutes.Delete("/AccountService/Accounts/{id}", a.DeleteAccount)

	e := httptest.New(t, mockApp)
	e.DELETE(
		"/redfish/v1/AccountService/Accounts/someID",
	).Expect().Status(http.StatusUnauthorized)
}
