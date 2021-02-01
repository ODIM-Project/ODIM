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

//Package rfphandler ..
package rfphandler

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	pluginConfig "github.com/ODIM-Project/ODIM/plugin-redfish/config"
	iris "github.com/kataras/iris/v12"
	"strings"
)

var (
	// In Channel
	In chan<- interface{}
	// Out Channel
	Out <-chan interface{}
)

// RedfishEvents receives the subscribed events from the south bound system
// Then it will send the received data and ip to publish method
func RedfishEvents(ctx iris.Context) {
	var req interface{}
	err := ctx.ReadJSON(&req)
	if err != nil {
		log.Error(err.Error())
		ctx.StatusCode(http.StatusBadRequest)
		ctx.WriteString("error: bad request")
		return
	}
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
	writeEventToJobQueue(event)
	ctx.StatusCode(http.StatusOK)

}

// writeEventToJobQueue will write events to worker pool
func writeEventToJobQueue(event common.Events) {
	var events []interface{}
	//events := make([]interface{}, 0)
	events = append(events, event)
	done := make(chan bool)
	go common.RunWriteWorkers(In, events, 1, done)
}
