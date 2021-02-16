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

package chassis

import (
	"net/http"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-systems/smodel"
	log "github.com/sirupsen/logrus"
)

func (f *fabricFactory) getFabricResource(rID string) response.RPC {
	var resp response.RPC
	ch := make(chan response.RPC)
	
	managers, err := f.getFabricManagers()
	if err != nil {
		log.Warn("while trying to collect fabric managers details from DB, got " + err.Error())
		return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, "", []interface{}{"Chassis", rID}, nil)
	}

	for _, manager := range managers {
		go f.getResource(manager, rID, ch)
	}

	for i := 0; i < len(managers); i++ {
		resp = <-ch
		if is2xx(int(resp.StatusCode)) {
			return resp
		}
	}

	return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, "", []interface{}{"Chassis", rID}, nil)
}

func (f *fabricFactory) getResource(plugin smodel.Plugin, rID string, ch chan response.RPC) {
	ch <- common.GeneralError(http.StatusNotFound, response.ResourceNotFound, "", []interface{}{"Chassis", rID}, nil)
}
