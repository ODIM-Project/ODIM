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

//Package chassis ...
package chassis

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	chassisproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/chassis"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-systems/scommon"
	"github.com/ODIM-Project/ODIM/svc-systems/smodel"

	log "github.com/sirupsen/logrus"
)

//PluginContact struct to inject the pmb client function into the handlers
type PluginContact struct {
	ContactClient   func(string, string, string, string, interface{}, map[string]string) (*http.Response, error)
	DecryptPassword func([]byte) ([]byte, error)
	GetPluginStatus func(smodel.Plugin) bool
}

// GetChassisResource is used to fetch resource data. The function is supposed to be used as part of RPC
// For getting chassis resource information,  parameters need to be passed Request .
// Request holds the  Uuid and Url ,
// Url will be parsed from that search key will created
// There will be two return values for the fuction. One is the RPC response, which contains the
// status code, status message, headers and body and the second value is error.
func (p *PluginContact) GetChassisResource(req *chassisproto.GetChassisRequest) (response.RPC, error) {
	var resp response.RPC
	resp.Header = map[string]string{
		"Allow":             `"GET"`,
		"Cache-Control":     "no-cache",
		"Connection":        "keep-alive",
		"Content-type":      "application/json; charset=utf-8",
		"Transfer-Encoding": "chunked",
		"OData-Version":     "4.0",
	}

	requestData := strings.Split(req.RequestParam, ":")
	if len(requestData) <= 1 {
		errorMessage := "error: SystemUUID not found"
		return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errorMessage, []interface{}{"Chassis", req.RequestParam}, nil), nil
	}
	uuid := requestData[0]
	urlData := strings.Split(req.URL, "/")
	//generating serachUrl which will be a part of key and also used in formatting  response
	var tableName string
	if req.ResourceID == "" {
		resourceName := urlData[len(urlData)-1]
		tableName = common.ChassisResource[resourceName]
	} else {
		tableName = urlData[len(urlData)-2]
	}
	data, gerr := smodel.GetResource(tableName, req.URL)
	if gerr != nil {
		log.Error("error getting system details : " + gerr.Error())
		errorMessage := gerr.Error()
		if errors.DBKeyNotFound == gerr.ErrNo() {
			var getDeviceInfoRequest = scommon.ResourceInfoRequest{
				URL:             req.URL,
				UUID:            uuid,
				SystemID:        requestData[1],
				ContactClient:   p.ContactClient,
				DevicePassword:  p.DecryptPassword,
				GetPluginStatus: p.GetPluginStatus,
			}
			log.Info("Request Url" + req.URL)
			var err error
			if data, err = scommon.GetResourceInfoFromDevice(getDeviceInfoRequest, true); err != nil {
				log.Error("error while getting resource: " + err.Error())
				errorMsg := err.Error()
				return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errorMsg, []interface{}{tableName, req.URL}, nil), nil
			}
		} else {
			log.Error("error while getting resource: " + errorMessage)
			return common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil), nil
		}
	}
	var resource map[string]interface{}
	json.Unmarshal([]byte(data), &resource)
	resp.Body = resource
	resp.StatusCode = http.StatusOK
	resp.StatusMessage = response.Success
	return resp, nil

}
