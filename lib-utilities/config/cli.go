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
	"fmt"
)

// cliModel holds the data passed as the command line argument
type cliModel struct {
	RegistryAddress string
	ServerAddress   string
}

// CLIData is for accessing the data passed as the command line argument
var CLIData cliModel

func collectCLIData() error {
	flag.StringVar(&CLIData.RegistryAddress, "registry_address", "", "address of the registry")
	if CLIData.RegistryAddress == "" {
		return fmt.Errorf("No CLI argument found for registry_address")
	}
	flag.StringVar(&CLIData.ServerAddress, "server_address", "", "address for the micro service")
	if CLIData.ServerAddress == "" {
		return fmt.Errorf("No CLI argument found for server_address")
	}
	return nil
}
