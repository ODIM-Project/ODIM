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
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	updateproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/update"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-update/update"
	log "github.com/sirupsen/logrus"
)

// Updater struct helps to register service
type Updater struct {
	connector *update.ExternalInterface
}

// GetUpdater intializes all the required connection functions for the updater execution
func GetUpdater() *Updater {
	return &Updater{
		connector: update.GetExternalInterface(),
	}
}

func generateResponse(input interface{}) []byte {
	bytes, err := json.Marshal(input)
	if err != nil {
		log.Warn("Unable to unmarshall response object from util-libs " + err.Error())
	}
	return bytes
}

func fillProtoResponse(resp *updateproto.UpdateResponse, data response.RPC) {
	resp.StatusCode = data.StatusCode
	resp.StatusMessage = data.StatusMessage
	resp.Body = generateResponse(data.Body)
	resp.Header = data.Header

}

func generateRPCResponse(rpcResp response.RPC, aggResp *updateproto.UpdateResponse) {
	bytes, _ := json.Marshal(rpcResp.Body)
	*aggResp = updateproto.UpdateResponse{
		StatusCode:    rpcResp.StatusCode,
		StatusMessage: rpcResp.StatusMessage,
		Header:        rpcResp.Header,
		Body:          bytes,
	}
}

func generateTaskRespone(taskID, taskURI string, rpcResp *response.RPC) {
	commonResponse := response.Response{
		OdataType:    common.TaskType,
		ID:           taskID,
		Name:         "Task " + taskID,
		OdataContext: "/redfish/v1/$metadata#Task.Task",
		OdataID:      taskURI,
	}
	commonResponse.MessageArgs = []string{taskID}
	commonResponse.CreateGenericResponse(rpcResp.StatusMessage)
	rpcResp.Body = commonResponse
}
