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

package config

import (
	"flag"

	log "github.com/sirupsen/logrus"
)

// cliModel holds the data passed as the command line argument
type clModel struct {
	ClientRequestTimeout string
	Registry             string
	RegistryAddress      string
	ServerAddress        string
	FrameWork            string
}

// CLArgs is for accessing the data passed as the command line argument
var CLArgs clModel

// CollectCLArgs will collect the command line arguments and fill CLArgs variable of type clModel
func CollectCLArgs() {
	flag.StringVar(&CLArgs.ClientRequestTimeout, "client_request_timeout", "", "maximum request time which client waits")
	flag.StringVar(&CLArgs.Registry, "registry", "", "service registry")
	flag.StringVar(&CLArgs.RegistryAddress, "registry_address", "", "address of the registry")
	flag.StringVar(&CLArgs.ServerAddress, "server_address", "", "address for the micro service")
	flag.StringVar(&CLArgs.FrameWork, "framework", "GRPC", "framework used for micro service communication")
	flag.Parse()
	if CLArgs.RegistryAddress == "" {
		log.Warn("No CLI argument found for registry_address")
	}
	if CLArgs.ServerAddress == "" {
		log.Warn("No CLI argument found for server_address")
	}
}
