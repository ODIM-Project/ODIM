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

package mgrmodel

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/stretchr/testify/assert"
)

var plugin = Plugin{
	IP:                "localhost",
	Port:              "9091",
	Username:          "admin",
	ID:                "ILO",
	PreferredAuthType: "XAuthToken",
}
var target = DeviceTarget{
	ManagerAddress: "10.24.0.14",
	Password:       []byte("Password"),
	UserName:       "admin",
	DeviceUUID:     "uuid",
	PluginID:       "GRF",
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

func getEncryptedKey(t *testing.T, key []byte) []byte {
	cryptedKey, err := common.EncryptWithPublicKey(key)
	if err != nil {
		t.Fatalf("error: failed to encrypt data: %v", err)
	}
	return cryptedKey
}

func mockPluginData(t *testing.T) error {
	plugin.Password = getEncryptedKey(t, []byte("12345"))
	connPool, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		return fmt.Errorf("error while trying to connecting to DB: %v", err.Error())
	}
	if err = connPool.Create("Plugin", "ILO", plugin); err != nil {
		return fmt.Errorf("error while trying to create new %v resource: %v", "Plugin", err.Error())
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

func TestGetTarget(t *testing.T) {
	config.SetUpMockConfig(t)
	defer func() {
		err := common.TruncateDB(common.OnDisk)
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
func TestGetSystemByUUID(t *testing.T) {
	body := `{"Id":"1","Status":{"State":"Enabled"}}`
	table := "ComputerSystem"
	key := "/redfish/v1/Systems/uuid.1"
	GenericSave([]byte(body), table, key)
	data, _ := GetSystemByUUID("/redfish/v1/Systems/uuid.1")
	assert.Equal(t, data, body, "should be same")

	_, err := GetSystemByUUID("/redfish/v1/Systems/uuid")
	assert.NotNil(t, err, "There should be an error")
}
