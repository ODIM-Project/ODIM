/* (C) Copyright [2022] Hewlett Packard Enterprise Development LP
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may
 * not use this file except in compliance with the License. You may obtain
 * a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
 * License for the specific language governing permissions and limitations
 * under the License.
 */

// Package system ...

package system

import (
	"testing"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/svc-aggregation/agcommon"
	"github.com/ODIM-Project/ODIM/svc-aggregation/agmodel"
	"github.com/stretchr/testify/assert"
)

func Test_checkPluginStatus(t *testing.T) {
	phc := &agcommon.PluginHealthCheckInterface{
		DecryptPassword: common.DecryptWithPrivateKey,
	}
	password, _ := stubDevicePassword([]byte("password"))
	plugindata := agmodel.Plugin{
		IP:                "duphost",
		Port:              "9091",
		Username:          "admin",
		Password:          password,
		PreferredAuthType: "BasicAuth",
		ManagerUUID:       "mgr-addr",
	}

	type args struct {
		phc    *agcommon.PluginHealthCheckInterface
		plugin agmodel.Plugin
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test",
			args: args{
				phc:    phc,
				plugin: plugindata,
			},
		},
	}
	ctx := mockContext()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checkPluginStatus(ctx, tt.args.phc, tt.args.plugin)
		})
	}
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
	}
	for _, plugin := range pluginArr {
		pl := "Plugin"
		//Save data into Database
		if err := connPool.Create(pl, plugin.ID, &plugin); err != nil {
			t.Fatalf("error: %v", err)
		}
	}
}

func TestPushPluginStartUpData(t *testing.T) {
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
	startUpData := &agmodel.PluginStartUpData{
		RequestType:           "full",
		ResyncEvtSubscription: true,
	}
	plugin := agmodel.Plugin{
		ID: "10.0.0.0",
	}
	startUpData1 := &agmodel.PluginStartUpData{}
	ctx := mockContext()
	PushPluginStartUpData(ctx, agmodel.Plugin{}, startUpData)
	PushPluginStartUpData(ctx, agmodel.Plugin{}, startUpData1)

	err := PushPluginStartUpData(ctx, plugin, startUpData)
	assert.NotNil(t, err, "There should be no error")

}

func Test_sendPluginStartupRequest(t *testing.T) {
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
	startUpData := agmodel.PluginStartUpData{
		RequestType:           "full",
		ResyncEvtSubscription: true,
	}
	ctx := mockContext()
	var startUpData1 interface{}
	_, err := sendPluginStartupRequest(ctx, agmodel.Plugin{}, startUpData1, "")
	assert.NotNil(t, err, "There should be error")
	_, err = sendPluginStartupRequest(ctx, agmodel.Plugin{}, startUpData1, "ILO_v1.0.0")
	assert.NotNil(t, err, "There should be error")
	_, err = sendPluginStartupRequest(ctx, agmodel.Plugin{}, startUpData, "ILO_v1.0.0")
	assert.NotNil(t, err, "There should be error")

}
func Test_sendFullPluginInventory(t *testing.T) {
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
	err := sendFullPluginInventory(ctx, "", agmodel.Plugin{})
	assert.Nil(t, err, "There should be no error")
	plugin := agmodel.Plugin{
		ID: "localhost",
	}

	mockPlugins(t)
	err = sendFullPluginInventory(ctx, "", plugin)
	assert.Nil(t, err, "There should be no error")

	err = sendFullPluginInventory(ctx, "10.0.0.0", plugin)
	assert.Nil(t, err, "There should be no error")

}

func Test_sharePluginInventory(t *testing.T) {
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
	ctx := mockContext()
	sharePluginInventory(ctx, agmodel.Plugin{}, false, "")

	sharePluginInventory(ctx, agmodel.Plugin{}, false, "ILO_v1.0.0")
	sharePluginInventory(ctx, agmodel.Plugin{}, true, "ILO_v1.0.0")
}

func TestSendPluginStartUpData(t *testing.T) {
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
	password := getEncryptedKey(t, []byte("Password"))

	plugin := agmodel.Plugin{
		IP:                "localhost",
		Port:              "1234",
		Password:          password,
		Username:          "admin",
		ID:                "GRF",
		PreferredAuthType: "BasicAuth",
		PluginType:        "GRF",
	}

	mockPlugins(t)
	ctx := mockContext()
	err := SendPluginStartUpData(ctx, "", agmodel.Plugin{})
	assert.Nil(t, err, "There should be no error")
	err = SendPluginStartUpData(ctx, "", plugin)
	assert.Nil(t, err, "There should be no error")

}
