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
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	pluginConfig "github.com/ODIM-Project/ODIM/plugin-dell/config"
	"github.com/ODIM-Project/ODIM/plugin-dell/dpmodel"
	"github.com/ODIM-Project/ODIM/plugin-dell/dpresponse"
	"github.com/ODIM-Project/ODIM/plugin-dell/dputilities"
	iris "github.com/kataras/iris/v12"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strings"
)

//GetManagersCollection  Fetches details of the given resource from the device
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
	var deviceDetails dpmodel.Device
	// if any error come while getting the device then request will be for  plugins manager
	ctx.ReadJSON(&deviceDetails)
	if deviceDetails.Host == "" {
		var members = []dpresponse.Link{
			dpresponse.Link{
				Oid: "/ODIM/v1/Managers/" + pluginConfig.Data.RootServiceUUID,
			},
		}

		managers := dpresponse.ManagersCollection{
			OdataContext: "/ODIM/v1/$metadata#ManagerCollection.ManagerCollection",
			//Etag:         "W/\"AA6D42B0\"",
			OdataID:      uri,
			OdataType:    "#ManagerCollection.ManagerCollection",
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

//GetManagersInfo Fetches details of the given resource from the device
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
	var deviceDetails dpmodel.Device
	// if any error come while getting the device then request will be for  plugins manager
	ctx.ReadJSON(&deviceDetails)
	if deviceDetails.Host == "" {
		managers := dpresponse.Manager{
			OdataContext: "/ODIM/v1/$metadata#Manager.Manager",
			//Etag:            "W/\"AA6D42B0\"",
			OdataID:         uri,
			OdataType:       "#Manager.v1_3_3.Manager",
			Name:            pluginConfig.Data.PluginConf.ID,
			ManagerType:     "Service",
			ID:              pluginConfig.Data.RootServiceUUID,
			UUID:            pluginConfig.Data.RootServiceUUID,
			FirmwareVersion: pluginConfig.Data.FirmwareVersion,
			Status: &dpresponse.ManagerStatus{
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

func getInfoFromDevice(uri string, deviceDetails dpmodel.Device, ctx iris.Context) {
	//replacing the request url with south bound translation URL
	for key, value := range pluginConfig.Data.URLTranslation.SouthBoundURL {
		uri = strings.Replace(uri, key, value, -1)
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

	if resp.StatusCode == 401 {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.WriteString("Authtication with the device failed")
		return
	}
	if resp.StatusCode >= 300 {
		log.Warn("Could not retreive generic resource for " + device.Host + ": " + string(body))
	}
	respData := string(body)
	//replacing the resposne with north bound translation URL
	for key, value := range pluginConfig.Data.URLTranslation.NorthBoundURL {
		respData = strings.Replace(respData, key, value, -1)
	}
	ctx.StatusCode(resp.StatusCode)
	ctx.Write([]byte(respData))
}

//VirtualMediaActions performs insert and eject virtual media operations on the device based on the request
func VirtualMediaActions(ctx iris.Context) {
	uri := ctx.Request().RequestURI
	uri = replaceURI(uri)
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
		PostBody: deviceDetails.PostBody,
	}

	// creating a eject virtual media payload for dell plugin
	if strings.Contains(uri, "VirtualMedia.EjectMedia") {
		payload := map[string]interface{}{}
		device.PostBody, err = json.Marshal(payload)
		if err != nil {
			log.Error(err.Error())
			ctx.StatusCode(http.StatusInternalServerError)
			ctx.WriteString(err.Error())
			return
		}
	}

	statusCode, _, body, err := queryDevice(uri, device, http.MethodPost)
	if err != nil {
		errMsg := "while performing actions on virtual media, got: " + err.Error()
		log.Error(errMsg)
		ctx.StatusCode(statusCode)
		ctx.WriteString(errMsg)
		return
	}

	if statusCode == http.StatusNoContent {
		log.Info("VirtualMediaActions is successful for URI : " + uri)
		statusCode = http.StatusOK
		body, err = createVirtMediaActionResponse()
		if err != nil {
			errMsg := "while creating a response for virtual media actions" + err.Error()
			log.Error(errMsg)
			ctx.StatusCode(http.StatusInternalServerError)
			ctx.WriteString(errMsg)
			return
		}
	} else {
		errResponse := string(body)
		log.Errorf("VirtualMediaActions is failed for the URI %s, getting response %v ", uri, errResponse)
	}

	ctx.StatusCode(statusCode)
	ctx.Write(body)
}

// createVirtMediaActionResponse is used for creating a final response for virtual media actions success scenario
func createVirtMediaActionResponse() ([]byte, error) {
	resp := dpresponse.ErrorResopnse{
		Error: dpresponse.Error{
			Code:    response.Success,
			Message: "See @Message.ExtendedInfo for more information.",
			MessageExtendedInfo: []dpresponse.MsgExtendedInfo{
				dpresponse.MsgExtendedInfo{
					MessageID:   response.Success,
					Message:     "Successfully performed virtual media actions",
					MessageArgs: []string{},
				},
			},
		},
	}
	body, err := json.Marshal(resp)
	return body, err
}
