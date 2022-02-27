//(C) Copyright [2019] Hewlett Packard Enterprise Development LP
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

// Package rfphandler ...
package rfphandler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	pluginConfig "github.com/ODIM-Project/ODIM/plugin-redfish/config"
	"github.com/ODIM-Project/ODIM/plugin-redfish/rfputilities"
	log "github.com/sirupsen/logrus"
)

//translateToSouthBoundURL replacing the request url with south bound translation URL
func translateToSouthBoundURL(uri string) string {
	for key, value := range pluginConfig.Data.URLTranslation.SouthBoundURL {
		uri = strings.Replace(uri, key, value, -1)
	}
	return uri
}

// queryDevice is for querying a Dell server
func queryDevice(uri string, device *rfputilities.RedfishDevice, method string) (int, http.Header, []byte, error) {
	redfishClient, err := rfputilities.GetRedfishClient()
	if err != nil {
		errMsg := "While trying to create the redfish client, got:" + err.Error()
		log.Error(errMsg)
		return http.StatusInternalServerError, nil, nil, fmt.Errorf(errMsg)
	}
	resp, err := redfishClient.DeviceCall(device, uri, method)
	if err != nil {
		log.Error(err.Error())
		if resp == nil {
			return http.StatusBadRequest, nil, nil, err
		}
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		errMsg := "While trying to read the response body, got: " + err.Error()
		log.Error(errMsg)
		return http.StatusInternalServerError, nil, nil, fmt.Errorf(errMsg)
	}
	return resp.StatusCode, resp.Header, body, nil
}
