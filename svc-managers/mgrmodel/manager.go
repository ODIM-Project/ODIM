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

//Package mgrmodel ....
package mgrmodel

import (
	"encoding/json"
	"fmt"

	dmtf "github.com/ODIM-Project/ODIM/lib-dmtf/model"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
)

// Manager struct for manager deta
type Manager struct {
	OdataContext       string            `json:"@odata.context"`
	Etag               string            `json:"@odata.etag,omitempty"`
	OdataID            string            `json:"@odata.id"`
	OdataType          string            `json:"@odata.type"`
	Name               string            `json:"Name"`
	ManagerType        string            `json:"ManagerType"`
	ID                 string            `json:"Id"`
	UUID               string            `json:"UUID"`
	FirmwareVersion    string            `json:"FirmwareVersion"`
	Status             *Status           `json:"Status,omitempty"`
	HostInterfaces     *OdataID          `json:"HostInterfaces,omitempty"`
	SerialInterface    *OdataID          `json:"SerialInterface,omitempty"`
	EthernetInterfaces *OdataID          `json:"EthernetInterfaces,omitempty"`
	LogServices        *OdataID          `json:"LogServices,omitempty"`
	NetworkProtocol    *OdataID          `json:"NetworkProtocol,omitempty"`
	VirtualMedia       *OdataID          `json:"VirtualMedia,omitempty"`
	CommandShell       *CommandShell     `json:"CommandShell,omitempty"`
	GraphicalConsole   *GraphicalConsole `json:"GraphicalConsole,omitempty"`
	Links              *Links            `json:"Links,omitempty"`
	Actions            *Actions          `json:"Actions,omitempty"`
}

// Status struct is to define the status of the manager
type Status struct {
	State string `json:"State"`
}

// OdataID is link
type OdataID struct {
	OdataID []dmtf.Link `json:"OdataID"`
}

// CommandShell service that manager provides.
type CommandShell struct {
	ConnectTypesSupported []string `json:"ConnectTypesSupported"`
	MaxConcurrentSessions int      `json:"MaxConcurrentSessions"`
	ServiceEnabled        bool     `json:"ServiceEnabled"`
}

// GraphicalConsole is the information about the graphical console (KVM-IP)service of the manager.
type GraphicalConsole struct {
	ConnectTypesSupported []string `json:"ConnectTypesSupported"`
	MaxConcurrentSessions int      `json:"MaxConcurrentSessions"`
	ServiceEnabled        bool     `json:"ServiceEnabled"`
}

// Links to other Resources that are related to this Resource.
type Links struct {
	ActiveSoftwareImage OdataID   `json:"ActiveSoftwareImage"`
	ManagerForChassis   []OdataID `json:"ManagerForChassis"`
	ManagerForServers   []OdataID `json:"ManagerForServers"`
	ManagerForSwitches  []OdataID `json:"ManagerForSwitches"`
	ManagerInChassis    OdataID   `json:"ManagerInChassis"`
}

// Actions struct for Actions to perform
type Actions struct {
	Reset Target `json:"#Manager.Reset"`
}

// Target ...
type Target struct {
	Target string `json:"target"`
}

// RAManager struct is to store odimra details into DB
type RAManager struct {
	ID              string `json:"ManagerID"`
	Name            string `json:"Name"`
	ManagerType     string `json:"ManagerType"`
	FirmwareVersion string `json:"FirmwareVersion"`
	UUID            string `json:"UUID"`
	State           string `json:"State"`
}

// VirtualMediaInsert struct is to store the insert virtual media request payload
type VirtualMediaInsert struct {
	Image                string `json:"Image" validate:"required"`
	Inserted             bool   `json:"Inserted"`
	WriteProtected       bool   `json:"WriteProtected"`
	Password             string `json:"Password,omitempty"`
	TransferMethod       string `json:"TransferMethod,omitempty"`
	TransferProtocolType string `json:"TransferProtocolType,omitempty"`
	UserName             string `json:"UserName,omitempty"`
}

//GetResource fetches a resource from database using table and key
func GetResource(Table, key string) (string, *errors.Error) {
	conn, err := common.GetDBConnection(common.InMemory)
	if err != nil {
		return "", err
	}
	resourceData, err := conn.Read(Table, key)
	if err != nil {
		return "", errors.PackError(err.ErrNo(), "unable to get resource details: ", err.Error())
	}
	var resource string
	if errs := json.Unmarshal([]byte(resourceData), &resource); errs != nil {
		return "", errors.PackError(errors.UndefinedErrorType, errs)
	}
	return resource, nil
}

//GetAllKeysFromTable fetches all keys in a given table
func GetAllKeysFromTable(table string) ([]string, error) {
	conn, err := common.GetDBConnection(common.InMemory)
	if err != nil {
		return nil, err
	}
	keysArray, err := conn.GetAllDetails(table)
	if err != nil {
		return nil, fmt.Errorf("unable to get all keys from table - %v: %v", table, err.Error())
	}
	return keysArray, nil
}

// GetManagerByURL fetches computer manager details by URL from database
func GetManagerByURL(url string) (string, *errors.Error) {
	var manager string
	conn, err := common.GetDBConnection(common.InMemory)
	if err != nil {
		// connection error
		return manager, err
	}
	managerData, err := conn.Read("Managers", url)
	if err != nil {
		return "", errors.PackError(err.ErrNo(), "unable to get managers details: ", err.Error())
	}
	if errs := json.Unmarshal([]byte(managerData), &manager); errs != nil {
		return "", errors.PackError(errors.UndefinedErrorType, errs)
	}
	return manager, nil
}

// UpdateData will modify the current details to given changes
func UpdateData(key string, updateData map[string]interface{}, table string) error {

	conn, err := common.GetDBConnection(common.InMemory)
	if err != nil {
		return fmt.Errorf("unable to connect DB: %v", err)
	}
	data, jerr := json.Marshal(updateData)
	if jerr != nil {
		return fmt.Errorf("unable to marshal data for updating: %v", jerr)
	}
	if _, err = conn.Update(table, key, string(data)); err != nil {
		return fmt.Errorf("unable to update details in DB: %v", err)
	}
	return nil
}

//GenericSave will save any resource data into the database
func GenericSave(body []byte, table string, key string) error {

	connPool, err := common.GetDBConnection(common.InMemory)
	if err != nil {
		return fmt.Errorf("unable to connect DB: %v", err.Error())
	}
	if err := connPool.Create(table, key, string(body)); err != nil {
		return fmt.Errorf("%v", err)
	}
	return nil
}

// AddManagertoDB will add odimra Manager details to DB
func AddManagertoDB(mgr RAManager) error {
	key := "/redfish/v1/Managers/" + mgr.UUID
	data, err := json.Marshal(mgr)
	if err != nil {
		return fmt.Errorf("unable to marshal manager data: %v", err)
	}
	connPool, connErr := common.GetDBConnection(common.InMemory)
	if connErr != nil {
		return fmt.Errorf("unable to connect DB: %v", connErr.Error())
	}
	if err := connPool.AddResourceData("Managers", key, string(data)); err != nil {
		return fmt.Errorf("%v", err.Error())
	}
	return nil
}
