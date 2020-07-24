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

//Package fabrics ...
package fabrics

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/bharath-b-hpe/odimra/lib-utilities/common"
	fabricsproto "github.com/bharath-b-hpe/odimra/lib-utilities/proto/fabrics"
	"github.com/bharath-b-hpe/odimra/svc-fabrics/fabmodel"
	"net/http"
	"testing"
)

func mockFabricData(fabricID, pluginID string) error {

	connPool, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		return fmt.Errorf("error while trying to connecting to DB: %v", err.Error())
	}
	fab := fabmodel.Fabric{
		FabricUUID: fabricID,
		PluginID:   pluginID,
	}
	if err = connPool.Create("Fabric", fabricID, fab); err != nil {
		return fmt.Errorf("error while trying to create new %v resource: %v", "fabric", err.Error())
	}
	return nil
}

func TestFabrics_WithInvalidPluginData(t *testing.T) {
	Token.Tokens = make(map[string]string)
	Token.Tokens["GRF"] = "234556"
	common.SetUpMockConfig()
	defer func() {
		err := common.TruncateDB(common.OnDisk)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()
	postData, _ := json.Marshal(map[string]interface{}{
		"@odata.id": "/redfish/v1/Fabrics",
	})
	err := mockPluginData(t, "GRF", "XAuthToken", "9091")
	if err != nil {
		t.Fatalf("Error in creating mock DeviceData :%v", err)
	}
	err = mockFabricData("fabid1", "CFM")
	if err != nil {
		t.Fatalf("Error in creating mockFabricData :%v", err)
	}

	var f = &Fabrics{
		Auth:          mockAuth,
		ContactClient: mockContactClient,
	}
	req := &fabricsproto.FabricRequest{
		SessionToken: "valid",
		Method:       http.MethodPost,
		URL:          "/redfish/v1/Fabrics/fabid1/Zones/Zone1",
		RequestBody:  postData,
	}
	resp := f.UpdateFabricResource(req)
	assert.Equal(t, int(resp.StatusCode), http.StatusNotFound, "should be same")

	req = &fabricsproto.FabricRequest{
		Method:       http.MethodDelete,
		SessionToken: "valid",
		URL:          "/redfish/v1/Fabrics/fabid1/Zones/Zone1",
	}

	resp = f.DeleteFabricResource(req)
	assert.Equal(t, int(resp.StatusCode), http.StatusNotFound, "should be same")

	req = &fabricsproto.FabricRequest{
		Method:       http.MethodGet,
		SessionToken: "valid",
		URL:          "/redfish/v1/Fabrics/fabid1/Zones/Zone1",
	}
	resp = f.GetFabricResource(req)
	assert.Equal(t, int(resp.StatusCode), http.StatusNotFound, "should be same")
}

func TestFabrics_WithInvalidURI(t *testing.T) {
	Token.Tokens = make(map[string]string)
	Token.Tokens["CFM"] = "123456"
	common.SetUpMockConfig()
	defer func() {
		err := common.TruncateDB(common.OnDisk)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()
	postData, _ := json.Marshal(map[string]interface{}{
		"@odata.id": "/redfish/v1/Fabrics",
	})
	err := mockPluginData(t, "CFM", "XAuthToken", "9091")
	if err != nil {
		t.Fatalf("Error in creating mock DeviceData :%v", err)
	}
	err = mockFabricData("fabid2", "CFM")
	if err != nil {
		t.Fatalf("Error in creating mockFabricData :%v", err)
	}
	var f = &Fabrics{
		Auth:          mockAuth,
		ContactClient: mockContactClient,
	}
	req := &fabricsproto.FabricRequest{
		SessionToken: "valid",
		URL:          "/redfish/v1/Fabrics/fabid2/Zones/Zone1",
		RequestBody:  postData,
		Method:       http.MethodPost,
	}
	resp := f.UpdateFabricResource(req)
	assert.Equal(t, int(resp.StatusCode), http.StatusNotFound, "should be same")

	req = &fabricsproto.FabricRequest{
		SessionToken: "valid",
		URL:          "/redfish/v1/Fabrics/fabid2/Zones/Zone1",
		Method:       http.MethodDelete,
	}

	resp = f.DeleteFabricResource(req)
	assert.Equal(t, int(resp.StatusCode), http.StatusNotFound, "should be same")

	req = &fabricsproto.FabricRequest{
		SessionToken: "valid",
		URL:          "/redfish/v1/Fabrics/fabid2",
		Method:       http.MethodGet,
	}

	resp = f.GetFabricResource(req)
	assert.Equal(t, int(resp.StatusCode), http.StatusNotFound, "should be same")
}

func TestFabrics_WithInvaliPluginCredentials(t *testing.T) {
	Token.Tokens = make(map[string]string)
	Token.Tokens["CFM"] = ""
	common.SetUpMockConfig()
	defer func() {
		err := common.TruncateDB(common.OnDisk)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()
	postData, _ := json.Marshal(map[string]interface{}{
		"@odata.id": "/redfish/v1/Fabrics",
	})
	err := mockPluginData(t, "CFM", "XAuthToken", "9092")
	if err != nil {
		t.Fatalf("Error in creating mock DeviceData :%v", err)
	}
	err = mockFabricData("fabid2", "CFM")
	if err != nil {
		t.Fatalf("Error in creating mockFabricData :%v", err)
	}
	var f = &Fabrics{
		Auth:          mockAuth,
		ContactClient: mockContactClient,
	}
	req := &fabricsproto.FabricRequest{
		SessionToken: "valid",
		URL:          "/redfish/v1/Fabrics/fabid2/Zones/Zone1",
		RequestBody:  postData,
		Method:       http.MethodPost,
	}
	resp := f.UpdateFabricResource(req)
	assert.Equal(t, int(resp.StatusCode), http.StatusUnauthorized, "should be same")

	req = &fabricsproto.FabricRequest{
		SessionToken: "valid",
		URL:          "/redfish/v1/Fabrics/fabid2/Zones/Zone1",
		Method:       http.MethodDelete,
	}

	resp = f.DeleteFabricResource(req)
	assert.Equal(t, int(resp.StatusCode), http.StatusUnauthorized, "should be same")

	req = &fabricsproto.FabricRequest{
		SessionToken: "valid",
		URL:          "/redfish/v1/Fabrics/fabid2",
		Method:       http.MethodGet,
	}

	resp = f.GetFabricResource(req)
	assert.Equal(t, int(resp.StatusCode), http.StatusUnauthorized, "should be same")
}

func TestFabrics_WithBasicAuth(t *testing.T) {
	Token.Tokens = make(map[string]string)

	common.SetUpMockConfig()
	defer func() {
		err := common.TruncateDB(common.OnDisk)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()
	err := mockPluginData(t, "CFM", "BasicAuth", "9093")
	if err != nil {
		t.Fatalf("Error in creating mock DeviceData :%v", err)
	}

	err = mockFabricData("fabid2", "CFM")
	if err != nil {
		t.Fatalf("Error in creating mockFabricData :%v", err)
	}
	var f = &Fabrics{
		Auth:          mockAuth,
		ContactClient: mockContactClient,
	}

	req := &fabricsproto.FabricRequest{
		SessionToken: "valid",
		URL:          "/redfish/v1/Fabrics/fabid2",
		Method:       http.MethodGet,
	}

	resp := f.GetFabricResource(req)
	assert.Equal(t, int(resp.StatusCode), http.StatusOK, "should be same")
}

func TestFabrics_WithInvalidData(t *testing.T) {
	Token.Tokens = make(map[string]string)
	Token.Tokens["CFM"] = "234556"

	common.SetUpMockConfig()
	defer func() {
		err := common.TruncateDB(common.OnDisk)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()
	err := mockPluginData(t, "CFM", "XAuthToken", "9095")
	if err != nil {
		t.Fatalf("Error in creating mock DeviceData :%v", err)
	}
	err = mockFabricData("d72dade0-c35a-984c-4859-1108132d72da", "CFM")
	if err != nil {
		t.Fatalf("Error in creating mockFabricData :%v", err)
	}
	postData, _ := json.Marshal(map[string]interface{}{
		"@odata.id": "/redfish/v1/Fabrics",
	})
	var f = &Fabrics{
		Auth:          mockAuth,
		ContactClient: mockContactClient,
	}

	req := &fabricsproto.FabricRequest{
		SessionToken: "valid",
		URL:          "/redfish/v1/Fabrics/d72dade0-c35a-984c-4859-1108132d72da/Zones/Zone1",
		RequestBody:  postData,
		Method:       "POST",
	}

	resp := f.UpdateFabricResource(req)
	assert.Equal(t, int(resp.StatusCode), http.StatusInternalServerError, "should be same")
}

func TestGetFabricID(t *testing.T) {
	url := "/redfish/v1/Fabrics/d72dade0-c35a-984c-4859-1108132d72da"
	fabID := getFabricID(url)
	assert.Equal(t, "d72dade0-c35a-984c-4859-1108132d72da", fabID, "fabric id should be d72dade0-c35a-984c-4859-1108132d72da")

	url = "/redfish/v1/Fabrics/d72dade0-c35a-984c-4859-1108132d72da/Zones/Zone1"
	fabID = getFabricID(url)
	assert.Equal(t, "d72dade0-c35a-984c-4859-1108132d72da", fabID, "fabric id should be d72dade0-c35a-984c-4859-1108132d72da")
}
