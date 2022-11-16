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

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/logs"
	accountproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/account"
	authproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/auth"
	roleproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/role"
	sessionproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/session"
	"github.com/ODIM-Project/ODIM/lib-utilities/services"
	"github.com/ODIM-Project/ODIM/svc-account-session/account"
	"github.com/ODIM-Project/ODIM/svc-account-session/rpc"
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
		log.Fatal(err.Error())
	}

	log.Logger.SetFormatter(&logs.SysLogFormatter{})
	log.Logger.SetOutput(os.Stdout)
	log.Logger.SetLevel(config.Data.LogLevel)

	// verifying the uid of the user
	if uid := os.Geteuid(); uid == 0 {
		log.Fatal("AccountSession Service should not be run as the root user")
	}

	config.CollectCLArgs()

	if err := common.CheckDBConnection(); err != nil {
		log.Fatal("Error while trying to check DB connection health: " + err.Error())
	}

	account.ConfigFilePath = os.Getenv("CONFIG_FILE_PATH")
	if account.ConfigFilePath == "" {
		log.Fatal("error: no value get the environment variable CONFIG_FILE_PATH")
	}
	// TrackConfigFileChanges monitors the odim config changes using fsnotfiy
	go account.TrackConfigFileChanges()

	if err := services.InitializeService(services.AccountSession); err != nil {
		log.Fatal("Error while trying to initialize the service: " + err.Error())
	}

	registerHandlers()
	if err := services.ODIMService.Run(); err != nil {
		log.Fatal("Failed to run a service: " + err.Error())
	}
}

func registerHandlers() {
	authproto.RegisterAuthorizationServer(services.ODIMService.Server(), new(rpc.Auth))
	sessionproto.RegisterSessionServer(services.ODIMService.Server(), new(rpc.Session))
	accountproto.RegisterAccountServer(services.ODIMService.Server(), new(rpc.Account))
	roleproto.RegisterRolesServer(services.ODIMService.Server(), new(rpc.Role))
}
