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
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	dc "github.com/ODIM-Project/ODIM/lib-messagebus/datacommunicator"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/logs"
	taskproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/task"
	"github.com/ODIM-Project/ODIM/lib-utilities/services"
	auth "github.com/ODIM-Project/ODIM/svc-task/tauth"
	"github.com/ODIM-Project/ODIM/svc-task/tcommon"
	"github.com/ODIM-Project/ODIM/svc-task/thandle"
	"github.com/ODIM-Project/ODIM/svc-task/tmessagebus"
	"github.com/ODIM-Project/ODIM/svc-task/tmodel"
	"github.com/ODIM-Project/ODIM/svc-task/tqueue"
)

func main() {
	// setting up the logging framework
	hostName := os.Getenv("HOST_NAME")
	podName := os.Getenv("POD_NAME")
	pid := os.Getpid()
	logs.Adorn(logrus.Fields{
		"host":   hostName,
		"procid": podName + fmt.Sprintf("_%d", pid),
	})

	// log should be initialized after Adorn is invoked
	// as Adorn will assign new pointer to Log variable in logs package.
	log := logs.Log
	configWarnings, err := config.SetConfiguration()
	if err != nil {
		log.Logger.SetFormatter(&logs.SysLogFormatter{})
		log.Fatal("Error while trying set up configuration: " + err.Error())
	}
	logs.SetFormatter(config.Data.LogFormat)
	log.Logger.SetOutput(os.Stdout)
	log.Logger.SetLevel(config.Data.LogLevel)

	// verifying the uid of the user
	if uid := os.Geteuid(); uid == 0 {
		log.Fatal("Task Service should not be run as the root user")
	}

	config.CollectCLArgs(&configWarnings)
	for _, warning := range configWarnings {
		log.Warn(warning)
	}

	if err := dc.SetConfiguration(config.Data.MessageBusConf.MessageBusConfigFilePath); err != nil {
		log.Fatal("error while trying to set messagebus configuration: " + err.Error())
	}
	if err := common.CheckDBConnection(); err != nil {
		log.Fatal("error while trying to check DB connection health: " + err.Error())
	}
	tcommon.ConfigFilePath = os.Getenv("CONFIG_FILE_PATH")
	if tcommon.ConfigFilePath == "" {
		log.Fatal("error: no value get the environment variable CONFIG_FILE_PATH")
	}

	errChan := make(chan error)
	// TrackConfigFileChanges monitors the odim config changes using fsnotfiy
	go tcommon.TrackConfigFileChanges(errChan)

	if err := services.InitializeService(services.Tasks, errChan); err != nil {
		log.Fatal("fatal: error while trying to initialize the service: " + err.Error())
	}

	tqueue.NewTaskQueue(config.Data.TaskQueueConf.QueueSize)

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
	task.UpdateTaskQueue = tqueue.EnqueueTask
	task.PersistTaskModel = tmodel.PersistTask
	task.ValidateTaskUserNameModel = tmodel.ValidateTaskUserName
	task.PublishToMessageBus = tmessagebus.Publish
	thandle.TaskCollection = thandle.TaskCollectionData{
		TaskCollection: make(map[string]int32),
		Lock:           sync.Mutex{},
	}
	taskproto.RegisterGetTaskServiceServer(services.ODIMService.Server(), task)

	tick := &tmodel.Tick{
		Ticker: time.NewTicker(time.Duration(config.Data.TaskQueueConf.DBCommitInterval) * time.Microsecond),
	}
	go tqueue.UpdateTasksWorker(tick)

	// Run server
	if err := services.ODIMService.Run(); err != nil {
		log.Fatal(err.Error())
	}
}
