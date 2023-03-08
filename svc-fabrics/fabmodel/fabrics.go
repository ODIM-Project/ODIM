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
// under the License

//Package fabmodel ...
package fabmodel

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	l "github.com/ODIM-Project/ODIM/lib-utilities/logs"
)

var (
	//GetDBConnectionFunc is pointer for common.GetDBConnection func
	GetDBConnectionFunc = common.GetDBConnection
)

// Plugin is the model for plugin information
type Plugin struct {
	IP                string
	Port              string
	Username          string
	Password          []byte
	ID                string
	PluginType        string
	PreferredAuthType string
}

// Fabric is the model to save fabric details
type Fabric struct {
	FabricUUID string
	PluginID   string
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

//GetAllFabricPluginDetails fetches all fabric plugin information from plugin table
func GetAllFabricPluginDetails(ctx context.Context) ([]string, error) {
	conn, err := GetDBConnectionFunc(common.OnDisk)
	if err != nil {
		return nil, err
	}
	keysArray, err := conn.GetAllDetails("Plugin")
	if err != nil {
		return nil, fmt.Errorf("error while trying to get data from table - Plugin : %v", err.Error())
	}
	l.LogWithFields(ctx).Debug("all fabric plugin details: ", keysArray)
	return keysArray, nil
}

// AddFabricData will add the fabric uuid and pluginid details into ondisk
func (fabric *Fabric) AddFabricData(fabuuid string) error {

	conn, err := GetDBConnectionFunc(common.OnDisk)
	if err != nil {
		return errors.PackError(errors.UndefinedErrorType, err)
	}
	//Create a header for data entry
	const table string = "Fabric"
	//Save data into Database
	if cerr := conn.Create(table, fabuuid, fabric); cerr != nil {
		if errors.DBKeyAlreadyExist != cerr.ErrNo() {
			return fmt.Errorf("error while trying to create new %v resource: %v", table, cerr.Error())
		}
		fmt.Printf("warning: skipped saving of duplicate data with key %v", fabuuid)
		return fmt.Errorf("warning: skipped saving of duplicate data with key %v", fabuuid)
	}
	return nil
}

// RemoveFabricData will remove the fabric uuid and pluginid details into ondisk
func (fabric *Fabric) RemoveFabricData(fabuuid string) error {
	conn, err := GetDBConnectionFunc(common.OnDisk)
	if err != nil {
		return errors.PackError(errors.UndefinedErrorType, err)
	}
	//Create a header for data entry
	const table string = "Fabric"
	//Save data into Database
	if cerr := conn.Delete(table, fabuuid); cerr != nil {
		return cerr
	}
	return nil
}

//GetManagingPluginIDForFabricID fetches the fabric details
func GetManagingPluginIDForFabricID(fabID string, ctx context.Context) (Fabric, error) {
	var fabric Fabric
	conn, err := GetDBConnectionFunc(common.OnDisk)
	if err != nil {
		return fabric, err
	}

	data, err := conn.Read("Fabric", fabID)
	if err != nil {
		return fabric, fmt.Errorf("error while trying to get fabric details: %v", err.Error())
	}

	if err := json.Unmarshal([]byte(data), &fabric); err != nil {
		return fabric, err
	}
	l.LogWithFields(ctx).Debug("Managing plugin for given fabricID: ", fabric)
	return fabric, nil
}

//GetAllTheFabrics fetches all the fabrics details
func GetAllTheFabrics(ctx context.Context) ([]Fabric, error) {
	conn, err := GetDBConnectionFunc(common.OnDisk)
	if err != nil {
		return nil, err
	}
	keys, err := conn.GetAllDetails("Fabric")
	l.LogWithFields(ctx).Debug("all keys from fabric database:",keys)
	if err != nil {
		return nil, err
	}
	var fabrics []Fabric
	for _, key := range keys {
		var fabric Fabric
		fabricData, err := conn.Read("Fabric", key)
		if err != nil {
			return nil, err
		}
		if errs := json.Unmarshal([]byte(fabricData), &fabric); errs != nil {
			return nil, errors.PackError(errors.UndefinedErrorType, errs)
		}
		fabrics = append(fabrics, fabric)

	}
	l.LogWithFields(ctx).Debug("All fabrics: ",fabrics)
	return fabrics, nil
}
