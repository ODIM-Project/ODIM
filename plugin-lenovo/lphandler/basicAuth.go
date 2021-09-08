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
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"time"

	pluginConfig "github.com/ODIM-Project/ODIM/plugin-lenovo/config"
	"github.com/ODIM-Project/ODIM/plugin-lenovo/lpmodel"
	"github.com/ODIM-Project/ODIM/plugin-lenovo/lpresponse"
	"github.com/ODIM-Project/ODIM/plugin-lenovo/lputilities"
	iris "github.com/kataras/iris/v12"
)

//TokenValidation validates sent token with the list of plugin generated tokens
func TokenValidation(token string) bool {
	var flag bool
	flag = false
	for _, v := range tokenDetails {
		if token == v.Token {
			flag = true
			if time.Since(v.LastUsed).Minutes() > pluginConfig.Data.SessionTimeoutInMinutes {
				return flag
			}
		}
	}
	return flag
}

//Validate does Basic authentication with device and returns UUID of device in response
func Validate(ctx iris.Context) {
	//Get token from Request
	if ctx.GetHeader("X-Auth-Token") == "" && ctx.GetHeader("Authorization") == "" {
		ctx.StatusCode(http.StatusUnauthorized)
		log.Error("No valid authorization")
		ctx.WriteString("No valid authorization")
		return
	}
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
	//Get device details from request
	err := ctx.ReadJSON(&deviceDetails)
	if err != nil {
		log.Error("Unable to collect data from request: " + err.Error())
		ctx.StatusCode(http.StatusBadRequest)
		ctx.WriteString("Unable to collect data from request: " + err.Error())
		return
	}
	device := &lputilities.RedfishDevice{
		Host:     deviceDetails.Host,
		Username: deviceDetails.Username,
		Password: string(deviceDetails.Password),
	}

	redfishClient, err := lputilities.GetRedfishClient()
	if err != nil {
		log.Error(err.Error())
		lpresponse.SetErrorResponse(ctx, err.Error(), http.StatusInternalServerError)
		return
	}

	//Get ServiceRoot of the device
	err = redfishClient.GetRootService(device)
	if err != nil {
		log.Error(err.Error())
		lpresponse.SetErrorResponse(ctx, err.Error(), http.StatusBadRequest)
		return
	}
	//Doing Get on device using basic Authentication
	resp, err := redfishClient.BasicAuthWithDevice(device, device.RootNode.Systems.Oid)
	if err != nil {
		log.Error(err.Error())
		lpresponse.SetErrorResponse(ctx, err.Error(), http.StatusBadRequest)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err.Error())
		lpresponse.SetErrorResponse(ctx, err.Error(), http.StatusBadRequest)
		return
	}

	defer resp.Body.Close()
	if resp.StatusCode == http.StatusUnauthorized {
		ctx.StatusCode(resp.StatusCode)
		ctx.JSON(string(body))
		return
	}
	if resp.StatusCode >= 300 {
		log.Error("Could not retrieve ComputerSystems for " + device.Host + " :" + string(body))
	}

	response := lpresponse.Device{
		ServerIP:   device.Host,
		Username:   device.Username,
		DeviceUUID: device.RootNode.UUID,
	}
	ctx.StatusCode(resp.StatusCode)
	ctx.JSON(response)
}
