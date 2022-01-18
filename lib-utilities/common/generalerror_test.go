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

package common

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/ODIM-Project/ODIM/lib-utilities/response"
)

// this flag is for checking whether the function call actually happened
var functionCallFlag bool

func mockUpdateTask(task TaskData) error {
	functionCallFlag = true
	return nil
}

func TestGeneralError(t *testing.T) {
	successRespArgs := response.Args{
		Code:    response.GeneralError,
		Message: "",
		ErrorArgs: []response.ErrArgs{
			response.ErrArgs{
				StatusMessage: response.PropertyMissing,
				ErrorMessage:  "error in processing the request",
				MessageArgs:   []interface{}{"token"},
			},
		},
	}
	type args struct {
		statusCode int32
		statusMsg  string
		errMsg     string
		msgArgs    []interface{}
		t          *TaskUpdateInfo
	}
	tests := []struct {
		name string
		args args
		want response.RPC
	}{
		{
			name: "General Error success case",
			args: args{
				statusCode: http.StatusBadRequest,
				statusMsg:  response.PropertyMissing,
				errMsg:     "error in processing the request",
				msgArgs:    []interface{}{"token"},
				t: &TaskUpdateInfo{
					TaskID:     "someID",
					TargetURI:  "someURI",
					UpdateTask: mockUpdateTask,
				},
			},
			want: response.RPC{
				StatusCode:    http.StatusBadRequest,
				StatusMessage: response.PropertyMissing,
				Body:          successRespArgs.CreateGenericErrorResponse(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GeneralError(tt.args.statusCode, tt.args.statusMsg, tt.args.errMsg, tt.args.msgArgs, tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GeneralError() = %v, want %v", got, tt.want)
			}
			if !functionCallFlag {
				t.Errorf("UpdateTask function was not called from GeneralError()")
			}
		})
	}
}
