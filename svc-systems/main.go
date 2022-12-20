//(C) Copyright [2020] Hewlett Packard Enterprise Development LP
//(C) Copyright 2020 Intel Corporation
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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/logs"
	chassisproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/chassis"
	systemsproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/systems"
	"github.com/ODIM-Project/ODIM/lib-utilities/services"
	"github.com/ODIM-Project/ODIM/svc-systems/chassis"
	"github.com/ODIM-Project/ODIM/svc-systems/plugin"
	"github.com/ODIM-Project/ODIM/svc-systems/rpc"
	"github.com/ODIM-Project/ODIM/svc-systems/scommon"
	"github.com/ODIM-Project/ODIM/svc-systems/smodel"
	"github.com/ODIM-Project/ODIM/svc-systems/systems"

	"github.com/sirupsen/logrus"
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
		log.Fatal("System Service should not be run as the root user")
	}

	config.CollectCLArgs(&configWarnings)
	for _, warning := range configWarnings {
		log.Warn(warning)
	}

	if err := common.CheckDBConnection(); err != nil {
		log.Fatal("error while trying to check DB connection health: " + err.Error())
	}

	chassis.Token.Tokens = make(map[string]string)

	schemaFile, err := ioutil.ReadFile(config.Data.SearchAndFilterSchemaPath)
	if err != nil {
		log.Fatal("Error while trying to read search/filter schema json: " + err.Error())
	}
	err = json.Unmarshal(schemaFile, &scommon.SF)
	if err != nil {
		log.Fatal("Error while trying to fetch search/filter schema json: " + err.Error())
	}

	configFilePath := os.Getenv("CONFIG_FILE_PATH")
	if configFilePath == "" {
		log.Fatal("error: no value get the environment variable CONFIG_FILE_PATH")
	}

	errChan := make(chan error)
	go scommon.TrackConfigFileChanges(configFilePath, errChan)

	err = services.InitializeService(services.Systems, errChan)
	if err != nil {
		log.Fatal("Error while trying to initialize the service: " + err.Error())
	}

	registerHandler()
	// Run server
	if err := services.ODIMService.Run(); err != nil {
		log.Fatal(err.Error())
	}
}

func registerHandler() {
	systemRPC := new(rpc.Systems)
	systemRPC.IsAuthorizedRPC = services.IsAuthorized
	systemRPC.GetSessionUserName = services.GetSessionUserName
	systemRPC.CreateTask = services.CreateTask
	systemRPC.UpdateTask = systems.UpdateTaskData

	systemRPC.EI = systems.GetExternalInterface()
	systemsproto.RegisterSystemsServer(services.ODIMService.Server(), systemRPC)

	pcf := plugin.NewClientFactory(config.Data.URLTranslation)
	chassisRPC := rpc.NewChassisRPC(
		services.IsAuthorized,
		chassis.NewCreateHandler(pcf),
		chassis.NewGetCollectionHandler(pcf, smodel.GetAllKeysFromTable),
		chassis.NewDeleteHandler(pcf, smodel.Find),
		chassis.NewGetHandler(pcf, smodel.Find),
		chassis.NewUpdateHandler(pcf),
	)

	chassisproto.RegisterChassisServer(services.ODIMService.Server(), chassisRPC)
}
