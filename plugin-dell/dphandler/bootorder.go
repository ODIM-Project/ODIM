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
	"net/http"
	"strings"

	l "github.com/ODIM-Project/ODIM/lib-utilities/logs"
	pluginConfig "github.com/ODIM-Project/ODIM/plugin-dell/config"
	"github.com/ODIM-Project/ODIM/plugin-dell/dpmodel"
	"github.com/ODIM-Project/ODIM/plugin-dell/dputilities"
	iris "github.com/kataras/iris/v12"
)

// SetDefaultBootOrder : sets the defult boot order
func SetDefaultBootOrder(ctx iris.Context) {
	ctxt := ctx.Request().Context()

	//Get token from Request
	token := ctx.GetHeader("X-Auth-Token")
	//Validating the token
	if token != "" {
		flag := TokenValidation(token)
		if !flag {
			l.LogWithFields(ctxt).Error("X-Auth-Token is either expired or invalid")
			ctx.StatusCode(http.StatusUnauthorized)
			ctx.WriteString("Invalid/Expired X-Auth-Token")
			return
		}
	}

	var deviceDetails dpmodel.Device
	uri := ctx.Request().RequestURI
	l.LogWithFields(ctxt).Debugf("incoming request received for setting boot order for URI %s method %s", uri, ctx.Request().Method)
	//replacing the request url with south bound translation URL
	for key, value := range pluginConfig.Data.URLTranslation.SouthBoundURL {
		uri = strings.Replace(uri, key, value, -1)
	}
	l.LogWithFields(ctxt).Debugf("the request URI for default boot setting has been replaced with southbound URI %s", uri)

	//Get device details from request
	err := ctx.ReadJSON(&deviceDetails)
	if err != nil {
		l.LogWithFields(ctxt).Error("error while trying to get device data from request " + err.Error())
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
		l.LogWithFields(ctxt).Error(errMsg)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.WriteString(errMsg)
		return
	}

	//Subscribe to Events
	resp, err := redfishClient.SetDefaultBootOrder(device, uri)
	if err != nil {
		errorMessage := "while trying to set default boot order, got: " + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		if resp == nil {
			ctx.StatusCode(http.StatusInternalServerError)
			ctx.WriteString(errorMessage)
			return
		}
	}
	defer resp.Body.Close()
	body, err := IoUtilReadAll(resp.Body)
	if err != nil {
		body = []byte("while trying to set default boot order, got: " + err.Error())
		l.LogWithFields(ctxt).Error(string(body))
	}

	ctx.StatusCode(resp.StatusCode)
	ctx.Write(body)
}
