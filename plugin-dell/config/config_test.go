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
	SetUpMockConfig(t)
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
                "StartUpResourceBatchSize": 10
        },
		"PluginTasksConf" : {
			"MonitorPluginTasksFrequencyInMins": 60
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
				"ConnectionMethodVariant":"Compute:GRF_v2.0.0"
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
			name:    "PLUGIN_CONFIG_FILE_PATH not set",
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
			os.Unsetenv("PLUGIN_CONFIG_FILE_PATH")
		case 1:
			cfgFilePath := filepath.Join(cwdDir, "")
			os.Setenv("PLUGIN_CONFIG_FILE_PATH", cfgFilePath)
		case 2:
			createFile(t, sampleConfigFile, "")
			cfgFilePath := sampleConfigFile
			os.Setenv("PLUGIN_CONFIG_FILE_PATH", cfgFilePath)
		case 3:
			createFile(t, sampleConfigFile, sampleConfig)
			cfgFilePath := sampleConfigFile
			Data.KeyCertConf.CertificatePath = sampleConfigFile
			Data.KeyCertConf.PrivateKeyPath = sampleConfigFile
			Data.KeyCertConf.RootCACertificatePath = sampleConfigFile
			os.Setenv("PLUGIN_CONFIG_FILE_PATH", cfgFilePath)
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

func TestValidateConfiguration(t *testing.T) {
	sampleFileForTest := filepath.Join(cwdDir, sampleFileName)
	createFile(t, sampleFileForTest, sampleFileContent)

	tests := []struct {
		name      string
		wantErr   bool
		setConfig func()
	}{
		{
			name: "Invalid RootServiceUUID",
			setConfig: func() {
				SetUpMockConfig(t)
				Data.RootServiceUUID = "invalid"
			},
			wantErr: true,
		},
		{
			name: "Positive Test case - Set Default FirmwareVersion",
			setConfig: func() {
				SetUpMockConfig(t)
				Data.FirmwareVersion = ""
			},
			wantErr: true,
		},
		{
			name: "Empty  RootServiceUUID",
			setConfig: func() {
				SetUpMockConfig(t)
				Data.RootServiceUUID = ""
			},
			wantErr: true,
		},
		{
			name: "Invalid RootServiceUUID",
			setConfig: func() {
				SetUpMockConfig(t)
				Data.RootServiceUUID = ""
			},
			wantErr: true,
		},
		{
			name: "Invalid SessionTimeoutInMinutes",
			setConfig: func() {
				SetUpMockConfig(t)
				Data.SessionTimeoutInMinutes = 0
			},
			wantErr: true,
		},
		{
			name: "Invalid Plugin Config",
			setConfig: func() {
				SetUpMockConfig(t)
				Data.PluginConf = nil
			},
			wantErr: true,
		},
		{
			name: "Invalid Event config",
			setConfig: func() {
				SetUpMockConfig(t)
				Data.EventConf = nil
			},
			wantErr: true,
		},
		{
			name: "Invalid MessageBusConf config",
			setConfig: func() {
				SetUpMockConfig(t)
				Data.MessageBusConf = nil
			},
			wantErr: true,
		},
		{
			name: "Invalid KeyCertConf config",
			setConfig: func() {
				SetUpMockConfig(t)
				Data.KeyCertConf = nil
			},
			wantErr: true,
		},
		{
			name: "Invalid TLS config",
			setConfig: func() {
				SetUpMockConfig(t)
				Data.TLSConf = nil
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setConfig()
			if err := ValidateConfiguration(); (err != nil) != tt.wantErr {
				t.Errorf("ValidateConfiguration() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_checkTLSConf(t *testing.T) {
	tests := []struct {
		name      string
		setConfig func()
		wantErr   bool
	}{
		{
			name:      "With Nil value",
			wantErr:   false,
			setConfig: func() {},
		}, {
			name:      "Positive case 1",
			wantErr:   false,
			setConfig: func() { SetUpMockConfig(t) },
		},
		{
			name:    "Invalid TLS Minimum version",
			wantErr: true,
			setConfig: func() {
				SetUpMockConfig(t)
				Data.TLSConf.MinVersion = "TLS_1.1.2"
			},
		}, {
			name:    "No value for minimum versions",
			wantErr: false,
			setConfig: func() {
				SetUpMockConfig(t)
				Data.TLSConf.MinVersion = ""
			},
		},
		{
			name:    "Invalid TLS max version",
			wantErr: true,
			setConfig: func() {
				SetUpMockConfig(t)
				Data.TLSConf.MaxVersion = "TLS_1.1.2"
			},
		},
		{
			name:    "No value for max versions",
			wantErr: false,
			setConfig: func() {
				SetUpMockConfig(t)
				Data.TLSConf.MaxVersion = ""
			},
		},
		{
			name:    "Invalid Cipher Suite ",
			wantErr: true,
			setConfig: func() {
				SetUpMockConfig(t)
				Data.TLSConf.PreferredCipherSuites = []string{"Invalid"}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setConfig()
			if err := checkTLSConf(); (err != nil) != tt.wantErr {
				t.Errorf("checkTLSConf() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_checkURLTranslationConf(t *testing.T) {
	tests := []struct {
		name      string
		setConfig func()
	}{
		{
			name: "Positive case ",
			setConfig: func() {
				SetUpMockConfig(t)
			},
		},
		{
			name: "Negative case - Nil URLTranslation ",
			setConfig: func() {
				SetUpMockConfig(t)
				Data.URLTranslation = nil
			},
		},
		{
			name: "Negative case - Setting zero NorthBond URL ",
			setConfig: func() {
				SetUpMockConfig(t)
				Data.URLTranslation = &URLTranslation{NorthBoundURL: map[string]string{}}

			},
		},
		{
			name: "Negative case - Setting zero SouthBond URL ",
			setConfig: func() {
				SetUpMockConfig(t)
				Data.URLTranslation = &URLTranslation{SouthBoundURL: map[string]string{}}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setConfig()
			checkURLTranslationConf()
		})
	}
}

func Test_checkCertsAndKeysConf(t *testing.T) {
	sampleFileForTest := filepath.Join(cwdDir, sampleFileName)
	createFile(t, sampleFileForTest, sampleFileContent)
	tests := []struct {
		name      string
		wantErr   bool
		setConfig func()
	}{
		{
			name: "Positive case ",
			setConfig: func() {
				SetUpMockConfig(t)
				Data.KeyCertConf.CertificatePath = sampleFileForTest
				Data.KeyCertConf.PrivateKeyPath = sampleFileForTest
				Data.KeyCertConf.RootCACertificatePath = sampleFileForTest
			},
			wantErr: false,
		},
		{
			name: "Negative case- nil ",
			setConfig: func() {
				SetUpMockConfig(t)
				Data.KeyCertConf = nil
			},
			wantErr: true,
		},
		{
			name: "Negative case- Invalid CertificatePath ",
			setConfig: func() {
				SetUpMockConfig(t)
				Data.KeyCertConf.CertificatePath = ""
			},
			wantErr: true,
		},
		{
			name: "Negative case- invalid PrivateKeyPath ",
			setConfig: func() {
				SetUpMockConfig(t)
				Data.KeyCertConf.CertificatePath = sampleFileForTest
				Data.KeyCertConf.PrivateKeyPath = ""
			},
			wantErr: true,
		},
		{
			name: "Negative case- invalid RootCACertificatePath ",
			setConfig: func() {
				SetUpMockConfig(t)
				Data.KeyCertConf.CertificatePath = sampleFileForTest
				Data.KeyCertConf.PrivateKeyPath = sampleFileForTest
				Data.KeyCertConf.RootCACertificatePath = ""
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setConfig()
			if err := checkCertsAndKeysConf(); (err != nil) != tt.wantErr {
				t.Errorf("checkCertsAndKeysConf() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_checkPluginConf(t *testing.T) {
	tests := []struct {
		name      string
		wantErr   bool
		setConfig func()
	}{
		{
			name:    "Data.PluginConf.ID with empty case ",
			wantErr: false,
			setConfig: func() {
				SetUpMockConfig(t)
				Data.PluginConf.ID = ""
			},
		},
		{
			name:    "Data.PluginConf.Host with empty case ",
			wantErr: true,
			setConfig: func() {
				SetUpMockConfig(t)
				Data.PluginConf.Host = ""
			},
		},
		{
			name:    "Data.PluginConf.Port with empty case ",
			wantErr: true,
			setConfig: func() {
				SetUpMockConfig(t)
				Data.PluginConf.Port = ""
			},
		},
		{
			name:    "Data.PluginConf.Password with empty case ",
			wantErr: true,
			setConfig: func() {
				SetUpMockConfig(t)
				Data.PluginConf.Password = ""
			},
		},
		{
			name:    "Data.PluginConf.UserName with empty case ",
			wantErr: true,
			setConfig: func() {
				SetUpMockConfig(t)
				Data.PluginConf.UserName = ""
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setConfig()
			if err := checkPluginConf(); (err != nil) != tt.wantErr {
				t.Errorf("checkPluginConf() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_checkLBConf(t *testing.T) {
	tests := []struct {
		name      string
		setConfig func()
	}{
		{
			name: "LoadBalancerConf nil",
			setConfig: func() {
				SetUpMockConfig(t)
				Data.LoadBalancerConf = nil
			},
		},
		{
			name: "LoadBalancerConf with Empty Host and Port ",
			setConfig: func() {
				SetUpMockConfig(t)
				Data.LoadBalancerConf.Host = ""
				Data.LoadBalancerConf.Port = ""
			},
		},
	}
	for _, tt := range tests {
		tt.setConfig()
		t.Run(tt.name, func(t *testing.T) {
			checkLBConf()
		})
	}
}

func Test_checkEventConf(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
		setConf func()
	}{
		{
			name:    "Negative Case : Empty Destination URI",
			wantErr: true,
			setConf: func() {
				SetUpMockConfig(t)
				Data.EventConf.DestURI = ""
			},
		},
		{
			name:    "Negative Case : Empty Listener Host",
			wantErr: true,
			setConf: func() {
				SetUpMockConfig(t)
				Data.EventConf.ListenerHost = ""
			},
		}, {
			name:    "Negative Case : Empty ListenerPort",
			wantErr: true,
			setConf: func() {
				SetUpMockConfig(t)
				Data.EventConf.ListenerPort = ""
			},
		},
	}
	for _, tt := range tests {
		tt.setConf()
		t.Run(tt.name, func(t *testing.T) {
			if err := checkEventConf(); (err != nil) != tt.wantErr {
				t.Errorf("checkEventConf() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_checkMessageBusConf(t *testing.T) {
	sampleFileForTest := filepath.Join(cwdDir, sampleFileName)
	createFile(t, sampleFileForTest, sampleFileContent)
	tests := []struct {
		name      string
		wantErr   bool
		setConfig func()
	}{
		{
			name: "Negative Empty EmbType",
			setConfig: func() {
				SetUpMockConfig(t)
				Data.MessageBusConf.EmbType = ""

			},
			wantErr: true,
		},
		{
			name: "Negative invalid EmbType",
			setConfig: func() {
				SetUpMockConfig(t)
				Data.MessageBusConf.EmbType = "Dummy"
				Data.MessageBusConf.MessageBusConfigFilePath = sampleFileForTest
			},
			wantErr: true,
		},
		{

			name: "Negative case EmbQueue is 0",
			setConfig: func() {
				SetUpMockConfig(t)
				Data.MessageBusConf.EmbQueue = []string{}
				Data.MessageBusConf.MessageBusConfigFilePath = sampleFileForTest
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setConfig()
			if err := checkMessageBusConf(); (err != nil) != tt.wantErr {
				t.Errorf("checkMessageBusConf() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
