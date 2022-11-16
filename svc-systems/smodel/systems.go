//(C) Copyright [2020] Hewlett Packard Enterprise Development LP
//(C) Copyright 2020 Intel Corporation
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

//Package smodel ....
package smodel

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ODIM-Project/ODIM/lib-persistence-manager/persistencemgr"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	l "github.com/ODIM-Project/ODIM/lib-utilities/logs"

	"github.com/gomodule/redigo/redis"
)

//Target is for sending the requst to south bound/plugin
type Target struct {
	ManagerAddress string `json:"ManagerAddress"`
	Password       []byte `json:"Password"`
	UserName       string `json:"UserName"`
	PostBody       []byte `json:"PostBody"`
	DeviceUUID     string `json:"DeviceUUID"`
	PluginID       string `json:"PluginID"`
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

var (
	// GetDBConnectionFunc  function pointer for the common.GetDBConnection
	GetDBConnectionFunc = common.GetDBConnection
	// JSONUnmarshalFunc  function pointer for the json.Unmarshal
	JSONUnmarshalFunc = json.Unmarshal
	// scanFunc function pointer for the scan
	scanFunc = scan
)

// Volume is for sending a volume's request to south bound
type Volume struct {
	RAIDType           string `json:"RAIDType"`
	Links              *Links `json:"Links"`
	OperationApplyTime string `json:"@Redfish.OperationApplyTime"`
	DisplayName        string `json:"DisplayName,omitempty"`
	WriteCachePolicy   string `json:"WriteCachePolicy,omitempty"`
	ReadCachePolicy    string `json:"ReadCachePolicy,omitempty"`
	IOPerfModeEnabled  bool   `json:"IOPerfModeEnabled,omitempty"`
}

// Links contains Drives resoruces info
type Links struct {
	Drives               []OdataIDLink `json:"Drives"`
	DedicatedSpareDrives []OdataIDLink `json:"DedicatedSpareDrives,omitempty"`
}

// OdataIDLink contains link to a resource
type OdataIDLink struct {
	OdataID string `json:"@odata.id"`
}

//GetSystemByUUID fetches computer system details by UUID from database
func GetSystemByUUID(systemUUID string) (string, *errors.Error) {
	var system string
	conn, err := GetDBConnectionFunc(common.InMemory)
	if err != nil {
		// connection error
		return system, err
	}
	systemData, err := conn.Read("ComputerSystem", systemUUID)
	if err != nil {
		return "", errors.PackError(err.ErrNo(), "error while trying to get system details: ", err.Error())
	}
	if errs := JSONUnmarshalFunc([]byte(systemData), &system); errs != nil {
		return "", errors.PackError(errors.UndefinedErrorType, errs)
	}
	return system, nil
}

//GetResource fetches a resource from database using table and key
func GetResource(Table, key string) (string, *errors.Error) {
	conn, err := GetDBConnectionFunc(common.InMemory)
	if err != nil {
		return "", err
	}
	resourceData, err := conn.Read(Table, key)
	if err != nil {
		return "", errors.PackError(err.ErrNo(), "error while trying to get resource details: ", err.Error())
	}
	var resource string
	if errs := JSONUnmarshalFunc([]byte(resourceData), &resource); errs != nil {
		return "", errors.PackError(errors.UndefinedErrorType, errs)
	}
	return resource, nil
}

// Find fetches a resource from database using table and key and store the data to an interface
func Find(table, key string, r interface{}) *errors.Error {
	conn, err := GetDBConnectionFunc(common.InMemory)
	if err != nil {
		return err
	}
	resourceData, err := conn.Read(table, key)
	if err != nil {
		return errors.PackError(err.ErrNo(), "error while trying to get resource details: ", err.Error())
	}

	var resourceAsString string
	if errs := JSONUnmarshalFunc([]byte(resourceData), &resourceAsString); errs != nil {
		return errors.PackError(errors.UndefinedErrorType, errs)
	}

	if errs := json.Unmarshal([]byte(resourceAsString), r); errs != nil {
		return errors.PackError(errors.UndefinedErrorType, errs)
	}
	return nil
}

// FindAll ...
func FindAll(table, key string) ([][]byte, error) {
	cp, cpErr := GetDBConnectionFunc(common.OnDisk)
	if cpErr != nil {
		return nil, cpErr
	}

	affectedKeys, err := scanFunc(cp, strings.Join([]string{table, key}, ":"))
	if err != nil {
		return nil, err
	}

	if len(affectedKeys) == 0 {
		return [][]byte{}, nil
	}

	conn := cp.ReadPool.Get()
	defer conn.Close()

	return redis.ByteSlices(conn.Do("MGET", affectedKeys...))
}

func scan(cp *persistencemgr.ConnPool, key string) ([]interface{}, error) {
	conn := cp.ReadPool.Get()
	defer conn.Close()

	var (
		cursor int64
		items  []interface{}
	)

	results := make([]interface{}, 0)

	for {
		values, err := redis.Values(conn.Do("SCAN", cursor, "MATCH", key))
		if err != nil {
			return nil, err
		}

		_, err = redis.Scan(values, &cursor, &items)
		if err != nil {
			return nil, err
		}
		results = append(results, items...)

		if cursor == 0 {
			break
		}
	}

	return results, nil
}

//GetAllKeysFromTable fetches all keys in a given table
func GetAllKeysFromTable(table string) ([]string, error) {
	conn, err := GetDBConnectionFunc(common.InMemory)
	if err != nil {
		return nil, err
	}
	keysArray, err := conn.GetAllDetails(table)
	if err != nil {
		return nil, fmt.Errorf("error while trying to get all keys from table - %v: %v", table, err.Error())
	}
	return keysArray, nil
}

//GetPluginData will fetch plugin details
func GetPluginData(pluginID string) (Plugin, *errors.Error) {
	var plugin Plugin

	conn, err := GetDBConnectionFunc(common.OnDisk)
	if err != nil {
		return plugin, err
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

//GetTarget fetches the System(Target Device Credentials) table details
func GetTarget(deviceUUID string) (*Target, *errors.Error) {
	var target Target
	conn, err := GetDBConnectionFunc(common.OnDisk)
	if err != nil {
		return nil, err
	}
	data, err := conn.Read("System", deviceUUID)
	if err != nil {
		return nil, errors.PackError(err.ErrNo(), "error while trying to get compute details: ", err.Error())
	}
	if errs := json.Unmarshal([]byte(data), &target); errs != nil {
		return nil, errors.PackError(errors.UndefinedErrorType, errs)
	}
	return &target, nil
}

//GenericSave will save any resource data into the database
func GenericSave(body []byte, table string, key string) error {
	connPool, err := GetDBConnectionFunc(common.InMemory)
	if err != nil {
		return fmt.Errorf("error while trying to connecting to DB: %v", err.Error())
	}
	if err = connPool.AddResourceData(table, key, string(body)); err != nil {
		if errors.DBKeyAlreadyExist == err.ErrNo() {
			return fmt.Errorf("error while trying to create new %v resource: %v", table, err.Error())
		}
		l.Log.Warn("Skipped saving of duplicate data with key " + key)
	}
	return nil
}

// GetStorageList is used to storage list of capacity
/*
1.index name to search with
2. condition is the value for condition operation
3. match is the search for list float type
*/
func GetStorageList(index, condition string, match float64, regexFlag bool) ([]string, error) {
	conn, dberr := GetDBConnectionFunc(common.InMemory)
	if dberr != nil {
		return nil, fmt.Errorf("error while trying to connecting to DB: %v", dberr.Error())
	}
	list, err := conn.GetStorageList(index, 0, match, condition, regexFlag)
	if err != nil {
		fmt.Println("error while getting the data", err)
		return []string{}, nil
	}
	return list, nil
}

// GetString is used to retrive index values of type string
/* Inputs:
1. index is the index name to search with
2. match is the value to match with
*/
func GetString(index, match string, regexFlag bool) ([]string, error) {
	conn, dberr := GetDBConnectionFunc(common.InMemory)
	if dberr != nil {
		return nil, fmt.Errorf("error while trying to connecting to DB: %v", dberr.Error())
	}
	list, err := conn.GetString(index, 0, "*"+match+"*", regexFlag)
	if err != nil {
		fmt.Println("error while getting the data", err)
		return []string{}, nil
	}
	return list, nil
}

// GetRange is used to retrive index values of type string
/* Inputs:
1. index is the index name to search with
2. min is the minimum  value passed
3. max is the max value passed
*/
func GetRange(index string, min, max int, regexFlag bool) ([]string, error) {
	conn, dberr := GetDBConnectionFunc(common.InMemory)
	if dberr != nil {
		return nil, fmt.Errorf("error while trying to connecting to DB: %v", dberr.Error())
	}
	list, err := conn.GetRange(index, min, max, regexFlag)
	if err != nil {
		fmt.Println("error while getting the data", err)
		return []string{}, nil
	}
	return list, nil
}

// AddSystemResetInfo connects to the persistencemgr and Add the system reset info to db
/* Inputs:
1.systemURI: computer system uri for which system operation is maintained
2.resetType : reset type which is performed
*/
func AddSystemResetInfo(systemID, resetType string) *errors.Error {

	conn, err := GetDBConnectionFunc(common.InMemory)
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

	conn, err := GetDBConnectionFunc(common.InMemory)
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

//DeleteVolume will delete the volume from InMemory
func DeleteVolume(key string) *errors.Error {
	connPool, err := GetDBConnectionFunc(common.InMemory)
	if err != nil {
		return errors.PackError(err.ErrNo(), "error while trying to connecting to DB: ", err.Error())
	}

	// Check key present in the DB
	if _, err = connPool.Read("Volumes", key); err != nil {
		return errors.PackError(err.ErrNo(), "error while trying to get voulme details: ", err.Error())
	}

	//Delete volume
	if err = connPool.Delete("Volumes", key); err != nil {
		return errors.PackError(err.ErrNo(), "error while trying to delete volume: ", err.Error())
	}
	return nil
}
