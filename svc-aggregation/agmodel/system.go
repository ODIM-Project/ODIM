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

// Package agmodel ...
package agmodel

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	dmtfmodel "github.com/ODIM-Project/ODIM/lib-dmtf/model"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	log "github.com/sirupsen/logrus"
)

//Schema model is used to iterate throgh the schema json for search/filter
type Schema struct {
	SearchKeys    []map[string]map[string]string `json:"searchKeys"`
	ConditionKeys []string                       `json:"conditionKeys"`
	QueryKeys     []string                       `json:"queryKeys"`
}

//SaveSystem model is used to save encrypted data into db
type SaveSystem struct {
	ManagerAddress string
	Password       []byte
	UserName       string
	DeviceUUID     string
	PluginID       string
}

// Plugin is the model for plugin information
type Plugin struct {
	IP                string
	Port              string
	Username          string
	Password          []byte
	ID                string
	PluginType        string
	PreferredAuthType string
	ManagerUUID       string
}

//Target is for sending the requst to south bound/plugin
type Target struct {
	ManagerAddress string `json:"ManagerAddress"`
	Password       []byte `json:"Password"`
	UserName       string `json:"UserName"`
	PostBody       []byte `json:"PostBody"`
	DeviceUUID     string `json:"DeviceUUID"`
	PluginID       string `json:"PluginID"`
}

//SystemOperation hold the value system operation(InventoryRediscovery or Delete)
type SystemOperation struct {
	Operation string
}

// AggregationSource  payload of adding a AggregationSource
type AggregationSource struct {
	HostName string
	UserName string
	Password []byte
	Links    interface{}
}

// Aggregate payload is used for perform the operations on Aggregate
type Aggregate struct {
	Elements []string `json:"Elements"`
}

// ConnectionMethod payload is used for perform the operations on connection method
type ConnectionMethod struct {
	ConnectionMethodType    string `json:"ConnectionMethodType"`
	ConnectionMethodVariant string `json:"ConnectionMethodVariant"`
	Links                   Links  `json:"Links"`
}

// Links is payload of aggregation resources
type Links struct {
	AggregationSources []OdataID `json:"AggregationSources"`
}

//OdataID struct definition for @odata.id
type OdataID struct {
	OdataID string `json:"@odata.id"`
}

//ServerInfo holds the details of the server
type ServerInfo SaveSystem

// PluginStartUpData holds the required data for plugin startup
type PluginStartUpData struct {
	RequestType           string
	ResyncEvtSubscription bool
	Devices               map[string]DeviceData
}

// DeviceData holds device credentials, event subcription and trigger details
type DeviceData struct {
	UserName              string
	Password              []byte
	Address               string
	Operation             string
	EventSubscriptionInfo *EventSubscriptionInfo
	TriggerInfo           *TriggerInfo
}

// EventSubscriptionInfo holds the event subscription details of a device
type EventSubscriptionInfo struct {
	EventTypes []string
	Location   string
}

// TriggerInfo holds the metric trigger info of a device
type TriggerInfo struct {
}

//PluginContactRequest holds the details required to contact the plugin
type PluginContactRequest struct {
	URL             string
	HTTPMethodType  string
	ContactClient   func(string, string, string, string, interface{}, map[string]string) (*http.Response, error)
	PostBody        interface{}
	LoginCredential map[string]string
	Token           string
	Plugin          Plugin
}

//GetResource fetches a resource from database using table and key
func GetResource(Table, key string) (string, *errors.Error) {
	conn, err := common.GetDBConnection(common.InMemory)
	if err != nil {
		return "", errors.PackError(err.ErrNo(), err)
	}
	resourceData, err := conn.Read(Table, key)
	if err != nil {
		return "", errors.PackError(err.ErrNo(), "error while trying to get resource details: ", err.Error())
	}
	var resource string
	if errs := json.Unmarshal([]byte(resourceData), &resource); errs != nil {
		return "", errors.PackError(errors.UndefinedErrorType, errs)
	}
	return resource, nil
}

// Create connects to the persistencemgr and creates a system in db
func (system *SaveSystem) Create(systemID string) *errors.Error {

	conn, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		log.Error("error while trying to get Db connection : " + err.Error())
		return err
	}
	//Create a header for data entry
	const table string = "System"
	//Save data into Database
	if err = conn.Create(table, systemID, system); err != nil {
		log.Error("error while trying to save system data in DB : " + err.Error())
		return err
	}
	return nil
}

//GetPluginData will fetch plugin details
func GetPluginData(pluginID string) (Plugin, *errors.Error) {
	var plugin Plugin

	conn, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		return plugin, errors.PackError(err.ErrNo(), "error while trying to connect to DB: ", err.Error())
	}

	plugindata, err := conn.Read("Plugin", pluginID)
	if err != nil {
		return plugin, errors.PackError(err.ErrNo(), "error while trying to fetch plugin data: ", err.Error())
	}

	if err := json.Unmarshal([]byte(plugindata), &plugin); err != nil {
		return plugin, errors.PackError(errors.JSONUnmarshalFailed, err)
	}

	bytepw, errs := common.DecryptWithPrivateKey([]byte(plugin.Password))
	if errs != nil {
		return Plugin{}, errors.PackError(errors.DecryptionFailed, "error: "+pluginID+" plugin password decryption failed: "+errs.Error())
	}
	plugin.Password = bytepw

	return plugin, nil
}

//GetComputeSystem will fetch the compute resource details
func GetComputeSystem(deviceUUID string) (dmtfmodel.ComputerSystem, error) {
	var compute dmtfmodel.ComputerSystem

	conn, err := common.GetDBConnection(common.InMemory)
	if err != nil {
		log.Error("GetComputeSystem : error while trying to get db conenction : " + err.Error())
		return compute, err
	}

	computeData, err := conn.Read("ComputerSystem", deviceUUID)
	if err != nil {
		return compute, fmt.Errorf("error while trying to get compute details: %v", err.Error())
	}

	if err := json.Unmarshal([]byte(computeData), &compute); err != nil {
		log.Error("GetComputeSystem : error while Unmarshaling data : " + err.Error())
		return compute, err
	}
	return compute, nil

}

//SaveComputeSystem will save the compute server complete details into the database
func SaveComputeSystem(computeServer dmtfmodel.ComputerSystem, deviceUUID string) error {
	//use dmtf logic to save data into database
	log.Info("Saving server details into database")
	err := computeServer.SaveInMemory(deviceUUID)
	if err != nil {
		log.Error("error while trying to save server details in DB : " + err.Error())
		return err
	}
	return nil
}

//SaveChassis will save the chassis details into the database
func SaveChassis(chassis dmtfmodel.Chassis, deviceUUID string) error {
	//use dmtf logic to save data into database
	log.Info("Saving chassis details into database")
	err := chassis.SaveInMemory(deviceUUID)
	if err != nil {
		log.Error("error while trying to save chassis details in DB : " + err.Error())
		return err
	}
	return nil
}

// GenericSave will save any resource data into the database
func GenericSave(body []byte, table string, key string) error {

	connPool, err := common.GetDBConnection(common.InMemory)
	if err != nil {
		log.Error("GenericSave : error while trying to get DB Connection : " + err.Error())
		return fmt.Errorf("error while trying to connecting to DB: %v", err.Error())
	}
	if err = connPool.AddResourceData(table, key, string(body)); err != nil {
		log.Error("GenericSave : error while trying to add resource date to DB: " + err.Error())
		return fmt.Errorf("error while trying to create new %v resource: %v", table, err.Error())
	}
	return nil
}

//SaveRegistryFile will save any Registry file in database OnDisk DB
func SaveRegistryFile(body []byte, table string, key string) error {

	connPool, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		return fmt.Errorf("error while trying to connecting to DB: %v", err.Error())
	}
	if err = connPool.Create(table, key, string(body)); err != nil {
		if errors.DBKeyAlreadyExist != err.ErrNo() {
			return fmt.Errorf("error while trying to create new %v resource: %v", table, err.Error())
		}
		log.Warn("Skipped saving of duplicate data with key " + key)
		return nil
	}
	return nil
}

//GetRegistryFile from Onisk DB
func GetRegistryFile(Table, key string) (string, *errors.Error) {
	conn, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		return "", errors.PackError(err.ErrNo(), err)
	}
	resourceData, err := conn.Read(Table, key)
	if err != nil {
		return "", errors.PackError(err.ErrNo(), "error while trying to get resource details: ", err.Error())
	}
	var resource string
	if errs := json.Unmarshal([]byte(resourceData), &resource); errs != nil {
		return "", errors.PackError(errors.UndefinedErrorType, errs)
	}
	return resource, nil
}

//DeleteComputeSystem will delete the compute system
func DeleteComputeSystem(index int, key string) *errors.Error {
	connPool, err := common.GetDBConnection(common.InMemory)
	if err != nil {
		return errors.PackError(err.ErrNo(), "error while trying to connecting to DB: ", err.Error())
	}

	// Check key present in the DB
	if _, err = connPool.Read("ComputerSystem", key); err != nil {
		return errors.PackError(err.ErrNo(), "error while trying to get compute details: ", err.Error())
	}
	var computeData, inventoryData []string
	editedKeyList := strings.Split(key, "/")
	editedKey := editedKeyList[len(editedKeyList)-1]
	systemID := strings.Split(editedKey, ":")[0]
	if computeData, err = connPool.GetAllMatchingDetails("ComputerSystem", systemID); err != nil {
		return errors.PackError(err.ErrNo(), "error while trying to get ComputerSystem details: ", err.Error())
	}
	if len(computeData) == 1 {
		if inventoryData, err = connPool.GetAllMatchingDetails("FirmwareInventory", systemID); err != nil {
			return errors.PackError(err.ErrNo(), "error while trying to get compute details: ", err.Error())
		}
		for _, value := range inventoryData {
			if err = connPool.Delete("FirmwareInventory", value); err != nil {
				return errors.PackError(err.ErrNo(), "error while trying to delete compute details: ", err.Error())
			}
		}
		if inventoryData, err = connPool.GetAllMatchingDetails("SoftwareInventory", systemID); err != nil {
			return errors.PackError(err.ErrNo(), "error while trying to get compute details: ", err.Error())
		}
		for _, value := range inventoryData {
			if err = connPool.Delete("SoftwareInventory", value); err != nil {
				return errors.PackError(err.ErrNo(), "error while trying to delete compute details: ", err.Error())
			}
		}
	}

	//Delete All resources
	deleteKey := "*" + systemID + "*"
	if err = connPool.DeleteServer(deleteKey); err != nil {
		return errors.PackError(err.ErrNo(), "error while trying to delete compute system: ", err.Error())
	}
	if errs := deletefilteredkeys(key); errs != nil {
		return errors.PackError(errors.UndefinedErrorType, errs)
	}

	return nil
}

func deletefilteredkeys(key string) error {
	var sf Schema
	schemaFile, ioErr := ioutil.ReadFile(config.Data.SearchAndFilterSchemaPath)
	if ioErr != nil {
		return fmt.Errorf("fatal: error while trying to read search/filter schema json: %v", ioErr)
	}
	jsonErr := json.Unmarshal(schemaFile, &sf)
	if jsonErr != nil {
		return fmt.Errorf("fatal: error while trying to fetch search/filter schema json: %v", jsonErr)
	}
	conn, err := common.GetDBConnection(common.InMemory)
	if err != nil {
		return fmt.Errorf("error while trying to connecting to DB: %v", err)
	}
	for _, value := range sf.SearchKeys {
		for k := range value {
			delErr := conn.Del(k, key)
			if delErr != nil {
				if delErr.Error() != "no data with ID found" {
					return fmt.Errorf("error while deleting data: " + delErr.Error())
				}
			}
		}
	}

	delErr := conn.Del("UUID", key)
	if delErr != nil {
		if delErr.Error() != "no data with ID found" {
			return fmt.Errorf("error while deleting data: " + delErr.Error())
		}
	}
	delErr = conn.Del("PowerState", key)
	if delErr != nil {
		if delErr.Error() != "no data with ID found" {
			return fmt.Errorf("error while deleting data: " + delErr.Error())
		}
	}
	return nil
}

//DeleteSystem will delete the system from OnDisk
func DeleteSystem(key string) *errors.Error {
	connPool, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		return errors.PackError(err.ErrNo(), "error while trying to connecting to DB: ", err.Error())
	}

	// Check key present in the DB
	if _, err = connPool.Read("System", key); err != nil {
		return errors.PackError(err.ErrNo(), "error while trying to get compute details: ", err.Error())
	}

	deleteKey := "System:" + key
	//Delete All resources
	if err = connPool.DeleteServer(deleteKey); err != nil {
		return errors.PackError(err.ErrNo(), "error while trying to delete compute system: ", err.Error())
	}
	return nil
}

//GetTarget fetches the System(Target Device Credentials) table details
func GetTarget(deviceUUID string) (*Target, error) {
	var target Target
	conn, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		return nil, err
	}

	data, err := conn.Read("System", deviceUUID)
	if err != nil {
		return nil, fmt.Errorf("error while trying to get compute details: %v", err.Error())
	}

	if err := json.Unmarshal([]byte(data), &target); err != nil {
		return nil, err
	}

	return &target, nil
}

//SaveIndex is used to create a
func SaveIndex(searchForm map[string]interface{}, table, uuid string) error {
	conn, err := common.GetDBConnection(common.InMemory)
	if err != nil {
		return fmt.Errorf("error while trying to connecting to DB: %v", err)
	}
	log.Info("Creating index")
	searchForm["UUID"] = uuid
	if err := conn.CreateIndex(searchForm, table); err != nil {
		return fmt.Errorf("error while trying to index the document: %v", err)
	}

	return nil

}

//SavePluginData will saves plugin on disk
func SavePluginData(plugin Plugin) *errors.Error {

	conn, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		return err
	}
	const table string = "Plugin"
	if err := conn.Create(table, plugin.ID, plugin); err != nil {
		return errors.PackError(err.ErrNo(), "error while trying to save plugin data: ", err.Error())
	}

	return nil
}

// GetAllSystems extracts all the computer systems saved in ondisk
func GetAllSystems() ([]Target, *errors.Error) {
	conn, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		return nil, err
	}
	keys, err := conn.GetAllDetails("System")
	if err != nil {
		return nil, err
	}
	var targets []Target
	for _, key := range keys {
		var target Target
		targetData, err := conn.Read("System", key)
		if err != nil {
			return nil, err
		}
		if errs := json.Unmarshal([]byte(targetData), &target); errs != nil {
			return nil, errors.PackError(errors.UndefinedErrorType, errs)
		}
		targets = append(targets, target)

	}
	return targets, nil
}

//DeletePluginData will delete the plugin entry from the database based on the uuid
func DeletePluginData(key string) *errors.Error {
	conn, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		return err
	}
	if err = conn.Delete("Plugin", key); err != nil {
		return err
	}
	return nil
}

//DeleteManagersData will delete the Managers entry from the database based on the uuid
func DeleteManagersData(key string) *errors.Error {
	conn, err := common.GetDBConnection(common.InMemory)
	if err != nil {
		return err
	}
	if err = conn.Delete("Managers", key); err != nil {
		return err
	}
	return nil
}

//UpdateIndex is used for updating an existing index
func UpdateIndex(searchForm map[string]interface{}, table, uuid string) error {
	conn, err := common.GetDBConnection(common.InMemory)
	if err != nil {
		return fmt.Errorf("error while trying to connecting to DB: %v", err)
	}
	searchForm["UUID"] = uuid
	if err := conn.UpdateResourceIndex(searchForm, table); err != nil {
		return fmt.Errorf("error while trying to update index: %v", err)
	}

	return nil
}

//UpdateComputeSystem is used for updating ComputerSystem table
func UpdateComputeSystem(key string, computeData interface{}) error {
	conn, err := common.GetDBConnection(common.InMemory)
	if err != nil {
		return err
	}
	marshaledData, errs := json.Marshal(computeData)
	if errs != nil {
		return errs
	}
	if _, err := conn.Update("ComputerSystem", key, string(marshaledData)); err != nil {
		return err
	}
	return nil
}

//GetResourceDetails fetches a resource from database using key
func GetResourceDetails(key string) (string, *errors.Error) {
	conn, err := common.GetDBConnection(common.InMemory)
	if err != nil {
		return "", errors.PackError(err.ErrNo(), err)
	}
	resourceData, err := conn.GetResourceDetails(key)
	if err != nil {
		return "", errors.PackError(err.ErrNo(), "error while trying to get resource details: ", err.Error())
	}
	var resource string
	if errs := json.Unmarshal([]byte(resourceData), &resource); errs != nil {
		return "", errors.PackError(errors.UndefinedErrorType, errs)
	}
	return resource, nil
}

// GetString is used to retrive index values of type string
/* Inputs:
1. index is the index name to search with
2. match is the value to match with
*/
func GetString(index, match string) ([]string, error) {
	conn, dberr := common.GetDBConnection(common.InMemory)
	if dberr != nil {
		return nil, fmt.Errorf("error while trying to connecting to DB: %v", dberr.Error())
	}
	list, err := conn.GetString(index, 0, "*"+match+"*", false)
	if err != nil && err.Error() != "no data with ID found" {
		fmt.Println("error while getting the data", err)
		return []string{}, nil
	}
	return list, nil
}

// AddSystemOperationInfo connects to the persistencemgr and Add the system operation info to db
/* Inputs:
1.systemURI: computer system uri for which system operation is maintained
*/
func (system *SystemOperation) AddSystemOperationInfo(systemID string) *errors.Error {

	conn, err := common.GetDBConnection(common.InMemory)
	if err != nil {
		return err
	}
	//Create a header for data entry
	const table string = "SystemOperation"
	//Save data into Database
	if err = conn.AddResourceData(table, systemID, system); err != nil {
		return err
	}
	return nil
}

//GetSystemOperationInfo fetches the system opeation info for the given systemURI
/* Inputs:
1.systemURI: computer system uri for which system operation is maintained
*/
func GetSystemOperationInfo(systemURI string) (SystemOperation, *errors.Error) {
	var systemOperation SystemOperation

	conn, err := common.GetDBConnection(common.InMemory)
	if err != nil {
		return systemOperation, err
	}

	plugindata, err := conn.Read("SystemOperation", systemURI)
	if err != nil {
		return systemOperation, errors.PackError(err.ErrNo(), "error while trying to fetch system operation data: ", err.Error())
	}

	if err := json.Unmarshal([]byte(plugindata), &systemOperation); err != nil {
		return systemOperation, errors.PackError(errors.JSONUnmarshalFailed, err)
	}
	return systemOperation, nil
}

//DeleteSystemOperationInfo will delete the system operation entry from the database based on the systemURI
func DeleteSystemOperationInfo(systemURI string) *errors.Error {
	conn, err := common.GetDBConnection(common.InMemory)
	if err != nil {
		return err
	}
	if err = conn.Delete("SystemOperation", systemURI); err != nil {
		return err
	}
	return nil
}

// AddSystemResetInfo connects to the persistencemgr and Add the system reset info to db
/* Inputs:
1.systemURI: computer system uri for which system operation is maintained
2.resetType : reset type which is performed
*/
func AddSystemResetInfo(systemID, resetType string) *errors.Error {

	conn, err := common.GetDBConnection(common.InMemory)
	if err != nil {
		return err
	}
	//Create a header for data entry
	const table string = "SystemReset"
	//Save data into Database
	if err = conn.AddResourceData(table, systemID, map[string]string{
		"ResetType": resetType,
	}); err != nil {
		return err
	}
	return nil
}

//GetSystemResetInfo fetches the system reset info for the given systemURI
/* Inputs:
1.systemURI: computer system uri for which system operation is maintained
*/
func GetSystemResetInfo(systemURI string) (map[string]string, *errors.Error) {
	var resetInfo map[string]string

	conn, err := common.GetDBConnection(common.InMemory)
	if err != nil {
		return resetInfo, err
	}

	plugindata, err := conn.Read("SystemReset", systemURI)
	if err != nil {
		return resetInfo, errors.PackError(err.ErrNo(), "error while trying to fetch system reset data: ", err.Error())
	}

	if err := json.Unmarshal([]byte(plugindata), &resetInfo); err != nil {
		return resetInfo, errors.PackError(errors.JSONUnmarshalFailed, err)
	}
	return resetInfo, nil
}

//DeleteSystemResetInfo will delete the system reset entry from the database based on the systemURI
func DeleteSystemResetInfo(systemURI string) *errors.Error {
	conn, err := common.GetDBConnection(common.InMemory)
	if err != nil {
		return err
	}
	if err = conn.Delete("SystemReset", systemURI); err != nil {
		return err
	}
	return nil
}

// AddAggregationSource connects to the persistencemgr and Add the AggregationSource to db
/* Inputs:
1.req: AggregationSource info
2.aggregationSourceURI : uri of AggregationSource
*/
func AddAggregationSource(req AggregationSource, aggregationSourceURI string) *errors.Error {
	conn, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		return err
	}
	//Create a header for data entry
	const table string = "AggregationSource"
	//Save data into Database
	if err = conn.Create(table, aggregationSourceURI, req); err != nil {
		return err
	}
	return nil
}

// GetAggregationSourceInfo fetches the AggregationSource info for the given aggregationSourceURI
func GetAggregationSourceInfo(aggregationSourceURI string) (AggregationSource, *errors.Error) {
	var aggregationSource AggregationSource

	conn, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		return aggregationSource, err
	}

	data, err := conn.Read("AggregationSource", aggregationSourceURI)
	if err != nil {
		return aggregationSource, errors.PackError(err.ErrNo(), "error: while trying to fetch Aggregation Source data: ", err.Error())
	}

	if err := json.Unmarshal([]byte(data), &aggregationSource); err != nil {
		return aggregationSource, errors.PackError(errors.JSONUnmarshalFailed, err)
	}
	return aggregationSource, nil
}

// UpdateSystemData updates the bmc details
func UpdateSystemData(system SaveSystem, key string) *errors.Error {
	conn, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		return err
	}
	if _, err := conn.Update("System", key, system); err != nil {
		return err
	}
	return nil
}

// UpdatePluginData updates the plugin details
func UpdatePluginData(plugin Plugin, key string) *errors.Error {
	conn, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		return err
	}
	if _, err := conn.Update("Plugin", key, plugin); err != nil {
		return err
	}
	return nil
}

// UpdateAggregtionSource updates the aggregation details
func UpdateAggregtionSource(aggregationSource AggregationSource, key string) *errors.Error {
	conn, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		return err
	}
	if _, err := conn.Update("AggregationSource", key, aggregationSource); err != nil {
		return err
	}
	return nil
}

//GetAllMatchingDetails accepts the table name ,pattern and DB type and return all the keys which mathces the pattern
func GetAllMatchingDetails(table, pattern string, dbtype common.DbType) ([]string, *errors.Error) {
	conn, err := common.GetDBConnection(dbtype)
	if err != nil {
		return []string{}, err
	}
	return conn.GetAllMatchingDetails(table, pattern)
}

//DeleteAggregationSource will delete the AggregationSource entry from the database based on the aggregtionSourceURI
func DeleteAggregationSource(aggregtionSourceURI string) *errors.Error {
	conn, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		return err
	}
	if err = conn.Delete("AggregationSource", aggregtionSourceURI); err != nil {
		return err
	}
	return nil
}

//GetComputerSystem fetches computer system details by UUID from database
func GetComputerSystem(systemid string) (string, *errors.Error) {
	var system string
	conn, err := common.GetDBConnection(common.InMemory)
	if err != nil {
		// connection error
		return system, err
	}
	systemData, err := conn.Read("ComputerSystem", systemid)
	if err != nil {
		return "", errors.PackError(err.ErrNo(), "error while trying to get system details: ", err.Error())
	}
	if errs := json.Unmarshal([]byte(systemData), &system); errs != nil {
		return "", errors.PackError(errors.UndefinedErrorType, errs)
	}
	return system, nil
}

//CreateAggregate will create aggregate on disk
func CreateAggregate(aggregate Aggregate, aggregateURI string) *errors.Error {

	conn, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		return err
	}
	const table string = "Aggregate"
	if err := conn.Create(table, aggregateURI, aggregate); err != nil {
		return errors.PackError(err.ErrNo(), "error while trying to create aggregate: ", err.Error())
	}

	return nil
}

// GetAggregate fetches the aggregate info for the given aggregateURI
func GetAggregate(aggregateURI string) (Aggregate, *errors.Error) {
	var aggregate Aggregate

	conn, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		return aggregate, err
	}
	const table string = "Aggregate"
	data, err := conn.Read(table, aggregateURI)
	if err != nil {
		return aggregate, errors.PackError(err.ErrNo(), "error: while trying to fetch connection method data: ", err.Error())
	}

	if err := json.Unmarshal([]byte(data), &aggregate); err != nil {
		return aggregate, errors.PackError(errors.JSONUnmarshalFailed, err)
	}
	return aggregate, nil
}

//DeleteAggregate will delete the aggregate
func DeleteAggregate(key string) *errors.Error {
	conn, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		return err
	}
	const table string = "Aggregate"
	if err = conn.Delete(table, key); err != nil {
		return err
	}
	return nil
}

//GetAllKeysFromTable retrun all matching data give table name
func GetAllKeysFromTable(table string) ([]string, error) {
	conn, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		return nil, err
	}
	keysArray, err := conn.GetAllDetails(table)
	if err != nil {
		return nil, fmt.Errorf("error while trying to get all keys from table - %v: %v", table, err.Error())
	}
	return keysArray, nil
}

//AddElementsToAggregate add elements to the aggregate
func AddElementsToAggregate(aggregate Aggregate, aggregateURL string) *errors.Error {
	conn, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		return err
	}
	agg, err := GetAggregate(aggregateURL)
	if err != nil {
		return err
	}
	aggregate.Elements = append(aggregate.Elements, agg.Elements...)
	const table string = "Aggregate"
	if _, err := conn.Update(table, aggregateURL, aggregate); err != nil {
		return err
	}
	return nil
}

//RemoveElementsFromAggregate remove elements from an aggregate
func RemoveElementsFromAggregate(aggregate Aggregate, aggregateURL string) *errors.Error {
	conn, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		return err
	}
	agg, err := GetAggregate(aggregateURL)
	if err != nil {
		return err
	}
	aggregate.Elements = removeElements(aggregate.Elements, agg.Elements)

	const table string = "Aggregate"
	if _, err := conn.Update(table, aggregateURL, aggregate); err != nil {
		return err
	}
	return nil
}

func removeElements(requestElements, presentElements []string) []string {
	newElements := []string{}
	var present bool
	for _, presentElement := range presentElements {
		front := 0
		rear := len(requestElements) - 1
		for front <= rear {
			if requestElements[front] == presentElement || requestElements[rear] == presentElement {
				present = true
			}
			front++
			rear--
		}
		if !present {
			newElements = append(newElements, presentElement)
		}
	}
	return newElements
}

//AddConnectionMethod will add connection methods on disk
func AddConnectionMethod(connectionMethod ConnectionMethod, connectionMethodURI string) *errors.Error {
	conn, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		return err
	}
	const table string = "ConnectionMethod"
	if err := conn.Create(table, connectionMethodURI, connectionMethod); err != nil {
		return errors.PackError(err.ErrNo(), "error while trying to create aggregate: ", err.Error())
	}

	return nil
}

// GetConnectionMethod fetches the connection method info for the given connection method uri
func GetConnectionMethod(connectionMethodURI string) (ConnectionMethod, *errors.Error) {
	var connectionMethod ConnectionMethod

	conn, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		return connectionMethod, err
	}
	const table string = "ConnectionMethod"
	data, err := conn.Read(table, connectionMethodURI)
	if err != nil {
		return connectionMethod, errors.PackError(err.ErrNo(), "error: while trying to fetch connection method data: ", err.Error())
	}

	if err := json.Unmarshal([]byte(data), &connectionMethod); err != nil {
		return connectionMethod, errors.PackError(errors.JSONUnmarshalFailed, err)
	}
	return connectionMethod, nil
}

// Delete will delete the data from the provided db with the provided table and key data
func Delete(table, key string, dbtype common.DbType) *errors.Error {
	conn, err := common.GetDBConnection(dbtype)
	if err != nil {
		return err
	}
	if err = conn.Delete(table, key); err != nil {
		return err
	}
	return nil
}

// UpdateConnectionMethod updates the Connection Method details
func UpdateConnectionMethod(connectionMethod ConnectionMethod, key string) *errors.Error {
	conn, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		return err
	}
	if _, err := conn.Update("ConnectionMethod", key, connectionMethod); err != nil {
		return err
	}
	return nil
}

// CheckActiveRequest will check the DB to see whether there are any active requests for the given key
// It will return true if there is an active request or false if not
// It will also through an error if any DB connection issues arise
func CheckActiveRequest(key string) (bool, *errors.Error) {
	conn, err := common.GetDBConnection(common.InMemory)
	if err != nil {
		return false, errors.PackError(err.ErrNo(), "error: while trying to create connection with DB: ", err.Error())
	}
	_, err = conn.Read("ActiveAddBMCRequest", key)
	if err != nil {
		if errors.DBKeyNotFound == err.ErrNo() {
			return false, nil
		}
		return false, errors.PackError(err.ErrNo(), "error: while trying to fetch active connection details: ", err.Error())
	}
	return true, nil
}

// DeleteActiveRequest deletes the active request key from the DB, return error if any
func DeleteActiveRequest(key string) *errors.Error {
	conn, err := common.GetDBConnection(common.InMemory)
	if err != nil {
		return errors.PackError(err.ErrNo(), "error: while trying to create connection with DB: ", err.Error())
	}
	err = conn.Delete("ActiveAddBMCRequest", key)
	if err != nil {
		return errors.PackError(err.ErrNo(), "error: while trying to delete active connection details: ", err.Error())
	}
	return nil
}

//SavePluginManagerInfo will save plugin manager  data into the database
func SavePluginManagerInfo(body []byte, table string, key string) error {

	conn, err := common.GetDBConnection(common.InMemory)
	if err != nil {
		return fmt.Errorf("Unable to save the plugin data with SavePluginManagerInfo: %v", err.Error())
	}
	if err := conn.Create(table, key, string(body)); err != nil {
		return errors.PackError(err.ErrNo(), "Unable to save the plugin data with SavePluginManagerInfo:  duplicate UUID: ", err.Error())
	}

	return nil
}

// GetDeviceSubscriptions is to get subscription details of device
func GetDeviceSubscriptions(hostIP string) (*common.DeviceSubscription, error) {
	conn, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		return nil, err
	}
	devSubscription, gerr := conn.GetDeviceSubscription(common.DeviceSubscriptionIndex, hostIP+"*")
	if gerr != nil {
		return nil, fmt.Errorf("error while trying to get device subscription details: %v", gerr.Error())
	}
	devSub := strings.Split(devSubscription[0], "::")
	var deviceSubscription = &common.DeviceSubscription{
		EventHostIP:     devSub[0],
		Location:        devSub[1],
		OriginResources: getSliceFromString(devSub[2]),
	}
	return deviceSubscription, nil
}

// getSliceFromString is to convert the string to array
func getSliceFromString(sliceString string) []string {
	// redis will store array as string enclosed in "[]"(ex "[alert statuschange]")
	// to convert to an array remove "[" ,"]" and create a slice
	sliceString = strings.TrimSuffix(strings.TrimPrefix(sliceString, "["), "]")
	if len(sliceString) < 1 {
		return []string{}
	}
	return strings.Fields(sliceString)
}

// GetEventSubscriptions is for getting the event subscription details
func GetEventSubscriptions(key string) ([]string, error) {
	conn, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		return nil, err
	}
	subscriptions, gerr := conn.GetEvtSubscriptions(common.SubscriptionIndex, "*"+key+"*")
	if gerr != nil {
		return nil, fmt.Errorf("error while trying to get event subsciption details: %v", gerr.Error())
	}
	return subscriptions, nil
}

// UpdateDeviceSubscription is to update subscription details of device
func UpdateDeviceSubscription(devSubscription common.DeviceSubscription) error {
	conn, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		return err
	}
	if err := conn.UpdateDeviceSubscription(common.DeviceSubscriptionIndex, devSubscription.EventHostIP, devSubscription.Location, devSubscription.OriginResources); err != nil {
		return fmt.Errorf("error while trying to update subscription of device %v", err.Error())
	}
	return nil
}

// CheckMetricRequest will check the DB to see whether there are any active requests for the given key
// It will return true if there is an active request or false if not
// It will also through an error if any DB connection issues arise
func CheckMetricRequest(key string) (bool, *errors.Error) {
	conn, err := common.GetDBConnection(common.InMemory)
	if err != nil {
		return false, errors.PackError(err.ErrNo(), "error: while trying to create connection with DB: ", err.Error())
	}
	_, err = conn.Read("ActiveMetricRequest", key)
	if err != nil {
		if errors.DBKeyNotFound == err.ErrNo() {
			return false, nil
		}
		return false, errors.PackError(err.ErrNo(), "error: while trying to fetch active connection details: ", err.Error())
	}
	return true, nil
}

// DeleteMetricRequest deletes the active request key from the DB, return error if any
func DeleteMetricRequest(key string) *errors.Error {
	conn, err := common.GetDBConnection(common.InMemory)
	if err != nil {
		return errors.PackError(err.ErrNo(), "error: while trying to create connection with DB: ", err.Error())
	}
	err = conn.Delete("ActiveMetricRequest", key)
	if err != nil {
		return errors.PackError(err.ErrNo(), "error: while trying to delete active connection details: ", err.Error())
	}
	return nil
}
