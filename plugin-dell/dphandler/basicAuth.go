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
	"io/ioutil"
	"net/http"
	"time"

	l "github.com/ODIM-Project/ODIM/lib-utilities/logs"
	pluginConfig "github.com/ODIM-Project/ODIM/plugin-dell/config"
	"github.com/ODIM-Project/ODIM/plugin-dell/dpmodel"
	"github.com/ODIM-Project/ODIM/plugin-dell/dpresponse"
	"github.com/ODIM-Project/ODIM/plugin-dell/dputilities"
	iris "github.com/kataras/iris/v12"
)

// TokenValidation validates sent token with the list of plugin generated tokens
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

var IoUtilReadAll = ioutil.ReadAll

// Validate does Basic authentication with device and returns UUID of device in response
func Validate(ctx iris.Context) {
	ctxt := ctx.Request().Context()
	//Get token from Request
	if ctx.GetHeader("X-Auth-Token") == "" && ctx.GetHeader("Authorization") == "" {
		ctx.StatusCode(http.StatusUnauthorized)
		l.LogWithFields(ctxt).Error("No valid authorization")
		ctx.WriteString("No valid authorization")
		return
	}
	token := ctx.GetHeader("X-Auth-Token")
	//Validating the token
	if token != "" {
		flag := TokenValidation(token)
		if !flag {
			l.LogWithFields(ctxt).Error("Invalid/Expired X-Auth-Token")
			ctx.StatusCode(http.StatusUnauthorized)
			ctx.WriteString("Invalid/Expired X-Auth-Token")
			return
		}
	}
	l.LogWithFields(ctxt).Debugf("incoming request received for the URI %s method %s", ctx.Request().RequestURI, ctx.Request().Method)
	var deviceDetails dpmodel.Device
	//Get device details from request
	err := ctx.ReadJSON(&deviceDetails)
	if err != nil {
		l.LogWithFields(ctxt).Error("While trying to collect data from request, got: " + err.Error())
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
		l.LogWithFields(ctxt).Error("error while trying to create the redfish client to connect to device " + err.Error())
		dpresponse.SetErrorResponse(ctx, err.Error(), http.StatusInternalServerError)
		return
	}
	//Get ServiceRoot of the device
	err = redfishClient.GetRootService(device)
	if err != nil {
		l.LogWithFields(ctxt).Error("error while trying to get RootService data of the device " + err.Error())
		dpresponse.SetErrorResponse(ctx, err.Error(), http.StatusBadRequest)
		return
	}
	//Doing Get on device using basic Authentication
	resp, err := redfishClient.BasicAuthWithDevice(ctxt, device, device.RootNode.Systems.Oid)
	if err != nil {
		l.LogWithFields(ctxt).Error("error while trying to authenticate device with basic auth " + err.Error())
		dpresponse.SetErrorResponse(ctx, err.Error(), http.StatusBadRequest)
		return
	}

	body, err := IoUtilReadAll(resp.Body)
	if err != nil {
		l.LogWithFields(ctxt).Error("error while trying to read the response body from the device" + err.Error())
		dpresponse.SetErrorResponse(ctx, err.Error(), http.StatusBadRequest)
		return
	}
	defer resp.Body.Close()
	l.LogWithFields(ctxt).Debugf("received response code %d for the URI %s with the response body %s from the southbound", resp.StatusCode, resp.Request.URL, string(body))

	if resp.StatusCode == http.StatusUnauthorized {
		l.LogWithFields(ctxt).Warn("received unauthorized error from the device while requesting for device data")
		ctx.StatusCode(resp.StatusCode)
		ctx.JSON(string(body))
		return
	}
	if resp.StatusCode >= 300 {
		l.LogWithFields(ctxt).Warn("Could not retrieve ComputerSystems for " + device.Host + ": " + string(body))
	}

	response := dpresponse.Device{
		ServerIP:   device.Host,
		Username:   device.Username,
		DeviceUUID: device.RootNode.UUID,
	}
	ctx.StatusCode(resp.StatusCode)
	ctx.JSON(response)
}
