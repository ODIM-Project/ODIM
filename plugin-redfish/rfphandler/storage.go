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
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strings"

	pluginConfig "github.com/ODIM-Project/ODIM/plugin-redfish/config"
	"github.com/ODIM-Project/ODIM/plugin-redfish/rfpmodel"
	"github.com/ODIM-Project/ODIM/plugin-redfish/rfputilities"
	iris "github.com/kataras/iris/v12"
)

//CreateVolume function is used for creating a volume under storage
func CreateVolume(ctx iris.Context) {
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

	var deviceDetails rfpmodel.Device

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
		PostBody: deviceDetails.PostBody,
	}

	redfishClient, err := rfputilities.GetRedfishClient()
	if err != nil {
		errMsg := "While trying to create the redfish client, got:" + err.Error()
		log.Error(errMsg)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.WriteString(errMsg)
		return
	}
	resp, err := redfishClient.DeviceCall(device, uri, http.MethodPost)
	if err != nil {
		errorMessage := "While trying to create volume, got:" + err.Error()
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
		body = []byte("While trying to read response body, got: " + err.Error())
		log.Error(string(body))
	}
	ctx.StatusCode(resp.StatusCode)
	ctx.Write(body)
}

// DeleteVolume function is used for deleting a volume under storage
func DeleteVolume(ctx iris.Context) {
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

	var deviceDetails rfpmodel.Device

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
		PostBody: deviceDetails.PostBody,
	}

	redfishClient, err := rfputilities.GetRedfishClient()
	if err != nil {
		errMsg := "While trying to create the redfish client, got:" + err.Error()
		log.Error(errMsg)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.WriteString(errMsg)
		return
	}
	resp, err := redfishClient.DeviceCall(device, uri, http.MethodDelete)
	if err != nil {
		errorMessage := "While trying to delete volume, got:" + err.Error()
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
		body = []byte("While trying to read the response body, got: " + err.Error())
		log.Error(string(body))
	}
	ctx.StatusCode(resp.StatusCode)
	ctx.Write(body)
}
