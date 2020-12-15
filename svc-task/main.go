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
	"sync"

	dc "github.com/ODIM-Project/ODIM/lib-messagebus/datacommunicator"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	taskproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/task"
	"github.com/ODIM-Project/ODIM/lib-utilities/services"
	auth "github.com/ODIM-Project/ODIM/svc-task/tauth"
	"github.com/ODIM-Project/ODIM/svc-task/thandle"
	"github.com/ODIM-Project/ODIM/svc-task/tmessagebus"
	"github.com/ODIM-Project/ODIM/svc-task/tmodel"
)

func main() {
	// verifying the uid of the user
	if uid := os.Geteuid(); uid == 0 {
		log.Fatalln("Task Service should not be run as the root user")
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
	configFilePath := os.Getenv("CONFIG_FILE_PATH")
	if configFilePath == "" {
		log.Fatalln("error: no value get the environment variable CONFIG_FILE_PATH")
	}
	eventChan := make(chan interface{})
	var lock sync.Mutex
	// TrackConfigFileChanges monitors the odim config changes using fsnotfiy
	go common.TrackConfigFileChanges(configFilePath, eventChan, &lock)

	if err := services.InitializeService(services.Tasks); err != nil {
		log.Fatalf("fatal: error while trying to initialize the service: %v", err)
	}

	task := new(thandle.TasksRPC)
	task.AuthenticationRPC = auth.Authentication
	task.GetSessionUserNameRPC = auth.GetSessionUserName
	task.GetTaskStatusModel = tmodel.GetTaskStatus
	task.GetAllTaskKeysModel = tmodel.GetAllTaskKeys
	task.TransactionModel = tmodel.Transaction
	task.OverWriteCompletedTaskUtilHelper = task.OverWriteCompletedTaskUtil
	task.CreateTaskUtilHelper = task.CreateTaskUtil
	task.GetCompletedTasksIndexModel = tmodel.GetCompletedTasksIndex
	task.DeleteTaskFromDBModel = tmodel.DeleteTaskFromDB
	task.DeleteTaskIndex = tmodel.DeleteTaskIndex
	task.UpdateTaskStatusModel = tmodel.UpdateTaskStatus
	task.PersistTaskModel = tmodel.PersistTask
	task.ValidateTaskUserNameModel = tmodel.ValidateTaskUserName
	task.PublishToMessageBus = tmessagebus.Publish

	taskproto.RegisterGetTaskServiceHandler(services.Service.Server(), task)

	// Run server
	if err := services.Service.Run(); err != nil {
		log.Fatal(err)
	}
}
