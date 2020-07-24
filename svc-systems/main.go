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
	"log"
	"os"

	"github.com/bharath-b-hpe/odimra/lib-utilities/common"
	"github.com/bharath-b-hpe/odimra/lib-utilities/config"
	chassisproto "github.com/bharath-b-hpe/odimra/lib-utilities/proto/chassis"
	systemsproto "github.com/bharath-b-hpe/odimra/lib-utilities/proto/systems"
	"github.com/bharath-b-hpe/odimra/lib-utilities/services"
	"github.com/bharath-b-hpe/odimra/svc-systems/rpc"
	"github.com/bharath-b-hpe/odimra/svc-systems/systems"
	"io/ioutil"
)

func main() {
	// verifying the uid of the user
	if uid := os.Geteuid(); uid == 0 {
		log.Fatalln("System Service should not be run as the root user")
	}

	if err := config.SetConfiguration(); err != nil {
		log.Fatalf("fatal: error while trying set up configuration: %v", err)
	}

	if err := common.CheckDBConnection(); err != nil {
		log.Fatalf("error while trying to check DB connection health: %v", err)
	}

	schemaFile, err := ioutil.ReadFile(config.Data.SearchAndFilterSchemaPath)
	if err != nil {
		log.Fatalf("fatal: error while trying to read search/filter schema json: %v", err)
	}
	err = json.Unmarshal(schemaFile, &systems.SF)
	if err != nil {
		log.Fatalf("fatal: error while trying to fetch search/filter schema json: %v", err)
	}

	err = services.InitializeService(services.Systems)
	if err != nil {
		log.Fatalf("fatal: error while trying to initialize the service: %v", err)
	}

	registerHandler()
	// Run server
	if err := services.Service.Run(); err != nil {
		log.Fatal(err)
	}
}

func registerHandler() {
	systemRPC := new(rpc.Systems)
	systemRPC.IsAuthorizedRPC = services.IsAuthorized
	systemsproto.RegisterSystemsHandler(services.Service.Server(), systemRPC)
	chassisRPC := new(rpc.ChassisRPC)
	chassisRPC.IsAuthorizedRPC = services.IsAuthorized
	chassisproto.RegisterChassisHandler(services.Service.Server(), chassisRPC)
}
