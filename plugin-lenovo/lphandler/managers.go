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

// Package lphandler ...
package lphandler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	dmtf "github.com/ODIM-Project/ODIM/lib-dmtf/model"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	pluginConfig "github.com/ODIM-Project/ODIM/plugin-lenovo/config"
	"github.com/ODIM-Project/ODIM/plugin-lenovo/lpmodel"
	"github.com/ODIM-Project/ODIM/plugin-lenovo/lputilities"
	iris "github.com/kataras/iris/v12"
	log "github.com/sirupsen/logrus"
)

// GetManagersCollection  Fetches details of the given resource from the device
func GetManagersCollection(ctx iris.Context) {
	//Get token from Request
	token := ctx.GetHeader("X-Auth-Token")
	uri := ctx.Request().RequestURI
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
	// if any error come while getting the device then request will be for  plugins manager
	ctx.ReadJSON(&deviceDetails)
	if deviceDetails.Host == "" {
		var members = []*dmtf.Link{
			&dmtf.Link{
				Oid: "/ODIM/v1/Managers/" + pluginConfig.Data.RootServiceUUID,
			},
		}

		managers := dmtf.Collection{
			ODataContext: "/ODIM/v1/$metadata#ManagerCollection.ManagerCollection",
			//ODataEtag:         "W/\"AA6D42B0\"",
			ODataID:      uri,
			ODataType:    "#ManagerCollection.ManagerCollection",
			Description:  "Managers view",
			Name:         "Managers",
			Members:      members,
			MembersCount: len(members),
		}
		ctx.StatusCode(http.StatusOK)
		ctx.JSON(managers)
		return
	}
	getInfoFromDevice(uri, deviceDetails, ctx)
	return

}

// GetManagersInfo Fetches details of the given resource from the device
func GetManagersInfo(ctx iris.Context) {
	//Get token from Request
	token := ctx.GetHeader("X-Auth-Token")
	uri := ctx.Request().RequestURI

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
	// if any error come while getting the device then request will be for  plugins manager
	ctx.ReadJSON(&deviceDetails)
	if deviceDetails.Host == "" {
		managers := dmtf.Manager{
			ODataContext: "/ODIM/v1/$metadata#Manager.Manager",
			//Etag:            "W/\"AA6D42B0\"",
			ODataID:         uri,
			ODataType:       common.ManagerType,
			Name:            pluginConfig.Data.PluginConf.ID,
			ManagerType:     "Service",
			ID:              pluginConfig.Data.RootServiceUUID,
			UUID:            pluginConfig.Data.RootServiceUUID,
			FirmwareVersion: pluginConfig.Data.FirmwareVersion,
			Status: &dmtf.ManagerStatus{
				State: "Enabled",
			},
		}
		ctx.StatusCode(http.StatusOK)
		ctx.JSON(managers)
		return
	}
	getInfoFromDevice(uri, deviceDetails, ctx)
	return

}

func getInfoFromDevice(uri string, deviceDetails lpmodel.Device, ctx iris.Context) {
	//replacing the request url with south bound translation URL
	for key, value := range pluginConfig.Data.URLTranslation.SouthBoundURL {
		uri = strings.Replace(uri, key, value, -1)
	}
	device := &lputilities.RedfishDevice{
		Host:     deviceDetails.Host,
		Username: deviceDetails.Username,
		Password: string(deviceDetails.Password),
	}
	redfishClient, err := lputilities.GetRedfishClient()
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
		errMsg := "While trying to read the response body, got: " + err.Error()
		log.Error(errMsg)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.WriteString(errMsg)
		return
	}

	if resp.StatusCode == 401 {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.WriteString("Authentication with the device failed")
		return
	}
	if resp.StatusCode >= 300 {
		log.Error("Could not retrieve generic resource for " + device.Host + ": \n" + string(body) + ":\n" + uri)
	}
	respData := string(body)
	//replacing the resposne with north bound translation URL
	for key, value := range pluginConfig.Data.URLTranslation.NorthBoundURL {
		respData = strings.Replace(respData, key, value, -1)
	}
	ctx.StatusCode(resp.StatusCode)
	ctx.Write([]byte(respData))
}

// VirtualMediaActions performs insert and eject virtual media operations on the device based on the request
func VirtualMediaActions(ctx iris.Context) {
	uri := ctx.Request().RequestURI
	//replacing the request url with south bound translation URL
	for key, value := range pluginConfig.Data.URLTranslation.SouthBoundURL {
		uri = strings.Replace(uri, key, value, -1)
	}
	var deviceDetails lpmodel.Device
	//Get device details from request
	err := ctx.ReadJSON(&deviceDetails)
	if err != nil {
		log.Error("While trying to collect data from request, got: " + err.Error())
		ctx.StatusCode(http.StatusBadRequest)
		ctx.WriteString("Error: bad request.")
		return
	}
	device := &lputilities.RedfishDevice{
		Host:     deviceDetails.Host,
		Username: deviceDetails.Username,
		Password: string(deviceDetails.Password),
		PostBody: deviceDetails.PostBody,
	}

	// creating a eject virtual media payload for lenovo plugin
	if strings.HasSuffix(uri, "VirtualMedia.EjectMedia") {
		//Creating a payload for eject virtual media
		payload := lpmodel.VirtualMediaEject{}
		log.Info("Payload for Eject virtual media ", payload)
		device.PostBody, err = json.Marshal(payload)
		if err != nil {
			log.Error(err.Error())
			ctx.StatusCode(http.StatusInternalServerError)
			ctx.WriteString(err.Error())
			return
		}
	}
	uri = convertToSouthBoundURI(uri)
	statusCode, _, body, err := queryDevice(uri, device, http.MethodPatch)
	if err != nil {
		errMsg := "while performing actions on virtual media, got: " + err.Error()
		log.Error(errMsg)
		ctx.StatusCode(statusCode)
		ctx.WriteString(errMsg)
		return
	}

	ctx.StatusCode(statusCode)
	ctx.Write(body)
}
