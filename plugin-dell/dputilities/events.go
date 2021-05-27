package dputilities

import (
	"encoding/json"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	pluginConfig "github.com/ODIM-Project/ODIM/plugin-dell/config"
	"strings"
)

var (
	In  chan<- interface{}
	Out <-chan interface{}
)

// manualEvents is used to generate an event based on the inputs provided
// It will send the received data and ip to publish method
func ManualEvents(req common.MessageData, hostAddress string) {
	request, _ := json.Marshal(req)
	reqData := string(request)
	//replacing the response with north bound translation URL
	for key, value := range pluginConfig.Data.URLTranslation.NorthBoundURL {
		reqData = strings.Replace(reqData, key, value, -1)
	}
	event := common.Events{
		IP:      hostAddress,
		Request: []byte(reqData),
	}
	// Call writeEventToJobQueue to write events to worker pool
	WriteEventToJobQueue(event)
}

// writeEventToJobQueue will write events to worker pool
func WriteEventToJobQueue(event common.Events) {
	var events []interface{}
	//events := make([]interface{}, 0)
	events = append(events, event)
	done := make(chan bool)
	go common.RunWriteWorkers(In, events, 1, done)
}
