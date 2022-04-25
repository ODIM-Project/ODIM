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

//Package fabrics ...
package fabrics

import (
	"net/http"
	"strings"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	fabricsproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/fabrics"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-fabrics/fabmodel"
	log "github.com/sirupsen/logrus"
)

//DeleteFabricResource holds the logic for deleting specfic fabric resource
// It accepts url and contacts the configured CFM plugin
// and deletes the metioned fabric resoure such as Endpoints,Ports
func RemoveFabric(req *fabricsproto.AddFabricRequest) response.RPC {
	var resp response.RPC
	origin := req.OriginResource
	uuid := origin[strings.LastIndexByte(origin, '/')+1:]
	var err error
	fab := fabmodel.Fabric{
		FabricUUID: uuid,
	}
	err = fab.RemoveFabricData(uuid)
	if err != nil {
		log.Error(err.Error())
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, err.Error(),
			[]interface{}{}, nil)
	}
	log.Info("Fabric Removed ", uuid)
	resp.StatusCode = http.StatusOK
	resp.StatusMessage = response.Success
	return resp
}
