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
// under the License

//Package fabmodel ...

package smodel

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/ODIM-Project/ODIM/lib-persistence-manager/persistencemgr"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	"github.com/stretchr/testify/assert"
)

var plugin = Plugin{
	IP:                "localhost",
	Port:              "9091",
	Username:          "admin",
	ID:                "CFM",
	PreferredAuthType: "XAuthTOken",
}

var target = Target{
	ManagerAddress: "10.24.0.14",
	Password:       []byte("Password"),
	UserName:       "admin",
	DeviceUUID:     "uuid",
	PluginID:       "GRF",
}

var invalidTarget = "Target"

func mockSystemIndex(uuid string) error {
	connPool, err := common.GetDBConnection(common.InMemory)
	if err != nil {
		return fmt.Errorf("error while trying to connecting to DB: %v", err.Error())
	}
	var indexData = map[string]interface{}{
		"ProcessorSummary/Model":  "Intel",
		"ProcessorSummary/Count":  2,
		"Storage/Drives/Capacity": []float64{40},
		"Storage/Drives/Type":     []string{"HDD", "HDD"},
	}
	if err := connPool.CreateIndex(indexData, "/redfish/v1/Systems/"+uuid); err != nil {
		return fmt.Errorf("error while creating  the index: %v", err.Error())
	}
	return nil

}

func mockPluginData(t *testing.T) error {
	plugin.Password = getEncryptedKey(t, []byte("12345"))
	connPool, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		return fmt.Errorf("error while trying to connecting to DB: %v", err.Error())
	}
	if err = connPool.Create("Plugin", "CFM", plugin); err != nil {
		return fmt.Errorf("error while trying to create new %v resource: %v", "Plugin", err.Error())
	}
	return nil
}

func mockInvalidPluginData() error {
	connPool, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		return fmt.Errorf("error while trying to connecting to DB: %v", err.Error())
	}
	if err = connPool.Create("Plugin", "Invalid", "plugin"); err != nil {
		return fmt.Errorf("error while trying to create new %v resource: %v", "Plugin", err.Error())
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

func mockTarget() error {
	connPool, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		return fmt.Errorf("error while trying to connecting to DB: %v", err.Error())
	}

	const table string = "System"
	//Save data into Database
	if err = connPool.Create(table, target.DeviceUUID, target); err != nil {
		return err
	}
	return nil
}
func mockInvalidTarget() error {
	connPool, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		return fmt.Errorf("error while trying to connecting to DB: %v", err.Error())
	}

	const table string = "System"
	//Save data into Database
	if err = connPool.Create(table, "falseData", "target"); err != nil {
		return err
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

func mockData(t *testing.T, dbType common.DbType, table, id string, data interface{}) {
	connPool, err := common.GetDBConnection(dbType)
	if err != nil {
		t.Fatalf("error: mockData() failed to DB connection: %v", err)
	}
	if err = connPool.Create(table, id, data); err != nil {
		t.Fatalf("error: mockData() failed to create entry %s-%s: %v", table, id, err)
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

func TestGetSystemByUUID(t *testing.T) {
	config.SetUpMockConfig(t)
	ctx := mockContext()
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
	body := `{"Id":"1","Status":{"State":"Enabled"}}`
	table := "ComputerSystem"
	key := "/redfish/v1/Systems/uuid.1"
	GenericSave(ctx, []byte(body), table, key)
	data, _ := GetSystemByUUID(ctx, "/redfish/v1/Systems/uuid.1")
	assert.Equal(t, data, body, "should be same")

	JSONUnmarshalFunc = func(data []byte, v interface{}) error {
		return &errors.Error{}
	}
	_, err := GetSystemByUUID(ctx, "/redfish/v1/Systems/uuid.1")
	assert.NotNil(t, err, "There should be an error")

	_, err = GetSystemByUUID(ctx, "/redfish/v1/Systems/uuid")
	assert.NotNil(t, err, "There should be an error")

	JSONUnmarshalFunc = func(data []byte, v interface{}) error {
		return json.Unmarshal(data, v)
	}
	GetDBConnectionFunc = func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
		return nil, &errors.Error{}

	}
	_, err = GetSystemByUUID(ctx, "/redfish/v1/Systems/uuid")
	assert.NotNil(t, err, "There should be an error")

	GetDBConnectionFunc = func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
		return common.GetDBConnection(dbFlag)

	}
}

func TestGetTarget(t *testing.T) {
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
	mockTarget()
	resp, _ := GetTarget("uuid")
	assert.Equal(t, resp.ManagerAddress, target.ManagerAddress, "should be same")

	_, err := GetTarget("uuid1")
	assert.NotNil(t, err, "There should be an error")
}

func TestGetTarget_negative(t *testing.T) {
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
	mockInvalidTarget()
	_, err := GetTarget("falseData")
	assert.NotNil(t, err, "There should be an error")
}

func TestGenericSave(t *testing.T) {
	config.SetUpMockConfig(t)
	ctx := mockContext()
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

	body := []byte(`body`)
	table := "EthernetInterfaces"
	key := "/redfish/v1/Managers/uuid.1/EthernetInterfaces/1"
	err := GenericSave(ctx, body, table, key)
	assert.Nil(t, err, "There should be no error")

	GetDBConnectionFunc = func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
		return nil, &errors.Error{}

	}
	err = GenericSave(ctx, body, table, key)
	assert.NotNil(t, err, "There should be an error")

	_, err = GetResource(ctx, table, key)
	assert.NotNil(t, err, "There should be an error")

	GetDBConnectionFunc = func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
		return common.GetDBConnection(dbFlag)

	}
	JSONUnmarshalFunc = func(data []byte, v interface{}) error {
		return &errors.Error{}
	}
	_, err = GetResource(ctx, table, key)
	assert.NotNil(t, err, "There should be an error")
	JSONUnmarshalFunc = func(data []byte, v interface{}) error {
		return json.Unmarshal(data, v)
	}

	data, err := GetResource(ctx, table, key)
	assert.Nil(t, err, "There should be no error")
	assert.Equal(t, data, string(body), "should be same")
}

func TestGetAllkeysFromTable(t *testing.T) {
	config.SetUpMockConfig(t)
	ctx := mockContext()
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

	body := []byte(`body`)
	table := "EthernetInterfaces"
	key := "/redfish/v1/Managers/uuid.1/EthernetInterfaces/1"
	err := GenericSave(ctx, body, table, key)
	assert.Nil(t, err, "There should be no error")

	allKeys, err := GetAllKeysFromTable(table)
	assert.Nil(t, err, "There should be no error")
	assert.Equal(t, len(allKeys), 1, "There should be one entry in DB")
}

func TestGetResourceNegativeTestCases(t *testing.T) {
	ctx := mockContext()
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
	// without db configuration
	table := "EthernetInterfaces"
	key := "/redfish/v1/Managers/uuid.1/EthernetInterfaces/1"

	_, err := GetResource(ctx, table, key)
	assert.NotNil(t, err, "There should be an error")

	// if key not present
	config.SetUpMockConfig(t)
	table = "Ethernet"
	key = "/redfish/v1/Managers/uuid.1/Ethernets/1"

	_, err = GetResource(ctx, table, key)
	assert.NotNil(t, err, "There should be an error")

}

func TestGetStorageList(t *testing.T) {
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
	err := mockSystemIndex("uuid.1")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	list, err := GetStorageList("Storage/Drives/Capacity", "le", 40, false)
	assert.Nil(t, err, "There should be no error")
	assert.Equal(t, len(list), 1, "Length should be equal")
	assert.Equal(t, list[0], "/redfish/v1/Systems/uuid.1", "system uri should be same")
	list, err = GetStorageList("Storage/Drives/Capacity", "eq", 40, false)
	assert.Nil(t, err, "There should be no error")
	assert.Equal(t, len(list), 1, "Length should be equal")
	assert.Equal(t, list[0], "/redfish/v1/Systems/uuid.1", "system uri should be same")

}

func TestGetString(t *testing.T) {
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
	err := mockSystemIndex("uuid.1")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	list, err := GetString("ProcessorSummary/Model", "Intel", false)
	assert.Nil(t, err, "There should be no error")
	assert.Equal(t, len(list), 1, "Length should be equal")
	assert.Equal(t, list[0], "/redfish/v1/Systems/uuid.1", "system uri should be same")
}

func TestGetRange(t *testing.T) {
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
	err := mockSystemIndex("uuid.1")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	list, err := GetRange("ProcessorSummary/Count", 2, 2, false)
	assert.Nil(t, err, "There should be no error")
	assert.Equal(t, len(list), 1, "Length should be equal")
	assert.Equal(t, list[0], "/redfish/v1/Systems/uuid.1", "system uri should be same")
}

func TestSystemReset(t *testing.T) {
	// testing  the Add SystemReset use case
	common.SetUpMockConfig()
	ctx := mockContext()
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
	systemURI := "/redfish/v1/System/uuid.1"
	err := AddSystemResetInfo(ctx, systemURI, "ForceRestart")
	assert.Nil(t, err, "err should be nil")

	// testing the get system operation
	data, err := GetSystemResetInfo(ctx, systemURI)
	assert.Nil(t, err, "err should be nil")
	assert.Equal(t, "ForceRestart", data["ResetType"])

	_, err = GetSystemResetInfo(ctx, "systemURI")
	assert.NotNil(t, err, "Error Should not be nil")

}

func TestDeleteVolume(t *testing.T) {
	ctx := mockContext()
	defer func() {
		common.TruncateDB(common.InMemory)
	}()
	mockData(t, common.InMemory, "Volumes", "/redfish/v1/Systems/ef83e569-7336-492a-aaee-31c02d9db831.1/Storage/1/Volume/1", "")
	tests := []struct {
		name string
		key  string
		want *errors.Error
	}{
		{
			name: "Positive case",
			key:  "/redfish/v1/Systems/ef83e569-7336-492a-aaee-31c02d9db831.1/Storage/1/Volume/1",
			want: nil,
		},
		{
			name: "not found",
			key:  "/redfish/v1/Systems/ef83e569-7336-492a-aaee-31c02d9db831.1/Storage/1/Volume/2",
			want: errors.PackError(errors.DBKeyNotFound, "error while trying to get voulme details: no data with the with key /redfish/v1/Systems/ef83e569-7336-492a-aaee-31c02d9db831.1/Storage/1/Volume/2 found"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DeleteVolume(ctx, tt.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteVolume() = %v, want %v", got, tt.want)
			}
		})
	}

}

func TestFind(t *testing.T) {
	defer func() {
		common.TruncateDB(common.InMemory)
	}()
	mockData(t, common.InMemory, "Volumes", "/redfish/v1/Systems/ef83e569-7336-492a-aaee-31c02d9db831.1/Storage/1/Volume/1", "")

	err := Find("Volumes", "/redfish/v1/Systems/ef83e569-7336-492a-aaee-31c02d9db831.1/Storage/1/Volume/1", "")
	assert.NotNil(t, err, "should be an error ")

	JSONUnmarshalFunc = func(data []byte, v interface{}) error {
		return &errors.Error{}
	}
	err = Find("Volumes", "/redfish/v1/Systems/ef83e569-7336-492a-aaee-31c02d9db831.1/Storage/1/Volume/1", "")
	assert.NotNil(t, err, "should be an error ")

	err = Find("Volumes", "", "")
	assert.NotNil(t, err, "should be an error, Invalid ID ")

	JSONUnmarshalFunc = func(data []byte, v interface{}) error {
		return json.Unmarshal(data, v)
	}
	GetDBConnectionFunc = func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
		return nil, &errors.Error{}
	}
	err = Find("Volumes", "/redfish/v1/Systems/ef83e569-7336-492a-aaee-31c02d9db831.1/Storage/1/Volume/1", "")
	assert.NotNil(t, err, "should be an error ")
	GetDBConnectionFunc = func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
		return common.GetDBConnection(dbFlag)
	}

}

func TestFindAll(t *testing.T) {
	defer func() {
		common.TruncateDB(common.InMemory)
	}()
	mockData(t, common.InMemory, "Volumes", "/redfish/v1/Systems/ef83e569-7336-492a-aaee-31c02d9db831.1/Storage/1/Volume/1", "")

	_, err := FindAll("Volumes", "/redfish/v1/Systems/ef83e569-7336-492a-aaee-31c02d9db831.1/Storage/1/Volume/1")
	assert.Nil(t, err, "should be no error ")

	GetDBConnectionFunc = func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
		return nil, &errors.Error{}
	}
	_, err = FindAll("Volumes", "/redfish/v1/Systems/ef83e569-7336-492a-aaee-31c02d9db831.1/Storage/1/Volume/1")
	assert.NotNil(t, err, "should be an error ")

	GetDBConnectionFunc = func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
		return common.GetDBConnection(dbFlag)
	}
	scanFunc = func(cp *persistencemgr.ConnPool, key string) ([]interface{}, error) {
		return nil, &errors.Error{}
	}
	_, err = FindAll("Volumes", "/redfish/v1/Systems/ef83e569-7336-492a-aaee-31c02d9db831.1/Storage/1/Volume/1")
	assert.NotNil(t, err, "should be an error ")

	scanFunc = func(cp *persistencemgr.ConnPool, key string) ([]interface{}, error) {
		return []interface{}{"dummy"}, nil
	}
	_, err = FindAll("Volumes", "/redfish/v1/Systems/ef83e569-7336-492a-aaee-31c02d9db831.1/Storage/1/Volume/1")
	assert.Nil(t, err, "should be no error ")

}

func TestGetAllKeysFromTable(t *testing.T) {
	ctx := mockContext()
	defer func() {
		common.TruncateDB(common.InMemory)
	}()
	mockData(t, common.InMemory, "Volumes", "/redfish/v1/Systems/ef83e569-7336-492a-aaee-31c02d9db831.1/Storage/1/Volume/1", "")
	GetDBConnectionFunc = func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
		return nil, &errors.Error{}
	}
	_, err := GetAllKeysFromTable("Volumes")
	assert.NotNil(t, err, "should be an error ")
	_, err = GetPluginData("Volumes")
	assert.NotNil(t, err, "should be an error ")
	_, err = GetTarget("Volumes")
	assert.NotNil(t, err, "should be an error ")
	_, err = GetStorageList("Volumes", "", 0.5, true)
	assert.NotNil(t, err, "should be an error ")
	_, err = GetString("Volumes", "", true)
	assert.NotNil(t, err, "should be an error ")

	_, err = GetRange("Volumes", 0, 100, true)
	assert.NotNil(t, err, "should be an error ")

	err = AddSystemResetInfo(ctx, "Volumes", "rese")
	assert.NotNil(t, err, "should be an error ")

	_, err = GetSystemResetInfo(ctx, "Volumes")
	assert.NotNil(t, err, "should be an error ")
	err = DeleteVolume(ctx, "Volumes")
	assert.NotNil(t, err, "should be an error ")
	GetDBConnectionFunc = func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
		return common.GetDBConnection(dbFlag)
	}

}
