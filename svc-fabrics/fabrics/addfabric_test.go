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

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"gotest.tools/assert"

	//	"github.com/ODIM-Project/ODIM/svc-fabrics/fabresponse"
	fabricsproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/fabrics"
	//"github.com/ODIM-Project/ODIM/svc-fabrics/fabmodel"
)

func TestAddFabric(t *testing.T) {
	config.SetUpMockConfig(t)
	defer func() {
		err := common.TruncateDB(common.OnDisk)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()
	err := mockPluginData(t, "CFM", "XAuthToken", "9091")
	if err != nil {
		t.Fatalf("Error in creating mock DeviceData :%v", err)
	}
	req := &fabricsproto.AddFabricRequest{
		OriginResource: "/redfish/v1/Fabrics/a926dec5-61eb-499b-988a-d45b45847466",
		Address:        "localhost",
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
	err := mockPluginData(t, "GRF", "XAuthToken", "9091")
	if err != nil {
		t.Fatalf("Error in creating mock Plugin Data :%v", err)
	}
	req := &fabricsproto.AddFabricRequest{
		OriginResource: "/redfish/v1/Fabrics/a926dec5-61eb-499b-988a-d45b45847466",
		Address:        "localhost",
	}

	resp := AddFabric(req)
	assert.Equal(t, int(resp.StatusCode), http.StatusOK, "should be same")

	resp = AddFabric(req)
	assert.Equal(t, int(resp.StatusCode), http.StatusInternalServerError, "should be same")
}
