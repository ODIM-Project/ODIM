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
	"strings"

	"encoding/json"
	dmtf "github.com/ODIM-Project/ODIM/lib-dmtf/model"
	pluginConfig "github.com/ODIM-Project/ODIM/plugin-lenovo/config"
	"github.com/ODIM-Project/ODIM/plugin-lenovo/lpmodel"
	"github.com/ODIM-Project/ODIM/plugin-lenovo/lputilities"
	iris "github.com/kataras/iris/v12"
)

//GetResource : Fetches details of the given resource from the device
func GetResource(ctx iris.Context) {
	//Get token from Request
	token := ctx.GetHeader("X-Auth-Token")
	uri := ctx.Request().RequestURI
	//replacing the reuest url with south bound translation URL
	for key, value := range pluginConfig.Data.URLTranslation.SouthBoundURL {
		uri = strings.Replace(uri, key, value, -1)
	} //Validating the token
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
	//plainText, err := lputilities.GetPlainText(device.Password)
	//device.Password = plainText
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
		errMsg := "While trying to read the response body, got: " + err.Error()
		log.Error(errMsg)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.WriteString(errMsg)
		return
	}

	if resp.StatusCode == http.StatusUnauthorized {
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.WriteString("Authentication with the device failed")
		return
	}
	if resp.StatusCode >= 300 {
		log.Error("Could not retrieve generic resource for" + device.Host + ": \n" + string(body) + ":\n" + uri)

	}
	respData := string(body)
	//replacing the resposne with north bound translation URL
	for key, value := range pluginConfig.Data.URLTranslation.NorthBoundURL {
		respData = strings.Replace(respData, key, value, -1)
	}

	respData = strings.Replace(respData, "/Bios/Pending", "/Bios/Settings", -1)

	//Adding actions links under virtual media get
	if strings.Contains(uri, "/VirtualMedia/") {
		vmActions := dmtf.VMActions{
			EjectMedia: dmtf.ActionTarget{
				Target: uri + "/Actions/VirtualMedia.EjectMedia",
			},
			InsertMedia: dmtf.ActionTarget{
				Target: uri + "/Actions/VirtualMedia.InsertMedia",
			},
		}
		var vm = dmtf.VirtualMedia{}
		err = json.Unmarshal([]byte(respData), &vm)
		if err != nil {
			errMsg := "while trying to unmarshal the response body" + err.Error()
			log.Error(errMsg)
			ctx.StatusCode(http.StatusInternalServerError)
			ctx.WriteString(errMsg)
			return
		}
		vm.Actions = vmActions

		out, err := json.Marshal(vm)
		if err != nil {
			errMsg := "while trying to marshal the response body" + err.Error()
			log.Error(errMsg)
			ctx.StatusCode(http.StatusInternalServerError)
			ctx.WriteString(errMsg)
			return
		}
		respData = string(out)
	}

	ctx.StatusCode(resp.StatusCode)
	ctx.Write([]byte(respData))
}
