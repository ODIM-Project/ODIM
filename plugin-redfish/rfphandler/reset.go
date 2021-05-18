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

//Package rfphandler ...
package rfphandler

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/ODIM-Project/ODIM/plugin-redfish/rfpmodel"
	"github.com/ODIM-Project/ODIM/plugin-redfish/rfputilities"
	iris "github.com/kataras/iris/v12"

	pluginConfig "github.com/ODIM-Project/ODIM/plugin-redfish/config"
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
	var deviceDetails rfpmodel.Device
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
	device := &rfputilities.RedfishDevice{
		Host:     deviceDetails.Host,
		Username: deviceDetails.Username,
		Password: string(deviceDetails.Password),
	}

	var request map[string]interface{}
	err = json.Unmarshal(deviceDetails.PostBody, &request)
	resetType := request["ResetType"].(string)
	device.PostBody, _ = json.Marshal(rfpmodel.ResetPostRequest{
		ResetType: resetType,
	})
	redfishClient, err := rfputilities.GetRedfishClient()
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

	ctx.StatusCode(resp.StatusCode)
	ctx.Write(body)
}
