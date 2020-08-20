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
// Package agmodel ...
package agmodel

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	dmtfmodel "github.com/ODIM-Project/ODIM/lib-dmtf/model"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
)

func getEncryptedKey(t *testing.T, key []byte) []byte {
	cryptedKey, err := common.EncryptWithPublicKey(key)
	if err != nil {
		t.Fatalf("error: failed to encrypt data: %v", err)
	}
	return cryptedKey
}

func mockData(t *testing.T, dbType common.DbType, table, id string, data interface{}) {
	connPool, err := common.GetDBConnection(dbType)
	if err != nil {
		t.Fatalf("error: mockData() failed to DB connection: %v", err)
	}
	if err = connPool.Create(table, id, data); err != nil {
		t.Fatalf("error: mockData() failed to create entry %s-%s: %v", table, id, err)
	}
}

func mockIndex(dbType common.DbType, index, key string) {
	connPool, _ := common.GetDBConnection(dbType)
	form := map[string]interface{}{index: "value", index: "value2"}
	connPool.CreateIndex(form, "/redfish/v1/systems/ef83e569-7336-492a-aaee-31c02d9db831:1")
}

func TestGetResource(t *testing.T) {
	config.SetUpMockConfig(t)
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()
	mockData(t, common.InMemory, "successTable", "successID", "successData")
	type args struct {
		Table string
		key   string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 *errors.Error
	}{
		{
			name: "success case",
			args: args{
				Table: "successTable",
				key:   "successID",
			},
			want:  "successData",
			want1: nil,
		},
		{
			name: "not found case",
			args: args{
				Table: "noTable",
				key:   "successID",
			},
			want:  "",
			want1: errors.PackError(errors.DBKeyNotFound, "error while trying to get resource details: no data with the with key successID found"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := GetResource(tt.args.Table, tt.args.key)
			if got != tt.want {
				t.Errorf("GetResource() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("GetResource() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestSaveSystem_Create(t *testing.T) {
	config.SetUpMockConfig(t)
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()
	type args struct {
		systemID string
	}
	tests := []struct {
		name   string
		system *SaveSystem
		args   args
		want   *errors.Error
	}{
		{
			name:   "positive case",
			system: &SaveSystem{ManagerAddress: "123"},
			args:   args{systemID: "xyz"},
			want:   nil,
		},
		{
			name:   "already exist",
			system: &SaveSystem{ManagerAddress: "123"},
			args:   args{systemID: "xyz"},
			want:   errors.PackError(errors.DBKeyAlreadyExist, "error: data with key xyz already exists"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.system.Create(tt.args.systemID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SaveSystem.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetPluginData(t *testing.T) {
	config.SetUpMockConfig(t)

	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()

	validPassword := []byte("password")
	invalidPassword := []byte("invalid")
	validPasswordEnc := getEncryptedKey(t, []byte("password"))

	pluginData := Plugin{
		IP:                "localhost",
		Port:              "45001",
		Username:          "admin",
		Password:          validPasswordEnc,
		ID:                "GRF",
		PluginType:        "RF-GENERIC",
		PreferredAuthType: "BasicAuth",
	}
	mockData(t, common.OnDisk, "Plugin", "validPlugin", pluginData)
	pluginData.Password = invalidPassword
	mockData(t, common.OnDisk, "Plugin", "invalidPassword", pluginData)
	mockData(t, common.OnDisk, "Plugin", "invalidPluginData", "pluginData")

	type args struct {
		pluginID string
	}
	tests := []struct {
		name    string
		args    args
		exec    func(*Plugin)
		want    Plugin
		wantErr bool
	}{
		{
			name: "Positive Case",
			args: args{pluginID: "validPlugin"},
			exec: func(want *Plugin) {
				want.Password = validPassword
			},
			want:    pluginData,
			wantErr: false,
		},
		{
			name:    "Negative Case - Non-existent plugin",
			args:    args{pluginID: "notFound"},
			exec:    nil,
			want:    Plugin{},
			wantErr: true,
		},
		{
			name:    "Negative Case - Invalid plugin data",
			args:    args{pluginID: "invalidPluginData"},
			exec:    nil,
			want:    Plugin{},
			wantErr: true,
		},
		{
			name:    "Negative Case - Plugin with invalid password",
			args:    args{pluginID: "invalidPassword"},
			exec:    nil,
			want:    Plugin{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		if tt.exec != nil {
			tt.exec(&tt.want)
		}
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetPluginData(tt.args.pluginID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPluginData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetPluginData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetComputeSystem(t *testing.T) {
	config.SetUpMockConfig(t)
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()
	mockData(t, common.InMemory, "ComputerSystem", "someID", dmtfmodel.ComputerSystem{ID: "someID"})
	mockData(t, common.InMemory, "ComputerSystem", "falseData", "some data")
	type args struct {
		deviceUUID string
	}
	tests := []struct {
		name    string
		args    args
		want    dmtfmodel.ComputerSystem
		wantErr bool
	}{
		{
			name: "positive case",
			args: args{
				deviceUUID: "someID",
			},
			want:    dmtfmodel.ComputerSystem{ID: "someID"},
			wantErr: false,
		},
		{
			name:    "not found",
			args:    args{deviceUUID: "invalid"},
			want:    dmtfmodel.ComputerSystem{},
			wantErr: true,
		},
		{
			name:    "invalid data",
			args:    args{deviceUUID: "falseData"},
			want:    dmtfmodel.ComputerSystem{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetComputeSystem(tt.args.deviceUUID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetComputeSystem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetComputeSystem() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSaveComputeSystem(t *testing.T) {
	config.SetUpMockConfig(t)
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()
	type args struct {
		computeServer dmtfmodel.ComputerSystem
		deviceUUID    string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "positive case",
			args:    args{computeServer: dmtfmodel.ComputerSystem{ID: "someID"}, deviceUUID: "someID"},
			wantErr: false,
		},
		{
			name:    "already exist",
			args:    args{computeServer: dmtfmodel.ComputerSystem{ID: "someID"}, deviceUUID: "someID"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SaveComputeSystem(tt.args.computeServer, tt.args.deviceUUID); (err != nil) != tt.wantErr {
				t.Errorf("SaveComputeSystem() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSaveChassis(t *testing.T) {
	config.SetUpMockConfig(t)
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()
	type args struct {
		chassis    dmtfmodel.Chassis
		deviceUUID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "positive case",
			args:    args{chassis: dmtfmodel.Chassis{ID: "someID"}, deviceUUID: "someID"},
			wantErr: false,
		},
		{
			name:    "already exist",
			args:    args{chassis: dmtfmodel.Chassis{ID: "someID"}, deviceUUID: "someID"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SaveChassis(tt.args.chassis, tt.args.deviceUUID); (err != nil) != tt.wantErr {
				t.Errorf("SaveChassis() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGenericSave(t *testing.T) {
	config.SetUpMockConfig(t)
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()
	type args struct {
		body  []byte
		table string
		key   string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "positive case",
			args:    args{body: []byte("someBody"), table: "someTable", key: "someKey"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := GenericSave(tt.args.body, tt.args.table, tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("GenericSave() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSaveRegistryFile(t *testing.T) {
	config.SetUpMockConfig(t)
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()
	type args struct {
		body  []byte
		table string
		key   string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "positive case",
			args:    args{body: []byte("someBody"), table: "someTable", key: "someKey"},
			wantErr: false,
		},
		{
			name:    "duplicate data case",
			args:    args{body: []byte("someBody"), table: "someTable", key: "someKey"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SaveRegistryFile(tt.args.body, tt.args.table, tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("SaveRegistryFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetRegistryFile(t *testing.T) {
	config.SetUpMockConfig(t)
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()
	mockData(t, common.OnDisk, "someTable", "someKey", "someData")
	type args struct {
		Table string
		key   string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 *errors.Error
	}{
		{
			name:  "positive case",
			args:  args{Table: "someTable", key: "someKey"},
			want:  "someData",
			want1: nil,
		},
		{
			name:  "not found",
			args:  args{Table: "notable", key: "someKey"},
			want:  "",
			want1: errors.PackError(errors.DBKeyAlreadyExist, "error while trying to get resource details: no data with the with key someKey found "),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := GetRegistryFile(tt.args.Table, tt.args.key)
			if got != tt.want {
				t.Errorf("GetRegistryFile() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeleteComputeSystem(t *testing.T) {
	config.SetUpMockConfig(t)

	sampleFile := filepath.Join(cwdDir, "sample.json")
	createFile(t, sampleFile, sampleData)
	config.Data.SearchAndFilterSchemaPath = sampleFile
	defer func() {
		os.Remove(sampleFile)
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()
	mockData(t, common.InMemory, "ComputerSystem", "/redfish/v1/systems/ef83e569-7336-492a-aaee-31c02d9db831:1", dmtfmodel.ComputerSystem{ID: "someID"})
	mockData(t, common.InMemory, "Systems", "/redfish/v1/systems/ef83e569-7336-492a-aaee-31c02d9db831:1", "some data")
	mockIndex(common.InMemory, "ProcessorSummary/Model", "ef83e569-7336-492a-aaee-31c02d9db831:1")
	type args struct {
		index int
		key   string
	}
	tests := []struct {
		name string
		args args
		want *errors.Error
	}{
		{
			name: "remove index",
			args: args{index: 19, key: "/redfish/v1/systems/ef83e569-7336-492a-aaee-31c02d9db831:1"},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DeleteComputeSystem(tt.args.index, tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteComputeSystem() = %v, want %v", got, tt.want)
			}
		})
	}

}

func TestDeleteSystem(t *testing.T) {
	config.SetUpMockConfig(t)
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()
	mockData(t, common.OnDisk, "System", "someKey", "some data")
	type args struct {
		key string
	}
	tests := []struct {
		name string
		args args
		want *errors.Error
	}{
		{
			name: "positive case",
			args: args{key: "someKey"},
			want: nil,
		},
		{
			name: "not found",
			args: args{key: "someOtherKey"},
			want: errors.PackError(errors.DBKeyNotFound, "error while trying to get compute details: no data with the with key someOtherKey found"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DeleteSystem(tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteSystem() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetTarget(t *testing.T) {
	config.SetUpMockConfig(t)
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()
	mockData(t, common.OnDisk, "System", "someKey", &Target{DeviceUUID: "someKey"})
	mockData(t, common.OnDisk, "System", "invalidData", "some data")
	type args struct {
		deviceUUID string
	}
	tests := []struct {
		name    string
		args    args
		want    *Target
		wantErr bool
	}{
		{
			name:    "positive case",
			args:    args{deviceUUID: "someKey"},
			want:    &Target{DeviceUUID: "someKey"},
			wantErr: false,
		},
		{
			name:    "not found",
			args:    args{deviceUUID: "noKey"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid data",
			args:    args{deviceUUID: "invalidData"},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetTarget(tt.args.deviceUUID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTarget() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTarget() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSaveIndex_WithError(t *testing.T) {
	config.SetUpMockConfig(t)
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()
	type args struct {
		searchForm map[string]interface{}
		table      string
		uuid       string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "document error",
			args:    args{searchForm: map[string]interface{}{"test": []int64{1}}, table: "test", uuid: "sample"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SaveIndex(tt.args.searchForm, tt.args.table, tt.args.uuid)
			if (err != nil) != tt.wantErr {
				t.Errorf("SaveIndex() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

var cwdDir, _ = os.Getwd()

func createFile(t *testing.T, fName, fContent string) {
	if err := ioutil.WriteFile(fName, []byte(fContent), 0644); err != nil {
		t.Fatal("error :failed to create a sample file for tests:", err)
	}
}

var sampleData = `{
	"searchKeys":[{"ProcessorSummary/Count":{"type":"float64"}}, {"ProcessorSummary/Model":{"type":"string"}}, {"SystemType":{"type":"string"}}, {"MemorySummary/TotalSystemMemoryGiB":{"type":"float64"}}, {"ProcessorSummary/sockets":{"type":"float64"}},
	   {"Processor/AccelerationFunctions/Compression":{"type":"string"}},{"Processor/AccelerationFunctions/AudioProcessing":{"type":"string"}},{"Processor/AccelerationFunctions/Encryption":{"type":"string"}},
	   {"Processor/AccelerationFunctions/PacketInspection":{"type":"string"}},{"Processor/AccelerationFunctions/PacketSwitch":{"type":"string"}},{"Processor/AccelerationFunctions/Scheduler":{"type":"string"}},
	   {"Processor/AccelerationFunctions/VideoProcessing":{"type":"string"}},{"NetworkInterfaces/NetworkPorts":{"type":"string"}},{"NetworkInterfaces/Model":{"type":"string"}},{"NetworkInterfaces/Bandwidth":{"type":"string"}},
	   {"FirmwareVersion":{"type":"string"}}],
	"conditionKeys":[ "eq","ne","gt","ge","lt","le"],
	"queryKeys":["filter"]
 }`

func TestSavePluginData(t *testing.T) {
	config.SetUpMockConfig(t)
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()
	type args struct {
		plugin Plugin
	}
	tests := []struct {
		name string
		args args
		want *errors.Error
	}{
		{
			name: "positive case",
			args: args{plugin: Plugin{ID: "someID"}},
			want: nil,
		},
		{
			name: "duplicate case",
			args: args{plugin: Plugin{ID: "someID"}},
			want: errors.PackError(errors.DBKeyAlreadyExist, "error while trying to save plugin data: error: data with key someID already exists"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SavePluginData(tt.args.plugin); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SavePluginData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAllSystems(t *testing.T) {
	config.SetUpMockConfig(t)
	mockData(t, common.OnDisk, "System", "someID", Target{DeviceUUID: "someID"})
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()
	tests := []struct {
		name  string
		want  []Target
		want1 *errors.Error
	}{
		{
			name:  "positive case",
			want:  []Target{Target{DeviceUUID: "someID"}},
			want1: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := GetAllSystems()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllSystems() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("GetAllSystems() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestDeletePluginData(t *testing.T) {
	config.SetUpMockConfig(t)
	mockData(t, common.OnDisk, "Plugin", "someID", "someData")
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()
	type args struct {
		key string
	}
	tests := []struct {
		name string
		args args
		want *errors.Error
	}{
		{
			name: "positive case",
			args: args{key: "someID"},
			want: nil,
		},
		{
			name: "no data found",
			args: args{key: "someOtherID"},
			want: errors.PackError(errors.DBKeyNotFound, "no data with the with key someOtherID found"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DeletePluginData(tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeletePluginData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeleteManagersData(t *testing.T) {
	config.SetUpMockConfig(t)
	mockData(t, common.InMemory, "Managers", "someID", "someData")
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()
	type args struct {
		key string
	}
	tests := []struct {
		name string
		args args
		want *errors.Error
	}{
		{
			name: "positive case",
			args: args{key: "someID"},
			want: nil,
		},
		{
			name: "no data found",
			args: args{key: "someOtherID"},
			want: errors.PackError(errors.DBKeyNotFound, "no data with the with key someOtherID found"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DeleteManagersData(tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteManagersData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpdateIndex(t *testing.T) {
	config.SetUpMockConfig(t)
	mockData(t, common.InMemory, "Systems", "someID", "someData")
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()
	type args struct {
		searchForm map[string]interface{}
		table      string
		uuid       string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "document indexing case",
			args:    args{searchForm: map[string]interface{}{"sample": "test"}, table: "sample", uuid: "test"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdateIndex(tt.args.searchForm, tt.args.table, tt.args.uuid); (err != nil) != tt.wantErr {
				t.Errorf("UpdateIndex() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateComputeSystem(t *testing.T) {
	config.SetUpMockConfig(t)
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()
	mockData(t, common.InMemory, "ComputerSystem", "/redfish/v1/systems/someID:1", dmtfmodel.ComputerSystem{ID: "someID"})
	type args struct {
		key         string
		computeData interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "positive case",
			args:    args{key: "/redfish/v1/systems/someID:1", computeData: dmtfmodel.ComputerSystem{ID: "someOtherID"}},
			wantErr: false,
		},
		{
			name:    "not found",
			args:    args{key: "/redfish/v1/systems/noID:1", computeData: dmtfmodel.ComputerSystem{ID: "someOtherID"}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdateComputeSystem(tt.args.key, tt.args.computeData); (err != nil) != tt.wantErr {
				t.Errorf("UpdateComputeSystem() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetResourceDetails(t *testing.T) {
	config.SetUpMockConfig(t)
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()
	mockData(t, common.InMemory, "ComputerSystem", "someKey", "someData")
	type args struct {
		key string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 *errors.Error
	}{
		{
			name:  "success case",
			args:  args{key: "someKey"},
			want:  "someData",
			want1: nil,
		},
		{
			name:  "not found",
			args:  args{key: "someOtherKey"},
			want:  "",
			want1: errors.PackError(errors.DBKeyNotFound, "error while trying to get resource details: no data with the with key someOtherKey found"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := GetResourceDetails(tt.args.key)
			if got != tt.want {
				t.Errorf("GetResourceDetails() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("GetResourceDetails() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestSystemOperation(t *testing.T) {
	// testing  the Add SystemOpearation use case
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
	var systemOperation = SystemOperation{
		Operation: "Rediscovery",
	}
	systemURI := "/redfish/v1/System/uuid:1"
	err := systemOperation.AddSystemOperationInfo(systemURI)
	assert.Nil(t, err, "err should be nil")

	// testing the get system operation
	data, err := GetSystemOperationInfo(systemURI)
	assert.Nil(t, err, "err should be nil")
	assert.Equal(t, "Rediscovery", data.Operation)

	_, err = GetSystemOperationInfo("systemURI")
	assert.NotNil(t, err, "Error Should not be nil")

	//testing the delete operation
	err = DeleteSystemOperationInfo(systemURI)
	assert.Nil(t, err, "err should be nil")

	err = DeleteSystemOperationInfo("systemURI")
	assert.NotNil(t, err, "Error Should not be nil")
}

func TestSystemReset(t *testing.T) {
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
	// testing  the Add SystemReset use case
	systemURI := "/redfish/v1/System/uuid:1"
	err := AddSystemResetInfo(systemURI, "ForceRestart")
	assert.Nil(t, err, "err should be nil")

	// testing the get system operation
	data, err := GetSystemResetInfo(systemURI)
	assert.Nil(t, err, "err should be nil")
	assert.Equal(t, "ForceRestart", data["ResetType"])

	_, err = GetSystemResetInfo("systemURI")
	assert.NotNil(t, err, "Error Should not be nil")

	//testing the delete operation
	err = DeleteSystemResetInfo(systemURI)
	assert.Nil(t, err, "err should be nil")

	err = DeleteSystemResetInfo("systemURI")
	assert.NotNil(t, err, "Error Should not be nil")
}
func TestAggregationSource(t *testing.T) {
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

	aggregationSourceURI := "/redfish/v1/AggregationService/AggregationSource/12345677651245-12341"
	req := AggregationSource{
		HostName: "localhost:9091",
		UserName: "admin",
		Password: []byte("password"),
		Links: map[string]interface{}{
			"Oem": map[string]string{
				"PluginID": "GRF",
			},
		},
	}
	err := AddAggregationSource(req, aggregationSourceURI)
	assert.Nil(t, err, "err should be nil")
	err = AddAggregationSource(req, aggregationSourceURI)
	assert.NotNil(t, err, "Error Should not be nil")
	keys, dbErr := GetAllKeysFromTable("AggregationSource")
	assert.Nil(t, dbErr, "err should be nil")
	assert.Equal(t, 1, len(keys), "length should be matching")
	data, err := GetAggregationSourceInfo(aggregationSourceURI)
	assert.Nil(t, err, "err should be nil")
	assert.Equal(t, data.HostName, req.HostName)
	assert.Equal(t, data.UserName, req.UserName)
	_, err = GetAggregationSourceInfo("/redfish/v1/AggregationService/AggregationSource/12345677651245-123433")
	assert.NotNil(t, err, "Error Should not be nil")
	err = UpdateAggregtionSource(req, aggregationSourceURI)
	assert.Nil(t, err, "err should be nil")
	err = UpdateAggregtionSource(req, "/redfish/v1/AggregationService/AggregationSource/12345677651245-123433")
	assert.NotNil(t, err, "Error Should not be nil")
	data, err = GetAggregationSourceInfo(aggregationSourceURI)
	assert.Nil(t, err, "err should be nil")
	assert.Equal(t, data.HostName, req.HostName)
	assert.Equal(t, data.UserName, req.UserName)
	keys, err = GetAllMatchingDetails("AggregationSource", "/redfish/v1/AggregationService/AggregationSource/12345677651245-", common.OnDisk)
	assert.Nil(t, err, "err should be nil")
	assert.Equal(t, 1, len(keys), "length should be matching")
	err = DeleteAggregationSource(aggregationSourceURI)
	assert.Nil(t, err, "err should be nil")
	err = DeleteAggregationSource(aggregationSourceURI)
	assert.NotNil(t, err, "Error Should not be nil")
}

func TestUpdatePluginData(t *testing.T) {
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
	validPasswordEnc := getEncryptedKey(t, []byte("password"))

	pluginData := Plugin{
		IP:                "localhost",
		Port:              "45001",
		Username:          "admin",
		Password:          validPasswordEnc,
		ID:                "GRF",
		PluginType:        "RF-GENERIC",
		PreferredAuthType: "BasicAuth",
		ManagerUUID:       "123453414-1223441",
	}
	mockData(t, common.OnDisk, "Plugin", "GRF", pluginData)
	pluginData.Username = "admin1"
	pluginData.IP = "9.9.9.0"
	err := UpdatePluginData(pluginData, "GRF")
	assert.Nil(t, err, "err should be nil")
	err = UpdatePluginData(pluginData, "GRF1")
	assert.NotNil(t, err, "Error Should not be nil")
}

func TestUpdateSystemData(t *testing.T) {
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
	req := SaveSystem{
		UserName:   "admin",
		Password:   []byte("12345"),
		DeviceUUID: "1234567678-12331",
		PluginID:   "GRF",
	}
	mockData(t, common.OnDisk, "System", "1234567678-12331", &req)
	req.UserName = "admin"
	req.Password = []byte("12346")
	dbErr := UpdateSystemData(req, "1234567678-12331")
	assert.Nil(t, dbErr, "err should be nil")
	data, err := GetTarget("1234567678-12331")
	assert.Nil(t, err, "err should be nil")
	assert.Equal(t, req.UserName, data.UserName, "UserName should be same")
	assert.Equal(t, req.Password, data.Password, "Password should be same")
	assert.Equal(t, req.PluginID, data.PluginID, "PluginID should be same")
	assert.Equal(t, req.DeviceUUID, data.DeviceUUID, "DeviceUUID should be same")

	dbErr = UpdateSystemData(req, "1234567678-12332")
	assert.NotNil(t, dbErr, "Error Should not be nil")
}
