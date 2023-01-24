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
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"testing"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
)

func Test_fabricFactory_updateFabricChassisResource(t *testing.T) {
	Token.Tokens = make(map[string]string)
	config.SetUpMockConfig(t)
	f := getFabricFactoryMock(nil)
	var r response.RPC
	successReq := json.RawMessage(`{"Name":"someNewName"}`)
	ctx := mockContext()
	initializeRPCResponse(&r, common.GeneralError(http.StatusOK, response.Success, "", nil, nil))

	errResp := response.RPC{
		StatusCode:    http.StatusNotFound,
		StatusMessage: response.ResourceNotFound,
	}
	json.Unmarshal([]byte(`{"MessageId":"`+response.GeneralError+`"}`), &errResp.Body)
	type args struct {
		url  string
		body *json.RawMessage
	}
	tests := []struct {
		name string
		f    *fabricFactory
		args args
		want response.RPC
	}{
		{
			name: "successful fabric resource update",
			f:    f,
			args: args{
				url:  "/redfish/v1/Chassis/valid_for_update",
				body: &successReq,
			},
			want: r,
		},
		{
			name: "fabric resource update with invalid url",
			f:    f,
			args: args{
				url:  "/redfish/v1/Chassis/invalid_for_update",
				body: &successReq,
			},
			want: errResp,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.updateFabricChassisResource(ctx, tt.args.url, tt.args.body); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fabricFactory.updateFabricChassisResource() = %v, want %v", got, tt.want)
			}
		})
	}

}

func Test_validateReqParamsCase(t *testing.T) {
	// config.SetUpMockConfig(t)

	JSONUnmarshalFunc = func(data []byte, v interface{}) error { return nil }
	ctx := mockContext()
	RequestParamsCaseValidatorFunc = func(rawRequestBody []byte, reqStruct interface{}) (string, error) {
		return "", errors.New("")
	}
	validateReqParamsCase(ctx, &json.RawMessage{})
	RequestParamsCaseValidatorFunc = func(rawRequestBody []byte, reqStruct interface{}) (string, error) {
		return "dummy", nil
	}
	validateReqParamsCase(ctx, &json.RawMessage{})
}
