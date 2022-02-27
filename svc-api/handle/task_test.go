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

	taskproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/task"
	iris "github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/httptest"
)

func mockGetTaskStatus(req *taskproto.GetTaskRequest) (*taskproto.TaskResponse, error) {
	var response = &taskproto.TaskResponse{}
	if req.TaskID == "1A" && req.SessionToken == "ValidToken" {
		response = &taskproto.TaskResponse{
			StatusCode:    200,
			StatusMessage: "Success",
			Body:          []byte(`{"Response":"Success"}`),
		}
	} else if req.TaskID == "1A" && req.SessionToken == "InvalidToken" {
		response = &taskproto.TaskResponse{
			StatusCode:    401,
			StatusMessage: "Unauthorized",
			Body:          []byte(`{"Response":"Unauthorized"}`),
		}
	} else if req.TaskID == "2A" && req.SessionToken == "ValidToken" {
		response = &taskproto.TaskResponse{
			StatusCode:    403,
			StatusMessage: "Forbidden",
			Body:          []byte(`{"Response":"Forbidden"}`),
		}
	} else if req.TaskID == "3A" {
		return &taskproto.TaskResponse{}, fmt.Errorf("RPC Error")

	}
	return response, nil
}
func mockTaskCollection(req *taskproto.GetTaskRequest) (*taskproto.TaskResponse, error) {
	var response = &taskproto.TaskResponse{}
	if req.SessionToken == "ValidToken" {
		response = &taskproto.TaskResponse{
			StatusCode:    200,
			StatusMessage: "Success",
			Body:          []byte(`{"Response":"Success"}`),
		}
	} else if req.SessionToken == "InvalidToken" {
		response = &taskproto.TaskResponse{
			StatusCode:    401,
			StatusMessage: "Unauthorized",
			Body:          []byte(`{"Response":"Unauthorized"}`),
		}
	} else if req.SessionToken == "token" {
		return &taskproto.TaskResponse{}, fmt.Errorf("RPC Error")
	}
	return response, nil
}
func mockGetTaskService(req *taskproto.GetTaskRequest) (*taskproto.TaskResponse, error) {
	var response = &taskproto.TaskResponse{}
	if req.SessionToken == "ValidToken" {
		response = &taskproto.TaskResponse{
			StatusCode:    200,
			StatusMessage: "Success",
			Body:          []byte(`{"Response":"Success"}`),
		}
	} else if req.SessionToken == "InvalidToken" {
		response = &taskproto.TaskResponse{
			StatusCode:    401,
			StatusMessage: "Unauthorized", Body: []byte(`{"Response":"Unauthorized"}`),
		}
	} else if req.SessionToken == "token" {
		return &taskproto.TaskResponse{}, fmt.Errorf("RPC Error")
	}
	return response, nil
}
func mockGetSubTasks(req *taskproto.GetTaskRequest) (*taskproto.TaskResponse, error) {
	var response = &taskproto.TaskResponse{}
	if req.SessionToken == "ValidToken" {
		response = &taskproto.TaskResponse{
			StatusCode:    200,
			StatusMessage: "Success",
			Body:          []byte(`{"Response":"Success"}`),
		}
	} else if req.SessionToken == "InvalidToken" {
		response = &taskproto.TaskResponse{
			StatusCode:    401,
			StatusMessage: "Unauthorized", Body: []byte(`{"Response":"Unauthorized"}`),
		}
	} else if req.SessionToken == "token" {
		return &taskproto.TaskResponse{}, fmt.Errorf("RPC Error")
	}
	return response, nil
}
func mockGetSubTask(req *taskproto.GetTaskRequest) (*taskproto.TaskResponse, error) {
	var response = &taskproto.TaskResponse{}
	if req.SessionToken == "ValidToken" {
		response = &taskproto.TaskResponse{
			StatusCode:    200,
			StatusMessage: "Success",
			Body:          []byte(`{"Response":"Success"}`),
		}
	} else if req.SessionToken == "InvalidToken" {
		response = &taskproto.TaskResponse{
			StatusCode:    401,
			StatusMessage: "Unauthorized", Body: []byte(`{"Response":"Unauthorized"}`),
		}
	} else if req.SessionToken == "token" {
		return &taskproto.TaskResponse{}, fmt.Errorf("RPC Error")
	}
	return response, nil
}
func mockGetTaskMonitor(req *taskproto.GetTaskRequest) (*taskproto.TaskResponse, error) {
	var response = &taskproto.TaskResponse{}
	if req.SessionToken == "ValidToken" {
		response = &taskproto.TaskResponse{
			StatusCode:    200,
			StatusMessage: "Success",
			Body:          []byte(`{"Response":"Success"}`),
		}
	} else if req.SessionToken == "" {
		response = &taskproto.TaskResponse{
			StatusCode:    401,
			StatusMessage: "Unauthorized", Body: []byte(`{"Response":"Unauthorized"}`),
		}
	} else if req.SessionToken == "token" {
		return &taskproto.TaskResponse{}, fmt.Errorf("RPC Error")
	}
	return response, nil
}

func TestGetTaskStatus_ValidTaskID(t *testing.T) {
	header["Allow"] = []string{"GET, DELETE"}
	defer delete(header, "Allow")
	var task TaskRPCs
	task.GetTaskRPC = mockGetTaskStatus
	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1/TaskService/Tasks")
	redfishRoutes.Get("/{TaskID}", task.GetTaskStatus)
	test := httptest.New(t, mockApp)
	test.GET(
		"/redfish/v1/TaskService/Tasks/1A",
	).WithHeader("X-Auth-Token", "ValidToken").Expect().Status(http.StatusOK)
}

func TestGetTaskStatus_InvalidTaskID(t *testing.T) {
	var task TaskRPCs
	task.GetTaskRPC = mockGetTaskStatus
	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1/TaskService/Tasks")
	redfishRoutes.Get("/{TaskID}", task.GetTaskStatus)
	test := httptest.New(t, mockApp)
	test.GET(
		"/redfish/v1/TaskService/Tasks/2A",
	).WithHeader("X-Auth-Token", "ValidToken").Expect().Status(http.StatusForbidden)
	test.GET(
		"/redfish/v1/TaskService/Tasks/3A",
	).WithHeader("X-Auth-Token", "ValidToken").Expect().Status(http.StatusInternalServerError).Headers().Equal(header)
}

func TestGetTaskStatus_InvalidToken(t *testing.T) {
	var task TaskRPCs
	task.GetTaskRPC = mockGetTaskStatus
	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1/TaskService/Tasks")
	redfishRoutes.Get("/{TaskID}", task.GetTaskStatus)
	test := httptest.New(t, mockApp)
	test.GET(
		"/redfish/v1/TaskService/Tasks/1A",
	).WithHeader("X-Auth-Token", "InvalidToken").Expect().Status(http.StatusUnauthorized)
}

func TestTaskCollection(t *testing.T) {
	var task TaskRPCs
	task.TaskCollectionRPC = mockTaskCollection
	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1/TaskService")
	redfishRoutes.Get("/Tasks", task.TaskCollection)
	test := httptest.New(t, mockApp)
	test.GET(
		"/redfish/v1/TaskService/Tasks",
	).WithHeader("X-Auth-Token", "ValidToken").Expect().Status(http.StatusOK)
	test.GET(
		"/redfish/v1/TaskService/Tasks",
	).WithHeader("X-Auth-Token", "").Expect().Status(http.StatusUnauthorized)
	test.GET(
		"/redfish/v1/TaskService/Tasks",
	).WithHeader("X-Auth-Token", "token").Expect().Status(http.StatusInternalServerError)
}

func TestGetTaskService_ValidToken(t *testing.T) {
	var task TaskRPCs
	task.GetTaskServiceRPC = mockGetTaskService
	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1/TaskService")
	redfishRoutes.Get("/", task.GetTaskService)
	test := httptest.New(t, mockApp)
	test.GET(
		"/redfish/v1/TaskService",
	).WithHeader("X-Auth-Token", "ValidToken").Expect().Status(http.StatusOK)
	test.GET(
		"/redfish/v1/TaskService",
	).WithHeader("X-Auth-Token", "token").Expect().Status(http.StatusInternalServerError)
}
func TestGetSubTasks(t *testing.T) {
	var task TaskRPCs
	task.GetSubTasksRPC = mockGetSubTasks
	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1/TaskService")
	redfishRoutes.Get("/Tasks/{id}/SubTasks", task.GetSubTasks)
	test := httptest.New(t, mockApp)
	test.GET(
		"/redfish/v1/TaskService/Tasks/1A/SubTasks",
	).WithHeader("X-Auth-Token", "ValidToken").Expect().Status(http.StatusOK)
	test.GET(
		"/redfish/v1/TaskService/Tasks/1A/SubTasks",
	).WithHeader("X-Auth-Token", "token").Expect().Status(http.StatusInternalServerError)
}
func TestGetSubTask(t *testing.T) {
	var task TaskRPCs
	task.GetSubTaskRPC = mockGetSubTask
	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1/TaskService")
	redfishRoutes.Get("/Tasks/{id}/SubTasks/{tid}", task.GetSubTask)
	test := httptest.New(t, mockApp)
	test.GET(
		"/redfish/v1/TaskService/Tasks/1A/SubTasks/1A",
	).WithHeader("X-Auth-Token", "ValidToken").Expect().Status(http.StatusOK)
	test.GET(
		"/redfish/v1/TaskService/Tasks/1A/SubTasks/1A",
	).WithHeader("X-Auth-Token", "token").Expect().Status(http.StatusInternalServerError)
}
func TestGetTaskMonitor(t *testing.T) {
	var task TaskRPCs
	task.GetTaskMonitorRPC = mockGetTaskMonitor
	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1/taskmon")
	redfishRoutes.Get("/1A", task.GetTaskMonitor)
	test := httptest.New(t, mockApp)
	test.GET(
		"/redfish/v1/taskmon/1A",
	).WithHeader("X-Auth-Token", "ValidToken").Expect().Status(http.StatusOK)
	test.GET(
		"/redfish/v1/taskmon/1A",
	).WithHeader("X-Auth-Token", "").Expect().Status(http.StatusUnauthorized)
	test.GET(
		"/redfish/v1/taskmon/1A",
	).WithHeader("X-Auth-Token", "token").Expect().Status(http.StatusInternalServerError)
}

func TestDeleteTask_ValidToken(t *testing.T) {
	var task TaskRPCs
	task.DeleteTaskRPC = mockGetTaskStatus
	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1/TaskService/Tasks")
	redfishRoutes.Delete("/{TaskID}", task.DeleteTask)
	test := httptest.New(t, mockApp)
	test.DELETE(
		"/redfish/v1/TaskService/Tasks/1A",
	).WithHeader("X-Auth-Token", "ValidToken").Expect().Status(http.StatusOK).Headers().Equal(header)
	test.DELETE(
		"/redfish/v1/TaskService/Tasks/3A",
	).WithHeader("X-Auth-Token", "ValidToken").Expect().Status(http.StatusInternalServerError).Headers().Equal(header)
}
