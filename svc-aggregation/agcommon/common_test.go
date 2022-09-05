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

package agcommon

import (
	"bytes"
	"fmt"
	"net"
	"testing"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	"github.com/ODIM-Project/ODIM/svc-aggregation/agmodel"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestAddConnectionMethods(t *testing.T) {
	var e = DBInterface{
		GetAllKeysFromTableInterface: stubGetAllkeys,
		GetConnectionMethodInterface: stubGetConnectionMethod,
		AddConnectionMethodInterface: stubAddConnectionMethod,
		DeleteInterface:              stubDeleteConnectionMethod,
	}
	config.SetUpMockConfig(t)
	err := e.AddConnectionMethods(config.Data.ConnectionMethodConf)
	assert.Nil(t, err, "err should be nil")
}

var connectionMethod = []string{"/redfish/v1/AggregationService/ConnectionMethods/1234567545691234f",
	"/redfish/v1/AggregationService/ConnectionMethods/1234567545691234g",
	"/redfish/v1/AggregationService/ConnectionMethods/1234567545691234h"}

func stubGetAllkeys(tableName string) ([]string, error) {
	return connectionMethod, nil
}

func stubGetConnectionMethod(key string) (agmodel.ConnectionMethod, *errors.Error) {
	if key == "/redfish/v1/AggregationService/ConnectionMethods/1234567545691234f" {
		return agmodel.ConnectionMethod{
			ConnectionMethodType:    "Redfish",
			ConnectionMethodVariant: "Compute:BasicAuth:GRF:1.0.0",
			Links: agmodel.Links{
				AggregationSources: []agmodel.OdataID{},
			},
		}, nil
	}

	if key == "/redfish/v1/AggregationService/ConnectionMethods/1234567545691234g" {
		return agmodel.ConnectionMethod{
			ConnectionMethodType:    "Redfish",
			ConnectionMethodVariant: "Fabric:XAuthToken:FabricPlugin:1.0.0",
			Links: agmodel.Links{
				AggregationSources: []agmodel.OdataID{
					agmodel.OdataID{OdataID: "/redfish/v1/AggregationService/AggregationSources/1234656881231fg1"},
				},
			},
		}, nil
	}
	return agmodel.ConnectionMethod{
		ConnectionMethodType:    "Redfish",
		ConnectionMethodVariant: "Storage:BasicAuth:Stg1:1.0.0",
		Links: agmodel.Links{
			AggregationSources: []agmodel.OdataID{},
		},
	}, nil
}

func stubAddConnectionMethod(data agmodel.ConnectionMethod, key string) *errors.Error {
	ConnectionMethod := agmodel.ConnectionMethod{
		ConnectionMethodType:    "Redfish",
		ConnectionMethodVariant: "Compute:BasicAuth:GRF:1.0.0",
		Links: agmodel.Links{
			AggregationSources: []agmodel.OdataID{},
		},
	}
	connectionMethodURI := "/redfish/v1/AggregationService/ConnectionMethods/" + uuid.NewV4().String()

	connPool, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		return errors.PackError(err.ErrNo(), "error while trying to connecting to DB: ", err.Error())

	}
	if err = connPool.Create("ConnectionMethod", connectionMethodURI, ConnectionMethod); err != nil {
		return errors.PackError(err.ErrNo(), "error while trying to create new  resource :ConnectionMethod  ", err.Error())
	}
	return nil
}

func stubDeleteConnectionMethod(table, key string, dbtype common.DbType) *errors.Error {

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

func TestGetStorageResources(t *testing.T) {
	config.SetUpMockConfig(t)
	storageURI := "/redfish/v1/Systems/12345677651245-12341/Storage"
	GetResourceDetailsFunc = func(key string) (string, *errors.Error) {
		return "", errors.PackError(0, "error while trying to connecting to DB: ")
	}
	resp := GetStorageResources(storageURI)
	assert.NotNil(t, resp, "There should be an error ")
	GetResourceDetailsFunc = func(key string) (string, *errors.Error) {
		return string([]byte(`{"user":"name"}`)), nil
	}
	resp = GetStorageResources(storageURI)
	assert.NotNil(t, resp, "There should be no error ")
}

func stubDevicePassword(password []byte) ([]byte, error) {
	if bytes.Compare(password, []byte("passwordWithInvalidEncryption")) == 0 {
		return []byte{}, fmt.Errorf("password decryption failed")
	}
	return password, nil
}

func TestPluginHealthCheckInterface_GetPluginStatus(t *testing.T) {
	type args struct {
		plugin agmodel.Plugin
	}
	PluginHealthCheck := &PluginHealthCheckInterface{
		DecryptPassword: common.DecryptWithPrivateKey,
	}
	password, _ := stubDevicePassword([]byte("password"))
	pluginData := agmodel.Plugin{
		IP:                "duphost",
		Port:              "9091",
		Username:          "admin",
		Password:          password,
		PreferredAuthType: "BasicAuth",
		ManagerUUID:       "mgr-addr",
	}
	var p = []string{}
	tests := []struct {
		name  string
		phc   *PluginHealthCheckInterface
		args  args
		want  bool
		want1 []string
	}{
		{
			name: "test1",
			phc:  PluginHealthCheck,
			args: args{
				plugin: pluginData,
			},
			want:  false,
			want1: p,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := tt.phc.GetPluginStatus(tt.args.plugin)
			if got != tt.want {
				t.Errorf("PluginHealthCheckInterface.GetPluginStatus() got = %v, want %v", got, tt.want)
			}

		})
	}
}

func TestLookupHost(t *testing.T) {
	config.SetUpMockConfig(t)

	ip, _, _, _ := LookupHost("10.0.0.0")
	assert.Equal(t, "10.0.0.0", ip, "Ip should be same")

	LookupIPfunc = func(host string) (ip []net.IP, err error) {
		err = fmt.Errorf("error")
		return
	}
	ip, _, _, err := LookupHost("10.0.0")
	assert.NotNil(t, err, "There should be an error")
	LookupIPfunc = func(host string) (ip []net.IP, err error) {
		return
	}
	ip, _, _, err = LookupHost("10.0.0.1")
	assert.NotNil(t, err, "There should be an error")

}
func mockGetPluginStatus(plugin agmodel.Plugin) bool {
	return true
}

func TestGetAllPlugins(t *testing.T) {
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
	mockPlugins(t)

	plugins, err := GetAllPlugins()
	assert.Nil(t, err, "Error Should be nil")
	assert.Equal(t, 3, len(plugins), "should be only 3 plugins")
}

func getEncryptedKey(t *testing.T, key []byte) []byte {
	cryptedKey, err := common.EncryptWithPublicKey(key)
	if err != nil {
		t.Fatalf("error: failed to encrypt data: %v", err)
	}
	return cryptedKey
}

func mockPlugins(t *testing.T) {
	connPool, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		t.Errorf("error while trying to connecting to DB: %v", err.Error())
	}

	password := getEncryptedKey(t, []byte("Password"))
	pluginArr := []agmodel.Plugin{
		{
			IP:                "localhost",
			Port:              "1234",
			Password:          password,
			Username:          "admin",
			ID:                "GRF",
			PreferredAuthType: "BasicAuth",
			PluginType:        "GRF",
		},
		{
			IP:                "localhost",
			Port:              "1234",
			Password:          password,
			Username:          "admin",
			ID:                "ILO",
			PreferredAuthType: "XAuthToken",
			PluginType:        "ILO",
		},
		{
			IP:                "localhost",
			Port:              "1234",
			Password:          password,
			Username:          "admin",
			ID:                "CFM",
			PreferredAuthType: "XAuthToken",
			PluginType:        "CFM",
		},
	}
	for _, plugin := range pluginArr {
		pl := "Plugin"
		//Save data into Database
		if err := connPool.Create(pl, plugin.ID, &plugin); err != nil {
			t.Fatalf("error: %v", err)
		}
	}
}
