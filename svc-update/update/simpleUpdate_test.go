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
package update

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"

	updateproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/update"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
)

func TestSimpleUpdate(t *testing.T) {
	errMsg := []string{"/redfish/v1/Systems/uuid:/target1"}
	errArg1 := response.Args{
		Code:    response.GeneralError,
		Message: "",
		ErrorArgs: []response.ErrArgs{
			response.ErrArgs{
				StatusMessage: response.ResourceNotFound,
				ErrorMessage:  "error: SystemUUID not found",
				MessageArgs:   []interface{}{"System", fmt.Sprintf("%v", errMsg)},
			},
		},
	}

	request1 := []byte(`{"ImageURI":"abc","Targets":["/redfish/v1/Systems/uuid:/target1"],"@Redfish.OperationApplyTimeSupport": {"@odata.type": "#Settings.v1_2_0.OperationApplyTimeSupport","SupportedValues": ["OnStartUpdate"]}}`)
	request3 := []byte(`{"ImageURI":"abc","Targets":["/redfish/v1/Systems/uuid:1/target1"],"@Redfish.OperationApplyTimeSupport": {"@odata.type": "#Settings.v1_2_0.OperationApplyTimeSupport","SupportedValues": ["OnStartUpdate"]}}`)
	output := map[string]interface{}{"Attributes": "sample"}
	tests := []struct {
		name string
		req  *updateproto.UpdateRequest
		want response.RPC
	}{
		{
			name: "uuid without system id",
			req: &updateproto.UpdateRequest{
				RequestBody:  request1,
				SessionToken: "token",
			},
			want: response.RPC{
				StatusCode:    http.StatusBadRequest,
				StatusMessage: response.ResourceNotFound,
				Header: map[string]string{
					"Content-type": "application/json; charset=utf-8",
				},
				Body: errArg1.CreateGenericErrorResponse(),
			},
		},
		{
			name: "Valid Request",
			req: &updateproto.UpdateRequest{
				RequestBody:  request3,
				SessionToken: "token",
			},
			want: response.RPC{
				StatusCode:    http.StatusOK,
				StatusMessage: response.Success,
				Header: map[string]string{
					"Cache-Control":     "no-cache",
					"Connection":        "keep-alive",
					"Content-type":      "application/json; charset=utf-8",
					"Transfer-Encoding": "chunked",
					"OData-Version":     "4.0",
				},
				Body: output,
			},
		},
	}
	e := mockGetExternalInterface()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := e.SimpleUpdate(tt.req); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SimpleUpdate() = %v, want %v", got, tt.want)
			}
		})
	}
}
