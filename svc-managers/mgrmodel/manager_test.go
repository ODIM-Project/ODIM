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
	"encoding/json"
	"testing"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/stretchr/testify/assert"
)

func TestGenericSave(t *testing.T) {
	common.SetUpMockConfig()
	defer func() {
		err := common.TruncateDB(common.InMemory)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()

	body := []byte(`body`)
	table := "EthernetInterfaces"
	key := "/redfish/v1/Managers/uuid:1/EthernetInterfaces/1"
	err := GenericSave(body, table, key)
	assert.Nil(t, err, "There should be no error")

	data, err := GetResource(table, key)
	assert.Nil(t, err, "There should be no error")
	assert.Equal(t, data, string(body), "should be same")
}

func TestManager_Update(t *testing.T) {
	common.SetUpMockConfig()
	defer func() {
		err := common.TruncateDB(common.InMemory)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()

	body := []byte(`{"Status":{"State":"Enabled"}}`)
	table := "Managers"
	key := "xyz"
	err := GenericSave(body, table, key)
	assert.Nil(t, err, "There should be no error while saving data")

	m := map[string]interface{}{
		"Status": map[string]string{
			"State": "Absent",
		},
	}
	err = UpdateData(key, m, "Managers")

	data, err := GetResource(table, key)
	assert.Nil(t, err, "There should be no error getting data")
	var mgr map[string]interface{}

	err = json.Unmarshal([]byte(data), &mgr)
	status := mgr["Status"].(map[string]interface{})
	state := status["State"].(string)
	assert.Equal(t, state, "Absent", "should be same")
}

func TestGetResourceNegativeTestCases(t *testing.T) {
	defer func() {
		err := common.TruncateDB(common.InMemory)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()
	// without db configuration
	table := "EthernetInterfaces"
	key := "/redfish/v1/Managers/uuid:1/EthernetInterfaces/1"

	_, err := GetResource(table, key)
	assert.NotNil(t, err, "There should be an error")

	// if key not present
	common.SetUpMockConfig()
	table = "Ethernet"
	key = "/redfish/v1/Managers/uuid:1/Ethernets/1"

	_, err = GetResource(table, key)
	assert.NotNil(t, err, "There should be an error")

}

func TestGetAllkeysFromTable(t *testing.T) {
	common.SetUpMockConfig()
	defer func() {
		err := common.TruncateDB(common.InMemory)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()

	body := []byte(`body`)
	table := "EthernetInterfaces"
	key := "/redfish/v1/Managers/uuid:1/EthernetInterfaces/1"
	err := GenericSave(body, table, key)
	assert.Nil(t, err, "There should be no error")

	allKeys, err := GetAllKeysFromTable(table)
	assert.Nil(t, err, "There should be no error")
	assert.Equal(t, len(allKeys), 1, "There should be one entry in DB")
}

func TestGetManagerByURL(t *testing.T) {
	common.SetUpMockConfig()
	defer func() {
		err := common.TruncateDB(common.InMemory)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()

	body := []byte(`body`)
	table := "Managers"
	key := "/redfish/v1/Managers/uuid:1"
	err := GenericSave(body, table, key)
	assert.Nil(t, err, "There should be no error")

	data, err := GetManagerByURL(key)
	assert.Nil(t, err, "There should be no error")
	assert.Equal(t, data, string(body), "should be same")
}

func TestGetManagerByURLNegativeTestCase(t *testing.T) {
	common.SetUpMockConfig()
	defer func() {
		err := common.TruncateDB(common.InMemory)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()
	// if key not present
	common.SetUpMockConfig()
	key := "/redfish/v1/Managers/uuid1:1"
	_, err := GetManagerByURL(key)
	assert.NotNil(t, err, "There should be an error")

}

func TestAddManagertoDB(t *testing.T) {
	common.SetUpMockConfig()
	defer func() {
		err := common.TruncateDB(common.InMemory)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()

	mngr := RAManager{
		Name:            "odimra",
		ManagerType:     "Service",
		FirmwareVersion: "1.0",
		ID:              "3bd1f589-117a-4cf9-89f2-da44ee8e012b",
		UUID:            "3bd1f589-117a-4cf9-89f2-da44ee8e012b",
		State:           "Enabled",
	}
	err := AddManagertoDB(mngr)
	assert.Nil(t, err, "There should be no error")

	data, err := GetManagerByURL("/redfish/v1/Managers/" + mngr.UUID)
	var manager RAManager
	json.Unmarshal([]byte(data), &manager)
	assert.Nil(t, err, "There should be no error")
	assert.Equal(t, manager.Name, "odimra", "Name should be odimra")
	assert.Equal(t, manager.FirmwareVersion, "1.0", "firmwareVersion should be 1.0")
	assert.Equal(t, manager.ManagerType, "Service", "managerType should be Service")
	assert.Equal(t, manager.ID, "3bd1f589-117a-4cf9-89f2-da44ee8e012b", "managerid should be 3bd1f589-117a-4cf9-89f2-da44ee8e012b")
	assert.Equal(t, manager.UUID, "3bd1f589-117a-4cf9-89f2-da44ee8e012b", "uuid should be 3bd1f589-117a-4cf9-89f2-da44ee8e012b")
	assert.Equal(t, manager.State, "Enabled", "state should be Enabled")
}
