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
package fabmodel

import (
	"context"
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
	PluginType:        "Fabric",
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

func getEncryptedKey(t *testing.T, key []byte) []byte {
	cryptedKey, err := common.EncryptWithPublicKey(key)
	if err != nil {
		t.Fatalf("error: failed to encrypt data: %v", err)
	}
	return cryptedKey
}

func mockFabricData(fabricID, pluginID string) error {

	connPool, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		return fmt.Errorf("error while trying to connecting to DB: %v", err.Error())
	}
	fab := Fabric{
		FabricUUID: fabricID,
		PluginID:   pluginID,
	}
	if err = connPool.Create("Fabric", fabricID, fab); err != nil {
		return fmt.Errorf("error while trying to create new %v resource: %v", "fabric", err.Error())
	}
	return nil
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
		PluginType:        "Fabric",
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
		name                string
		args                args
		exec                func(*Plugin)
		want                Plugin
		GetDBConnectionFunc func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error)
		wantErr             bool
	}{
		{
			name: "Positive Case",
			args: args{pluginID: "validPlugin"},
			exec: func(want *Plugin) {
				want.Password = validPassword
			},
			want:    pluginData,
			wantErr: false,
			GetDBConnectionFunc: func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
				return common.GetDBConnection(dbFlag)
			},
		},
		{
			name:    "Negative Case - Non-existent plugin",
			args:    args{pluginID: "notFound"},
			exec:    nil,
			want:    Plugin{},
			wantErr: true,
			GetDBConnectionFunc: func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
				return common.GetDBConnection(dbFlag)
			},
		},
		{
			name:    "Negative Case - Invalid plugin data",
			args:    args{pluginID: "invalidPluginData"},
			exec:    nil,
			want:    Plugin{},
			wantErr: true,
			GetDBConnectionFunc: func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
				return common.GetDBConnection(dbFlag)
			},
		},
		{
			name:                "Negative Case - Get DB connection",
			args:                args{pluginID: "invalidPassword"},
			exec:                nil,
			want:                Plugin{},
			wantErr:             true,
			GetDBConnectionFunc: func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) { return nil, &errors.Error{} },
		},
		{
			name: "Negative Case - Plugin with invalid password",
			args: args{pluginID: "invalidPassword"},
			exec: nil,
			want: Plugin{},
			GetDBConnectionFunc: func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
				return common.GetDBConnection(dbFlag)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		if tt.exec != nil {
			tt.exec(&tt.want)
		}
		t.Run(tt.name, func(t *testing.T) {
			GetDBConnectionFunc = tt.GetDBConnectionFunc
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

func TestGetAllFabricPluginDetails(t *testing.T) {
	config.SetUpMockConfig(t)
	defer func() {
		err := common.TruncateDB(common.OnDisk)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()
	mockPluginData(t)
	resp, _ := GetAllFabricPluginDetails(context.TODO())
	assert.Equal(t, len(resp), 1, "should be same")
	GetDBConnectionFunc = func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
		return nil, &errors.Error{}
	}
	_, err := GetAllFabricPluginDetails(context.TODO())
	assert.NotNil(t, err, "There should be an error")
	GetDBConnectionFunc = func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
		return common.GetDBConnection(dbFlag)
	}

}

func TestAddFabricData(t *testing.T) {
	config.SetUpMockConfig(t)
	defer func() {
		err := common.TruncateDB(common.OnDisk)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()
	var fab1 = Fabric{
		FabricUUID: "12345",
		PluginID:   "CFM",
	}
	err := fab1.AddFabricData("12345")
	assert.Equal(t, nil, err, "should be same")

	// Adding Duplicate Fabrics Data
	err = fab1.AddFabricData("12345")
	assert.Equal(t, "warning: skipped saving of duplicate data with key 12345", err.Error(), "should be same")
	GetDBConnectionFunc = func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
		return nil, &errors.Error{}
	}
	err = fab1.AddFabricData("12345")
	assert.NotNil(t, err, "There should be an error")
	GetDBConnectionFunc = func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
		return common.GetDBConnection(dbFlag)
	}
}

func TesGetManagingPluginIDForFabricID(t *testing.T) {
	config.SetUpMockConfig(t)
	defer func() {
		common.TruncateDB(common.OnDisk)
	}()
	var fab = Fabric{
		FabricUUID: "12345",
		PluginID:   "CFM",
	}
	err := fab.AddFabricData("12345")
	assert.Equal(t, nil, err, "there should be no error")

	// positive test case
	fabric, err := GetManagingPluginIDForFabricID(fab.FabricUUID,context.TODO())
	assert.Equal(t, nil, err, "there should be no error")
	assert.Equal(t, "12345", fabric.FabricUUID, "fabric uuid should be 12345")
	assert.Equal(t, "CFM", fabric.PluginID, "plugin id should be CFM")

	// negative test case
	fabric, err = GetManagingPluginIDForFabricID("54321",context.TODO())
	assert.NotNil(t, err, "there should be an error")

	GetDBConnectionFunc = func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
		return nil, &errors.Error{}
	}
	_, err = GetManagingPluginIDForFabricID("54321",context.TODO())

	assert.NotNil(t, err, "There should be an error")
	GetDBConnectionFunc = func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
		return common.GetDBConnection(dbFlag)
	}

}

func TestGetAllFabrics(t *testing.T) {
	config.SetUpMockConfig(t)
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()
	err := mockFabricData("d72dade0-c35a-984c-4859-1108132d72da", "CFM")
	if err != nil {
		t.Fatalf("Error in creating mockFabricData :%v", err)
	}
	tests := []struct {
		name string
		want []Fabric
	}{
		{
			name: "positive case",
			want: []Fabric{Fabric{FabricUUID: "d72dade0-c35a-984c-4859-1108132d72da", PluginID: "CFM"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := GetAllTheFabrics(context.TODO())
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllTheFabrics() got = %v, want %v", got, tt.want)
			}
		})
	}
	GetDBConnectionFunc = func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
		return nil, &errors.Error{}
	}
	_, err = GetAllTheFabrics(context.TODO())
	assert.NotNil(t, err, "There should be an error")
	GetDBConnectionFunc = func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
		return common.GetDBConnection(dbFlag)
	}

}

func TestFabric_RemoveFabricData(t *testing.T) {
	var fab1 = Fabric{
		FabricUUID: "12345",
		PluginID:   "CFM",
	}
	config.SetUpMockConfig(t)
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()

	fab1.AddFabricData("12345")
	fab1.RemoveFabricData("12345")

	_, err := GetManagingPluginIDForFabricID("12345",context.TODO())
	assert.NotNil(t, err, "There should be an error ")
	mockFabricData("12345", "AFC")
	_, err = GetManagingPluginIDForFabricID("12345",context.TODO())
	assert.Nil(t, err, "There should no error ")

	GetDBConnectionFunc = func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
		return nil, &errors.Error{}
	}
	err = fab1.RemoveFabricData("12345")
	assert.NotNil(t, err, "There should be an error")
	GetDBConnectionFunc = func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
		return common.GetDBConnection(dbFlag)
	}

}
