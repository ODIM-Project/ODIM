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

//Package models ...
package models

import (
	"encoding/json"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	log "github.com/sirupsen/logrus"
)

//GetRegistryFile fetches a resource from database using table and key
func GetRegistryFile(Table, key string) ([]byte, *errors.Error) {
	conn, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		return nil, errors.PackError(err.ErrNo(), err)
	}
	resourceData, err := conn.Read(Table, key)
	log.Info("Table Name: " + Table + ", Key : " + key)
	if err != nil {
		return nil, errors.PackError(err.ErrNo(), "error while trying to get resource details: ", err.Error())
	}

	var resource string
	if errs := json.Unmarshal([]byte(resourceData), &resource); errs != nil {
		return nil, errors.PackError(errors.UndefinedErrorType, errs)
	}
	return []byte(resource), nil
}

//GetAllRegistryFileNamesFromDB return all key in given table
func GetAllRegistryFileNamesFromDB(table string) ([]string, *errors.Error) {

	conn, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		return nil, err
	}
	return conn.GetAllDetails(table)
}
