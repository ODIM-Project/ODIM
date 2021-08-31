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

//Package rpc ...
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
	"github.com/ODIM-Project/ODIM/svc-fabrics/fabmodel"
	"github.com/ODIM-Project/ODIM/svc-fabrics/fabrics"

	fabricsproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/fabrics"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
)

func mockAuth(sessionToken string, privileges []string, oemPrivileges []string) response.RPC {
	if sessionToken == "valid" {
		return common.GeneralError(http.StatusOK, response.Success, "", nil, nil)
	} else if sessionToken == "invalid" {
		return common.GeneralError(http.StatusUnauthorized, response.NoValidSession, "error while trying to authenticate session", nil, nil)
	}
	return common.GeneralError(http.StatusForbidden, response.InsufficientPrivilege, "error while trying to authenticate session", nil, nil)
}

func getEncryptedKey(t *testing.T, key []byte) []byte {
	cryptedKey, err := common.EncryptWithPublicKey(key)
	if err != nil {
		t.Fatalf("error: failed to encrypt data: %v", err)
	}
	return cryptedKey
}

func mockPluginData(t *testing.T) error {
	password := getEncryptedKey(t, []byte("12345"))
	var plugin = fabmodel.Plugin{
		IP:                "localhost",
		Port:              "9091",
		Username:          "admin",
		Password:          password,
		ID:                "CFM",
		PreferredAuthType: "XAuthTOken",
	}
	connPool, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		return fmt.Errorf("error while trying to connecting to DB: %v", err.Error())
	}
	if err = connPool.Create("Plugin", "CFM", plugin); err != nil {
		return fmt.Errorf("error while trying to create new %v resource: %v", "Plugin", err.Error())
	}
	return nil
}

func mockContactClient(url, method, token string, odataID string, body interface{}, loginCredential map[string]string) (*http.Response, error) {
	if url == "https://localhost:9091/ODIM/v1/Sessions" {
		body := `{"Token": "12345"}`
		return &http.Response{
			StatusCode: http.StatusCreated,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
			Header: http.Header{
				"X-Auth-Token": []string{"12345"},
			},
		}, nil
	} else if url == "https://localhost:9092/ODIM/v1/Sessions" {
		body := `{"Token": ""}`
		return &http.Response{
			StatusCode: http.StatusUnauthorized,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil
	}
	if (url == "https://localhost:9091/ODIM/v1/Fabrics/fabid1" ||
		url == "https://localhost:9091/ODIM/v1/Fabrics/fabid1/Zones/Zone1" ||
		url == "https://localhost:9091/ODIM/v1/Fabrics") && token == "12345" && method == "POST" {
		body := `{"@odata.id": "/ODIM/v1/Fabrics/fabid1"}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
			Header: http.Header{
				"Location": []string{"12345"},
			},
		}, nil
	} else if (url == "https://localhost:9091/ODIM/v1/Fabrics/fabid1" ||
		url == "https://localhost:9091/ODIM/v1/Fabrics/fabid1/Zones/Zone1" ||
		url == "https://localhost:9091/ODIM/v1/Fabrics") && token == "12345" {
		body := `{"@odata.id": "/ODIM/v1/Fabrics/fabid1"}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil
	} else if (url == "https://localhost:9091/ODIM/v1/Fabrics/fabid2" || url == "https://localhost:9091/ODIM/v1/Fabrics/fabid2/Zones/Zone1") && token == "12345" {
		body := `{"error": "invalidrequest"}`
		return &http.Response{
			StatusCode: http.StatusNotFound,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil

	} else if url == "https://localhost:9093/ODIM/v1/Fabrics/fabid2" {
		body := `{"@odata.id": "/ODIM/v1/Fabrics/fabid1"}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil

	} else if url == "https://localhost:9095/ODIM/v1/Fabrics/fabid1/Zones/Zone1" {
		body := `{"@odata.id": "/ODIM/v1/Fabrics/fabid1}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil

	} else if token != "12345" {
		body := `{"error": "invalidsession"}`
		return &http.Response{
			StatusCode: http.StatusUnauthorized,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil
	}
	return nil, fmt.Errorf("InvalidRequest")
}
func TestFabrics_GetFabricResource(t *testing.T) {
	fabrics.Token.Tokens = make(map[string]string)
	config.SetUpMockConfig(t)
	defer func() {
		err := common.TruncateDB(common.OnDisk)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()
	err := mockPluginData(t)
	if err != nil {
		t.Fatalf("Error in creating mock DeviceData :%v", err)
	}
	var fabricsData = &Fabrics{
		IsAuthorizedRPC:  mockAuth,
		ContactClientRPC: mockContactClient,
	}
	type args struct {
		ctx context.Context
		req *fabricsproto.FabricRequest
	}
	tests := []struct {
		name    string
		f       *Fabrics
		args    args
		wantErr bool
	}{
		{
			name: "Postive Test Case",
			f:    fabricsData,
			args: args{
				req: &fabricsproto.FabricRequest{
					SessionToken: "valid",
					URL:          "/redfish/v1/Fabrics/fabid1",
					Method:       "GET",
				},
			}, wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.f.GetFabricResource(tt.args.ctx, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("Fabrics.GetFabricResource() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFabrics_UpdateFabricResource(t *testing.T) {
	fabrics.Token.Tokens = make(map[string]string)
	postData, _ := json.Marshal(map[string]interface{}{
		"@odata.id": "/redfish/v1/Fabrics",
	})
	config.SetUpMockConfig(t)
	defer func() {
		err := common.TruncateDB(common.OnDisk)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()
	err := mockPluginData(t)
	if err != nil {
		t.Fatalf("Error in creating mock DeviceData :%v", err)
	}
	var fabricsData = &Fabrics{
		IsAuthorizedRPC:  mockAuth,
		ContactClientRPC: mockContactClient,
	}
	type args struct {
		ctx context.Context
		req *fabricsproto.FabricRequest
	}
	tests := []struct {
		name    string
		f       *Fabrics
		args    args
		wantErr bool
	}{
		{
			name: "Postive Test Case",
			f:    fabricsData,
			args: args{
				req: &fabricsproto.FabricRequest{
					SessionToken: "valid",
					URL:          "/redfish/v1/Fabrics/fabid1/Zones/Zone1",
					RequestBody:  postData,
					Method:       "POST",
				},
			}, wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.f.UpdateFabricResource(tt.args.ctx, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("Fabrics.UpdateFabricResource() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFabrics_DeleteFabricResource(t *testing.T) {
	fabrics.Token.Tokens = make(map[string]string)
	config.SetUpMockConfig(t)
	defer func() {
		err := common.TruncateDB(common.OnDisk)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()
	err := mockPluginData(t)
	if err != nil {
		t.Fatalf("Error in creating mock DeviceData :%v", err)
	}
	var fabricsData = &Fabrics{
		IsAuthorizedRPC:  mockAuth,
		ContactClientRPC: mockContactClient,
	}
	type args struct {
		ctx context.Context
		req *fabricsproto.FabricRequest
	}
	tests := []struct {
		name    string
		f       *Fabrics
		args    args
		wantErr bool
	}{
		{
			name: "Postive Test Case",
			f:    fabricsData,
			args: args{
				req: &fabricsproto.FabricRequest{
					SessionToken: "valid",
					URL:          "/redfish/v1/Fabrics/fabid1",
					Method:       "DELETE",
				},
			}, wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.f.DeleteFabricResource(tt.args.ctx, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("Fabrics.DeleteFabricResource() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFabrics_AddFabric(t *testing.T) {
	config.SetUpMockConfig(t)
	defer func() {
		err := common.TruncateDB(common.OnDisk)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()
	err := mockPluginData(t)
	if err != nil {
		t.Fatalf("Error in creating mock DeviceData :%v", err)
	}
	var fabricsData = &Fabrics{
		IsAuthorizedRPC:  mockAuth,
		ContactClientRPC: mockContactClient,
	}

	type args struct {
		ctx context.Context
		req *fabricsproto.AddFabricRequest
	}
	tests := []struct {
		name    string
		f       *Fabrics
		args    args
		wantErr bool
	}{
		{
			name: "Postive Test Case",
			f:    fabricsData,
			args: args{
				req: &fabricsproto.AddFabricRequest{
					OriginResource: "/redfish/v1/Fabrics/a926dec5-61eb-499b-988a-d45b45847466",
					Address:        "localhost",
				},
			}, wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.f.AddFabric(tt.args.ctx, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("Fabrics.AddFabric() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
