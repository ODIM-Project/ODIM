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

package update

import (
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	"github.com/ODIM-Project/ODIM/lib-utilities/services"
	"github.com/ODIM-Project/ODIM/svc-plugin-rest-client/pmbhandle"
	"github.com/ODIM-Project/ODIM/svc-update/umodel"
	"net/http"
)

//Device struct to define the response from plugin for UUID
type Device struct {
	ServerIP   string `json:"ServerIP"`
	Username   string `json:"Username"`
	DeviceUUID string `json:"device_UUID"`
}

// ExternalInterface struct holds the structs to which hold function pointers to outboud calls
type ExternalInterface struct {
	External External
	DB       DB
}

// External struct holds the function pointers all outboud services
type External struct {
	ContactClient  func(string, string, string, string, interface{}, map[string]string) (*http.Response, error)
	Auth           func(string, []string, []string) (int32, string)
	DevicePassword func([]byte) ([]byte, error)
}

type responseStatus struct {
	StatusCode    int32
	StatusMessage string
	MsgArgs       []interface{}
}

// DB struct holds the function pointers to database operations
type DB struct {
	GetAllKeysFromTable func(string) ([]string, error)
	GetResource         func(string, string) (string, *errors.Error)
}

// GetExternalInterface retrieves all the external connections update package functions uses
func GetExternalInterface() *ExternalInterface {
	return &ExternalInterface{
		External: External{
			ContactClient:  pmbhandle.ContactPlugin,
			Auth:           services.IsAuthorized,
			DevicePassword: common.DecryptWithPrivateKey,
		},
		DB: DB{
			GetAllKeysFromTable: umodel.GetAllKeysFromTable,
			GetResource:         umodel.GetResource,
		},
	}
}
