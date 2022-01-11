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
	var respArgs response.Args
	respArgs = response.Args{
		Code:    response.Success,
		Message: "Request completed successfully",
	}
	body := respArgs.CreateGenericErrorResponse()
	tests := []struct {
		name string
		args args
		want response.RPC
	}{
		{
			name: "Valid Request",
			args: args{
				taskID: "someID", sessionUserName: "someUser",
				req: &updateproto.UpdateRequest{
					SessionToken: "validToken",
				},
			},
			want: response.RPC{
				StatusCode:    http.StatusOK,
				StatusMessage: response.Success,
				Body:          body,
			},
		},
	}
	e := mockGetExternalInterface()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := e.StartUpdate(tt.args.taskID, tt.args.sessionUserName, tt.args.req); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StartUpdate() = %v, want %v", got, tt.want)
			}
		})
	}
}
