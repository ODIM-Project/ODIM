//(C) Copyright [2022] Hewlett Packard Enterprise Development LP
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

package licenses

import (
	"context"
	"net/http"

	"github.com/ODIM-Project/ODIM/lib-persistence-manager/persistencemgr"
	"github.com/ODIM-Project/ODIM/lib-rest-client/pmbhandle"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/lib-utilities/services"
	lcommon "github.com/ODIM-Project/ODIM/svc-licenses/lcommon"
	"github.com/ODIM-Project/ODIM/svc-licenses/model"
)

// ExternalInterface struct holds the structs to which hold function pointers to outboud calls
type ExternalInterface struct {
	External External
	DB       DB
}

// External struct holds the function pointers all outboud services
type External struct {
	ContactClient      func(context.Context, string, string, string, string, interface{}, map[string]string) (*http.Response, error)
	Auth               func(string, []string, []string) (response.RPC, error)
	DevicePassword     func([]byte) ([]byte, error)
	GetPluginData      func(string) (*model.Plugin, *errors.Error)
	ContactPlugin      func(model.PluginContactRequest, string) ([]byte, string, model.ResponseStatus, error)
	GetTarget          func(string) (*model.Target, *errors.Error)
	GetSessionUserName func(string) (string, error)
	GenericSave        func([]byte, string, string) error
}

// DB struct holds the function pointers to database operations
type DB struct {
	GetAllKeysFromTable func(string, persistencemgr.DbType) ([]string, error)
	GetResource         func(string, string, persistencemgr.DbType) (interface{}, *errors.Error)
}

// GetExternalInterface retrieves all the external connections update package functions uses
func GetExternalInterface() *ExternalInterface {
	return &ExternalInterface{
		External: External{
			ContactClient:      pmbhandle.ContactPlugin,
			Auth:               services.IsAuthorized,
			DevicePassword:     common.DecryptWithPrivateKey,
			GetPluginData:      lcommon.GetPluginData,
			ContactPlugin:      lcommon.ContactPlugin,
			GetTarget:          lcommon.GetTarget,
			GetSessionUserName: services.GetSessionUserName,
			GenericSave:        lcommon.GenericSave,
		},
		DB: DB{
			GetAllKeysFromTable: lcommon.GetAllKeysFromTable,
			GetResource:         lcommon.GetResource,
		},
	}
}
