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
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/plugin-lenovo/lpmodel"
	"github.com/ODIM-Project/ODIM/plugin-lenovo/lpresponse"
	"github.com/ODIM-Project/ODIM/plugin-lenovo/lputilities"
	iris "github.com/kataras/iris/v12"

	pluginConfig "github.com/ODIM-Project/ODIM/plugin-lenovo/config"
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
	resetType := request["ResetType"].(string)
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

	//Subscribe to Events
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
	body, err := ioutil.ReadAll(resp.Body)
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
