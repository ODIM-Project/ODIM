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
	redis_password = "Tr/72W+ie8WUdQFNQpa/BctrpbD9NOZiEkfl7OIBpOB0oLe14JWcfL7CVHs5WZ+" +
		"4qUtPV0cjSDnsDRCJAS/gXpNB8Mf2Cv/YI6uUZxlSh2GEJXfgbHaTxzFVVNyQ9x10vPi6hIVkIT9MzX9fE8" +
		"udLsi1X3R9yVjy6HxJ/VUH87nl4uP+d6MBhkK2/KCLYjeW/NZzH04SI8kYljJnuxWIZIM8GC4DiWpJD3A2r" +
		"jxs0oBQwpVfkzw/Wedf/qAoAOZq8eMK/ScxI6lmsGuoaPosvfelegtzsZr61SWVhcv72qg1OwxGmfcDgznjVGZ" +
		"huLSva40XOdlBEhUY074iiGL8a4Krun1AZlqIXY+irSzwNyhW9706pO2AxOsqDmQCJOD6QObYX3jzraM9o0/r" +
		"C9A9jiBgPbt5GnsiQfpv0rLcJP4ekbs2yfHsNh7P3Oj+jcjkOvENAtMmtTj6qLL5NxIa/HNTMU8Vf7g3HB4+" +
		"8ezS2r7xTO/5D6mr+/r4T8zlOXY9zjlWQYDP1SxXamymjexoof7KawIJY1NhU7cFbD1iUABd+AzFqJfQL0dbf" +
		"47xGD/eL1ZIL6ZJB2WNA5SVKRZz2JhW1jWLlmjqrmN1tvtsif4OVlLtoYqCT1YtiRlyAsdAt92C4Q/vv/" +
		"msUvhH9n8NfRRSwC47OcmNMEc92MHx1z0="
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
	createFile(t, RedisInMemoryPasswordFilePath, redis_password)
	createFile(t, RedisOnDiskPasswordFilePath, redis_password)

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
      			"MessageBusQueue": ["REDFISH-EVENTS-TOPIC"]
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
			if err := SetConfiguration(); (err != nil) != tt.wantErr {
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
			Data.MessageBusConf.MessageBusQueue = []string{"REDFISH-EVENTS-TOPIC"}
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
			if err := ValidateConfiguration(); (err != nil) != tt.wantErr {
				t.Errorf("TestValidateConfigurationForEventConf()  = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	os.Remove(sampleFileForTest)
}
