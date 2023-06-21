//(C) Copyright [2023] Hewlett Packard Enterprise Development LP
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

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func TestGetIPFromHostName(t *testing.T) {
	type args struct {
		fqdn string
	}
	tests := []struct {
		name  string
		args  args
		want  []string
		want1 bool
	}{
		{
			name: "success case",
			args: args{
				fqdn: "localhost:8080",
			},
			want:  []string{"127.0.0.1", "::1"},
			want1: false,
		},
		{
			name: "success case 2",
			args: args{
				fqdn: "127.0.0.1",
			},
			want:  []string{"127.0.0.1", "::1"},
			want1: false,
		},
		{
			name: "success case 2",
			args: args{
				fqdn: "invalid",
			},
			want:  []string{""},
			want1: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := GetIPFromHostName(tt.args.fqdn)
			if !contains(tt.want, got) {
				t.Errorf("GetIPFromHostName() got = %v, want one of %v", got, tt.want)
			}
			if (got1 != nil) != tt.want1 {
				t.Errorf("GetIPFromHostName() got = %v, want %v", got1 != nil, tt.want1)
			}
		})
	}
}
