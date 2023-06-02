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
2F/fX7Q=
-----END CERTIFICATE-----`)
	hostCert = []byte(`-----BEGIN CERTIFICATE-----
MIIHHjCCBQagAwIBAgIUWq4pRQucFydxhbi8V9ldsCUkrDswDQYJKoZIhvcNAQEL
BQAwcDELMAkGA1UEBhMCVVMxCzAJBgNVBAgMAkNBMRMwEQYDVQQHDApDYWxpZm9y
bmlhMQwwCgYDVQQKDANIUEUxGDAWBgNVBAsMD1RlbGNvIFNvbHV0aW9uczEXMBUG
A1UEAwwOT0RJTVJBX1JPT1RfQ0EwHhcNMjMwNTMxMDkzNDA2WhcNMzMwNTI4MDkz
NDA2WjBwMQswCQYDVQQGEwJVUzELMAkGA1UECAwCQ0ExEzARBgNVBAcMCkNhbGlm
b3JuaWExDDAKBgNVBAoMA0hQRTEYMBYGA1UECwwPVGVsY28gU29sdXRpb25zMRcw
FQYDVQQDDA5PRElNUkFfU1ZDX0NSVDCCAiIwDQYJKoZIhvcNAQEBBQADggIPADCC
AgoCggIBAOkd1KX8H2cnRQm7E+YZSEPbgw03o76ms5erVNNjrngWg2/MwcAyHSuZ
qhviKcgATrusPAWeUILyh0kv3avg0xjIv9YZSO12fFn/IoDgok+Q7Xge4YdbZvZS
UA6CKJA7/WxaE8Xj8xb+iNnyOYFEoiKsDhH1vcML2Aqo0I0b25I6CQqEyn9lhZLR
Y+O2WN7RoB4Tg4A1f3E0eKF6xA4fkKBrurHquXlzCJN0MLNWT92Cm99q+KaCcDyt
mYyLQREFVdlhKhowmE49rUZn5b8nYlXOcxPXt2tCAAy76raE9/MvaAYYbV6JNmFg
suuvxYFKbc40LbfxB7ImaP36J/5b4QCFSsE7fzb8efaG8lFSHl8PJBkuyojWv0/8
vK4e7UadzhHTagRwz1CSajhdgc6qD6MBqjqgcX+zcpfboRCMtXf9sVN8lYAH35uw
HES0K9npjMsdtC9XQH4Q8giQpONNoYl//UohSuFu0a47aXSvge3XVO9rtWven6gH
xBCDfo7eAyTtq2H70r0x5DF9gRnfGHPxgmNCb2x6/FzYb4gh3GGoKmmnmFKyU1R2
NIZX9WYEiJkIZE/XPbv8H+m9adlEwYl76+rvoMXNyMsRtCpi6JQknHdsvmInSNoS
yJPckWSqUcMW7C6QQtBug3oIri44hZ0yfckP4dZ9eLIbNPsiGI8LAgMBAAGjggGu
MIIBqjAdBgNVHQ4EFgQUR7gpyNJqJP5UIYPi6YQ8Wzd9Kjkwga0GA1UdIwSBpTCB
ooAUP8I+omfJRhcuUxu2aEKK4pqtuRmhdKRyMHAxCzAJBgNVBAYTAlVTMQswCQYD
VQQIDAJDQTETMBEGA1UEBwwKQ2FsaWZvcm5pYTEMMAoGA1UECgwDSFBFMRgwFgYD
VQQLDA9UZWxjbyBTb2x1dGlvbnMxFzAVBgNVBAMMDk9ESU1SQV9ST09UX0NBghQ1
h3TAmq2bikSVzIG7Yn79aUNePzAOBgNVHQ8BAf8EBAMCBeAwHQYDVR0lBBYwFAYI
KwYBBQUHAwIGCCsGAQUFBwMBMIGpBgNVHREEgaEwgZ6CEG9kaW0uZXhhbXBsZS5j
b22CDnJlZGlzLWlubWVtb3J5ggxyZWRpcy1vbmRpc2uCCHVycGx1Z2luggNhcGmC
CmRlbGxwbHVnaW6CEWRlbGxwbHVnaW4tZXZlbnRzggxsZW5vdm9wbHVnaW6CE2xl
bm92b3BsdWdpbi1ldmVudHOCCWlsb3BsdWdpboIQaWxvcGx1Z2luLWV2ZW50czAN
BgkqhkiG9w0BAQsFAAOCAgEAlL7sb9eyOi5S8K+gq418RxB7QOehaI0KfCgzjs2V
lkuRSvjDE1dShyHjyMWcKU3E6YtQyEf6sy+Os/FLU7JO2cghvzIVZu6t3DTwcHRX
d6VO0snaJ/em0+4FYs2H3QE1JetWDY5/s7tlIhMHFWJqNEsSB+wgfG8QMZe4wCtv
h2A1WMKMz0X/urmGg153rPG9f4on5c2/N4KE8erSb1YIVAhmA/Dr/A5/iipPvwOY
uuQVlUNLD7Af+hfSd3UnM8May62z7pa/yrc1E7nLUTTw+VzCPwNcelEZyxeqttKq
6Fx4MDfvIeuB/ZCS+qXWhBPkEyklvuKIVy8qq8MBCZ+HyerUspHaVXMaas6ruzZP
Swco/oSaLvwv9GgotZyRvI+OcsR4BMN/khLcPFv9eCaHyhA6vG95/0dUS7Fb/94i
aaTHCPNwFID4TKAwaOGNsQK9g7G5CxQW7hp4WzFMYDG68zv92MqhKbgjegrjNF1I
kLc8zmLQfR5oUG7H93KuVSNZ4LnIsIsW/30Ry4Pv8XmKjvC6Yn+mJcs6od83AT8n
ngHgGECdv+QIg3TxOeKFLKoVFXwtzFi7GmcmURDTC6rMjGSDcK2sI/onjexog2+s
R1eF+XD/yJlP7Nz8pESQzw4TPpS5i80lMLa2YvtoVRsXgz4BEr+gNK+tt07BnZ5E
FbE=
-----END CERTIFICATE-----`)
	hostPrivKey = []byte(`-----BEGIN PRIVATE KEY-----
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
NspBONT5pKNsdIy3ntWiF+ixkKfe1iM=
-----END PRIVATE KEY-----`)
	hostPubKey = []byte(`-----BEGIN PUBLIC KEY-----
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAvJ9wBDVsCPgglMmHbw/8
Rjr+tbgAdum3cIff3Rrvx2JZ9Gnus63BPrZgTNAiAC9n0I+QP2WgzDU3e6rNl8kZ
zuWwmIY5emReBeJkTpzVD/G35rMlWGXFw0U96b6BAon5Ku4jgt853zvoXkP7sqg4
3G8lHMxETt1TXliSOdNbu4bSTAv1i7UcmGzp0D47i4Ep1DECK46Kc+qjtrDBJehb
DzIhf8k5wtqrmehHHLdYkLyCsfvbe9JZ9m/FEFzZ1lQrB0r/kt3vKrquy15PsNTy
uIQhdnqLhsFhCPWFJfYVgrgudanBcSBviLv5S95DXpc5WIH9qRxSfX95FMi4rOv9
RvZpQmK4Izo/HyJYDT8ajbvsUD7EIh100PvtB2SJ7VzOdV5j0Ek7Qg19l01rnoNb
Z7SGezKskVqqotlwfYuQB+LgI+o3IUM5n7UAX6yfLrPMppF6s2Xfj3gX53omJdRN
knKX300FHpWNijqu8x1H9TTosSF281OVwq52wXa95k1n6dXBl0erX5ozWEXUHTkj
7y54H/9OKj7hPhQkx6nOsAAg47GOJOWIRnSKI30d5T9Pry8hVYXXSsTFoqerqty1
Wg7HWMFTAc97sjcl9nV2alCA/2gBJk8YWnvAHkP5XR11+odT93hBdEEJy3GID0kF
kYwtUoCKTua/vXWa6XSujT0CAwEAAQ==
-----END PUBLIC KEY-----`)
	hostRSAPrivKey = []byte(`-----BEGIN PRIVATE KEY-----
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
NspBONT5pKNsdIy3ntWiF+ixkKfe1iM=
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
