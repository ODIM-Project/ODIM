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
//under the License.

// Package config ...
package config

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net"
	"net/http"
	"sync"
	"time"
)

// HTTPConfig is for passing required info to create a http.Server object
type HTTPConfig struct {
	// Certificate contains the certifcate data to be loaded
	Certificate *[]byte
	// PrivateKey contains the private key data to be loaded
	PrivateKey *[]byte
	// CACertificate contains the CA certificate data to be loaded
	CACertificate *[]byte
	// ServerAddress contains the IP/FQDN address of the server
	ServerAddress string
	// ServerPort contains the port of the server
	ServerPort string
	// loadCertificates is for marking to load CA cert only or not
	loadCertificates bool
}

var (
	// TLSConfMutex is used for avoiding race conditions
	TLSConfMutex = &sync.RWMutex{}
	// configuredCipherSuiteList contains the list of configured cipher suites
	configuredCipherSuiteList = make([]uint16, 0)
	// configuredTLSMinVersion is the configured TLS minimum version
	configuredTLSMinVersion uint16
	// configuredTLSMaxVersion is the configured TLS maximum version
	configuredTLSMaxVersion uint16
	// verifyPeer is for verify peers(Server/Client)
	verifyPeer bool
	// DefaultHTTPClient is the global reusable client instance for contacting a server
	DefaultHTTPClient = &http.Client{
		Timeout: time.Duration(DefaultHTTPConnTimeout) * time.Second,
	}
	// DefaultHTTPTransport is the global resuable tranport instance for contacting a server
	DefaultHTTPTransport = &http.Transport{
		MaxIdleConns:          DefaultHTTPMaxIdleConns,
		IdleConnTimeout:       time.Duration(DefaultHTTPIdleConnTimeout) * time.Second,
		TLSHandshakeTimeout:   time.Duration(DefaultTLSHandShakeTimeout) * time.Second,
		DisableKeepAlives:     !DefaultHTTPUseKeepAlive,
		MaxIdleConnsPerHost:   DefaultHTTPMaxIdleConnPerHost,
		ExpectContinueTimeout: time.Duration(DefaultHTTPExpectContinueTimeout) * time.Second,
	}
)

// GetHTTPClientObj is for obtaining a client instance for making http(s) queries
func (config *HTTPConfig) GetHTTPClientObj() (*http.Client, error) {
	tlsConfig := &tls.Config{}
	if err := config.LoadCertificates(tlsConfig); err != nil {
		return nil, err
	}
	if DefaultHTTPClient == nil || DefaultHTTPClient.Transport == nil || DefaultHTTPTransport.TLSClientConfig == nil {
		TLSConfMutex.Lock()
		Client.SetTLSConfig(tlsConfig)
		DefaultHTTPTransport.TLSClientConfig = tlsConfig
		DefaultHTTPClient.Transport = DefaultHTTPTransport
		TLSConfMutex.Unlock()
	}
	return DefaultHTTPClient, nil
}

// GetHTTPServerObj is for obtaining a server instance to start a service using iris helper
func (config *HTTPConfig) GetHTTPServerObj() (*http.Server, error) {
	config.loadCertificates = true
	tlsConfig := &tls.Config{}
	if err := config.LoadCertificates(tlsConfig); err != nil {
		return nil, err
	}
	Server.SetTLSConfig(tlsConfig)

	return &http.Server{
		Addr:      net.JoinHostPort(config.ServerAddress, config.ServerPort),
		TLSConfig: tlsConfig,
	}, nil
}

// LoadCertificates is for including passed certificates in tls.Config
func (config *HTTPConfig) LoadCertificates(tlsConfig *tls.Config) error {
	// for client mode interaction certificates will not be required and
	// just CA certificate needs to be loaded for server validation
	if config.loadCertificates {
		cert, err := tls.X509KeyPair(*config.Certificate, *config.PrivateKey)
		if err != nil {
			return fmt.Errorf("error: failed to load key pair: %v", err)
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
		tlsConfig.BuildNameToCertificate()
	}

	capool := x509.NewCertPool()
	if !capool.AppendCertsFromPEM(*config.CACertificate) {
		return fmt.Errorf("error: failed to load CA certificate")
	}

	tlsConfig.RootCAs = capool
	tlsConfig.ClientCAs = capool
	return nil
}

// SetTLSConfig is for setting updating common fields of tls.Config
func (host Host) SetTLSConfig(tlsConfig *tls.Config) {
	tlsConfig.MinVersion = configuredTLSMinVersion
	tlsConfig.MaxVersion = configuredTLSMaxVersion
	if !verifyPeer {
		tlsConfig.InsecureSkipVerify = true
	}
	if len(configuredCipherSuiteList) != 0 {
		tlsConfig.CipherSuites = configuredCipherSuiteList
		if Server == host {
			tlsConfig.PreferServerCipherSuites = true
		}
	}
}

// SetDefaultTLSConf is for updating TLS conf with default values
func SetDefaultTLSConf() {
	TLSConfMutex.RLock()
	defer TLSConfMutex.RUnlock()
	verifyPeer = DefaultTLSServerVerify
	configuredTLSMinVersion = DefaultTLSMinVersion
	configuredTLSMaxVersion = DefaultTLSMaxVersion
	configuredCipherSuiteList = DefaultCipherSuiteList
}

// SetVerifyPeer is for updating verifyPeer
func SetVerifyPeer(val bool) {
	verifyPeer = val
}

// SetTLSMinVersion is for setting configuredTLSMinVersion
func SetTLSMinVersion(version string) error {
	if version == "" {
		log.Warn("TLS MinVersion is not provided, setting default value")
		configuredTLSMinVersion = DefaultTLSMinVersion
		return nil
	}
	var valid bool
	if configuredTLSMinVersion, valid = SupportedTLSVersions[version]; !valid {
		return fmt.Errorf("error: invalid TLS MinVersion %s set", version)
	}
	return nil
}

// SetTLSMaxVersion is for setting configuredTLSMaxVersion
func SetTLSMaxVersion(version string) error {
	if version == "" {
		log.Warn("TLS MaxVersion is not provided, setting default value")
		configuredTLSMaxVersion = DefaultTLSMaxVersion
		return nil
	}
	var valid bool
	if configuredTLSMaxVersion, valid = SupportedTLSVersions[version]; !valid {
		return fmt.Errorf("error: invalid TLS MaxVersion %s set", version)
	}
	return nil
}

// SetPreferredCipherSuites is for setting configuredCipherSuiteList
func SetPreferredCipherSuites(cipherList []string) error {

	if len(cipherList) == 0 {
		configuredCipherSuiteList = DefaultCipherSuiteList
		return nil
	}

	for _, cipher := range cipherList {
		goCipher, exist := SupportedCipherSuitesList[cipher]
		if !exist {
			return fmt.Errorf("error: PreferredTLSCipherSuites contains unknown cipher %s", cipher)
		}
		configuredCipherSuiteList = append(configuredCipherSuiteList, goCipher)
	}
	return nil
}

// ValidateConfiguredTLSVersions is for valdiating TLS versions configured
func ValidateConfiguredTLSVersions() error {
	if configuredTLSMinVersion < DefaultTLSMinVersion {
		log.Warn("TLS MinVersion set is lower than suggested version")
	}
	if configuredTLSMinVersion > DefaultTLSMinVersion {
		log.Warn("TLS MinVersion set is higher than supported version, setting default value")
		configuredTLSMinVersion = DefaultTLSMinVersion
	}
	if configuredTLSMaxVersion < configuredTLSMinVersion {
		return fmt.Errorf("error: TLS MaxVersion cannot be lower than MinVersion")
	}
	if configuredTLSMaxVersion > DefaultTLSMaxVersion {
		log.Warn("TLS MaxVersion set is higher than supported version, setting default value")
		configuredTLSMaxVersion = DefaultTLSMaxVersion
	}
	return nil
}
