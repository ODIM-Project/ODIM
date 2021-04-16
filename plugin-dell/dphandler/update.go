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
	pluginConfig "github.com/ODIM-Project/ODIM/plugin-dell/config"
	"github.com/ODIM-Project/ODIM/plugin-dell/dpmodel"
	"github.com/ODIM-Project/ODIM/plugin-dell/dputilities"
	iris "github.com/kataras/iris/v12"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strings"
)

// SimpleUpdate updates the BMC resources
func SimpleUpdate(ctx iris.Context) {
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
	//Get device details from request
	err := ctx.ReadJSON(&deviceDetails)
	if err != nil {
		log.Error("While trying to collect data from request, got: " + err.Error())
		ctx.StatusCode(http.StatusBadRequest)
		ctx.WriteString("Error: bad request.")
		return
	}

	var reqData string
	//replacing the request url with south bound translation URL
	for key, value := range pluginConfig.Data.URLTranslation.SouthBoundURL {
		uri = strings.Replace(uri, key, value, -1)
		reqData = strings.Replace(string(deviceDetails.PostBody), key, value, -1)
	}

	device := &dputilities.RedfishDevice{
		Host:     deviceDetails.Host,
		Username: deviceDetails.Username,
		Password: string(deviceDetails.Password),
		PostBody: []byte(reqData),
	}

	redfishClient, err := dputilities.GetRedfishClient()
	if err != nil {
		errMsg := "While trying to create the redfish client, got:" + err.Error()
		log.Error(errMsg)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.WriteString(errMsg)
		return
	}
	//Update BMC resource
	resp, err := redfishClient.DeviceCall(device, uri, http.MethodPost)
	if err != nil {
		errorMessage := "While trying to update BMC resource: " + err.Error()
		log.Error(errorMessage)
		if resp == nil {
			ctx.StatusCode(http.StatusInternalServerError)
			ctx.WriteString(errorMessage)
			return
		}
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		errorMessage := "While trying to read the response body, got: " + err.Error()
		log.Error(errorMessage)
		ctx.WriteString(errorMessage)
	}

	ctx.StatusCode(resp.StatusCode)
	ctx.Write(body)
}

// StartUpdate updates the BMC resources
func StartUpdate(ctx iris.Context) {
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

	redfishClient, err := dputilities.GetRedfishClient()
	if err != nil {
		errMsg := "While trying to create the redfish client, got:" + err.Error()
		log.Error(errMsg)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.WriteString(errMsg)
		return
	}
	//Update BMC resource
	resp, err := redfishClient.DeviceCall(device, uri, http.MethodPost)
	if err != nil {
		errorMessage := "While trying to update BMC resource: " + err.Error()
		log.Error(errorMessage)
		if resp == nil {
			ctx.StatusCode(http.StatusInternalServerError)
			ctx.WriteString(errorMessage)
			return
		}
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		errorMessage := "While trying to read response body, got: " + err.Error()
		log.Error(errorMessage)
		ctx.WriteString(errorMessage)
	}

	ctx.StatusCode(resp.StatusCode)
	ctx.Write(body)
}
