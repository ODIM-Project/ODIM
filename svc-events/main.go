// (C) Copyright [2020] Hewlett Packard Enterprise Development LP
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.
package main

import (
	"context"
	"fmt"
	"os"

	dc "github.com/ODIM-Project/ODIM/lib-messagebus/datacommunicator"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/logs"
	eventsproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/events"
	"github.com/ODIM-Project/ODIM/lib-utilities/services"
	"github.com/ODIM-Project/ODIM/svc-events/consumer"
	"github.com/ODIM-Project/ODIM/svc-events/evcommon"
	"github.com/ODIM-Project/ODIM/svc-events/rpc"
	"github.com/sirupsen/logrus"
)

var (
	processControlActionID   = "219"
	processControlActionName = "ProcessControl"
)

func main() {
	// setting up the logging framework
	hostName := os.Getenv("HOST_NAME")
	podName := os.Getenv("POD_NAME")
	pid := os.Getpid()
	logs.Adorn(logrus.Fields{
		"host":       hostName,
		"process_id": podName + fmt.Sprintf("_%d", pid),
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
		log.Fatal("Event Service should not be run as the root user")
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

	errChan := make(chan error)
	if err := services.InitializeService(services.Events, errChan); err != nil {
		log.Fatal("fatal: error while trying to initialize the service: " + err.Error())
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, common.ProcessName, podName)
	ctx = context.WithValue(ctx, common.ThreadName, common.EventService)
	ctx = context.WithValue(ctx, common.ThreadID, common.DefaultThreadID)
	// Initializing the TopicsList
	evcommon.EMBTopics.TopicsList = make(map[string]bool)
	// Initializing plugin token
	evcommon.Token.Tokens = make(map[string]string)

	// register handlers
	events := rpc.GetPluginContactInitializer()
	eventsproto.RegisterEventsServer(services.ODIMService.Server(), events)

	// Load event cache
	if err = events.Connector.LoadSubscriptionData(ctx); err != nil {
		log.Fatal("fatal: error while trying to load cache : " + err.Error())
	}
	// CreateJobQueue defines the queue which will act as an infinite buffer
	// In channel is an entry or input channel and the Out channel is an exit or output channel
	jobQueueSize := 10
	consumer.In, consumer.Out = common.CreateJobQueue(jobQueueSize)
	// RunReadWorkers will create a worker pool for doing a specific task
	// which is passed to it as PublishEventsToDestination method after reading the data from the channel.
	common.RunReadWorkers(ctx, consumer.Out, events.Connector.PublishEventsToDestination, 5)

	// CreateJobQueue defines the queue which will act as an infinite buffer
	// In channel is an entry or input channel and the Out channel is an exit or output channel
	ctrlMsgProcQueueSize := 1
	consumer.CtrlMsgRecvQueue, consumer.CtrlMsgProcQueue = common.CreateJobQueue(ctrlMsgProcQueueSize)
	// RunReadWorkers will create a worker pool for doing a specific task
	// which is passed to it as ProcessCtrlMsg method after reading the data from the channel.
	ctx = context.WithValue(ctx, common.ActionName, processControlActionName)
	ctx = context.WithValue(ctx, common.ActionID, processControlActionID)
	common.RunReadWorkers(ctx, consumer.CtrlMsgProcQueue, evcommon.ProcessCtrlMsg, 1)

	evcommon.ConfigFilePath = os.Getenv("CONFIG_FILE_PATH")
	if evcommon.ConfigFilePath == "" {
		log.Fatal("error: no value get the environment variable CONFIG_FILE_PATH")
	}
	// TrackConfigFileChanges monitors the odim config changes using fsnotfiy
	go evcommon.TrackConfigFileChanges(ctx, errChan)

	// Subscribe to inter communication message bus queue
	go consumer.SubscribeCtrlMsgQueue(config.Data.MessageBusConf.OdimControlMessageQueue)

	// Subscribe to EMBs of all the available plugins
	startUPInterface := evcommon.StartUpInterface{
		DecryptPassword: common.DecryptWithPrivateKey,
		EMBConsume:      consumer.Consume,
	}
	go startUPInterface.SubscribePluginEMB(ctx)

	// Run server
	if err := services.ODIMService.Run(); err != nil {
		log.Fatal(err.Error())
	}
}
