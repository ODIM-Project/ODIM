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

package smodel

import (
	"context"
	"fmt"
	"testing"

	"github.com/ODIM-Project/ODIM/lib-persistence-manager/persistencemgr"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	"github.com/stretchr/testify/assert"
)

func mockFabricData(t *testing.T, table, id string, data interface{}) error {
	connPool, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		return fmt.Errorf("error while trying to connecting to DB: %v", err)
	}
	if err = connPool.Create(table, id, data); err != nil {
		return fmt.Errorf("error: mockData() failed to create entry %s-%s: %v", table, id, err)
	}
	return nil
}

func mockPlugin(t *testing.T, pluginID, PreferredAuthType, port string) error {
	password := getEncryptedKey(t, []byte("$2a$10$OgSUYvuYdI/7dLL5KkYNp.RCXISefftdj.MjbBTr6vWyNwAvht6ci"))
	plugin := Plugin{
		IP:                "10.10.10.10",
		Port:              port,
		Username:          "admin",
		Password:          password,
		ID:                pluginID,
		PluginType:        "Fabric",
		PreferredAuthType: PreferredAuthType,
	}
	connPool, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		return fmt.Errorf("error while trying to connecting to DB: %v", err.Error())
	}
	if err = connPool.Create("Plugin", pluginID, plugin); err != nil {
		return fmt.Errorf("error while trying to create new %v resource: %v", "Plugin", err.Error())
	}
	return nil
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

func TestGetFabricManagers(t *testing.T) {
	config.SetUpMockConfig(t)
	ctx := mockContext()
	defer func() {
		common.TruncateDB(common.OnDisk)
	}()
	err := mockPlugin(t, "GRF", "XAuthToken", "9091")
	if err != nil {
		t.Fatalf("Error in creating mock Plugin Data :%v", err)
	}
	f := map[string]interface{}{"FabricUUID": "794ef789-f54b-460c-8f30-03779d8403bc", "PluginID": "AFC_v6.0.1"}
	mockFabricData(t, "Fabric", "794ef789-f54b-460c-8f30-03779d8403bc", f)

	_, err = GetFabricManagers(ctx)
	assert.Nil(t, err, "should be no error ")

	GetDBConnectionFunc = func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
		return nil, &errors.Error{}
	}
	_, err = GetFabricManagers(ctx)
	assert.NotNil(t, err, "should be an error ")

	GetDBConnectionFunc = func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
		return common.GetDBConnection(dbFlag)
	}
	GetPluginDataFunc = func(pluginID string) (Plugin, *errors.Error) {
		return Plugin{}, nil
	}
	_, err = GetFabricManagers(ctx)
	assert.Nil(t, err, "should be no error ")

}
