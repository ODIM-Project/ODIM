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

package account

import (
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	l "github.com/ODIM-Project/ODIM/lib-utilities/logs"
	"github.com/ODIM-Project/ODIM/svc-account-session/asmodel"
)

const (
	defaultAdminAccount = "admin"
)

var (
	// ConfigFilePath holds the value of odim config file path
	ConfigFilePath string
)

// ExternalInterface holds all the external connections account package functions uses
type ExternalInterface struct {
	CreateUser         func(asmodel.User) *errors.Error
	GetUserDetails     func(string) (asmodel.User, *errors.Error)
	GetRoleDetailsByID func(string) (asmodel.Role, *errors.Error)
	UpdateUserDetails  func(asmodel.User, asmodel.User) *errors.Error
}

// GetExternalInterface retrieves all the external connections account package functions uses
func GetExternalInterface() *ExternalInterface {
	return &ExternalInterface{
		CreateUser:         asmodel.CreateUser,
		GetUserDetails:     asmodel.GetUserDetails,
		GetRoleDetailsByID: asmodel.GetRoleDetailsByID,
		UpdateUserDetails:  asmodel.UpdateUserDetails,
	}
}

// TrackConfigFileChanges monitors the host file using fsnotfiy
func TrackConfigFileChanges(errChan chan error) {
	eventChan := make(chan interface{})
	format := config.Data.LogFormat
	go common.TrackConfigFileChanges(ConfigFilePath, eventChan, errChan)
	for {
		select {
		case info := <-eventChan:
			l.Log.Info(info) // new data arrives through eventChan channel
			if l.Log.Level != config.Data.LogLevel {
				l.Log.Info("Log level is updated, new log level is ", config.Data.LogLevel)
				l.Log.Logger.SetLevel(config.Data.LogLevel)
			}
			if format != config.Data.LogFormat {
				l.SetFormatter(config.Data.LogFormat)
				format = config.Data.LogFormat
				l.Log.Info("Log format is updated, new log format is ", config.Data.LogFormat)
			}
		case err := <-errChan:
			l.Log.Error(err)
		}
	}
}
