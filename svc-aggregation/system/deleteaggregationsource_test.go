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
	"reflect"
	"strings"
	"testing"

	dmtf "github.com/ODIM-Project/ODIM/lib-dmtf/model"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	aggregatorproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/aggregator"
	eventsproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/events"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-aggregation/agmodel"
)

func deleteComputeforTest(index int, key string) *errors.Error {
	if key == "/redfish/v1/Systems/del-comp-internal-err.1" {
		return errors.PackError(errors.UndefinedErrorType, "some internal error happed")
	}
	if key != "/redfish/v1/Systems/ef83e569-7336-492a-aaee-31c02d9db831.1" && key != "/redfish/v1/Systems/" &&
		key != "/redfish/v1/Systems/del-sys-internal-err.1" && key != "/redfish/v1/Systems/sys-not-found.1" {
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
	if uuid == "/redfish/v1/Systems/ef83e569-7336-492a-aaee-31c02d9db832.1" {
		return nil, fmt.Errorf("error while trying to delete event subcription")
	} else if uuid == "/redfish/v1/Systems/unexpected-statuscode.1" {
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

func mockAggregateData(id string, data map[string]interface{}) error {
	connPool, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		return fmt.Errorf("error while trying to connecting to DB: %v", err.Error())
	}
	if err = connPool.Create("Aggregate", id, data); err != nil {
		return fmt.Errorf("error while trying to create new %v resource: %v", "Managaers", err.Error())
	}
	return nil
}

func getDataFromDB(table, key string, db common.DbType) (string, error) {
	connPool, err := common.GetDBConnection(db)
	if err != nil {
		return "", fmt.Errorf("error while trying to connecting to DB: %v", err.Error())
	}
	data, err := connPool.Read(table, key)
	if err != nil {
		return "", fmt.Errorf("error while trying to create new %v resource: %v", "Managaers", err.Error())
	}
	return data, nil
}

func mockContactClientForDelete(url, method, token string, odataID string, body interface{}, credentials map[string]string) (*http.Response, error) {
	if url == "https://localhost:9092/ODIM/v1/Status" || (strings.Contains(url, "/ODIM/v1/Status") && credentials["UserName"] == "noStatusUser") {
		body := `{"MessageId": "` + response.Success + `"}`
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
	return systemOperation.AddSystemOperationInfo("/redfish/v1/Systems/ef83e569-7336-492a-aaee-31c02d9db831.1")
}

func mockLogServicesCollectionData(id string, data map[string]interface{}) error {
	reqData, _ := json.Marshal(data)

	connPool, err := common.GetDBConnection(common.InMemory)
	if err != nil {
		return fmt.Errorf("error while trying to connecting to DB: %v", err.Error())
	}

	if err = connPool.Create("LogServicesCollection", id, string(reqData)); err != nil {
		return fmt.Errorf("error while trying to create new %v resource: %v", "LogServicesCollectionData", err.Error())
	}
	return nil
}
func TestDeleteAggregationSourceWithRediscovery(t *testing.T) {
	d := getMockExternalInterface()
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

	mockPluginData(t, "GRF_v1.0.0")
	mockManagersData("/redfish/v1/Managers/1234877451-1234", map[string]interface{}{
		"Name": "GRF_v1.0.0",
		"UUID": "1234877451-1234",
	})
	reqManagerGRF := agmodel.AggregationSource{
		HostName: "100.0.0.1:50000",
		UserName: "admin",
		Password: []byte("admin12345"),
		Links: map[string]interface{}{
			"ConnectionMethod": map[string]interface{}{
				"@odata.id": "/redfish/v1/AggregationService/ConnectionMethods/7ff3bd97-c41c-5de0-937d-85d390691b73",
			},
		},
	}
	reqSuccess := agmodel.AggregationSource{
		HostName: "100.0.0.1",
		UserName: "admin",
		Password: []byte("admin12345"),
		Links: map[string]interface{}{
			"ConnectionMethod": map[string]interface{}{
				"@odata.id": "/redfish/v1/AggregationService/ConnectionMethods/7ff3bd97-c41c-5de0-937d-85d390691b73",
			},
		},
	}
	device := agmodel.Target{
		ManagerAddress: "100.0.0.1",
		Password:       []byte("imKp3Q6Cx989b6JSPHnRhritEcXWtaB3zqVBkSwhCenJYfgAYBf9FlAocE"),
		UserName:       "admin",
		DeviceUUID:     "ef83e569-7336-492a-aaee-31c02d9db831",
		PluginID:       "GRF_v1.0.0",
	}
	mockDeviceData("ef83e569-7336-492a-aaee-31c02d9db831", device)
	mockSystemData("/redfish/v1/Systems/ef83e569-7336-492a-aaee-31c02d9db831.1")

	err := agmodel.AddAggregationSource(reqManagerGRF, "/redfish/v1/AggregationService/AggregationSources/123456")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	err = agmodel.AddAggregationSource(reqSuccess, "/redfish/v1/AggregationService/AggregationSources/ef83e569-7336-492a-aaee-31c02d9db831")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	err = mockSystemOperationInfo()
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
					URL:          "/redfish/v1/AggregationService/AggregationSources/ef83e569-7336-492a-aaee-31c02d9db831",
				},
			},
			want: http.StatusNotAcceptable,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := d.DeleteAggregationSource(tt.args.req)
			if got.StatusCode != tt.want {
				t.Errorf("DeleteAggregationSource() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExternalInterface_DeleteAggregationSourceManager(t *testing.T) {
	d := getMockExternalInterface()
	d.ContactClient = mockContactClientForDelete
	common.MuxLock.Lock()
	config.SetUpMockConfig(t)
	common.MuxLock.Unlock()
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
		PluginID:       "ILO_v1.0.0",
	}

	mockPluginData(t, "GRF_v1.0.0")
	mockPluginData(t, "ILO_v1.0.0")
	mockPluginData(t, "NoStatusPlugin_v1.0.0")
	mockDeviceData("24b243cf-f1e3-5318-92d9-2d6737d6b0b9", device1)
	mockManagersData("/redfish/v1/Managers/1234877451-1234", map[string]interface{}{
		"Name": "GRF_v1.0.0",
		"UUID": "1234877451-1234",
	})
	mockManagersData("/redfish/v1/Managers/1234877451-1233", map[string]interface{}{
		"Name": "ILO_v1.0.0",
		"UUID": "1234877451-1233",
		"LogService": map[string]interface{}{
			"@odata.id": "/redfish/v1/Managers/1234877451-1233/LogServices",
		},
	})
	mockLogServicesCollectionData("/redfish/v1/Managers/1234877451-1233/LogServices", map[string]interface{}{
		"ODataContext": "/redfish/v1/$metadata#LogServiceCollection.LogServiceCollection",
		"ODataID":      "/redfish/v1/Managers/1234877451-1233/LogServices",
		"ODataType":    "#LogServiceCollection.LogServiceCollection",
		"Description":  "Logs view",
		"Members":      map[string]interface{}{},
		"MembersCount": 0,
		"Name":         "Logs",
	})

	mockManagersData("/redfish/v1/Managers/1234877451-1235", map[string]interface{}{
		"Name": "NoStatusPlugin_v1.0.0",
		"UUID": "1234877451-1235",
	})
	reqManagerGRF := agmodel.AggregationSource{
		HostName: "100.0.0.1:50000",
		UserName: "admin",
		Password: []byte("admin12345"),
		Links: map[string]interface{}{
			"ConnectionMethod": map[string]interface{}{
				"@odata.id": "/redfish/v1/AggregationService/ConnectionMethods/7ff3bd97-c41c-5de0-937d-85d390691b73",
			},
		},
	}
	reqManagerILO := agmodel.AggregationSource{
		HostName: "100.0.0.1:50001",
		UserName: "admin",
		Password: []byte("admin12345"),
		Links: map[string]interface{}{
			"ConnectionMethod": map[string]interface{}{
				"@odata.id": "/redfish/v1/AggregationService/ConnectionMethods/c41cbd97-937d-1b73-c41c-1b7385d39069",
			},
		},
	}
	req1 := agmodel.AggregationSource{
		HostName: "100.0.0.1:50002",
		UserName: "admin",
		Password: []byte("admin12345"),
		Links: map[string]interface{}{
			"ConnectionMethod": map[string]interface{}{
				"@odata.id": "/redfish/v1/AggregationService/ConnectionMethods/e85bd91f-b257-4db8-b049-171099f3beec",
			},
		},
	}
	err := agmodel.AddAggregationSource(reqManagerILO, "/redfish/v1/AggregationService/AggregationSources/123455")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	err = agmodel.AddAggregationSource(reqManagerGRF, "/redfish/v1/AggregationService/AggregationSources/123456")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	err = agmodel.AddAggregationSource(req1, "/redfish/v1/AggregationService/AggregationSources/123457")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
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
					URL:          "/redfish/v1/AggregationService/AggregationSources/123456",
				},
			},
			want: http.StatusNoContent,
		},
		{
			name: "deletion of plugin with mangaged devices",
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "SessionToken",
					URL:          "/redfish/v1/AggregationService/AggregationSources/123455",
				},
			},
			want: http.StatusNotAcceptable,
		},
		{
			name: "deletion of plugin with invalid aggregation source id",
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "SessionToken",
					URL:          "/redfish/v1/AggregationService/AggregationSources/123434",
				},
			},
			want: http.StatusNotFound,
		},
		{
			name: "plugin status check failure",
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "SessionToken",
					URL:          "/redfish/v1/AggregationService/AggregationSources/123457",
				},
			},
			want: http.StatusNotAcceptable,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := d.DeleteAggregationSource(tt.args.req)
			if got.StatusCode != tt.want {
				t.Errorf("DeleteAggregationSource() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExternalInterface_DeleteBMC(t *testing.T) {
	d := getMockExternalInterface()
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
	mockPluginData(t, "GRF_v1.0.0")
	mockManagersData("/redfish/v1/Managers/1234877451-1234", map[string]interface{}{
		"Name": "GRF_v1.0.0",
		"UUID": "1234877451-1234",
	})
	reqManagerGRF := agmodel.AggregationSource{
		HostName: "100.0.0.1:50000",
		UserName: "admin",
		Password: []byte("admin12345"),
		Links: map[string]interface{}{
			"ConnectionMethod": map[string]interface{}{
				"@odata.id": "/redfish/v1/AggregationService/ConnectionMethods/7ff3bd97-c41c-5de0-937d-85d390691b73",
			},
		},
	}
	reqSuccess := agmodel.AggregationSource{
		HostName: "100.0.0.1",
		UserName: "admin",
		Password: []byte("admin12345"),
		Links: map[string]interface{}{
			"ConnectionMethod": map[string]interface{}{
				"@odata.id": "/redfish/v1/AggregationService/ConnectionMethods/7ff3bd97-c41c-5de0-937d-85d390691b73",
			},
		},
	}
	reqFailure := agmodel.AggregationSource{
		HostName: "100.0.0.2",
		UserName: "admin",
		Password: []byte("admin12345"),
		Links: map[string]interface{}{
			"ConnectionMethod": map[string]interface{}{
				"@odata.id": "/redfish/v1/AggregationService/ConnectionMethods/7ff3bd97-c41c-5de0-937d-85d390691b73",
			},
		},
	}
	device1 := agmodel.Target{
		ManagerAddress: "100.0.0.1",
		Password:       []byte("imKp3Q6Cx989b6JSPHnRhritEcXWtaB3zqVBkSwhCenJYfgAYBf9FlAocE"),
		UserName:       "admin",
		DeviceUUID:     "ef83e569-7336-492a-aaee-31c02d9db831",
		PluginID:       "GRF_v1.0.0",
	}
	device2 := agmodel.Target{
		ManagerAddress: "100.0.0.1",
		Password:       []byte("imKp3Q6Cx989b6JSPHnRhritEcXWtaB3zqVBkSwhCenJYfgAYBf9FlAocE"),
		UserName:       "admin",
		DeviceUUID:     "ef83e569-7336-492a-aaee-31c02d9db832",
		PluginID:       "GRF_v1.0.0",
	}

	mockDeviceData("ef83e569-7336-492a-aaee-31c02d9db831", device1)
	mockDeviceData("ef83e569-7336-492a-aaee-31c02d9db832", device2)
	mockSystemData("/redfish/v1/Systems/ef83e569-7336-492a-aaee-31c02d9db831.1")
	mockSystemData("/redfish/v1/Systems/ef83e569-7336-492a-aaee-31c02d9db832.1")

	err := agmodel.AddAggregationSource(reqManagerGRF, "/redfish/v1/AggregationService/AggregationSources/123456")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	err = agmodel.AddAggregationSource(reqSuccess, "/redfish/v1/AggregationService/AggregationSources/ef83e569-7336-492a-aaee-31c02d9db831")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	err = agmodel.AddAggregationSource(reqFailure, "/redfish/v1/AggregationService/AggregationSources/ef83e569-7336-492a-aaee-31c02d9db832")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	type args struct {
		req *aggregatorproto.AggregatorRequest
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
					URL:          "/redfish/v1/AggregationService/AggregationSources/ef83e569-7336-492a-aaee-31c02d9db831",
				},
			},
			want: http.StatusNoContent,
		},
		{
			name: "delete subscription failure",
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "SessionToken",
					URL:          "/redfish/v1/AggregationService/AggregationSources/ef83e569-7336-492a-aaee-31c02d9db832",
				},
			},
			want: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := d.DeleteAggregationSource(tt.args.req)
			if got.StatusCode != tt.want {
				t.Errorf("DeleteAggregationSource() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_deleteLinkDetails(t *testing.T) {
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
	var chassisLink []string
	chassisLink = append(chassisLink, "/redfish/v1/Managers/uuid.1")
	type args struct {
		managerData map[string]interface{}
		systemID    string
		chassisList []string
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		{
			name: "Test",
			args: args{
				managerData: map[string]interface{}{
					"Links": map[string]interface{}{
						"@odata.id": "/redfish/v1/Managers/uuid.1",
					},
				},
				systemID:    "/redfish/v1/Systems/ef83e569-7336-492a-aaee-31c02d9db831.1",
				chassisList: chassisLink,
			},
			want: map[string]interface{}{
				"Links": map[string]interface{}{
					"@odata.id": "/redfish/v1/Managers/uuid.1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := deleteLinkDetails(tt.args.managerData, tt.args.systemID, tt.args.chassisList); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("deleteLinkDetails() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_checkAndRemoveWildCardValueG(t *testing.T) {
	var values []string
	checkAndRemoveWildCardValue("", values)
}

func TestExternalInterface_deleteWildCardValues(t *testing.T) {
	p := getMockExternalInterface()
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
	type args struct {
		systemID string
	}
	tests := []struct {
		name string
		p    *ExternalInterface
		args args
	}{
		{
			name: "deleteWildCardValues",
			p:    p,
			args: args{
				systemID: "/redfish/v1/Systems/ef83e569-7336-492a-aaee-31c02d9db831.1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.deleteWildCardValues(tt.args.systemID)
		})
	}
}

func Test_checkAndRemoveWildCardValue(t *testing.T) {
	var values []string
	var want []string
	Values1 := []string{
		"45201b16-5305-49f0-846b-4597e982f6f8.1",
		"64992250-2a1a-41c6-82c6-b046140d615d.1",
	}
	type args struct {
		val    string
		values []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "nagative case",
			args: args{
				val:    "wildcard123",
				values: values,
			},
			want: want,
		},
		{
			name: "positive case",
			args: args{
				val:    "45201b16-5305-49f0-846b-4597e982f6f8.1",
				values: Values1,
			},
			want: []string{"64992250-2a1a-41c6-82c6-b046140d615d.1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkAndRemoveWildCardValue(tt.args.val, tt.args.values); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("checkAndRemoveWildCardValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExternalInterface_updateMemberCollection(t *testing.T) {
	p := getMockExternalInterface()
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
	type args struct {
		resName string
		odataID string
	}
	tests := []struct {
		name string
		e    *ExternalInterface
		args args
	}{
		{
			name: "test",
			e:    p,
			args: args{
				resName: "Collection",
				odataID: "/redfish/v1/Systems/ef83e569-7336-492a-aaee-31c02d9db831.1/storage/"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.e.updateMemberCollection(tt.args.resName, tt.args.odataID)
		})
	}
}

func Test_removeMemberFromCollection(t *testing.T) {
	type args struct {
		collectionOdataID string
		telemetryInfo     []*dmtf.Link
	}
	tests := []struct {
		name string
		args args
		want []*dmtf.Link
	}{
		{
			name: "test case",
			args: args{},
			want: []*dmtf.Link{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := removeMemberFromCollection(tt.args.collectionOdataID, tt.args.telemetryInfo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("removeMemberFromCollection() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_removeAggregationSourceFromAggregates(t *testing.T) {
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

	data := map[string]interface{}{
		"Elements": []map[string]interface{}{
			{"@odata.id": "/redfish/v1/Systems/12345678-1234-5678-9123-123456789012.1"},
			{"@odata.id": "/redfish/v1/Systems/12345678-1234-5678-9123-123456789123.1"},
			{"@odata.id": "/redfish/v1/Systems/12345678-1234-5678-9123-123456789124.1"},
		},
	}
	mockAggregateData("/redfish/v1/AggregationService/Aggregates/1234877451-1234", data)
	type args struct {
		systemURIs []string
	}
	wantFunc := func(data map[string]interface{}) string {
		jsonData, _ := json.Marshal(data)
		return string(jsonData)
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "no system URIs",
			args: args{},
			want: wantFunc(data),
		},
		{
			name: "invalid system URIs",
			args: args{systemURIs: []string{"/redfish/v1/Systems/12345678-1234-5678-9123-987654321012.1"}},
			want: wantFunc(data),
		},
		{
			name: "valid system URIs",
			args: args{
				systemURIs: []string{"/redfish/v1/Systems/12345678-1234-5678-9123-123456789012.1",
					"/redfish/v1/Systems/12345678-1234-5678-9123-123456789123.1",
				}},
			want: wantFunc(map[string]interface{}{
				"Elements": []map[string]interface{}{
					{"@odata.id": "/redfish/v1/Systems/12345678-1234-5678-9123-123456789124.1"},
				},
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			removeAggregationSourceFromAggregates(tt.args.systemURIs)
			elements, _ := getDataFromDB("Aggregate", "/redfish/v1/AggregationService/Aggregates/1234877451-1234", common.OnDisk)
			if tt.want != elements {
				t.Errorf("removeAggregationSourceFromAggregates() = %v, want %v", elements, tt.want)
			}
		})
	}
}
