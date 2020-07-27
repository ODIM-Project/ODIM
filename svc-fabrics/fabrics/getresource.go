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
	fabricsproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/fabrics"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"net/http"
)

//GetFabricResource holds the logic for getting  specfic fabric resource
//It accepts url and contacts the configured CFM plugin
// and gets the metioned fabric resoure such as Endpoints,Ports
func (f *Fabrics) GetFabricResource(req *fabricsproto.FabricRequest) response.RPC {
	var resp response.RPC
	var contactRequest pluginContactRequest
	var err error
	req.Method = http.MethodGet
	contactRequest, resp, err = f.parseFabricsRequest(req)
	if err != nil {
		return resp
	}
	if req.URL == "/redfish/v1/Fabrics" {
		return resp
	}
	return f.parseFabricsResponse(contactRequest, req.URL)
}
