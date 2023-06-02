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
	"os"
	"strings"
	"testing"
)

var (
	localhost = "127.0.0.1"
	hostCA    = []byte(`-----BEGIN CERTIFICATE-----
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
0Gs8qmk=
-----END CERTIFICATE-----`)
	hostCert = []byte(`-----BEGIN CERTIFICATE-----
MIIHKTCCBRGgAwIBAgIUc175RPwiWfyW/DU3mPEwn9yTqt0wDQYJKoZIhvcNAQEL
BQAwcDELMAkGA1UEBhMCVVMxCzAJBgNVBAgMAkNBMRMwEQYDVQQHDApDYWxpZm9y
bmlhMQwwCgYDVQQKDANIUEUxGDAWBgNVBAsMD1RlbGNvIFNvbHV0aW9uczEXMBUG
A1UEAwwOT0RJTVJBX1JPT1RfQ0EwHhcNMjMwNjAyMDczNzQwWhcNMzMwNTMwMDcz
NzQwWjBwMQswCQYDVQQGEwJVUzELMAkGA1UECAwCQ0ExEzARBgNVBAcMCkNhbGlm
b3JuaWExDDAKBgNVBAoMA0hQRTEYMBYGA1UECwwPVGVsY28gU29sdXRpb25zMRcw
FQYDVQQDDA5PRElNUkFfU1ZDX0NSVDCCAiIwDQYJKoZIhvcNAQEBBQADggIPADCC
AgoCggIBAMA/+7F1q7DVLcsZn4gpfFfV41JFM6lSpOzUsqoJB4r0+wxi+pADia/j
obuTPnOoMdtm4wEUeWKzsTHg3iGo4I1ktUSvHdI+7LzxbE0N6buISv+Ric+5YFba
eHJaOLF+8N6EodFTrgPmgDLezEmNtgN6eaMAGFjem7gC+inMhDgiZMF7oJTa+bP2
F92w6phJjjzipkuAOTFEiA5B8SL9TRGsODdgyOH1jKSmHb99Z6IjbatYXgFqTlr3
8vNi59/bZw2yRIGBEiQ7fY488UOQxQ8maUvx81eZ4R4zfrIpQuIlhMWcYgXPmeer
1lnjkBc90NxjnQaPSI1FigCshvDIIWy2m8k8AzWdskRvqs1on1NWsOqvvdrXxW/L
QXw9GWTfkOWv1HE/zgf/U1MNvHbyB2SQp6NoHtvwd/o47bxn3R24x9kGCtlyf3pE
gHB8kNAFLZi/0kPv5f0wt3+udpw0hVnUtCM72zDumWcCsrCOXyPAXS0sKpk/2y1y
VR3opdqPJ7/zOMLKR23+RFz3ypyEuXvhkg/rN9W1mM41tEr0aJ6ov4kDeSKOFRBp
SKlKpdN1c2043c/SHdz6BIZ+n1yjDP06nHzZ+3W9Q/G8+p6VGLYl766Rcfzx5zEK
DKEKVgsHJDgjoyogAzw7tMf+uEsBWbZLfrswxnhqjlH/rhwNAR1TAgMBAAGjggG5
MIIBtTAdBgNVHQ4EFgQUP0DguuY3ZxAcgMl0T+XkmQRl2I8wga0GA1UdIwSBpTCB
ooAUTE5x7fEbNijbFRY2xR807zgF1NihdKRyMHAxCzAJBgNVBAYTAlVTMQswCQYD
VQQIDAJDQTETMBEGA1UEBwwKQ2FsaWZvcm5pYTEMMAoGA1UECgwDSFBFMRgwFgYD
VQQLDA9UZWxjbyBTb2x1dGlvbnMxFzAVBgNVBAMMDk9ESU1SQV9ST09UX0NBghRc
+fecgWeTFLPfXyVq82YARJAArTAOBgNVHQ8BAf8EBAMCBeAwHQYDVR0lBBYwFAYI
KwYBBQUHAwIGCCsGAQUFBwMBMIG0BgNVHREEgawwgamCEG9kaW0uZXhhbXBsZS5j
b22CDnJlZGlzLWlubWVtb3J5ggxyZWRpcy1vbmRpc2uCCHVycGx1Z2luggNhcGmC
CmRlbGxwbHVnaW6CEWRlbGxwbHVnaW4tZXZlbnRzggxsZW5vdm9wbHVnaW6CE2xl
bm92b3BsdWdpbi1ldmVudHOCCWlsb3BsdWdpboIQaWxvcGx1Z2luLWV2ZW50c4IJ
bG9jYWxob3N0MA0GCSqGSIb3DQEBCwUAA4ICAQCLA/Af5606Glq89Th1G9+FNrsn
ujc/3EXtCFHNsVZDJ7b/EApCWDHibOtoo3CIIxYvyjAPdNijx9FaAfuf2Lyt+aCT
cR96N0MLh3gCznmjvu/lknIG196cckY5ZIMPucjMOlYJo5jjx7aburgeUA/NaGRT
8NC43uzBJ6KMjs7AaKpI9LYl68ishP+n9caBxohxujd8xMD0T/N3KjSNqbNu1sKI
YK3mlrn0LUAVacPJNCWFQTKNjBeBLCZo3rDHe3CzQCwzasJxYnKP2V1JW/hObmCP
XDXgRCzfz2hHx8QfReV5jd3CkxX0Cp/k3SWItYbQVBmdyf0Nqz96ZLJ05k7Of7s8
HH8YG8T9YtDSj4WWxpVrgtuijDnHQq3YOcOMdKyBA5jkEtA+9kFg8oSW/Qi56Naz
9cMIpiapftAbhG1fBHGf0MXHrZQJA/4UuIVXroBb23IGl+x3gNNwlK/Ji/BD65Z7
XDHF76jOu5I8Dq44ACwv5/FX5F8zvZTlPDwADA0mgT6a6ZYZKGtN1cCLj16cbit0
XXf8m2Dyhq9hqU0R2zCnBVhgkZZ1AQzb5FA1+rmHtPepHAcy2H9D3i8hewZrTDK4
56gX1x+VY6QRI+w7tKMRexPyaKLDYvOsDeUiOM0ihum5WRSquTgU2nTEJoeHH4Z6
2IPTgssoAAEoW9EgtA==
-----END CERTIFICATE-----`)
	hostPrivKey = []byte(`-----BEGIN PRIVATE KEY-----
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
P3y4SIcsoq4kiR2L9mSuG0iPbtuN
-----END PRIVATE KEY-----`)
	hostPubKey = []byte(`-----BEGIN PUBLIC KEY-----
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAkUuiBkDMvX1bnBoyAvzP
bclO5VzSnYZdfRCZtRsbzSkwDDW1Yrt+JMGz0MB3bAx9quI+PIBzkEIO4Cni1pbO
1rlgE10oiIB2nhpiqMvwTlGk+2cwTEkhTp4Lertx+dhraK8KiJVuUEx8GlOuZ2sp
oOZNxaEsiinCmY0WqJyJXiW3YFjZ73HMbdEGO2WEuek0Wt11xDiVnBeTiWgdepyX
5XFYtIGjbSeiSGK+X6CiohlHzBQoKQv8ZYBfkXVun2dfS5YVl/m7FTg0MSjFy2Lw
UM458sPOhw1CrL0wrKzVV/jWE/jo5FS5gM6SS43YifnyuQGc2TG3ymAECUl4qhSa
/jVfSOKPC53+0O7sEY4+PyJhjqBrEZfZSXD1LrkrlUQZkTAT0x41aFrszFtfIBYN
O1emvSM/wRSmt6nzOpot5uposQ34AhNVB6PlXwKdeEwWrtbUdzlzKHwEFqZ+6OGp
x7/8F5hdFuTWh7u4vV5GaDIwt9nkmvbOEaYojEQpWMm/bLKtHWzbk03MTpY7jVcx
K3DyPYIFN9n7+KCnvrSg5uhoR1Z9BTulpd0lSRqAKYU5n5QdEfWzxYyyUD3OErbo
bOXXzi3vKRB7BLFLTELdqCSz4QYdTZ6OcF6ngExN39/61LRTHzAnDs7rVVW/2Cmv
DQsjLLB3DMffD/zqLo3nyYsCAwEAAQ==
-----END PUBLIC KEY-----`)
	hostRSAPrivKey = []byte(`-----BEGIN PRIVATE KEY-----
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
P3y4SIcsoq4kiR2L9mSuG0iPbtuN
-----END PRIVATE KEY-----`)
)

// SetUpMockConfig set ups a mock ration for unit testing
func SetUpMockConfig(t *testing.T) error {
	workingDir, _ := os.Getwd()

	Data.RootServiceUUID = "3bd1f589-117a-4cf9-89f2-da44ee8e012b"
	Data.FirmwareVersion = "1.0"
	Data.SouthBoundRequestTimeoutInSecs = 10
	Data.ServerRediscoveryBatchSize = 10
	path := strings.SplitAfter(workingDir, "ODIM")
	var basePath string
	if len(path) > 2 {
		for i := 0; i < len(path)-1; i++ {
			basePath = basePath + path[i]
		}
	} else {
		basePath = path[0]
	}
	Data.EventForwardingWorkerPoolCount = 1
	Data.EventSaveWorkerPoolCount = 1
	Data.RegistryStorePath = basePath + "/lib-utilities/etc/"
	Data.LocalhostFQDN = "odim.test.com"
	Data.EnabledServices = []string{"SessionService", "AccountService", "EventService"}
	Data.DBConf = &DBConf{
		Protocol:              "tcp",
		InMemoryHost:          localhost,
		InMemoryPort:          "6379",
		OnDiskHost:            localhost,
		OnDiskPort:            "6380",
		MaxIdleConns:          10,
		MaxActiveConns:        120,
		RedisInMemoryPassword: []byte("redis_password"),
		RedisOnDiskPassword:   []byte("redis_password"),
	}
	Data.MessageBusConf = &MessageBusConf{
		MessageBusType:          "Kafka",
		OdimControlMessageQueue: "odim-control-messages",
	}
	Data.MessageBusConf.MessageBusConfigFilePath = "mockfilepath"
	Data.KeyCertConf = &KeyCertConf{
		RootCACertificate: hostCA,
		RPCPrivateKey:     hostPrivKey,
		RPCCertificate:    hostCert,
		RSAPublicKey:      hostPubKey,
		RSAPrivateKey:     hostRSAPrivKey,
	}
	Data.AuthConf = &AuthConf{
		SessionTimeOutInMins:            30,
		ExpiredSessionCleanUpTimeInMins: 15,
		PasswordRules: &PasswordRules{
			MinPasswordLength:       12,
			MaxPasswordLength:       16,
			AllowedSpecialCharcters: "~!@#$%^&*-+_|(){}:;<>,.?/",
		},
	}
	Data.APIGatewayConf = &APIGatewayConf{
		Port:        "9090",
		Host:        localhost,
		PrivateKey:  hostPrivKey,
		Certificate: hostCert,
	}
	Data.AddComputeSkipResources = &AddComputeSkipResources{
		SkipResourceListUnderSystem: []string{
			"Chassis",
			"LogServices",
			"Managers",
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
	Data.URLTranslation = &URLTranslation{
		NorthBoundURL: map[string]string{
			"ODIM": "redfish",
		},
		SouthBoundURL: map[string]string{
			"redfish": "ODIM",
		},
	}
	Data.PluginStatusPolling = &PluginStatusPolling{
		MaxRetryAttempt:          1,
		RetryIntervalInMins:      1,
		ResponseTimeoutInSecs:    1,
		StartUpResourceBatchSize: 1,
		PollingFrequencyInMins:   1,
	}
	Data.ExecPriorityDelayConf = &ExecPriorityDelayConf{
		MinResetPriority:    1,
		MaxResetPriority:    10,
		MaxResetDelayInSecs: 36000,
	}
	Data.TLSConf = &TLSConf{
		VerifyPeer: true,
		MinVersion: "TLS_1.2",
		MaxVersion: "TLS_1.2",
		PreferredCipherSuites: []string{
			"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256",
			"TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384",
			"TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256",
			"TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384",
			"TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256",
		},
	}
	Data.SupportedPluginTypes = []string{"Compute", "Fabric"}
	Data.ConnectionMethodConf = []ConnectionMethodConf{
		{
			ConnectionMethodType:    "Redfish",
			ConnectionMethodVariant: "Compute:BasicAuth:GRF:1.0.0",
		},
		{
			ConnectionMethodType:    "Redfish",
			ConnectionMethodVariant: "Storage:BasicAuth:STG:1.0.0",
		},
	}
	Data.EventConf = &EventConf{
		DeliveryRetryAttempts:        1,
		DeliveryRetryIntervalSeconds: 1,
	}
	Data.TaskQueueConf = &TaskQueueConf{
		QueueSize:        1000,
		DBCommitInterval: 1000,
		RetryInterval:    1000,
	}
	Data.PluginTasksConf = &PluginTasksConf{
		MonitorPluginTasksFrequencyInMins: 60,
	}
	SetVerifyPeer(Data.TLSConf.VerifyPeer)
	SetTLSMinVersion(Data.TLSConf.MinVersion, &WarningList{})
	SetTLSMaxVersion(Data.TLSConf.MaxVersion, &WarningList{})
	SetPreferredCipherSuites(Data.TLSConf.PreferredCipherSuites)
	return nil
}
