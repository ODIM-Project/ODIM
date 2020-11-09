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
	"log"
	"os"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	updateproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/update"
	"github.com/ODIM-Project/ODIM/lib-utilities/services"
	"github.com/ODIM-Project/ODIM/svc-update/rpc"
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

	err := services.InitializeService(services.Update)
	if err != nil {
		log.Fatalf("fatal: error while trying to initialize the service: %v", err)
	}
	registerHandlers()
	// Run server
	if err := services.Service.Run(); err != nil {
		log.Fatal(err)
	}

}

func registerHandlers() {
	err := services.InitializeService(services.Update)
	if err != nil {
		log.Fatalf("fatal: error while trying to initialize service: %v", err)
	}
	updater := rpc.GetUpdater()
	updateproto.RegisterUpdateHandler(services.Service.Server(), updater)
}
