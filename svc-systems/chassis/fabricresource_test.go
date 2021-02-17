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

package chassis

import (
	"net/http"
	"reflect"
	"testing"

	dmtfmodel "github.com/ODIM-Project/ODIM/lib-dmtf/model"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
)

func Test_fabricFactory_getFabricChassisResource(t *testing.T) {
	Token.Tokens = make(map[string]string)
	config.SetUpMockConfig(t)
	f := getFabricFactoryMock(nil)
	var r response.RPC
	initializeRPCResponse(
		&r,
		dmtfmodel.Chassis{
			ChassisType:  "valid_type",
			SerialNumber: "valid_serial_number",
		},
	)
	type args struct {
		rID string
	}
	tests := []struct {
		name string
		f    *fabricFactory
		args args
		want response.RPC
	}{
		{
			name: "successful GET on fabric chassis resource",
			f:    f,
			args: args{
				rID: "valid",
			},
			want: r,
		},
		{
			name: "GET with invalid resource id",
			f:    f,
			args: args{
				rID: "invalid",
			},
			want: common.GeneralError(http.StatusNotFound, response.ResourceNotFound, "", []interface{}{"Chassis", "invalid"}, nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.getFabricChassisResource(tt.args.rID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fabricFactory.getFabricChassisResource() = %v, want %v", got, tt.want)
			}
		})
	}
}
