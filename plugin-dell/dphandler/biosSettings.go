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
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	pluginConfig "github.com/ODIM-Project/ODIM/plugin-dell/config"
	"github.com/ODIM-Project/ODIM/plugin-dell/dpmodel"
	"github.com/ODIM-Project/ODIM/plugin-dell/dpresponse"
	"github.com/ODIM-Project/ODIM/plugin-dell/dputilities"
	iris "github.com/kataras/iris/v12"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strings"
)

// DeviceClient struct to call device for the operation
type DeviceClient struct {
	ChangeSettingsOnDevice func(device *dputilities.RedfishDevice, url, method string) (*http.Response, error)
	DecryptPassword        func(password []byte) ([]byte, error)
}

//ChangeSettings is generic function where we can do following operations on different call
// 1. change bios settings
// 2. change boot order settings
func ChangeSettings(ctx iris.Context) {
	//Get token from Request
	token := ctx.GetHeader("X-Auth-Token")
	uri := ctx.Request().RequestURI
	//replacing the request url with south bound translation URL
	for key, value := range pluginConfig.Data.URLTranslation.SouthBoundURL {
		uri = strings.Replace(uri, key, value, -1)
	}
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
		PostBody: deviceDetails.PostBody,
	}

	redfishClient, err := dputilities.GetRedfishClient()
	if err != nil {
		errMsg := "While trying to create the redfish client, got:" + err.Error()
		log.Error(errMsg)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.WriteString(errMsg)
		return
	}
	respBody, err := redfishClient.DeviceCall(device, uri, http.MethodPatch)
	if err != nil {
		errMsg := "While trying to change bios settings, got: " + err.Error()
		log.Error(errMsg)
		if respBody == nil {
			ctx.StatusCode(http.StatusInternalServerError)
			ctx.WriteString(errMsg)
			return
		}
	}
	var resp []byte
	statusCode := respBody.StatusCode
	var errorMessage string
	if strings.Contains(uri, "/Bios/Settings") && statusCode == http.StatusOK {
		statusCode, resp, errorMessage = changeBiosSettings(uri, device)
		if statusCode != http.StatusOK {
			ctx.StatusCode(statusCode)
			ctx.WriteString(errorMessage)
			return
		}
	} else {
		defer respBody.Body.Close()
		resp, err = ioutil.ReadAll(respBody.Body)
		if err != nil {
			errMsg := "While reading the response body, got" + err.Error()
			log.Error(errMsg)
			ctx.WriteString(errMsg)
		}
	}
	ctx.StatusCode(statusCode)
	ctx.Write(resp)
}

//changeBiosSettings contains the logic for changing the bios settings
func changeBiosSettings(uri string, device *dputilities.RedfishDevice) (int, []byte, string) {
	var errorMessage string
	statusCode, _, resp, err := queryDevice(uri, device, http.MethodGet)
	if err != nil {
		errorMessage = "While trying to retrieve bios settings details, got: " + err.Error()
		log.Error(errorMessage)
		return statusCode, nil, errorMessage
	}
	var biosSetting dpmodel.BiosSettings
	err = json.Unmarshal(resp, &biosSetting)
	if err != nil {
		errorMessage = "While trying to unmarshal bios settings data, got: " + err.Error()
		log.Error(errorMessage)
		return http.StatusInternalServerError, nil, errorMessage
	}

	jobsURI := biosSetting.Oem.Dell.Jobs.OdataID
	if jobsURI != "" {
		reqPostBody := map[string]interface{}{"TargetSettingsURI": uri}
		req, _ := json.Marshal(reqPostBody)
		device.PostBody = req
		statusCode, _, resp, err = queryDevice(jobsURI, device, http.MethodPost)
		if err != nil {
			errorMessage = "While trying to create a job for updating the Bios settings, got: " + err.Error()
			log.Error(errorMessage)
			return statusCode, nil, errorMessage
		}
		if statusCode == http.StatusOK {
			log.Info("Creation of job for bios settings is successful with body: " + string(resp))
			resp = createBiosResponse()
		} else {
			errorMessage = "Unable to create a job for applying bios settings"
			log.Error(errorMessage)
		}
	} else {
		errorMessage := "Unable to get the job URI from bios settings details"
		log.Error(errorMessage)
		return http.StatusInternalServerError, nil, errorMessage
	}
	return statusCode, resp, errorMessage
}

// createBiosResponse is used for creating a final response for bios settings
func createBiosResponse() []byte {
	resp := dpresponse.ErrorResopnse{
		Error: dpresponse.Error{
			Code:    response.Success,
			Message: "See @Message.ExtendedInfo for more information.",
			MessageExtendedInfo: []dpresponse.MsgExtendedInfo{
				dpresponse.MsgExtendedInfo{
					MessageID:   response.Success,
					Message:     "A system reset is required for BIOS settings changes to get affected",
					MessageArgs: []string{},
				},
			},
		},
	}
	body, _ := json.Marshal(resp)
	return body
}
