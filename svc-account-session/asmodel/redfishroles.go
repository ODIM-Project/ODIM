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

// Package asmodel ...
package asmodel

import (
	"encoding/json"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
)

//RedfishRoles struct definition
type RedfishRoles struct {
	List []string
}

//GetRedfishRoles retrives the privileges from database
func GetRedfishRoles() (RedfishRoles, *errors.Error) {
	var redfishRoles RedfishRoles
	conn, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		return redfishRoles, err
	}
	predefinedRoles, err := conn.Read("roles", "redfishdefined")
	if err != nil {
		return redfishRoles, errors.PackError(err.ErrNo(), "error while trying to get redfish roles: ", err.Error())
	}
	if jerr := json.Unmarshal([]byte(predefinedRoles), &redfishRoles); jerr != nil {
		return redfishRoles, errors.PackError(errors.UndefinedErrorType, jerr)
	}
	return redfishRoles, nil
}

// Create method is to insert the privileges list to database
func (r *RedfishRoles) Create() *errors.Error {

	conn, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		return err
	}
	if err = conn.Create("roles", "redfishdefined", r); err != nil {
		return errors.PackError(err.ErrNo(), "error while trying to create redfishroles: ", err.Error())
	}

	return nil
}
