//(C) Copyright [2022] Hewlett Packard Enterprise Development LP
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
	licenseproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/licenses"
	"github.com/ODIM-Project/ODIM/lib-utilities/services"
	lcommon "github.com/ODIM-Project/ODIM/svc-licenses/lcommon"
	"github.com/ODIM-Project/ODIM/svc-licenses/rpc"

	"github.com/sirupsen/logrus"
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

	config.CollectCLArgs(&configWarnings)
	for _, warning := range configWarnings {
		log.Warn(warning)
	}

	errChan := make(chan error)
	if err := common.CheckDBConnection(); err != nil {
		log.Error("error while trying to check DB connection health: " + err.Error())
	}
	lcommon.ConfigFilePath = os.Getenv("CONFIG_FILE_PATH")
	if lcommon.ConfigFilePath == "" {
		log.Fatal("error: no value get the environment variable CONFIG_FILE_PATH")
	}
	// TrackConfigFileChanges monitors the odim config changes using fsnotfiy
	go lcommon.TrackConfigFileChanges(errChan)

	registerHandlers(errChan)

	if err := services.ODIMService.Run(); err != nil {
		log.Error(err)
	}
}

func registerHandlers(errChan chan error) {
	log := logs.Log
	if err := services.InitializeService(services.Licenses, errChan); err != nil {
		log.Error("fatal: error while trying to initialize service: " + err.Error())
	}
	licenses := rpc.GetLicense()
	licenseproto.RegisterLicensesServer(services.ODIMService.Server(), licenses)
}
