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

// Package evcommon ...
package evcommon

import (
	"context"
	"net/http"
	"testing"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-events/evmodel"
	"github.com/stretchr/testify/assert"
)

func TestMockIsAuthorized(t *testing.T) {
	type args struct {
		sessionToken  string
		privileges    []string
		oemPrivileges []string
	}
	tests := []struct {
		name string
		args args
		want response.RPC
	}{
		{
			name: "Test case 1 ",
			args: args{
				sessionToken:  "dummy",
				privileges:    []string{},
				oemPrivileges: []string{},
			},
			want: response.RPC{StatusCode: http.StatusUnauthorized},
		},
		{
			name: "Test case 2 ",
			args: args{
				sessionToken:  "validToken",
				privileges:    []string{},
				oemPrivileges: []string{},
			},
			want: response.RPC{StatusCode: http.StatusOK},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := MockIsAuthorized(tt.args.sessionToken, tt.args.privileges, tt.args.oemPrivileges); got.StatusCode != tt.want.StatusCode {
				t.Errorf("MockIsAuthorized() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMockGetSessionUserName(t *testing.T) {
	type args struct {
		sessionToken string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "valid token",
			args: args{sessionToken: "validToken"},
			want: "admin",
		},
		{
			name: "non-admin token",
			args: args{sessionToken: "token"},
			want: "non-admin",
		},
		{
			name: "Empty",
			args: args{sessionToken: ""},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := MockGetSessionUserName(tt.args.sessionToken)

			if got != tt.want {
				t.Errorf("MockGetSessionUserName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMockCreateTask(t *testing.T) {
	type args struct {
		sessionusername string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{

			name: "non-admin",
			args: args{sessionusername: "non-admin"},
			want: "",
		},
		{

			name: "non-admin",
			args: args{sessionusername: "/redfish/v1/tasks/123"},
			want: "/redfish/v1/tasks/123",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := MockCreateTask(context.TODO(), tt.args.sessionusername)
			if got != tt.want {
				t.Errorf("MockCreateTask() = %v, want %v", got, tt.want)
			}
		})
	}
	stubEMBConsume("")
	GetEncryptedKey([]byte{11})
}

func TestMockContactClient(t *testing.T) {
	type args struct {
		url         string
		method      string
		token       string
		odataID     string
		body        interface{}
		credentials map[string]string
	}
	tests := []struct {
		name string
		args args
		want *http.Response
	}{
		{
			name: "Positive Test ",
			args: args{url: "https://localhost:1234/ODIM/v1/Subscriptions",
				method: http.MethodDelete,
				body:   evmodel.Target{},
			},
			want: &http.Response{StatusCode: http.StatusNoContent},
		},
		{
			name: "Positive Test ",
			args: args{url: "https://localhost:1234/ODIM/v1/Subscriptions",
				method: http.MethodPost,
				body: &evmodel.Target{
					DeviceUUID: "d72dade0-c35a-984c-4859-1108132d72da",
				},
			},
			want: &http.Response{StatusCode: http.StatusBadRequest},
		},
		{
			name: "Positive Test ",
			args: args{url: "https://localhost:1234/ODIM/v1/Subscriptions",
				method: http.MethodPost,
				body: &evmodel.Target{
					DeviceUUID: "d72dade0-c35a-984c-4859-1108132d72daa",
				},
			},
			want: &http.Response{StatusCode: http.StatusCreated},
		},
		{
			name: "Positive Test ",
			args: args{url: "https://localhost:1234/ODIM/v1/Sessions",
				method: http.MethodPost,
				body: &evmodel.Target{
					DeviceUUID: "d72dade0-c35a-984c-4859-1108132d72daa",
				},
			},
			want: &http.Response{StatusCode: http.StatusCreated},
		},
		{
			name: "Positive Test ",
			args: args{url: "https://10.10.10.23:4321/ODIM/v1/Sessions",
				method: http.MethodPost,
				body: &evmodel.Target{
					DeviceUUID: "d72dade0-c35a-984c-4859-1108132d72daa",
				},
			},
			want: &http.Response{StatusCode: http.StatusCreated},
		},
		{
			name: "Positive Test ",
			args: args{url: "https://10.10.10.23:4321/ODIM/v1/Subscriptions",
				method: http.MethodPost,
				body: &evmodel.Target{
					DeviceUUID: "d72dade0-c35a-984c-4859-1108132d72daa",
				},
			},
			want: &http.Response{StatusCode: http.StatusUnauthorized},
		},
		{
			name: "Positive Test ",
			args: args{url: "https://odim.controller.com:1234/ODIM/v1/Subscriptions/123",
				method: http.MethodPost,
				body: &evmodel.Target{
					DeviceUUID: "d72dade0-c35a-984c-4859-1108132d72daa",
				},
			},
			want: &http.Response{StatusCode: http.StatusOK},
		},
		{
			name: "Positive Test ",
			args: args{url: "https://localhost:1234/ODIM/v1/Subscriptions/12345",
				method: http.MethodPost,
				body: &evmodel.Target{
					DeviceUUID: "d72dade0-c35a-984c-4859-1108132d72daa",
				},
			},
			want: &http.Response{StatusCode: http.StatusOK},
		},

		{
			name: "Positive Test ",
			args: args{url: "https://10.10.1.6:4321/ODIM/v1/Subscriptions",
				method: http.MethodPost,
				body: &evmodel.Target{
					DeviceUUID: "d72dade0-c35a-984c-4859-1108132d72daa",
				},
			},
			want: &http.Response{StatusCode: http.StatusCreated},
		},

		{
			name: "Positive Test ",
			args: args{url: "https://10.10.1.6:4321/ODIM/v1/Subscriptions/12345",
				method: http.MethodPost,
				body: &evmodel.Target{
					DeviceUUID: "d72dade0-c35a-984c-4859-1108132d72daa",
				},
			},
			want: &http.Response{StatusCode: http.StatusOK},
		},
		{
			name: "Invalid",
			args: args{url: "Invalid",
				method: http.MethodPost,
				body: &evmodel.Target{
					DeviceUUID: "d72dade0-c35a-984c-4859-1108132d72daa",
				},
			},
			want: &http.Response{StatusCode: 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := MockContactClient(context.TODO(), tt.args.url, tt.args.method, tt.args.token, tt.args.odataID, tt.args.body, tt.args.credentials)
			if got.StatusCode != tt.want.StatusCode {
				t.Errorf("MockContactClient() = %v, want %v", got, tt.want)
			}
		})
	}
	_, err := MockCreateChildTask(context.TODO(), "", destinationIP)
	assert.Nil(t, err)
	err = MockUpdateTask(context.TODO(), common.TaskData{})
	assert.Nil(t, err)
}

func TestMockGetPluginData(t *testing.T) {
	config.SetUpMockConfig(t)
	type args struct {
		pluginID string
	}
	tests := []struct {
		name string
		args args
		want *evmodel.Plugin
	}{
		{
			name: "GRF plugin",
			args: args{pluginID: "GRF"},
			want: &evmodel.Plugin{ID: "GRF"},
		},
		{
			name: "ILO plugin",
			args: args{pluginID: "ILO"},
			want: &evmodel.Plugin{ID: "ILO"},
		},
		{
			name: "CFM plugin",
			args: args{pluginID: "CFM"},
			want: &evmodel.Plugin{ID: "CFM"},
		},
		{
			name: "CFMPlugin plugin",
			args: args{pluginID: "CFMPlugin"},
			want: &evmodel.Plugin{ID: "CFMPlugin"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := MockGetPluginData(tt.args.pluginID)
			if got.ID != tt.want.ID {
				t.Errorf("MockGetPluginData() got = %v, want %v", got, tt.want)
			}

		})
	}
	_, err := MockGetPluginData("")
	assert.NotNil(t, err)
	MockGetSingleSystem("")

}

func TestMockGetTarget(t *testing.T) {
	config.SetUpMockConfig(t)
	type args struct {
		uuid string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Positive test case ",
			args:    args{uuid: "6d4a0a66-7efa-578e-83cf-44dc68d2874e"},
			wantErr: false,
		},
		{
			name:    "Positive test case ",
			args:    args{uuid: "11081de0-4859-984c-c35a-6c50732d72da"},
			wantErr: false,
		},
		{
			name:    "Positive test case ",
			args:    args{uuid: "d72dade0-c35a-984c-4859-1108132d72da"},
			wantErr: false,
		},
		{
			name:    "Positive test case ",
			args:    args{uuid: "110813e0-4859-984c-984c-d72da32d72da"},
			wantErr: false,
		},
		{
			name:    "Positive test case ",
			args:    args{uuid: "abab09db-e7a9-4352-8df0-5e41315a2a4c"},
			wantErr: false,
		},
		{
			name:    "Positive test case ",
			args:    args{uuid: "6d4a0a66-7efa-578e-83cf-44dc68d2874d"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := MockGetTarget(tt.args.uuid)
			if (err != nil) != tt.wantErr {
				t.Errorf("MockGetTarget() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
	_, err := MockGetTarget("default")
	assert.NotNil(t, err)
}

func TestMockGetFabricData(t *testing.T) {
	config.SetUpMockConfig(t)
	type args struct {
		fabricID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Positive Case ",
			args: args{
				fabricID: "123456",
			},
			wantErr: false,
		},
		{
			name: "Positive Case ",
			args: args{
				fabricID: "6d4a0a66-7efa-578e-83cf-44dc68d2874e",
			},
			wantErr: false,
		},
		{
			name: "Positive Case ",
			args: args{
				fabricID: "11081de0-4859-984c-c35a-6c50732d72da",
			},
			wantErr: false,
		},
		{
			name: "Positive Case ",
			args: args{
				fabricID: "48591de0-4859-1108-c35a-6c50110872da",
			},
			wantErr: false,
		},
		{
			name: "Nagative Case ",
			args: args{
				fabricID: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := MockGetFabricData(tt.args.fabricID)
			if (err != nil) != tt.wantErr {
				t.Errorf("MockGetFabricData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func TestMockGetAggregateDatacData(t *testing.T) {
	type args struct {
		aggregateID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Positive Test case ",
			args:    args{aggregateID: "6d4a0a66-7efa-578e-83cf-44dc68d2874e"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := MockGetAggregateDatacData(tt.args.aggregateID)
			if (err != nil) != tt.wantErr {
				t.Errorf("MockGetAggregateDatacData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestMockGetEvtSubscriptions(t *testing.T) {
	type args struct {
		searchKey string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Positive test",
			args: args{
				searchKey: "81de0110-c35a-4859-984c-072d6c5a32d7",
			},
			wantErr: false,
		},
		{
			name: "Positive test",
			args: args{
				searchKey: "/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874d.1",
			},
			wantErr: false,
		},
		{
			name: "Positive test",
			args: args{
				searchKey: "11081de0-4859-984c-c35a-6c50732d72da",
			},
			wantErr: false,
		},
		{
			name: "Positive test",
			args: args{
				searchKey: "71de0110-c35a-4859-984c-072d6c5a32d8",
			},
			wantErr: false,
		},
		{
			name: "Positive test",
			args: args{
				searchKey: "71de0110-c35a-4859-984c-072d6c5a32d9",
			},
			wantErr: false,
		},
		{
			name: "Positive test",
			args: args{
				searchKey: "5a321010-c35a-4859-984c-072d6c",
			},
			wantErr: false,
		},
		{
			name: "Positive test",
			args: args{
				searchKey: "71de0110-c35a-4859-984c-072d6c5a3211",
			},
			wantErr: false,
		},
		{
			name: "Positive test",
			args: args{
				searchKey: "81de0110-c35a-4859-984c-072d6c5a32d8",
			},
			wantErr: false,
		},
		{
			name: "Nagative",
			args: args{
				searchKey: "default",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := MockGetEvtSubscriptions(tt.args.searchKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("MockGetEvtSubscriptions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestMockGetDeviceSubscriptions(t *testing.T) {
	type args struct {
		hostIP string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Positive Test",
			args:    args{hostIP: "100.100.100.100"},
			wantErr: false,
		},
		{
			name:    "Positive Test",
			args:    args{hostIP: "100.100.100.101"},
			wantErr: false,
		},
		{
			name:    "Positive Test",
			args:    args{hostIP: "odim.ip.com"},
			wantErr: false,
		},
		{
			name:    "Positive Test",
			args:    args{hostIP: "odim.controller.com"},
			wantErr: false,
		},
		{
			name:    "Positive Test",
			args:    args{hostIP: "localhost"},
			wantErr: false,
		},
		{
			name:    "Positive Test",
			args:    args{hostIP: "odim.t.com"},
			wantErr: false,
		},
		{
			name:    "Positive Test",
			args:    args{hostIP: "*/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e.1"},
			wantErr: false,
		},
		{
			name:    "Nagative Test",
			args:    args{hostIP: "default"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := MockGetDeviceSubscriptions(tt.args.hostIP)
			if (err != nil) != tt.wantErr {
				t.Errorf("MockGetDeviceSubscriptions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
func TestFunc(t *testing.T) {
	err := MockSaveEventSubscription(evmodel.Subscription{})
	assert.Nil(t, err)

	err = MockUpdateEventSubscription(evmodel.Subscription{})
	assert.Nil(t, err)

	err = MockDeleteEvtSubscription("")
	assert.Nil(t, err)

	err = MockDeleteDeviceSubscription("")
	assert.Nil(t, err)
	err = MockUpdateDeviceSubscriptionLocation(evmodel.DeviceSubscription{})
	assert.Nil(t, err)
	_, err = MockGetAllKeysFromTable("")
	assert.Nil(t, err)
	_, err = MockGetAllFabrics()
	assert.Nil(t, err)
	_, err = MockGetAllMatchingDetails("", "*", common.InMemory)
	assert.Nil(t, err)
	_, err = MockGetAggregateHosts("")
	assert.Nil(t, err)
	err = MockSaveAggregateSubscription("", []string{})
	assert.Nil(t, err)

	err = MockSaveUndeliveredEvents("", []byte{})
	assert.Nil(t, err)

	err = MockSaveDeviceSubscription(evmodel.DeviceSubscription{})
	assert.Nil(t, err)

	_, err = MockGetUndeliveredEvents("")
	assert.Nil(t, err)
	_, err = MockGetUndeliveredEventsFlag("")
	assert.Nil(t, err)

	err = MockSetUndeliveredEventsFlag("")
	assert.Nil(t, err)

	err = MockDeleteUndeliveredEventsFlag("")
	assert.Nil(t, err)
	err = MockDeleteUndeliveredEvents("")
	assert.Nil(t, err)

}
