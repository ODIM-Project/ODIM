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

// Package dphandler ...
package dphandler

import (
	"encoding/json"
	"fmt"
	taskproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/task"
	"github.com/ODIM-Project/ODIM/plugin-dell/dputilities"
	"github.com/google/uuid"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/ODIM-Project/ODIM/plugin-dell/config"
	"github.com/ODIM-Project/ODIM/plugin-dell/dpmodel"
	"github.com/ODIM-Project/ODIM/plugin-dell/dpresponse"
	iris "github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/httptest"
)

func mockDevice(username, password, url string, w http.ResponseWriter) {
	var volume = dpmodel.VolumesCollection{
		OdataContext: "/redfish/v1/$metadata#VolumeCollection.VolumeCollection",
		OdataID:      "/redfish/v1/Systems/1/Storage/ArrayControllers-0/Volumes",
		OdataType:    "#VolumeCollection.VolumeCollection",
		Description:  "Volume Collection view",
		Members: []dpmodel.OdataIDLink{
			{
				OdataID: "/redfish/v1/Systems/1/Storage/ArrayControllers-0/Volumes/1",
			},
		},
		MembersCount: 1,
		Name:         "Volume Collection",
	}

	firmware := dpmodel.FirmwareVersion{
		FirmwareVersion: "4.40.10.00",
	}

	firmwareOld := dpmodel.FirmwareVersion{
		FirmwareVersion: "4.39.10.00",
	}

	volumeTask := dpmodel.Task{
		TaskState: "Completed",
	}

	if url == "/redfish/v1/Managers/1" {
		e, _ := json.Marshal(firmware)
		w.WriteHeader(http.StatusOK)
		w.Write(e)
		return
	}

	if url == "/redfish/v1/Managers/2" {
		e, _ := json.Marshal(firmwareOld)
		w.WriteHeader(http.StatusOK)
		w.Write(e)
		return
	}

	if url == "/ODIM/v1/Systems/1/Storage/1/Volumes" && username == "admin" {
		e, _ := json.Marshal(volume)
		w.Header().Add("Location", "/taskmon/1")
		w.WriteHeader(http.StatusOK)
		w.Write(e)
		return
	}
	if url == "/ODIM/v1/Systems/1/Storage/1/Volumes" && username != "admin" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if url == "/ODIM/v1/Systems/1/Storage/1/Volumes/1" && username == "admin" {
		w.WriteHeader(http.StatusOK)
		return
	}
	if url == "/ODIM/v1/Systems/1/Storage/1/Volumes/1" && username != "admin" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if url == "/ODIM/v1/Systems/1/Storage/Volumes/1" && username == "admin" {
		w.Header().Add("Location", "/taskmon/1")
		w.WriteHeader(http.StatusOK)
		return
	}
	if url == "/taskmon/1" && username == "admin" {
		e, _ := json.Marshal(volumeTask)
		w.WriteHeader(http.StatusOK)
		w.Write(e)
		return
	}
	return
}

func TestCreateVolume(t *testing.T) {
	config.SetUpMockConfig(t)
	deviceHost := "localhost"
	devicePort := "1234"
	ts := startTestServer(mockDevice)
	// Start the server.
	ts.StartTLS()
	defer ts.Close()

	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1")

	redfishRoutes.Post("/Systems/{id}/Storage/{rid}/Volumes", CreateVolume)

	dpresponse.PluginToken = "token"

	e := httptest.New(t, mockApp)

	reqPostBody := map[string]interface{}{
		"Name":       "Volume_Test1",
		"RAIDType":   "RAID0",
		"VolumeType": "NonRedundant",
		"Drives":     []dpmodel.OdataIDLink{{OdataID: "/ODIM/v1/Systems/5a9e8356-265c-413b-80d2-58210592d931:1/Storage/ArrayControllers-0/Drives/0"}},
	}
	reqBodyBytes, _ := json.Marshal(reqPostBody)
	requestBody := map[string]interface{}{
		"ManagerAddress": fmt.Sprintf("%s:%s", deviceHost, devicePort),
		"UserName":       "admin",
		"Password":       []byte("P@$$w0rd"),
		"PostBody":       reqBodyBytes,
	}

	taskService = newTaskServiceMock()
	hardwareMock := newVolumeOnHardwareServiceMock()
	hardwareService = hardwareMock

	// Unit Test for success scenario
	hardwareMock.wg.Add(1)
	e.POST("/redfish/v1/Systems/1/Storage/1/Volumes").WithJSON(requestBody).Expect().Status(http.StatusAccepted)
	hardwareMock.wg.Wait()

	// Case for invalid token
	e.POST("/redfish/v1/Systems/1/Storage/1/Volumes").WithHeader("X-Auth-Token", "token").WithJSON(requestBody).Expect().Status(http.StatusUnauthorized)

	// Unit test for bad request scenario
	invalidRequestBody := "invalid"
	e.POST("/redfish/v1/Systems/1/Storage/1/Volumes").WithJSON(invalidRequestBody).Expect().Status(http.StatusBadRequest)

	// Unit test for firmware version less than 4.40
	reqPostBody = map[string]interface{}{
		"Name":       "Volume_Test2",
		"RAIDType":   "RAID0",
		"VolumeType": "NonRedundant",
		"Drives":     []dpmodel.OdataIDLink{{OdataID: "/ODIM/v1/Systems/5a9e8356-265c-413b-80d2-58210592d931:2/Storage/ArrayControllers-0/Drives/0"}},
	}
	reqBodyBytes, _ = json.Marshal(reqPostBody)
	requestBody = map[string]interface{}{
		"ManagerAddress": fmt.Sprintf("%s:%s", deviceHost, devicePort),
		"UserName":       "admin",
		"Password":       []byte("P@$$w0rd"),
		"PostBody":       reqBodyBytes,
	}
	// Unit Test for firmware version less than 4.40 scenario
	e.POST("/redfish/v1/Systems/2/Storage/1/Volumes").WithJSON(requestBody).Expect().Status(http.StatusBadRequest)
}

func TestDeleteVolume(t *testing.T) {
	config.SetUpMockConfig(t)
	deviceHost := "localhost"
	devicePort := "1234"
	ts := startTestServer(mockDevice)
	// Start the server.
	ts.StartTLS()
	defer ts.Close()

	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1")

	redfishRoutes.Delete("/Systems/{id}/Storage/{id2}/Volumes/{rid}", DeleteVolume)

	dpresponse.PluginToken = "token"

	e := httptest.New(t, mockApp)

	requestBody := map[string]interface{}{
		"ManagerAddress": fmt.Sprintf("%s:%s", deviceHost, devicePort),
		"UserName":       "admin",
		"Password":       []byte("P@$$w0rd"),
	}

	taskService = newTaskServiceMock()
	hardwareMock := newVolumeOnHardwareServiceMock()
	hardwareService = hardwareMock

	//Unit Test for success scenario
	hardwareMock.wg.Add(1)
	e.DELETE("/redfish/v1/Systems/1/Storage/1/Volumes/1").WithJSON(requestBody).Expect().Status(http.StatusAccepted)
	hardwareMock.wg.Wait()

	//Case for invalid token
	e.DELETE("/redfish/v1/Systems/1/Storage/1/Volumes/1").WithHeader("X-Auth-Token", "token").WithJSON(requestBody).Expect().Status(http.StatusUnauthorized)
}

type volumeOnHardwareServiceMock struct {
	wg sync.WaitGroup
}

func (hm *volumeOnHardwareServiceMock) createVolume(device *dputilities.RedfishDevice, taskID, uri, requestBody string) {
	defer hm.wg.Done()
	createVolume(device, taskID, uri, requestBody)
}

func (hm *volumeOnHardwareServiceMock) deleteVolume(device *dputilities.RedfishDevice, taskID, uri string) {
	defer hm.wg.Done()
	deleteVolume(device, taskID, uri)
}

func newVolumeOnHardwareServiceMock() *volumeOnHardwareServiceMock {
	return &volumeOnHardwareServiceMock{}
}

type taskServiceMock struct {
	dputilities.TaskService
}

func (ts *taskServiceMock) CreateTask() (string, error) {
	return "task" + uuid.New().String(), nil
}

func (ts *taskServiceMock) UpdateTask(taskID, host string, taskState dputilities.TaskState, taskStatus dputilities.TaskStatus,
	percentComplete int32, payLoad *taskproto.Payload, endTime time.Time) error {
	return nil
}

func (ts *taskServiceMock) GetTaskState(state string) (dputilities.TaskState, error) {
	return dputilities.Completed, nil
}

func (ts *taskServiceMock) GetTaskStatus(status string) (dputilities.TaskStatus, error) {
	return dputilities.Ok, nil
}

func newTaskServiceMock() *taskServiceMock {
	return &taskServiceMock{}
}
