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
	pluginConfig "github.com/ODIM-Project/ODIM/plugin-dell/config"
	"github.com/ODIM-Project/ODIM/plugin-dell/dpmodel"
	"github.com/ODIM-Project/ODIM/plugin-dell/dputilities"
	iris "github.com/kataras/iris/v12"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strings"
)

//GetResource : Fetches details of the given resource from the device
func GetResource(ctx iris.Context) {
	//Get token from Request
	token := ctx.GetHeader("X-Auth-Token")
	uri := ctx.Request().RequestURI
	//replacing the request url with south bound translation URL
	for key, value := range pluginConfig.Data.URLTranslation.SouthBoundURL {
		uri = strings.Replace(uri, key, value, -1)
	}
	//Getting storageinstance from URI
	storageInstance := ctx.Params().Get("id2")

	uri = convertToSouthBoundURI(uri, storageInstance)
	// Transforming NetworkAdapters URI
	if strings.Contains(uri, "/Chassis/") && strings.Contains(uri, "NetworkAdapters") {
		uri = strings.Replace(uri, "/Chassis/", "/Systems/", -1)
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
	}

	redfishClient, err := dputilities.GetRedfishClient()
	if err != nil {
		errMsg := "While trying to create the redfish client, got:" + err.Error()
		log.Error(errMsg)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.WriteString(errMsg)
		return
	}

	//Fetching generic resource details from the device
	resp, err := redfishClient.GetWithBasicAuth(device, uri)
	if err != nil {
		errMsg := "Authentication failed: " + err.Error()
		log.Error(errMsg)
		if resp == nil {
			ctx.StatusCode(http.StatusInternalServerError)
			ctx.WriteString(errMsg)
			return
		}
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("While trying to read the response body, got: " + err.Error())
		return
	}

	if resp.StatusCode == http.StatusUnauthorized {
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.WriteString("Authentication with the device failed")
		return
	}
	if resp.StatusCode >= 300 {
		log.Warn("Could not retreive generic resource for " + device.Host + ": " + string(body))
	}
	respData := string(body)
	//replacing the response with north bound translation URL
	for key, value := range pluginConfig.Data.URLTranslation.NorthBoundURL {
		respData = strings.Replace(respData, key, value, -1)
	}
	respData = convertToNorthBoundURI(respData, storageInstance)

	//Transforming NetworkAdapters URI's
	if strings.Contains(uri, "/Chassis/") && strings.Contains(respData, "NetworkAdapters") {
		var respMap map[string]interface{}
		err := json.Unmarshal([]byte(respData), &respMap)
		if err != nil {
			errMsg := "While trying to unmarshal Chassis data, got: " + err.Error()
			log.Error(errMsg)
			ctx.StatusCode(http.StatusInternalServerError)
			ctx.WriteString(errMsg)
			return
		}
		netAdapterLink := respMap["NetworkAdapters"]
		var networkAdapterURI interface{}
		if netAdapterLink != nil {
			networkAdapterURI = netAdapterLink.(map[string]interface{})["@odata.id"]
		}
		netAdapterTransitionURI := strings.Replace(networkAdapterURI.(string), "Systems", "Chassis", -1)
		respData = strings.Replace(respData, networkAdapterURI.(string), netAdapterTransitionURI, -1)
	} else if strings.Contains(uri, "/Systems/") && strings.Contains(uri, "NetworkAdapters") {
		respData = strings.Replace(respData, "/Systems/", "/Chassis/", -1)
	}
	ctx.StatusCode(resp.StatusCode)
	ctx.Write([]byte(respData))
}
