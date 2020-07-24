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
	fabricsproto "github.com/bharath-b-hpe/odimra/lib-utilities/proto/fabrics"
	"github.com/bharath-b-hpe/odimra/lib-utilities/response"
	"net/http"
)

//DeleteFabricResource holds the logic for deleting specfic fabric resource
// It accepts url and contacts the configured CFM plugin
// and deletes the metioned fabric resoure such as Endpoints,Ports
func (f *Fabrics) DeleteFabricResource(req *fabricsproto.FabricRequest) response.RPC {
	var resp response.RPC
	var contactRequest pluginContactRequest
	req.Method = http.MethodDelete
	var err error
	contactRequest, resp, err = f.parseFabricsRequest(req)
	if err != nil {
		return resp
	}
	return f.parseFabricsResponse(contactRequest, req.URL)
}
