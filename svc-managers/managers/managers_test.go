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
package managers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	managersproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/managers"
	"github.com/ODIM-Project/ODIM/svc-managers/mgrcommon"
	"github.com/ODIM-Project/ODIM/svc-managers/mgrmodel"
	"github.com/ODIM-Project/ODIM/svc-managers/mgrresponse"
	"github.com/stretchr/testify/assert"
)

func TestGetManagersCollection(t *testing.T) {
	config.SetUpMockConfig(t)
	common.TruncateDB(common.InMemory)
	defer func() {
		err := common.TruncateDB(common.InMemory)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()
	body := []byte(`body`)
	table := "Managers"
	key := "/redfish/v1/Managers/uuid:1"
	mgrmodel.GenericSave(body, table, key)

	req := &managersproto.ManagerRequest{}
	response, err := GetManagersCollection(req)
	assert.Nil(t, err, "There should be no error")

	manager := response.Body.(mgrresponse.ManagersCollection)
	assert.Equal(t, int(response.StatusCode), http.StatusOK, "Status code should be StatusOK.")
	assert.Equal(t, manager.MembersCount, 2, "Status code should be StatusOK.")
}

func mockGetDeviceInfo(req mgrcommon.ResourceInfoRequest) (string, error) {
	if req.URL == "/redfish/v1/Managers/deviceAbsent:1" || req.URL == "/redfish/v1/Managers/uuid1:1/Ethernet" {
		return "", fmt.Errorf("error")
	}
	manager := mgrmodel.Manager{
		Status: &mgrmodel.Status{
			State: "Enabled",
		},
	}
	dataByte, err := json.Marshal(manager)
	return string(dataByte), err
}

func mockContactClient(url, method, token string, odataID string, body interface{}, loginCredential map[string]string) (*http.Response, error) {

	if url == "https://localhost:9091/ODIM/v1/Sessions" {
		body := `{"Token": "12345"}`
		return &http.Response{
			StatusCode: http.StatusCreated,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
			Header: http.Header{
				"X-Auth-Token": []string{"12345"},
			},
		}, nil
	} else if url == "https://localhost:9092/ODIM/v1/Sessions" {
		body := `{"Token": ""}`
		return &http.Response{
			StatusCode: http.StatusUnauthorized,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil
	}
	if url == "https://localhost:9091/ODIM/v1/Managers/uuid/EthernetInterfaces" && token == "12345" {
		body := `{"data": "/ODIM/v1/Managers/uuid/EthernetInterfaces"}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil
	} else if url == "https://localhost:9093/ODIM/v1/Managers/uuid1/EthernetInterfaces" {
		body := `{"data": "/ODIM/v1/Managers/uuid/EthernetInterfaces"}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil
	} else if url == "https://localhost:9092/ODIM/v1/Managers/uuid/EthernetInterfaces" && token == "23456" {
		body := `{"data": "/ODIM/v1/Managers/uuid/EthernetInterfaces"}`
		return &http.Response{
			StatusCode: http.StatusUnauthorized,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil
	}
	return nil, fmt.Errorf("InvalidRequest")
}

func getEncryptedKey(t *testing.T, key []byte) []byte {
	cryptedKey, err := common.EncryptWithPublicKey(key)
	if err != nil {
		t.Fatalf("error: failed to encrypt data: %v", err)
	}
	return cryptedKey
}

func mockPluginData(t *testing.T, pluginID, PreferredAuthType, port string) error {
	password := getEncryptedKey(t, []byte("$2a$10$OgSUYvuYdI/7dLL5KkYNp.RCXISefftdj.MjbBTr6vWyNwAvht6ci"))
	plugin := mgrmodel.Plugin{
		IP:                "localhost",
		Port:              port,
		Username:          "admin",
		Password:          password,
		ID:                pluginID,
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

func TestGetManagerRootUUIDNotFound(t *testing.T) {
	config.SetUpMockConfig(t)
	defer func() {
		err := common.TruncateDB(common.InMemory)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()

	req := &managersproto.ManagerRequest{
		ManagerID: config.Data.RootServiceUUID,
	}
	d := DeviceContact{
		GetDeviceInfo: mockGetDeviceInfo,
		ContactClient: mockContactClient,
	}
	response := d.GetManagers(req)

	assert.Equal(t, int(response.StatusCode), http.StatusNotFound, "Status code should be StatusNotFound")
}

func TestGetManager(t *testing.T) {
	config.SetUpMockConfig(t)
	defer func() {
		err := common.TruncateDB(common.InMemory)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()
	mngr := mgrmodel.RAManager{
		Name:            "odimra",
		ManagerType:     "Service",
		FirmwareVersion: config.Data.FirmwareVersion,
		ID:              config.Data.RootServiceUUID,
		UUID:            config.Data.RootServiceUUID,
		State:           "Enabled",
	}
	mngr.AddManagertoDB()
	req := &managersproto.ManagerRequest{
		ManagerID: "3bd1f589-117a-4cf9-89f2-da44ee8e012b",
	}
	d := DeviceContact{
		GetDeviceInfo: mockGetDeviceInfo,
		ContactClient: mockContactClient,
	}
	response := d.GetManagers(req)

	var manager mgrmodel.Manager
	data, _ := json.Marshal(response.Body)
	json.Unmarshal(data, &manager)

	assert.Equal(t, int(response.StatusCode), http.StatusOK, "Status code should be StatusOK.")
	assert.Equal(t, manager.Name, "odimra", "Status code should be StatusOK.")
	assert.Equal(t, manager.ManagerType, "Service", "Status code should be StatusOK.")
	assert.Equal(t, manager.ID, req.ManagerID, "Status code should be StatusOK.")
	assert.Equal(t, manager.FirmwareVersion, "1.0", "Status code should be StatusOK.")

}

func TestGetManagerWithDeviceAbsent(t *testing.T) {
	config.SetUpMockConfig(t)
	defer func() {
		err := common.TruncateDB(common.InMemory)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()
	body := []byte(`{"ManagerType":"BMC","Status":{"State":"Enabled"}}`)
	table := "Managers"
	key := "/redfish/v1/Managers/deviceAbsent:1"
	mgrmodel.GenericSave(body, table, key)

	req := &managersproto.ManagerRequest{
		ManagerID: "uuid:1",
		URL:       "/redfish/v1/Managers/deviceAbsent:1",
	}
	d := DeviceContact{
		GetDeviceInfo: mockGetDeviceInfo,
		ContactClient: mockContactClient,
	}
	response := d.GetManagers(req)

	var manager mgrmodel.Manager
	data, _ := json.Marshal(response.Body)
	json.Unmarshal(data, &manager)

	assert.Equal(t, int(response.StatusCode), http.StatusOK, "Status code should be StatusOK.")
	assert.Equal(t, manager.Status.State, "Absent", "Status state should be Absent.")

}

func TestGetManagerwithInvalidURL(t *testing.T) {
	config.SetUpMockConfig(t)
	defer func() {
		err := common.TruncateDB(common.InMemory)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()

	body := []byte(`{"Status":{"State":"Enabled"}}`)
	table := "Managers"
	key := "/redfish/v1/Managers/uuid:1"
	mgrmodel.GenericSave(body, table, key)

	req := &managersproto.ManagerRequest{
		ManagerID: "uuid:1",
		URL:       "/redfish/v1/Managers/uuid1:1",
	}
	d := DeviceContact{
		GetDeviceInfo: mockGetDeviceInfo,
		ContactClient: mockContactClient,
	}
	response := d.GetManagers(req)
	assert.Equal(t, int(response.StatusCode), http.StatusNotFound, "Status code should be StatusOK.")

}

func TestGetManagerwithValidURL(t *testing.T) {
	config.SetUpMockConfig(t)
	defer func() {
		err := common.TruncateDB(common.InMemory)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()

	body := []byte(`{"ManagerType":"BMC","Status":{"State":"Enabled"}}`)
	table := "Managers"
	key := "/redfish/v1/Managers/uuid:1"
	mgrmodel.GenericSave(body, table, key)

	req := &managersproto.ManagerRequest{
		ManagerID: "uuid:1",
		URL:       "/redfish/v1/Managers/uuid:1",
	}
	d := DeviceContact{
		GetDeviceInfo: mockGetDeviceInfo,
		ContactClient: mockContactClient,
	}
	response := d.GetManagers(req)
	assert.Equal(t, http.StatusOK, int(response.StatusCode), "Status code should be StatusOK.")

}

func TestGetManagerInvalidID(t *testing.T) {
	config.SetUpMockConfig(t)
	defer func() {
		err := common.TruncateDB(common.InMemory)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()

	req := &managersproto.ManagerRequest{
		ManagerID: "invalidID",
	}
	d := DeviceContact{
		GetDeviceInfo: mockGetDeviceInfo,
		ContactClient: mockContactClient,
	}
	response := d.GetManagers(req)

	assert.Equal(t, int(response.StatusCode), http.StatusNotFound, "Status code should be StatusNotFound")
}

func TestGetManagerResourcewithBadManagerID(t *testing.T) {
	config.SetUpMockConfig(t)
	defer func() {
		err := common.TruncateDB(common.InMemory)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()

	req := &managersproto.ManagerRequest{
		ManagerID: "uuid",
		URL:       "/redfish/v1/Managers/uuid",
	}
	d := DeviceContact{
		GetDeviceInfo: mockGetDeviceInfo,
		ContactClient: mockContactClient,
	}
	response := d.GetManagersResource(req)
	assert.Equal(t, http.StatusNotFound, int(response.StatusCode), "Status code should be StatusBadRequest.")
}

func TestGetManagerResourcewithValidURL(t *testing.T) {
	config.SetUpMockConfig(t)
	defer func() {
		err := common.TruncateDB(common.InMemory)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()

	body := []byte(`body`)
	table := "Managers"
	key := "/redfish/v1/Managers/uuid:1"
	mgrmodel.GenericSave(body, table, key)

	body = []byte(`body`)
	table = "EthernetInterfacesCollection"
	key = "/redfish/v1/Managers/uuid:1/EthernetInterfaces"
	mgrmodel.GenericSave(body, table, key)

	body = []byte(`body`)
	table = "EthernetInterfaces"
	key = "/redfish/v1/Managers/uuid:1/EthernetInterfaces/1"
	mgrmodel.GenericSave(body, table, key)

	req := &managersproto.ManagerRequest{
		ManagerID: "uuid:1",
		URL:       "/redfish/v1/Managers/uuid:1/EthernetInterfaces",
	}
	d := DeviceContact{
		GetDeviceInfo: mockGetDeviceInfo,
		ContactClient: mockContactClient,
	}
	response := d.GetManagersResource(req)
	assert.Equal(t, http.StatusOK, int(response.StatusCode), "Status code should be StatusOK.")

	req = &managersproto.ManagerRequest{
		ManagerID:  "uuid:1",
		ResourceID: "1",
		URL:        "/redfish/v1/Managers/uuid:1/EthernetInterfaces/1",
	}
	d = DeviceContact{
		GetDeviceInfo: mockGetDeviceInfo,
		ContactClient: mockContactClient,
	}
	response = d.GetManagersResource(req)
	assert.Equal(t, http.StatusOK, int(response.StatusCode), "Status code should be StatusOK.")

}

func TestGetManagerResourcewithInvalidURL(t *testing.T) {
	config.SetUpMockConfig(t)
	defer func() {
		err := common.TruncateDB(common.InMemory)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()
	body := []byte(`body`)
	table := "EthernetInterfacesCollection"
	key := "/redfish/v1/Managers/uuid:1/EthernetInterfaces"
	mgrmodel.GenericSave(body, table, key)

	req := &managersproto.ManagerRequest{
		ManagerID: "uuid1:1",
		URL:       "/redfish/v1/Managers/uuid1:1/Ethernet",
	}
	d := DeviceContact{
		GetDeviceInfo: mockGetDeviceInfo,
		ContactClient: mockContactClient,
	}
	response := d.GetManagersResource(req)
	assert.Equal(t, http.StatusNotFound, int(response.StatusCode), "Status code should be StatusNotFound.")
}

func TestGetPluginManagerResource(t *testing.T) {
	mgrcommon.Token.Tokens = make(map[string]string)

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
	err := mockPluginData(t, "CFM", "XAuthToken", "9091")
	if err != nil {
		t.Fatalf("Error in creating mock DeviceData :%v", err)
	}
	err = mockPluginData(t, "GRF", "BasicAuth", "9093")
	if err != nil {
		t.Fatalf("Error in creating mock DeviceData :%v", err)
	}
	body, _ := json.Marshal(map[string]string{
		"Name": "CFM",
	})
	table := "Managers"
	key := "/redfish/v1/Managers/uuid"
	mgrmodel.GenericSave(body, table, key)

	body, _ = json.Marshal(map[string]string{
		"Name": "GRF",
	})
	table = "Managers"
	key = "/redfish/v1/Managers/uuid1"
	mgrmodel.GenericSave(body, table, key)

	req := &managersproto.ManagerRequest{
		ManagerID: "uuid",
		URL:       "/redfish/v1/Managers/uuid/EthernetInterfaces",
	}
	d := DeviceContact{
		GetDeviceInfo: mockGetDeviceInfo,
		ContactClient: mockContactClient,
	}
	response := d.GetManagersResource(req)
	assert.Equal(t, http.StatusOK, int(response.StatusCode), "Status code should be StatusOK.")

	req = &managersproto.ManagerRequest{
		ManagerID:  "uuid1",
		ResourceID: "1",
		URL:        "/redfish/v1/Managers/uuid1/EthernetInterfaces",
	}
	d = DeviceContact{
		GetDeviceInfo: mockGetDeviceInfo,
		ContactClient: mockContactClient,
	}
	response = d.GetManagersResource(req)
	assert.Equal(t, http.StatusOK, int(response.StatusCode), "Status code should be StatusOK.")

}

func TestGetPluginManagerResourceInvalidPlugin(t *testing.T) {
	mgrcommon.Token.Tokens = make(map[string]string)

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
	err := mockPluginData(t, "CFM1", "XAuthToken", "9091")
	if err != nil {
		t.Fatalf("Error in creating mock DeviceData :%v", err)
	}
	body, _ := json.Marshal(map[string]string{
		"Name": "CFM",
	})
	table := "Managers"
	key := "/redfish/v1/Managers/uuid"
	mgrmodel.GenericSave(body, table, key)

	req := &managersproto.ManagerRequest{
		ManagerID: "uuid",
		URL:       "/redfish/v1/Managers/uuid/EthernetInterfaces",
	}
	d := DeviceContact{
		GetDeviceInfo: mockGetDeviceInfo,
		ContactClient: mockContactClient,
	}
	response := d.GetManagersResource(req)
	assert.Equal(t, http.StatusNotFound, int(response.StatusCode), "Status code should be StatusOK.")
}

func TestGetPluginManagerResourceInvalidPluginSessions(t *testing.T) {
	mgrcommon.Token.Tokens = make(map[string]string)

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
	err := mockPluginData(t, "CFM", "XAuthToken", "9092")
	if err != nil {
		t.Fatalf("Error in creating mock DeviceData :%v", err)
	}
	body, _ := json.Marshal(map[string]string{
		"Name": "CFM",
	})
	table := "Managers"
	key := "/redfish/v1/Managers/uuid"
	mgrmodel.GenericSave(body, table, key)

	req := &managersproto.ManagerRequest{
		ManagerID: "uuid",
		URL:       "/redfish/v1/Managers/uuid/EthernetInterfaces",
	}
	d := DeviceContact{
		GetDeviceInfo: mockGetDeviceInfo,
		ContactClient: mockContactClient,
	}
	response := d.GetManagersResource(req)
	assert.Equal(t, http.StatusUnauthorized, int(response.StatusCode), "Status code should be StatusOK.")
	mgrcommon.Token.Tokens = map[string]string{
		"CFM": "23456",
	}
	response = d.GetManagersResource(req)
	assert.Equal(t, http.StatusUnauthorized, int(response.StatusCode), "Status code should be StatusOK.")

}