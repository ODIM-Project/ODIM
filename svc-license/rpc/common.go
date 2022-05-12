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

package rpc

import (
	"encoding/json"

	licenseproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/license"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-license/license"
	log "github.com/sirupsen/logrus"
)

// License struct helps to register service
type License struct {
	connector *license.ExternalInterface
}

// GetTele intializes all the required connection functions for the telemetry execution
func GetLicense() *License {
	return &License{
		connector: license.GetExternalInterface(),
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
