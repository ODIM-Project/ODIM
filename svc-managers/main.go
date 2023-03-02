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
	"os"

	"github.com/sirupsen/logrus"

	"fmt"

	dmtf "github.com/ODIM-Project/ODIM/lib-dmtf/model"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/logs"
	managersproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/managers"
	"github.com/ODIM-Project/ODIM/lib-utilities/services"
	"github.com/ODIM-Project/ODIM/svc-managers/managers"
	"github.com/ODIM-Project/ODIM/svc-managers/mgrcommon"
	"github.com/ODIM-Project/ODIM/svc-managers/mgrmodel"
	"github.com/ODIM-Project/ODIM/svc-managers/rpc"
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
		log.Fatal("Manager Service should not be run as the root user")
	}

	config.CollectCLArgs(&configWarnings)
	for _, warning := range configWarnings {
		log.Warn(warning)
	}

	if err := common.CheckDBConnection(); err != nil {
		log.Fatal(err.Error())
	}

	var managerInterface = mgrcommon.DBInterface{
		AddManagertoDBInterface: mgrmodel.AddManagertoDB,
		GenericSave:             mgrmodel.GenericSave,
	}
	err = addManagertoDB(managerInterface)
	if err != nil {
		log.Fatal(err.Error())
	}
	configFilePath := os.Getenv("CONFIG_FILE_PATH")
	if configFilePath == "" {
		log.Fatal("error: no value get the environment variable CONFIG_FILE_PATH")
	}

	errChan := make(chan error)
	go mgrcommon.TrackConfigFileChanges(configFilePath, managerInterface, errChan)

	err = services.InitializeService(services.Managers, errChan)
	if err != nil {
		log.Fatal("fatal: error while trying to initialize service: %v" + err.Error())
	}
	mgrcommon.Token.Tokens = make(map[string]string)
	registerHandlers()
	if err = services.ODIMService.Run(); err != nil {
		log.Fatal("failed to run a service: " + err.Error())
	}
}

func registerHandlers() {
	manager := new(rpc.Managers)

	manager.IsAuthorizedRPC = services.IsAuthorized
	manager.EI = managers.GetExternalInterface()

	managersproto.RegisterManagersServer(services.ODIMService.Server(), manager)
}
func addManagertoDB(managerInterface mgrcommon.DBInterface) error {
	mgr := mgrmodel.RAManager{
		Name:            "odimra",
		ManagerType:     "Service",
		FirmwareVersion: config.Data.FirmwareVersion,
		ID:              config.Data.RootServiceUUID,
		UUID:            config.Data.RootServiceUUID,
		State:           "Enabled",
		Health:          "OK",
		Description:     "Odimra Manager",
		LogServices: &dmtf.Link{
			Oid: "/redfish/v1/Managers/" + config.Data.RootServiceUUID + "/LogServices",
		},
		Model:      "ODIMRA" + " " + config.Data.FirmwareVersion,
		PowerState: "On",
	}
	managerInterface.AddManagertoDBInterface(mgr)

	//adding LogeSrvice Collection
	data := dmtf.Collection{
		ODataContext: "/redfish/v1/$metadata#LogServiceCollection.LogServiceCollection",
		ODataID:      "/redfish/v1/Managers/" + config.Data.RootServiceUUID + "/LogServices",
		ODataType:    "#LogServiceCollection.LogServiceCollection",
		Description:  "Logs view",
		Members: []*dmtf.Link{
			&dmtf.Link{
				Oid: "/redfish/v1/Managers/" + config.Data.RootServiceUUID + "/LogServices/SL",
			},
		},
		MembersCount: 1,
		Name:         "Logs",
	}
	dbdata, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("unable to marshal manager data: %v", err)
	}
	key := "/redfish/v1/Managers/" + config.Data.RootServiceUUID + "/LogServices"
	mgrmodel.GenericSave([]byte(dbdata), "LogServicesCollection", key)

	//adding LogService Members
	logEntrydata := dmtf.LogServices{
		Ocontext:    "/redfish/v1/$metadata#LogServiceCollection.LogServiceCollection",
		Oid:         "/redfish/v1/Managers/" + config.Data.RootServiceUUID + "/LogServices/SL",
		Otype:       "#LogService.v1_3_0.LogService",
		Description: "Logs view",
		Entries: &dmtf.Entries{
			Oid: "/redfish/v1/Managers/" + config.Data.RootServiceUUID + "/LogServices/SL/Entries",
		},
		ID:              "SL",
		Name:            "Security Log",
		OverWritePolicy: "WrapsWhenFull",
	}
	dbdata, err = json.Marshal(logEntrydata)
	if err != nil {
		return fmt.Errorf("unable to marshal manager data: %v", err)
	}
	key = "/redfish/v1/Managers/" + config.Data.RootServiceUUID + "/LogServices/SL"
	mgrmodel.GenericSave([]byte(dbdata), "LogServices", key)

	// adding empty logservice entry collection
	entriesdata := dmtf.Collection{
		ODataContext: "/redfish/v1/$metadata#LogServiceCollection.LogServiceCollection",
		ODataID:      "/redfish/v1/Managers/" + config.Data.RootServiceUUID + "/LogServices/SL/Entries",
		ODataType:    "#LogEntryCollection.LogEntryCollection",
		Description:  "Security Logs view",
		Members:      []*dmtf.Link{},
		MembersCount: 0,
		Name:         "Security Logs",
	}
	dbentriesdata, err := json.Marshal(entriesdata)
	if err != nil {
		return fmt.Errorf("unable to marshal manager data: %v", err)
	}
	key = "/redfish/v1/Managers/" + config.Data.RootServiceUUID + "/LogServices/SL/Entries"
	mgrmodel.GenericSave([]byte(dbentriesdata), "EntriesCollection", key)

	return nil
}
