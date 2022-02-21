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

// Role struct definition
type Role struct {
	ID                 string   `json:"RoleId"`
	IsPredefined       bool     `json:"IsPredefined"`
	AssignedPrivileges []string `json:"AssignedPrivileges"`
	OEMPrivileges      []string `json:"OemPrivileges"`
}

// Create method is to insert the role details into database
func (r *Role) Create() *errors.Error {

	conn, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		return err
	}
	const table string = "role"
	if err := conn.Create(table, r.ID, r); err != nil {
		return errors.PackError(err.ErrNo(), "error while trying to create role: ", err.Error())
	}

	return nil
}

// GetRoleDetailsByID retrives the privileges for a role from database
func GetRoleDetailsByID(roleID string) (Role, *errors.Error) {
	var role Role
	conn, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		return role, err
	}
	roleData, err := conn.Read("role", roleID)
	if err != nil {
		return role, errors.PackError(err.ErrNo(), "error while trying to get role details: ", err.Error())
	}
	if jerr := json.Unmarshal([]byte(roleData), &role); jerr != nil {
		return role, errors.PackError(errors.UndefinedErrorType, jerr)
	}
	return role, nil
}

//UpdateRoleDetails will modify the current details to given changes
func (r *Role) UpdateRoleDetails() *errors.Error {

	conn, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		return err
	}
	if _, err = conn.Update("role", r.ID, r); err != nil {
		return err
	}
	return nil
}

//GetAllRoles gets all the roles from the db
func GetAllRoles() ([]Role, *errors.Error) {
	conn, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		return nil, err
	}
	keys, err := conn.GetAllDetails("role")
	if err != nil {
		return nil, err
	}
	var roles []Role
	//roles := make(map[string]Role)
	for _, key := range keys {
		var role Role
		roledata, err := conn.Read("role", key)
		if err != nil {
			return nil, err
		}
		if jerr := json.Unmarshal([]byte(roledata), &role); jerr != nil {
			return nil, errors.PackError(errors.UndefinedErrorType, jerr)
		}
		roles = append(roles, role)

	}
	return roles, nil
}

//Delete will delete the role entry from the database based on the uuid
func (r *Role) Delete() *errors.Error {
	conn, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		return err
	}
	if err = conn.Delete("role", r.ID); err != nil {
		return err
	}
	return nil
}
