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
	"context"
	"encoding/json"
	"log"
	"net/http"

	updateproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/update"
	"github.com/ODIM-Project/ODIM/svc-update/update"
)

// GetUpdateService is an rpc handler, it gets invoked during GET on UpdateService API (/redfis/v1/UpdateService/)
func (a *Updater) GetUpdateService(ctx context.Context, req *updateproto.UpdateRequest, resp *updateproto.UpdateResponse) error {
	response := update.GetUpdateService()
	body, err := json.Marshal(response.Body)
	if err != nil {
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMessage = "error while trying marshal the response body for get update service: " + err.Error()
		log.Printf(response.StatusMessage)
		return nil
	}
	resp.StatusCode = response.StatusCode
	resp.StatusMessage = response.StatusMessage
	resp.Header = response.Header
	resp.Body = body
	return nil
}
