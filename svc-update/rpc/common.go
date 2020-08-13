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
	updateproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/update"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/lib-utilities/services"
	"github.com/ODIM-Project/ODIM/svc-update/update"
)

// Updater struct helps to register service
type Updater struct {
	connector *update.ExternalInterface
}

func generateResponse(rpcResp response.RPC, uResp *updateproto.UpdateResponse) {
	bytes, _ := json.Marshal(rpcResp.Body)
	*uResp = updateproto.UpdateResponse{
		StatusCode:    rpcResp.StatusCode,
		StatusMessage: rpcResp.StatusMessage,
		Header:        rpcResp.Header,
		Body:          bytes,
	}
}

// GetUpdater intializes all the required connection functions for the updater execution
func GetUpdater() *Updater {
	return &Updater{
		connector: &update.ExternalInterface{
			Auth: services.IsAuthorized,
		},
	}
}
