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
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/ODIM-Project/ODIM/lib-dmtf/model"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	pluginConfig "github.com/ODIM-Project/ODIM/plugin-dell/config"
	"github.com/ODIM-Project/ODIM/plugin-dell/dpmodel"
	"github.com/ODIM-Project/ODIM/plugin-dell/dputilities"
	iris "github.com/kataras/iris/v12"
	log "github.com/sirupsen/logrus"
)

// QueryDevice alias queryDevice
var QueryDevice = queryDevice

// ResetComputerSystem : reset computer system
func ResetComputerSystem(ctx iris.Context) {

	//Get token from Request
	token := ctx.GetHeader("X-Auth-Token")
	//Validating the token
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
	uri := ctx.Request().RequestURI
	//replacing the request url with south bound translation URL
	for key, value := range pluginConfig.Data.URLTranslation.SouthBoundURL {
		uri = strings.Replace(uri, key, value, -1)
	}
	//Get device details from request
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
	}

	var request map[string]interface{}
	err = json.Unmarshal(deviceDetails.PostBody, &request)
	if err != nil {
		log.Error("While trying to unmarshal data : " + err.Error())
		ctx.StatusCode(http.StatusBadRequest)
		ctx.WriteString("Error: bad request.")
		return
	}
	resetType := request["ResetType"].(string)
	systemURI := strings.Split(uri, "Actions")[0]
	statusCode, _, body, err := QueryDevice(systemURI, device, http.MethodGet)
	if err != nil {
		errMsg := "error while getting system data, got: " + err.Error()
		log.Error(errMsg)
		ctx.StatusCode(int(statusCode))
		ctx.WriteString(errMsg)
		return
	}
	var respData model.ComputerSystem
	if err := json.Unmarshal(body, &respData); err != nil {
		log.Warn("While unmarshaling the bios settings response from device, got: " + err.Error())
		ctx.StatusCode(http.StatusBadRequest)
		ctx.WriteString("Error: bad request.")
		return
	}
	respBody, statuscode, err := checkPowerState(resetType, respData.PowerState)
	if err != nil {
		log.Error(err.Error())
		ctx.StatusCode(int(statuscode))
		ctx.Write(respBody)
		return
	}
	device.PostBody, _ = json.Marshal(dpmodel.ResetPostRequest{
		ResetType: resetType,
	})
	redfishClient, err := dputilities.GetRedfishClient()
	if err != nil {
		errMsg := "While trying to create the redfish client, got:" + err.Error()
		log.Error(errMsg)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.WriteString(errMsg)
		return
	}

	resp, err := redfishClient.ResetComputerSystem(device, uri)
	if err != nil {
		errorMessage := "While trying to reset the computer system, got: " + err.Error()
		log.Error(errorMessage)
		if resp == nil {
			ctx.WriteString(errorMessage)
			ctx.StatusCode(http.StatusInternalServerError)
			return
		}
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		errorMessage := "While trying to read the response body, got: " + err.Error()
		log.Error(errorMessage)
		ctx.WriteString(errorMessage)
	}
	if resp.StatusCode == http.StatusNoContent {
		resp.StatusCode = http.StatusOK
		body = updateResetResponse()
	}
	ctx.StatusCode(resp.StatusCode)
	ctx.Write(body)
}

func updateResetResponse() []byte {
	resp := response.Args{
		Code:      response.Success,
		Message:   "Request completed successfully",
		ErrorArgs: []response.ErrArgs{},
	}
	resp.CreateGenericErrorResponse()
	body, _ := json.Marshal(resp)
	return body
}

func checkPowerState(resetType, powerState string) ([]byte, int32, error) {
	resp := response.Args{
		Code:    response.NoOperation,
		Message: "",
		ErrorArgs: []response.ErrArgs{
			{
				StatusMessage: response.NoOperation,
			},
		},
	}
	r := resp.CreateGenericErrorResponse()
	body, _ := json.Marshal(r)

	switch powerState {
	case "On":
		if resetType == "On" {
			return body, http.StatusOK, fmt.Errorf("power state is on")
		}
	case "Off":
		if resetType == "ForceOff" {
			return body, http.StatusOK, fmt.Errorf("power state is Off")
		}
		if resetType != "On" {
			errorMessage := "Can't reset, power is in off state"
			resp := common.GeneralError(http.StatusConflict, response.PropertyValueConflict, errorMessage, []interface{}{resetType, powerState}, nil)
			body, _ := json.Marshal(resp.Body)
			return body, http.StatusConflict, fmt.Errorf("power state is Off")
		}
	}
	return body, 0, nil
}
