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
6ATq
-----END CERTIFICATE-----`)
	hostCert = []byte(`-----BEGIN CERTIFICATE-----
MIIFsjCCA5qgAwIBAgIUNN0LnDt13/3nM1GipqN9Pku+ZdEwDQYJKoZIhvcNAQEN
BQAwUjELMAkGA1UEBhMCVVMxCzAJBgNVBAgMAkNBMRMwEQYDVQQHDApDYWxpZm9y
bmlhMQ8wDQYDVQQKDAZPRElNUkExEDAOBgNVBAMMB1Jvb3QgQ0EwHhcNMjExMDI2
MDkzNDE5WhcNMjMwMzEwMDkzNDE5WjBBMQswCQYDVQQGEwJVUzELMAkGA1UECAwC
Q0ExDzANBgNVBAoMBk9ESU1SQTEUMBIGA1UEAwwLU2VydmVyIENlcnQwggIiMA0G
CSqGSIb3DQEBAQUAA4ICDwAwggIKAoICAQC1vI9gNVOQnpoSsm2Bonk41xbcJ1CR
5tj6wswEOBaDNXHtYOq85vMT29ou4r5DYanN05Fhd+uJC22Gf4nbGoLkkt1AuIIb
d2U3PY9YJvNfQ2R2GAvTSBTuB5wZ6Co7JEX9kRlVbCcQrF1TJjTI0MgSsHf5U4rg
pr70YN2WecpUieQ3GG+x94NpupwwGJWeh8JIgDxAdh8m/cebOT1VO5xxtPpcd8n4
TxlroPLbgz3UFyMiElmih3+w+L89pOuuPgivBdIVKNsYLYICFw9s09cgwlsi6qhy
Esa0yMjLARuaxz2VFVPwDH9+JmGHWWIUfjXnQ3WixGTrpAI1iboN24DrFtjdIRfP
RgVpuDP7HSN2/CNIXxQ2RQF0fB/z4g1KJS56UgGFCvgDl2BG3aOj+tlQofRvgxT9
jdBBkeiCLgixkVs4nu73rB4Bq4RHRosrpt4kT4nmioMuUcLJRsGB1zt0+xQ+DR6L
vsnzbD4bS3fYADuTzDzXjBem9Z2RYUX9wUtXeiOqqNXepb/AeHmZLzUlVlSu5h6c
PxR9adSVw8NBSR1JXn0hqLFjOVEDk3wMicQnzRXYUrpJgh0tUqu6gHs7DHUsYQPj
qzXxMZbpJyMZSJirqu43VkyVhyIswnzGErwDs6gEM9qUlrQ7L4Zcpy+Af9SqpZPG
aB6p6/m1a4T7XwIDAQABo4GQMIGNMAkGA1UdEwQCMAAwCwYDVR0PBAQDAgXgMB0G
A1UdJQQWMBQGCCsGAQUFBwMBBggrBgEFBQcDAjAdBgNVHQ4EFgQUdWXGZkyYvxcX
A8WCUJ9Df/KA4w0wHwYDVR0jBBgwFoAU2x5lPlMmcsAOk5p96W+fpo2kPeYwFAYD
VR0RBA0wC4IJbG9jYWxob3N0MA0GCSqGSIb3DQEBDQUAA4ICAQBpMuSlw82RHEFJ
YOFy8iPM8bxTThrZRzCpubePKdRPOKEor8WXuHyiKNhbaaPdvgHgPvTBYjFGOsJn
gjWbUc8WwhyCj/1U+m/xSiOTzgbs2sE5Qa4/v4s7Xai3GRAeTC3Y+QDTpaZQWkRT
mSoZQbELzjPYtyxcblUUzLjjupjjld5vlNKKoPyl6WdtRrnkkvgABEFAysgWyW1l
yKO/xfe6sBS7H/UWqZOZZQQIy8svmCzNEPXsr1UwjQiJB74KgZrRuomm/FaA+fVW
jpAQ+LsN6Dhe9dv8J4e4MgHA9uaiFRwI1z+xeQnUHXJUbWHQNWpbI62ZTH3LftFT
+f54Obh3DdVyt+zugBZoIHUaAB+eKcTcruVEwDcDRFS0ObyqBkIej2ED+1iTFw3l
VoASse8XMdn734a3Uso65yUoTCj/Z54/N09m24ZRXyoyahgWwl0Kc+2lINcEEj7l
zAAjO94sWIWvFnA1c2rgx9zFQd7PU2zIHOP/A3tse3yyoszvS38eiTZkvP2Tj0mu
AzjfunVr6C7YjfYWnAa3tGYTpSlKnKVnCl7vIfNtthC2i3EcXACQcCecQOBpp6A2
7pCccBHHUvPMuETjEnCC5wSrSYA0MSHOreS0to7w+NGAc5/YjINN+7Zo6JooALNk
KoQRsFQ0udOGCu+Z5uZfMr0nre618w==
-----END CERTIFICATE-----`)
	hostPrivKey = []byte(`-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----`)
	hostPubKey = []byte(`-----BEGIN PUBLIC KEY-----
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEA+Q01x6+Jgpi8Tdo3ul3P
iuacjeJsFr32gmdo1l0oJp4FweMtiTm9k/pggQMrUGcPqV0ANI4h0DVHx+RdR27R
DV1bqMJqOghTAKbCOw4Wh0X9vbITiuxhLPsYMhSsOxY0Au/YVJd8/ZZQ7QHqKf5V
SvF0ickFrKP4Rkd9LecAXTlXwVxBYBMw8Y1OZVuv272soFMVt8BblCKxNEM8pX8H
KcFoXMrWZr2tgrxDJi7Ry1zNJ4S/gBkPYJfNj8+lPOwwc1nUKIIbzAGkN67h3Q9I
ZRlyyM8D7ayZEKk3tfhNvD9lAip24yORWQDocQ8+wsjerXtTJU/bdqDpLPeAvTdv
QViOzzrMvikIpw9YzbRN6i17jA26BEI0yOgtLLcHOA2ah+K/0kDpHINR3YX0TNfL
SeMEWylQq2Sv6cRO9u0iaRih4GHfWOkc0R/4VaRftS2TJmGEVKT2h7XEnlspCnDw
OdPafLtPKL+aNAdoDnS9fALkAb/lGskmsM/tmSrS4tjSgYhdsYMQp2FseyhtHfsC
4hLW0AnQn6ckdlr2kwXGOc+kpDcoWtc9V16rkaCoNjTu76P8nWIjEasNBZvm3unV
o79i25P4izfyzAG/tdOA+NVbArJEjBaHge0ekJKajHFPLMaaJ1kptItWS1PGVy5d
ZlgyKGJ8O0RB8M1vofMdLfMCAwEAAQ==
-----END PUBLIC KEY-----`)
	hostRSAPrivKey = []byte(`-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----`)
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
	Data.RegistryStorePath = basePath + "/lib-utilities/etc/"
	Data.LocalhostFQDN = "odim.test.com"
	Data.EnabledServices = []string{"SessionService", "AccountService", "EventService"}
	Data.DBConf = &DBConf{
		Protocol:       "tcp",
		InMemoryHost:   localhost,
		InMemoryPort:   "6379",
		OnDiskHost:     localhost,
		OnDiskPort:     "6380",
		MaxIdleConns:   10,
		MaxActiveConns: 120,
	}
	Data.MessageBusConf = &MessageBusConf{
		MessageBusType:  "Kafka",
		MessageBusQueue: []string{"redfish-topic"},
	}
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
		MaxRetryAttempt:         1,
		RetryIntervalInMins:     1,
		ResponseTimeoutInSecs:   1,
		StartUpResouceBatchSize: 1,
		PollingFrequencyInMins:  1,
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
		DeliveryRetryAttempts:                 1,
		DeliveryRetryIntervalSeconds:          1,
	}
	SetVerifyPeer(Data.TLSConf.VerifyPeer)
	SetTLSMinVersion(Data.TLSConf.MinVersion)
	SetTLSMaxVersion(Data.TLSConf.MaxVersion)
	SetPreferredCipherSuites(Data.TLSConf.PreferredCipherSuites)
	return nil
}
