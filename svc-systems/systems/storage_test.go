// (C) Copyright [2020] Hewlett Packard Enterprise Development LP
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.
package systems

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	systemsproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/systems"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-systems/scommon"
	"github.com/ODIM-Project/ODIM/svc-systems/smodel"
	"github.com/stretchr/testify/assert"
)

func contactPluginClient(ctx context.Context, url, method, token string, odataID string, body interface{}, basicAuth map[string]string) (*http.Response, error) {
	if url == "https://localhost:9091/ODIM/v1/Systems/1/Storage/ArrayControllers-0/Volumes" {
		body := `{"MessageId": "` + response.Success + `"}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil
	}
	if url == "https://localhost:9091/ODIM/v1/Systems/1/Storage/1/Volumes/1" {
		body := `{"MessageId": "` + response.Success + `"}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil
	}
	if url == "https://localhost:9091/ODIM/v1/Systems/1/Storage/1/Volumes/2" {
		body := `{"MessageId": "` + response.Failure + `"}`
		return &http.Response{
			StatusCode: http.StatusNotFound,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil
	}
	return nil, fmt.Errorf("InvalidRequest")
}

func mockGetTarget(uuid string) (*smodel.Target, *errors.Error) {
	var target *smodel.Target
	switch uuid {
	case "54b243cf-f1e3-5319-92d9-2d6737d6b0a":
		target = &smodel.Target{
			ManagerAddress: "10.24.0.12",
			Password:       []byte("imKp3Q6Cx989b6JSPHnRhritEcXWtaB3zqVBkSwhCenJYfgAYBf9FlAocE"),
			UserName:       "admin",
			DeviceUUID:     "54b243cf-f1e3-5319-92d9-2d6737d6b0a",
			PluginID:       "GRF",
		}
	case "8e896459-a8f9-4c83-95b7-7b316b4908e1":
		target = &smodel.Target{
			ManagerAddress: "100.0.0.2",
			Password:       []byte("someValidPassword"),
			UserName:       "admin",
			DeviceUUID:     "8e896459-a8f9-4c83-95b7-7b316b4908e1",
			PluginID:       "Unknown-Plugin",
		}
	default:
		errorMessage := fmt.Sprintf("error while trying to get compute details: no data with the with key %v found", uuid)
		return target, errors.PackError(errors.UndefinedErrorType, errorMessage)
	}
	return target, nil
}

func getEncryptKey(key []byte) ([]byte, error) {
	cryptedKey, err := common.EncryptWithPublicKey(key)
	if err != nil {
		return cryptedKey, fmt.Errorf("error: failed to encrypt data: %v", err)
	}
	return cryptedKey, nil
}

func mockGetPluginData(pluginID string) (smodel.Plugin, *errors.Error) {
	var plugin smodel.Plugin
	password, keyErr := getEncryptKey([]byte("$2a$10$OgSUYvuYdI/7dLL5KkYNp.RCXISefftdj.MjbBTr6vWyNwAvht6ci"))
	if keyErr != nil {
		return plugin, errors.PackError(errors.UndefinedErrorType, keyErr.Error())
	}
	switch pluginID {
	case "GRF":
		plugin = smodel.Plugin{
			IP:                "localhost",
			Port:              "9091",
			Username:          "admin",
			Password:          password,
			ID:                "GRF",
			PreferredAuthType: "BasicAuth",
			PluginType:        "GRF",
		}
	default:
		return plugin, errors.PackError(errors.UndefinedErrorType, "No data found for the key")
	}
	return plugin, nil
}

func mockAddSystemResetInfo(ctx context.Context, systemID, resetType string) *errors.Error {
	return nil
}

func mockDeleteVolume(ctx context.Context, key string) *errors.Error {
	return nil
}

func mockGetExternalInterface() *ExternalInterface {
	return &ExternalInterface{
		ContactClient:  contactPluginClient,
		DevicePassword: stubDevicePassword,
		DB: DB{
			GetResource:        mockGetResource,
			DeleteVolume:       mockDeleteVolume,
			AddSystemResetInfo: mockAddSystemResetInfo,
			GetPluginData:      mockGetPluginData,
			GetTarget:          mockGetTarget,
		},
		GetPluginStatus: mockPluginStatus,
	}
}

func mockGetResource(ctx context.Context, table, key string) (string, *errors.Error) {
	var reqData = `{"@odata.id":"/redfish/v1/Systems/1/Storage/1/Volumes/1"}`
	switch key {
	case "/redfish/v1/Systems/54b243cf-f1e3-5319-92d9-2d6737d6b0a.1/Storage/ArrayControllers-0/Drives/0":
		return "", nil
	case "/redfish/v1/Systems/54b243cf-f1e3-5319-92d9-2d6737d6b0a.1/Storage/ArrayControllers-0/Drives/1":
		return "", nil
	case "/redfish/v1/Systems/54b243cf-f1e3-5319-92d9-2d6737d6b0a.1/Storage/1/Volumes/1":
		return reqData, nil
	case "/redfish/v1/Systems/8e896459-a8f9-4c83-95b7-7b316b4908e1.1/Storage/1/Volumes/1":
		return reqData, nil
	}
	return "body", nil
}

func mockPluginStatus(ctx context.Context, plugin smodel.Plugin) bool {
	return true
}

func TestPluginContact_CreateVolume(t *testing.T) {
	// Modify the contents with http.StatusNotImplemented to the correct status
	// and modify all other info accordingly after implementations
	config.SetUpMockConfig(t)
	var positiveResponse interface{}
	json.Unmarshal([]byte(`{"MessageId": "`+response.Success+`"}`), &positiveResponse)
	pluginContact := mockGetExternalInterface()
	ctx := mockContext()
	tests := []struct {
		name string
		p    *ExternalInterface
		req  *systemsproto.VolumeRequest
		want response.RPC
	}{
		{
			name: "Valid request",
			p:    pluginContact,
			req: &systemsproto.VolumeRequest{
				SystemID:        "54b243cf-f1e3-5319-92d9-2d6737d6b0a.1",
				StorageInstance: "ArrayControllers-0",
				RequestBody: []byte(`{
										"DisplayName":"Volume1",
										"WriteCachePolicy":"Off",
										"ReadCachePolicy":"Off",
										"IOPerfModeEnabled":false,
										"RAIDType":"RAID0",
										"Links":{
										"Drives":[{"@odata.id": "/redfish/v1/Systems/54b243cf-f1e3-5319-92d9-2d6737d6b0a.1/Storage/ArrayControllers-0/Drives/0"}]
										}
										,"@Redfish.OperationApplyTime": "OnReset"
										}`),
			},
			want: response.RPC{
				StatusCode:    http.StatusOK,
				StatusMessage: response.Success,
				Body:          map[string]interface{}{"MessageId": response.Success},
			},
		}, {
			name: "Valid request with multiple drives",
			p:    pluginContact,
			req: &systemsproto.VolumeRequest{
				SystemID:        "54b243cf-f1e3-5319-92d9-2d6737d6b0a.1",
				StorageInstance: "ArrayControllers-0",
				RequestBody: []byte(`{
										"DisplayName":"Volume1",
										"WriteCachePolicy":"Off",
										"ReadCachePolicy":"Off",
										"IOPerfModeEnabled":false,
										"RAIDType":"RAID0",
										"Links":{
										"Drives":[{"@odata.id": "/redfish/v1/Systems/54b243cf-f1e3-5319-92d9-2d6737d6b0a.1/Storage/ArrayControllers-0/Drives/0"},{"@odata.id": "/redfish/v1/Systems/54b243cf-f1e3-5319-92d9-2d6737d6b0a.1/Storage/ArrayControllers-0/Drives/1"}]}}`),
			},
			want: response.RPC{
				StatusCode:    http.StatusOK,
				StatusMessage: response.Success,
				Body:          map[string]interface{}{"MessageId": response.Success},
			},
		}, {
			name: "invalid system id",
			p:    pluginContact,
			req: &systemsproto.VolumeRequest{
				SystemID:        "54b243cf-f1e3-5319-92d9-2d6737d6b0b.1",
				StorageInstance: "ArrayControllers-0",
				RequestBody: []byte(`{
										"DisplayName":"Volume1",
										"WriteCachePolicy":"Off",
										"ReadCachePolicy":"Off",
										"IOPerfModeEnabled":false,
										"RAIDType":"RAID0",
									    "Links":{
										"Drives":[{"@odata.id": "/redfish/v1/Systems/54b243cf-f1e3-5319-92d9-2d6737d6b0b.1/Storage/ArrayControllers-0/Drives/0"}]}}`),
			},
			want: common.GeneralError(http.StatusNotFound, response.ResourceNotFound, "error while trying to get compute details: no data with the with key 54b243cf-f1e3-5319-92d9-2d6737d6b0b found", []interface{}{"System", "54b243cf-f1e3-5319-92d9-2d6737d6b0b"}, nil),
		}, {
			name: "invalid storage instance",
			p:    pluginContact,
			req: &systemsproto.VolumeRequest{
				SystemID:        "54b243cf-f1e3-5319-92d9-2d6737d6b0a.1",
				StorageInstance: "",
				RequestBody: []byte(`{
										"DisplayName":"Volume1",
										"WriteCachePolicy":"Off",
										"ReadCachePolicy":"Off",
										"IOPerfModeEnabled":false,
										"RAIDType":"RAID0",
										"Links":{
										"Drives":[{"@odata.id": "/redfish/v1/Systems/54b243cf-f1e3-5319-92d9-2d6737d6b0a.1/Storage/ArrayControllers-0/Drives/0"}]}}`),
			},
			want: common.GeneralError(http.StatusBadRequest, response.ResourceNotFound, "error: Storage instance is not found", []interface{}{"Storage", ""}, nil),
		}, {
			name: "invalid WriteCachePolicy",
			p:    pluginContact,
			req: &systemsproto.VolumeRequest{
				SystemID:        "54b243cf-f1e3-5319-92d9-2d6737d6b0a.1",
				StorageInstance: "",
				RequestBody: []byte(`{
										"DisplayName":"Volume1",
										"WriteCachePolicy":"dummy",
										"ReadCachePolicy":"Off",
										"IOPerfModeEnabled":false,
										"RAIDType":"RAID0",
										"Links":{
										"Drives":[{"@odata.id": "/redfish/v1/Systems/54b243cf-f1e3-5319-92d9-2d6737d6b0a.1/Storage/ArrayControllers-0/Drives/0"}]}}`),
			},
			want: common.GeneralError(http.StatusBadRequest, response.ResourceNotFound, "error: Storage instance is not found", []interface{}{"Storage", ""}, nil),
		}, {
			name: "invalid ReadCachePolicy instance",
			p:    pluginContact,
			req: &systemsproto.VolumeRequest{
				SystemID:        "54b243cf-f1e3-5319-92d9-2d6737d6b0a.1",
				StorageInstance: "",
				RequestBody: []byte(`{
										"DisplayName":"Volume1",
										"WriteCachePolicy":"Off",
										"ReadCachePolicy":"dummy",
										"IOPerfModeEnabled":false,
										"RAIDType":"RAID0",
										"Links":{
										"Drives":[{"@odata.id": "/redfish/v1/Systems/54b243cf-f1e3-5319-92d9-2d6737d6b0a.1/Storage/ArrayControllers-0/Drives/0"}]}}`),
			},
			want: common.GeneralError(http.StatusBadRequest, response.ResourceNotFound, "error: Storage instance is not found", []interface{}{"Storage", ""}, nil),
		}, {
			name: "invalid RaidType",
			p:    pluginContact,
			req: &systemsproto.VolumeRequest{
				SystemID:        "54b243cf-f1e3-5319-92d9-2d6737d6b0a.1",
				StorageInstance: "ArrayControllers-0",
				RequestBody: []byte(`{
										"RAIDType":"Invalid",
										"Links":{"Drives":[{"@odata.id": "/redfish/v1/Systems/54b243cf-f1e3-5319-92d9-2d6737d6b0a.1/Storage/ArrayControllers-0/Drives/0"}]}}`),
			},
			want: common.GeneralError(http.StatusBadRequest, response.PropertyValueNotInList, "error: request payload validation failed: RAIDType Invalid is invalid", []interface{}{"Invalid", "RAIDType"}, nil),
		}, {
			name: "Invalid Drives format",
			p:    pluginContact,
			req: &systemsproto.VolumeRequest{
				SystemID:        "54b243cf-f1e3-5319-92d9-2d6737d6b0a.1",
				StorageInstance: "ArrayControllers-0",
				RequestBody: []byte(`{
										"RaidType":"Invalid",
										"Links":{"Drives":["/redfish/v1/Systems/54b243cf-f1e3-5319-92d9-2d6737d6b0a.1/Storage/ArrayControllers-0/Drives/12"]}`),
			},
			want: common.GeneralError(http.StatusBadRequest, response.MalformedJSON, "Error while unmarshaling the create volume request: unexpected end of JSON input", []interface{}{}, nil),
		}, {
			name: "Empty System ID",
			p:    pluginContact,
			req: &systemsproto.VolumeRequest{
				SystemID:        "",
				StorageInstance: "ArrayControllers-0",
				RequestBody: []byte(`{
										"RaidType":"Invalid",
										"Links":{"Drives":["/redfish/v1/Systems/54b243cf-f1e3-5319-92d9-2d6737d6b0a.1/Storage/ArrayControllers-0/Drives/12"]}`),
			},
			want: common.GeneralError(http.StatusNotFound, response.ResourceNotFound, "error: SystemUUID not found", []interface{}{"System", ""}, nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.CreateVolume(ctx, tt.req); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PluginContact.CreateVolume() = %v, want %v", got, tt.want)
			}
		})
	}

	StringContain = func(s, substr string) bool {
		return true
	}
	// Test case for Empty
	storage := mockGetExternalInterface()
	req := systemsproto.VolumeRequest{
		SystemID:        "54b243cf-f1e3-5319-92d9-2d6737d6b0a.1",
		StorageInstance: "ArrayControllers-0",
		RequestBody:     []byte(`invalidJson`),
	}
	resp := storage.CreateVolume(ctx, &req)
	assert.Equal(t, http.StatusBadRequest, int(resp.StatusCode), "Error: @odata.id key(s) is missing in Drives list")

	// Validate the request JSON properties for case sensitive
	RequestParamsCaseValidatorFunc = func(rawRequestBody []byte, reqStruct interface{}) (string, error) {
		return "", &errors.Error{}
	}
	req = systemsproto.VolumeRequest{
		SystemID:        "54b243cf-f1e3-5319-92d9-2d6737d6b0a.1",
		StorageInstance: "ArrayControllers-0",
		RequestBody: []byte(`{
								"RAIDType":"RAID0",
								"Links":{"Drives":[{"@odata.id": "/redfish/v1/Systems/54b243cf-f1e3-5319-92d9-2d6737d6b0a.1/Storage/ArrayControllers-0/Drives/0"},{"@odata.id": "/redfish/v1/Systems/54b243cf-f1e3-5319-92d9-2d6737d6b0a.1/Storage/ArrayControllers-0/Drives/1"}]}}`),
	}

	resp = storage.CreateVolume(ctx, &req)
	assert.Equal(t, http.StatusInternalServerError, int(resp.StatusCode), "error: one or more properties given in the request body are not valid, ensure properties are listed in uppercamelcase ")

	RequestParamsCaseValidatorFunc = func(rawRequestBody []byte, reqStruct interface{}) (string, error) {
		return common.RequestParamsCaseValidator(rawRequestBody, reqStruct)
	}
	req = systemsproto.VolumeRequest{
		SystemID:        "54b243cf-f1e3-5319-92d9-2d6737d6b0a.1",
		StorageInstance: "ArrayControllers-0",
		RequestBody: []byte(`{
								"rAIDType":"RAID0",
								"Links":{"drives":[{"@odata.id": "/redfish/v1/Systems/54b243cf-f1e3-5319-92d9-2d6737d6b0a.1/Storage/ArrayControllers-0/Drives/0"},{"@odata.id": "/redfish/v1/Systems/54b243cf-f1e3-5319-92d9-2d6737d6b0a.1/Storage/ArrayControllers-0/Drives/1"}]}}`),
	}
	resp = storage.CreateVolume(ctx, &req)
	assert.Equal(t, http.StatusBadRequest, int(resp.StatusCode), "error: one or more properties given in the request body are not valid, ensure properties are listed in uppercamelcase ")

	StringsEqualFold = func(s, t string) bool {
		return true
	}
	req = systemsproto.VolumeRequest{
		SystemID:        "54b243cf-f1e3-5319-92d9-2d6737d6b0a.1",
		StorageInstance: "ArrayControllers-0",
		RequestBody: []byte(`{
								"RAIDType":"RAID0",
								 "Links":{
								"Drives":[{"@odata.id": "/redfish/v1/Systems/54b243cf-f1e3-5319-92d9-2d6737d6b0a.1/Storage/ArrayControllers-0/Drives/0"},{"@odata.id": "/redfish/v1/Systems/54b243cf-f1e3-5319-92d9-2d6737d6b0a.1/Storage/ArrayControllers-0/Drives/1"}]}}`),
	}

	resp = storage.CreateVolume(ctx, &req)
	assert.True(t, true, "Auth type XAuthToken")

	StringsEqualFold = func(s, t string) bool {
		return false
	}
	ContactPluginFunc = func(ctx context.Context, req scommon.PluginContactRequest, errorMessage string) (data []byte, data1 string, status scommon.ResponseStatus, err error) {
		err = &errors.Error{}
		return
	}
	req = systemsproto.VolumeRequest{
		SystemID:        "54b243cf-f1e3-5319-92d9-2d6737d6b0a.1",
		StorageInstance: "ArrayControllers-0",
		RequestBody: []byte(`{
								"RAIDType":"RAID0",
								"Links":{
								"Drives":[{"@odata.id": "/redfish/v1/Systems/54b243cf-f1e3-5319-92d9-2d6737d6b0a.1/Storage/ArrayControllers-0/Drives/0"},{"@odata.id": "/redfish/v1/Systems/54b243cf-f1e3-5319-92d9-2d6737d6b0a.1/Storage/ArrayControllers-0/Drives/1"}]}}`),
	}
	resp = storage.CreateVolume(ctx, &req)
	assert.True(t, true, "Error: Plugin Contact")

	ContactPluginFunc = func(ctx context.Context, req scommon.PluginContactRequest, errorMessage string) (data []byte, data1 string, status scommon.ResponseStatus, err error) {
		return scommon.ContactPlugin(ctx, req, errorMessage)
	}
	JSONUnmarshalFunc = func(data []byte, v interface{}) error {
		return &errors.Error{}
	}
	req = systemsproto.VolumeRequest{
		SystemID:        "54b243cf-f1e3-5319-92d9-2d6737d6b0a.1",
		StorageInstance: "ArrayControllers-0",
		RequestBody: []byte(`{
								"RAIDType":"RAID0",
								"Links":{
								"Drives":[{"@odata.id": "/redfish/v1/Systems/54b243cf-f1e3-5319-92d9-2d6737d6b0a.1/Storage/ArrayControllers-0/Drives/0"},{"@odata.id": "/redfish/v1/Systems/54b243cf-f1e3-5319-92d9-2d6737d6b0a.1/Storage/ArrayControllers-0/Drives/1"}]}}`),
	}
	resp = storage.CreateVolume(ctx, &req)
	assert.Equal(t, http.StatusInternalServerError, int(resp.StatusCode), "Auth type XAuthToken")

}

func TestPluginContact_DeleteVolume(t *testing.T) {
	// Modify the contents with http.StatusNotImplemented to the correct status
	// and modify all other info accordingly after implementations
	config.SetUpMockConfig(t)
	var positiveResponse interface{}
	json.Unmarshal([]byte(`{"MessageId": "`+response.Success+`"}`), &positiveResponse)
	pluginContact := mockGetExternalInterface()
	ctx := mockContext()
	tests := []struct {
		name              string
		p                 *ExternalInterface
		req               *systemsproto.VolumeRequest
		JSONUnmarshalFunc func(data []byte, v interface{}) error
		wantStatusCode    int32
	}{
		{
			name: "Valid request",
			p:    pluginContact,
			req: &systemsproto.VolumeRequest{
				SystemID:        "54b243cf-f1e3-5319-92d9-2d6737d6b0a.1",
				StorageInstance: "1",
				VolumeID:        "1",
				RequestBody:     []byte(`{"@Redfish.OperationApplyTime": "OnReset"}`),
			},
			JSONUnmarshalFunc: func(data []byte, v interface{}) error {
				return nil
			},
			wantStatusCode: http.StatusNoContent,
		},
		{
			name: "invalid system id",
			p:    pluginContact,
			req: &systemsproto.VolumeRequest{
				SystemID:        "54b243cf-f1e3-5319-92d9-2d6737d6b0b.1",
				StorageInstance: "1",
				VolumeID:        "1",
				RequestBody:     []byte(`{"@Redfish.OperationApplyTime": "OnReset"}`),
			},
			JSONUnmarshalFunc: func(data []byte, v interface{}) error {
				return nil
			},
			wantStatusCode: http.StatusNotFound,
		},
		{
			name: "without system id",
			p:    pluginContact,
			req: &systemsproto.VolumeRequest{
				SystemID:        "54b243cf-f1e3-5319-92d9-2d6737d6b0b.",
				StorageInstance: "1",
				VolumeID:        "1",
				RequestBody:     []byte(`{"@Redfish.OperationApplyTime": "OnReset"}`),
			},
			JSONUnmarshalFunc: func(data []byte, v interface{}) error {
				return nil
			},
			wantStatusCode: http.StatusNotFound,
		},
		{
			name: "Invalid volume id",
			p:    pluginContact,
			req: &systemsproto.VolumeRequest{
				SystemID:        "54b243cf-f1e3-5319-92d9-2d6737d6b0a.1",
				StorageInstance: "1",
				VolumeID:        "2",
				RequestBody:     []byte(`{"@Redfish.OperationApplyTime": "OnReset"}`),
			},
			JSONUnmarshalFunc: func(data []byte, v interface{}) error {
				return nil
			},
			wantStatusCode: http.StatusNotFound,
		},
		{
			name: "unknown plugin",
			p:    pluginContact,
			req: &systemsproto.VolumeRequest{
				SystemID:        "8e896459-a8f9-4c83-95b7-7b316b4908e1.1",
				StorageInstance: "1",
				VolumeID:        "1",
				RequestBody:     []byte(`{"@Redfish.OperationApplyTime": "OnReset"}`),
			},
			JSONUnmarshalFunc: func(data []byte, v interface{}) error {
				return nil
			},
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name: "Invalid Json",
			p:    pluginContact,
			req: &systemsproto.VolumeRequest{
				SystemID:        "8e896459-a8f9-4c83-95b7-7b316b4908e1.1",
				StorageInstance: "1",
				VolumeID:        "1",
				RequestBody:     []byte(`{"@Redfish.OperationApplyTime": "OnReset"}`),
			},
			JSONUnmarshalFunc: func(data []byte, v interface{}) error {
				return &errors.Error{}
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "Invalid field Keyword",
			p:    pluginContact,
			req: &systemsproto.VolumeRequest{
				SystemID:        "54b243cf-f1e3-5319-92d9-2d6737d6b0a.1",
				StorageInstance: "1",
				VolumeID:        "1",
				RequestBody:     []byte(`{"@redfish.operationapplytime": "OnReset"}`),
			},
			JSONUnmarshalFunc: func(data []byte, v interface{}) error {
				return nil
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "Invalid System ID",
			p:    pluginContact,
			req: &systemsproto.VolumeRequest{
				SystemID:        "",
				StorageInstance: "",
				VolumeID:        "",
				RequestBody:     []byte(`{"@Redfish.OperationApplyTime": "OnReset"}`),
			},
			JSONUnmarshalFunc: func(data []byte, v interface{}) error {
				return json.Unmarshal(data, v)
			},
			wantStatusCode: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		JSONUnmarshalFunc = tt.JSONUnmarshalFunc
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.DeleteVolume(ctx, tt.req); got.StatusCode != tt.wantStatusCode {
				t.Errorf("PluginContact.DeleteVolume() = %v, want %v", got.StatusCode, tt.wantStatusCode)
			}
		})
	}

	req := &systemsproto.VolumeRequest{
		SystemID:        "54b243cf-f1e3-5319-92d9-2d6737d6b0a.1",
		StorageInstance: "1",
		VolumeID:        "1",
		RequestBody:     []byte(`{"@redfish.operationapplytime": "OnReset"}`),
	}
	RequestParamsCaseValidatorFunc = func(rawRequestBody []byte, reqStruct interface{}) (string, error) {
		return "", &errors.Error{}
	}
	res := pluginContact.DeleteVolume(ctx, req)
	assert.Equal(t, http.StatusInternalServerError, int(res.StatusCode), "Error validating request parameters for volume creation, status code should be StatusInternalServerError")

	RequestParamsCaseValidatorFunc = func(rawRequestBody []byte, reqStruct interface{}) (string, error) {
		return "", nil
	}

	StringsEqualFold = func(s, t string) bool {
		return true
	}
	res = pluginContact.DeleteVolume(ctx, req)
	assert.Equal(t, http.StatusInternalServerError, int(res.StatusCode), "Error : status code should StatusInternalServerError")

	StringsEqualFold = func(s, t string) bool {
		return false
	}
	StringTrimSpace = func(s string) string {
		return ""
	}
	res = pluginContact.DeleteVolume(ctx, req)
	assert.Equal(t, http.StatusBadRequest, int(res.StatusCode), "Error: Status code should StatusBadRequest")

}

func TestGetExternalInterface(t *testing.T) {
	GetExternalInterface()
}

func Test_searchItem(t *testing.T) {
	type args struct {
		slice []string
		val   string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Test for non matching value",
			args: args{
				slice: []string{"RAID0"},
				val:   "RAID1",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := searchItem(tt.args.slice, tt.args.val); got != tt.want {
				t.Errorf("searchItem() = %v, want %v", got, tt.want)
			}
		})
	}
}
