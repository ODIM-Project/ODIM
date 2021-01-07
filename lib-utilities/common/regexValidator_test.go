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

// Package common ...
package common

import (
	"testing"
)

func TestURIValidator(t *testing.T) {

	tests := []struct {
		name string
		uri  string
		want bool
	}{
		{
			name: "Valid one word hostname",
			uri:  "https://URP:45003/EventService/Events",
			want: true,
		},
		{
			name: "Valid localhost",
			uri:  "https://localhost/EventService/Events",
			want: true,
		},
		{
			name: "Valid URI with IPv4 address",
			uri:  "https://1.2.3.4:5678/redfish/v1",
			want: true,
		},
		{
			name: "Valid URI with IPv6 address in square brackets",
			uri:  "https://[fe80::1234]:5678/redfish/v1",
			want: true,
		},
		{
			name: "Valid URI with fqdn",
			uri:  "https://odim.net/redfish/v1",
			want: true,
		},
		{
			name: "Valid URI with fqdn with port",
			uri:  "https://odim.net:12345/redfish/v1",
			want: true,
		},
		{
			name: "Invalid URI with just endpoint",
			uri:  "/redfish/v1",
			want: false,
		},
		{
			name: "Invalid URI missing IP and Port",
			uri:  "https:///redfish/v1",
			want: false,
		},
		{
			name: "Invalid URI with wrong IPv6 address",
			uri:  "https://[[fe80:1234]:5678/redfish/v1",
			want: false,
		},
		{
			name: "Invalid URI with http as scheme",
			uri:  "http://[fe80::1234]:5678/redfish/v1",
			want: false,
		},
		{
			name: "Invalid URI with wrong IPv4 address",
			uri:  "http://10.20:5678/redfish/v1",
			want: false,
		},
		{
			name: "URI with IPv4 address and invalid port",
			uri:  "https://127.0.0.1:/redfish/v1",
			want: false,
		},
		{
			name: "URI with fqdn and invalid port",
			uri:  "https://odim.net:123456/redfish/v1",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := URIValidator(tt.uri)
			if got != tt.want {
				t.Errorf("TestURIValidator got = %v, want %v", got, tt.want)
			}
		})
	}
}
