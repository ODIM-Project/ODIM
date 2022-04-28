#(C) Copyright [2022] American Megatrends International LLC
#
#Licensed under the Apache License, Version 2.0 (the "License"); you may
#not use this file except in compliance with the License. You may obtain
#a copy of the License at
#
#    http:#www.apache.org/licenses/LICENSE-2.0
#
#Unless required by applicable law or agreed to in writing, software
#distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
#WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
#License for the specific language governing permissions and limitations
# under the License.

import unittest
import os
import json
from config.config import set_configuraion


class TestConfiguration(unittest.TestCase):

    def test_update_config_data(self):

        sample_config = {
            "RootServiceUUID": "a9762fb2-b9dd-4ce8-818a-af6833ba19f6",
            "LocalhostFQDN": "test.odim.local",
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
                "PasswordRules": {
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
            "SupportedPluginTypes": ["Compute", "Fabric", "Storage"],
            "ConnectionMethodConf": [
                {
                    "ConnectionMethodType": "Redfish",
                    "ConnectionMethodVariant": "Compute:GRF_v1.0.0"
                },
                {
                    "ConnectionMethodType": "Redfish",
                    "ConnectionMethodVariant": "Storage:STG_v1.0.0"
                }
            ],
            "EventConf": {
                "DeliveryRetryAttempts": 1,
                "DeliveryRetryIntervalSeconds": 1,
                "RetentionOfUndeliveredEventsInMinutes": 1
            }
        }
        conf_file = "{cwd}/{file}".format(cwd=os.getcwd(), file="sample.json")
        if not os.path.exists(conf_file):
            with open(conf_file, "w") as f:
                json.dump(sample_config, f)

        set_configuraion(conf_file)
