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

// Package tmodel ....
package tmodel

import (
	"context"
	"testing"

	"github.com/ODIM-Project/ODIM/lib-persistence-manager/persistencemgr"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	"github.com/stretchr/testify/assert"
)

func TestGetResource(t *testing.T) {
	config.SetUpMockConfig(t)
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()
	mockData(t, common.InMemory, "someTable", "someKey", "someData")
	mockData(t, common.InMemory, "someTable", "invalidData", 235)

	type args struct {
		Table string
		key   string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Positive Test",
			args: args{key: "someKey", Table: "someTable"},
			want: "someData",
		},
		{
			name: "Negative Test",
			args: args{key: "invalidData", Table: "someTable"},
			want: "",
		},

		{
			name: "Negative Test",
			args: args{key: "invalid", Table: "someTable"},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := GetResource(tt.args.Table, tt.args.key, common.InMemory)
			if got != tt.want {
				t.Errorf("GetResource() got = %v, want %v", got, tt.want)
			}

		})
	}
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

func TestGetAllKeysFromTable(t *testing.T) {
	config.SetUpMockConfig(t)
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()

	// validPassword := []byte("password")
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

	type args struct {
		table  string
		dbtype common.DbType
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Positive Case ",
			args: args{table: "Plugin", dbtype: common.OnDisk},
			want: []string{"validPlugin"},
		},
		{
			name: "Negative Case ",
			args: args{table: "", dbtype: common.OnDisk},
			want: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := GetAllKeysFromTable(tt.args.table, tt.args.dbtype)

			if len(got) != len(tt.want) {
				t.Errorf("GetAllKeysFromTable() = %v, want %v", got, tt.want)
			}
		})
	}
	GetDBConnectionFunc = func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
		return nil, &errors.Error{}
	}
	_, err := GetAllKeysFromTable("System", common.OnDisk)
	assert.NotNil(t, err, "There should be an error ")

	_, err = GetResource("System", "dummy", common.OnDisk)
	GetDBConnectionFunc = func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
		return common.GetDBConnection(dbFlag)
	}

}
func getEncryptedKey(t *testing.T, key []byte) []byte {
	cryptedKey, err := common.EncryptWithPublicKey(key)
	if err != nil {
		t.Fatalf("error: failed to encrypt data: %v", err)
	}
	return cryptedKey
}

func TestGetPluginData(t *testing.T) {
	config.SetUpMockConfig(t)
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()

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

	_, err := GetPluginData("validPlugin")
	assert.Nil(t, err, "There should be no error ")

	// Invalid Plugin id
	_, err = GetPluginData("validPlugin1")
	assert.NotNil(t, err, "There should be an error ")

	// Invalid password
	pluginData.Password = invalidPassword
	mockData(t, common.OnDisk, "Plugin", "invalidPassword", pluginData)
	_, err = GetPluginData("invalidPassword")
	assert.NotNil(t, err, "There should be an error ")

	// Invalid password
	pluginData.Password = invalidPassword
	mockData(t, common.OnDisk, "Plugin", "invalidData", "dummy")
	_, err = GetPluginData("invalidData")
	assert.NotNil(t, err, "There should be an error ")

	GetDBConnectionFunc = func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
		return nil, &errors.Error{}
	}
	_, err = GetPluginData("invalidData")
	assert.NotNil(t, err, "There should be an error ")
	GetDBConnectionFunc = func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
		return common.GetDBConnection(dbFlag)
	}

}

func TestGetTarget(t *testing.T) {
	config.SetUpMockConfig(t)
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
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
	}
	mockData(t, common.OnDisk, "System", "system_id", pluginData)
	_, err := GetTarget("system_id")
	assert.Nil(t, err, "There should be no error ")

	_, err = GetTarget("invalid")
	assert.NotNil(t, err, "There should be an error ")
	mockData(t, common.OnDisk, "System", "invalid", "dummy")
	_, err = GetTarget("invalid")
	assert.NotNil(t, err, "There should be no error ")
	GetDBConnectionFunc = func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
		return nil, &errors.Error{}
	}
	_, err = GetTarget("system_id")
	assert.NotNil(t, err, "There should be an error ")
	GetDBConnectionFunc = func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
		return common.GetDBConnection(dbFlag)
	}

}

func TestGenericSave(t *testing.T) {
	config.SetUpMockConfig(t)
	ctx := mockContext()
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()

	// validPassword := []byte("password")
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
	mockData(t, common.OnDisk, "System", "system_id", pluginData)
	err := GenericSave(ctx, []byte("system_id"), "System", "dummy")
	assert.Nil(t, err, "There should be no error ")

	GetDBConnectionFunc = func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
		return nil, &errors.Error{}
	}
	GenericSave(ctx, []byte("system_id"), "System", "dummy")

	GetDBConnectionFunc = func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
		return common.GetDBConnection(dbFlag)
	}

}

func mockContext() context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, common.TransactionID, "xyz")
	ctx = context.WithValue(ctx, common.ActionID, "001")
	ctx = context.WithValue(ctx, common.ActionName, "xyz")
	ctx = context.WithValue(ctx, common.ThreadID, "0")
	ctx = context.WithValue(ctx, common.ThreadName, "xyz")
	ctx = context.WithValue(ctx, common.ProcessName, "xyz")
	return ctx
}
