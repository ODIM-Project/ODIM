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

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/logs"
	teleproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/telemetry"
	"github.com/ODIM-Project/ODIM/lib-utilities/services"
	"github.com/ODIM-Project/ODIM/svc-telemetry/rpc"
	"github.com/ODIM-Project/ODIM/svc-telemetry/tcommon"

	"github.com/sirupsen/logrus"
)

func main() {
	// setting up the logging framework
	hostName := os.Getenv("HOST_NAME")
	podName := os.Getenv("POD_NAME")
	pid := os.Getpid()
	log := logs.Log
	logs.Adorn(logrus.Fields{
		"host":   hostName,
		"procid": podName + fmt.Sprintf("_%d", pid),
	})

	if err := config.SetConfiguration(); err != nil {
		log.Logger.SetFormatter(&logs.SysLogFormatter{})
		log.Error("fatal: error while trying set up configuration: " + err.Error())
	}

	log.Logger.SetFormatter(&logs.SysLogFormatter{})
	log.Logger.SetOutput(os.Stdout)
	log.Logger.SetLevel(config.Data.LogLevel)

	// verifying the uid of the user
	if uid := os.Geteuid(); uid == 0 {
		log.Error("Telemetry Service should not be run as the root user")
	}

	if err := config.SetConfiguration(); err != nil {
		log.Error("fatal: error while trying set up configuration: " + err.Error())
	}

	config.CollectCLArgs()
	if err := common.CheckDBConnection(); err != nil {
		log.Error("error while trying to check DB connection health: " + err.Error())
	}
	tcommon.ConfigFilePath = os.Getenv("CONFIG_FILE_PATH")
	if tcommon.ConfigFilePath == "" {
		log.Fatal("error: no value get the environment variable CONFIG_FILE_PATH")
	}
	// TrackConfigFileChanges monitors the odim config changes using fsnotfiy
	go tcommon.TrackConfigFileChanges()

	registerHandlers()
	// Run server
	if err := services.ODIMService.Run(); err != nil {
		log.Error(err)
	}

}

func registerHandlers() {
	log := logs.Log
	err := services.InitializeService(services.Telemetry)
	if err != nil {
		log.Error("fatal: error while trying to initialize service: " + err.Error())
	}
	tele := rpc.GetTele()
	teleproto.RegisterTelemetryServer(services.ODIMService.Server(), tele)
}
