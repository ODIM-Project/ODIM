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

	"github.com/ODIM-Project/ODIM/lib-rest-client/pmbhandle"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	aggregatorproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/aggregator"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/lib-utilities/services"
	"github.com/ODIM-Project/ODIM/svc-aggregation/agcommon"
	"github.com/ODIM-Project/ODIM/svc-aggregation/agmessagebus"
	"github.com/ODIM-Project/ODIM/svc-aggregation/agmodel"
	"github.com/ODIM-Project/ODIM/svc-aggregation/system"
)

// Aggregator struct helps to register service
type Aggregator struct {
	connector *system.ExternalInterface
}

// GetAggregator intializes all the required connection functions for the aggregation execution
func GetAggregator() *Aggregator {
	return &Aggregator{
		connector: &system.ExternalInterface{
			ContactClient:            pmbhandle.ContactPlugin,
			Auth:                     services.IsAuthorized,
			GetSessionUserName:       services.GetSessionUserName,
			CreateTask:               services.CreateTask,
			CreateChildTask:          services.CreateChildTask,
			UpdateTask:               system.UpdateTaskData,
			CreateSubcription:        system.CreateDefaultEventSubscription,
			PublishEvent:             system.PublishEvent,
			GetPluginStatus:          agcommon.GetPluginStatus,
			SubscribeToEMB:           services.SubscribeToEMB,
			EncryptPassword:          common.EncryptWithPublicKey,
			DecryptPassword:          common.DecryptWithPrivateKey,
			DeleteComputeSystem:      agmodel.DeleteComputeSystem,
			DeleteSystem:             agmodel.DeleteSystem,
			DeleteEventSubscription:  services.DeleteSubscription,
			EventNotification:        agmessagebus.Publish,
			GetAllKeysFromTable:      agmodel.GetAllKeysFromTable,
			GetConnectionMethod:      agmodel.GetConnectionMethod,
			UpdateConnectionMethod:   agmodel.UpdateConnectionMethod,
			GetPluginMgrAddr:         agmodel.GetPluginData,
			GetAggregationSourceInfo: agmodel.GetAggregationSourceInfo,
			GenericSave:              agmodel.GenericSave,
			CheckActiveRequest:       agmodel.CheckActiveRequest,
			DeleteActiveRequest:      agmodel.DeleteActiveRequest,
			GetAllMatchingDetails:    agmodel.GetAllMatchingDetails,
			CheckMetricRequest:       agmodel.CheckMetricRequest,
			DeleteMetricRequest:      agmodel.DeleteMetricRequest,
			GetResource:              agmodel.GetResource,
			Delete:                   agmodel.Delete,
		},
	}
}

func generateResponse(rpcResp response.RPC, aggResp *aggregatorproto.AggregatorResponse) {
	bytes, _ := json.Marshal(rpcResp.Body)
	*aggResp = aggregatorproto.AggregatorResponse{
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
