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
	if Data.SouthBoundRequestTimeoutInSecs == 0 {
		t.Error("error: Data.SouthBoundRequestTimeoutInSecs is not initialized")
	}
	if Data.ServerRediscoveryBatchSize == 0 {
		t.Error("error: Data.ServerRediscoveryBatchSize is not initialized")
	}
	if Data.FirmwareVersion == "" {
		t.Error("error: Data.FirmwareVersion is not initialized")
	}
	if Data.RootServiceUUID == "" {
		t.Error("error: Data.RootServiceUUID is not initialized")
	}
	if Data.RegistryStorePath == "" {
		t.Error("error: Data.RegistryStorePath is not initialized")
	}
	if Data.LocalhostFQDN == "" {
		t.Error("error: Data.LocalhostFQDN is not initialized")
	}
	if len(Data.EnabledServices) == 0 {
		t.Error("error: Data.EnabledServices is not initialized")
	}
	if Data.DBConf == nil {
		t.Error("error: Data.DBConf is not initialized")
	}
	if Data.KeyCertConf == nil {
		t.Error("error: Data.KeyCertConf is not initialized")
	}
	if Data.AuthConf == nil {
		t.Error("error: Data.AuthConf is not initialized")
	}
	if Data.AuthConf.PasswordRules == nil {
		t.Error("error: Data.AuthConf.PasswordRules is not initialized")
	}
	if Data.APIGatewayConf == nil {
		t.Error("error: Data.APIGatewayConf is not initialized")
	}
	if Data.AddComputeSkipResources == nil {
		t.Error("error: Data.AddComputeSkipResources is not initialized")
	}
	if Data.URLTranslation == nil {
		t.Error("error: Data.URLTranslation is not initialized")
	}
	if Data.PluginStatusPolling == nil {
		t.Error("error: Data.PluginStatusPolling is not initialized")
	}
	if Data.ExecPriorityDelayConf == nil {
		t.Error("error: Data.ExecPriorityDelayConf is not initialized")
	}
	if Data.TLSConf == nil {
		t.Error("error: Data.TLSConf is not initialized")
	}
	if len(Data.SupportedPluginTypes) == 0 {
		t.Error("error: Data.SupportedPluginTypes is not initialized")
	}
	if len(Data.ConnectionMethodConf) == 0 {
		t.Error("error: Data.ConnectionMethodConf is not initialized")
	}
}
