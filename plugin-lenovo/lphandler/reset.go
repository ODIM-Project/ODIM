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

//Package lphandler ...
package lphandler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/ODIM-Project/ODIM/lib-dmtf/model"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	pluginConfig "github.com/ODIM-Project/ODIM/plugin-lenovo/config"
	"github.com/ODIM-Project/ODIM/plugin-lenovo/lpmodel"
	"github.com/ODIM-Project/ODIM/plugin-lenovo/lpresponse"
	"github.com/ODIM-Project/ODIM/plugin-lenovo/lputilities"
	iris "github.com/kataras/iris/v12"
	log "github.com/sirupsen/logrus"
)

//ResetComputerSystem : reset computer system
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
	var deviceDetails lpmodel.Device
	uri := ctx.Request().RequestURI
	//replacing the request url with south bound translation URL
	for key, value := range pluginConfig.Data.URLTranslation.SouthBoundURL {
		uri = strings.Replace(uri, key, value, -1)
	}
	//Get device details from request
	err := ctx.ReadJSON(&deviceDetails)
	if err != nil {
		errMsg := "Unable to collect data from request: " + err.Error()
		log.Error(errMsg)
		ctx.StatusCode(http.StatusBadRequest)
		ctx.WriteString(errMsg)
		return
	}
	device := &lputilities.RedfishDevice{
		Host:     deviceDetails.Host,
		Username: deviceDetails.Username,
		Password: string(deviceDetails.Password),
	}

	var request map[string]interface{}
	err = json.Unmarshal(deviceDetails.PostBody, &request)
	if err != nil {
		errMsg := "Unable to unmarshal request body in reset operation: " + err.Error()
		log.Error(errMsg)
		ctx.StatusCode(http.StatusBadRequest)
		ctx.WriteString(errMsg)
		return
	}
	if err != nil {
		log.Error("While trying to unmarshall data : " + err.Error())
		ctx.StatusCode(http.StatusBadRequest)
		ctx.WriteString("Error: bad request.")
		return
	}
	resetType := request["ResetType"].(string)
	systemURI := strings.Split(uri, "Actions")[0]
	statusCode, _, body, err := queryDevice(systemURI, device, http.MethodGet)
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
	device.PostBody, _ = json.Marshal(lpmodel.ResetPostRequest{
		ResetType: resetType,
	})
	redfishClient, err := lputilities.GetRedfishClient()
	if err != nil {
		errMsg := "While trying to create the redfish client, got:" + err.Error()
		log.Error(errMsg)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.WriteString(errMsg)
		return
	}

	resp, err := redfishClient.ResetComputerSystem(device, uri)
	if err != nil {
		errorMessage := "While trying to reset, got: " + err.Error()
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
		body = []byte("While trying to read response body, got: " + err.Error())
		log.Error(string(body))
	}

	// If the response code is http.StatusNoContent then converting it to http.StatusOK
	if resp.StatusCode == http.StatusNoContent {
		resp.StatusCode = http.StatusOK
		body, err = createResetActionResponse()
		if err != nil {
			errMsg := "while creating a response for computersystem reset action" + err.Error()
			log.Error(errMsg)
			ctx.StatusCode(http.StatusInternalServerError)
			ctx.WriteString(errMsg)
			return
		}
	}

	ctx.StatusCode(resp.StatusCode)
	ctx.Write(body)
}

// createResetActionResponse is used for creating a final response for reset action success scenario
func createResetActionResponse() ([]byte, error) {
	resp := lpresponse.ErrorResopnse{
		Error: lpresponse.Error{
			Code:    response.Success,
			Message: "See @Message.ExtendedInfo for more information.",
			MessageExtendedInfo: []lpresponse.MsgExtendedInfo{
				lpresponse.MsgExtendedInfo{
					MessageID:   response.Success,
					Message:     "Reset initiated successfully",
					MessageArgs: []string{},
				},
			},
		},
	}
	body, err := json.Marshal(resp)
	return body, err
}

func checkPowerState(resetType, powerState string) ([]byte, int32, error) {
	resp := response.Args{
		Code:    response.NoOperation,
		Message: "",
		ErrorArgs: []response.ErrArgs{
			response.ErrArgs{
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
