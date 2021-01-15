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
// package ...
package config

import (
	"testing"
)

func TestSetUpMockConfig(t *testing.T) {
	if err := SetUpMockConfig(t); err != nil {
		t.Error("error: SetUpMockConfig failed with", err)
	}
	if Data.FirmwareVersion == "" {
		t.Error("error: Data.FirmwareVersion is not initialized")
	}
	if Data.RootServiceUUID == "" {
		t.Error("error: Data.RootServiceUUID is not initialized")
	}
	if Data.SessionTimeoutInMinutes == 0 {
		t.Error("error: Data.SessionTimeoutInMinutes is not initialized")
	}
	if Data.PluginConf == nil {
		t.Error("error: Data.PluginConf is not initialized")
	}
	if Data.KeyCertConf == nil {
		t.Error("error: Data.KeyCertConf is not initialized")
	}
	if Data.LoadBalancerConf == nil {
		t.Error("error: Data.LoadBalancerConf is not initialized")
	}
	if Data.EventConf == nil {
		t.Error("error: Data.EventConf is not initialized")
	}
	if Data.MessageBusConf == nil {
		t.Error("error: Data.MessageBusConfis not initialized")
	}
	if Data.URLTranslation == nil {
		t.Error("error: Data.URLTranslation is not initialized")
	}
	if Data.TLSConf == nil {
		t.Error("error: Data.TLSConf is not initialized")
	}
}
