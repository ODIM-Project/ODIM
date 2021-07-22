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

// Package common ...
package common

import (
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"os"
	"strings"
)

const localhost = "127.0.0.1"

// SetUpMockConfig set ups a mock configuration for unit testing
func SetUpMockConfig() error {
	workingDir, _ := os.Getwd()

	path := strings.SplitAfter(workingDir, "ODIM")
	var basePath string
	if len(path) > 2 {
		for i := 0; i < len(path)-1; i++ {
			basePath = basePath + path[i]
		}
	} else {
		basePath = path[0]
	}
	config.Data.RegistryStorePath = basePath + "/lib-utilities/etc/"
	config.Data.DBConf = &config.DBConf{
		InMemoryPort:   "6379",
		OnDiskPort:     "6380",
		Protocol:       "tcp",
		OnDiskHost:     localhost,
		InMemoryHost:   localhost,
		MaxIdleConns:   10,
		MaxActiveConns: 120,
	}
	config.Data.AuthConf = &config.AuthConf{
		SessionTimeOutInMins:            30,
		ExpiredSessionCleanUpTimeInMins: 15,
	}
	config.Data.APIGatewayConf = &config.APIGatewayConf{
		Port: "9090",
		Host: localhost,
	}
	config.Data.EnabledServices = []string{"SessionService", "AccountService", "EventService"}
	config.Data.URLTranslation = &config.URLTranslation{
		NorthBoundURL: map[string]string{
			"ODIM": "redfish",
		},
		SouthBoundURL: map[string]string{
			"redfish": "ODIM",
		},
	}
	config.Data.AddComputeSkipResources = &config.AddComputeSkipResources{
		SkipResourceListUnderSystem: []string{
			"Chassis",
			"Managers",
			"LogServices",
		},
		SkipResourceListUnderManager: []string{
			"Systems",
			"Chassis",
			"LogServices",
		},
		SkipResourceListUnderChassis: []string{
			"Managers",
			"Systems",
			"Devices",
		},
		SkipResourceListUnderOthers: []string{
			"Power",
			"Thermal",
			"SmartStorage",
			"LogServices",
		},
	}
	config.Data.AuthConf.PasswordRules = &config.PasswordRules{
		MinPasswordLength:       12,
		MaxPasswordLength:       16,
		AllowedSpecialCharcters: "~!@#$%^&*-+_|(){}:;<>,.?/",
	}
	config.Data.RootServiceUUID = "3bd1f589-117a-4cf9-89f2-da44ee8e012b"
	config.Data.FirmwareVersion = "1.0"

	config.Data.ExecPriorityDelayConf = &config.ExecPriorityDelayConf{
		MinResetPriority:    1,
		MaxResetPriority:    10,
		MaxResetDelayInSecs: 36000,
	}
	config.Data.PluginStatusPolling = &config.PluginStatusPolling{
		MaxRetryAttempt:         1,
		RetryIntervalInMins:     1,
		ResponseTimeoutInSecs:   30,
		StartUpResouceBatchSize: 10,
	}
	config.Data.AddComputeSkipResources = &config.AddComputeSkipResources{
		SkipResourceListUnderOthers: []string{"Power", "Thermal", "SmartStorage", "LogServices"},
	}
	return nil
}
