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
package main

import (
	"log"
	"os"

	dc "github.com/ODIM-Project/ODIM/lib-messagebus/datacommunicator"
	"github.com/ODIM-Project/ODIM/lib-rest-client/pmbhandle"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	eventsproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/events"
	"github.com/ODIM-Project/ODIM/lib-utilities/services"
	"github.com/ODIM-Project/ODIM/svc-events/consumer"
	"github.com/ODIM-Project/ODIM/svc-events/evcommon"
	evt "github.com/ODIM-Project/ODIM/svc-events/events"
	"github.com/ODIM-Project/ODIM/svc-events/rpc"
)

func main() {
	// verifying the uid of the user
	if uid := os.Geteuid(); uid == 0 {
		log.Fatalln("Event Service should not be run as the root user")
	}

	if err := config.SetConfiguration(); err != nil {
		log.Fatalf("fatal: error while trying set up configuration: %v", err)
	}

	if err := dc.SetConfiguration(config.Data.MessageQueueConfigFilePath); err != nil {
		log.Fatalf("error while trying to set messagebus configuration: %v", err)
	}

	if err := common.CheckDBConnection(); err != nil {
		log.Fatalf("error while trying to check DB connection health: %v", err)
	}

	if err := services.InitializeService(services.Events); err != nil {
		log.Fatalf("fatal: error while trying to initialize the service: %v", err)
	}

	// Intializing the TopicsList
	evcommon.EMBTopics.TopicsList = make(map[string]bool)
	// Intializing plugin token
	evcommon.Token.Tokens = make(map[string]string)
	registerHandler()

	// CreateJobQueue defines the queue which will act as an infinite buffer
	// In channel is an entry or input channel and the Out channel is an exit or output channel
	consumer.In, consumer.Out = common.CreateJobQueue()

	// RunReadWorkers will create a worker pool for doing a specific task
	// which is passed to it as PublishEventsToDestination method after reading the data from the channel.
	common.RunReadWorkers(consumer.Out, evt.PublishEventsToDestination, 1)
	startUPInterface := evcommon.StartUpInteraface{
		DecryptPassword: common.DecryptWithPrivateKey,
		EMBConsume:      consumer.Consume,
	}
	go startUPInterface.GetAllPluginStatus()
	// Run server
	if err := services.Service.Run(); err != nil {
		log.Fatal(err)
	}
}

func registerHandler() {
	events := new(rpc.Events)
	events.IsAuthorizedRPC = services.IsAuthorized
	events.GetSessionUserNameRPC = services.GetSessionUserName
	events.ContactClientRPC = pmbhandle.ContactPlugin
	events.CreateTaskRPC = services.CreateTask
	events.UpdateTaskRPC = evt.UpdateTaskData
	events.CreateChildTaskRPC = services.CreateChildTask
	eventsproto.RegisterEventsHandler(services.Service.Server(), events)
}
