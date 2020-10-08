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
	"net/http"

	"github.com/ODIM-Project/ODIM/lib-rest-client/pmbhandle"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	"github.com/ODIM-Project/ODIM/lib-utilities/services"
	"github.com/ODIM-Project/ODIM/svc-update/ucommon"
	"github.com/ODIM-Project/ODIM/svc-update/umodel"
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

// External struct holds the function pointers all outboud services
type External struct {
	ContactClient  func(string, string, string, string, interface{}, map[string]string) (*http.Response, error)
	Auth           func(string, []string, []string) (int32, string)
	DevicePassword func([]byte) ([]byte, error)
	GetPluginData  func(string) (umodel.Plugin, *errors.Error)
	ContactPlugin  func(ucommon.PluginContactRequest, string) ([]byte, string, ucommon.ResponseStatus, error)
	GetTarget      func(string) (*umodel.Target, *errors.Error)
}

type responseStatus struct {
	StatusCode    int32
	StatusMessage string
	MsgArgs       []interface{}
}

// DB struct holds the function pointers to database operations
type DB struct {
	GetAllKeysFromTable func(string, common.DbType) ([]string, error)
	GetResource         func(string, string, common.DbType) (string, *errors.Error)
}

// UpdateRequestBody struct defines the request body for update action
type UpdateRequestBody struct {
	ImageURI                         string                            `json:"ImageURI"`
	Password                         string                            `json:"Password,omitempty"`
	Targets                          []string                          `json:"Targets"`
	TransferProtocol                 string                            `json:"TransferProtocol,omitempty"`
	Username                         string                            `json:"Username,omitempty"`
	RedfishOperationApplyTimeSupport *RedfishOperationApplyTimeSupport `json:"@Redfish.OperationApplyTimeSupport,omitempty"`
}

// RedfishOperationApplyTimeSupport struct defines the apply time for the action in place
type RedfishOperationApplyTimeSupport struct {
	OdataType       string   `json:"@odata.type,omitempty"`
	SupportedValues []string `json:"SupportedValues,omitempty"`
}

// GetExternalInterface retrieves all the external connections update package functions uses
func GetExternalInterface() *ExternalInterface {
	return &ExternalInterface{
		External: External{
			ContactClient:  pmbhandle.ContactPlugin,
			Auth:           services.IsAuthorized,
			DevicePassword: common.DecryptWithPrivateKey,
			GetPluginData:  umodel.GetPluginData,
			ContactPlugin:  ucommon.ContactPlugin,
			GetTarget:      umodel.GetTarget,
		},
		DB: DB{
			GetAllKeysFromTable: umodel.GetAllKeysFromTable,
			GetResource:         umodel.GetResource,
		},
	}
}
