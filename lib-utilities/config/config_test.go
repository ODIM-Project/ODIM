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
	rsaPrivateKey     = `-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----`
	redisPassword = "F2n0YuRgavd/tYeanHI94cJR5r/C+FUaGQJBetQOxed1pLXxnWKAMmVLjs+jCBGq" +
		"C66YfEZ+DK5ZIg9QmuQGEoahwSVWC+Pa8hNrqIDgBYXP4cyyEAE0XE0j8amyf049aqhxxTYXfzov5Km4t/" +
		"Tzqru3CJ2CcUGzRmq1WfbfuMqx+tAZGw4UY1SW9IDoHwXaqsKld9uwiYq6lBqJpYzNGcCgrVyHwQg" +
		"hTrYlypQocsDdVY7/bFzg3amIHdStmzF+mvpNolhtkgrXeq1ov7stepdgpzOF39Fe5DDO+OG53wyR" +
		"4OMBAZ2NjX5LLQkhNEUpAA7GM8ajtOuJGecO506St2ASatcojJqRDHbIzNhzAxY3wtB0bx1S5TS1jl" +
		"kW1VTFXqFNjnKd3j7Q/YZXJY5a+zX/PhIZLnCp+yWY2/qU7s4BZjex8jNRikFTRzhqDGfKP1hFar8" +
		"qLr2D0FSRrDK4NtViUMUv5PaWygHtRk8e0fSnNhTSGv5kzr/fwEE4S/ayo5+9pqjgjr4iu+d6oRSq" +
		"2dQVQIfdm25Lqfw8RnmeveVVKuQk7xT/T0pcKmwfYuN0R4UjJ44BiBXgI9e/pgVTHzOrzHfLT5ekk1" +
		"eSx/fIuTMe38lxCD+L6mfhi4zsI/IkaQsjgvR70n5RlpsT8ndNLBtmfNS4NB3Ls2Cbp0AFyTY="
)

var cwdDir, _ = os.Getwd()

func createFile(t *testing.T, fName, fContent string) {
	if err := ioutil.WriteFile(fName, []byte(fContent), 0644); err != nil {
		t.Fatal("error :failed to create a sample file for tests:", err)
	}
}

func TestSetConfiguration(t *testing.T) {
	var (
		sampleConfigFile              = filepath.Join(cwdDir, "sample.json")
		sampleTestFile                = "/tmp/testFile.dat"
		sampleRSAPrivateKeyFile       = "/tmp/rsa.private"
		RedisInMemoryPasswordFilePath = "/tmp/redis_inmemory_password"
		RedisOnDiskPasswordFilePath   = "/tmp/redis_ondisk_password"
	)

	createFile(t, sampleTestFile, sampleFileContent)
	createFile(t, sampleRSAPrivateKeyFile, rsaPrivateKey)
	createFile(t, RedisInMemoryPasswordFilePath, redisPassword)
	createFile(t, RedisOnDiskPasswordFilePath, redisPassword)

	var sampleConfig = `{
        "RootServiceUUID": "a9762fb2-b9dd-4ce8-818a-af6833ba19f6",
        "LocalhostFQDN": "test.odim.local",
        "SearchAndFilterSchemaPath": "/tmp/testFile.dat",
        "RegistryStorePath": "/tmp",
        "KeyCertConf": {
                "RootCACertificatePath": "/tmp/testFile.dat",
                "RPCPrivateKeyPath": "/tmp/testFile.dat",
                "RPCCertificatePath": "/tmp/testFile.dat",
                "RSAPublicKeyPath": "/tmp/testFile.dat",
                "RSAPrivateKeyPath": "/tmp/rsa.private"
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
                "MaxActiveConns": 120,
                "RedisInMemoryPasswordFilePath": "/tmp/redis_inmemory_password",
                "RedisOnDiskPasswordFilePath": "/tmp/redis_ondisk_password"
        },
       	"MessageBusConf": {
      			"MessageBusConfigFilePath": "/tmp/testFile.dat",
	            "MessageBusType": "Kafka",
				"OdimControlMessageQueue":"ODIM-CONTROL-MESSAGES"
	      },
		"TaskQueueConf" : {
			"QueueSize": 20000,
			"DBCommitInterval": 1000,
			"RetryInterval": 5000
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
                "Managers",
		"CompositionService"
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
		],
    "EventConf": {
  		"DeliveryRetryAttempts" : 1,
  		"DeliveryRetryIntervalSeconds" : 1
    }
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
			if _, err := SetConfiguration(); (err != nil) != tt.wantErr {
				t.Errorf("SetConfiguration() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	os.Remove(sampleConfigFile)
	os.Remove(sampleConfigFile)
	os.Remove(sampleTestFile)
	os.Remove(sampleRSAPrivateKeyFile)
	os.Remove(RedisInMemoryPasswordFilePath)
	os.Remove(RedisOnDiskPasswordFilePath)
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
			Data.SearchAndFilterSchemaPath = sampleFileForTest
		case 4:
			Data.RegistryStorePath = cwdDir
		case 5:
			Data.EnabledServices = []string{"API"}
		case 6:
			Data.SupportedPluginTypes = []string{"plugin"}
		case 7:
			Data.DBConf = &DBConf{
				Protocol: "tcp",
			}
		}
		t.Run(tt.name, func(t *testing.T) {
			if _, err := ValidateConfiguration(); (err != nil) != tt.wantErr {
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
		{
			name:    "Invalid value for MessageBusConfigFilePath",
			wantErr: true,
		},
		{
			name:    "Invalid value for MessageBusType",
			wantErr: true,
		},
		{
			name:    "Invalid value for MessageBusQueue",
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
		case 10:
			Data.MessageBusConf = &MessageBusConf{
				MessageBusConfigFilePath: sampleFileForTest,
			}
		case 11:
			Data.MessageBusConf.MessageBusType = "Kafka"
		case 12:
			Data.MessageBusConf.OdimControlMessageQueue = "ODIM-CONTROL-MESSAGES"
		}
		t.Run(tt.name, func(t *testing.T) {
			if _, err := ValidateConfiguration(); (err != nil) != tt.wantErr {
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
		Data.TaskQueueConf = &TaskQueueConf{
			QueueSize:        1000,
			DBCommitInterval: 1000,
			RetryInterval:    5000,
		}
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
			if _, err := ValidateConfiguration(); (err != nil) != tt.wantErr {
				t.Errorf("TestValidateConfigurationGroup3() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	os.Remove(sampleFileForTest)
}

func TestValidateConfigurationForEventConf(t *testing.T) {
	sampleFileForTest := filepath.Join(cwdDir, sampleFileName)
	createFile(t, sampleFileForTest, sampleFileContent)
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "Empty event conf",
			wantErr: false,
		},
		{
			name:    "Zero value configured, setting to default",
			wantErr: false,
		},
	}
	for num, tt := range tests {
		switch num {
		case 0:
			Data.EventConf = &EventConf{}
		case 1:
			Data.EventConf = &EventConf{
				DeliveryRetryAttempts:        0,
				DeliveryRetryIntervalSeconds: 0,
			}
		}
		t.Run(tt.name, func(t *testing.T) {
			if _, err := ValidateConfiguration(); (err != nil) != tt.wantErr {
				t.Errorf("TestValidateConfigurationForEventConf()  = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	os.Remove(sampleFileForTest)
}
