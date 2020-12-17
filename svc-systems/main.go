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
	"encoding/json"
	"github.com/sirupsen/logrus"
	"os"

	"io/ioutil"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	chassisproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/chassis"
	systemsproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/systems"
	"github.com/ODIM-Project/ODIM/lib-utilities/services"
	"github.com/ODIM-Project/ODIM/svc-systems/rpc"
	"github.com/ODIM-Project/ODIM/svc-systems/systems"
)

var log = logrus.New()

func main() {
	// verifying the uid of the user
	if uid := os.Geteuid(); uid == 0 {
		log.Error("System Service should not be run as the root user")
	}

	if err := config.SetConfiguration(); err != nil {
		log.Error("fatal: error while trying set up configuration: " + err.Error())
	}

	if err := common.CheckDBConnection(); err != nil {
		log.Error("error while trying to check DB connection health: " + err.Error())
	}

	schemaFile, err := ioutil.ReadFile(config.Data.SearchAndFilterSchemaPath)
	if err != nil {
		log.Error("fatal: error while trying to read search/filter schema json: " + err.Error())
	}
	err = json.Unmarshal(schemaFile, &systems.SF)
	if err != nil {
		log.Error("fatal: error while trying to fetch search/filter schema json: " + err.Error())
	}

	err = services.InitializeService(services.Systems)
	if err != nil {
		log.Error("fatal: error while trying to initialize the service: " + err.Error())
	}

	registerHandler()
	// Run server
	if err := services.Service.Run(); err != nil {
		log.Error(err.Error())
	}
}

func registerHandler() {
	systemRPC := new(rpc.Systems)
	systemRPC.IsAuthorizedRPC = services.IsAuthorized
	systemRPC.EI = systems.GetExternalInterface()
	systemsproto.RegisterSystemsHandler(services.Service.Server(), systemRPC)
	chassisRPC := new(rpc.ChassisRPC)
	chassisRPC.IsAuthorizedRPC = services.IsAuthorized
	chassisproto.RegisterChassisHandler(services.Service.Server(), chassisRPC)
}
