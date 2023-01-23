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
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"

	"gotest.tools/assert"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	fabricsproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/fabrics"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-fabrics/fabmodel"
	"github.com/ODIM-Project/ODIM/svc-fabrics/fabresponse"
)

func getEncryptedKey(t *testing.T, key []byte) []byte {
	cryptedKey, err := common.EncryptWithPublicKey(key)
	if err != nil {
		t.Fatalf("error: failed to encrypt data: %v", err)
	}
	return cryptedKey
}

func mockPluginData(t *testing.T, pluginID, PreferredAuthType, port string) error {
	password := getEncryptedKey(t, []byte("$2a$10$OgSUYvuYdI/7dLL5KkYNp.RCXISefftdj.MjbBTr6vWyNwAvht6ci"))
	plugin := fabmodel.Plugin{
		IP:                "localhost",
		Port:              port,
		Username:          "admin",
		Password:          password,
		ID:                pluginID,
		PluginType:        "Fabric",
		PreferredAuthType: PreferredAuthType,
	}
	connPool, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		return fmt.Errorf("error while trying to connecting to DB: %v", err.Error())
	}
	if err = connPool.Create("Plugin", pluginID, plugin); err != nil {
		return fmt.Errorf("error while trying to create new %v resource: %v", "Plugin", err.Error())
	}
	return nil
}

func mockContactClient(ctx context.Context, url, method, token string, odataID string, body interface{}, loginCredential map[string]string) (*http.Response, error) {
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
	if (url == "https://localhost:9091/ODIM/v1/Fabrics/d72dade0-c35a-984c-4859-1108132d72da" ||
		url == "https://localhost:9091/ODIM/v1/Fabrics/d72dade0-c35a-984c-4859-1108132d72da/Zones/Zone1" ||
		url == "https://localhost:9091/ODIM/v1/Fabrics") && token == "12345" && method == "POST" {
		body := `{"@odata.id": "/ODIM/v1/Fabrics/d72dade0-c35a-984c-4859-1108132d72da"}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
			Header: http.Header{
				"Location": []string{"12345"},
			},
		}, nil
	} else if (url == "https://localhost:9091/ODIM/v1/Fabrics/d72dade0-c35a-984c-4859-1108132d72da" ||
		url == "https://localhost:9091/ODIM/v1/Fabrics/d72dade0-c35a-984c-4859-1108132d72da/Zones/Zone1" ||
		url == "https://localhost:9091/ODIM/v1/Fabrics") && token == "12345" {
		body := `{"@odata.id": "/ODIM/v1/Fabrics/d72dade0-c35a-984c-4859-1108132d72da"}`
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
		body := `{"@odata.id": "/ODIM/v1/Fabrics/d72dade0-c35a-984c-4859-1108132d72da"}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil

	} else if url == "https://localhost:9095/ODIM/v1/Fabrics/d72dade0-c35a-984c-4859-1108132d72da/Zones/Zone1" {
		body := `{"@odata.id": "/ODIM/v1/Fabrics/d72dade0-c35a-984c-4859-1108132d72da}`
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
	Token.Tokens = make(map[string]string)
	config.SetUpMockConfig(t)
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
					Method:       "GET",
				},
			},
			want: response.RPC{
				StatusCode:    http.StatusOK,
				StatusMessage: response.Success,
				Header: map[string]string{
					"Allow": `"GET", "PUT", "POST", "PATCH", "DELETE"`,
				},
				Body: map[string]interface{}{
					"@odata.id": "/redfish/v1/Fabrics/d72dade0-c35a-984c-4859-1108132d72da",
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
					Method:       "GET",
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
					Method:       "GET",
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
			if got := tt.f.GetFabricResource(tt.args.req); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Fabrics.GetFabricResource() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFabrics_GetFabricResourceWithNoValidSession(t *testing.T) {
	Token.Tokens = make(map[string]string)
	Token.Tokens["CFM"] = "234556"
	config.SetUpMockConfig(t)
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
					Method:       "GET",
				},
			},
			want: response.RPC{
				StatusCode:    http.StatusOK,
				StatusMessage: response.Success,
				Header: map[string]string{
					"Allow": `"GET", "PUT", "POST", "PATCH", "DELETE"`,
				},
				Body: map[string]interface{}{
					"@odata.id": "/redfish/v1/Fabrics/d72dade0-c35a-984c-4859-1108132d72da",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.GetFabricResource(tt.args.req); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Fabrics.GetFabricResource() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFabricsCollection_WithInvalidPlugin(t *testing.T) {
	Token.Tokens = make(map[string]string)
	Token.Tokens["CFM"] = "234556"
	config.SetUpMockConfig(t)
	defer func() {
		err := common.TruncateDB(common.OnDisk)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()
	err := mockPluginData(t, "CFM", "XAuthToken", "9094")
	if err != nil {
		t.Fatalf("Error in creating mock DeviceData :%v", err)
	}

	var f = &Fabrics{
		Auth:          mockAuth,
		ContactClient: mockContactClient,
	}

	req := &fabricsproto.FabricRequest{
		SessionToken: "valid",
		URL:          "/redfish/v1/Fabrics",
		Method:       http.MethodGet,
	}

	resp := f.GetFabricResource(req)
	assert.Equal(t, int(resp.StatusCode), http.StatusOK, "should be same")
}

func TestFabricsCollection_emptyCollection(t *testing.T) {
	config.SetUpMockConfig(t)
	defer func() {
		err := common.TruncateDB(common.OnDisk)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()

	var f = &Fabrics{
		Auth:          mockAuth,
		ContactClient: mockContactClient,
	}

	req := &fabricsproto.FabricRequest{
		SessionToken: "valid",
		URL:          "/redfish/v1/Fabrics",
		Method:       http.MethodGet,
	}

	resp := f.GetFabricResource(req)
	assert.Equal(t, int(resp.StatusCode), http.StatusOK, "StatusCode should be statusok")
	response := resp.Body.(fabresponse.FabricCollection)
	assert.Equal(t, response.MembersCount, 0, "MembersCount should be 0")
}

func TestFabricsCollection_Collection(t *testing.T) {
	config.SetUpMockConfig(t)
	defer func() {
		err := common.TruncateDB(common.OnDisk)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()

	err := mockFabricData("d72dade0-c35a-984c-4859-1108132d72da", "CFM")
	if err != nil {
		t.Fatalf("Error in creating mockFabricData :%v", err)
	}

	var f = &Fabrics{
		Auth:          mockAuth,
		ContactClient: mockContactClient,
	}

	req := &fabricsproto.FabricRequest{
		SessionToken: "valid",
		URL:          "/redfish/v1/Fabrics",
		Method:       http.MethodGet,
	}

	resp := f.GetFabricResource(req)
	assert.Equal(t, int(resp.StatusCode), http.StatusOK, "StatusCode should be statusok")
	response := resp.Body.(fabresponse.FabricCollection)
	assert.Equal(t, response.MembersCount, 1, "MembersCount should be 0")
	assert.Equal(t, response.Members[0].Oid, "/redfish/v1/Fabrics/d72dade0-c35a-984c-4859-1108132d72da", "odataid should be /redfish/v1/Fabrics/d72dade0-c35a-984c-4859-1108132d72da")
}
