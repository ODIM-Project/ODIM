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

	"github.com/sirupsen/logrus"

	dc "github.com/ODIM-Project/ODIM/lib-messagebus/datacommunicator"
	"github.com/ODIM-Project/ODIM/lib-rest-client/pmbhandle"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/logs"
	aggregatorproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/aggregator"
	"github.com/ODIM-Project/ODIM/lib-utilities/services"
	"github.com/ODIM-Project/ODIM/svc-aggregation/agcommon"
	"github.com/ODIM-Project/ODIM/svc-aggregation/agmessagebus"
	"github.com/ODIM-Project/ODIM/svc-aggregation/agmodel"
	"github.com/ODIM-Project/ODIM/svc-aggregation/rpc"
	"github.com/ODIM-Project/ODIM/svc-aggregation/system"
)

// Schema is a struct to define search, condition and query keys
type Schema struct {
	SearchKeys    []string `json:"searchKeys"`
	ConditionKeys []string `json:"conditionKeys"`
	QueryKeys     []string `json:"queryKeys"`
}

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
		log.Fatal("Aggregation Service should not be run as the root user")
	}

	config.CollectCLArgs(&configWarnings)
	for _, warning := range configWarnings {
		log.Warn(warning)
	}

	if err := dc.SetConfiguration(config.Data.MessageBusConf.MessageBusConfigFilePath); err != nil {
		log.Fatal("error while trying to set message bus configuration: " + err.Error())
	}
	if err := common.CheckDBConnection(); err != nil {
		log.Fatal("error while trying to check DB connection health: " + err.Error())
	}
	var connectionMethodInterface = agcommon.DBInterface{
		GetAllKeysFromTableInterface: agmodel.GetAllKeysFromTable,
		GetConnectionMethodInterface: agmodel.GetConnectionMethod,
		AddConnectionMethodInterface: agmodel.AddConnectionMethod,
		DeleteInterface:              agmodel.Delete,
	}
	if err := connectionMethodInterface.AddConnectionMethods(config.Data.ConnectionMethodConf); err != nil {
		log.Fatal("error while trying add connection method: " + err.Error())
	}

	errChan := make(chan error)
	if err := services.InitializeService(services.Aggregator, errChan); err != nil {
		log.Fatal("fatal: error while trying to initialize service: " + err.Error())
	}

	aggregator := rpc.GetAggregator()
	aggregatorproto.RegisterAggregatorServer(services.ODIMService.Server(), aggregator)

	// Rediscover the Resources by looking in OnDisk DB, populate the resources in InMemory DB
	//This happens only if the InMemory DB lost it contents due to DB reboot or host VM reboot.
	p := system.ExternalInterface{
		ContactClient:   pmbhandle.ContactPlugin,
		Auth:            services.IsAuthorized,
		PublishEventMB:  agmessagebus.Publish,
		GetPluginStatus: agcommon.GetPluginStatus,
		SubscribeToEMB:  services.SubscribeToEMB,
		DecryptPassword: common.DecryptWithPrivateKey,
		UpdateTask:      system.UpdateTaskData,
	}

	go p.RediscoverResources()

	agcommon.ConfigFilePath = os.Getenv("CONFIG_FILE_PATH")
	if agcommon.ConfigFilePath == "" {
		log.Fatal("error: no value get the environment variable CONFIG_FILE_PATH")
	}
	go agcommon.TrackConfigFileChanges(connectionMethodInterface, errChan)

	go system.PerformPluginHealthCheck()

	if err := services.ODIMService.Run(); err != nil {
		log.Fatal("failed to run a service: " + err.Error())
	}
}
