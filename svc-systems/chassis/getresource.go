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
	"log"
	"net/http"
	"strings"

	dmtf "github.com/ODIM-Project/ODIM/lib-dmtf/model"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	chassisproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/chassis"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-systems/scommon"
	"github.com/ODIM-Project/ODIM/svc-systems/smodel"
	"github.com/ODIM-Project/ODIM/svc-systems/sresponse"
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
func (p *PluginContact) GetChassisResource(req *chassisproto.GetChassisRequest) response.RPC {
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
		return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errorMessage, []interface{}{"Chassis", req.RequestParam}, nil)
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
		log.Printf("error getting system details : %v", gerr.Error())
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
			log.Println("Request Url", req.URL)
			var err error
			if data, err = scommon.GetResourceInfoFromDevice(getDeviceInfoRequest, true); err != nil {
				log.Printf("error while getting resource: %v", err)
				errorMsg := err.Error()
				return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errorMsg, []interface{}{tableName, req.URL}, nil)
			}
		} else {
			log.Printf("error while getting resource: %v", errorMessage)
			return common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		}
	}
	var resource map[string]interface{}
	json.Unmarshal([]byte(data), &resource)
	resp.Body = resource
	resp.StatusCode = http.StatusOK
	resp.StatusMessage = response.Success
	return resp

}

//GetChassisCollection is to fetch all the Systems uri's and retruns with created collection
// of Chassis data of odimra
func GetChassisCollection(req *chassisproto.GetChassisRequest) response.RPC {
	var resp response.RPC
	resp.Header = map[string]string{
		"Allow":             `"GET"`,
		"Cache-Control":     "no-cache",
		"Connection":        "keep-alive",
		"Content-type":      "application/json; charset=utf-8",
		"Transfer-Encoding": "chunked",
		"OData-Version":     "4.0",
	}
	chassisCollectionKeysArray, err := smodel.GetAllKeysFromTable("Chassis")
	if err != nil {
		log.Printf("error getting all keys of ChassisCollection table : %v", err)
		errorMessage := err.Error()
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
	}
	chassisCollection := sresponse.Collection{
		OdataContext: "/redfish/v1/$metadata#ChassisCollection.ChassisCollection",
		OdataID:      "/redfish/v1/Chassis/",
		OdataType:    "#ChassisCollection.ChassisCollection",
		Description:  "Computer System Chassis view",
		Name:         "Computer System Chassis",
	}
	var members = []dmtf.Link{}
	for _, key := range chassisCollectionKeysArray {
		members = append(members, dmtf.Link{Oid: key})
	}
	chassisCollection.Members = members
	chassisCollection.MembersCount = len(members)
	resp.Body = chassisCollection
	resp.StatusCode = http.StatusOK
	resp.StatusMessage = response.Success

	return resp
}

// GetChassisInfo is used to fetch resource data. The function is supposed to be used as part of RPC
// For getting chassis resource information,  parameters need to be passed Request .
// Request holds the  Uuid and Url ,
// Url will be parsed from that search key will created
// There will be two return values for the fuction. One is the RPC response, which contains the
// status code, status message, headers and body and the second value is error.
func GetChassisInfo(req *chassisproto.GetChassisRequest) response.RPC {
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
		return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errorMessage, []interface{}{"Chassis", req.RequestParam}, nil)
	}
	data, gerr := smodel.GetResource("Chassis", req.URL)
	if gerr != nil {
		log.Printf("error getting system details : %v", gerr.Error())
		errorMessage := gerr.Error()
		if errors.DBKeyNotFound == gerr.ErrNo() {
			resp.StatusCode = http.StatusNotFound
			resp.StatusMessage = errors.ResourceNotFound
			return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errorMessage, []interface{}{"Chassis", req.URL}, nil)
		}
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
	}
	var resource map[string]interface{}
	json.Unmarshal([]byte(data), &resource)
	resp.Body = resource
	resp.StatusCode = http.StatusOK
	resp.StatusMessage = response.Success

	return resp

}
