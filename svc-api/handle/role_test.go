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

	roleproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/role"
	iris "github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/httptest"
)

func mockGetAllRolesRPC(roleproto.GetRoleRequest) (*roleproto.RoleResponse, error) {
	return &roleproto.RoleResponse{
		StatusCode: http.StatusOK,
	}, nil
}
func mockGetAllRolesRPCWithRPCError(roleproto.GetRoleRequest) (*roleproto.RoleResponse, error) {
	return &roleproto.RoleResponse{}, errors.New("Unable to RPC Call")
}

func mockCreateRoleRPC(roleproto.RoleRequest) (*roleproto.RoleResponse, error) {
	return &roleproto.RoleResponse{
		StatusCode: http.StatusCreated,
	}, nil
}

func mockCreateRoleRPCWithRPCError(roleproto.RoleRequest) (*roleproto.RoleResponse, error) {
	return &roleproto.RoleResponse{}, errors.New("Unable to RPC Call")
}

func mockGetRoleRPC(roleproto.GetRoleRequest) (*roleproto.RoleResponse, error) {
	return &roleproto.RoleResponse{
		StatusCode: http.StatusOK,
	}, nil
}

func mockGetRoleRPCWithRPCError(roleproto.GetRoleRequest) (*roleproto.RoleResponse, error) {
	return &roleproto.RoleResponse{}, errors.New("Unable to RPC Call")
}

func mockUpdateRoleRPC(roleproto.UpdateRoleRequest) (*roleproto.RoleResponse, error) {
	return &roleproto.RoleResponse{
		StatusCode: http.StatusOK,
	}, nil
}

func mockUpdateRoleRPCWithRPCError(roleproto.UpdateRoleRequest) (*roleproto.RoleResponse, error) {
	return &roleproto.RoleResponse{}, errors.New("Unable to RPC Call")
}
func mockDeleteRoleRPC(roleproto.DeleteRoleRequest) (*roleproto.RoleResponse, error) {
	return &roleproto.RoleResponse{
		StatusCode: http.StatusOK,
	}, nil
}

func TestRoleRPCs_GetAllRoles(t *testing.T) {
	header["Allow"] = []string{"GET, POST"}
	defer delete(header, "Allow")
	var r RoleRPCs
	r.GetAllRolesRPC = mockGetAllRolesRPC
	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1")
	redfishRoutes.Get("/AccountService/Roles", r.GetAllRoles)

	e := httptest.New(t, mockApp)
	e.GET(
		"/redfish/v1/AccountService/Roles",
	).WithHeader("X-Auth-Token", "token").Expect().Status(http.StatusOK).Headers().Equal(header)
}

func TestRoleRPCs_GetAllRolesWithRPCError(t *testing.T) {
	var r RoleRPCs
	r.GetAllRolesRPC = mockGetAllRolesRPCWithRPCError
	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1")
	redfishRoutes.Get("/AccountService/Roles", r.GetAllRoles)

	e := httptest.New(t, mockApp)
	e.GET(
		"/redfish/v1/AccountService/Roles",
	).WithHeader("X-Auth-Token", "token").Expect().Status(http.StatusInternalServerError).Headers().Equal(header)
}

func TestRoleRPCs_GetAllRolesWithoutToken(t *testing.T) {
	var r RoleRPCs
	r.GetAllRolesRPC = mockGetAllRolesRPC
	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1")
	redfishRoutes.Get("/AccountService/Roles", r.GetAllRoles)

	e := httptest.New(t, mockApp)
	e.GET(
		"/redfish/v1/AccountService/Roles",
	).Expect().Status(http.StatusUnauthorized).Headers().Equal(header)
}

func TestRoleRPCs_CreateRole(t *testing.T) {
	var r RoleRPCs
	r.CreateRoleRPC = mockCreateRoleRPC
	body := map[string]interface{}{
		"RoleId":             "someRole",
		"AssignedPrivileges": []string{"SomePrivilege"},
	}

	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1")
	redfishRoutes.Post("/AccountService/Roles", r.CreateRole)

	e := httptest.New(t, mockApp)
	e.POST(
		"/redfish/v1/AccountService/Roles",
	).WithHeader("X-Auth-Token", "token").WithJSON(body).Expect().Status(http.StatusCreated).Headers().Equal(header)
	e.POST(
		"/redfish/v1/AccountService/Roles",
	).WithHeader("X-Auth-Token", "token").Expect().Status(http.StatusBadRequest)
	e.POST(
		"/redfish/v1/AccountService/Roles",
	).WithHeader("X-Auth-Token", "").WithJSON(body).Expect().Status(http.StatusUnauthorized)
}

func TestRoleRPCs_CreateRoleWithRPCError(t *testing.T) {
	var r RoleRPCs
	r.CreateRoleRPC = mockCreateRoleRPCWithRPCError
	body := map[string]interface{}{
		"RoleId":             "someRole",
		"AssignedPrivileges": []string{"SomePrivilege"},
	}

	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1")
	redfishRoutes.Post("/AccountService/Roles", r.CreateRole)

	e := httptest.New(t, mockApp)
	e.POST(
		"/redfish/v1/AccountService/Roles",
	).WithHeader("X-Auth-Token", "token").WithJSON(body).Expect().Status(http.StatusInternalServerError)
}

func TestRoleRPCs_GetRole(t *testing.T) {
	header["Allow"] = []string{"GET, PATCH, DELETE"}
	defer delete(header, "Allow")
	var r RoleRPCs
	r.GetRoleRPC = mockGetRoleRPC
	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1")
	redfishRoutes.Get("/AccountService/Roles/{id}", r.GetRole)

	e := httptest.New(t, mockApp)
	e.GET(
		"/redfish/v1/AccountService/Roles/someID",
	).WithHeader("X-Auth-Token", "token").Expect().Status(http.StatusOK).Headers().Equal(header)
}

func TestRoleRPCs_GetRoleWithRPCError(t *testing.T) {
	var r RoleRPCs
	r.GetRoleRPC = mockGetRoleRPCWithRPCError
	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1")
	redfishRoutes.Get("/AccountService/Roles/{id}", r.GetRole)

	e := httptest.New(t, mockApp)
	e.GET(
		"/redfish/v1/AccountService/Roles/someID",
	).WithHeader("X-Auth-Token", "token").Expect().Status(http.StatusInternalServerError)
}
func TestRoleRPCs_GetRoleWithoutToken(t *testing.T) {
	var r RoleRPCs
	r.GetRoleRPC = mockGetRoleRPC
	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1")
	redfishRoutes.Get("/AccountService/Roles/{id}", r.GetRole)

	e := httptest.New(t, mockApp)
	e.GET(
		"/redfish/v1/AccountService/Roles/someID",
	).Expect().Status(http.StatusUnauthorized)
}

func TestRoleRPCs_UpdateRole(t *testing.T) {
	var r RoleRPCs
	r.UpdateRoleRPC = mockUpdateRoleRPC
	body := map[string]interface{}{
		"AssignedPrivileges": []string{"SomePrivilege"},
	}

	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1")
	redfishRoutes.Patch("/AccountService/Roles/{id}", r.UpdateRole)

	e := httptest.New(t, mockApp)
	e.PATCH(
		"/redfish/v1/AccountService/Roles/someID",
	).WithHeader("X-Auth-Token", "token").WithJSON(body).Expect().Status(http.StatusOK).Headers().Equal(header)
	e.PATCH(
		"/redfish/v1/AccountService/Roles/someID",
	).Expect().Status(http.StatusBadRequest)
	e.PATCH(
		"/redfish/v1/AccountService/Roles/someID",
	).WithHeader("X-Auth-Token", "").WithJSON(body).Expect().Status(http.StatusUnauthorized)
}

func TestRoleRPCs_UpdateRoleWithRPCError(t *testing.T) {
	var r RoleRPCs
	r.UpdateRoleRPC = mockUpdateRoleRPCWithRPCError
	body := map[string]interface{}{
		"AssignedPrivileges": []string{"SomePrivilege"},
	}

	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1")
	redfishRoutes.Patch("/AccountService/Roles/{id}", r.UpdateRole)

	e := httptest.New(t, mockApp)
	e.PATCH(
		"/redfish/v1/AccountService/Roles/someID",
	).WithHeader("X-Auth-Token", "token").WithJSON(body).Expect().Status(http.StatusInternalServerError)
}

func TestDeleteRole(t *testing.T) {
	var r RoleRPCs
	r.DeleteRoleRPC = mockDeleteRoleRPC

	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1")
	redfishRoutes.Delete("/Role/{id}", r.DeleteRole)

	e := httptest.New(t, mockApp)

	// rpc call with proper details
	e.DELETE(
		"/redfish/v1/Role/SomeID",
	).WithHeader("X-Auth-Token", "token").Expect().Status(http.StatusOK).Headers().Equal(header)

	//rpc call with invalid session token
	e.DELETE(
		"/redfish/v1/Role/SomeID",
	).WithHeader("X-Auth-Token", "").Expect().Status(http.StatusUnauthorized)
}
