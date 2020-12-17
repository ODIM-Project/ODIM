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

// Package role ...
package role

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"

	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-account-session/asmodel"
)

// Status to handle the error code and message
type Status struct {
	Code    int
	Message string
}

//ValidateAssignedPrivileges  provides functionality which verifies user provided privileges
// with configured privileges
// It accepts  user provided privileges for the role as request and returns  status and error
func validateAssignedPrivileges(assignedPrivileges []string) (*Status, []interface{}, error) {
	//Get privilege registry from database
	privilegeRegistry, err := asmodel.GetPrivilegeRegistry()
	if err != nil {
		log.Error("Unable to get Privileges: " + err.Error())
		return &Status{
			Code:    http.StatusInternalServerError,
			Message: response.InternalError,
		}, []interface{}{}, fmt.Errorf("error while getting Privileges: %v", err)
	}

	//Check if requested privileges are in privilege registry
	if len(assignedPrivileges) != 0 {
		for _, userPrivilege := range assignedPrivileges {
			flag := false
			for _, redfishPrivilege := range privilegeRegistry.List {
				if userPrivilege == redfishPrivilege {
					flag = true
					break
				}
			}
			if !flag {
				log.Error("Requested Redfish predefined privilege is not correct")
				return &Status{Code: http.StatusBadRequest, Message: response.PropertyValueNotInList}, []interface{}{userPrivilege, "AssignedPrivileges"}, fmt.Errorf("Requested Redfish predefined privilege is not correct")
			}
		}
	}
	return nil, []interface{}{}, nil
}

//ValidateOEMPrivileges provides functionality which verifies user provided OEMprivileges
// with configured OEMprivileges
// It accepts  user provided OEMprivileges for the role as request and returns Status and error as response
func validateOEMPrivileges(oemPrivileges []string) (*Status, []interface{}, error) {
	//Get OEM privileges from database
	oemPrivilegeRegistry, err := asmodel.GetOEMPrivileges()
	if err != nil {
		log.Error("error getting OEM Privileges: " + err.Error())
		return &Status{Code: http.StatusInternalServerError, Message: response.InternalError}, []interface{}{}, fmt.Errorf("error getting OEM Privileges: %v", err)
	}

	//Check if requested privileges are OEM privileges
	if len(oemPrivileges) != 0 {
		for _, userPrivilege := range oemPrivileges {
			flag := false
			for _, oemPrivilege := range oemPrivilegeRegistry.List {
				if userPrivilege == oemPrivilege {
					flag = true
					break
				}
			}
			if !flag {
				log.Error("Requested OEM privilege is not correct")
				return &Status{Code: http.StatusBadRequest, Message: response.PropertyValueNotInList}, []interface{}{userPrivilege, "OemPrivileges"}, fmt.Errorf("Requested OEM privilege is not correct")
			}
		}
	}
	return nil, []interface{}{}, nil
}

//CheckForPrivilege checks given  privilege  in session object
//It accepts session object  and privilege as request and returns  Status and error as response
func checkForPrivilege(session *asmodel.Session, privilege string) (*Status, error) {
	//Get session object
	//check if user has ConfigureUsers privilege
	//	log.Println(session)
	if !session.Privileges[privilege] {
		log.Error("InsufficientPrivilege")
		return &Status{Code: http.StatusForbidden, Message: response.InsufficientPrivilege}, fmt.Errorf("InsufficientPrivilege")
	}
	return nil, nil
}
