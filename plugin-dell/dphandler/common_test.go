//(C) Copyright [2019] Hewlett Packard Enterprise Development LP
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

// Package dphandler ...
package dphandler

import (
	"net/http"
	"testing"

	"github.com/ODIM-Project/ODIM/plugin-dell/config"
	"github.com/ODIM-Project/ODIM/plugin-dell/dputilities"
	"github.com/stretchr/testify/assert"
)

func Test_convertToNorthBoundURI(t *testing.T) {
	type args struct {
		req             string
		storageInstance string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Positive test case ",
			args: args{
				req: "/ODIM/PCIeDevice",
			},
			want: "/ODIM/PCIeDevice",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertToNorthBoundURI(tt.args.req, tt.args.storageInstance); got != tt.want {
				t.Errorf("convertToNorthBoundURI() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_queryDevice(t *testing.T) {
	config.SetUpMockConfig(t)

	device := dputilities.RedfishDevice{
		Username: "Localhost",
		Password: "password",
	}
	ctxt := mockContext()
	config.Data.KeyCertConf.RootCACertificate = nil
	statusCode, _, _, err := queryDevice(ctxt, "/ODIM", &device, http.MethodPost)
	assert.NotNil(t, err, "should contain error code StatusInternalServerError")
	assert.Equal(t, http.StatusInternalServerError, statusCode)

	device = dputilities.RedfishDevice{
		Username: "Localhost",
		Password: "password",
	}
	config.SetUpMockConfig(t)
	statusCode, _, _, err = queryDevice(ctxt, "/ODIM", &device, http.MethodPost)
	assert.NotNil(t, err, "should contain error code StatusBadRequest")
	assert.Equal(t, http.StatusBadRequest, statusCode)

}
