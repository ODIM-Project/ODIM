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
	"crypto/tls"
	"net/http"
	"reflect"
	"testing"
)

type modifyTLSData func(tlsConfig *tls.Config)

var nonX509Certificate = []byte("non x509 certificate")

func TestGetHTTPClientObj(t *testing.T) {
	if err := SetUpMockConfig(t); err != nil {
		t.Fatal("error: SetUpMockConfig failed with", err)
	}

	httpConf := &HTTPConfig{
		Certificate:   &Data.APIGatewayConf.Certificate,
		PrivateKey:    &Data.APIGatewayConf.PrivateKey,
		CACertificate: &Data.KeyCertConf.RootCACertificate,
	}

	tlsConfig := &tls.Config{}
	httpConf.LoadCertificates(tlsConfig)
	Client.SetTLSConfig(tlsConfig)

	tests := []struct {
		name    string
		exec    modifyTLSData
		want    *tls.Config
		wanterr bool
	}{
		{
			name:    "Positive Case",
			exec:    nil,
			want:    tlsConfig,
			wanterr: false,
		},
		{
			name: "Verify Peer set",
			exec: func(tlsConfig *tls.Config) {
				SetVerifyPeer(true)
				Client.SetTLSConfig(tlsConfig)
				tlsConfig.InsecureSkipVerify = false
			},
			want:    tlsConfig,
			wanterr: false,
		},
		{
			name: "Ciphers not configured",
			exec: func(tlsConfig *tls.Config) {
				Client.SetTLSConfig(tlsConfig)
			},
			want:    tlsConfig,
			wanterr: false,
		},
		{
			name: "Invalid CA certificate",
			exec: func(tlsConfig *tls.Config) {
				httpConf.CACertificate = &nonX509Certificate
			},
			want:    nil,
			wanterr: true,
		},
	}
	for _, tt := range tests {
		if tt.exec != nil {
			tt.exec(tt.want)
		}
		t.Run(tt.name, func(t *testing.T) {
			httpClient, err := httpConf.GetHTTPClientObj()
			if (err != nil) != tt.wanterr {
				t.Errorf("GetHTTPClientObj() err = %v, wanterr %v", err, tt.wanterr)
			}
			if httpClient == nil && tt.want != nil {
				t.Errorf("GetHTTPClientObj() : Expected valid HTTP client object")
				return
			}
			if tt.want != nil {
				gotTLSConfig := httpClient.Transport.(*http.Transport).TLSClientConfig
				if gotTLSConfig == nil {
					t.Errorf("GetHTTPClientObj() : Expected HTTP client object to contain valid tls.Config")
					return
				}
				if gotTLSConfig.MinVersion != tt.want.MinVersion {
					t.Errorf("GetHTTPClientObj() gotTLSConfig.MinVersion = %+v, wantTLSConfig.MinVersion = %+v",
						gotTLSConfig.MinVersion, tt.want.MinVersion)
				}
				if gotTLSConfig.MaxVersion != tt.want.MaxVersion {
					t.Errorf("GetHTTPClientObj() gotTLSConfig.MaxVersion = %+v, wantTLSConfig.MaxVersion = %+v",
						gotTLSConfig.MaxVersion, tt.want.MaxVersion)
				}
				if gotTLSConfig.InsecureSkipVerify != tt.want.InsecureSkipVerify {
					t.Errorf("GetHTTPClientObj() gotTLSConfig.InsecureSkipVerify = %+v, wantTLSConfig.InsecureSkipVerify = %+v",
						gotTLSConfig.InsecureSkipVerify, tt.want.InsecureSkipVerify)
				}
				if !reflect.DeepEqual(gotTLSConfig.CipherSuites, tt.want.CipherSuites) {
					t.Errorf("GetHTTPClientObj() gotTLSConfig.CipherSuites = %+v, wantTLSConfig.CipherSuites = %+v",
						gotTLSConfig.CipherSuites, tt.want.CipherSuites)
				}
				if gotTLSConfig.PreferServerCipherSuites != tt.want.PreferServerCipherSuites {
					t.Errorf("GetHTTPClientObj() gotTLSConfig.PreferServerCipherSuites = %+v, wantTLSConfig.PreferServerCipherSuites = %+v",
						gotTLSConfig.PreferServerCipherSuites, tt.want.PreferServerCipherSuites)
				}
			}
		})
	}
}

func TestGetHTTPServerObj(t *testing.T) {
	SetUpMockConfig(t)

	httpConf := &HTTPConfig{
		Certificate:   &Data.APIGatewayConf.Certificate,
		PrivateKey:    &Data.APIGatewayConf.PrivateKey,
		CACertificate: &Data.KeyCertConf.RootCACertificate,
	}

	tlsConfig := &tls.Config{}
	httpConf.LoadCertificates(tlsConfig)
	Server.SetTLSConfig(tlsConfig)

	tests := []struct {
		name    string
		exec    modifyTLSData
		want    *tls.Config
		wanterr bool
	}{
		{
			name:    "Positive Case",
			exec:    nil,
			want:    tlsConfig,
			wanterr: false,
		},
		{
			name: "Verify Peer set",
			exec: func(tlsConfig *tls.Config) {
				SetVerifyPeer(true)
				Server.SetTLSConfig(tlsConfig)
				tlsConfig.InsecureSkipVerify = false
			},
			want:    tlsConfig,
			wanterr: false,
		},
		{
			name: "Ciphers not configured",
			exec: func(tlsConfig *tls.Config) {
				SetPreferredCipherSuites([]string{})
				Server.SetTLSConfig(tlsConfig)
				tlsConfig.CipherSuites = DefaultCipherSuiteList
				tlsConfig.PreferServerCipherSuites = true
			},
			want:    tlsConfig,
			wanterr: false,
		},
		{
			name: "Invalid certificates",
			exec: func(tlsConfig *tls.Config) {
				httpConf.Certificate = &nonX509Certificate
			},
			want:    nil,
			wanterr: true,
		},
	}
	for _, tt := range tests {
		if tt.exec != nil {
			tt.exec(tt.want)
		}
		t.Run(tt.name, func(t *testing.T) {
			httpServer, err := httpConf.GetHTTPServerObj()
			if (err != nil) != tt.wanterr {
				t.Errorf("GetHTTPServerObj() err = %v, wanterr %v", err, tt.wanterr)
			}
			if httpServer == nil && tt.want != nil {
				t.Errorf("GetHTTPServerObj() : Expected valid HTTP server object")
				return
			}
			if tt.want != nil {
				gotTLSConfig := httpServer.TLSConfig
				if gotTLSConfig == nil {
					t.Errorf("GetHTTPServerObj() : Expected HTTP server object to contain valid tls.Config")
					return
				}
				if gotTLSConfig.MinVersion != tt.want.MinVersion {
					t.Errorf("GetHTTPServerObj() gotTLSConfig.MinVersion = %+v, wantTLSConfig.MinVersion = %+v",
						gotTLSConfig.MinVersion, tt.want.MinVersion)
				}
				if gotTLSConfig.MaxVersion != tt.want.MaxVersion {
					t.Errorf("GetHTTPServerObj() gotTLSConfig.MaxVersion = %+v, wantTLSConfig.MaxVersion = %+v",
						gotTLSConfig.MaxVersion, tt.want.MaxVersion)
				}
				if gotTLSConfig.InsecureSkipVerify != tt.want.InsecureSkipVerify {
					t.Errorf("GetHTTPServerObj() gotTLSConfig.InsecureSkipVerify = %+v, wantTLSConfig.InsecureSkipVerify = %+v",
						gotTLSConfig.InsecureSkipVerify, tt.want.InsecureSkipVerify)
				}
				if !reflect.DeepEqual(gotTLSConfig.CipherSuites, tt.want.CipherSuites) {
					t.Errorf("GetHTTPServerObj() gotTLSConfig.CipherSuites = %+v, wantTLSConfig.CipherSuites = %+v",
						gotTLSConfig.CipherSuites, tt.want.CipherSuites)
				}
				if gotTLSConfig.PreferServerCipherSuites != tt.want.PreferServerCipherSuites {
					t.Errorf("GetHTTPServerObj() gotTLSConfig.PreferServerCipherSuites = %+v, wantTLSConfig.PreferServerCipherSuites = %+v",
						gotTLSConfig.PreferServerCipherSuites, tt.want.PreferServerCipherSuites)
				}
			}
		})
	}
}

func TestSetDefaultTLSConf(t *testing.T) {
	configuredTLSMinVersion = uint16(0)
	configuredTLSMaxVersion = uint16(0)
	tests := []struct {
		name string
	}{
		{
			name: "Set default TLS versions",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetDefaultTLSConf()
			if configuredTLSMinVersion != DefaultTLSMinVersion {
				t.Errorf("SetDefaultTLSConf() configuredTLSMinVersion = %v want = %v", configuredTLSMinVersion, DefaultTLSMinVersion)
			}
			if configuredTLSMaxVersion != DefaultTLSMaxVersion {
				t.Errorf("SetDefaultTLSConf() configuredTLSMaxVersion = %v want = %v", configuredTLSMaxVersion, DefaultTLSMaxVersion)
			}
			if verifyPeer != DefaultTLSServerVerify {
				t.Errorf("SetDefaultTLSConf() verifyPeer = %v want = %v", verifyPeer, DefaultTLSServerVerify)
			}
			if !reflect.DeepEqual(configuredCipherSuiteList, DefaultCipherSuiteList) {
				t.Errorf("SetPreferredCipherSuites() configuredCipherSuiteList = %v want = %v",
					configuredCipherSuiteList, DefaultCipherSuiteList)
			}
		})
	}
}

func TestSetTLSMinVersion(t *testing.T) {
	tests := []struct {
		name    string
		arg     string
		want    uint16
		wanterr bool
	}{
		{
			name:    "TLS Min Verison not set",
			arg:     "",
			want:    DefaultTLSMinVersion,
			wanterr: false,
		},
		{
			name:    "Valid TLS Min Version",
			arg:     "TLS_1.2",
			want:    tls.VersionTLS12,
			wanterr: false,
		},
		{
			name:    "Invalid TLS Min Version",
			arg:     "TLS",
			want:    0,
			wanterr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SetTLSMinVersion(tt.arg)
			if (err != nil) != tt.wanterr {
				t.Errorf("SetTLSMinVersion() err = %v tt.wanterr = %v", err, tt.wanterr)
			}
			if configuredTLSMinVersion != tt.want {
				t.Errorf("SetTLSMinVersion() configuredTLSMinVersion = %v tt.want = %v", configuredTLSMinVersion, tt.want)
			}
		})
	}
}

func TestSetTLSMaxVersion(t *testing.T) {
	tests := []struct {
		name    string
		arg     string
		want    uint16
		wanterr bool
	}{
		{
			name:    "TLS Max Verison not set",
			arg:     "",
			want:    DefaultTLSMaxVersion,
			wanterr: false,
		},
		{
			name:    "Valid TLS Max Version",
			arg:     "TLS_1.2",
			want:    tls.VersionTLS12,
			wanterr: false,
		},
		{
			name:    "Invalid TLS Max Version",
			arg:     "TLS",
			want:    0,
			wanterr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SetTLSMaxVersion(tt.arg)
			if (err != nil) != tt.wanterr {
				t.Errorf("SetTLSMaxVersion() err = %v tt.wanterr = %v", err, tt.wanterr)
			}
			if configuredTLSMaxVersion != tt.want {
				t.Errorf("SetTLSMaxVersion() configuredTLSMaxVersion = %v tt.want = %v", configuredTLSMaxVersion, tt.want)
			}
		})
	}
}

func TestValidateConfiguredTLSVersions(t *testing.T) {
	SetDefaultTLSConf()
	tests := []struct {
		name    string
		exec    func()
		wantMin uint16
		wantMax uint16
		wanterr bool
	}{
		{
			name: "TLS Min version is lesser than default",
			exec: func() {
				configuredTLSMinVersion = tls.VersionTLS10
			},
			wantMin: tls.VersionTLS10,
			wantMax: DefaultTLSMaxVersion,
			wanterr: false,
		},
		{
			name: "TLS Min version is higher than default",
			exec: func() {
				configuredTLSMinVersion = tls.VersionTLS13
			},
			wantMin: DefaultTLSMinVersion,
			wantMax: DefaultTLSMaxVersion,
			wanterr: false,
		},
		{
			name: "TLS Max version is lesser than default",
			exec: func() {
				configuredTLSMinVersion = tls.VersionTLS10
				configuredTLSMaxVersion = tls.VersionTLS10
			},
			wantMin: tls.VersionTLS10,
			wantMax: tls.VersionTLS10,
			wanterr: false,
		},
		{
			name: "TLS Max version is lesser than Min version",
			exec: func() {
				configuredTLSMinVersion = tls.VersionTLS12
				configuredTLSMaxVersion = tls.VersionTLS11
			},
			wantMin: tls.VersionTLS12,
			wantMax: tls.VersionTLS11,
			wanterr: true,
		},
		{
			name: "TLS Max version is higher than default",
			exec: func() {
				configuredTLSMaxVersion = tls.VersionTLS13
			},
			wantMin: tls.VersionTLS12,
			wantMax: DefaultTLSMaxVersion,
			wanterr: false,
		},
	}
	for _, tt := range tests {
		if tt.exec != nil {
			tt.exec()
		}
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateConfiguredTLSVersions()
			if (err != nil) != tt.wanterr {
				t.Errorf("ValidateConfiguredTLSVersions() goterr = %v, wanterr = %v", err, tt.wanterr)
			}
			if configuredTLSMinVersion != tt.wantMin {
				t.Errorf("ValidateConfiguredTLSVersions() reqMinVersion = %v, setMinVersion= %v", tt.wantMin, configuredTLSMinVersion)
			}
			if configuredTLSMaxVersion != tt.wantMax {
				t.Errorf("ValidateConfiguredTLSVersions() reqMaxVersion = %v, setMaxVersion = %v", tt.wantMax, configuredTLSMaxVersion)
			}
		})
	}
}

func TestSetPreferredCipherSuites(t *testing.T) {
	tests := []struct {
		name    string
		exec    func()
		arg     []string
		want    []uint16
		wanterr bool
	}{
		{
			name:    "TLS Cipher not set",
			arg:     []string{},
			want:    DefaultCipherSuiteList,
			wanterr: false,
		},
		{
			name: "Valid TLS Cipher Set",
			exec: func() {
				configuredCipherSuiteList = []uint16{}
			},
			arg:     []string{"TLS_RSA_WITH_AES_256_GCM_SHA384"},
			want:    []uint16{tls.TLS_RSA_WITH_AES_256_GCM_SHA384},
			wanterr: false,
		},
		{
			name: "Invalid TLS Cipher Set",
			exec: func() {
				configuredCipherSuiteList = []uint16{}
			},
			arg:     []string{"invalid"},
			want:    []uint16{},
			wanterr: true,
		},
	}
	for _, tt := range tests {
		if tt.exec != nil {
			tt.exec()
		}
		t.Run(tt.name, func(t *testing.T) {
			err := SetPreferredCipherSuites(tt.arg)
			if (err != nil) != tt.wanterr {
				t.Errorf("SetPreferredCipherSuites() err = %v tt.wanterr = %v", err, tt.wanterr)
			}
			if !reflect.DeepEqual(configuredCipherSuiteList, tt.want) {
				t.Errorf("SetPreferredCipherSuites() configuredCipherSuiteList = %v tt.want = %v", configuredCipherSuiteList, tt.want)
			}
		})
	}
}
