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
	"os"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	teleproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/telemetry"
	"github.com/ODIM-Project/ODIM/lib-utilities/services"
	"github.com/ODIM-Project/ODIM/svc-telemetry/rpc"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func main() {
	//log.SetFormatter(&log.TextFormatter{})
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
	configFilePath := os.Getenv("CONFIG_FILE_PATH")
	if configFilePath == "" {
		log.Fatal("error: no value get the environment variable CONFIG_FILE_PATH")
	}
	eventChan := make(chan interface{})
	// TrackConfigFileChanges monitors the odim config changes using fsnotfiy
	go common.TrackConfigFileChanges(configFilePath, eventChan)

	registerHandlers()
	// Run server
	if err := services.ODIMService.Run(); err != nil {
		log.Error(err)
	}

}

func registerHandlers() {
	err := services.InitializeService(services.Telemetry)
	if err != nil {
		log.Error("fatal: error while trying to initialize service: " + err.Error())
	}
	tele := rpc.GetTele()
	teleproto.RegisterTelemetryServer(services.ODIMService.Server(), tele)
}
