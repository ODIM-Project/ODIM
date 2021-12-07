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

package fabrics

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	fabricsproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/fabrics"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
)

func TestFabrics_DeleteFabricResource(t *testing.T) {
	Token.Tokens = make(map[string]string)
	common.SetUpMockConfig()
	defer func() {
		err := common.TruncateDB(common.OnDisk)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()
	err := mockPluginData(t, "CFM", "XAuthToken", "9091")
	if err != nil {
		t.Fatalf("Error in creating mock DeviceData :%v", err)
	}
	err = mockFabricData("d72dade0-c35a-984c-4859-1108132d72da", "CFM")
	if err != nil {
		t.Fatalf("Error in creating mockFabricData :%v", err)
	}
	errResp1 := response.Args{
		Code:    response.GeneralError,
		Message: "",
		ErrorArgs: []response.ErrArgs{
			response.ErrArgs{
				StatusMessage: response.InsufficientPrivilege,
				ErrorMessage:  "error while trying to authenticate session",
				MessageArgs:   []interface{}{},
			},
		},
	}
	errResp2 := response.Args{
		Code:    response.GeneralError,
		Message: "",
		ErrorArgs: []response.ErrArgs{
			response.ErrArgs{
				StatusMessage: response.NoValidSession,
				ErrorMessage:  "error while trying to authenticate session",
				MessageArgs:   []interface{}{},
			},
		},
	}

	type args struct {
		req *fabricsproto.FabricRequest
	}
	tests := []struct {
		name string
		f    *Fabrics
		args args
		want response.RPC
	}{
		{
			name: "positive case",
			f: &Fabrics{
				Auth:          mockAuth,
				ContactClient: mockContactClient,
			},
			args: args{
				req: &fabricsproto.FabricRequest{
					SessionToken: "valid",
					URL:          "/redfish/v1/Fabrics/d72dade0-c35a-984c-4859-1108132d72da",
					Method:       "DELETE",
				},
			},
			want: response.RPC{
				StatusCode:    http.StatusOK,
				StatusMessage: response.Success,
				Header: map[string]string{
					"Allow":             `"GET", "PUT", "POST", "PATCH", "DELETE"`,
				},
			},
		}, {
			name: "insufficient privilege",
			f: &Fabrics{
				Auth: mockAuth,
			},
			args: args{
				req: &fabricsproto.FabricRequest{
					SessionToken: "sometoken",
					Method:       "DELETE",
				},
			},
			want: response.RPC{
				StatusCode:    http.StatusForbidden,
				StatusMessage: response.InsufficientPrivilege,
				Body:          errResp1.CreateGenericErrorResponse(),
			},
		},
		{
			name: "invalid token",
			f: &Fabrics{
				Auth: mockAuth,
			},
			args: args{
				req: &fabricsproto.FabricRequest{
					SessionToken: "invalid",
					Method:       "DELETE",
				},
			},
			want: response.RPC{
				StatusCode:    http.StatusUnauthorized,
				StatusMessage: response.NoValidSession,
				Body:          errResp2.CreateGenericErrorResponse(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.DeleteFabricResource(tt.args.req); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Fabrics.DeleteFabricResource() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFabrics_DeleteFabricResourceWithNoValidSession(t *testing.T) {
	Token.Tokens = make(map[string]string)
	Token.Tokens["CFM"] = "234556"
	common.SetUpMockConfig()
	defer func() {
		err := common.TruncateDB(common.OnDisk)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()
	err := mockPluginData(t, "CFM", "XAuthToken", "9091")
	if err != nil {
		t.Fatalf("Error in creating mock DeviceData :%v", err)
	}
	err = mockFabricData("d72dade0-c35a-984c-4859-1108132d72da", "CFM")
	if err != nil {
		t.Fatalf("Error in creating mockFabricData :%v", err)
	}

	type args struct {
		req *fabricsproto.FabricRequest
	}
	tests := []struct {
		name string
		f    *Fabrics
		args args
		want response.RPC
	}{
		{
			name: "positive case",
			f: &Fabrics{
				Auth:          mockAuth,
				ContactClient: mockContactClient,
			},
			args: args{
				req: &fabricsproto.FabricRequest{
					SessionToken: "valid",
					URL:          "/redfish/v1/Fabrics/d72dade0-c35a-984c-4859-1108132d72da",
					Method:       "DELETE",
				},
			},
			want: response.RPC{
				StatusCode:    http.StatusOK,
				StatusMessage: response.Success,
				Header: map[string]string{
					"Allow":             `"GET", "PUT", "POST", "PATCH", "DELETE"`,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.DeleteFabricResource(tt.args.req); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Fabrics.GetFabricResource() = %v, want %v", got, tt.want)
			}
		})
	}
}
