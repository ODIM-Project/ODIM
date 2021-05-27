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
	"io/ioutil"
	"log"
	"net/http"

	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	pluginConfig "github.com/ODIM-Project/ODIM/plugin-dell/config"
	"github.com/ODIM-Project/ODIM/plugin-dell/dpmodel"
	"github.com/ODIM-Project/ODIM/plugin-dell/dpresponse"
	"github.com/ODIM-Project/ODIM/plugin-dell/dputilities"
	iris "github.com/kataras/iris/v12"
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
			log.Println("Invalid/Expired X-Auth-Token")
			ctx.StatusCode(http.StatusUnauthorized)
			ctx.WriteString("Invalid/Expired X-Auth-Token")
			return
		}
	}

	var deviceDetails dpmodel.Device

	//Get device details from request
	err := ctx.ReadJSON(&deviceDetails)
	if err != nil {
		log.Println("Error while trying to collect data from request: ", err)
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
		errMsg := "error: internal processing error: " + err.Error()
		log.Println(errMsg)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.WriteString(errMsg)
		return
	}
	respBody, err := redfishClient.DeviceCall(device, uri, http.MethodPatch)
	if err != nil {
		errorMessage := err.Error()
		fmt.Println(err)
		if respBody == nil {
			ctx.StatusCode(http.StatusInternalServerError)
			ctx.WriteString("error while trying to change bios settings: " + errorMessage)
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
			errorMessage := err.Error()
			log.Println(err)
			ctx.WriteString("Error while trying to change bios settings: " + errorMessage)
		}
	}

	log.Println("Response body: ", string(resp))
	ctx.StatusCode(statusCode)
	ctx.Write(resp)
}

//changeBiosSettings contains the logic for changing the bios settings
func changeBiosSettings(uri string, device *dputilities.RedfishDevice) (int, []byte, string) {
	var errorMessage string
	statusCode, _, resp, err := queryDevice(uri, device, http.MethodGet)
	if err != nil {
		errorMessage := "Error while trying to retrieve bios settings details: " + err.Error()
		log.Println(errorMessage)
		return statusCode, nil, errorMessage
	}
	var biosSetting dpmodel.BiosSettings
	err = json.Unmarshal(resp, &biosSetting)
	if err != nil {
		errorMessage := "error while trying to unmarshal bios settings data: " + err.Error()
		log.Println(errorMessage)
		return http.StatusInternalServerError, nil, errorMessage
	}

	jobsURI := biosSetting.Oem.Dell.Jobs.OdataID
	if jobsURI != "" {
		reqPostBody := map[string]interface{}{"TargetSettingsURI": uri}
		req, _ := json.Marshal(reqPostBody)
		device.PostBody = req
		statusCode, _, resp, err = queryDevice(jobsURI, device, http.MethodPost)
		if err != nil {
			errorMessage := "Error while trying to create a job for updating the Bios settings: " + err.Error()
			log.Println(errorMessage)
			return statusCode, nil, errorMessage
		}
		if statusCode == http.StatusOK {
			log.Println("Creation of job for bios settings is successful", string(resp))
			resp = createBiosResponse()
		} else {
			errorMessage = "error : Unable to create a job for applying bios settings"
			log.Println(errorMessage)
		}
	} else {
		errorMessage := "error : Unable to get the job URI from bios settings details"
		log.Println(errorMessage)
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
				{
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
