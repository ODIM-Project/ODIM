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

//Package fabrics ...
package fabrics

import (
	"net/http"
	"testing"

	"fmt"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	fabricsproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/fabrics"
	"github.com/ODIM-Project/ODIM/svc-fabrics/fabmodel"
	"gotest.tools/assert"
)

func TestAddFabric(t *testing.T) {
	config.SetUpMockConfig(t)
	defer func() {
		err := common.TruncateDB(common.OnDisk)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()
	err := mockPlugin(t, "CFM", "XAuthToken", "9091")
	if err != nil {
		t.Fatalf("Error in creating mock DeviceData :%v", err)
	}
	req := &fabricsproto.AddFabricRequest{
		OriginResource: "/redfish/v1/Fabrics/a926dec5-61eb-499b-988a-d45b45847466",
		Address:        "10.10.10.10",
	}

	resp := AddFabric(req)
	assert.Equal(t, int(resp.StatusCode), http.StatusOK, "should be same")

	resp = AddFabric(req)
	assert.Equal(t, int(resp.StatusCode), http.StatusInternalServerError, "should be same")
}

func TestAddFabricInvalidPluginID(t *testing.T) {
	config.SetUpMockConfig(t)
	defer func() {
		err := common.TruncateDB(common.OnDisk)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()
	err := mockPlugin(t, "GRF", "XAuthToken", "9091")
	if err != nil {
		t.Fatalf("Error in creating mock Plugin Data :%v", err)
	}
	req := &fabricsproto.AddFabricRequest{
		OriginResource: "/redfish/v1/Fabrics/a926dec5-61eb-499b-988a-d45b45847466",
		Address:        "10.10.10.10",
	}

	resp := AddFabric(req)
	assert.Equal(t, int(resp.StatusCode), http.StatusOK, "should be same")

	resp = AddFabric(req)
	assert.Equal(t, int(resp.StatusCode), http.StatusInternalServerError, "should be same")
}

func mockPlugin(t *testing.T, pluginID, PreferredAuthType, port string) error {
	password := getEncryptedKey(t, []byte("$2a$10$OgSUYvuYdI/7dLL5KkYNp.RCXISefftdj.MjbBTr6vWyNwAvht6ci"))
	plugin := fabmodel.Plugin{
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
