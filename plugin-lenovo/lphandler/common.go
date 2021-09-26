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

// Package lphandler ...
package lphandler

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strings"

	pluginConfig "github.com/ODIM-Project/ODIM/plugin-lenovo/config"
	"github.com/ODIM-Project/ODIM/plugin-lenovo/lputilities"
)

// convertToSouthBoundURI searches the key in an array and return bool
func convertToSouthBoundURI(req string) string {
	req = strings.Replace(req, "/Actions/VirtualMedia.EjectMedia", "", -1)
	req = strings.Replace(req, "/Actions/VirtualMedia.InsertMedia", "", -1)
	return req
}

// queryDevice is for querying a Dell server
func queryDevice(uri string, device *lputilities.RedfishDevice, method string) (int, http.Header, []byte, error) {
	redfishClient, err := lputilities.GetRedfishClient()
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

//translateToSouthBoundURL replacing the request url with south bound translation URL
func translateToSouthBoundURL(uri string) string {
	for key, value := range pluginConfig.Data.URLTranslation.SouthBoundURL {
		uri = strings.Replace(uri, key, value, -1)
	}
	return uri
}
