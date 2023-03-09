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

// Package fabrics ...
package fabrics

import (
	"context"
	"net/http"
	"strings"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	l "github.com/ODIM-Project/ODIM/lib-utilities/logs"
	fabricsproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/fabrics"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-fabrics/fabmodel"
)

// RemoveFabric holds the logic for deleting specfic fabric resource
func RemoveFabric(ctx context.Context, req *fabricsproto.AddFabricRequest) response.RPC {
	var resp response.RPC
	origin := req.OriginResource
	uuid := origin[strings.LastIndexByte(origin, '/')+1:]
	var err error
	fab := fabmodel.Fabric{
		FabricUUID: uuid,
	}
	err = fab.RemoveFabricData(ctx,uuid)
	if err != nil {
		l.LogWithFields(ctx).Error(err.Error())
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, err.Error(),
			[]interface{}{}, nil)
	}
	l.LogWithFields(ctx).Info("Fabric Removed ", uuid)
	resp.StatusCode = http.StatusOK
	resp.StatusMessage = response.Success
	return resp
}
