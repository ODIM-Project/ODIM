// (C) Copyright [2020] Hewlett Packard Enterprise Development LP
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package systems

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	dmtf "github.com/ODIM-Project/ODIM/lib-dmtf/model"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	systemsproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/systems"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-systems/scommon"
	"github.com/ODIM-Project/ODIM/svc-systems/smodel"
	"github.com/ODIM-Project/ODIM/svc-systems/sresponse"
)

func mockSystemIndex(table, uuid string, indexData map[string]interface{}) error {
	scommon.SF.QueryKeys = []string{"filter"}
	scommon.SF.ConditionKeys = []string{"eq", "gt", "lt", "ge", "le", "ne"}
	scommon.SF.SearchKeys = []map[string]map[string]string{
		{
			"ProcessorSummary/Count": {
				"type": "float64",
			},
		},
		{
			"ProcessorSummary/Model": {
				"type": "string",
			},
		},
		{
			"Storage/Drives/Capacity": {
				"type": "[]float64",
			},
		},
		{
			"Storage/Drives/Type": {
				"type": "[]string",
			},
		},
	}

	connPool, err := common.GetDBConnection(common.InMemory)
	if err != nil {
		return fmt.Errorf("error while trying to connecting to DB: %v", err.Error())
	}

	if err := connPool.CreateIndex(indexData, "/redfish/v1/Systems/"+uuid); err != nil {
		return fmt.Errorf("error while creating  the index: %v", err.Error())
	}

	return nil

}

func mockSystemResourceData(body []byte, table, key string) error {
	connPool, err := common.GetDBConnection(common.InMemory)
	if err != nil {
		return fmt.Errorf("error while trying to connecting to DB: %v", err.Error())
	}
	if err = connPool.Create(table, key, string(body)); err != nil {
		return fmt.Errorf("error while trying to create new %v resource: %v", table, err.Error())
	}
	return nil
}

func getEncryptedKey(t *testing.T, key []byte) []byte {
	cryptedKey, err := common.EncryptWithPublicKey(key)
	if err != nil {
		t.Fatalf("error: failed to encrypt data: %v", err)
	}
	return cryptedKey
}

func mockTargetandPlugin(t *testing.T) error {
	connPool, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		return fmt.Errorf("error while trying to connecting to DB: %v", err.Error())
	}
	target := &smodel.Target{
		ManagerAddress: "10.24.0.14",
		Password:       []byte("Password"),
		UserName:       "admin",
		DeviceUUID:     "6d4a0a66-7efa-578e-83cf-44dc68d2874e",
		PluginID:       "GRF",
	}
	const table string = "System"
	//Save data into Database
	if err = connPool.Create(table, target.DeviceUUID, target); err != nil {
		return err
	}

	password := getEncryptedKey(t, []byte("Password"))
	plugin := &smodel.Plugin{
		IP:                "localhost",
		Port:              "9098",
		Password:          password,
		Username:          "admin",
		ID:                "GRF",
		PreferredAuthType: "basicauth",
		PluginType:        "GRF",
	}
	ptable := "Plugin"
	//Save data into Database
	if err = connPool.Create(ptable, "GRF", plugin); err != nil {
		return err
	}
	return nil
}

func mockGetDeviceInfo(req scommon.ResourceInfoRequest) (string, error) {
	if req.URL == "/redfish/v1/Systems/uuid:1/Storage" {
		return "", fmt.Errorf("error")
	}
	reqData := []byte(`\"@odata.id\":\"/redfish/v1/Systems/uuid:1/Storage\"`)
	return string(reqData), nil
}

func TestGetAllSystems(t *testing.T) {
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
	reqData, _ := json.Marshal(map[string]interface{}{"@odata.id": "/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1"})
	err := mockSystemResourceData(reqData, "ComputerSystem", "/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1")
	if err != nil {
		t.Fatalf("Error in creating mock resource data :%v", err)
	}
	var indexData = map[string]interface{}{
		"ProcessorSummary/Model":  "Intel",
		"ProcessorSummary/Count":  2,
		"Storage/Drives/Capacity": []float64{40},
		"Storage/Drives/Type":     []string{"HDD", "HDD"},
	}

	err = mockSystemIndex("/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1", "6d4a0a66-7efa-578e-83cf-44dc68d2874e:1", indexData)
	if err != nil {
		t.Fatalf("Error while creating mock index: %v", err)
	}
	header := map[string]string{
		"Allow":             `"GET"`,
		"Cache-Control":     "no-cache",
		"Connection":        "keep-alive",
		"Content-type":      "application/json; charset=utf-8",
		"Transfer-Encoding": "chunked",
		"OData-Version":     "4.0",
	}
	systemsCollection := sresponse.Collection{
		OdataContext: "/redfish/v1/$metadata#ComputerSystemCollection.ComputerSystemCollection",
		OdataID:      "/redfish/v1/Systems",
		OdataType:    "#ComputerSystemCollection.ComputerSystemCollection",
		Description:  "Computer Systems view",
		Name:         "Computer Systems",
	}
	systemsCollection.Members = []dmtf.Link{dmtf.Link{Oid: "/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1"}}
	systemsCollection.MembersCount = 1

	type args struct {
		req *systemsproto.GetSystemsRequest
	}
	tests := []struct {
		name    string
		args    args
		want    response.RPC
		wantErr bool
	}{
		{
			name: "successful get data",
			args: args{
				req: &systemsproto.GetSystemsRequest{
					URL: "/redfish/v1/Systems/",
				},
			},
			want: response.RPC{
				Header:        header,
				StatusCode:    http.StatusOK,
				StatusMessage: response.Success,
				Body:          systemsCollection,
			},
			wantErr: false,
		},
		{
			name: "successful get data with filter",
			args: args{
				req: &systemsproto.GetSystemsRequest{
					URL: "/redfish/v1/Systems?$filter=(ProcessorSummary/Model%20eq%20Intel%20or%20ProcessorSummary/Model%20eq%20amd)%20and%20(ProcessorSummary/Model%20eq%20Int*)",
				},
			},
			want: response.RPC{
				Header:        header,
				StatusCode:    http.StatusOK,
				StatusMessage: response.Success,
				Body:          systemsCollection,
			},
			wantErr: false,
		},
		{
			name: "successful get data with filter 2",
			args: args{
				req: &systemsproto.GetSystemsRequest{
					URL: "/redfish/v1/Systems?$filter=ProcessorSummary/Model%20ne%20AMD",
				},
			},
			want: response.RPC{
				Header:        header,
				StatusCode:    http.StatusOK,
				StatusMessage: response.Success,
				Body:          systemsCollection,
			},
			wantErr: false,
		},
		{
			name: "successful get data with filter 3",
			args: args{
				req: &systemsproto.GetSystemsRequest{
					URL: "/redfish/v1/Systems?$filter=ProcessorSummary/Count%20gt%201",
				},
			},
			want: response.RPC{
				Header:        header,
				StatusCode:    http.StatusOK,
				StatusMessage: response.Success,
				Body:          systemsCollection,
			},
			wantErr: false,
		},
		{
			name: "successful get data with filter 4",
			args: args{
				req: &systemsproto.GetSystemsRequest{
					URL: "/redfish/v1/Systems?$filter=ProcessorSummary/Count%20ge%202",
				},
			},
			want: response.RPC{
				Header:        header,
				StatusCode:    http.StatusOK,
				StatusMessage: response.Success,
				Body:          systemsCollection,
			},
			wantErr: false,
		},
		{
			name: "successful get data with filter 5",
			args: args{
				req: &systemsproto.GetSystemsRequest{
					URL: "/redfish/v1/Systems?$filter=ProcessorSummary/Count%20le%202",
				},
			},
			want: response.RPC{
				Header:        header,
				StatusCode:    http.StatusOK,
				StatusMessage: response.Success,
				Body:          systemsCollection,
			},
			wantErr: false,
		},
		{
			name: "successful get data with filter 6",
			args: args{
				req: &systemsproto.GetSystemsRequest{
					URL: "/redfish/v1/Systems?$filter=ProcessorSummary/Count%20lt%203",
				},
			},
			want: response.RPC{
				Header:        header,
				StatusCode:    http.StatusOK,
				StatusMessage: response.Success,
				Body:          systemsCollection,
			},
			wantErr: false,
		},
		{
			name: "successful get data with storage type ne SDD",
			args: args{
				req: &systemsproto.GetSystemsRequest{
					URL: "/redfish/v1/Systems?$filter=Storage/Drives/Type%20ne%20SDD",
				},
			},
			want: response.RPC{
				Header:        header,
				StatusCode:    http.StatusOK,
				StatusMessage: response.Success,
				Body:          systemsCollection,
			},
			wantErr: false,
		},
		{
			name: "successful get data with storage type HDD",
			args: args{
				req: &systemsproto.GetSystemsRequest{
					URL: "/redfish/v1/Systems?$filter=Storage/Drives/Type%20eq%20HDD",
				},
			},
			want: response.RPC{
				Header:        header,
				StatusCode:    http.StatusOK,
				StatusMessage: response.Success,
				Body:          systemsCollection,
			},
			wantErr: false,
		},
		{
			name: "successful get data with storage capcity greater than the value",
			args: args{
				req: &systemsproto.GetSystemsRequest{
					URL: "/redfish/v1/Systems?$filter=Storage/Drives/Capacity%20gt%2020",
				},
			},
			want: response.RPC{
				Header:        header,
				StatusCode:    http.StatusOK,
				StatusMessage: response.Success,
				Body:          systemsCollection,
			},
			wantErr: false,
		},
		{
			name: "successful get data with storage capcity greater than or equal to the value",
			args: args{
				req: &systemsproto.GetSystemsRequest{
					URL: "/redfish/v1/Systems?$filter=Storage/Drives/Capacity%20ge%2040",
				},
			},
			want: response.RPC{
				Header:        header,
				StatusCode:    http.StatusOK,
				StatusMessage: response.Success,
				Body:          systemsCollection,
			},
			wantErr: false,
		},
		{
			name: "successful get data with storage capcity less than the value",
			args: args{
				req: &systemsproto.GetSystemsRequest{
					URL: "/redfish/v1/Systems?$filter=Storage/Drives/Capacity%20lt%2050",
				},
			},
			want: response.RPC{
				Header:        header,
				StatusCode:    http.StatusOK,
				StatusMessage: response.Success,
				Body:          systemsCollection,
			},
			wantErr: false,
		},
		{
			name: "successful get data with storage capcity less than or equal to a value",
			args: args{
				req: &systemsproto.GetSystemsRequest{
					URL: "/redfish/v1/Systems?$filter=Storage/Drives/Capacity%20le%2040",
				},
			},
			want: response.RPC{
				Header:        header,
				StatusCode:    http.StatusOK,
				StatusMessage: response.Success,
				Body:          systemsCollection,
			},
			wantErr: false,
		},
		{
			name: "successful get data with not regular expression",
			args: args{
				req: &systemsproto.GetSystemsRequest{
					URL: "/redfish/v1/Systems?$filter=not%20Storage/Drives/Capacity%20eq%2030",
				},
			},
			want: response.RPC{
				Header:        header,
				StatusCode:    http.StatusOK,
				StatusMessage: response.Success,
				Body:          systemsCollection,
			},
			wantErr: false,
		},
		{
			name: "successful get data with storage capcity equal to the value",
			args: args{
				req: &systemsproto.GetSystemsRequest{
					URL: "/redfish/v1/Systems?$filter=Storage/Drives/Capacity%20eq%2040",
				},
			},
			want: response.RPC{
				Header:        header,
				StatusCode:    http.StatusOK,
				StatusMessage: response.Success,
				Body:          systemsCollection,
			},
			wantErr: false,
		},
		{
			name: "successful get data with or logical operation",
			args: args{
				req: &systemsproto.GetSystemsRequest{
					URL: "/redfish/v1/Systems?$filter=Storage/Drives/Capacity%20eq%2040%20or%20ProcessorSummary/Count%20lt%203",
				},
			},
			want: response.RPC{
				Header:        header,
				StatusCode:    http.StatusOK,
				StatusMessage: response.Success,
				Body:          systemsCollection,
			},
			wantErr: false,
		},
		{
			name: "successful get data with and logical operation",
			args: args{
				req: &systemsproto.GetSystemsRequest{
					URL: "/redfish/v1/Systems?$filter=Storage/Drives/Capacity%20eq%2040%20and%20ProcessorSummary/Count%20lt%203",
				},
			},
			want: response.RPC{
				Header:        header,
				StatusCode:    http.StatusOK,
				StatusMessage: response.Success,
				Body:          systemsCollection,
			},
			wantErr: false,
		},
		{
			name: "successful get data with  logical operation",
			args: args{
				req: &systemsproto.GetSystemsRequest{
					URL: "/redfish/v1/Systems?$filter=Storage/Drives/Capacity%20eq%2040%20and%20ProcessorSummary/Count%20lt%203",
				},
			},
			want: response.RPC{
				Header:        header,
				StatusCode:    http.StatusOK,
				StatusMessage: response.Success,
				Body:          systemsCollection,
			},
			wantErr: false,
		},
		{
			name: "successful get data with and logical operation",
			args: args{
				req: &systemsproto.GetSystemsRequest{
					URL: "/redfish/v1/Systems?$filter=(Storage/Drives/Capacity%20eq%2040%20and%20ProcessorSummary/Count%20lt%203)%20or%20(ProcessorSummary/Model%20ne%20AMD)",
				},
			},
			want: response.RPC{
				Header:        header,
				StatusCode:    http.StatusOK,
				StatusMessage: response.Success,
				Body:          systemsCollection,
			},
			wantErr: false,
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetSystemsCollection(tt.args.req)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSystemsInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSystems(t *testing.T) {
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
	reqData, _ := json.Marshal(map[string]interface{}{"@odata.id": "/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1"})
	err := mockSystemResourceData(reqData, "ComputerSystem", "/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1")
	if err != nil {
		t.Fatalf("Error in creating mock resource data :%v", err)
	}
	errArgs := response.Args{
		Code:    response.GeneralError,
		Message: "",
		ErrorArgs: []response.ErrArgs{
			response.ErrArgs{
				StatusMessage: response.ResourceNotFound,
				ErrorMessage:  "error: SystemUUID not found",
				MessageArgs:   []interface{}{"ComputerSystem", "6d4a0a66-7efa-578e-83cf-44dc68d2874e"},
			},
		},
	}
	errArgs1 := response.Args{
		Code:    response.GeneralError,
		Message: "",
		ErrorArgs: []response.ErrArgs{
			response.ErrArgs{
				StatusMessage: response.ResourceNotFound,
				ErrorMessage:  "error while trying to get system details: no data with the with key /redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e1:1 found",
				MessageArgs:   []interface{}{"ComputerSystem", "6d4a0a66-7efa-578e-83cf-44dc68d2874e1:1"},
			},
		},
	}
	pluginContact := PluginContact{
		ContactClient:  mockContactClient,
		DevicePassword: stubDevicePassword,
	}
	type args struct {
		req *systemsproto.GetSystemsRequest
	}
	tests := []struct {
		name    string
		p       *PluginContact
		args    args
		want    response.RPC
		wantErr bool
	}{
		{
			name: "successful get data",
			p:    &pluginContact,
			args: args{
				req: &systemsproto.GetSystemsRequest{
					RequestParam: "6d4a0a66-7efa-578e-83cf-44dc68d2874e:1",
					URL:          "/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1",
				},
			},
			want: response.RPC{
				Header: map[string]string{
					"Allow":             `"GET"`,
					"Cache-Control":     "no-cache",
					"Connection":        "keep-alive",
					"Content-type":      "application/json; charset=utf-8",
					"Transfer-Encoding": "chunked",
					"OData-Version":     "4.0",
				},
				StatusCode:    http.StatusOK,
				StatusMessage: response.Success,
				Body:          map[string]interface{}{"@odata.id": "/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1"},
			},
			wantErr: false,
		},
		{
			name: "invalid request param",
			p:    &pluginContact,
			args: args{
				req: &systemsproto.GetSystemsRequest{
					RequestParam: "6d4a0a66-7efa-578e-83cf-44dc68d2874e",
					URL:          "/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1",
				},
			},
			want: response.RPC{
				Header: map[string]string{
					"Content-type": "application/json; charset=utf-8",
				},
				StatusCode:    http.StatusNotFound,
				StatusMessage: response.ResourceNotFound,
				Body:          errArgs.CreateGenericErrorResponse(),
			},
			wantErr: true,
		},
		{
			name: "invalid url",
			p:    &pluginContact,
			args: args{
				req: &systemsproto.GetSystemsRequest{
					RequestParam: "6d4a0a66-7efa-578e-83cf-44dc68d2874e1:1",
					URL:          "/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e1:1",
				},
			},
			want: response.RPC{
				Header: map[string]string{
					"Content-type": "application/json; charset=utf-8",
				},
				StatusCode:    http.StatusNotFound,
				StatusMessage: response.ResourceNotFound,
				Body:          errArgs1.CreateGenericErrorResponse(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.p.GetSystems(tt.args.req)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSystems() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPluginContact_GetSystemResource(t *testing.T) {
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
	reqData, _ := json.Marshal(map[string]interface{}{"@odata.id": "/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1/SecureBoot"})
	err := mockSystemResourceData(reqData, "SecureBoot", "/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1/SecureBoot")
	if err != nil {
		t.Fatalf("Error in creating mock resource data :%v", err)
	}
	header := map[string]string{
		"Allow":             `"GET"`,
		"Cache-Control":     "no-cache",
		"Connection":        "keep-alive",
		"Content-type":      "application/json; charset=utf-8",
		"Transfer-Encoding": "chunked",
		"OData-Version":     "4.0",
	}
	pluginContact := PluginContact{
		ContactClient:  mockContactClient,
		DevicePassword: stubDevicePassword,
	}
	errArgs := response.Args{
		Code:    response.GeneralError,
		Message: "",
		ErrorArgs: []response.ErrArgs{
			response.ErrArgs{
				StatusMessage: response.ResourceNotFound,
				ErrorMessage:  "error: SystemUUID not found",
				MessageArgs:   []interface{}{"ComputerSystem", "6d4a0a66-7efa-578e-83cf-44dc68d2874e"},
			},
		},
	}
	errArgs1 := response.Args{
		Code:    response.GeneralError,
		Message: "",
		ErrorArgs: []response.ErrArgs{
			response.ErrArgs{
				StatusMessage: response.ResourceNotFound,
				ErrorMessage:  "error while trying to get compute details: no data with the with key 6d4a0a66-7efa-578e-83cf-44dc68d2874e found",
				MessageArgs:   []interface{}{"ComputerSystem", "/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1/SecureBoot1"},
			},
		},
	}
	type args struct {
		req *systemsproto.GetSystemsRequest
	}
	tests := []struct {
		name    string
		p       *PluginContact
		args    args
		want    response.RPC
		wantErr bool
	}{
		{
			name: "successful get data",
			p:    &pluginContact,
			args: args{
				req: &systemsproto.GetSystemsRequest{
					RequestParam: "6d4a0a66-7efa-578e-83cf-44dc68d2874e:1",
					URL:          "/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1/SecureBoot",
				},
			},
			want: response.RPC{
				Header:        header,
				StatusCode:    http.StatusOK,
				StatusMessage: response.Success,
				Body:          map[string]interface{}{"@odata.id": "/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1/SecureBoot"},
			},
			wantErr: false,
		},
		{
			name: "invalid request param",
			p:    &pluginContact,
			args: args{
				req: &systemsproto.GetSystemsRequest{
					RequestParam: "6d4a0a66-7efa-578e-83cf-44dc68d2874e",
					URL:          "/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1/SecureBoot",
				},
			},
			want: response.RPC{
				Header: map[string]string{
					"Content-type": "application/json; charset=utf-8",
				},
				StatusCode:    http.StatusNotFound,
				StatusMessage: response.ResourceNotFound,
				Body:          errArgs.CreateGenericErrorResponse(),
			},
			wantErr: true,
		},
		{
			name: "url if no present in data base",
			p:    &pluginContact,
			args: args{
				req: &systemsproto.GetSystemsRequest{
					RequestParam: "6d4a0a66-7efa-578e-83cf-44dc68d2874e:1",
					URL:          "/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1/SecureBoot1",
				},
			},
			want: response.RPC{
				Header: map[string]string{
					"Content-type": "application/json; charset=utf-8",
				},
				StatusCode:    http.StatusNotFound,
				StatusMessage: response.ResourceNotFound,
				Body:          errArgs1.CreateGenericErrorResponse(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.p.GetSystemResource(tt.args.req)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PluginContact.GetSystemResource() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestGetAllSystemsWithMultipleIndexData(t *testing.T) {
	t.Skip("skipping test")
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
	reqData, _ := json.Marshal(map[string]interface{}{"@odata.id": "/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1"})
	err := mockSystemResourceData(reqData, "ComputerSystem", "/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1")
	if err != nil {
		t.Fatalf("Error in creating mock resource data :%v", err)
	}
	err = mockSystemResourceData(reqData, "ComputerSystem", "/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874f:1")
	if err != nil {
		t.Fatalf("Error in creating mock resource data :%v", err)
	}
	err = mockSystemResourceData(reqData, "ComputerSystem", "/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874g:1")
	if err != nil {
		t.Fatalf("Error in creating mock resource data :%v", err)
	}

	var indexData1 = map[string]interface{}{
		"ProcessorSummary/Model":  "Intel",
		"ProcessorSummary/Count":  2,
		"Storage/Drives/Capacity": []float64{40},
		"Storage/Drives/Type":     []string{"HDD", "HDD"},
	}

	err = mockSystemIndex("/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1", "6d4a0a66-7efa-578e-83cf-44dc68d2874e:1", indexData1)
	if err != nil {
		t.Fatalf("Error while creating mock index: %v", err)
	}
	var indexData2 = map[string]interface{}{
		"ProcessorSummary/Model":  "amd",
		"ProcessorSummary/Count":  3,
		"Storage/Drives/Capacity": []float64{35},
		"Storage/Drives/Type":     []string{"HDD", "HDD"},
	}

	err = mockSystemIndex("/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874f:1", "6d4a0a66-7efa-578e-83cf-44dc68d2874f:1", indexData2)
	if err != nil {
		t.Fatalf("Error while creating mock index: %v", err)
	}

	var indexData3 = map[string]interface{}{
		"ProcessorSummary/Model":  "Intel",
		"ProcessorSummary/Count":  5,
		"Storage/Drives/Capacity": []float64{45},
		"Storage/Drives/Type":     []string{"HDD", "HDD"},
	}

	err = mockSystemIndex("/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874g:1", "6d4a0a66-7efa-578e-83cf-44dc68d2874g:1", indexData3)
	if err != nil {
		t.Fatalf("Error while creating mock index: %v", err)
	}

	header := map[string]string{
		"Allow":             `"GET"`,
		"Cache-Control":     "no-cache",
		"Connection":        "keep-alive",
		"Content-type":      "application/json; charset=utf-8",
		"Transfer-Encoding": "chunked",
		"OData-Version":     "4.0",
	}
	systemsCollection := sresponse.Collection{
		OdataContext: "/redfish/v1/$metadata#ComputerSystemCollection.ComputerSystemCollection",
		OdataID:      "/redfish/v1/Systems",
		OdataType:    "#ComputerSystemCollection.ComputerSystemCollection",
		Description:  "Computer Systems view",
		Name:         "Computer Systems",
	}
	systemsCollection.Members = []dmtf.Link{dmtf.Link{Oid: "/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1"},
		dmtf.Link{Oid: "/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874f:1"},
		dmtf.Link{Oid: "/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874g:1"}}
	systemsCollection.MembersCount = 3

	resp1 := systemsCollection
	resp1.Members = []dmtf.Link{dmtf.Link{Oid: "/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1"},
		dmtf.Link{Oid: "/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874g:1"}}
	resp1.MembersCount = len(resp1.Members)

	resp2 := systemsCollection
	resp2.Members = []dmtf.Link{dmtf.Link{Oid: "/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874g:1"},
		dmtf.Link{Oid: "/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1"}}
	resp2.MembersCount = len(resp2.Members)

	resp3 := systemsCollection
	resp3.Members = []dmtf.Link{dmtf.Link{Oid: "/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1"}}
	resp3.MembersCount = len(resp3.Members)

	resp4 := systemsCollection
	resp4.Members = []dmtf.Link{dmtf.Link{Oid: "/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1"},
		{Oid: "/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874f:1"}}
	resp4.MembersCount = len(resp4.Members)

	resp5 := systemsCollection
	resp5.Members = []dmtf.Link{dmtf.Link{Oid: "/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874f:1"},
		{Oid: "/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874g:1"},
		{Oid: "/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1"}}
	resp5.MembersCount = len(resp5.Members)
	type args struct {
		req *systemsproto.GetSystemsRequest
	}
	tests := []struct {
		name    string
		args    args
		want    response.RPC
		wantErr bool
	}{
		{
			name: "successful get data with filter1",
			args: args{
				req: &systemsproto.GetSystemsRequest{
					URL: "/redfish/v1/Systems?$filter=(ProcessorSummary/Model%20eq%20Intel%20or%20ProcessorSummary/Model%20eq%20amd)%20and%20(ProcessorSummary/Model%20eq%20Int*)",
				},
			},
			want: response.RPC{
				Header:        header,
				StatusCode:    http.StatusOK,
				StatusMessage: response.Success,
				Body:          resp2,
			},
			wantErr: false,
		},
		{
			name: "successful get data with filter 2",
			args: args{
				req: &systemsproto.GetSystemsRequest{
					URL: "/redfish/v1/Systems?$filter=ProcessorSummary/Model%20ne%20amd",
				},
			},
			want: response.RPC{
				Header:        header,
				StatusCode:    http.StatusOK,
				StatusMessage: response.Success,
				Body:          resp2,
			},
			wantErr: false,
		},
		{
			name: "successful get data with filter 3",
			args: args{
				req: &systemsproto.GetSystemsRequest{
					URL: "/redfish/v1/Systems?$filter=ProcessorSummary/Count%20gt%201",
				},
			},
			want: response.RPC{
				Header:        header,
				StatusCode:    http.StatusOK,
				StatusMessage: response.Success,
				Body:          systemsCollection,
			},
			wantErr: false,
		},
		{
			name: "successful get data with filter 4",
			args: args{
				req: &systemsproto.GetSystemsRequest{
					URL: "/redfish/v1/Systems?$filter=ProcessorSummary/Count%20ge%202",
				},
			},
			want: response.RPC{
				Header:        header,
				StatusCode:    http.StatusOK,
				StatusMessage: response.Success,
				Body:          systemsCollection,
			},
			wantErr: false,
		},
		{
			name: "successful get data with filter 5",
			args: args{
				req: &systemsproto.GetSystemsRequest{
					URL: "/redfish/v1/Systems?$filter=ProcessorSummary/Count%20le%202",
				},
			},
			want: response.RPC{
				Header:        header,
				StatusCode:    http.StatusOK,
				StatusMessage: response.Success,
				Body:          resp3,
			},
			wantErr: false,
		},
		{
			name: "successful get data with filter 6",
			args: args{
				req: &systemsproto.GetSystemsRequest{
					URL: "/redfish/v1/Systems?$filter=ProcessorSummary/Count%20lt%204",
				},
			},
			want: response.RPC{
				Header:        header,
				StatusCode:    http.StatusOK,
				StatusMessage: response.Success,
				Body:          resp4,
			},
			wantErr: false,
		},
		{
			name: "successful get data with filter 7",
			args: args{
				req: &systemsproto.GetSystemsRequest{
					URL: "/redfish/v1/Systems?$filter=ProcessorSummary/Count%20ne%200",
				},
			},
			want: response.RPC{
				Header:        header,
				StatusCode:    http.StatusOK,
				StatusMessage: response.Success,
				Body:          systemsCollection,
			},
			wantErr: false,
		},
		{
			name: "successful get data with filter 8",
			args: args{
				req: &systemsproto.GetSystemsRequest{
					URL: "/redfish/v1/Systems?$filter=ProcessorSummary/Count%20ne%203",
				},
			},
			want: response.RPC{
				Header:        header,
				StatusCode:    http.StatusOK,
				StatusMessage: response.Success,
				Body:          resp1,
			},
			wantErr: false,
		},
		{
			name: "successful get data with storage capcity greater than or equal to the value",
			args: args{
				req: &systemsproto.GetSystemsRequest{
					URL: "/redfish/v1/Systems?$filter=Storage/Drives/Capacity%20ge%2040",
				},
			},
			want: response.RPC{
				Header:        header,
				StatusCode:    http.StatusOK,
				StatusMessage: response.Success,
				Body:          resp1,
			},
			wantErr: false,
		},
		{
			name: "successful get data with not logical expression",
			args: args{
				req: &systemsproto.GetSystemsRequest{
					URL: "/redfish/v1/Systems?$filter=not%20Storage/Drives/Capacity%20eq%2030",
				},
			},
			want: response.RPC{
				Header:        header,
				StatusCode:    http.StatusOK,
				StatusMessage: response.Success,
				Body:          resp5,
			},
			wantErr: false,
		},
		{
			name: "successful get data with storage capcity equal to the value",
			args: args{
				req: &systemsproto.GetSystemsRequest{
					URL: "/redfish/v1/Systems?$filter=Storage/Drives/Capacity%20eq%2040",
				},
			},
			want: response.RPC{
				Header:        header,
				StatusCode:    http.StatusOK,
				StatusMessage: response.Success,
				Body:          resp3,
			},
			wantErr: false,
		},
		{
			name: "successful get data with or logical operation",
			args: args{
				req: &systemsproto.GetSystemsRequest{
					URL: "/redfish/v1/Systems?$filter=Storage/Drives/Capacity%20eq%2040%20or%20ProcessorSummary/Count%20lt%203",
				},
			},
			want: response.RPC{
				Header:        header,
				StatusCode:    http.StatusOK,
				StatusMessage: response.Success,
				Body:          resp3,
			},
			wantErr: false,
		},
		{
			name: "successful get data with and logical operation",
			args: args{
				req: &systemsproto.GetSystemsRequest{
					URL: "/redfish/v1/Systems?$filter=Storage/Drives/Capacity%20eq%2040%20and%20ProcessorSummary/Count%20lt%203",
				},
			},
			want: response.RPC{
				Header:        header,
				StatusCode:    http.StatusOK,
				StatusMessage: response.Success,
				Body:          resp3,
			},
			wantErr: false,
		},
		{
			name: "successful get data with  logical operation",
			args: args{
				req: &systemsproto.GetSystemsRequest{
					URL: "/redfish/v1/Systems?$filter=Storage/Drives/Capacity%20eq%2040%20and%20ProcessorSummary/Count%20lt%205",
				},
			},
			want: response.RPC{
				Header:        header,
				StatusCode:    http.StatusOK,
				StatusMessage: response.Success,
				Body:          resp3,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetSystemsCollection(tt.args.req)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllSystemsWithMultipleIndexData = %v, want %v", got, tt.want)
			}
		})
	}
}
