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

//Package dphandler ...
package dphandler

import (
	"encoding/json"
	"fmt"
	taskproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/task"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	pluginConfig "github.com/ODIM-Project/ODIM/plugin-dell/config"
	"github.com/ODIM-Project/ODIM/plugin-dell/dpmodel"
	"github.com/ODIM-Project/ODIM/plugin-dell/dputilities"
	iris "github.com/kataras/iris/v12"
	log "github.com/sirupsen/logrus"
)

//CreateVolume function is used for creating a volume under storage
func CreateVolume(ctx iris.Context) {
	// Get token from Request
	token := ctx.GetHeader("X-Auth-Token")
	uri := ctx.Request().RequestURI

	// Replacing the request url with south bound translation URL
	for key, value := range pluginConfig.Data.URLTranslation.SouthBoundURL {
		uri = strings.Replace(uri, key, value, -1)
	}

	// Validating the token
	if token != "" {
		flag := TokenValidation(token)
		if !flag {
			log.Error("Invalid/Expired X-Auth-Token")
			ctx.StatusCode(http.StatusUnauthorized)
			ctx.WriteString("Invalid/Expired X-Auth-Token")
			return
		}
	}

	// Create new task
	taskURI, err := dputilities.CreateTask()
	if err != nil {
		log.Errorf("Unable to create the task: %s", err.Error())
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.WriteString("Unable to create the task")
		return
	}

	var deviceDetails dpmodel.Device
	// Get device details from request
	err = ctx.ReadJSON(&deviceDetails)
	if err != nil {
		log.Error("While trying to collect data from request: " + err.Error())
		ctx.StatusCode(http.StatusBadRequest)
		ctx.WriteString("Error: bad request.")
		return
	}

	// Transforming the payload
	reqPostBody := string(deviceDetails.PostBody)
	var reqBody dpmodel.Volume
	err = json.Unmarshal(deviceDetails.PostBody, &reqBody)
	if err != nil {
		errMsg := "While unmarshalling the create volume request to the device, got: " + err.Error()
		log.Error(errMsg)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.WriteString(errMsg)
		return
	}

	systemID := ctx.Params().Get("id")
	storageInstance := ctx.Params().Get("id2")

	isInvalid, res := validateRequest(reqBody)
	if isInvalid {
		ctx.StatusCode(int(res.StatusCode))
		_, _ = ctx.JSON(res.Body)
		return
	}

	driveURI := reqBody.Drives[0].OdataID
	s := strings.Split(driveURI, "/")
	driveSystemID := s[4]
	reqPostBody = strings.Replace(reqPostBody, driveSystemID, systemID, -1)
	reqPostBody = convertToSouthBoundURI(reqPostBody, storageInstance)
	device := &dputilities.RedfishDevice{
		Host:     deviceDetails.Host,
		Username: deviceDetails.Username,
		Password: string(deviceDetails.Password),
		PostBody: []byte(reqPostBody),
	}

	// Getting the firmware version of server before creating a new volume
	managersURI := "/redfish/v1/Managers/" + systemID
	managersURI = strings.Replace(managersURI, "System", "iDRAC", -1)
	statusCode, verErrMsg := getFirmwareVersion(managersURI, device)
	if statusCode != http.StatusOK {
		log.Error(verErrMsg)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.WriteString(verErrMsg)
		return
	}

	taskID := retrieveTaskID(taskURI)
	go createVolume(device, taskID, uri, reqPostBody)

	ctx.Header("Location", "/taskmon/"+taskID)
	ctx.StatusCode(http.StatusAccepted)
}

func createVolume(device *dputilities.RedfishDevice, taskID, uri, requestBody string) {
	updateTask(taskID, device.Host, uri, requestBody, dputilities.Running, dputilities.Ok, 0, http.MethodPost)

	// Getting the list of volumes before creating a new volume
	volumeListBeforeCreate, err := getVolumeCollection(uri, device)
	if err != nil {
		updateTaskWithException(taskID, device.Host, uri, requestBody, http.MethodPost)
		return
	}

	// Send create Volume request to BMC
	statusCode, header, _, err := queryDevice(uri, device, http.MethodPost)
	if err != nil {
		log.Errorf("Error while creating volume. StatusCode: %d, msg: %s", statusCode, err.Error())
		updateTaskWithException(taskID, device.Host, uri, requestBody, http.MethodPost)
		return
	}

	taskURI := header.Get("Location")
	if taskURI == "" {
		log.Errorf("missing location in volume create response header. Unable to track task - create volume might or might not finish successfully")
		updateTask(taskID, device.Host, uri, requestBody, dputilities.Completed, dputilities.Warning, 100, http.MethodPost)
		return
	}
	// Wait for create Volume task to change its state
	err = waitForTaskToFinish(taskURI, device, taskID, uri, requestBody, http.MethodPost)
	if err != nil {
		return
	}

	// Getting the list of volumes after creating a new volume
	volumeListAfterCreate, err := getVolumeCollection(uri, device)
	if err != nil {
		updateTaskWithException(taskID, device.Host, uri, requestBody, http.MethodPost)
		return
	}

	log.Info("volume was created successfully.")
	updateTask(taskID, device.Host, uri, requestBody, dputilities.Completed, dputilities.Ok, 100, http.MethodPost)

	// Getting the origin of condition for event
	oriOfCondition := compareCollection(volumeListBeforeCreate, volumeListAfterCreate)
	event := createEvent("Volume created Event", "ResourceAdded",
		"Volume is created successfully", "ResourceEvent.1.0.3.ResourceCreated", oriOfCondition)
	dputilities.ManualEvents(event, device.Host)
}

func waitForTaskToFinish(taskURI string, device *dputilities.RedfishDevice, taskID string, uri string, requestBody string,
	httpMethod string) error {
	for {
		time.Sleep(5 * time.Second)
		statusCode, _, body, err := queryDevice(taskURI, device, http.MethodGet)
		if err != nil {
			log.Errorf("Error while retrieving volume task. StatusCode: %d, msg: %s", statusCode, err.Error())
			updateTaskWithException(taskID, device.Host, uri, requestBody, httpMethod)
			return err
		}

		volumeTask := new(dpmodel.Task)
		err = json.Unmarshal(body, &volumeTask)
		if err != nil {
			log.Errorf("error while trying to unmarshal response data: " + err.Error())
			updateTaskWithException(taskID, device.Host, uri, requestBody, httpMethod)
			return err
		}

		state, err := dputilities.GetTaskState(volumeTask.TaskState)
		if err != nil {
			log.Errorf("error while trying to get task state from task: " + err.Error())
			updateTaskWithException(taskID, device.Host, uri, requestBody, httpMethod)
			return err
		}

		switch state {
		case dputilities.New, dputilities.Starting, dputilities.Running:
			continue
		case dputilities.Completed:
			log.Infof("volume task is completed!")
			return nil
		default:
			errorMsg := fmt.Sprintf("volume task finished with state %s, status code: %d", state.String(), statusCode)
			log.Errorf(errorMsg)
			updateTaskWithException(taskID, device.Host, uri, requestBody, http.MethodPost)
			return fmt.Errorf(errorMsg)
		}
	}
}

// DeleteVolume function is used for deleting a volume under storage
func DeleteVolume(ctx iris.Context) {
	// Get token from Request
	token := ctx.GetHeader("X-Auth-Token")
	uri := ctx.Request().RequestURI
	// Replacing the request url with south bound translation URL
	for key, value := range pluginConfig.Data.URLTranslation.SouthBoundURL {
		uri = strings.Replace(uri, key, value, -1)
	}

	storageInstance := ctx.Params().Get("id2")
	uri = convertToSouthBoundURI(uri, storageInstance)

	// Validating the token
	if token != "" {
		flag := TokenValidation(token)
		if !flag {
			log.Error("Invalid/Expired X-Auth-Token")
			ctx.StatusCode(http.StatusUnauthorized)
			ctx.WriteString("Invalid/Expired X-Auth-Token")
			return
		}
	}

	var deviceDetails dpmodel.Device

	// Get device details from request
	err := ctx.ReadJSON(&deviceDetails)
	if err != nil {
		log.Error("While trying to collect data from request, got: " + err.Error())
		ctx.StatusCode(http.StatusBadRequest)
		ctx.WriteString("Error: bad request.")
		return
	}
	device := &dputilities.RedfishDevice{
		Host:     deviceDetails.Host,
		Username: deviceDetails.Username,
		Password: string(deviceDetails.Password),
		PostBody: deviceDetails.PostBody,
	}

	taskURI, err := dputilities.CreateTask()
	if err != nil {
		log.Errorf("Unable to create the task: %s", err.Error())
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.WriteString("Unable to create the task")
		return
	}

	taskID := retrieveTaskID(taskURI)
	go deleteVolume(device, taskID, uri)

	ctx.StatusCode(http.StatusAccepted)
	ctx.Header("Location", "/taskmon/"+taskID)
}

func deleteVolume(device *dputilities.RedfishDevice, taskID, uri string) {
	updateTask(taskID, device.Host, uri, "", dputilities.Running, dputilities.Ok, 0, http.MethodDelete)

	// Send delete Volume request to BMC
	statusCode, header, _, err := queryDevice(uri, device, http.MethodDelete)
	if err != nil {
		log.Errorf("Error while deleting volume. StatusCode: %d, msg: %s", statusCode, err.Error())
		updateTaskWithException(taskID, device.Host, uri, "", http.MethodPost)
		return
	}

	taskURI := header.Get("Location")
	if taskURI == "" {
		log.Errorf("missing location in volume delete response header. Unable to track task - delete volume might or might not finish successfully")
		updateTask(taskID, device.Host, uri, "", dputilities.Completed, dputilities.Warning, 100, http.MethodPost)
		return
	}

	// Wait for delete volume task to complete.
	err = waitForTaskToFinish(taskURI, device, taskID, uri, "", http.MethodDelete)
	if err != nil {
		return
	}

	log.Infof("volume was deleted successfully.")
	updateTask(taskID, device.Host, uri, "", dputilities.Completed, dputilities.Ok, 100, http.MethodDelete)

	event := createEvent("Volume removed event", "ResourceRemoved", "Volume is deleted successfully",
		"ResourceEvent.1.0.3.ResourceRemoved", uri)
	dputilities.ManualEvents(event, device.Host)
}

// getVolumeCollection lists all the available volumes in the device
func getVolumeCollection(uri string, device *dputilities.RedfishDevice) ([]string, error) {
	// Getting the list of volumes already exist in the server
	statusCode, _, resp, err := queryDevice(uri, device, http.MethodGet)
	if err != nil {
		log.Errorf("Error while fetching volume collection. StatusCode: %d, msg: %s", statusCode, err.Error())
		return nil, err
	}

	var volumes dpmodel.VolumesCollection
	err = json.Unmarshal(resp, &volumes)
	if err != nil {
		log.Errorf("error while trying to unmarshal response data in create volume: " + err.Error())
		return nil, err
	}

	var list []string
	for _, member := range volumes.Members {
		list = append(list, member.OdataID)
	}
	return list, nil
}

// compareCollection will compare 2 slices and return the unique item from list2
func compareCollection(list1, list2 []string) string {
	var result string
	if len(list1) == 0 && len(list2) > 0 {
		result = list2[0]
	} else {
	outer:
		for i := len(list2) - 1; i >= 0; i-- {
		inner:
			for _, item := range list1 {
				if list2[i] == item {
					break inner
				} else {
					result = list2[i]
					break outer
				}
			}
		}
	}
	return result
}

func updateTask(taskID, host, targetURI, request string, taskState dputilities.TaskState, taskStatus dputilities.TaskStatus,
	percentComplete int32, httpMethod string) {
	payLoad := &taskproto.Payload{
		HTTPOperation: httpMethod,
		JSONBody:      request,
		TargetURI:     targetURI,
	}

	err := dputilities.UpdateTask(taskID, host, taskState, taskStatus, percentComplete, payLoad, time.Now())
	if err != nil {
		log.Errorf("Unable to update task with ID: %s", taskID)
	}
}

func retrieveTaskID(taskURI string) string {
	strArray := strings.Split(taskURI, "/")
	if strings.HasSuffix(taskURI, "/") {
		return strArray[len(strArray)-2]
	} else {
		return strArray[len(strArray)-1]
	}
}

func updateTaskWithException(taskID, host, uri, requestBody, method string) {
	updateTask(taskID, host, uri, requestBody, dputilities.Exception, dputilities.Critical, 100, method)
}

func validateRequest(requestBody dpmodel.Volume) (bool, response.RPC) {
	if requestBody.VolumeType == "" {
		return true, common.GeneralError(http.StatusBadRequest, response.PropertyMissing, "", []interface{}{"VolumeType"}, nil)
	}

	if requestBody.Drives == nil || len(requestBody.Drives) == 0 {
		return true, common.GeneralError(http.StatusBadRequest, response.PropertyMissing, "", []interface{}{"Drives"}, nil)
	}

	if requestBody.Drives[0].OdataID == "" {
		return true, common.GeneralError(http.StatusBadRequest, response.PropertyMissing, "", []interface{}{"@odata.id"}, nil)
	}

	return false, response.RPC{}
}

func createEvent(name string, eventType string, message string, messageID string, origin string) common.MessageData {
	return common.MessageData{
		OdataType: "#Event.v1_2_1.Event",
		Name:      name,
		Context:   "/redfish/v1/$metadata#Event.Event",
		Events: []common.Event{
			{
				EventType:      eventType,
				Severity:       "Critical",
				EventTimestamp: time.Now().String(),
				Message:        message,
				MessageID:      messageID,
				OriginOfCondition: &common.Link{
					Oid: origin,
				},
			},
		},
	}
}

// getFirmwareVersion of the device
func getFirmwareVersion(uri string, device *dputilities.RedfishDevice) (int, string) {
	// Getting the firmware version of a server
	statusCode, _, resp, err := queryDevice(uri, device, http.MethodGet)
	if err != nil {
		errMsg := "While getting firmware version details during create volume, got: " + err.Error()
		log.Error(errMsg)
		return statusCode, errMsg
	}
	var firmware dpmodel.FirmwareVersion
	err = json.Unmarshal(resp, &firmware)
	if err != nil {
		errMsg := "While trying to unmarshal response data in create volume, got: " + err.Error()
		log.Error(errMsg)
		return http.StatusInternalServerError, errMsg
	}

	verSplit := strings.SplitN(firmware.FirmwareVersion, ".", 3)
	version := strings.Join(verSplit[:2], "")
	res, _ := strconv.Atoi(version)

	if res < 440 {
		return http.StatusBadRequest, "Unsupported Firmware version, Firmware version should be >= 4.40"
	}
	return http.StatusOK, ""
}
