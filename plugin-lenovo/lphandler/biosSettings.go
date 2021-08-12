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

	pluginConfig "github.com/ODIM-Project/ODIM/plugin-lenovo/config"
	"github.com/ODIM-Project/ODIM/plugin-lenovo/lpmodel"
	"github.com/ODIM-Project/ODIM/plugin-lenovo/lputilities"
	iris "github.com/kataras/iris/v12"
	"strings"
)

// DeviceClient struct to call device for the operation
type DeviceClient struct {
	ChangeSettingsOnDevice func(device *lputilities.RedfishDevice, url, method string) (*http.Response, error)
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

	var deviceDetails lpmodel.Device

	//Get device details from request
	err := ctx.ReadJSON(&deviceDetails)
	if err != nil {
		errorMessage := "Unable to collect data from request: " + err.Error()
		log.Error(errorMessage)
		ctx.StatusCode(http.StatusBadRequest)
		ctx.WriteString(errorMessage)
		return
	}
	device := &lputilities.RedfishDevice{
		Host:     deviceDetails.Host,
		Username: deviceDetails.Username,
		Password: string(deviceDetails.Password),
		PostBody: deviceDetails.PostBody,
	}

	redfishClient, err := lputilities.GetRedfishClient()
	if err != nil {
		errMsg := "While trying to create the redfish client, got:" + err.Error()
		log.Error(errMsg)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.WriteString(errMsg)
		return
	}

	if strings.HasSuffix(uri, "/Bios/Settings") {
		uri = strings.Replace(uri, "/Bios/Settings", "/Bios/Pending", -1)
	}

	resp, err := redfishClient.DeviceCall(device, uri, http.MethodPatch)
	if err != nil {
		errorMessage := "While trying to change bios settings, got: " + err.Error()
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
		body = []byte("While trying to change bios settings, got: " + err.Error())
		log.Error(string(body))
	}

	// Replace response body with the standard redfish specification
	respData := string(body)
	respData = strings.Replace(respData, "/Bios/Pending", "/Bios/Settings", -1)

	ctx.StatusCode(resp.StatusCode)
	ctx.Write([]byte(respData))
}
