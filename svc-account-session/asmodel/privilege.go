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

//Privileges struct definition
type Privileges struct {
	List []string
}

//GetPrivilegeRegistry retrives the privileges from database
func GetPrivilegeRegistry() (Privileges, *errors.Error) {
	var privileges Privileges
	conn, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		return privileges, err
	}
	privilegeregistry, err := conn.Read("registry", "assignedprivileges")
	if err != nil {
		return privileges, errors.PackError(err.ErrNo(), "error while trying to get Redfish predefined privileges: ", err.Error())
	}
	if jerr := json.Unmarshal([]byte(privilegeregistry), &privileges); jerr != nil {
		return privileges, errors.PackError(errors.UndefinedErrorType, jerr)
	}
	return privileges, nil
}

// Create method is to insert the privileges list to database
func (p *Privileges) Create() *errors.Error {

	conn, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		return err
	}
	if err = conn.Create("registry", "assignedprivileges", p); err != nil {
		return errors.PackError(err.ErrNo(), "error while trying to create privileges: ", err.Error())
	}

	return nil
}

//OEMPrivileges struct definition
type OEMPrivileges struct {
	List []string
}

//GetOEMPrivileges retrives the privileges from database
func GetOEMPrivileges() (OEMPrivileges, *errors.Error) {
	var oemPrivileges OEMPrivileges
	conn, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		return oemPrivileges, err
	}
	oemPrivilegeList, err := conn.Read("registry", "oemprivileges")
	if err != nil {
		return oemPrivileges, errors.PackError(err.ErrNo(), "error while trying to get OEM privileges: ", err.Error())
	}
	if jerr := json.Unmarshal([]byte(oemPrivilegeList), &oemPrivileges); jerr != nil {
		return oemPrivileges, errors.PackError(errors.UndefinedErrorType, jerr)
	}
	return oemPrivileges, nil
}

// Create method is to insert the oemprivileges list to database
func (p *OEMPrivileges) Create() *errors.Error {
	conn, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		return err
	}
	if err = conn.Create("registry", "oemprivileges", p); err != nil {
		return errors.PackError(err.ErrNo(), "error creating OEM privileges: ", err.Error())
	}
	return nil
}
