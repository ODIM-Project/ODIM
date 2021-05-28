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
	teleproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/telemetry"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-telemetry/telemetry"
	log "github.com/sirupsen/logrus"
)

// Telemetry struct helps to register service
type Telemetry struct {
	connector *telemetry.ExternalInterface
}

// GetTele intializes all the required connection functions for the telemetry execution
func GetTele() *Telemetry {
	return &Telemetry{
		connector: telemetry.GetExternalInterface(),
	}
}

func generateResponse(input interface{}) []byte {
	bytes, err := json.Marshal(input)
	if err != nil {
		log.Warn("Unable to unmarshall response object from util-libs " + err.Error())
	}
	return bytes
}

func fillProtoResponse(resp *teleproto.TelemetryResponse, data response.RPC) {
	resp.StatusCode = data.StatusCode
	resp.StatusMessage = data.StatusMessage
	resp.Body = generateResponse(data.Body)
	resp.Header = data.Header

}

func generateRPCResponse(rpcResp response.RPC, teleResp *teleproto.TelemetryResponse) {
	bytes, _ := json.Marshal(rpcResp.Body)
	*teleResp = teleproto.TelemetryResponse{
		StatusCode:    rpcResp.StatusCode,
		StatusMessage: rpcResp.StatusMessage,
		Header:        rpcResp.Header,
		Body:          bytes,
	}
}

func generateTaskRespone(taskID, taskURI string, rpcResp *response.RPC) {
	commonResponse := response.Response{
		OdataType:    "#Task.v1_4_2.Task",
		ID:           taskID,
		Name:         "Task " + taskID,
		OdataContext: "/redfish/v1/$metadata#Task.Task",
		OdataID:      taskURI,
	}
	commonResponse.MessageArgs = []string{taskID}
	commonResponse.CreateGenericResponse(rpcResp.StatusMessage)
	rpcResp.Body = commonResponse
}
