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
	"strings"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	managersproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/managers"
	"github.com/ODIM-Project/ODIM/lib-utilities/services"
	"github.com/ODIM-Project/ODIM/svc-managers/mgrcommon"
	"github.com/ODIM-Project/ODIM/svc-managers/mgrmodel"
	"github.com/ODIM-Project/ODIM/svc-managers/rpc"
	"github.com/ODIM-Project/ODIM/svc-plugin-rest-client/pmbhandle"
	"github.com/ODIM-Project/ODIM/svc-managers/managers"
)

func main() {

	// verifying the uid of the user
	if uid := os.Geteuid(); uid == 0 {
		log.Fatalln("Fabric Service should not be run as the root user")
	}

	if err := config.SetConfiguration(); err != nil {
		log.Fatalf("fatal: error while trying set up configuration: %v", err)
	}

	if err := common.CheckDBConnection(); err != nil {
		log.Fatalf("error while trying to check DB connection health: %v", err)
	}

	err := addManagertoDB()
	if err != nil && !(strings.Contains(err.Error(), "already exists")) {
		log.Fatalf("error while trying to add  manager details into DB: %v", err)
	}
	err = services.InitializeService(services.Managers)
	if err != nil {
		log.Fatalf("fatal: error while trying to initialize service: %v", err)
	}
	mgrcommon.Token.Tokens = make(map[string]string)
	registerHandlers()
	if err = services.Service.Run(); err != nil {
		log.Fatal("failed to run a service: ", err)
	}
}

func registerHandlers() {
	managers := new(rpc.Managers)

	managers.IsAuthorizedRPC = services.IsAuthorized
	managers.ContactClientRPC = pmbhandle.ContactPlugin
	managers.EI = managers.GetExternalInterface()

	managersproto.RegisterManagersHandler(services.Service.Server(), managers)
}

func addManagertoDB() error {
	mgr := mgrmodel.RAManager{
		Name:            "odimra",
		ManagerType:     "Service",
		FirmwareVersion: config.Data.FirmwareVersion,
		ID:              config.Data.RootServiceUUID,
		UUID:            config.Data.RootServiceUUID,
		State:           "Enabled",
	}
	return mgr.AddManagertoDB()
}
