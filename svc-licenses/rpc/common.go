//(C) Copyright [2022] Hewlett Packard Enterprise Development LP
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

package rpc

import (
	"encoding/json"

	licenseproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/licenses"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-licenses/licenses"

	log "github.com/sirupsen/logrus"
)

// Licenses struct helps to register service
type Licenses struct {
	connector *licenses.ExternalInterface
}

// GetLicense intializes all the required connection
func GetLicense() *Licenses {
	return &Licenses{
		connector: licenses.GetExternalInterface(),
	}
}

func generateResponse(input interface{}) []byte {
	bytes, err := json.Marshal(input)
	if err != nil {
		log.Warn("Unable to unmarshall response object from util-libs " + err.Error())
	}
	return bytes
}

func fillProtoResponse(resp *licenseproto.GetLicenseResponse, data response.RPC) {
	resp.StatusCode = data.StatusCode
	resp.StatusMessage = data.StatusMessage
	resp.Body = generateResponse(data.Body)
	resp.Header = data.Header

}

func generateRPCResponse(rpcResp response.RPC, licenseResp *licenseproto.GetLicenseResponse) {
	bytes, _ := json.Marshal(rpcResp.Body)
	*licenseResp = licenseproto.GetLicenseResponse{
		StatusCode:    rpcResp.StatusCode,
		StatusMessage: rpcResp.StatusMessage,
		Header:        rpcResp.Header,
		Body:          bytes,
	}
}
