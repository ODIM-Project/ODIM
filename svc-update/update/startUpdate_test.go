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
	"net/http"
	"reflect"
	"testing"

	updateproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/update"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
)

func TestStartUpdate(t *testing.T) {

	output := map[string]interface{}{"Attributes": "sample"}
	tests := []struct {
		name string
		req  *updateproto.UpdateRequest
		want response.RPC
	}{
		{
			name: "Valid Request",
			req: &updateproto.UpdateRequest{
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
			if got := e.StartUpdate(tt.req); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StartUpdate() = %v, want %v", got, tt.want)
			}
		})
	}
}
