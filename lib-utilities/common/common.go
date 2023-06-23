// (C) Copyright [2023] Hewlett Packard Enterprise Development LP
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.
// Package common ...
package common

import (
	"fmt"
	"net"
)

// GetIPFromHostName - look up the ip from the fqdn
func GetIPFromHostName(fqdn string) (string, error) {
	host, _, err := net.SplitHostPort(fqdn)
	if err != nil {
		host = fqdn
	}
	addr, err := net.LookupIP(host)
	if err != nil || len(addr) < 1 {
		errorMessage := "Can't lookup the ip from host name"
		if err != nil {
			errorMessage = "Can't lookup the ip from host name " + err.Error()
		}
		return "", fmt.Errorf(errorMessage)
	}
	return fmt.Sprintf("%v", addr[0]), nil
}
