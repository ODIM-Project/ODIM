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
package system

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	aggregatorproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/aggregator"
	eventsproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/events"
	"github.com/ODIM-Project/ODIM/svc-aggregation/agmodel"
)

var successReq = map[string]interface{}{"Parameters": []Parameters{{Name: "/redfish/v1/Systems/ef83e569-7336-492a-aaee-31c02d9db831:1"}}}

func deleteComputeforTest(index int, key string) *errors.Error {
	if key == "/redfish/v1/Systems/del-comp-internal-err:1" {
		return errors.PackError(errors.UndefinedErrorType, "some internal error happed")
	}
	if key != "/redfish/v1/Systems/ef83e569-7336-492a-aaee-31c02d9db831:1" && key != "/redfish/v1/Systems/" &&
		key != "/redfish/v1/Systems/del-sys-internal-err:1" && key != "/redfish/v1/Systems/sys-not-found:1" {
		return errors.PackError(errors.DBKeyNotFound, "error while trying to get compute details: no data with the with key "+key+" found")
	}
	return nil
}

func deleteSystemforTest(key string) *errors.Error {
	if key == "del-sys-internal-err" {
		return errors.PackError(errors.UndefinedErrorType, "some internal error happed")
	}
	if key != "ef83e569-7336-492a-aaee-31c02d9db831" {
		return errors.PackError(errors.DBKeyNotFound, "error while trying to get compute details: no data with the with key "+key+" found")
	}
	return nil
}

func mockDeleteSubscription(uuid string) (*eventsproto.EventSubResponse, error) {
	if uuid == "/redfish/v1/Systems/ef83e569-7336-492a-aaee-31c02d9db832:1" {
		return nil, fmt.Errorf("error while trying to delete event subcription")
	} else if uuid == "/redfish/v1/Systems/unexpected-statuscode:1" {
		return &eventsproto.EventSubResponse{
			StatusCode: http.StatusCreated,
		}, nil
	}
	return &eventsproto.EventSubResponse{
		StatusCode: http.StatusNoContent,
	}, nil
}

func mockEventNotification(systemID, eventType, collectionType string) {
	return
}
func mockManagersData(id string, data map[string]interface{}) error {
	reqData, _ := json.Marshal(data)

	connPool, err := common.GetDBConnection(common.InMemory)
	if err != nil {
		return fmt.Errorf("error while trying to connecting to DB: %v", err.Error())
	}
	if err = connPool.Create("Managers", id, string(reqData)); err != nil {
		return fmt.Errorf("error while trying to create new %v resource: %v", "Managaers", err.Error())
	}
	return nil
}

func mockContactClientForDelete(url, method, token string, odataID string, body interface{}, credentials map[string]string) (*http.Response, error) {
	if url == "https://localhost:9092/ODIM/v1/Status" || (strings.Contains(url, "/ODIM/v1/Status") && credentials["UserName"] == "noStatusUser") {
		body := `{"MessageId": "Base.1.0.Success"}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil
	}
	return nil, fmt.Errorf("InvalidRequest")
}

func mockSystemOperationInfo() *errors.Error {
	systemOperation := agmodel.SystemOperation{
		Operation: "InventoryRediscovery ",
	}
	return systemOperation.AddSystemOperationInfo("/redfish/v1/Systems/ef83e569-7336-492a-aaee-31c02d9db831:1")
}
func TestDeleteCompute(t *testing.T) {
	d := &ExternalInterface{
		DeleteComputeSystem:     deleteComputeforTest,
		DeleteSystem:            deleteSystemforTest,
		DeleteEventSubscription: mockDeleteSubscription,
		EventNotification:       mockEventNotification,
		DecryptPassword:         stubDevicePassword,
	}
	successReq, _ := json.Marshal(successReq)
	subDeleteErrReq, _ := json.Marshal(map[string]interface{}{"Parameters": []Parameters{{Name: "/redfish/v1/Systems/ef83e569-7336-492a-aaee-31c02d9db832:1"}}})
	delCompSysInternalErrReq, _ := json.Marshal(map[string]interface{}{"Parameters": []Parameters{{Name: "/redfish/v1/Systems/del-comp-internal-err:1"}}})
	keyWithoutUUIDReq, _ := json.Marshal(map[string]interface{}{"Parameters": []Parameters{{Name: "/redfish/v1/Systems/"}}})
	delSysInternalErrReq, _ := json.Marshal(map[string]interface{}{"Parameters": []Parameters{{Name: "/redfish/v1/Systems/del-sys-internal-err:1"}}})
	sysNotFoundErrReq, _ := json.Marshal(map[string]interface{}{"Parameters": []Parameters{{Name: "/redfish/v1/Systems/sys-not-found:1"}}})
	invalidJSONErrReq, _ := json.Marshal(map[string]interface{}{"sample": "test"})
	type args struct {
		req *aggregatorproto.AggregatorRequest
	}
	tests := []struct {
		name string
		args args
		want int32
	}{
		{
			name: "compute deletion without request body",
			args: args{
				req: &aggregatorproto.AggregatorRequest{},
			},
			want: http.StatusBadRequest,
		},
		{
			name: "successful compute deletion",
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "SessionToken",
					RequestBody:  successReq,
				},
			},
			want: http.StatusOK,
		},
		{
			name: "delete subscription failure",
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "SessionToken",
					RequestBody:  subDeleteErrReq,
				},
			},
			want: http.StatusInternalServerError,
		},
		{
			name: "delete system internal error",
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "SessionToken",
					RequestBody:  delCompSysInternalErrReq,
				},
			},
			want: http.StatusInternalServerError,
		},
		{
			name: "key without UUID",
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "SessionToken",
					RequestBody:  keyWithoutUUIDReq,
				},
			},
			want: http.StatusInternalServerError,
		},
		{
			name: "System delete with UUID internal error",
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "SessionToken",
					RequestBody:  delSysInternalErrReq,
				},
			},
			want: http.StatusInternalServerError,
		},
		{
			name: "System not found with UUID",
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "SessionToken",
					RequestBody:  sysNotFoundErrReq,
				},
			},
			want: http.StatusNotFound,
		},
		{
			name: "malformed JSON",
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "SessionToken",
					RequestBody:  invalidJSONErrReq,
				},
			},
			want: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := d.DeleteCompute(tt.args.req)
			if got.StatusCode != tt.want {
				t.Errorf("DeleteCompute() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeleteComputeInvalidUUID(t *testing.T) {
	d := &ExternalInterface{
		DeleteComputeSystem:     deleteComputeforTest,
		DeleteSystem:            deleteSystemforTest,
		DeleteEventSubscription: mockDeleteSubscription,
		EventNotification:       mockEventNotification,
		DecryptPassword:         stubDevicePassword,
	}
	successReq["Parameters"] = []Parameters{{Name: "/redfish/v1/Systems/uuid:1"}}
	req, _ := json.Marshal(successReq)
	deleteReq := &aggregatorproto.AggregatorRequest{
		SessionToken: "SessionToken",
		RequestBody:  req,
	}
	resp := d.DeleteCompute(deleteReq)
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("Status code should be 404 but got %v", resp.StatusCode)
	}
}

func TestDeletePlugin(t *testing.T) {
	d := &ExternalInterface{
		EventNotification: mockEventNotification,
		ContactClient:     mockContactClientForDelete,
		DecryptPassword:   stubDevicePassword,
	}
	config.SetUpMockConfig(t)
	defer func() {
		err := common.TruncateDB(common.OnDisk)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		err = common.TruncateDB(common.InMemory)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()
	device1 := agmodel.Target{
		ManagerAddress: "100.0.0.1",
		Password:       []byte("imKp3Q6Cx989b6JSPHnRhritEcXWtaB3zqVBkSwhCenJYfgAYBf9FlAocE"),
		UserName:       "admin",
		DeviceUUID:     "24b243cf-f1e3-5318-92d9-2d6737d6b0b9",
		PluginID:       "ILO",
	}

	mockPluginData(t, "GRF")
	mockPluginData(t, "ILO")
	mockPluginData(t, "NoStatusPlugin")
	mockDeviceData("24b243cf-f1e3-5318-92d9-2d6737d6b0b9", device1)
	mockManagersData("/redfish/v1/Managers/1234877451-1234", map[string]interface{}{
		"Name": "GRF",
		"UUID": "1234877451-1234",
	})
	mockManagersData("/redfish/v1/Managers/1234877451-1233", map[string]interface{}{
		"Name": "ILO",
		"UUID": "1234877451-1233",
	})
	mockManagersData("/redfish/v1/Managers/invalid-plugin", map[string]interface{}{
		"Name": "invalid-Plugin",
		"UUID": "1234877451-1233",
	})
	mockManagersData("/redfish/v1/Managers/1234877451-1233", map[string]interface{}{
		"Name": "NoStatusPlugin",
		"UUID": "1234877451-1233",
	})

	successReq["Parameters"] = []Parameters{{Name: "/redfish/v1/Managers/1234877451-1234"}}
	successReq, _ := json.Marshal(successReq)
	managedDeviceReq, _ := json.Marshal(map[string]interface{}{"Parameters": []Parameters{{Name: "/redfish/v1/Managers/1234877451-1233"}}})
	invalidPluginReq, _ := json.Marshal(map[string]interface{}{"Parameters": []Parameters{{Name: "/redfish/v1/Managers/invalid-plugin"}}})
	invalidManagerReq, _ := json.Marshal(map[string]interface{}{"Parameters": []Parameters{{Name: "/redfish/v1/Managers/invalid-manager"}}})
	noStatusPluginReq, _ := json.Marshal(map[string]interface{}{"Parameters": []Parameters{{Name: "/redfish/v1/Managers/1234877451-1233"}}})

	type args struct {
		req *aggregatorproto.AggregatorRequest
	}
	tests := []struct {
		name string
		args args
		want int32
	}{
		{
			name: "successful plugin deletion",
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "SessionToken",
					RequestBody:  successReq,
				},
			},
			want: http.StatusOK,
		},
		{
			name: "deletion of plugin with mangaged devices",
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "SessionToken",
					RequestBody:  managedDeviceReq,
				},
			},
			want: http.StatusNotAcceptable,
		},
		{
			name: "invalid plugin",
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "SessionToken",
					RequestBody:  invalidPluginReq,
				},
			},
			want: http.StatusNotFound,
		},
		{
			name: "invalid msnsger",
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "SessionToken",
					RequestBody:  invalidManagerReq,
				},
			},
			want: http.StatusNotFound,
		},
		{
			name: "plugin status check failure",
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "SessionToken",
					RequestBody:  noStatusPluginReq,
				},
			},
			want: http.StatusNotAcceptable,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := d.DeleteCompute(tt.args.req)
			if got.StatusCode != tt.want {
				t.Errorf("DeleteCompute() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeleteComputeWithRediscovery(t *testing.T) {
	d := &ExternalInterface{
		DeleteComputeSystem:     deleteComputeforTest,
		DeleteSystem:            deleteSystemforTest,
		DeleteEventSubscription: mockDeleteSubscription,
		EventNotification:       mockEventNotification,
		DecryptPassword:         stubDevicePassword,
	}
	successReq = map[string]interface{}{"Parameters": []Parameters{{Name: "/redfish/v1/Systems/ef83e569-7336-492a-aaee-31c02d9db831:1"}}}

	successReq, _ := json.Marshal(successReq)
	type args struct {
		req *aggregatorproto.AggregatorRequest
	}
	config.SetUpMockConfig(t)
	defer func() {
		err := common.TruncateDB(common.OnDisk)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		err = common.TruncateDB(common.InMemory)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()
	err := mockSystemOperationInfo()
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	tests := []struct {
		name string
		args args
		want int32
	}{

		{
			name: "successful compute deletion",
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "SessionToken",
					RequestBody:  successReq,
				},
			},
			want: http.StatusNotAcceptable,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := d.DeleteCompute(tt.args.req)
			if got.StatusCode != tt.want {
				t.Errorf("DeleteCompute() = %v, want %v", got, tt.want)
			}
		})
	}
}
