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
MIIFFDCCAvygAwIBAgIUSRrUi66SkjZFWwLVDhdyrFX65pYwDQYJKoZIhvcNAQEN
BQAwFDESMBAGA1UEAwwJbG9jYWxob3N0MB4XDTIwMDQxODEyMjYwN1oXDTMwMDQx
NjEyMjYwN1owFDESMBAGA1UEAwwJbG9jYWxob3N0MIICIjANBgkqhkiG9w0BAQEF
AAOCAg8AMIICCgKCAgEAtl+m9W7LvxhV9ghNu8Rr2VMUCha3LB+k4XOiey0Wt/DF
e1D36tJBbfNHfe+6m4/GHM4e4Tp+h8inCyBCCaHLPXfcA0jY3bHOC0cYjRjhx4rP
TG06ok2YIhovfJIZcANse8nSCenAFSSJqdAAKlsNH9z5lT+JhPr9kG+48WnvDPeX
DNWNf/VUqkYjon7q5hLQL9ImtGSFmjPMKFjNACSDZfIyb/x6TUuTVHBmtKT25+rp
4FHa3mxBYNLisavyNlg59U189FISAGYjC+fOlYQBL7RXb2Xrrd4DgGJC1N8Hr0To
z9zmFj0u4sJ5RFdOabDCNZgEx3tGD1fYqjEdSuspRzMu6CUl8PjN/C0Ml0NOHyHn
7x/0VKxOwGUW0eoGE9BNGfSONrX6Wt3Ej3uKkyDPxIX0PhbgHewQF/jF+698hv10
uGoeAZLhyZlFhCxPCZe4EtBOQ7LXWcEroneoVj4CWY0tm1NVvDY7PVw4MsnloHW/
/qfjgi9yQbaJ7IllSNGaZhqYvtc8bXH5xwd4/F9wOkKC6p63RqB0Pom7lDOBeiha
uhE8gsqg0oCJCkYyKdQ8Ye6xiF+cfmnV5VT+S9XLyTza3utuuJPG8qg6xY5aVBkg
P5JrVaiDTQg19FopUFHfUVY8K121s4WMS02XjySpLyY11ZGrqp3hqIRL+kUJX7sC
AwEAAaNeMFwwDwYDVR0TAQH/BAUwAwEB/zAdBgNVHQ4EFgQUE3wuzAFexqxv1/4o
HH2HNVut8q4wCwYDVR0PBAQDAgHmMB0GA1UdJQQWMBQGCCsGAQUFBwMCBggrBgEF
BQcDATANBgkqhkiG9w0BAQ0FAAOCAgEATgew4Rh3VEgOm6OF5od83AUm/R2nomMe
FXNAw70duP3miYsQbeE9jiLL5vF3Rwv9X3EQ7yBo+EAnKaaMha5Ee9IT8Yu9LnFj
qYKkPlK/fslFROFASZCyLltwJ3KBwfNcUBh+55IrYjvEzh1NUOK/9r+Lt7fn6KVq
a0vVqF9HZjnmr7I4yOw142QMS7nPmQ/PF7ql/Q8IuAuFhJA3pFPK5YO0p0cdaWmS
X0RSmxr0ANkbPWp82lonX4aYz+2uqXDHWBv+kQBMdhU400s9MC+q1vjMD+C+qG3g
BO1bt8/p/DDFnZPwwUyHvU8VIxL99oYDi4nk02YRJrXpsWhHogeD9m8Kv3BcVFx2
7kvuhV+2UiX6IMyricN4EQzG+XS+Pu6UkxiHrxIl1BNVjc3pfRjPUOIU8wt1nBHS
v3jiYtR7fAqWHdfiqz9LDjfxam1QqOt3l6CVMnuQvaZF+nzC/Aa7l09ElrsSWrF5
sp5MSJplugY7r3qG0/evQ/HzkThhVD+0HMwUGx9xa+Z96kRzQ5qYEaj03IU8qZve
wMEvS6jE4aNQ5PFZ/9HKIJ7LkcubuBpaMGTeNzO6vn+9XvxxHyTlMWLpfm04ukBI
3ylqIAA630pHjCUJGPHBiYXw9DD7MRGczqkOLVPE3sOHuTiyL/DNQ9kSReyMG8Ss
8Uv+gESYkqk=
-----END CERTIFICATE-----`)
	hostCert = []byte(`-----BEGIN CERTIFICATE-----
MIIFVjCCAz6gAwIBAgIUKcTsAJKTQSMIKWDWYdZf/Y0uxLUwDQYJKoZIhvcNAQEL
BQAwFDESMBAGA1UEAwwJbG9jYWxob3N0MB4XDTIwMDQxODEyMjk1OFoXDTMwMDQx
NTEyMjk1OFowFDESMBAGA1UEAwwJbG9jYWxob3N0MIICIjANBgkqhkiG9w0BAQEF
AAOCAg8AMIICCgKCAgEA5eYgv0jdi0GY+x+fs9bFwEX6Knam7cWhxUjAMUs00ii8
8n6o5XE0FXbfLs44xf/FpYkpZtf+qXqPXLQCohfeUUiOXGRTGUkQA5hVVxj6gE+U
NqrOSocxwep3KpACPYIniAB2hSz0Dt6fxAdM29SONtNHI9xihuGm9deN+5Dx+ZgW
gcaMfio62wlwTdMcD6VZT3oaYog/+dvS7tqU55tROC335GEGYRg2E7RP+elsvESK
wm1x8US5cfPWgcY2P/hGfbrnhsqW2XVaVrA7KQzWS9/1WRrq7j41sIbG5WHydMLI
kmaXnMuZ+xz4tcEsnt2MOlL/izXJs4H+pUpL+iTF/i9MSRfw2Bi18fVH1CitNRdo
M5SrAmt+YFeDw8N6SuXtIu8l2fQecCeieCSbHraW+VGqY+n9Em4MNgw+tSxl3lsi
j4ijvbptGjbvs2Vkr37MwYqfwIWfLEOdYJVcNgJI72T+mhpoHtKTZ50Ugk4MOkNL
WWuX8wX7UHswwuSTv2jORBXdoyq3mz6PqNrnXOGzmujhVOyp9PxV/AT66Px5GIWy
VpC4MMHiJCNSGP/ocsKepAg23fmboQ3OzjB9Y+AHUnp8ML/h2HhCKMSREj/BilRO
K+Om3ZA/baCngtEvna2pXKuHFKHiHILpyaGrlhRm55j094VEHrNfNHfCJ2+Qf10C
AwEAAaOBnzCBnDAdBgNVHQ4EFgQUqNezn0zjq7KHDGV0IZy94danHiwwTwYDVR0j
BEgwRoAUE3wuzAFexqxv1/4oHH2HNVut8q6hGKQWMBQxEjAQBgNVBAMMCWxvY2Fs
aG9zdIIUSRrUi66SkjZFWwLVDhdyrFX65pYwCwYDVR0PBAQDAgWgMB0GA1UdJQQW
MBQGCCsGAQUFBwMCBggrBgEFBQcDATANBgkqhkiG9w0BAQsFAAOCAgEAiqyKAGvP
HWt8boQe6KkiH5tmPjrd9U6YSsB3qzzfac20BIuY7ZSdpNf2REt/pAUKUyO6mYyV
uuIHtwd58Yuvdop3WD+K70qTKqXB8wNeyfQNeWFhEDumJikbjcRRNn9P+GhCTATI
wp+U7zxxXOxWWlm1sDsZmAuouaFVDYG0F4AUdHP16gLrZEkzup9r7VnUFpw/0tkB
k1PR8RKgpz4r/x5QtCft+FFcyPYqRZNiA/EJ9o+ObtZe2ig5QTpzJhGSwTRm9emt
PxLKv672deyjyOypIgwa94dZRtaiZWK8apDMmzx/ab1Q0rMyZzaVp/EHmEyDQNDd
zqiapUOu2t4vu3NhqgEjLTfSgQR7Kj5BguzQfZ3oJr4zn8JGKfkPdW0OV0LFG9o3
C0HvinWBueGW8dpiHGbOVH57iOx8CU+AesDK5T2rD3kryJ+fos0fSO0fWwmkNYMZ
WkH7hmtVup2jCQy97T0Jlpwi9/gkhnJW4XV1HbAtdkX5QiHC0xs7oTlTpSjm409d
8Zqevusshp2FLkZC4NypW34dJhkrJl/vFJN+zRlwNQUu4fEVFJUb3tTkPWGqYWwz
Ne4ovENCT7e8oBqFDJrATPbHs06b1GNk1nxO1BtAgnba37/7WyQxXK9aD5j8w3pK
dIRctXGgrIroKS/k4eFQa99kbmFl5Khb8Wo=
-----END CERTIFICATE-----`)
	hostPrivKey = []byte(`-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----`)
	hostPubKey = []byte(`-----BEGIN PUBLIC KEY-----
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEA5eYgv0jdi0GY+x+fs9bF
wEX6Knam7cWhxUjAMUs00ii88n6o5XE0FXbfLs44xf/FpYkpZtf+qXqPXLQCohfe
UUiOXGRTGUkQA5hVVxj6gE+UNqrOSocxwep3KpACPYIniAB2hSz0Dt6fxAdM29SO
NtNHI9xihuGm9deN+5Dx+ZgWgcaMfio62wlwTdMcD6VZT3oaYog/+dvS7tqU55tR
OC335GEGYRg2E7RP+elsvESKwm1x8US5cfPWgcY2P/hGfbrnhsqW2XVaVrA7KQzW
S9/1WRrq7j41sIbG5WHydMLIkmaXnMuZ+xz4tcEsnt2MOlL/izXJs4H+pUpL+iTF
/i9MSRfw2Bi18fVH1CitNRdoM5SrAmt+YFeDw8N6SuXtIu8l2fQecCeieCSbHraW
+VGqY+n9Em4MNgw+tSxl3lsij4ijvbptGjbvs2Vkr37MwYqfwIWfLEOdYJVcNgJI
72T+mhpoHtKTZ50Ugk4MOkNLWWuX8wX7UHswwuSTv2jORBXdoyq3mz6PqNrnXOGz
mujhVOyp9PxV/AT66Px5GIWyVpC4MMHiJCNSGP/ocsKepAg23fmboQ3OzjB9Y+AH
Unp8ML/h2HhCKMSREj/BilROK+Om3ZA/baCngtEvna2pXKuHFKHiHILpyaGrlhRm
55j094VEHrNfNHfCJ2+Qf10CAwEAAQ==
-----END PUBLIC KEY-----`)
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
	Data.KeyCertConf = &KeyCertConf{
		RootCACertificate: hostCA,
		RPCPrivateKey:     hostPrivKey,
		RPCCertificate:    hostCert,
		RSAPublicKey:      hostPubKey,
		RSAPrivateKey:     hostPrivKey,
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
		SystemCollection: []string{
			"Chassis",
			"LogServices",
		},
		ChassisCollection: []string{
			"Managers",
			"Systems",
			"Devices",
		},
		OtherCollection: []string{
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
			"TLS_RSA_WITH_AES_128_GCM_SHA256",
			"TLS_RSA_WITH_AES_256_GCM_SHA384",
			"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256",
			"TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384",
			"TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256",
			"TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384",
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
	SetVerifyPeer(Data.TLSConf.VerifyPeer)
	SetTLSMinVersion(Data.TLSConf.MinVersion)
	SetTLSMaxVersion(Data.TLSConf.MaxVersion)
	SetPreferredCipherSuites(Data.TLSConf.PreferredCipherSuites)
	return nil
}
