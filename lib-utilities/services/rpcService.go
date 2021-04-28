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
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"time"

	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/coreos/etcd/clientv3"
	micro "github.com/micro/go-micro"
	"github.com/micro/go-micro/transport"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// TODO: remove line 33 to 87 after the completion of go micro to gRPC migration.

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

// gRPC implementation starts here. TODO: remove this comment after the removal of go micro implementation.

type serviceType int

const (
	serverService serviceType = iota
	clientService
)

var (
	clientTransportCreds credentials.TransportCredentials
	// Server is for bringing up gRPC micro services
	Server *grpc.Server
)

// InitializeMicroService register the micro service and initializes the gRPC client transport
// and server. Function returns error at the failure of server or client transport creation
func InitializeMicroService(serviceName string) error {
	tlsServerCredentials, err := loadTLSCredentials(serverService)
	if err != nil {
		return fmt.Errorf("While trying to setup TLS transport layer for gRPC client, got: %v", err)
	}
	err = registerService(serviceName)
	if err != nil {
		return fmt.Errorf("While trying to register the service in the registry, got: %v", err)
	}

	_, err = loadTLSCredentials(clientService)
	if err != nil {
		return fmt.Errorf("While trying to setup TLS transport layer for gRPC client, got: %v", err)
	}
	Server = grpc.NewServer(
		grpc.Creds(tlsServerCredentials),
	)
	return nil
}

// GetServiceClient will return the gRPC client connection for the requested service.
// Function will get the service address from the service registry
// with the name privided for establishing connection.
// IMPORTANT: the connection created with this function must be close by the user.
// usage:
// conn, err := GetServiceClient(AccountSession)
// defer conn.Close()
func GetServiceClient(serviceName string) (*grpc.ClientConn, error) {
	serviceAddress, err := getServiceAddress(serviceName)
	if err != nil {
		return nil, fmt.Errorf("While trying to get the service address from registry, got: %v", err)
	}

	tlsCredentials, err := loadTLSCredentials(clientService)
	if err != nil {
		return nil, fmt.Errorf("Failed to load TLS credentials: %v", err)
	}

	return grpc.Dial(
		serviceAddress,
		grpc.WithTransportCredentials(tlsCredentials),
	)
}

func getServiceAddress(serviceName string) (string, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{config.CLIData.RegistryAddress},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return "", fmt.Errorf("While trying to create registry client, got: %v", err)
	}
	defer cli.Close()
	kv := clientv3.NewKV(cli)
	resp, err := kv.Get(context.TODO(), serviceName, clientv3.WithPrefix())
	if err != nil {
		return "", fmt.Errorf("While trying to get the service from registry, got: %v", err)
	}
	if len(resp.Kvs) == 0 {
		return "", fmt.Errorf("No service with %s found in the service registry", serviceName)
	}
	return string(resp.Kvs[0].Value), nil
}

func registerService(serviceName string) error {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{config.CLIData.RegistryAddress},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("While trying to create registry client, got: %v", err)
	}
	defer cli.Close()
	kv := clientv3.NewKV(cli)
	_, err = kv.Put(context.TODO(), serviceName, config.CLIData.ServerAddress)
	if err != nil {
		return fmt.Errorf("While trying to register the service, got: %v", err)
	}
	return nil
}

func loadTLSCredentials(st serviceType) (credentials.TransportCredentials, error) {
	switch st {
	case serverService:
		tlsConfig, err := loadServerTLSConfig()
		if err != nil {
			return nil, fmt.Errorf("While trying to load server TLS config, got: %v", err)
		}
		return credentials.NewTLS(tlsConfig), nil
	case clientService:
		if clientTransportCreds == nil {
			tlsConfig, err := loadClientTLSConfig()
			if err != nil {
				return nil, fmt.Errorf("While trying to load client TLS config, got: %v", err)
			}
			clientTransportCreds = credentials.NewTLS(tlsConfig)
		}
		return clientTransportCreds, nil
	}
	return nil, nil
}

func loadServerTLSConfig() (*tls.Config, error) {
	cer, err := tls.X509KeyPair(
		config.Data.KeyCertConf.RPCCertificate,
		config.Data.KeyCertConf.RPCPrivateKey,
	)
	if err != nil {
		return nil, fmt.Errorf("While trying to load x509 key pair, got: %v", err)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{cer},
		ClientAuth:   tls.NoClientCert,
		ServerName:   config.Data.LocalhostFQDN,
	}, nil
}

func loadClientTLSConfig() (*tls.Config, error) {
	certPool := x509.NewCertPool()
	block, _ := pem.Decode(config.Data.KeyCertConf.RootCACertificate)
	if block == nil {
		return nil, fmt.Errorf("Failed in decoding ca file")
	}
	if block.Type != "CERTIFICATE" || len(block.Headers) != 0 {
		return nil, fmt.Errorf("Failed in decoding ca file")
	}
	certificate, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("While ParseCertificate ca block file, got: %v", err)
	}
	certPool.AddCert(certificate)
	return &tls.Config{
		RootCAs:    certPool,
		ServerName: config.Data.LocalhostFQDN,
	}, nil
}
