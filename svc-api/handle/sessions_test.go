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

	sessionproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/session"
	iris "github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/httptest"
)

func mockCreateSessionRPC(req sessionproto.SessionCreateRequest) (*sessionproto.SessionCreateResponse, error) {
	return &sessionproto.SessionCreateResponse{
		StatusCode: http.StatusCreated,
		Header:     map[string]string{"Location": "sample"},
	}, nil
}

func mockCreateSessionRPCError(req sessionproto.SessionCreateRequest) (*sessionproto.SessionCreateResponse, error) {
	return nil, errors.New("RPC Error")
}

func mockDeleteSessionRPC(sessionID, sessionToken string) (*sessionproto.SessionResponse, error) {
	return &sessionproto.SessionResponse{
		StatusCode: http.StatusOK,
	}, nil
}

func mockSessionRPCError(sessionID, sessionToken string) (*sessionproto.SessionResponse, error) {
	return nil, errors.New("RPC Error")
}

func mockGetSessionRPC(sessionID, sessionToken string) (*sessionproto.SessionResponse, error) {
	return &sessionproto.SessionResponse{
		StatusCode: http.StatusOK,
		Header:     map[string]string{"Location": "sample"},
	}, nil
}

func mockGetAllActiveSessionsRPC(sessionID, sessionToken string) (*sessionproto.SessionResponse, error) {
	return &sessionproto.SessionResponse{
		StatusCode: http.StatusOK,
	}, nil
}

func mockGetSessionServiceRPC() (*sessionproto.SessionResponse, error) {
	return &sessionproto.SessionResponse{
		StatusCode: http.StatusOK,
	}, nil
}

func mockGetSessionServiceRPCError() (*sessionproto.SessionResponse, error) {
	return nil, errors.New("RPC Error")
}

func TestSessionRPCs_CreateSession(t *testing.T) {
	var s SessionRPCs
	s.CreateSessionRPC = mockCreateSessionRPC

	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1")
	redfishRoutes.Post("/SessionService/Sessions", s.CreateSession)
	e := httptest.New(t, mockApp)
	e.POST(
		"/redfish/v1/SessionService/Sessions",
	).WithJSON(map[string]string{"admin": "Password"}).Expect().Status(http.StatusCreated)
}

func TestSessionRPCs_CreateSessionRPCError(t *testing.T) {
	var s SessionRPCs
	s.CreateSessionRPC = mockCreateSessionRPCError

	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1")
	redfishRoutes.Post("/SessionService/Sessions", s.CreateSession)
	e := httptest.New(t, mockApp)
	e.POST(
		"/redfish/v1/SessionService/Sessions",
	).WithJSON(map[string]string{"admin": "Password"}).Expect().Status(http.StatusInternalServerError).Headers().Equal(header)
}

func TestSessionRPCs_DeleteSession(t *testing.T) {
	var s SessionRPCs
	s.DeleteSessionRPC = mockDeleteSessionRPC

	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1")
	redfishRoutes.Delete("/SessionService/Sessions/{id}", s.DeleteSession)
	e := httptest.New(t, mockApp)
	e.DELETE(
		"/redfish/v1/SessionService/Sessions/123",
	).WithHeader("X-Auth-Token", "token").Expect().Status(http.StatusOK).Headers().Equal(header)
	e.DELETE(
		"/redfish/v1/SessionService/Sessions/123",
	).Expect().Status(http.StatusUnauthorized)
}

func TestSessionRPCs_DeleteSessionRPCError(t *testing.T) {
	var s SessionRPCs
	s.DeleteSessionRPC = mockSessionRPCError

	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1")
	redfishRoutes.Delete("/SessionService/Sessions/{id}", s.DeleteSession)
	e := httptest.New(t, mockApp)
	e.DELETE(
		"/redfish/v1/SessionService/Sessions/123",
	).WithHeader("X-Auth-Token", "token").Expect().Status(http.StatusInternalServerError)
}

func TestSessionRPCs_GetSession(t *testing.T) {
	var s SessionRPCs
	s.GetSessionRPC = mockGetSessionRPC

	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1")
	redfishRoutes.Get("/SessionService/Sessions/{id}", s.GetSession)
	e := httptest.New(t, mockApp)
	e.GET(
		"/redfish/v1/SessionService/Sessions/123",
	).WithHeader("X-Auth-Token", "token").WithHeader("sessionID", "token").Expect().Status(http.StatusOK)
}

func TestSessionRPCs_GetSessionRPCError(t *testing.T) {
	var s SessionRPCs
	s.GetSessionRPC = mockSessionRPCError

	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1")
	redfishRoutes.Get("/SessionService/Sessions/{id}", s.GetSession)
	e := httptest.New(t, mockApp)
	e.GET(
		"/redfish/v1/SessionService/Sessions/123",
	).WithHeader("X-Auth-Token", "token").WithHeader("sessionID", "token").Expect().Status(http.StatusInternalServerError)
}

func TestSessionRPCs_GetAllAciveSessions(t *testing.T) {
	header["Allow"] = []string{"GET, POST"}
	defer delete(header, "Allow")
	var s SessionRPCs
	s.GetAllActiveSessionsRPC = mockGetAllActiveSessionsRPC

	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1")
	redfishRoutes.Get("/SessionService/Sessions", s.GetAllActiveSessions)
	e := httptest.New(t, mockApp)
	e.GET(
		"/redfish/v1/SessionService/Sessions",
	).WithHeader("X-Auth-Token", "token").Expect().Status(http.StatusOK).Headers().Equal(header)
}

func TestSessionRPCs_GetAllAciveSessionsRPCError(t *testing.T) {
	var s SessionRPCs
	s.GetAllActiveSessionsRPC = mockSessionRPCError

	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1")
	redfishRoutes.Get("/SessionService/Sessions", s.GetAllActiveSessions)
	e := httptest.New(t, mockApp)
	e.GET(
		"/redfish/v1/SessionService/Sessions",
	).WithHeader("X-Auth-Token", "token").Expect().Status(http.StatusInternalServerError)
}

func TestSessionRPCs_GetSessionService(t *testing.T) {
	header["Allow"] = []string{"GET"}
	defer delete(header, "Allow")
	var s SessionRPCs
	s.GetSessionServiceRPC = mockGetSessionServiceRPC

	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1")
	redfishRoutes.Get("/SessionService", s.GetSessionService)
	e := httptest.New(t, mockApp)
	e.GET(
		"/redfish/v1/SessionService",
	).Expect().Status(http.StatusOK).Headers().Equal(header)
}

func TestSessionRPCs_GetSessionServiceRPCError(t *testing.T) {
	var s SessionRPCs
	s.GetSessionServiceRPC = mockGetSessionServiceRPCError

	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1")
	redfishRoutes.Get("/SessionService", s.GetSessionService)
	e := httptest.New(t, mockApp)
	e.GET(
		"/redfish/v1/SessionService",
	).Expect().Status(http.StatusInternalServerError)
}
