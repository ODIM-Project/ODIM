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
package common

import (
	"testing"

	"github.com/ODIM-Project/ODIM/lib-utilities/config"
)

func TestSetUpMockConfig(t *testing.T) {
	SetUpMockConfig()

	if config.Data.DBConf.InMemoryHost == "" {
		t.Error("InMemoryHost is not set")
	}
	if config.Data.DBConf.OnDiskHost == "" {
		t.Error("OnDiskHost is not set")
	}
	if config.Data.DBConf.OnDiskPort == "" {
		t.Error("OnDiskPort is not set")
	}
	if config.Data.DBConf.InMemoryPort == "" {
		t.Error("InMemoryPort is not set")
	}
	if config.Data.DBConf.Protocol == "" {
		t.Error("Protocol is not set")
	}
	if config.Data.DBConf.MaxActiveConns == 0 {
		t.Error("MaxActiveConns is not set")
	}
	if config.Data.DBConf.MaxIdleConns == 0 {
		t.Error("MaxIdleConns is not set")
	}
	if config.Data.AuthConf.SessionTimeOutInMins == 0 {
		t.Error("SessionTimeOutInMins is not set")
	}
	if config.Data.APIGatewayConf.Host == "" {
		t.Error("APIGatewayHost is not set")
	}
	if config.Data.APIGatewayConf.Port == "" {
		t.Error("APIGatewayPort is not set")
	}
	if len(config.Data.EnabledServices) == 0 {
		t.Error("EnabledServices is not set")
	}
	if config.Data.RootServiceUUID == "" {
		t.Error("RootServiceUUID is not set")
	}
	if config.Data.FirmwareVersion == "" {
		t.Error("Version is not set")
	}
}
