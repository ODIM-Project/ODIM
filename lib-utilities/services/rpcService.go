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

package services

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"github.com/bharath-b-hpe/odimra/lib-utilities/config"
	micro "github.com/micro/go-micro"
	"github.com/micro/go-micro/transport"
)

// Service holds the microservice instance
var Service micro.Service

// InitializeService will initialize a new micro.Service.
// Service will be initialized here itself, so the Server() and Client()
// called easily.
func InitializeService(serviceName string) error {

	cer, err := tls.X509KeyPair(
		config.Data.KeyCertConf.RPCCertificate,
		config.Data.KeyCertConf.RPCPrivateKey,
	)
	if err != nil {
		return fmt.Errorf("error while trying to load x509 key pair: %v", err)
	}

	certPool := x509.NewCertPool()

	block, _ := pem.Decode(config.Data.KeyCertConf.RootCACertificate)
	if block == nil {
		return fmt.Errorf("error while decoding ca file")
	}
	if block.Type != "CERTIFICATE" || len(block.Headers) != 0 {
		return fmt.Errorf("error while decoding ca file")
	}

	certificate, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return fmt.Errorf("error while ParseCertificate ca block file: %v", err)
	}

	certPool.AddCert(certificate)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cer},
		RootCAs:      certPool,
		ServerName:   config.Data.LocalhostFQDN,
	}
	config.Server.SetTLSConfig(tlsConfig)

	Service = micro.NewService(
		micro.Name(serviceName),
		micro.Transport(
			transport.NewTransport(
				transport.Secure(true),
				transport.TLSConfig(tlsConfig),
			),
		),
	)
	Service.Init()
	return nil

}
