//(C) Copyright [2022] Hewlett Packard Enterprise Development LP
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

// Package dpresponse ...
package dpresponse

import (
	"net/http/httptest"
	"reflect"
	"testing"

	iris "github.com/kataras/iris/v12"
)

func TestCreateErrorResponse(t *testing.T) {
	type args struct {
		errs string
	}
	tests := []struct {
		name string
		args args
		want ErrorResponse
	}{
		{
			name: "positive",
			args: args{
				errs: "fake error",
			},
			want: ErrorResponse{
				Error{
					Code:    "Base.1.13.0.GeneralError",
					Message: "See @Message.ExtendedInfo for more information.",
					MessageExtendedInfo: []MsgExtendedInfo{
						{
							MessageID: "Base.1.13.0.GeneralError",
							Message:   "fake error",
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateErrorResponse(tt.args.errs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateErrorResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetErrorResponse(t *testing.T) {
	// Create a new instance of iris.Context
	app := iris.New()
	ctx := app.ContextPool.Acquire(httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil))

	// Call the SetErrorResponse function
	SetErrorResponse(ctx, "test error", 500)

	// Check the response status code
	if ctx.ResponseWriter().StatusCode() != 500 {
		t.Errorf("Unexpected status code: %d", ctx.ResponseWriter().StatusCode())
	}
}
