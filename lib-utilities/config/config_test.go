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

// Package config ...
package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

const (
	sampleFileContent = "A sample file for UT"
	sampleFileName    = "file_for_ut.test"
)

var cwdDir, _ = os.Getwd()

func createFile(t *testing.T, fName, fContent string) {
	if err := ioutil.WriteFile(fName, []byte(fContent), 0644); err != nil {
		t.Fatal("error :failed to create a sample file for tests:", err)
	}
}

func TestSetConfiguration(t *testing.T) {
	var (
		sampleConfigFile = filepath.Join(cwdDir, "sample.json")
		sampleTestFile   = "/tmp/testFile.dat"
	)

	createFile(t, sampleTestFile, sampleFileContent)

	var sampleConfig = `{
        "RootServiceUUID": "a9762fb2-b9dd-4ce8-818a-af6833ba19f6",
        "LocalhostFQDN": "test.odim.local",
        "MessageQueueConfigFilePath": "/tmp/testFile.dat",
        "SearchAndFilterSchemaPath": "/tmp/testFile.dat",
        "RegistryStorePath": "/tmp",
        "KeyCertConf": {
                "RootCACertificatePath": "/tmp/testFile.dat",
                "RPCPrivateKeyPath": "/tmp/testFile.dat",
                "RPCCertificatePath": "/tmp/testFile.dat",
                "RSAPublicKeyPath": "/tmp/testFile.dat",
                "RSAPrivateKeyPath": "/tmp/testFile.dat"
        },
        "APIGatewayConf": {
                "Host": "localhost",
                "Port": "9091",
                "PrivateKeyPath": "/tmp/testFile.dat",
                "CertificatePath": "/tmp/testFile.dat"
        },
        "DBConf": {
                "Protocol": "tcp",
                "InMemoryHost": "localhost",
                "InMemoryPort": "6379",
                "OnDiskHost": "localhost",
                "OnDiskPort": "6380",
                "MaxIdleConns": 10,
                "MaxActiveConns": 120
        },
        "FirmwareVersion": "1.0",
        "SouthBoundRequestTimeoutInSecs": 10,
        "ServerRediscoveryBatchSize": 30,
        "AuthConf": {
                "SessionTimeOutInMins": 30,
                "ExpiredSessionCleanUpTimeInMins": 15,
                "PasswordRules":{
                        "MinPasswordLength": 12,
                        "MaxPasswordLength": 16,
                        "AllowedSpecialCharcters": "~!@#$%^&*-+_|(){}:;<>,.?/"
                }
        },
        "AddComputeSkipResources": {
                "SkipResourceListUnderSystem": [
                        "Chassis",
                        "LogServices",
						"Managers"
                ],
				"SkipResourceListUnderManager": [
                        "Chassis",
                        "Systems",
                        "LogServices"
                ],
                "SkipResourceListUnderChassis": [
                        "Managers",
                        "Systems",
                        "Devices"
                ],
                "SkipResourceListUnderOthers": [
                        "Power",
                        "Thermal",
                        "SmartStorage",
                        "LogServices"
                ]
        },
        "URLTranslation": {
                "NorthBoundURL": {
                        "ODIM": "redfish"
                },
                "SouthBoundURL": {
                        "redfish": "ODIM"
                }
        },
        "PluginStatusPolling": {
                "PollingFrequencyInMins": 30,
                "MaxRetryAttempt": 3,
                "RetryIntervalInMins": 2,
                "ResponseTimeoutInSecs": 30,
                "StartUpResouceBatchSize": 10
        },
        "ExecPriorityDelayConf": {
                "MinResetPriority": 1,
                "MaxResetPriority": 10,
                "MaxResetDelayInSecs": 36000
        },
        "EnabledServices": [
                "SessionService",
                "AccountService",
                "AggregationService",
                "Systems",
                "Chassis",
                "TaskService",
                "EventService",
                "Fabrics",
                "Managers"
		],
		"SupportedPluginTypes" : ["Compute", "Fabric", "Storage"],
		"ConnectionMethodConf":[
		  {
				"ConnectionMethodType":"Redfish",
				"ConnectionMethodVariant":"Compute:GRF_v1.0.0"
			},
		  {
				"ConnectionMethodType":"Redfish",
				"ConnectionMethodVariant":"Storage:STG_v1.0.0"
			}
		]
}`

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "CONFIG_FILE_PATH not set",
			wantErr: true,
		},
		{
			name:    "Config file doesn't exist",
			wantErr: true,
		},
		{
			name:    "Empty config file",
			wantErr: true,
		},
		{
			name:    "Valid config file",
			wantErr: false,
		},
	}
	for num, tt := range tests {
		switch num {
		case 0:
			os.Unsetenv("CONFIG_FILE_PATH")
		case 1:
			cfgFilePath := filepath.Join(cwdDir, "")
			os.Setenv("CONFIG_FILE_PATH", cfgFilePath)
		case 2:
			createFile(t, sampleConfigFile, "")
			cfgFilePath := sampleConfigFile
			os.Setenv("CONFIG_FILE_PATH", cfgFilePath)
		case 3:
			createFile(t, sampleConfigFile, sampleConfig)
			cfgFilePath := sampleConfigFile
			os.Setenv("CONFIG_FILE_PATH", cfgFilePath)
		}
		t.Run(tt.name, func(t *testing.T) {
			if err := SetConfiguration(); (err != nil) != tt.wantErr {
				t.Errorf("SetConfiguration() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	os.Remove(sampleConfigFile)
	os.Remove(sampleConfigFile)
	os.Remove(sampleTestFile)
}

func TestValidateConfigurationGroup1(t *testing.T) {
	sampleFileForTest := filepath.Join(cwdDir, sampleFileName)
	createFile(t, sampleFileForTest, sampleFileContent)

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "0) negative case",
			wantErr: true,
		},
		{
			name:    "Invalid value for config data",
			wantErr: true,
		},
		{
			name:    "Invalid value for RootServiceUUID",
			wantErr: true,
		},
		{
			name:    "Invalid value for LocalhostFQDN",
			wantErr: true,
		},
		{
			name:    "Invalid value for MessageQueueConfigFilePath",
			wantErr: true,
		},
		{
			name:    "Invalid value for SearchAndFilterSchemaPath",
			wantErr: true,
		},
		{
			name:    "Invalid value for RegistryStorePath",
			wantErr: true,
		},
		{
			name:    "Invalid value for EnabledServices",
			wantErr: true,
		},
		{
			name:    "Invalid value for Plugin types",
			wantErr: true,
		},
		{
			name:    "Invalid value for DBConf",
			wantErr: true,
		},
	}
	for num, tt := range tests {
		switch num {
		case 0:
			Data = configModel{}
		case 1:
			Data.FirmwareVersion = "someVal"
			Data.RootServiceUUID = "someVal"
		case 2:
			Data.RootServiceUUID = "a9762fb2-b9dd-4ce8-818a-af6833ba19f6"
			Data.SouthBoundRequestTimeoutInSecs = 10
			Data.LocalhostFQDN = "test.odim.local"
		case 3:
			Data.MessageQueueConfigFilePath = sampleFileForTest
		case 4:
			Data.SearchAndFilterSchemaPath = sampleFileForTest
		case 5:
			Data.RegistryStorePath = cwdDir
		case 6:
			Data.EnabledServices = []string{"API"}
		case 7:
			Data.SupportedPluginTypes = []string{"plugin"}
		case 8:
			Data.DBConf = &DBConf{
				Protocol: "tcp",
			}
		}
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateConfiguration(); (err != nil) != tt.wantErr {
				t.Errorf("TestValidateConfigurationGroup1() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	os.Remove(sampleFileForTest)
}

func TestValidateConfigurationGroup2(t *testing.T) {
	sampleFileForTest := filepath.Join(cwdDir, sampleFileName)
	createFile(t, sampleFileForTest, sampleFileContent)

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "Invalid value for InMemoryHost",
			wantErr: true,
		},
		{
			name:    "Invalid value for InMemoryPort",
			wantErr: true,
		},
		{
			name:    "Invalid value for OnDiskHost",
			wantErr: true,
		},
		{
			name:    "Invalid value for OnDiskPort",
			wantErr: true,
		},
		{
			name:    "Invalid value for MaxActiveConns and MaxIdleConns",
			wantErr: true,
		},
		{
			name:    "Invalid value for KeyCertConf",
			wantErr: true,
		},
		{
			name:    "Invalid value for RPCPrivateKeyPath",
			wantErr: true,
		},
		{
			name:    "Invalid value for RPCCertificatePath",
			wantErr: true,
		},
		{
			name:    "Invalid value for RSAPublicKeyPath",
			wantErr: true,
		},
		{
			name:    "Invalid value for RSAPrivateKeyPath",
			wantErr: true,
		},
	}
	for num, tt := range tests {
		switch num {
		case 0:
			Data.DBConf.InMemoryHost = "localhost"
		case 1:
			Data.DBConf.InMemoryPort = "someport"
		case 2:
			Data.DBConf.OnDiskHost = "localhost"
		case 3:
			Data.DBConf.OnDiskPort = "someport"
		case 4:
			Data.DBConf.MaxActiveConns = 120
			Data.DBConf.MaxIdleConns = 10
		case 5:
			Data.KeyCertConf = &KeyCertConf{
				RootCACertificatePath: sampleFileForTest,
			}
		case 6:
			Data.KeyCertConf.RPCPrivateKeyPath = sampleFileForTest
		case 7:
			Data.KeyCertConf.RPCCertificatePath = sampleFileForTest
		case 8:
			Data.KeyCertConf.RSAPublicKeyPath = sampleFileForTest
		case 9:
			Data.KeyCertConf.RSAPrivateKeyPath = sampleFileForTest
		}
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateConfiguration(); (err != nil) != tt.wantErr {
				t.Errorf("TestValidateConfigurationGroup2() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	os.Remove(sampleFileForTest)
}

func TestValidateConfigurationGroup3(t *testing.T) {
	sampleFileForTest := filepath.Join(cwdDir, sampleFileName)
	createFile(t, sampleFileForTest, sampleFileContent)
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "Invalid value for AuthConf",
			wantErr: true,
		},
		{
			name:    "Invalid value for ExpiredSessionCleanUpTimeInMins",
			wantErr: true,
		},
		{
			name:    "Invalid value for PasswordRules",
			wantErr: true,
		},
		{
			name:    "Invalid value for APIGatewayConf",
			wantErr: true,
		},
		{
			name:    "Invalid value for APIGatewayConf.Port",
			wantErr: true,
		},
		{
			name:    "Invalid value for APIGatewayConf.PrivateKeyPath",
			wantErr: true,
		},
		{
			name:    "Invalid value for ConnectionMethodConf",
			wantErr: true,
		},
		{
			name:    "Invalid value for APIGatewayConf.CertificatePath",
			wantErr: false,
		},
		{
			name:    "Invalid value for AddComputeSkipResources",
			wantErr: false,
		},
		{
			name:    "Invalid value for AddComputeSkipResources.SkipResourceListUnderChassis",
			wantErr: false,
		},
		{
			name:    "Invalid value for AddComputeSkipResources.SkipResourceListUnderOthers",
			wantErr: false,
		},
		{
			name:    "Invalid value for URLTranslation",
			wantErr: false,
		},
		{
			name:    "Invalid value for AddComputeSkipResources.SkipResourceListUnderManager",
			wantErr: false,
		},
	}
	for num, tt := range tests {
		switch num {
		case 0:
			Data.AuthConf = &AuthConf{
				SessionTimeOutInMins: 30,
			}
		case 1:
			Data.AuthConf.ExpiredSessionCleanUpTimeInMins = 30
		case 2:
			Data.AuthConf.PasswordRules = &PasswordRules{
				MinPasswordLength:       0,
				MaxPasswordLength:       0,
				AllowedSpecialCharcters: "",
			}
		case 3:
			Data.APIGatewayConf = &APIGatewayConf{
				Host: "localhost",
			}
		case 4:
			Data.APIGatewayConf.Port = "someport"
		case 5:
			Data.APIGatewayConf.PrivateKeyPath = sampleFileForTest
		case 6:
			Data.ConnectionMethodConf = []ConnectionMethodConf{
				{
					ConnectionMethodType:    "Redfish",
					ConnectionMethodVariant: "GRF_v1.0.0",
				},
			}
		case 7:
			Data.APIGatewayConf.CertificatePath = sampleFileForTest
		case 8:
			Data.AddComputeSkipResources = &AddComputeSkipResources{
				SkipResourceListUnderSystem: []string{"Chassis", "LogServices", "Manager"},
			}
		case 9:
			Data.AddComputeSkipResources.SkipResourceListUnderChassis = []string{"Managers", "Systems", "Devices"}
		case 10:
			Data.AddComputeSkipResources.SkipResourceListUnderOthers = []string{"Power", "Thermal", "SmartStorage"}
		case 11:
			Data.URLTranslation = &URLTranslation{
				NorthBoundURL: map[string]string{},
				SouthBoundURL: map[string]string{},
			}
			Data.PluginStatusPolling = &PluginStatusPolling{}
		case 12:
			Data.AddComputeSkipResources.SkipResourceListUnderManager = []string{"Chassis", "Systems", "LogServices"}
		}
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateConfiguration(); (err != nil) != tt.wantErr {
				t.Errorf("TestValidateConfigurationGroup3() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	os.Remove(sampleFileForTest)
}
