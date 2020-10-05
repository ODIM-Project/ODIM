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
	"fmt"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	fabricsproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/fabrics"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-fabrics/fabmodel"
	"log"
	"net"
	"net/http"
	"strings"
)

//AddFabric holds the logic for Adding fabric
//It accepts post body and store the fabric details in DB
func AddFabric(req *fabricsproto.AddFabricRequest) response.RPC {
	var resp response.RPC
	origin := req.OriginResource
	address := req.Address
	uuid := origin[strings.LastIndexByte(origin, '/')+1:]

	pluginDetails, err := fabmodel.GetAllFabricPluginDetails()
	if err != nil {
		log.Printf(err.Error())
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, err.Error(),
			[]interface{}{}, nil)
	}
	var pluginID string
	for _, pluginkey := range pluginDetails {

		plugin, errs := fabmodel.GetPluginData(pluginkey)
		if errs != nil {
			log.Printf(errs.Error())
			return common.GeneralError(http.StatusInternalServerError, response.InternalError, errs.Error(),
				[]interface{}{}, nil)
		}

		// get the ip address from the host name
		addr, err := net.LookupIP(plugin.IP)
		if err != nil || len(addr) < 1 {
			errorMessage := "Can't lookup the ip from host name"
			if err != nil {
				errorMessage = "Can't lookup the ip from host name" + err.Error()
			}
			log.Printf(errorMessage)
			return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errs.Error(),
				[]interface{}{"IP Address", plugin.IP}, nil)
		}
		deviceIPAddress := fmt.Sprintf("%v", addr[0])

		if deviceIPAddress == address {
			pluginID = plugin.ID
			break
		}
	}
	if pluginID == "" {
		log.Printf("error: plugin ID is empty")
		return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, "error: no match found for plugin ID",
			[]interface{}{"IP Address", address}, nil)
	}
	fab := fabmodel.Fabric{
		FabricUUID: uuid,
		PluginID:   pluginID,
	}

	err = fab.AddFabricData(uuid)
	if err != nil {
		log.Printf(err.Error())
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, err.Error(),
			[]interface{}{}, nil)
	}
	resp.Header = map[string]string{
		"Allow":             `"GET"`,
		"Cache-Control":     "no-cache",
		"Connection":        "keep-alive",
		"Content-type":      "application/json; charset=utf-8",
		"Transfer-Encoding": "chunked",
		"OData-Version":     "4.0",
	}
	log.Println("Fabric Added")
	resp.StatusCode = http.StatusOK
	resp.StatusMessage = response.Success
	return resp
}
