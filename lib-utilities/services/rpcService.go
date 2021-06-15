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
	"net"
	"time"

	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/coreos/etcd/clientv3"
	micro "github.com/micro/go-micro"
	"github.com/micro/go-micro/transport"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type serviceType int

const (
	serverService serviceType = iota
	clientService
)

// odimService holds the components for bringing up and communicating with a micro service
type odimService struct {
	clientTransportCreds credentials.TransportCredentials
	etcdTLSConfig        *tls.Config
	registryAddress      string
	server               *grpc.Server
	serverAddress        string
	serverName           string
	serverTransportCreds credentials.TransportCredentials
}

// ODIMService holds the initialized instance of odimService
var ODIMService odimService

// Service holds the microservice instance
var Service micro.Service

// InitializeService will initialize a new micro service with the selected framework.
func InitializeService(serviceName string) error {
	switch config.CLArgs.FrameWork {
	case "GRPC":
		err := ODIMService.Init(serviceName)
		if err != nil {
			return fmt.Errorf("While trying to initiate ODIMService model, got: %v", err)
		}
		err = ODIMService.registerService()
		if err != nil {
			return fmt.Errorf("While trying to register the service in the registry, got: %v", err)
		}

	case "GOMICRO":
		tlsConfig, err := getGoMicroTLSConfig()
		if err != nil {
			return fmt.Errorf("Failed to load TLS config for go micro: %v", err)
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
	}
	return nil
}

// InitializeClient will initialize a client for micro service communication.
func InitializeClient(serviceName string) error {
	err := ODIMService.Init(serviceName)
	if err != nil {
		return fmt.Errorf("While trying to initiate ODIMService model, got: %v", err)
	}
	return nil
}

// Server returns the gRPC server, which helps in bringing up the gRPC microservices
func (s *odimService) Server() *grpc.Server {
	return s.server
}

// Client will return the gRPC client connection for the requested service.
// Function will get the service address from the service registry
// with the name privided for establishing connection.
// IMPORTANT: the connection created with this function must be close by the user.
// usage:
// conn, err := ODIMService.Client(AccountSession)
// defer conn.Close()
func (s *odimService) Client(clientName string) (*grpc.ClientConn, error) {
	clientAddress, err := s.getServiceAddress(clientName)
	if err != nil {
		return nil, fmt.Errorf("While trying to get the service address from registry, got: %v", err)
	}

	err = s.loadTLSCredentials(clientService)
	if err != nil {
		return nil, fmt.Errorf("Failed to load TLS credentials: %v", err)
	}

	return grpc.Dial(
		clientAddress,
		grpc.WithTransportCredentials(s.clientTransportCreds),
	)
}

// Run will make the gRPC microservice up and running
func (s *odimService) Run() error {
	l, err := net.Listen("tcp", s.serverAddress)
	if err != nil {
		return fmt.Errorf("While trying to get listen for the grpc, got: %v", err)
	}
	s.server.Serve(l)
	return nil
}

// Init initializes the ODIMService with server and client TLS, server and registry details etc.
// It also initialize ODIMService.server which will help in bring up a microservice
func (s *odimService) Init(serviceName string) error {
	s.serverName = serviceName
	s.registryAddress = config.CLArgs.RegistryAddress
	if s.registryAddress == "" {
		return fmt.Errorf("RegistryAddress not found")
	}
	s.serverAddress = config.CLArgs.ServerAddress
	if s.serverAddress == "" && s.serverName != APIClient {
		return fmt.Errorf("ServerAddress not found")
	}
	err := s.loadTLSCredentials(serverService)
	if err != nil {
		return fmt.Errorf("While trying to setup TLS transport layer for gRPC client, got: %v", err)
	}
	err = s.loadTLSCredentials(clientService)
	if err != nil {
		return fmt.Errorf("While trying to setup TLS transport layer for gRPC client, got: %v", err)
	}
	s.etcdTLSConfig, err = getTLSConfig()
	if err != nil {
		return fmt.Errorf("While trying to get tls for etcd, got: %v", err)
	}
	ODIMService.server = grpc.NewServer(
		grpc.Creds(s.serverTransportCreds),
	)
	return nil
}

func (s *odimService) getServiceAddress(serviceName string) (string, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{s.registryAddress},
		DialTimeout: 5 * time.Second,
		TLS:         s.etcdTLSConfig,
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

func (s *odimService) registerService() error {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{s.registryAddress},
		DialTimeout: 5 * time.Second,
		TLS:         s.etcdTLSConfig,
	})
	if err != nil {
		return fmt.Errorf("While trying to create registry client, got: %v", err)
	}
	defer cli.Close()
	kv := clientv3.NewKV(cli)
	_, err = kv.Put(context.TODO(), s.serverName, s.serverAddress)
	if err != nil {
		return fmt.Errorf("While trying to register the service, got: %v", err)
	}
	return nil
}

func (s *odimService) loadTLSCredentials(st serviceType) error {
	switch st {
	case serverService:
		tlsConfig, err := loadServerTLSConfig()
		if err != nil {
			return fmt.Errorf("While trying to load server TLS config, got: %v", err)
		}
		s.serverTransportCreds = credentials.NewTLS(tlsConfig)
	case clientService:
		if s.clientTransportCreds == nil {
			tlsConfig, err := loadClientTLSConfig()
			if err != nil {
				return fmt.Errorf("While trying to load client TLS config, got: %v", err)
			}
			s.clientTransportCreds = credentials.NewTLS(tlsConfig)
		}
	}
	return nil
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

func getGoMicroTLSConfig() (*tls.Config, error) {
	goMicroTLS, err := getTLSConfig()
	if err != nil {
		return nil, fmt.Errorf("Failed to load certificates for GoMicro: %v", err)
	}
	goMicroTLS.ServerName = config.Data.LocalhostFQDN
	return goMicroTLS, nil
}

func getTLSConfig() (*tls.Config, error) {
	serverTLS, err := loadServerTLSConfig()
	if err != nil {
		return nil, fmt.Errorf("Failed to load server tls: %v", err)
	}
	clientTLS, err := loadClientTLSConfig()
	if err != nil {
		return nil, fmt.Errorf("Failed to load client tls: %v", err)
	}
	return &tls.Config{
		RootCAs:      clientTLS.RootCAs,
		Certificates: serverTLS.Certificates,
	}, nil
}
