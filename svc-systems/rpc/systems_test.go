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
package rpc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	systemsproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/systems"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-systems/smodel"
	"github.com/ODIM-Project/ODIM/svc-systems/systems"
)

func mockPluginClientData(t *testing.T) error {
	password := getEncryptedKey(t, []byte("$2a$10$OgSUYvuYdI/7dLL5KkYNp.RCXISefftdj.MjbBTr6vWyNwAvht6ci"))
	plugin := smodel.Plugin{
		IP:                "localhost",
		Port:              "9091",
		Username:          "admin",
		Password:          password,
		ID:                "GRF",
		PreferredAuthType: "BasicAuth",
	}
	connPool, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		return fmt.Errorf("error while trying to connecting to DB: %v", err)
	}
	// Clear previously created key, if exists
	connPool.Delete("Plugin", "GRF")
	if err = connPool.Create("Plugin", "GRF", plugin); err != nil {
		return fmt.Errorf("error while trying to create new %v resource: %v", "Plugin", err.Error())
	}
	return nil
}

func mockGetExternalInterface() *systems.ExternalInterface {
	return &systems.ExternalInterface{
		ContactClient:   contactPluginClient,
		DevicePassword:  stubDevicePassword,
		GetPluginStatus: mockPluginStatus,
	}
}
func getEncryptedKey(t *testing.T, key []byte) []byte {
	cryptedKey, err := common.EncryptWithPublicKey(key)
	if err != nil {
		t.Fatalf("error: failed to encrypt data: %v", err)
	}
	return cryptedKey
}

func mockDeviceData(uuid string, device smodel.Target) error {

	connPool, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		return fmt.Errorf("error while trying to connecting to DB: %v", err)
	}
	if err = connPool.Create("System", uuid, device); err != nil {
		return fmt.Errorf("error while trying to create new %v resource: %v", "System", err.Error())
	}
	return nil
}

func stubDevicePassword(password []byte) ([]byte, error) {
	return password, nil
}

func mockPluginStatus(plugin smodel.Plugin) bool {
	return true
}

func contactPluginClient(url, method, token string, odataID string, body interface{}, basicAuth map[string]string) (*http.Response, error) {
	if url == "https://localhost:9091/ODIM/v1/Systems/1/Storage/1/Volumes/1" {
		body := `{"MessageId": "` + response.Success + `"}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil
	}
	return nil, fmt.Errorf("InvalidRequest")
}

func mockSystemData(systemID string) error {
	reqData, _ := json.Marshal(&map[string]interface{}{
		"Manufacturer": "HPE",
		"Id":           "1",
	})

	connPool, err := common.GetDBConnection(common.InMemory)
	if err != nil {
		return fmt.Errorf("error while trying to connecting to DB: %v", err.Error())
	}
	if err = connPool.Create("ComputerSystem", systemID, string(reqData)); err != nil {
		return fmt.Errorf("error while trying to create new %v resource: %v", "System", err.Error())
	}
	return nil
}

func mockSystemResourceData(body []byte, table, key string) error {
	connPool, err := common.GetDBConnection(common.InMemory)
	if err != nil {
		return err
	}
	if err = connPool.Create(table, key, string(body)); err != nil {
		return err
	}
	return nil
}

func TestSystems_GetSystemResource(t *testing.T) {
	common.SetUpMockConfig()
	defer func() {
		err := common.TruncateDB(common.InMemory)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		err = common.TruncateDB(common.OnDisk)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()
	reqData := []byte(`\"@odata.id\":\"/redfsh/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1/SecureBoot\"`)
	err := mockResourceData(reqData, "SecureBoot", "/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1/SecureBoot")
	if err != nil {
		t.Fatalf("Error in creating mock resource data :%v", err)
	}
	sys := new(Systems)
	sys.IsAuthorizedRPC = mockIsAuthorized

	type args struct {
		ctx  context.Context
		req  *systemsproto.GetSystemsRequest
		resp *systemsproto.SystemsResponse
	}
	tests := []struct {
		name    string
		s       *Systems
		args    args
		wantErr bool
	}{
		{
			name: "Request with valid token",
			s:    sys,
			args: args{
				req: &systemsproto.GetSystemsRequest{
					RequestParam: "6d4a0a66-7efa-578e-83cf-44dc68d2874e:1",
					URL:          "/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1/SecureBoot",
					SessionToken: "validToken",
				},
				resp: &systemsproto.SystemsResponse{},
			}, wantErr: false,
		},
		{
			name: "Request with invalid token",
			s:    sys,
			args: args{
				req: &systemsproto.GetSystemsRequest{
					RequestParam: "6d4a0a66-7efa-578e-83cf-44dc68d2874e:1",
					URL:          "/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1/SecureBoot",
					SessionToken: "invalidToken",
				},
				resp: &systemsproto.SystemsResponse{},
			}, wantErr: false,
		},
		{
			name: "Request with invalid url",
			s:    sys,
			args: args{
				req: &systemsproto.GetSystemsRequest{
					RequestParam: "6d4a0a66-7efa-578e-83cf-44dc68d2874e:1",
					URL:          "/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1/SecureBoot1",
					SessionToken: "validToken",
				},
				resp: &systemsproto.SystemsResponse{},
			}, wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.s.GetSystemResource(tt.args.ctx, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("Systems.GetSystemResource() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSystems_GetAllSystems(t *testing.T) {
	common.SetUpMockConfig()
	defer func() {
		err := common.TruncateDB(common.InMemory)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		err = common.TruncateDB(common.OnDisk)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()
	reqData := []byte(`\"@odata.id\":\"/redfsh/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1\"`)
	err := mockResourceData(reqData, "ComputerSystem", "/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1")
	if err != nil {
		t.Fatalf("Error in creating mock resource data :%v", err)
	}
	sys := new(Systems)
	sys.IsAuthorizedRPC = mockIsAuthorized

	type args struct {
		ctx  context.Context
		req  *systemsproto.GetSystemsRequest
		resp *systemsproto.SystemsResponse
	}
	tests := []struct {
		name    string
		s       *Systems
		args    args
		wantErr bool
	}{
		{
			name: "Request with valid token",
			s:    sys,
			args: args{
				req: &systemsproto.GetSystemsRequest{
					URL:          "/redfish/v1/Systems",
					SessionToken: "validToken",
				},
				resp: &systemsproto.SystemsResponse{},
			}, wantErr: false,
		},
		{
			name: "Request with invalid token",
			s:    sys,
			args: args{
				req: &systemsproto.GetSystemsRequest{
					URL:          "/redfish/v1/Systems",
					SessionToken: "invalidToken",
				},
				resp: &systemsproto.SystemsResponse{},
			}, wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.s.GetSystemsCollection(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Systems.GetSystemsCollection() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSystems_GetSystems(t *testing.T) {
	common.SetUpMockConfig()
	defer func() {
		err := common.TruncateDB(common.InMemory)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		err = common.TruncateDB(common.OnDisk)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()
	reqData := []byte(`\"@odata.id\":\"/redfsh/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1\"`)
	err := mockResourceData(reqData, "ComputerSystem", "/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1")
	if err != nil {
		t.Fatalf("Error in creating mock resource data :%v", err)
	}
	sys := new(Systems)
	sys.IsAuthorizedRPC = mockIsAuthorized

	type args struct {
		ctx  context.Context
		req  *systemsproto.GetSystemsRequest
		resp *systemsproto.SystemsResponse
	}
	tests := []struct {
		name    string
		s       *Systems
		args    args
		wantErr bool
	}{
		{
			name: "Request with valid token",
			s:    sys,
			args: args{
				req: &systemsproto.GetSystemsRequest{
					URL:          "/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1",
					SessionToken: "validToken",
				},
				resp: &systemsproto.SystemsResponse{},
			}, wantErr: false,
		},
		{
			name: "Request with invalid token",
			s:    sys,
			args: args{
				req: &systemsproto.GetSystemsRequest{
					URL:          "/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1",
					SessionToken: "invalidToken",
				},
				resp: &systemsproto.SystemsResponse{},
			}, wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.s.GetSystems(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Systems.GetSystems() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSystems_ComputerSystemReset(t *testing.T) {
	common.SetUpMockConfig()
	sys := new(Systems)
	sys.IsAuthorizedRPC = mockIsAuthorized

	type args struct {
		ctx  context.Context
		req  *systemsproto.ComputerSystemResetRequest
		resp *systemsproto.SystemsResponse
	}
	tests := []struct {
		name    string
		s       *Systems
		args    args
		wantErr bool
	}{
		{
			name: "Request with valid token",
			s:    sys,
			args: args{
				req: &systemsproto.ComputerSystemResetRequest{
					RequestBody:  []byte(`{"ResetType": "ForceRestart"}`),
					SystemID:     "6d4a0a66-7efa-578e-83cf-44dc68d2874e:1",
					SessionToken: "validToken",
				},
				resp: &systemsproto.SystemsResponse{},
			}, wantErr: false,
		},
		{
			name: "Request with invalid token",
			s:    sys,
			args: args{
				req: &systemsproto.ComputerSystemResetRequest{
					RequestBody:  []byte(`{"ResetType": "ForceRestart"}`),
					SystemID:     "6d4a0a66-7efa-578e-83cf-44dc68d2874e:1",
					SessionToken: "invalidToken",
				},
				resp: &systemsproto.SystemsResponse{},
			}, wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.s.ComputerSystemReset(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Systems.ComputerSystemReset() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSystems_SetDefaultBootOrder(t *testing.T) {
	common.SetUpMockConfig()
	defer func() {
		err := common.TruncateDB(common.InMemory)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		err = common.TruncateDB(common.OnDisk)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()
	sys := new(Systems)
	sys.IsAuthorizedRPC = mockIsAuthorized

	type args struct {
		ctx  context.Context
		req  *systemsproto.DefaultBootOrderRequest
		resp *systemsproto.SystemsResponse
	}
	tests := []struct {
		name    string
		s       *Systems
		args    args
		wantErr bool
	}{
		{
			name: "Request with valid token",
			s:    sys,
			args: args{
				req: &systemsproto.DefaultBootOrderRequest{
					SystemID:     "6d4a0a66-7efa-578e-83cf-44dc68d2874e:1",
					SessionToken: "validToken",
				},
				resp: &systemsproto.SystemsResponse{},
			}, wantErr: false,
		},
		{
			name: "Request with invalid token",
			s:    sys,
			args: args{
				req: &systemsproto.DefaultBootOrderRequest{
					SystemID:     "6d4a0a66-7efa-578e-83cf-44dc68d2874e:1",
					SessionToken: "invalidToken",
				},
				resp: &systemsproto.SystemsResponse{},
			}, wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.s.SetDefaultBootOrder(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Systems.SetDefaultBootOrder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSystems_ChangeBiosSettings(t *testing.T) {
	common.SetUpMockConfig()
	defer func() {
		err := common.TruncateDB(common.InMemory)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		err = common.TruncateDB(common.OnDisk)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()
	sys := new(Systems)
	sys.IsAuthorizedRPC = mockIsAuthorized

	type args struct {
		ctx  context.Context
		req  *systemsproto.BiosSettingsRequest
		resp *systemsproto.SystemsResponse
	}
	tests := []struct {
		name    string
		s       *Systems
		args    args
		wantErr bool
	}{
		{
			name: "Request with valid token",
			s:    sys,
			args: args{
				req: &systemsproto.BiosSettingsRequest{
					SystemID:     "6d4a0a66-7efa-578e-83cf-44dc68d2874e:1",
					SessionToken: "validToken",
				},
				resp: &systemsproto.SystemsResponse{},
			}, wantErr: false,
		},
		{
			name: "Request with invalid token",
			s:    sys,
			args: args{
				req: &systemsproto.BiosSettingsRequest{
					SystemID:     "6d4a0a66-7efa-578e-83cf-44dc68d2874e:1",
					SessionToken: "invalidToken",
				},
				resp: &systemsproto.SystemsResponse{},
			}, wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.s.ChangeBiosSettings(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Systems.ChangeBiosSettings() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSystems_ChangeBootOrderSettings(t *testing.T) {
	common.SetUpMockConfig()
	defer func() {
		err := common.TruncateDB(common.InMemory)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		err = common.TruncateDB(common.OnDisk)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()
	sys := new(Systems)
	sys.IsAuthorizedRPC = mockIsAuthorized

	type args struct {
		ctx  context.Context
		req  *systemsproto.BootOrderSettingsRequest
		resp *systemsproto.SystemsResponse
	}
	tests := []struct {
		name    string
		s       *Systems
		args    args
		wantErr bool
	}{
		{
			name: "Request with valid token",
			s:    sys,
			args: args{
				req: &systemsproto.BootOrderSettingsRequest{
					SystemID:     "6d4a0a66-7efa-578e-83cf-44dc68d2874e:1",
					SessionToken: "validToken",
				},
				resp: &systemsproto.SystemsResponse{},
			}, wantErr: false,
		},
		{
			name: "Request with invalid token",
			s:    sys,
			args: args{
				req: &systemsproto.BootOrderSettingsRequest{
					SystemID:     "6d4a0a66-7efa-578e-83cf-44dc68d2874e:1",
					SessionToken: "invalidToken",
				},
				resp: &systemsproto.SystemsResponse{},
			}, wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.s.ChangeBootOrderSettings(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Systems.ChangeBootOrderSettings() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSystems_CreateVolume(t *testing.T) {
	common.SetUpMockConfig()
	defer func() {
		err := common.TruncateDB(common.InMemory)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		err = common.TruncateDB(common.OnDisk)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()
	sys := new(Systems)
	sys.IsAuthorizedRPC = mockIsAuthorized

	type args struct {
		ctx  context.Context
		req  *systemsproto.VolumeRequest
		resp *systemsproto.SystemsResponse
	}
	tests := []struct {
		name    string
		s       *Systems
		args    args
		wantErr bool
	}{
		{
			name: "Request with valid token",
			s:    sys,
			args: args{
				req: &systemsproto.VolumeRequest{
					SystemID:     "6d5a0a66-7efa-578e-83cf-44dc68d2874e:1",
					SessionToken: "validToken",
				},
				resp: &systemsproto.SystemsResponse{},
			}, wantErr: false,
		},
		{
			name: "Request with invalid token",
			s:    sys,
			args: args{
				req: &systemsproto.VolumeRequest{
					SystemID:     "6d5a0a66-7efa-578e-83cf-44dc68d2874e:1",
					SessionToken: "invalidToken",
				},
				resp: &systemsproto.SystemsResponse{},
			}, wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.s.CreateVolume(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Systems.CreateVolume() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSystems_DeleteVolume(t *testing.T) {
	config.SetUpMockConfig(t)
	defer func() {
		err := common.TruncateDB(common.InMemory)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		err = common.TruncateDB(common.OnDisk)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()
	device1 := smodel.Target{
		ManagerAddress: "10.24.0.12",
		Password:       []byte("imKp3Q6Cx989b6JSPHnRhritEcXWtaB3zqVBkSwhCenJYfgAYBf9FlAocE"),
		UserName:       "admin",
		DeviceUUID:     "6d5a0a66-7efa-578e-83cf-44dc68d2874e",
		PluginID:       "GRF",
	}
	mockPluginClientData(t)
	mockDeviceData("6d5a0a66-7efa-578e-83cf-44dc68d2874e", device1)
	mockSystemData("/redfish/v1/Systems/6d5a0a66-7efa-578e-83cf-44dc68d2874e:1")
	var reqData = `{"@odata.id":"/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1/Storage/1/Volumes/1"}`
	mockSystemResourceData([]byte(reqData), "Volumes", "/redfish/v1/Systems/6d5a0a66-7efa-578e-83cf-44dc68d2874e:1/Storage/1/Volumes/1")

	sys := new(Systems)
	sys.IsAuthorizedRPC = mockIsAuthorized
	sys.EI = mockGetExternalInterface()

	type args struct {
		ctx  context.Context
		req  *systemsproto.VolumeRequest
		resp *systemsproto.SystemsResponse
	}
	tests := []struct {
		name           string
		s              *Systems
		args           args
		wantStatusCode int32
	}{
		{
			name: "Request with valid token",
			s:    sys,
			args: args{
				req: &systemsproto.VolumeRequest{
					SystemID:        "6d5a0a66-7efa-578e-83cf-44dc68d2874e:1",
					SessionToken:    "validToken",
					StorageInstance: "1",
					VolumeID:        "1",
					RequestBody:     []byte(`{"@Redfish.OperationApplyTime": "OnReset"}`),
				},
				resp: &systemsproto.SystemsResponse{},
			},
			wantStatusCode: http.StatusNoContent,
		},
		{
			name: "Request with invalid token",
			s:    sys,
			args: args{
				req: &systemsproto.VolumeRequest{
					SystemID:        "6d5a0a66-7efa-578e-83cf-44dc68d2874e:1",
					SessionToken:    "invalidToken",
					StorageInstance: "1",
					VolumeID:        "1",
					RequestBody:     []byte(`{"@Redfish.OperationApplyTime": "OnReset"}`),
				},
				resp: &systemsproto.SystemsResponse{},
			},
			wantStatusCode: http.StatusUnauthorized,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, _ := tt.s.DeleteVolume(tt.args.ctx, tt.args.req)
			if resp.StatusCode != tt.wantStatusCode {
				t.Errorf("Systems.DeleteVolume() = %v, want %v", resp.StatusCode, tt.wantStatusCode)
			}
		})
	}
}
