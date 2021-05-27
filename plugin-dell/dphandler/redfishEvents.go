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

//Package dphandler ..
package dphandler

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	pluginConfig "github.com/ODIM-Project/ODIM/plugin-dell/config"
	"github.com/ODIM-Project/ODIM/plugin-dell/dputilities"
	iris "github.com/kataras/iris/v12"
)

// RedfishEvents receives the subscribed events from the south bound system
// Then it will send the received data and ip to publish method
func RedfishEvents(ctx iris.Context) {
	var req interface{}
	err := ctx.ReadJSON(&req)
	if err != nil {
		log.Println(err)
		ctx.StatusCode(http.StatusBadRequest)
		ctx.WriteString("error: bad request")
		return
	}
	log.Println("Event Request", req)
	remoteAddr := ctx.RemoteAddr()
	request, _ := json.Marshal(req)

	reqData := string(request)
	//replacing the resposne with north bound translation URL
	for key, value := range pluginConfig.Data.URLTranslation.NorthBoundURL {
		reqData = strings.Replace(reqData, key, value, -1)
	}
	event := common.Events{
		IP:      remoteAddr,
		Request: []byte(reqData),
	}

	// Call writeEventToJobQueue to write events to worker pool
	dputilities.WriteEventToJobQueue(event)
	ctx.StatusCode(http.StatusOK)

}
